resource "kubernetes_secret" "master_ignition" {
  count = var.master_count

  metadata {
    name      = "${var.cluster_id}-master-${count.index}-ignition"
    namespace = var.namespace
    labels    = var.labels
  }
  data = {
    "userdata" = var.ignition_data
  }
}

resource "kubevirt_virtual_machine" "master_vm" {
  wait = true

  count                = var.master_count
  name                 = "${var.cluster_id}-master-${count.index}"
  namespace            = var.namespace
  labels               = var.labels
  storage_size         = var.storage
  memory               = var.memory
  cpu                  = var.cpu
  storage_class_name   = var.storage_class
  network_name         = var.network_name
  access_mode          = var.pv_access_mode
  ignition_secret_name = kubernetes_secret.master_ignition[count.index].metadata[0].name
  image_url            = var.image_url
  // pvc_name = var.pvc_name
}
