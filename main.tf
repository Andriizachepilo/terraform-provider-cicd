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
  step_1 = "npm test" //check if mvn is installed, if it's not ask if u want to and if yes automatically install it 
  path = "/Users/andriizachepilo/CopylotProject/terraform-provider-cicd/java" //error like it's not found or smth else (looking for json or pom.xml)
 
}