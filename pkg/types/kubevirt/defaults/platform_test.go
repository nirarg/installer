package defaults

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/kubevirt"
)

func defaultInstallConfig() *types.InstallConfig {
	ic := &types.InstallConfig{
		ControlPlane: &types.MachinePool{
			Name: "master",
		},
		Compute: []types.MachinePool{
			{
				Name: "worker",
			},
		},
	}
	return ic
}

func expectedInstallConfig() *types.InstallConfig {
	ic := defaultInstallConfig()
	ic.ControlPlane.Platform.Kubevirt = &kubevirt.MachinePool{
		CPU:         8,
		Memory:      "16G",
		StorageSize: "35Gi",
	}
	ic.Compute[0].Platform = types.MachinePoolPlatform{
		Kubevirt: &kubevirt.MachinePool{
			CPU:         4,
			Memory:      "10G",
			StorageSize: "35Gi",
		},
	}
	return ic
}

func TestSetPlatformDefaults(t *testing.T) {
	cases := []struct {
		name     string
		ic       *types.InstallConfig
		expected *types.InstallConfig
	}{
		{
			name:     "happy_fllow",
			ic:       defaultInstallConfig(),
			expected: expectedInstallConfig(),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			SetPlatformDefaults(tc.ic.Platform.Kubevirt, tc.ic.ControlPlane, tc.ic.Compute)
			assert.Equal(t, tc.expected, tc.ic, "unexpected InstallConfig")
		})
	}
}
