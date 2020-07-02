// Package kubevirt extracts Kubevirt metadata from install configurations.
package kubevirt

import (
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/kubevirt"
)

// Metadata converts an install configuration to kubevirt metadata.
func Metadata(config *types.InstallConfig) *kubevirt.Metadata {
	return &kubevirt.Metadata{
		KubeConfig:    config.Kubevirt.InfraClusterKubeConfig,
		IngressDomain: config.Kubevirt.InfraClusterIngressDomain,
	}
}
