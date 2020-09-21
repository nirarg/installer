/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kubevirt

import (
	"context"
	"fmt"
	"log"

	nadv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	kubevirtapiv1 "kubevirt.io/client-go/api/v1"
	cdiv1 "kubevirt.io/containerized-data-importer/pkg/apis/core/v1alpha1"
)

//go:generate mockgen -source=./client.go -destination=./mock/client_generated.go -package=mock

type Client interface {
	GetVirtualMachine(namespace string, name string) (*kubevirtapiv1.VirtualMachine, error)
	ListVirtualMachine(namespace string) (*kubevirtapiv1.VirtualMachineList, error)
	DeleteVirtualMachine(namespace string, name string) error
	GetDataVolume(namespace string, name string) (*cdiv1.DataVolume, error)
	ListDataVolume(namespace string) (*cdiv1.DataVolumeList, error)
	DeleteDataVolume(namespace string, name string) error
	GetSecret(namespace string, name string) (*corev1.Secret, error)
	ListSecret(namespace string) (*corev1.SecretList, error)
	DeleteSecret(namespace string, name string) error
	GetStorageClass(ctx context.Context, name string) (*storagev1.StorageClass, error)
	GetNetworkAttachmentDefinition(ctx context.Context, name string, namespace string) (*unstructured.Unstructured, error)
}

type client struct {
	kubernetesClient *kubernetes.Clientset
	dynamicClient    dynamic.Interface
}

// NewClient creates our client wrapper object for the actual kubeVirt and kubernetes clients we use.
func NewClient() (Client, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	// if you want to change the loading rules (which files in which order), you can do so here

	configOverrides := &clientcmd.ConfigOverrides{}
	// if you want to change override values or bind them to flags, there are methods to help you

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	restClientConfig, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	result := &client{}

	if result.kubernetesClient, err = kubernetes.NewForConfig(restClientConfig); err != nil {
		return nil, err
	}
	if result.dynamicClient, err = dynamic.NewForConfig(restClientConfig); err != nil {
		return nil, err
	}
	return result, nil
}

// VirtualMachine

func (c *client) GetVirtualMachine(namespace string, name string) (*kubevirtapiv1.VirtualMachine, error) {
	resp, err := c.getResource(namespace, name, vmRes())
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, err
		}
		msg := fmt.Sprintf("Failed to get VirtualMachine, with error: %v", err)
		log.Printf("[Error] %s", msg)
		return nil, fmt.Errorf(msg)
	}
	var vm kubevirtapiv1.VirtualMachine
	err = c.fromUnstructedToInterface(*resp, &vm, "VirtualMachine")
	return &vm, err
}

func (c *client) ListVirtualMachine(namespace string) (*kubevirtapiv1.VirtualMachineList, error) {
	resp, err := c.listResource(namespace, vmRes())
	if err != nil {
		msg := fmt.Sprintf("Failed to list VirtualMachine, with error: %v", err)
		log.Printf("[Error] %s", msg)
		return nil, fmt.Errorf(msg)
	}
	var vmList kubevirtapiv1.VirtualMachineList
	err = c.fromUnstructedListToInterface(*resp, &vmList, "VirtualMachineList")
	return &vmList, err
}

func (c *client) DeleteVirtualMachine(namespace string, name string) error {
	return c.deleteResource(namespace, name, vmRes())
}

func vmRes() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    kubevirtapiv1.GroupVersion.Group,
		Version:  kubevirtapiv1.GroupVersion.Version,
		Resource: "virtualmachines",
	}

}

// DataVolume

func (c *client) GetDataVolume(namespace string, name string) (*cdiv1.DataVolume, error) {
	resp, err := c.getResource(namespace, name, dvRes())
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, err
		}
		msg := fmt.Sprintf("Failed to get DataVolume, with error: %v", err)
		log.Printf("[Error] %s", msg)
		return nil, fmt.Errorf(msg)
	}
	var dv cdiv1.DataVolume
	err = c.fromUnstructedToInterface(*resp, &dv, "DataVolume")
	return &dv, nil
}

