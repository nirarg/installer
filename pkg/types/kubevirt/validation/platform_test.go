package validation

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/kubevirt"
)

var (
	existingFilePath string
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	existingFilePath = filepath.Join(path, "test-file")
	if _, err := os.Create(existingFilePath); err != nil {
		log.Fatal(err)
	}
}

func shutdown() {
	if err := os.Remove(existingFilePath); err != nil {
		log.Fatal(err)
	}
}

func validPlatform() *kubevirt.Platform {
	// TODO <nargaman> find a way to use existing file for test environments
	return &kubevirt.Platform{
		InfraClusterKubeConfig:     existingFilePath,
		InfraClusterNamespace:      "test-namespace",
		InfraClusterStorageClass:   "",
		NetworkName:                "test network",
		APIVIP:                     "10.0.0.1",
		IngressVIP:                 "10.0.0.3",
		PersistentVolumeAccessMode: "ReadWriteMany",
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
			name: "empty network name",
			platform: func() *kubevirt.Platform {
				p := validPlatform()
				p.NetworkName = ""
				return p
			}(),
			valid: false,
		},
		{
			name: "empty API VIP",
			platform: func() *kubevirt.Platform {
				p := validPlatform()
				p.APIVIP = ""
				return p
			}(),
			valid: false,
		},
		{
			name: "invalid API VIP",
			platform: func() *kubevirt.Platform {
				p := validPlatform()
				p.APIVIP = "invalid API VIP"
				return p
			}(),
			valid: false,
		},
		{
			name: "empty ingress VIP",
			platform: func() *kubevirt.Platform {
				p := validPlatform()
				p.IngressVIP = ""
				return p
			}(),
			valid: false,
		},
		{
			name: "invalid ingress VIP",
			platform: func() *kubevirt.Platform {
				p := validPlatform()
				p.IngressVIP = "invalid ingress VIP"
				return p
			}(),
			valid: false,
		},
		{
			name: "empty access mode",
			platform: func() *kubevirt.Platform {
				p := validPlatform()
				p.PersistentVolumeAccessMode = ""
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
