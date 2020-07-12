variable "kubevirt_namespace" {
  type        = string
  description = "The namespace in the underkube to install the overkube cluster within"
}

variable "kubevirt_source_pvc_name" {
  type        = string
  description = ""
}

variable "kubevirt_master_storage" {
  type        = string
  description = ""
}

variable "kubevirt_image_url" {
  type        = string
  description = ""
}

variable "kubevirt_master_memory" {
  type        = string
  description = ""
}

variable "kubevirt_master_cpu" {
  type        = string
  description = ""
}

variable "kubevirt_storage_class" {
  type        = string
  description = ""
}
