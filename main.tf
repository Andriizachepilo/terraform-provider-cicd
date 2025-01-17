terraform {
  required_providers {
    cicd = {
      source  = "local/cicd"
      version = "0.1.0"
    }
  }
}

provider "cicd" {
  # Configuration options
}

resource "cicd_example" "example" {
  programming_language     = "java"
  build_tool               = "maven"
  docker_tag               = "myapp:123"
  registry_url             = "myregistry.com"
  registry_credentials     = "MY_REGISTRY_PASSWORD"

  skip_build               = false
  skip_test                = false
  skip_build_docker_image  = false
  skip_push_docker_image   = false

  custom_build_command     = null
  custom_test_command      = null

  environment_variables = {
    KEY1 = "value1"
    KEY2 = "value2"
  }
}