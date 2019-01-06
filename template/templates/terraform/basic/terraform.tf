# Service Principle
variable "subscription_id" {}
variable "client_id" {}
variable "client_secret" {}
variable "tenant_id" {}

# Resource Group
variable "resource_group" {
  default = "{{.PackageName}}-rg"
}
variable "location" {
  default = "japaneast"
}
variable "tag_name" {
  default = "{{.PackageName}}"
}

variable "environment_base_name" {
  default = "{{.PackageName}}"
}

# Repository Settings

variable "repository_base_uri" {
    default = "https://strikesrepoe9eej5x3.blob.core.windows.net/repository/"
}

# Pacakge settings

variable "package_name" {
  # example
  # default = "hello-world"
}

variable "package_version" {
  # example
  # default = 1.0.0
}

variable "package_zip_name_0" {
    # example
    # default = "hello.zip"
}

# Specify the language of the package
variable "language" {
    default = "dotnet"
}

provider "azurerm" {
   client_id = "${var.client_id}"
   client_secret = "${var.client_secret}"
   subscription_id = "${var.subscription_id}"
   tenant_id = "${var.tenant_id}"
}

resource "azurerm_resource_group" "test" {
  name     = "${var.resource_group}"
  location = "${var.location}"
  tags {
    name = "${var.tag_name}"
  }
}

resource "random_string" "suffix" {
  length = 4
  special = false
  upper = false
}

resource "azurerm_storage_account" "test" {
  name                     = "${replace(var.environment_base_name,"-","")}sa${random_string.suffix.result}"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  depends_on = ["random_string.suffix", "azurerm_resource_group.test"]
}

resource "azurerm_app_service_plan" "test" {
  name                = "${var.environment_base_name}plan"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  kind                = "FunctionApp"

  sku {
    tier = "Dynamic"
    size = "Y1"
  }
    depends_on = ["azurerm_resource_group.test"]
}

resource "azurerm_application_insights" "test" {
  name                = "${var.environment_base_name}ai"
  location            = "eastus"
  resource_group_name = "${azurerm_resource_group.test.name}"
  application_type    = "Web"
  
  depends_on = ["azurerm_resource_group.test"]
}

resource "azurerm_function_app" "test" {
  name                      = "${var.environment_base_name}app"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"
  version = "~2"

  app_settings {
    "APPINSIGHTS_INSTRUMENTATIONKEY" = "${azurerm_application_insights.test.instrumentation_key}"
    "FUNCTIONS_WORKER_RUNTIME" = "${var.language}"
    # This going to change into WEBSITE_RUN_FROM_PACKAGE
    "WEBSITE_RUN_FROM_PACKAGE" = "${var.repository_base_uri}${var.package_name}/${var.package_version}/package/${var.package_zip_name_0}"
  }

  depends_on = ["azurerm_app_service_plan.test", "azurerm_resource_group.test", "azurerm_storage_account.test"] 
}
