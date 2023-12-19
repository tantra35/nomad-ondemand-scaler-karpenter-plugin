package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	awsclient "github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/karpenter-core/pkg/apis/v1alpha5"
	"github.com/aws/karpenter/pkg/apis/settings"
	"github.com/aws/karpenter/pkg/apis/v1alpha1"
	awscache "github.com/aws/karpenter/pkg/cache"
	"github.com/aws/karpenter/pkg/cloudprovider"
	"github.com/aws/karpenter/pkg/providers/amifamily"
	"github.com/aws/karpenter/pkg/providers/instance"
	"github.com/aws/karpenter/pkg/providers/instancetype"
	"github.com/aws/karpenter/pkg/providers/launchtemplate"
	"github.com/aws/karpenter/pkg/providers/pricing"
	"github.com/aws/karpenter/pkg/providers/securitygroup"
	"github.com/aws/karpenter/pkg/providers/subnet"
	"github.com/aws/karpenter/pkg/utils"
	"github.com/hashicorp/go-plugin"
	"github.com/patrickmn/go-cache"
	"google.golang.org/grpc"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"playrix.com/it/nomad-cluster-scalerv2/karpenterprovidergrpc/karpenterprovidergrpc"
)

const (
	RESCOURCERESERVATIONPERSENT = 0.025
)

type server struct {
	karpenterprovidergrpc.UnimplementedKarpenterServiceServer
	clusterName            string
	clusterEndpoint        string
	fakek8sAPI             *fakek8sClient
	karpenderCloudProvider *cloudprovider.CloudProvider
}

func (s *server) ListInstances(_c context.Context, _r *karpenterprovidergrpc.ListInstancesRequest) (*karpenterprovidergrpc.ListInstancesResponse, error) {
	lsettings := &settings.Settings{ //заполняется кластером k8s, По имени клатера нужно огранизовывать пулы
		ClusterName:                s.clusterName + "/" + _r.PoolName,
		ClusterEndpoint:            s.clusterEndpoint,
		DefaultInstanceProfile:     "",
		EnablePodENI:               false,
		EnableENILimitedPodDensity: true,
		IsolatedVPC:                false,
		VMMemoryOverheadPercent:    RESCOURCERESERVATIONPERSENT,
		InterruptionQueueName:      "",
		Tags:                       map[string]string{},
		ReservedENIs:               0,
	}
	execCtx := context.WithValue(_c, settings.ContextKey, lsettings)

	//---------------------------------------------------------------------------
	nodetemplate := &v1alpha1.AWSNodeTemplate{}
	execCtx = context.WithValue(execCtx, fakek8sClientTemplateKey, nodetemplate)

	//---------------------------------------------------------------------------
	lmachines, lerr := s.karpenderCloudProvider.List(execCtx)
	if lerr != nil {
		return nil, lerr
	}

	linstanseIds := make([]string, 0, len(lmachines))
	for _, lmachine := range lmachines {
		linstanseId, lerr := utils.ParseInstanceID(lmachine.Status.ProviderID)
		if lerr != nil {
			continue
		}

		linstanseIds = append(linstanseIds, linstanseId)
	}

	return &karpenterprovidergrpc.ListInstancesResponse{
		Instanseids: linstanseIds,
	}, nil
}

