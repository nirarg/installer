resource "kubernetes_secret" "bootstrap_ignition" {
  provider = kubernetes

  metadata {
    name      = "${var.cluster_id}-bootstrap-ignition"
    namespace = var.namespace
    labels    = var.labels
  }
  data = {
    "userdata" = var.ignition_data
  }
}

resource "kubernetes_manifest" "bootstrap_vm" {
  provider = kubernetes-alpha

  open_api_path = "io.kubevirt.v1alpha3.VirtualMachine"

  manifest = {
    "apiVersion" = "kubevirt.io/v1alpha3"
    "kind" = "VirtualMachine"
    "metadata" = {
      "name" = "${var.cluster_id}-bootstrap"
      "namespace" = var.namespace
      "labels" = var.labels
    }
    "spec" = {
      "dataVolumeTemplates" = [
        {
          "metadata" = {
            "name" = "${var.cluster_id}-bootstrap-bootvolume"
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
              "name" = "${var.cluster_id}-bootstrap"
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
                "name" = "${var.cluster_id}-bootstrap-bootvolume"
              }
              "name" = "datavolumedisk1"
            },
            {
              "cloudInitConfigDrive" = {
                "secretRef" = {
                  "name" = "${var.cluster_id}-bootstrap-ignition"
                }
              }
              "name" = "cloudinitdisk"
            },
          ]
        }
      }
    }
  }
}
