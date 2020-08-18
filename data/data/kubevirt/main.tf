provider "kubernetes" {
  config_path = "kubeconfig"
}

provider "kubernetes-alpha" {
  server_side_planning = false
  config_path = "kubeconfig"
}

resource "kubernetes_manifest" "datavolume" {
  provider = kubernetes-alpha

  open_api_path = "io.kubevirt.cdi.v1alpha1.DataVolume"

  manifest = {
    "apiVersion" = "cdi.kubevirt.io/v1alpha1"
    "kind" = "DataVolume"
    "metadata" = {
      "name" = var.kubevirt_source_pvc_name
      "namespace" = var.kubevirt_namespace
      "labels" = var.kubevirt_labels
    }
    "spec" = {
      "pvc" = {
        "storageClassName" = var.kubevirt_storage_class
        "accessModes" = [
          var.kubevirt_pv_access_mode,
        ]
        "resources" = {
          "requests" = {
            "storage" = "20Gi"
          }
        }
      }
      "source" = {
        "http" = {
          "url" = var.kubevirt_image_url
        }
      }
    }
  }
}

module "masters" {
  source         = "./master"
  master_count   = var.master_count
  cluster_id     = var.cluster_id
  ignition_data  = var.ignition_master
  namespace      = var.kubevirt_namespace
  storage        = var.kubevirt_master_storage
  memory         = var.kubevirt_master_memory
  cpu            = var.kubevirt_master_cpu
  image_url      = var.kubevirt_image_url
  storage_class  = var.kubevirt_storage_class
  network_name   = var.kubevirt_network_name
  pv_access_mode = var.kubevirt_pv_access_mode
  labels         = var.kubevirt_labels
}

module "bootstrap" {
  source         = "./bootstrap"
  cluster_id     = var.cluster_id
  ignition_data  = var.ignition_bootstrap
  namespace      = var.kubevirt_namespace
  storage        = "35Gi"
  memory         = "4G"
  cpu            = "2"
  image_url      = var.kubevirt_image_url
  storage_class  = var.kubevirt_storage_class
  network_name   = var.kubevirt_network_name
  pv_access_mode = var.kubevirt_pv_access_mode
  labels         = var.kubevirt_labels
}
