resource "kubernetes_secret" "bootstrap_ignition" {
  provider = kubernetes

  metadata {
    name      = "bootstrap-ignition"
    namespace = var.namespace
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
    }
    "spec" = {
      "dataVolumeTemplates" = [
        {
          "metadata" = {
            "name" = "${var.cluster_id}-bootstrap"
          }
          "spec" = {
            "pvc" = {
              "storageClassName" = var.storage_class
              "accessModes" = [
                "ReadWriteOnce",
              ]
              "resources" = {
                "requests" = {
                  "storage" = var.storage_size
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
            }
            "machine" = {
              "type" = ""
            }
            "resources" = {
              "requests" = {
                "cpu" = var.cpu_size
                "memory" = var.memory_size
              }
            }
          }
          "terminationGracePeriodSeconds" = 0
          "volumes" = [
            {
              "dataVolume" = {
                "name" = "${var.cluster_id}-bootstrap"
              }
              "name" = "datavolumedisk1"
            },
            {
              "cloudInitConfigDrive" = {
                "secretRef" = {
                  "name" = "bootstrap-ignition"
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
