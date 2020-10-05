package defaults

import (
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/kubevirt"
)

// SetPlatformDefaults sets the defaults for the platform.
func SetPlatformDefaults(p *kubevirt.Platform, controlPlane *types.MachinePool, compute []types.MachinePool) {
	if controlPlane.Platform.Kubevirt == nil {
		controlPlane.Platform.Kubevirt = &kubevirt.MachinePool{
			CPU:         8,
			Memory:      "16G",
			StorageSize: "120Gi",
		}
	}
	for i := range compute {
		if compute[i].Platform.Kubevirt == nil {
			compute[i].Platform.Kubevirt = &kubevirt.MachinePool{
				CPU:         4,
				Memory:      "10G",
				StorageSize: "120Gi",
			}
		}
	}
}
