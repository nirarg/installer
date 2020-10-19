data "ignition_file" "hostname" {
  count = var.master_count
  mode  = "420"
  path  = "/etc/hostname"

  content {
    content = <<EOF
${var.cluster_id}-master-${count.index}
EOF
  }
}

data "ignition_config" "master_ignition_config" {
  count = var.master_count

  merge {
    source = "data:text/plain;charset=utf-8;base64,${base64encode(var.ignition_data)}"
  }

  files = [
    element(data.ignition_file.hostname.*.rendered, count.index)
  ]
}

resource "kubernetes_secret" "master_ignition" {
  count = var.master_count

  metadata {
    name      = "${var.cluster_id}-master-${count.index}-ignition"
    namespace = var.namespace
    labels    = var.labels
  }
  data = {
    "userdata" = element(
      data.ignition_config.master_ignition_config.*.rendered,
      count.index,
    )
  }
}

locals {
  anti_affinity_label = {
    "anti-affinity-tag-${var.cluster_id}" = "master"
  }
}

resource "kubevirt_virtual_machine" "master_vm" {
  wait = true

  count                      = var.master_count
  name                       = "${var.cluster_id}-master-${count.index}"
  namespace                  = var.namespace
  labels                     = merge(var.labels, local.anti_affinity_label)
  storage_size               = var.storage
  memory                     = var.memory
  cpu                        = var.cpu
  storage_class_name         = var.storage_class
  network_name               = var.network_name
  access_mode                = var.pv_access_mode
  ignition_secret_name       = kubernetes_secret.master_ignition[count.index].metadata[0].name
  image_url                  = var.image_url
  anti_affinity_match_labels = local.anti_affinity_label
  anti_affinity_topology_key = "kubernetes.io/hostname"
  // pvc_name = var.pvc_name
}