func (s *server) AddInstances(_c context.Context, _a *karpenterprovidergrpc.AddInstancesRequest) (*karpenterprovidergrpc.AddInstancesResponse, error) {
	nodeTemplate := &v1alpha1.AWSNodeTemplate{
		ObjectMeta: metav1.ObjectMeta{
			UID: "123",
		},

		Spec: v1alpha1.AWSNodeTemplateSpec{
			AWS: v1alpha1.AWS{
				InstanceProfile: aws.String(_a.Spec.InstanceProfile),
				SubnetSelector:  _a.Spec.Subnets,

				SecurityGroupSelector: _a.Spec.SecurityGroups,

				LaunchTemplate: v1alpha1.LaunchTemplate{
					LaunchTemplateName: _a.Spec.LaunchTemplate,
				},
			},

			AMISelector: _a.Spec.Ami,
		},
	}
	execCtx := context.WithValue(_c, fakek8sClientTemplateKey, nodeTemplate)

	linstanceSpec := _a.Spec
	lexecsettings := &settings.Settings{ //заполняется кластером k8s, не то чтобы очень полезная инфа
		ClusterName:                s.clusterName + "/" + _a.PoolName,
		ClusterEndpoint:            s.clusterEndpoint,
		DefaultInstanceProfile:     "",
		EnablePodENI:               false,
		EnableENILimitedPodDensity: true,
		IsolatedVPC:                false,
		VMMemoryOverheadPercent:    RESCOURCERESERVATIONPERSENT,
		InterruptionQueueName:      "",
		Tags:                       map[string]string{},
		ReservedENIs:               0,
	}

	linstances := make([]string, 0, _a.Count)

	for i := 0; i < int(_a.Count); i++ {
		lrequirements := make([]v1.NodeSelectorRequirement, 0, len(linstanceSpec.GetRequirements()))
		for _, lreq := range linstanceSpec.GetRequirements() {
			lrequirements = append(lrequirements, v1.NodeSelectorRequirement{
				Key:      lreq.Key,
				Operator: v1.NodeSelectorOperator(lreq.Operator),
				Values:   lreq.Values,
			})
		}

		var lresources v1.ResourceList
		if linstanceSpec.GetResources() != nil {
			lresources = v1.ResourceList{}
			for lresname, lresval := range linstanceSpec.GetResources() {
				quantity, lerr := resource.ParseQuantity(lresval)
				if lerr != nil {
					return nil, fmt.Errorf("error parsing quantity: %s", lerr)
				}

				switch lresname {
				case "memory":
					lresources[v1.ResourceMemory] = quantity
				case "cpu	":
					lresources[v1.ResourceCPU] = quantity
				case "disk":
					lresources[v1.ResourceStorage] = quantity
				}
			}
		}

		lmachine := &v1alpha5.Machine{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					v1alpha5.ProvisionerNameLabelKey: "playrix", // Задаем здесь имя которое нельзя найти через k8spai, karpenter фалбечит этот случай и List(пвозвращает список инстансов) отрабатывает без ошибок
				},
			},
			Spec: v1alpha5.MachineSpec{
				Kubelet: &v1alpha5.KubeletConfiguration{
					KubeReserved: v1.ResourceList{
						v1.ResourceMemory: *resource.NewQuantity(0, resource.DecimalSI),
						v1.ResourceCPU:    *resource.NewQuantity(0, resource.DecimalSI),
					},

					EvictionHard: map[string]string{
						"memory.available": "0",
					},

					EvictionSoft: map[string]string{
						"memory.available": "0",
					},
				},

				Requirements: lrequirements,

				Resources: v1alpha5.ResourceRequirements{
					Requests: lresources,
				},

				MachineTemplateRef: &v1alpha5.MachineTemplateRef{
					Name: "template",
				},
			},
		}

		execCtx = context.WithValue(execCtx, settings.ContextKey, lexecsettings)
		linstanceInfo, lerr := s.karpenderCloudProvider.Create(execCtx, lmachine)
		if lerr != nil {
			return &karpenterprovidergrpc.AddInstancesResponse{
				Instanseids: linstances,
				Reason:      lerr.Error(),
			}, nil
		}

		linstanceid, _ := utils.ParseInstanceID(linstanceInfo.Status.ProviderID)
		linstances = append(linstances, linstanceid)
	}

	return &karpenterprovidergrpc.AddInstancesResponse{Instanseids: linstances}, nil
}

func (s *server) RemoveInstances(_ctx context.Context, _d *karpenterprovidergrpc.DeleteInstancesRequest) (*karpenterprovidergrpc.DeleteInstancesResponse, error) {
	lExecsettings := &settings.Settings{ //заполняется кластером k8s, не то чтобы очень полезная инфа
		ClusterName:                s.clusterName + "/" + _d.PoolName,
		ClusterEndpoint:            s.clusterEndpoint,
		DefaultInstanceProfile:     "",
		EnablePodENI:               false,
		EnableENILimitedPodDensity: true,
		IsolatedVPC:                false,
		VMMemoryOverheadPercent:    RESCOURCERESERVATIONPERSENT,
		InterruptionQueueName:      "",
		Tags:                       map[string]string{},
		ReservedENIs:               0,
	}
	execCtx := context.WithValue(_ctx, settings.ContextKey, lExecsettings)

	for _, linstanceId := range _d.Instanseids {
		lmachine := &v1alpha5.Machine{
			Status: v1alpha5.MachineStatus{
				ProviderID: fmt.Sprintf("aws:///region/%s", linstanceId),
			},
		}

		s.karpenderCloudProvider.Delete(execCtx, lmachine)
	}

	return &karpenterprovidergrpc.DeleteInstancesResponse{}, nil
}

type K8sKapenterProviderPlugin struct {
	plugin.Plugin
	impl *server
}

