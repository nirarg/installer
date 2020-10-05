package defaults

import (
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/kubevirt"
)

// SetPlatformDefaults sets the defaults for the platform.
func SetPlatformDefaults(p *kubevirt.Platform, controlPlane *types.MachinePool, compute []types.MachinePool) {
	controlPlane.Platform = types.MachinePoolPlatform{
		Kubevirt: &kubevirt.MachinePool{
			CPU:         8,
			Memory:      "16G",
			StorageSize: "35Gi",
		},
	}
	for i := range compute {
		compute[i].Platform = types.MachinePoolPlatform{
			Kubevirt: &kubevirt.MachinePool{
				CPU:         4,
				Memory:      "10G",
				StorageSize: "35Gi",
			},
		}
	}
}
