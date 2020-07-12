variable "cluster_id" {
  description = "The ID of Openshift cluster"
}

variable "master_count" {
  description = "The number of master vm instances"
}

variable "namespace" {
  type        = string
  description = "The namespace in the underkube to install the overkube cluster within"
}

variable "master_storage" {
  type        = string
  description = ""
}

variable "image_url" {
  type        = string
  description = ""
}

variable "master_memory" {
  type        = string
  description = ""
}

variable "master_cpu" {
  type        = string
  description = ""
}

variable "ignition_data" {
  type        = string
  description = ""
}

variable "storage_class" {
  type        = string
  description = ""
}
