package kubevirt

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/kubevirt"
	"github.com/openshift/installer/pkg/types/kubevirt/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// Validate executes kubevirt specific validation
func Validate(ic *types.InstallConfig) error {
	kubevirtPlatformPath := field.NewPath("platform", "kubevirt")

	if ic.Platform.Kubevirt == nil {
		return errors.New(field.Required(
			kubevirtPlatformPath,
			"validation requires a Engine platform configuration").Error())
	}

	return validatePlatform(ic.Platform.Kubevirt, ic.MachineNetwork, kubevirtPlatformPath).ToAggregate()
}

func validatePlatform(kubevirtPlatform *kubevirt.Platform, machineNetworkEntryList []types.MachineNetworkEntry, fldPath *field.Path) field.ErrorList {
	allErrs := validation.ValidatePlatform(kubevirtPlatform, fldPath)
	ctx := context.Background()

	client, resultErrs := validateInfraClusterReachable(ctx, kubevirtPlatform.Namespace, fldPath)
	allErrs = append(allErrs, resultErrs...)
	if client != nil {
		allErrs = append(allErrs, validateStorageClassExistsInInfraCluster(ctx, kubevirtPlatform.StorageClass, client, fldPath)...)
		allErrs = append(allErrs, validateNetworkAttachmentDefinitionExistsInInfraCluster(ctx, kubevirtPlatform.NetworkName, kubevirtPlatform.Namespace, client, fldPath)...)
	}
	allErrs = append(allErrs, validateIPsInMachineNetworkEntryList(machineNetworkEntryList, kubevirtPlatform.APIVIP, kubevirtPlatform.IngressVIP, fldPath)...)

	return allErrs
}

// validateInfraClusterReachable validates the following:
// 1. Client can be created -  login to cluster is done or KUBECONFIG environment variable was exported
// 2. Cluster is reachable - valid user
// 3. Namespace exists
// 4. Kubevirt installed
// 5. User can list VMs in namespace
func validateInfraClusterReachable(ctx context.Context, namespace string, fieldPath *field.Path) (Client, field.ErrorList) {
	allErrs := field.ErrorList{}
	client, err := NewClient()
	if err != nil {
		detailedErr := fmt.Errorf("failed to create InfraCluster client with error: %v", err)
		allErrs = append(allErrs, field.Invalid(fieldPath.Child("InfraClusterReachable"), "InfraCluster", detailedErr.Error()))

		return nil, allErrs
	}

	if _, err := client.ListVirtualMachine(namespace); err != nil {
		detailedErr := fmt.Errorf("failed to access to InfraCluster with error: %v", err)
		allErrs = append(allErrs, field.Invalid(fieldPath.Child("InfraClusterReachable"), "InfraCluster", detailedErr.Error()))

		return nil, allErrs
	}

	return client, allErrs
}

func validateStorageClassExistsInInfraCluster(ctx context.Context, name string, client Client, fieldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if _, err := client.GetStorageClass(ctx, name); err != nil {
		detailedErr := fmt.Errorf("failed to get storageClass %s from InfraCluster, with error: %v", name, err)
		allErrs = append(allErrs, field.Invalid(fieldPath.Child("StorageClassExistsInInfraCluster"), name, detailedErr.Error()))
	}

	return allErrs
}

func validateNetworkAttachmentDefinitionExistsInInfraCluster(ctx context.Context, name string, namespace string, client Client, fieldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if _, err := client.GetNetworkAttachmentDefinition(ctx, name, namespace); err != nil {
		detailedErr := fmt.Errorf("failed to get network-attachment-definition %s from InfraCluster, with error: %v", name, err)
		allErrs = append(allErrs, field.Invalid(fieldPath.Child("NetworkAttachmentDefinitionExistsInInfraCluster"), name, detailedErr.Error()))
	}

	return allErrs
}

func validateIPsInMachineNetworkEntryList(machineNetworkEntryList []types.MachineNetworkEntry, apiVIP string, ingressVIP string, fieldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if err := assertIPInMachineNetworkEntryList(machineNetworkEntryList, apiVIP); err != nil {
		detailedErr := fmt.Errorf("validation of apiVIP %s in cider %s failed, with error: %v", apiVIP, machineNetworkEntryList, err)
		allErrs = append(allErrs, field.Invalid(fieldPath.Child("IPsInCIDR"), apiVIP, detailedErr.Error()))
	}

	if err := assertIPInMachineNetworkEntryList(machineNetworkEntryList, ingressVIP); err != nil {
		detailedErr := fmt.Errorf("validation of ingressVIP %s in cider %s failed, with error: %v", ingressVIP, machineNetworkEntryList, err)
		allErrs = append(allErrs, field.Invalid(fieldPath.Child("IPsInCIDR"), ingressVIP, detailedErr.Error()))
	}

	return allErrs
}

func assertIPInMachineNetworkEntryList(machineNetworkEntryList []types.MachineNetworkEntry, ip string) error {
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return fmt.Errorf("ip %s is not valid IP address", ip)
	}
	for _, machineNetworkEntry := range machineNetworkEntryList {
		if machineNetworkEntry.CIDR.Contains(ipAddr) {
			return nil
		}
		return fmt.Errorf("ip %s not in machineNetworkEntryList %s", ip, machineNetworkEntryList)
	}
	return nil
}
