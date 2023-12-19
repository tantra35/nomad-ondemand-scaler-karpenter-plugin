package main

import (
	discovery "k8s.io/client-go/discovery"
	admissionregistrationv1 "k8s.io/client-go/kubernetes/typed/admissionregistration/v1"
	admissionregistrationv1beta1 "k8s.io/client-go/kubernetes/typed/admissionregistration/v1beta1"
	internalv1alpha1 "k8s.io/client-go/kubernetes/typed/apiserverinternal/v1alpha1"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	appsv1beta1 "k8s.io/client-go/kubernetes/typed/apps/v1beta1"
	appsv1beta2 "k8s.io/client-go/kubernetes/typed/apps/v1beta2"
	authenticationv1 "k8s.io/client-go/kubernetes/typed/authentication/v1"
	authenticationv1beta1 "k8s.io/client-go/kubernetes/typed/authentication/v1beta1"
	authorizationv1 "k8s.io/client-go/kubernetes/typed/authorization/v1"
	authorizationv1beta1 "k8s.io/client-go/kubernetes/typed/authorization/v1beta1"
	autoscalingv1 "k8s.io/client-go/kubernetes/typed/autoscaling/v1"
	autoscalingv2 "k8s.io/client-go/kubernetes/typed/autoscaling/v2"
	autoscalingv2beta1 "k8s.io/client-go/kubernetes/typed/autoscaling/v2beta1"
	autoscalingv2beta2 "k8s.io/client-go/kubernetes/typed/autoscaling/v2beta2"
	batchv1 "k8s.io/client-go/kubernetes/typed/batch/v1"
	batchv1beta1 "k8s.io/client-go/kubernetes/typed/batch/v1beta1"
	certificatesv1 "k8s.io/client-go/kubernetes/typed/certificates/v1"
	certificatesv1beta1 "k8s.io/client-go/kubernetes/typed/certificates/v1beta1"
	coordinationv1 "k8s.io/client-go/kubernetes/typed/coordination/v1"
	coordinationv1beta1 "k8s.io/client-go/kubernetes/typed/coordination/v1beta1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	discoveryv1 "k8s.io/client-go/kubernetes/typed/discovery/v1"
	discoveryv1beta1 "k8s.io/client-go/kubernetes/typed/discovery/v1beta1"
	eventsv1 "k8s.io/client-go/kubernetes/typed/events/v1"
	eventsv1beta1 "k8s.io/client-go/kubernetes/typed/events/v1beta1"
	extensionsv1beta1 "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
	flowcontrolv1alpha1 "k8s.io/client-go/kubernetes/typed/flowcontrol/v1alpha1"
	flowcontrolv1beta1 "k8s.io/client-go/kubernetes/typed/flowcontrol/v1beta1"
	flowcontrolv1beta2 "k8s.io/client-go/kubernetes/typed/flowcontrol/v1beta2"
	networkingv1 "k8s.io/client-go/kubernetes/typed/networking/v1"
	networkingv1alpha1 "k8s.io/client-go/kubernetes/typed/networking/v1alpha1"
	networkingv1beta1 "k8s.io/client-go/kubernetes/typed/networking/v1beta1"
	nodev1 "k8s.io/client-go/kubernetes/typed/node/v1"
	nodev1alpha1 "k8s.io/client-go/kubernetes/typed/node/v1alpha1"
	nodev1beta1 "k8s.io/client-go/kubernetes/typed/node/v1beta1"
	policyv1 "k8s.io/client-go/kubernetes/typed/policy/v1"
	policyv1beta1 "k8s.io/client-go/kubernetes/typed/policy/v1beta1"
	rbacv1 "k8s.io/client-go/kubernetes/typed/rbac/v1"
	rbacv1alpha1 "k8s.io/client-go/kubernetes/typed/rbac/v1alpha1"
	rbacv1beta1 "k8s.io/client-go/kubernetes/typed/rbac/v1beta1"
	schedulingv1 "k8s.io/client-go/kubernetes/typed/scheduling/v1"
	schedulingv1alpha1 "k8s.io/client-go/kubernetes/typed/scheduling/v1alpha1"
	schedulingv1beta1 "k8s.io/client-go/kubernetes/typed/scheduling/v1beta1"
	storagev1 "k8s.io/client-go/kubernetes/typed/storage/v1"
	storagev1alpha1 "k8s.io/client-go/kubernetes/typed/storage/v1alpha1"
	storagev1beta1 "k8s.io/client-go/kubernetes/typed/storage/v1beta1"
)

type fakeKubernetesInterface struct {
}

func (f *fakeKubernetesInterface) Discovery() discovery.DiscoveryInterface {
	return &fakeKubernetesDiscoveryInterface{}
}

func (f *fakeKubernetesInterface) AdmissionregistrationV1() admissionregistrationv1.AdmissionregistrationV1Interface {
	return nil
}

func (f *fakeKubernetesInterface) AdmissionregistrationV1beta1() admissionregistrationv1beta1.AdmissionregistrationV1beta1Interface {
	return nil
}

func (f *fakeKubernetesInterface) InternalV1alpha1() internalv1alpha1.InternalV1alpha1Interface {
	return nil
}

func (f *fakeKubernetesInterface) AppsV1() appsv1.AppsV1Interface {
	return nil
}

func (f *fakeKubernetesInterface) AppsV1beta1() appsv1beta1.AppsV1beta1Interface {
	return nil
}

func (f *fakeKubernetesInterface) AppsV1beta2() appsv1beta2.AppsV1beta2Interface {
	return nil
}

