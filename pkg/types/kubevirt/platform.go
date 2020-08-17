package kubevirt

// Platform stores all the global configuration that all
// machinesets use.
type Platform struct {
	// The path to the kubeconfig file of the infra cluster
	InfraClusterKubeConfig string `json:"kubevirt_kubeconfig"`

	// The Ingress Domain of the infra cluster
	InfraClusterIngressDomain string `json:"ingress_domain"`

	// The Namespace in the infra cluster, which the control plane (master vms)
	// and the compute (worker vms) are installed in
	InfraClusterNamespace string `json:"namespace"`

	// The Storage Class used in the infra cluster
	InfraClusterStorageClass string `json:"storage_class"`

	// The nameof the PVCof the data volume created by the installer
	// To be used by cluster-api-provider-kubevirt, for the workers installation
	SourcePvcName string `json:"SourcePvcName,omitempty"`

	// The namespace in the underkube cluster to install the overkube in it
	// Namespace string `json:"kubevirt_namespace"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on kubevirt for machine pools which do not define their
	// own platform configuration.
	// Default will set the image field to the latest RHCOS image.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`
}
