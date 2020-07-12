provider "kubernetes" {
  config_path = "/home/nargaman/.kube/config"
}

provider "kubernetes-alpha" {
  server_side_planning = false
  config_path = "/home/nargaman/.kube/config"
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
    }
    "spec" = {
      "pvc" = {
        "storageClassName" = var.kubevirt_storage_class
        "accessModes" = [
          "ReadWriteMany",
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
  source          = "./master"
  master_count    = var.master_count
  cluster_id      = var.cluster_id
  ignition_data   = var.ignition_master
  namespace       = var.kubevirt_namespace
  master_storage  = var.kubevirt_master_storage
  master_memory   = var.kubevirt_master_memory
  master_cpu      = var.kubevirt_master_cpu
  image_url       = var.kubevirt_image_url
  storage_class   = var.kubevirt_storage_class
}

module "bootstrap" {
  source          = "./bootstrap"
  cluster_id      = var.cluster_id
  ignition_data   = var.ignition_bootstrap
  namespace       = var.kubevirt_namespace
  storage_size    = "35Gi"
  memory_size     = "4G"
  cpu_size        = "2"
  image_url       = var.kubevirt_image_url
  storage_class   = var.kubevirt_storage_class
}

resource "kubernetes_service" "api_service" {
  provider = kubernetes

  metadata {
    name = "api"
    namespace = var.kubevirt_namespace
  }
  spec {
    selector = {
      api-vm = ""
    }
    port {
      name = "api-server"
      port = 6443
      protocol = "TCP"
    }
  }
}

resource "kubernetes_manifest" "api_route" {
  provider = kubernetes-alpha

  open_api_path = "com.github.openshift.api.route.v1.Route"

  manifest = {
    "apiVersion" = "route.openshift.io/v1"
    "kind" = "Route"
    "metadata" = {
      "name" = "api"
      "namespace" = var.kubevirt_namespace
    }
    "spec" = {
      "path" = ""
      "to" = {
        "kind" = "Service"
        "name" = "api"
      }
      "port" = {
        "targetPort" = "api-server"
      }
      "tls" = {
        "termination" = "passthrough"
      }
    }
  }
}

resource "kubernetes_service" "api-int_service" {
  provider = kubernetes

  metadata {
    name = "api-int"
    namespace = var.kubevirt_namespace
  }
  spec {
    selector = {
      api-vm = ""
    }
    port {
      name = "api-server"
      port = 6443
      protocol = "TCP"
    }
    port {
      name = "machine-config"
      port = 2222
      protocol = "TCP"
    }
  }
}
