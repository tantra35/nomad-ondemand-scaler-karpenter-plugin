package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/dustin/go-humanize"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/api/resource"
	"playrix.com/it/nomad-cluster-scalerv2/karpenterprovidergrpc/karpenterprovidergrpc"
)

type TestK8sKapenterProviderPluginInterface interface {
	ListInstances(string) ([]string, error)
	AddInstances(string, int, *karpenterprovidergrpc.AddInstancesSpec) ([]string, error)
	RemoveInstances(string, []string) error
}

type TestK8sKapenterProviderPlugin struct {
	plugin.Plugin

	Impl TestK8sKapenterProviderPluginInterface
}

func (p *TestK8sKapenterProviderPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	return fmt.Errorf("server not allowed here")
}

func (p *TestK8sKapenterProviderPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &TestK8sKapenterProviderClient{
		client: karpenterprovidergrpc.NewKarpenterServiceClient(c),
	}, nil
}

type TestK8sKapenterProviderClient struct {
	client karpenterprovidergrpc.KarpenterServiceClient
}

func (k *TestK8sKapenterProviderClient) ListInstances(_poolName string) ([]string, error) {
	resp, lerr := k.client.ListInstances(context.Background(), &karpenterprovidergrpc.ListInstancesRequest{
		PoolName: _poolName,
	})
	if lerr != nil {
		return nil, lerr
	}

	return resp.Instanseids, nil
}

func (k *TestK8sKapenterProviderClient) AddInstances(_poolName string, _count int, _spec *karpenterprovidergrpc.AddInstancesSpec) ([]string, error) {
	resp, lerr := k.client.AddInstances(context.Background(), &karpenterprovidergrpc.AddInstancesRequest{
		PoolName: _poolName,
		Count:    int32(_count),
		Spec:     _spec,
	})

	return resp.Instanseids, lerr
}

func (k *TestK8sKapenterProviderClient) RemoveInstances(_poolName string, _instanses []string) error {
	_, lerr := k.client.RemoveInstances(context.Background(), &karpenterprovidergrpc.DeleteInstancesRequest{
		PoolName:    _poolName,
		Instanseids: _instanses,
	})

	return lerr
}

// -----------------------------------------------------------------------------
var createPluginClientOnce sync.Once
var rpcClient plugin.ClientProtocol

// -----------------------------------------------------------------------------
func TestPluginListInstances(t *testing.T) {
	// vault read secrets/aws/plr/atf01/creds/full
	os.Setenv("AWS_ACCESS_KEY_ID", "AKIAZ2XHWNGPTOCC344B")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "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")
	os.Setenv("AWS_DEFAULT_REGION", "eu-central-1")

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			// This isn't required when using VersionedPlugins
			ProtocolVersion:  1,
			MagicCookieKey:   "BASIC_PLUGIN",
			MagicCookieValue: "hello",
		},
		Plugins: map[string]plugin.Plugin{
			"grpc": &TestK8sKapenterProviderPlugin{},
		},
		Cmd: exec.Command("./karpenterprovider_plugin.exe", "somenamedcluster"),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
	})

	rpcClient, lerr := client.Client()
	if lerr != nil {
		t.Fatal(lerr)
	}

	raw, lerr := rpcClient.Dispense("grpc")
	if lerr != nil {
		t.Fatal(lerr)
	}

	svc := raw.(TestK8sKapenterProviderPluginInterface)
	instances, lerr := svc.ListInstances("testpool")
	if lerr != nil {
		t.Fatal(lerr)
	}

	t.Logf("instances: %v", instances)
}

