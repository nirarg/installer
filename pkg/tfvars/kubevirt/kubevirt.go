// Package kubevirt contains kubevirt-specific Terraform-variable logic.
package kubevirt

import (
	"encoding/json"

	v1 "github.com/openshift/cluster-api-provider-kubevirt/pkg/apis/kubevirtprovider/v1alpha1"
	// "github.com/openshift/installer/pkg/rhcos"
	// "github.com/openshift/installer/pkg/tfvars/internal/cache"
)

type config struct {
	Namespace                  string            `json:"kubevirt_namespace"`
	ImageURL                   string            `json:"kubevirt_image_url"`
	SourcePvcName              string            `json:"kubevirt_source_pvc_name"`
	SecretName                 string            `json:"kubevirt_secret_name"`
	RequestedMemory            string            `json:"kubevirt_master_memory"`
	RequestedCPU               uint32            `json:"kubevirt_master_cpu"`
	RequestedStorage           string            `json:"kubevirt_master_storage"`
	MachineType                string            `json:"kubevirt_master_machine_type"`
	StorageClass               string            `json:"kubevirt_storage_class"`
	NetworkName                string            `json:"kubevirt_network_name"`
	PersistentVolumeAccessMode string            `json:"kubevirt_pv_access_mode"`
	ResourcesLabels            map[string]string `json:"kubevirt_labels"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	MasterSpecs     []*v1.KubevirtMachineProviderSpec
	ImageURL        string
	Namespace       string
	ResourcesLabels map[string]string
}

// TFVars generates kubevirt-specific Terraform variables.
func TFVars(sources TFVarsSources) ([]byte, error) {
	masterSpec := sources.MasterSpecs[0]

	// For optional parametes, set only if not nil
	cfg := config{
		Namespace:                  sources.Namespace,
		ImageURL:                   sources.ImageURL,
		SourcePvcName:              masterSpec.SourcePvcName,
		RequestedMemory:            masterSpec.RequestedMemory,
		RequestedCPU:               masterSpec.RequestedCPU,
		RequestedStorage:           masterSpec.RequestedStorage,
		StorageClass:               masterSpec.StorageClassName,
		NetworkName:                masterSpec.NetworkName,
		PersistentVolumeAccessMode: masterSpec.PersistentVolumeAccessMode,
		ResourcesLabels:            sources.ResourcesLabels,
	}

	return json.MarshalIndent(cfg, "", "  ")
}
