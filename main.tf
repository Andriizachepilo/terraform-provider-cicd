resource "cicd_example" "creating" {
  working_directory = "/Users/andriizachepilo/CopylotProject/terraform-provider-cicd/java"
  build             = "npm run build"
  test              = "npm test"
  # dockerfile_directory = "if not the same as working dir"
  docker_build = "andrey342/day4:hopew" // path to our dockerfile ?
  # docker_tag        = "mannually along with build? docker tag name:tag registryURL/repo/name:tag"
  # docker_tag1       = "automatically along with push ?"
  docker_push = "andrey342/day4:hopew"
  # container_registry_url = "acrukwestuniq.azurecr.io"
  container_registry_url      = "B"
  container_registry_password = "A"
}


// if build or test != "" ask working dir not to be emtpy
// auto tag +1 number etc ?

// 0) working_dir - wrong path
// 1) build -  no file found, dependencies are not installed, failed build ?
// 2) test -  no test file found, dependencies are not isntalled, failed test

// 3) docker_build - no dockerfile found, docker is not installed/not running (if dockerfile is in a different dir - add to the code if docker_dir != "" {cd into it})
// 4) registry - for docker_push we need creds for registry, export secrets, check if dependencies are installed like aws, azure etc, try to log into it, if wrong - 3 more attempts do not start from the beginning 
// 4) docker_push - image does not exist
// autotag ?