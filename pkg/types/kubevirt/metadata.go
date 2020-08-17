package kubevirt

// Metadata contains kubevirt metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	KubeConfig    string `json:"kubeconfig"`
	IngressDomain string `json:"ingress_domain"`
}
