variable "cluster_id" {
  description = "The ID of Openshift cluster"
}

variable "namespace" {
  type        = string
  description = "The namespace in the underkube to install the overkube cluster within"
}

variable "storage_size" {
  type        = string
  description = ""
}

variable "image_url" {
  type        = string
  description = ""
}

variable "memory_size" {
  type        = string
  description = ""
}

variable "cpu_size" {
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
