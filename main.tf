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
  step_1 = "/Users/andriizachepilo/.nvm/versions/node/v21.5.0/bin/npm test"
  path   = "/Users/andriizachepilo/CopylotProject/terraform-provider-cicd/java"
}