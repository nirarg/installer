package kubevirt

// Platform stores all the global configuration that all
// machinesets use.
type Platform struct {
	// The path to the kubeconfig file of the infra cluster
	InfraClusterKubeConfig string `json:"kubeconfig"`

	// The Namespace in the infra cluster, which the control plane (master vms)
	// and the compute (worker vms) are installed in
	InfraClusterNamespace string `json:"namespace"`

	// The Storage Class used in the infra cluster
	InfraClusterStorageClass string `json:"storageClass"`

	// NetworkName is the target network of all the network interfaces of the nodes.
	NetworkName string `json:"networkName,omitempty"`

	// APIVIP is an IP which will be served by bootstrap and then pivoted masters, using keepalived
	APIVIP string `json:"apiVIP"`

	// IngressIP is an external IP which routes to the default ingress controller.
	IngressVIP string `json:"IngressVIP"`

	// PersistentVolumeAccessMode is the access mode should be use with the persistent volumes
	PersistentVolumeAccessMode string `json:"persistentVolumeAccessMode"`
}
