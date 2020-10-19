resource "kubevirt_data_volume" "datavolume" {
  name               = var.pvc_name
  namespace          = var.namespace
  labels             = var.labels
  storage_size       = var.storage
  access_mode        = var.pv_access_mode
  storage_class_name = var.storage_class
  image_url          = var.image_url
}
