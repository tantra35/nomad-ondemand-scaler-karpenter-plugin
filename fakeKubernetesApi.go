package main

import (
	"context"
	"fmt"

	"github.com/aws/karpenter/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type fakek8sClientProvisioner struct {
}

var fakek8sClientProvisionerKey = fakek8sClientProvisioner{}

type fakek8sClientTemplate struct {
}

var fakek8sClientTemplateKey = fakek8sClientTemplate{}

type fakek8sClient struct {
}

type fakek8ApiErrorApiStatus struct {
	status metav1.Status
}

func (s *fakek8ApiErrorApiStatus) Error() string {
	return ""
}

func (s *fakek8ApiErrorApiStatus) Status() metav1.Status {
	return s.status
}

func (f *fakek8sClient) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	switch key.String() {
	case "/template":
		nodetemplate := ctx.Value(fakek8sClientTemplateKey).(*v1alpha1.AWSNodeTemplate)
		lobj, lok := obj.(*v1alpha1.AWSNodeTemplate)
		if !lok {
			return fmt.Errorf("not an AWSNodeTemplate")
		}

		nodetemplate.DeepCopyInto(lobj)
		return nil
	}

	return &fakek8ApiErrorApiStatus{metav1.Status{Reason: metav1.StatusReasonNotFound}}
}

func (f *fakek8sClient) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	return fmt.Errorf("not implemented")
}

func (f *fakek8sClient) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	return fmt.Errorf("not implemented")
}

func (f *fakek8sClient) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
	return fmt.Errorf("not implemented")
}

func (f *fakek8sClient) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	return fmt.Errorf("not implemented")
}

func (f *fakek8sClient) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	return fmt.Errorf("not implemented")
}

func (f *fakek8sClient) DeleteAllOf(ctx context.Context, obj client.Object, opts ...client.DeleteAllOfOption) error {
	return fmt.Errorf("not implemented")
}

func (f *fakek8sClient) Status() client.StatusWriter {
	return nil
}

// Scheme returns the scheme this client is using.
func (f *fakek8sClient) Scheme() *runtime.Scheme {
	return nil
}

func (f *fakek8sClient) RESTMapper() meta.RESTMapper {
	return nil
}
