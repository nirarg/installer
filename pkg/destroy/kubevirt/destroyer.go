package kubevirt

import (
	"github.com/sirupsen/logrus"

	ickubevirt "github.com/openshift/installer/pkg/asset/installconfig/kubevirt"
	"github.com/openshift/installer/pkg/destroy/providers"
	"github.com/openshift/installer/pkg/types"
)

// ClusterUninstaller holds the various options for the cluster we want to delete.
type ClusterUninstaller struct {
	Metadata types.ClusterMetadata
	Logger   logrus.FieldLogger
}

// Run is the entrypoint to start the uninstall process.
func (uninstaller *ClusterUninstaller) Run() error {
	namespace := uninstaller.Metadata.Kubevirt.Namespace
	labels := uninstaller.Metadata.Kubevirt.Labels

	kubevirtClient, err := ickubevirt.NewClient()
	if err != nil {
		return err
	}
	if err := uninstaller.deleteAllVMs(namespace, labels, kubevirtClient); err != nil {
		return err
	}
	if err := uninstaller.deleteAllDVs(namespace, labels, kubevirtClient); err != nil {
		return err
	}
	if err := uninstaller.deleteAllSecrets(namespace, labels, kubevirtClient); err != nil {
		return err
	}
	return nil
}

func (uninstaller *ClusterUninstaller) deleteAllVMs(namespace string, labels map[string]string, kubevirtClient ickubevirt.Client) error {
	vmList, err := kubevirtClient.ListVirtualMachine(namespace)
	if err != nil {
		return err
	}

	var destroyList []string
	for _, d := range vmList.Items {
		existLabels := d.GetLabels()
		for k, v := range labels {
			if existVal, ok := existLabels[k]; ok && existVal == v {
				destroyList = append(destroyList, d.GetName())
				break
			}
		}
	}

	uninstaller.Logger.Infof("List tenant cluster's VMs (in namespace %s) return: %s", namespace, destroyList)
	for _, vmName := range destroyList {
		uninstaller.Logger.Infof("Delete VM %s", vmName)
		if err := kubevirtClient.DeleteVirtualMachine(namespace, vmName); err != nil {
			return err
		}
		ickubevirt.WaitForDeletionComplete(vmName, func() error {
			_, err := kubevirtClient.GetVirtualMachine(namespace, vmName)
			return err
		})
	}
	return nil
}

func (uninstaller *ClusterUninstaller) deleteAllDVs(namespace string, labels map[string]string, kubevirtClient ickubevirt.Client) error {
	dvList, err := kubevirtClient.ListDataVolume(namespace)
	if err != nil {
		return err
	}

	var destroyList []string
	for _, d := range dvList.Items {
		existLabels := d.GetLabels()
		for k, v := range labels {
			if existVal, ok := existLabels[k]; ok && existVal == v {
				destroyList = append(destroyList, d.GetName())
				break
			}
		}
	}

	uninstaller.Logger.Infof("List tenant cluster's DVs (in namespace %s) return: %s", namespace, destroyList)
	for _, dvName := range destroyList {
		uninstaller.Logger.Infof("Delete DV %s", dvName)
		if err := kubevirtClient.DeleteDataVolume(namespace, dvName); err != nil {
			return err
		}
		ickubevirt.WaitForDeletionComplete(dvName, func() error {
			_, err := kubevirtClient.GetDataVolume(namespace, dvName)
			return err
		})
	}
	return nil
}

func (uninstaller *ClusterUninstaller) deleteAllSecrets(namespace string, labels map[string]string, kubevirtClient ickubevirt.Client) error {
	secretList, err := kubevirtClient.ListSecret(namespace)
	if err != nil {
		return err
	}

	var destroyList []string
	for _, d := range secretList.Items {
		existLabels := d.GetLabels()
		for k, v := range labels {
			if existVal, ok := existLabels[k]; ok && existVal == v {
				destroyList = append(destroyList, d.GetName())
				break
			}
		}
	}
	uninstaller.Logger.Infof("List tenant cluster's secrets (in namespace %s) return: %s", namespace, destroyList)
	for _, secretName := range destroyList {
		uninstaller.Logger.Infof("Delete secret %s", secretName)
		if err := kubevirtClient.DeleteSecret(namespace, secretName); err != nil {
			return err
		}
		ickubevirt.WaitForDeletionComplete(secretName, func() error {
			_, err := kubevirtClient.GetSecret(namespace, secretName)
			return err
		})
	}
	return nil
}

// New returns oVirt Uninstaller from ClusterMetadata.
func New(logger logrus.FieldLogger, metadata *types.ClusterMetadata) (providers.Destroyer, error) {
	return &ClusterUninstaller{
		Metadata: *metadata,
		Logger:   logger,
	}, nil
}
