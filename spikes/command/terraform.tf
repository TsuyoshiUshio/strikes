# Service Principle
variable "subscription_id" {}
variable "client_id" {}
variable "client_secret" {}
variable "tenant_id" {}

# Resource Group
variable "resource_group" {
  default = "hello-world-rg"
}
variable "location" {
  default = "japaneast"
}
variable "tag_name" {
  default = "hello-world"
}

variable "environment_base_name" {
  default = "hello-world"
}

# Repository Settings

variable "repository_base_uri" {
    default = "https://asset.simplearchitect.club/"
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

variable "package_zip_name" {
    # example
    # default = "hello.zip"
}

# Specify the language of the package
variable "language" {
    default = "dotnet"
}



resource "azurerm_resource_group" "test" {
  name     = "${var.resource_group}"
  location = "${var.location}"
  tags {
    name = "${var.tag_name}"
  }
}

resource "random_string" "suffix" {
  length = 8
  special = false
  upper = false
}

resource "azurerm_storage_account" "test" {
  name                     = "${var.environment_base_name}sa${random_string.suffix.result}"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
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
}

resource "azurerm_application_insights" "test" {
  name                = "${var.environment_base_name}ai"
  location            = "eastus"
  resource_group_name = "${azurerm_resource_group.test.name}"
  application_type    = "Web"
}

# CosmosDB

resource "random_integer" "ri" {
    min = 10000
    max = 99999
}

resource "azurerm_function_app" "test" {
  name                      = "${var.environment_base_name}app"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  app_settings {
    "APPINSIGHTS_INSTRUMENTATIONKEY" = "${azurerm_application_insights.test.instrumentation_key}"
    "FUNCTIONS_EXTENSION_VERSION" = "beta"
    "FUNCTIONS_WORKER_RUNTIME" = "${var.language}"
    # This going to change into WEBSITE_RUN_FROM_PACKAGE
    "WEBSITE_RUN_FROM_ZIP" = "${var.repository_base_uri}${var.package_name}/${var.package_zip_name}/${var.package_zip_name}"
  }
}
