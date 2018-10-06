variable "environment_base_name" {
  default = "{{.PackageName}}"
}
variable "resource_group" {
  default = "{{.PackageName}}-rg" // by default, resource_group is the environment_base_name + -rg
}

variable "location" {
  default = "japaneast"
}

variable "tag_name" {
  default = "{{.PackageName}}"
}

variable "language" {
    default = "dotnet"
}