func TestPluginListAddInstances(t *testing.T) {
	// vault read secrets/aws/plr/atf01/creds/full
	os.Setenv("AWS_ACCESS_KEY_ID", "AKIAZ2XHWNGPSMQSNSNU")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "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")
	os.Setenv("AWS_DEFAULT_REGION", "eu-central-1")

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			// This isn't required when using VersionedPlugins
			ProtocolVersion:  1,
			MagicCookieKey:   "BASIC_PLUGIN",
			MagicCookieValue: "hello",
		},
		Plugins: map[string]plugin.Plugin{
			"grpc": &TestK8sKapenterProviderPlugin{},
		},
		Cmd: exec.Command("./karpenterprovider_plugin.exe", "somenamedcluster"),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
	})

	rpcClient, lerr := client.Client()
	if lerr != nil {
		t.Fatal(lerr)
	}

	defer rpcClient.Close()

	raw, lerr := rpcClient.Dispense("grpc")
	if lerr != nil {
		t.Fatal(lerr)
	}

	uintB, _ := humanize.ParseBytes("3.8GiB")

	svc := raw.(TestK8sKapenterProviderPluginInterface)
	_, lerr = svc.AddInstances("testpool", 3, &karpenterprovidergrpc.AddInstancesSpec{
		/*
			SecurityGroups: map[string]string{
				"Name": "vpc+playrix",
			},

			Ami: map[string]string{
				"aws::owners": "self",
				"Name":        "binddnscisco",
				"Version":     "38",
			},

			InstanceProfile: "Docker",
		*/

		LaunchTemplate: aws.String("test-nomadscalerv2-karpenter"),

		Subnets: map[string]string{
			"aws-ids": "subnet-428e112a,subnet-0442eb1148e03dc43",
		},

		Requirements: []*karpenterprovidergrpc.Requirement{
			/*
				{
					Key:      "node.kubernetes.io/instance-type",
					Operator: "In",
					Values:   []string{"t3a.medium", "t2.medium"},
				},
			*/

			{
				Key:      "karpenter.k8s.aws/instance-family",
				Operator: "In",
				Values:   []string{"t3a", "t2"},
			},

			{
				Key:      "karpenter.sh/capacity-type",
				Operator: "In",
				Values:   []string{"spot", "on-demand"},
			},
		},

		Resources: map[string]string{
			"memory": fmt.Sprintf("%dMi", uintB/1024/1024),
		},
	})

	if lerr != nil {
		t.Fatal(lerr)
	}
}

func TestPluginListAddInstancesSparkDriver(t *testing.T) {
	// vault read secrets/aws/plr/amrv01/creds/full
	os.Setenv("AWS_ACCESS_KEY_ID", "AKIAQGLVH5TLCIGIHX7Y")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "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")
	os.Setenv("AWS_DEFAULT_REGION", "us-east-1")

	createPluginClientOnce.Do(func() {
		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: plugin.HandshakeConfig{
				// This isn't required when using VersionedPlugins
				ProtocolVersion:  1,
				MagicCookieKey:   "BASIC_PLUGIN",
				MagicCookieValue: "hello",
			},
			Plugins: map[string]plugin.Plugin{
				"grpc": &TestK8sKapenterProviderPlugin{},
			},
			Cmd: exec.Command("./karpenterprovider_plugin.exe", "somenamedcluster"),
			AllowedProtocols: []plugin.Protocol{
				plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
		})

		lrpcClient, lerr := client.Client()
		if lerr != nil {
			t.Fatal(lerr)
		}

		rpcClient = lrpcClient
	})

	defer rpcClient.Close()

	raw, lerr := rpcClient.Dispense("grpc")
	if lerr != nil {
		t.Fatal(lerr)
	}

	uintB, _ := humanize.ParseBytes("31948Mib")
	unitCPU := fmt.Sprintf("8000m")
	cpuq, lerr := resource.ParseQuantity(unitCPU)
	if lerr != nil {
		t.Fatal(lerr)
	}

	t.Logf("Cpu: %s", &cpuq)

	svc := raw.(TestK8sKapenterProviderPluginInterface)
	linstancelist, lerr := svc.AddInstances("sparkdriver", 1, &karpenterprovidergrpc.AddInstancesSpec{
		LaunchTemplate: aws.String("dockerworker-marketingrr-sparkdriver-t3.2xlarge"),

		Subnets: map[string]string{
			"Name": "natted-us-east-1a",
		},

		Requirements: []*karpenterprovidergrpc.Requirement{
			{
				Key:      "karpenter.k8s.aws/instance-family",
				Operator: "In",
				Values:   []string{"t3a", "t3"},
			},

			{
				Key:      "karpenter.sh/capacity-type",
				Operator: "In",
				Values:   []string{"on-demand"},
			},
		},

		Resources: map[string]string{
			"cpu":    unitCPU,
			"memory": fmt.Sprintf("%dMi", uintB/1024/1024),
		},
	})

	if lerr != nil {
		t.Fatal(lerr)
	}

	t.Logf("linstancelist: %v", linstancelist)
}

