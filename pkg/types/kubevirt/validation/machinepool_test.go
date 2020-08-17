package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/kubevirt"
)

func TestValidateMachinePool(t *testing.T) {
	cases := []struct {
		name  string
		pool  *kubevirt.MachinePool
		valid bool
	}{
		{
			name: "valid",
			pool: &kubevirt.MachinePool{
				CPU:         4,
				Memory:      "5G",
				StorageSize: "100Gi",
			},
			valid: true,
		},
		{
			name: "invalid cpu",
			pool: &kubevirt.MachinePool{
				CPU:         0,
				Memory:      "5G",
				StorageSize: "100Gi",
			},
			valid: false,
		},
		{
			name: "empty memory",
			pool: &kubevirt.MachinePool{
				CPU:         4,
				Memory:      "5G",
				StorageSize: "",
			},
			valid: false,
		},
		{
			name: "invalid memory",
			pool: &kubevirt.MachinePool{
				CPU:         4,
				Memory:      "5G",
				StorageSize: "invalid string",
			},
			valid: false,
		},
		{
			name: "empty storageSize",
			pool: &kubevirt.MachinePool{
				CPU:         4,
				Memory:      "",
				StorageSize: "100Gi",
			},
			valid: false,
		},
		{
			name: "invalid memory",
			pool: &kubevirt.MachinePool{
				CPU:         4,
				Memory:      "invalid string",
				StorageSize: "100Gi",
			},
			valid: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateMachinePool(tc.pool, field.NewPath("test-path")).ToAggregate()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
