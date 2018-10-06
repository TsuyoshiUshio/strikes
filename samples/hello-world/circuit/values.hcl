variable "environment_base_name" {
  default = "hello-world"
}
variable "resource_group" {
  default = "hello-world-rg" // by default, resource_group is the environment_base_name + -rg
}

variable "location" {
  default = "japaneast"
}

variable "tag_name" {
  default = "hello-world"
}

variable "language" {
    default = "dotnet"
}

variable "packages_sub_dir" {
    default = "hello-world/1.0.0/hello.zip"
}