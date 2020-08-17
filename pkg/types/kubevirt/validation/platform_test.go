package validation

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/kubevirt"
)

func validPlatform() *kubevirt.Platform {
	// TODO <nargaman> find a way to use existing file for test environments
	existsFile := filepath.Join(os.Getenv("HOME"), "/.kube/config")
	return &kubevirt.Platform{
		InfraClusterIngressDomain: "localhost",
		InfraClusterKubeConfig:    existsFile,
		InfraClusterNamespace:     "test-namespace",
		InfraClusterStorageClass:  "",
		SourcePvcName:             "test-pvc",
		DefaultMachinePlatform:    nil,
	}
}

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name     string
		platform *kubevirt.Platform
		valid    bool
	}{
		{
			name:     "valid",
			platform: validPlatform(),
			valid:    true,
		},
		{
			name: "empty ingress domain",
			platform: func() *kubevirt.Platform {
				p := validPlatform()
				p.InfraClusterIngressDomain = ""
				return p
			}(),
			valid: false,
		},
		{
			name: "invalid ingress domain",
			platform: func() *kubevirt.Platform {
				p := validPlatform()
				p.InfraClusterIngressDomain = "invalid ingress domain"
				return p
			}(),
			valid: false,
		},
		{
			name: "empty kubeconfig",
			platform: func() *kubevirt.Platform {
				p := validPlatform()
				p.InfraClusterKubeConfig = ""
				return p
			}(),
			valid: false,
		},
		{
			name: "invalid kubeconfig",
			platform: func() *kubevirt.Platform {
				p := validPlatform()
				p.InfraClusterKubeConfig = "invalid ingress domain"
				return p
			}(),
			valid: false,
		},
		{
			name: "empty namespace",
			platform: func() *kubevirt.Platform {
				p := validPlatform()
				p.InfraClusterNamespace = ""
				return p
			}(),
			valid: false,
		},
		{
			name: "empty source pvc",
			platform: func() *kubevirt.Platform {
				p := validPlatform()
				p.SourcePvcName = ""
				return p
			}(),
			valid: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePlatform(tc.platform, field.NewPath("test-path")).ToAggregate()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
