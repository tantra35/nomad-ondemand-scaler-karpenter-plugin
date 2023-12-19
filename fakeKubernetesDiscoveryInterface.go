package main

import (
	openapi_v2 "github.com/google/gnostic/openapiv2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/openapi"
	restclient "k8s.io/client-go/rest"
)

type fakeKubernetesDiscoveryInterface struct {
}

func (f *fakeKubernetesDiscoveryInterface) RESTClient() restclient.Interface {
	return nil
}

func (f *fakeKubernetesDiscoveryInterface) ServerGroups() (*metav1.APIGroupList, error) {
	return nil, nil
}

func (f *fakeKubernetesDiscoveryInterface) ServerResourcesForGroupVersion(groupVersion string) (*metav1.APIResourceList, error) {
	return nil, nil
}

func (f *fakeKubernetesDiscoveryInterface) ServerGroupsAndResources() ([]*metav1.APIGroup, []*metav1.APIResourceList, error) {
	return nil, nil, nil
}

func (f *fakeKubernetesDiscoveryInterface) ServerPreferredResources() ([]*metav1.APIResourceList, error) {
	return nil, nil
}

func (f *fakeKubernetesDiscoveryInterface) ServerPreferredNamespacedResources() ([]*metav1.APIResourceList, error) {
	return nil, nil
}

func (f *fakeKubernetesDiscoveryInterface) ServerVersion() (*version.Info, error) {
	vi := &version.Info{
		Major: "1",
		Minor: "24",
	}

	return vi, nil
}

func (f *fakeKubernetesDiscoveryInterface) OpenAPISchema() (*openapi_v2.Document, error) {
	return nil, nil
}

func (f *fakeKubernetesDiscoveryInterface) OpenAPIV3() openapi.Client {
	return nil
}
