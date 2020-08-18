variable "cluster_id" {
  description = "The ID of Openshift cluster"
}

variable "namespace" {
  type        = string
  description = "The namespace/project in the infracluster which all the tenantcluster resources should be created in"
}

variable "storage" {
  type        = string
  description = "bootstrap VM disk size, of type Quantity (see: https://github.com/kubernetes/apimachinery/blob/master/pkg/api/resource/quantity.go)"
}

variable "image_url" {
  type        = string
  description = "TODO nargaman remove (OCPCNV-100)"
}

variable "memory" {
  type        = string
  description = "Bootstrap VM memory size, of type Quantity (see: https://github.com/kubernetes/apimachinery/blob/master/pkg/api/resource/quantity.go)"
}

variable "cpu" {
  type        = string
  description = "Bootstrap VM number of cores"
}

variable "ignition_data" {
  type        = string
  description = "Ignition config file contents of the bootstrap VM"
}

variable "storage_class" {
  type        = string
  description = "The \"class\" of the storage located in the infracluster"
}

variable "network_name" {
  type        = string
  description = "The name of the sub network created in the infracluster which should be used by the tenantcluster resources"
}

variable "pv_access_mode" {
  type        = string
  description = "The access mode which all the persistant volumes should be created with [ReadWriteOnce,ReadOnlyMany,ReadWriteMany]"
}

variable "labels" {
  type = map(string)

  description = <<EOF
(optional) Labels to be applied to created resources.

Example: `{ "key" = "value", "foo" = "bar" }`
EOF

  default = {}
}
