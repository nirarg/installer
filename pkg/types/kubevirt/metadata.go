package kubevirt

// Metadata contains kubevirt metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	KubeConfig string            `json:"kubeconfig"`
	Namespace  string            `json:"namespace"`
	Labels     map[string]string `json:"labels"`
}
