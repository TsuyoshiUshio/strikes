variable "environment_base_name" {
  default = "{{.PackageName}}"
}
variable "resource_group" {
  default = "{{.PackageName}}-rg" // by default, resource_group is the environment_base_name + -rg
}

variable "repository_base_uri" {
    default = "https://strikesrepoe9eej5x3.blob.core.windows.net/repository/"
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