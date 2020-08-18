resource "kubernetes_manifest" "master_vm" {
  provider = kubernetes-alpha

  open_api_path = "io.kubevirt.v1alpha3.VirtualMachine"
  count = var.master_count

  manifest = {
    "apiVersion" = "kubevirt.io/v1alpha3"
    "kind" = "VirtualMachine"
    "metadata" = {
      "name" = "${var.cluster_id}-master-${count.index}"
      "namespace" = var.namespace
      "labels" = var.labels
    }
    "spec" = {
      "dataVolumeTemplates" = [
        {
          "metadata" = {
            "name" = "${var.cluster_id}-master-${count.index}-bootvolume"
          }
          "spec" = {
            "pvc" = {
              "storageClassName" = var.storage_class
              "accessModes" = [
          	var.pv_access_mode,
              ]
              "resources" = {
                "requests" = {
                  "storage" = var.storage
                }
              }
            }
	    "source" = {
	      "http" = {
	        "url" = var.image_url
	      }
	    }
          }
          "status" = {}
        },
      ]
      "running" = true
      "template" = {
        "metadata" = {
          "labels" = {
              "api-vm" = ""
              "name" = "${var.cluster_id}-master-${count.index}"
	  }
        }
        "spec" = {
          "domain" = {
            "devices" = {
              "disks" = [
                {
                  "disk" = {
                    "bus" = "virtio"
                  }
                  "name" = "datavolumedisk1"
                },
                {
                  "disk" = {
                    "bus" = "virtio"
                  }
                  "name" = "cloudinitdisk"
                },
              ]
              "interfaces" = [
                {
                  "name" = "main"
                  "bridge" = {}
                },
		{
		  "name" = "default"
		  "masquerade" = {}
		},
              ]
            }
            "machine" = {
              "type" = ""
            }
            "resources" = {
              "requests" = {
                "cpu" = var.cpu
                "memory" = var.memory
              }
            }
          }
          "terminationGracePeriodSeconds" = 0
	  "networks" = [
            {
              "multus" = {
                "networkName" = var.network_name
              }
              "name" = "main"
            },
	    {
	      "pod" = {}
  	      "name"= "default"
	    },
          ]
          "volumes" = [
            {
              "dataVolume" = {
                "name" = "${var.cluster_id}-master-${count.index}-bootvolume"
              }
              "name" = "datavolumedisk1"
            },
            {
              "cloudInitConfigDrive" = {
                "userData" = var.ignition_data
              }
              "name" = "cloudinitdisk"
            },
          ]
        }
      }
    }
  }
}
