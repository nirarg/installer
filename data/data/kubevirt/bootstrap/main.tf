data "ignition_file" "hostname" {
  mode  = "420"
  path  = "/etc/hostname"

  content {
    content = <<EOF
${var.cluster_id}-bootstrap
EOF
  }
}

data "ignition_config" "bootstrap_ignition_config" {

  merge {
    source = "data:text/plain;charset=utf-8;base64,${base64encode(var.ignition_data)}"
  }

  files = [
    element(data.ignition_file.hostname.*.rendered, 0)
  ]
}

resource "kubernetes_secret" "bootstrap_ignition" {
  metadata {
    name      = "${var.cluster_id}-bootstrap-ignition"
    namespace = var.namespace
    labels    = var.labels
  }
  data = {
    "userdata" = element(
      data.ignition_config.bootstrap_ignition_config.*.rendered,
      0,
    )
  }
}

resource "kubevirt_virtual_machine" "bootstrap_vm" {
  wait = true

  name                 = "${var.cluster_id}-bootstrap"
  namespace            = var.namespace
  labels               = var.labels
  storage_size         = var.storage
  memory               = var.memory
  cpu                  = var.cpu
  storage_class_name   = var.storage_class
  network_name         = var.network_name
  access_mode          = var.pv_access_mode
  ignition_secret_name = kubernetes_secret.bootstrap_ignition.metadata[0].name
  image_url            = var.image_url
  // pvc_name = var.pvc_name
}