func (p *K8sKapenterProviderPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	karpenterprovidergrpc.RegisterKarpenterServiceServer(s, p.impl)
	return nil
}

func (p *K8sKapenterProviderPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return nil, fmt.Errorf("client not allowed here")
}

func NewK8sKapenterProviderPlugin(_clusterName string) *K8sKapenterProviderPlugin {
	lclusternameEndpoint := fmt.Sprintf("%s.endpoint", _clusterName)

	config := &aws.Config{
		Region:              aws.String(os.Getenv("AWS_DEFAULT_REGION")),
		STSRegionalEndpoint: endpoints.RegionalSTSEndpoint,
	}

	sess := withUserAgent(session.Must(session.NewSession(
		request.WithRetryer(
			config,
			awsclient.DefaultRetryer{NumMaxRetries: awsclient.DefaultRetryerMaxNumRetries},
		),
	)))

	//---------------------------------------------------------------------------
	lsettings := &settings.Settings{ //заполняется кластером k8s
		ClusterName:                _clusterName, // не похоже что при создании экземпляра cloudprovider это используется
		ClusterEndpoint:            lclusternameEndpoint,
		DefaultInstanceProfile:     "",
		EnablePodENI:               false,
		EnableENILimitedPodDensity: true,
		IsolatedVPC:                false,
		VMMemoryOverheadPercent:    RESCOURCERESERVATIONPERSENT,
		InterruptionQueueName:      "",
		Tags:                       map[string]string{},
		ReservedENIs:               0,
	}

	creatCtx := context.WithValue(context.TODO(), settings.ContextKey, lsettings)
	ec2api := ec2.New(sess)
	ssmapi := ssm.New(sess)
	unavailableOfferingsCache := awscache.NewUnavailableOfferings()

	//---------------------------------------------------------------------------
	subnetProvider := subnet.NewProvider(ec2api, cache.New(awscache.DefaultTTL, awscache.DefaultCleanupInterval))

	//---------------------------------------------------------------------------
	securityGroupProvider := securitygroup.NewProvider(ec2api, cache.New(awscache.DefaultTTL, awscache.DefaultCleanupInterval))

	//---------------------------------------------------------------------------
	fakek8s := &fakeKubernetesInterface{}
	amiProvider := amifamily.NewProvider(nil, fakek8s, ssmapi, ec2api, cache.New(awscache.DefaultTTL, awscache.DefaultCleanupInterval), cache.New(awscache.DefaultTTL, awscache.DefaultCleanupInterval), cache.New(awscache.DefaultTTL, awscache.DefaultCleanupInterval))

	//---------------------------------------------------------------------------
	pricingProvider := pricing.NewProvider(
		creatCtx,
		pricing.NewAPI(sess, *sess.Config.Region),
		ec2api,
		*sess.Config.Region,
	)

	//---------------------------------------------------------------------------
	instanceTypeProvider := instancetype.NewProvider(
		*sess.Config.Region,
		cache.New(awscache.InstanceTypesAndZonesTTL, awscache.DefaultCleanupInterval),
		ec2api,
		subnetProvider,
		unavailableOfferingsCache,
		pricingProvider,
	)

	//---------------------------------------------------------------------------
	amiResolver := amifamily.New(amiProvider)

	leaderCh := make(chan struct{})
	close(leaderCh)

	launchTemplateProvider := launchtemplate.NewProvider(
		creatCtx,
		cache.New(awscache.DefaultTTL, awscache.DefaultCleanupInterval),
		ec2api,
		amiResolver,
		securityGroupProvider,
		subnetProvider,
		nil,
		leaderCh,
		nil,
		lclusternameEndpoint,
	)

	//---------------------------------------------------------------------------
	instanceProvider := instance.NewProvider(
		creatCtx,
		aws.StringValue(sess.Config.Region),
		ec2api,
		unavailableOfferingsCache,
		instanceTypeProvider,
		subnetProvider,
		launchTemplateProvider,
	)

	//---------------------------------------------------------------------------
	fakek8sclient := &fakek8sClient{}
	cloudProvider := cloudprovider.New(instanceTypeProvider, instanceProvider, fakek8sclient, amiProvider, securityGroupProvider, subnetProvider)

	return &K8sKapenterProviderPlugin{
		impl: &server{
			clusterName:            _clusterName,
			clusterEndpoint:        lclusternameEndpoint,
			fakek8sAPI:             fakek8sclient,
			karpenderCloudProvider: cloudProvider,
		},
	}
}