func createPluginClient(t *testing.T) (plugin.ClientProtocol, TestK8sKapenterProviderPluginInterface) {
	createPluginClientOnce.Do(func() {
		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: plugin.HandshakeConfig{
				// This isn't required when using VersionedPlugins
				ProtocolVersion:  1,
				MagicCookieKey:   "BASIC_PLUGIN",
				MagicCookieValue: "hello",
			},
			Plugins: map[string]plugin.Plugin{
				"grpc": &TestK8sKapenterProviderPlugin{},
			},
			Cmd: exec.Command("./karpenterprovider_plugin.exe", "somenamedcluster"),
			AllowedProtocols: []plugin.Protocol{
				plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
		})

		lrpcClient, lerr := client.Client()
		if lerr != nil {
			t.Fatal(lerr)
		}

		rpcClient = lrpcClient
	})

	raw, lerr := rpcClient.Dispense("grpc")
	if lerr != nil {
		t.Fatal(lerr)
	}

	return rpcClient, raw.(TestK8sKapenterProviderPluginInterface)
}

func TestPluginSingleton(t *testing.T) {
	// vault read secrets/aws/plr/atf01/creds/full
	os.Setenv("AWS_ACCESS_KEY_ID", "AKIAZ2XHWNGPSMQSNSNU")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "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")
	os.Setenv("AWS_DEFAULT_REGION", "eu-central-1")

	c1, intf1 := createPluginClient(t)
	defer c1.Close()
	t.Logf("c1: %v", intf1)

	c2, intf2 := createPluginClient(t)
	defer c2.Close()
	t.Logf("c2: %v", intf2)

	t.Logf("c1 == c2: %v", c1 == c2)
}

func TestPluginDirectAddinstnce(t *testing.T) {
	// vault read secrets/aws/plr/atf01/creds/full
	os.Setenv("AWS_ACCESS_KEY_ID", "AKIAZ2XHWNGPWZSMGWCI")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "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")
	os.Setenv("AWS_DEFAULT_REGION", "eu-central-1")

	directPlugin := NewK8sKapenterProviderPlugin("gotesting")

	uintB, _ := humanize.ParseBytes("3.8GiB")
	ctx := context.Background()
	lreq := &karpenterprovidergrpc.AddInstancesRequest{
		PoolName: "gotestingpool",
		Count:    1,
		Spec: &karpenterprovidergrpc.AddInstancesSpec{

			LaunchTemplate: aws.String("test-nomadscalerv2-karpenter"),

			Subnets: map[string]string{
				"aws-ids": "subnet-428e112a,subnet-0442eb1148e03dc43",
			},

			Requirements: []*karpenterprovidergrpc.Requirement{
				/*
					{
						Key:      "node.kubernetes.io/instance-type",
						Operator: "In",
						Values:   []string{"t3a.medium", "t2.medium"},
					},
				*/

				{
					Key:      "karpenter.k8s.aws/instance-family",
					Operator: "In",
					Values:   []string{"t3a", "t2"},
				},

				{
					Key:      "karpenter.sh/capacity-type",
					Operator: "In",
					Values:   []string{"spot", "on-demand"},
				},
			},

			Resources: map[string]string{
				"memory": fmt.Sprintf("%dMi", uintB/1024/1024),
			},
		},
	}

	resp, lerr := directPlugin.impl.AddInstances(ctx, lreq)
	if lerr != nil {
		t.Fatalf("AddInstances error due: %v", lerr)
	}

	t.Logf("Added instances: %v", resp.GetInstanseids())
}