func (f *fakeKubernetesInterface) AuthenticationV1() authenticationv1.AuthenticationV1Interface {
	return nil
}

func (f *fakeKubernetesInterface) AuthenticationV1beta1() authenticationv1beta1.AuthenticationV1beta1Interface {
	return nil
}

func (f *fakeKubernetesInterface) AuthorizationV1() authorizationv1.AuthorizationV1Interface {
	return nil
}

func (f *fakeKubernetesInterface) AuthorizationV1beta1() authorizationv1beta1.AuthorizationV1beta1Interface {
	return nil
}

func (f *fakeKubernetesInterface) AutoscalingV1() autoscalingv1.AutoscalingV1Interface {
	return nil
}

func (f *fakeKubernetesInterface) AutoscalingV2() autoscalingv2.AutoscalingV2Interface {
	return nil
}

func (f *fakeKubernetesInterface) AutoscalingV2beta1() autoscalingv2beta1.AutoscalingV2beta1Interface {
	return nil
}

func (f *fakeKubernetesInterface) AutoscalingV2beta2() autoscalingv2beta2.AutoscalingV2beta2Interface {
	return nil
}

func (f *fakeKubernetesInterface) BatchV1() batchv1.BatchV1Interface {
	return nil
}

func (f *fakeKubernetesInterface) BatchV1beta1() batchv1beta1.BatchV1beta1Interface {
	return nil
}

func (f *fakeKubernetesInterface) CertificatesV1() certificatesv1.CertificatesV1Interface {
	return nil
}

func (f *fakeKubernetesInterface) CertificatesV1beta1() certificatesv1beta1.CertificatesV1beta1Interface {
	return nil
}

func (f *fakeKubernetesInterface) CoordinationV1beta1() coordinationv1beta1.CoordinationV1beta1Interface {
	return nil
}

func (f *fakeKubernetesInterface) CoordinationV1() coordinationv1.CoordinationV1Interface {
	return nil
}

func (f *fakeKubernetesInterface) CoreV1() corev1.CoreV1Interface {
	return nil
}

func (f *fakeKubernetesInterface) DiscoveryV1() discoveryv1.DiscoveryV1Interface {
	return nil
}

func (f *fakeKubernetesInterface) DiscoveryV1beta1() discoveryv1beta1.DiscoveryV1beta1Interface {
	return nil
}

func (f *fakeKubernetesInterface) EventsV1() eventsv1.EventsV1Interface {
	return nil
}

func (f *fakeKubernetesInterface) EventsV1beta1() eventsv1beta1.EventsV1beta1Interface {
	return nil
}

func (f *fakeKubernetesInterface) ExtensionsV1beta1() extensionsv1beta1.ExtensionsV1beta1Interface {
	return nil
}

func (f *fakeKubernetesInterface) FlowcontrolV1alpha1() flowcontrolv1alpha1.FlowcontrolV1alpha1Interface {
	return nil
}

func (f *fakeKubernetesInterface) FlowcontrolV1beta1() flowcontrolv1beta1.FlowcontrolV1beta1Interface {
	return nil
}

func (f *fakeKubernetesInterface) FlowcontrolV1beta2() flowcontrolv1beta2.FlowcontrolV1beta2Interface {
	return nil
}

func (f *fakeKubernetesInterface) NetworkingV1() networkingv1.NetworkingV1Interface {
	return nil
}

func (f *fakeKubernetesInterface) NetworkingV1alpha1() networkingv1alpha1.NetworkingV1alpha1Interface {
	return nil
}

func (f *fakeKubernetesInterface) NetworkingV1beta1() networkingv1beta1.NetworkingV1beta1Interface {
	return nil
}

func (f *fakeKubernetesInterface) NodeV1() nodev1.NodeV1Interface {
	return nil
}

func (f *fakeKubernetesInterface) NodeV1alpha1() nodev1alpha1.NodeV1alpha1Interface {
	return nil
}

func (f *fakeKubernetesInterface) NodeV1beta1() nodev1beta1.NodeV1beta1Interface {
	return nil
}

func (f *fakeKubernetesInterface) PolicyV1() policyv1.PolicyV1Interface {
	return nil
}

func (f *fakeKubernetesInterface) PolicyV1beta1() policyv1beta1.PolicyV1beta1Interface {
	return nil
}

func (f *fakeKubernetesInterface) RbacV1() rbacv1.RbacV1Interface {
	return nil
}

func (f *fakeKubernetesInterface) RbacV1beta1() rbacv1beta1.RbacV1beta1Interface {
	return nil
}

func (f *fakeKubernetesInterface) RbacV1alpha1() rbacv1alpha1.RbacV1alpha1Interface {
	return nil
}

func (f *fakeKubernetesInterface) SchedulingV1alpha1() schedulingv1alpha1.SchedulingV1alpha1Interface {
	return nil
}

func (f *fakeKubernetesInterface) SchedulingV1beta1() schedulingv1beta1.SchedulingV1beta1Interface {
	return nil
}

func (f *fakeKubernetesInterface) SchedulingV1() schedulingv1.SchedulingV1Interface {
	return nil
}

func (f *fakeKubernetesInterface) StorageV1beta1() storagev1beta1.StorageV1beta1Interface {
	return nil
}

func (f *fakeKubernetesInterface) StorageV1() storagev1.StorageV1Interface {
	return nil
}

func (f *fakeKubernetesInterface) StorageV1alpha1() storagev1alpha1.StorageV1alpha1Interface {
	return nil
}
