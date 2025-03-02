resource "cicd_example" "creating" {
  working_directory = "/Users/andriizachepilo/CopylotProject/terraform-provider-cicd/java"
  build             = "npm run build"
  test              = "npm test"
  # dockerfile_directory = "/Users/andriizachepilo/CopylotProject/terraform-provider-cicd/java"
  # docker_build = "docker build -t 18:feb ." 
  container_registry_url = "docker.io/andrey342"
  # docker_tag        = "mannually along with build? docker tag name:tag registryURL/repo/name:tag"
  # docker_tag1       = "automatically along with push ?"
  # docker_push = "andrey342/day4:how"
}


// if build or test != "" ask working dir not to be emtpy
// auto tag +1 number etc ?

// 0) working_dir - wrong path
// 1) build -  no file found, dependencies are not installed, failed build ?
// 2) test -  no file found, dependencies are not isntalled, failed test

// 3) docker_build - no dockerfile found, docker is not installed/not running (if dockerfile is in a different dir - add to the code if docker_dir != "" {cd into it})
// 4) registry - for docker_push we need creds for registry, export secrets, check if dependencies are installed like aws, azure etc, try to log into it, if wrong - 3 more attempts do not start from the beginning/ private and public registry ?
// 4) docker_push - image does not exist, delete all images after pushing ?
// autotag ?

//rename feedback function, add error handling and other structure of my project, get rid of dots


// authentication probles = docker is not running, token expired, command not found, incorrect url for security 
// malicious handling, regex for injections for each step ?
// install dependencies 
// token for cr's ?
// check if acr exists ?