func (c *client) ListDataVolume(namespace string) (*cdiv1.DataVolumeList, error) {
	resp, err := c.listResource(namespace, dvRes())
	if err != nil {
		msg := fmt.Sprintf("Failed to list DataVolume, with error: %v", err)
		log.Printf("[Error] %s", msg)
		return nil, fmt.Errorf(msg)
	}
	var dvList cdiv1.DataVolumeList
	err = c.fromUnstructedListToInterface(*resp, &dvList, "DataVolumeList")
	return &dvList, err
}

func (c *client) DeleteDataVolume(namespace string, name string) error {
	return c.deleteResource(namespace, name, dvRes())
}

func dvRes() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    cdiv1.SchemeGroupVersion.Group,
		Version:  cdiv1.SchemeGroupVersion.Version,
		Resource: "datavolumes",
	}
}

// Secret

func (c *client) GetSecret(namespace string, name string) (*corev1.Secret, error) {
	return c.kubernetesClient.CoreV1().Secrets(namespace).Get(context.Background(), name, metav1.GetOptions{})
}

func (c *client) ListSecret(namespace string) (*corev1.SecretList, error) {
	return c.kubernetesClient.CoreV1().Secrets(namespace).List(context.Background(), metav1.ListOptions{})
}

func (c *client) DeleteSecret(namespace string, name string) error {
	return c.kubernetesClient.CoreV1().Secrets(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
}

// StorageClass

func (c *client) GetStorageClass(ctx context.Context, name string) (*storagev1.StorageClass, error) {
	return c.kubernetesClient.StorageV1().StorageClasses().Get(ctx, name, metav1.GetOptions{})
}

// NetworkAttachmentDefinition

func (c *client) GetNetworkAttachmentDefinition(ctx context.Context, name string, namespace string) (*unstructured.Unstructured, error) {
	nadRes := schema.GroupVersionResource{
		Group:    nadv1.SchemeGroupVersion.Group,
		Version:  nadv1.SchemeGroupVersion.Version,
		Resource: "network-attachment-definitions",
	}
	return c.getResource(namespace, name, nadRes)
}

// dynamicClient resources

func (c *client) createResource(obj interface{}, namespace string, resource schema.GroupVersionResource) error {
	resultMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		msg := fmt.Sprintf("Failed to translate %s to Unstructed (for create operation), with error: %v", resource.Resource, err)
		log.Printf("[Error] %s", msg)
		return fmt.Errorf(msg)
	}
	input := unstructured.Unstructured{}
	input.SetUnstructuredContent(resultMap)
	resp, err := c.dynamicClient.Resource(resource).Namespace(namespace).Create(context.Background(), &input, meta_v1.CreateOptions{})
	if err != nil {
		msg := fmt.Sprintf("Failed to create %s, with error: %v", resource.Resource, err)
		log.Printf("[Error] %s", msg)
		return fmt.Errorf(msg)
	}
	unstructured := resp.UnstructuredContent()
	return runtime.DefaultUnstructuredConverter.FromUnstructured(unstructured, obj)
}

func (c *client) getResource(namespace string, name string, resource schema.GroupVersionResource) (*unstructured.Unstructured, error) {
	return c.dynamicClient.Resource(resource).Namespace(namespace).Get(context.Background(), name, metav1.GetOptions{})
}

func (c *client) deleteResource(namespace string, name string, resource schema.GroupVersionResource) error {
	return c.dynamicClient.Resource(resource).Namespace(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
}

func (c *client) listResource(namespace string, resource schema.GroupVersionResource) (*unstructured.UnstructuredList, error) {
	return c.dynamicClient.Resource(resource).Namespace(namespace).List(context.Background(), metav1.ListOptions{})
}

func (c *client) fromUnstructedToInterface(src unstructured.Unstructured, dst interface{}, interfaceType string) error {
	unstructured := src.UnstructuredContent()
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructured, dst); err != nil {
		msg := fmt.Sprintf("Failed to translate unstructed to %s, with error: %v", interfaceType, err)
		log.Printf("[Error] %s", msg)
		return fmt.Errorf(msg)
	}
	return nil
}

func (c *client) fromUnstructedListToInterface(src unstructured.UnstructuredList, dst interface{}, interfaceType string) error {
	unstructured := src.UnstructuredContent()
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructured, dst); err != nil {
		msg := fmt.Sprintf("Failed to translate unstructed to %s, with error: %v", interfaceType, err)
		log.Printf("[Error] %s", msg)
		return fmt.Errorf(msg)
	}
	return nil
}
