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

resource "cicd_example" "creating" {
  path   = "/Users/andriizachepilo/CopylotProject/terraform-provider-cicd/java"
  step_1 = "npm test"
}