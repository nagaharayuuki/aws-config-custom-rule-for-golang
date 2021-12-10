variable "resource_name_prefix" {
  description = "All the resources will be prefixed with this varible"
  default     = "nagahara-test"
}

variable "exclude_buckets" {
  description = "Skip compliance check for listed buckets"
  default     = ""
}

locals {
  config_rule_name = "${var.resource_name_prefix}-s3-lifecycle"
}
