resource "kubernetes_secret" "master_ignition" {
  provider = kubernetes

  metadata {
    name      = "master-ignition"
    namespace = var.namespace
  }
  data = {
    "userdata" = var.ignition_data
  }
}

resource "kubernetes_service" "master_service" {
  provider = kubernetes

  count = var.master_count

  metadata {
    name = "${var.cluster_id}-master-${count.index}"
    namespace = var.namespace
  }
  spec {
    cluster_ip = "None"
    selector = {
      name = "${var.cluster_id}-master-${count.index}"
    }
  }
}

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
    }
    "spec" = {
      "dataVolumeTemplates" = [
        {
          "metadata" = {
            "name" = "${var.cluster_id}-master-${count.index}"
          }
          "spec" = {
            "pvc" = {
              "storageClassName" = var.storage_class
              "accessModes" = [
                "ReadWriteOnce",
              ]
              "resources" = {
                "requests" = {
                  "storage" = var.master_storage
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
            }
            "machine" = {
              "type" = ""
            }
            "resources" = {
              "requests" = {
                "cpu" = var.master_cpu
                "memory" = var.master_memory
              }
            }
          }
          "terminationGracePeriodSeconds" = 0
          "volumes" = [
            {
              "dataVolume" = {
                "name" = "${var.cluster_id}-master-${count.index}"
              }
              "name" = "datavolumedisk1"
            },
            {
              "cloudInitConfigDrive" = {
                "secretRef" = {
                  "name" = "master-ignition"
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
