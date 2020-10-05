package kubevirt

import (
	"fmt"
	"os"

	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/types/kubevirt"
)

// Platform collects kubevirt-specific configuration.
func Platform() (*kubevirt.Platform, error) {
	var (
		kubeConfig, namespace, apiVIP, ingressVIP, networkName, storageClass, persistentVolumeAccessMode string
		err                                                                                              error
	)

	if kubeConfig, err = selectKubeConfig(); err != nil {
		return nil, err
	}

	if namespace, err = selectNamespace(); err != nil {
		return nil, err
	}

	if apiVIP, err = selectAPIVIP(); err != nil {
		return nil, err
	}

	if ingressVIP, err = selectIngressVIP(); err != nil {
		return nil, err
	}

	if networkName, err = selectNetworkName(); err != nil {
		return nil, err
	}

	if storageClass, err = selectStorageClass(); err != nil {
		return nil, err
	}

	if persistentVolumeAccessMode, err = selectPersistentVolumeAccessMode(); err != nil {
		return nil, err
	}

	return &kubevirt.Platform{
		InfraClusterKubeConfig:     kubeConfig,
		InfraClusterNamespace:      namespace,
		InfraClusterStorageClass:   storageClass,
		APIVIP:                     apiVIP,
		IngressVIP:                 ingressVIP,
		NetworkName:                networkName,
		PersistentVolumeAccessMode: persistentVolumeAccessMode,
	}, nil
}

func selectKubeConfig() (string, error) {
	// Didn't find a way to use this value in terraform "provider" declaration
	// Current solution is that the user should put the kubeconfig in running directory
	// And it declared hard coded in the provider
	// Need to research for alternative way
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/kubeconfig", dir), nil
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

func selectAPIVIP() (string, error) {
	var selectedAPIVIP string

	defaultValue := ""

	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "API VIP",
				Help:    "An IP which will be served by bootstrap and then pivoted masters, using keepalived.",
				Default: defaultValue,
			},
		},
	}, &selectedAPIVIP)

	return selectedAPIVIP, err
}

func selectIngressVIP() (string, error) {
	var selectedIngressVIP string

	defaultValue := ""

	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Ingress VIP",
				Help:    "An external IP which routes to the default ingress controller.",
				Default: defaultValue,
			},
		},
	}, &selectedIngressVIP)

	return selectedIngressVIP, err
}

func selectNetworkName() (string, error) {
	var selectedNetworkName string

	defaultValue := ""

	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Network Name",
				Help:    "The target network of all the network interfaces of the nodes.",
				Default: defaultValue,
			},
		},
	}, &selectedNetworkName)

	return selectedNetworkName, err
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

func selectPersistentVolumeAccessMode() (string, error) {
	var selectedPersistentVolumeAccessMode string

	defaultValue := "ReadWriteMany"

	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Persistent Volume Access Mode",
				Help:    "The Access Mode should be used with the Persistent volumes [ReadWriteOnce,ReadOnlyMany,ReadWriteMany].",
				Default: defaultValue,
			},
		},
	}, &selectedPersistentVolumeAccessMode)

	return selectedPersistentVolumeAccessMode, err
}
