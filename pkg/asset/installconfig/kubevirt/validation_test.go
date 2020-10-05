package kubevirt

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/asset/installconfig/kubevirt/mock"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/kubevirt"
)

var (
	validKubeconfigPath   = ""
	validNamespace        = "valid-namespace"
	validStorageClass     = "valid-storage-class"
	validNetworkName      = "valid-network-name"
	validAPIVIP           = "192.168.123.15"
	validIngressVIP       = "192.168.123.20"
	validAccessMode       = "valid-access-mode"
	validMachineCIDR      = "192.168.123.0/24"
	invalidKubeconfigPath = "invalid-kubeconfig-path"
	invalidNamespace      = "invalid-namespace"
	invalidStorageClass   = "invalid-storage-class"
	invalidNetworkName    = "invalid-network-name"
	invalidAPIVIP         = "invalid-api-vip"
	invalidIngressVIP     = "invalid-ingress-vip"
	invalidAccessMode     = "invalid-access-mode"
	invalidMachineCIDR    = "10.0.0.0/16"
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
	validKubeconfigPath = filepath.Join(path, "test-file")
	if _, err := os.Create(validKubeconfigPath); err != nil {
		log.Fatal(err)
	}
}

func shutdown() {
	if err := os.Remove(validKubeconfigPath); err != nil {
		log.Fatal(err)
	}
}

func validInstallConfig() *types.InstallConfig {
	return &types.InstallConfig{
		Networking: &types.Networking{
			MachineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR(validMachineCIDR)},
			},
		},
		Platform: types.Platform{
			Kubevirt: &kubevirt.Platform{
				InfraClusterKubeConfig:     validKubeconfigPath,
				InfraClusterNamespace:      validNamespace,
				InfraClusterStorageClass:   validStorageClass,
				NetworkName:                validNetworkName,
				APIVIP:                     validAPIVIP,
				IngressVIP:                 validIngressVIP,
				PersistentVolumeAccessMode: validAccessMode,
			},
		},
	}
}

func TestKubevirtInstallConfigValidation(t *testing.T) {
	cases := []struct {
		name             string
		edit             func(ic *types.InstallConfig)
		expectedError    bool
		expectedErrMsg   string
		clientBuilderErr error
	}{
		{
			name:           "valid",
			edit:           nil,
			expectedError:  false,
			expectedErrMsg: "",
		},
		{
			name:           "invalid empty platform",
			edit:           func(ic *types.InstallConfig) { ic.Platform.Kubevirt = nil },
			expectedError:  true,
			expectedErrMsg: "platform.kubevirt: Required value: validation requires a Engine platform configuration",
		},
		{
			name:           "invalid empty kubeconfig",
			edit:           func(ic *types.InstallConfig) { ic.Platform.Kubevirt.InfraClusterKubeConfig = "" },
			expectedError:  true,
			expectedErrMsg: "platform.kubevirt.InfraClusterKubeConfig: Invalid value: \"\": stat : no such file or directory",
		},
		{
			name:             "invalid client builder error",
			expectedError:    true,
			expectedErrMsg:   fmt.Sprintf("platform.kubevirt.InfraClusterReachable: Invalid value: \"%s\": failed to create InfraCluster client with error: test", validKubeconfigPath),
			clientBuilderErr: errors.New("test"),
		},
		{
			name:           "invalid namespace",
			edit:           func(ic *types.InstallConfig) { ic.Platform.Kubevirt.InfraClusterNamespace = invalidNamespace },
			expectedError:  true,
			expectedErrMsg: "platform.kubevirt.NamespaceExistsInInfraCluster: Invalid value: \"invalid-namespace\": failed to get namespace invalid-namespace from InfraCluster, with error: test",
		},
		{
			name:           "invalid storage class",
			edit:           func(ic *types.InstallConfig) { ic.Platform.Kubevirt.InfraClusterStorageClass = invalidStorageClass },
			expectedError:  true,
			expectedErrMsg: "platform.kubevirt.StorageClassExistsInInfraCluster: Invalid value: \"invalid-storage-class\": failed to get storageClass invalid-storage-class from InfraCluster, with error: test",
		},
		{
			name:           "invalid network name",
			edit:           func(ic *types.InstallConfig) { ic.Platform.Kubevirt.NetworkName = invalidNetworkName },
			expectedError:  true,
			expectedErrMsg: "platform.kubevirt.NetworkAttachmentDefinitionExistsInInfraCluster: Invalid value: \"invalid-network-name\": failed to get network-attachment-definition invalid-network-name from InfraCluster, with error: test",
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			installConfig := validInstallConfig()
			if tc.edit != nil {
				tc.edit(installConfig)
			}

			kubevirtClient := mock.NewMockClient(mockCtrl)
			if installConfig.Platform.Kubevirt != nil {
				kubevirtClient.EXPECT().ListNamespace(gomock.Any()).Return(nil, nil).AnyTimes()
				if installConfig.Platform.Kubevirt.InfraClusterNamespace == validNamespace {
					kubevirtClient.EXPECT().GetNamespace(gomock.Any(), validNamespace).Return(nil, nil).AnyTimes()
					if installConfig.Platform.Kubevirt.NetworkName == validNetworkName {
						kubevirtClient.EXPECT().GetNetworkAttachmentDefinition(gomock.Any(), validNetworkName, validNamespace).Return(nil, nil).AnyTimes()
					} else {
						kubevirtClient.EXPECT().GetNetworkAttachmentDefinition(gomock.Any(), invalidNetworkName, validNamespace).Return(nil, errors.New("test")).AnyTimes()
					}
				} else {
					kubevirtClient.EXPECT().GetNamespace(gomock.Any(), invalidNamespace).Return(nil, errors.New("test")).AnyTimes()
				}
				if installConfig.Platform.Kubevirt.InfraClusterStorageClass == validStorageClass {
					kubevirtClient.EXPECT().GetStorageClass(gomock.Any(), validStorageClass).Return(nil, nil).AnyTimes()
				} else {
					kubevirtClient.EXPECT().GetStorageClass(gomock.Any(), invalidStorageClass).Return(nil, errors.New("test")).AnyTimes()
				}
			}

			errs := Validate(installConfig, func(kubeconfig string) (Client, error) { return kubevirtClient, tc.clientBuilderErr })
			if tc.expectedError {
				assert.Regexp(t, tc.expectedErrMsg, errs)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}
