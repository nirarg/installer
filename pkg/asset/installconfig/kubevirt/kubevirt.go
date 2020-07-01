package kubevirt

import (
	"os"
	"path/filepath"

	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/types/kubevirt"
)

// Platform collects kubevirt-specific configuration.
func Platform() (*kubevirt.Platform, error) {
	var (
		kubeConfig, namespace, ingressDomain, storageClass string
		err                                                error
	)

	if kubeConfig, err = selectKubeConfig(); err != nil {
		return nil, err
	}

	if namespace, err = selectNamespace(); err != nil {
		return nil, err
	}

	if ingressDomain, err = selectIngressDomain(); err != nil {
		return nil, err
	}

	if storageClass, err = selectStorageClass(); err != nil {
		return nil, err
	}

	return &kubevirt.Platform{
		InfraClusterKubeConfig:    kubeConfig,
		InfraClusterIngressDomain: ingressDomain,
		InfraClusterNamespace:     namespace,
		InfraClusterStorageClass:  storageClass,
	}, nil
}

func selectKubeConfig() (string, error) {
	var selectedKubeConfig string

	defaultValue := ""
	home := os.Getenv("HOME")
	if home != "" {
		defaultValue = filepath.Join(home, "/.kube/config")
	}

	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "KubeConfig path",
				Help:    "The KubeConfig file path of the underkube cluster.",
				Default: defaultValue,
			},
		},
	}, &selectedKubeConfig)

	return selectedKubeConfig, err
}

func selectNamespace() (string, error) {
	var selectedNamespace string

	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Namespace",
				Help:    "The namespace, in the undercluster, where all the resources of the overcluster would be created.",
			},
		},
	}, &selectedNamespace)

	return selectedNamespace, err
}

func selectIngressDomain() (string, error) {
	var selectedIngressDomain string

	defaultValue := ""

	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Ingress Domain",
				Help:    "The ingress domain in the underkube would be use to access the apps.",
				Default: defaultValue,
			},
		},
	}, &selectedIngressDomain)

	return selectedIngressDomain, err
}

func selectStorageClass() (string, error) {
	var selectedStorageClass string

	defaultValue := ""

	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Storage Class",
				Help:    "The name of the storage class used in the infra ocp cluster.",
				Default: defaultValue,
			},
		},
	}, &selectedStorageClass)

	return selectedStorageClass, err
}
