package cicd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceCICD() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCICDCreate,
		ReadContext:   resourceCICDRead,
		UpdateContext: resourceCICDUpdate,
		DeleteContext: resourceCICDDelete,
		CustomizeDiff: customizeDiff,

		Schema: map[string]*schema.Schema{
			"build": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"test": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"working_directory": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"build_and_test": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"docker_build": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"docker_push": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"container_registry_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dockerfile_directory": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func customizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// Set a new timestamp to force the resource to be updated
	d.SetNewComputed("timestamp")
	return nil
}

func executeCommand(command string, workDir string) (error, string) {
	if workDir == "" {
		return fmt.Errorf("working diretory is not specified"), "."
	} else {
		err := os.Chdir(workDir)
		if err != nil {
			return fmt.Errorf("failed to change directory: %v", err), "."
		}
		output, err := exec.Command("sh", "-c", command).CombinedOutput()
		if err != nil {
			dependenciesErr := ""
			if strings.Contains(string(output), "command not found") {
				dependenciesErr = "* Make sure all necessary dependencies are installed *"
			}
			return fmt.Errorf("command failed with error: %v\nOutput: %s\n%s", err, string(output), dependenciesErr), "."
		}
		return nil, string(output)
	}

}

func resourceCICDCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	working_directory := d.Get("working_directory").(string)
	build := d.Get("build").(string)
	test := d.Get("test").(string)
	build_and_test := d.Get("build_and_test").(string)
	dockerfile_dir := d.Get("dockerfile_directory").(string)
	docker_build := d.Get("docker_build").(string)
	cr_url := d.Get("container_registry_url").(string)
	// docker_push := d.Get("docker_push").(string)

	feedback := func(processName, output string) diag.Diagnostics {
		return diag.Diagnostics{
			{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Step %v completed!\n%v", processName, output),
			},
		}
	}

	steps := []struct {
		name string
		cmd  string
	}{
		{"build", build},
		{"test", test},
		{"build_and_test", build_and_test},
		{"docker_build", docker_build},
	}

	for _, step := range steps {
		if step.cmd != "" {
			var err error
			var output string

			if dockerfile_dir != "" && step.name == "docker_build" {
				err, output = executeCommand(step.cmd, dockerfile_dir)
			} else {
				err, output = executeCommand(step.cmd, working_directory)
			}

			if err != nil {
				return diag.FromErr(err)
			}
			diags = append(diags, feedback(step.name, output)...)
		}
	}

	// // Execute the Docker build command if provided
	// if docker_build != "" {
	// 	if dockerfile_dir != "" {
	// 		err := os.Chdir(dockerfile_dir)
	// 		if err != nil {
	// 			return diag.FromErr(fmt.Errorf("failed to change directory: %v", err))
	// 		}
	// 	}
	// 	output, err := exec.Command("sh", "-c", fmt.Sprintf("docker build -t %v .", docker_build)).CombinedOutput()
	// 	if err != nil {
	// 		dockerErr := ""
	// 		if strings.Contains(string(output), "Is the docker daemon running?") {
	// 			dockerErr = "Docker is not running, please start docker and try again"
	// 		} else if strings.Contains(string(output), "Not found") {
	// 			dockerErr = "Docker is not installed, please install Docker and try again"
	// 		} else if strings.Contains(string(output), "failed to read dockerfile") {
	// 			dockerErr = "Dockerfile is not found"
	// 		}
	// 		return diag.FromErr(fmt.Errorf("docker build failed with error: %v\n\nOutput: %v\n\n%s", err, dockerErr, string(output)))
	//     } else {
	// 		return feedback(diags, "docker_build", string(output))
	// 	}
	// }

	// <aws_account_id>.dkr.ecr.<region>.amazonaws.com
	// 987654321098.dkr.ecr.eu-west-1.amazonaws.com

	//     urlFormat := func(aws string, azure string) (bool) {
	//      if aws != "" {
	// 		return regexp.MustCompile(`^\d{12}\.dkr\.ecr\.[a-z0-9-]+\.amazonaws\.com$`).MatchString(aws)
	//      } else if azure != "" {
	//         return regexp.MustCompile(`^[a-zA-Z0-9]+\.azurecr\.io$`).MatchString(azure)
	// 	 }
	// 	 return false
	// }

	if cr_url != "" {
		var input string
		if strings.Contains(cr_url, "amazonaws.com") {
			check_url_format := regexp.MustCompile(`^\d{12}\.dkr\.ecr\.[a-z0-9-]+\.amazonaws\.com$`).MatchString(cr_url)
			if check_url_format {
				err := exec.Command("aws", "sts", "get-caller-identity").Run()
				if err == nil {
					regionRegex := regexp.MustCompile(`(?:[^\.]*\.){3}([^\.]*)`).FindStringSubmatch(cr_url)[1]
					input = fmt.Sprintf("aws ecr get-login-password --region %s | docker login --username AWS --password-stdin %s", regionRegex, cr_url)
				} else {
					return diag.FromErr(fmt.Errorf("unable to locate credentials. You can configure credentials by running 'aws configure' "))
				}
			} else {
				return diag.FromErr(fmt.Errorf("invalid ECR URL format. Expected format: <aws_account_id>.dkr.ecr.<region>.amazonaws.com"))
			}
			//myproject.azurecr.io
		} else if strings.Contains(cr_url, "azurecr.io") {
			check_url_format := regexp.MustCompile(`^[a-zA-Z0-9]+\.azurecr\.io$`).MatchString(cr_url)
			if check_url_format {
				err := exec.Command("az", "account", "show").Run()
				if err == nil {
					cr_name := regexp.MustCompile(`[^.]+`).FindStringSubmatch(cr_url)[1]
					input = fmt.Sprintf("az acr login --name %s", cr_name)
				} else {
					return diag.FromErr(fmt.Errorf("unable to locate Azure credentials. Please log in using: az login"))
				}
			}
			//gcr.io/[PROJECT-ID]
		} else if strings.Contains(cr_url, "gcr.io") {
			check_url_format := regexp.MustCompile(`^gcr\.io/[a-zA-Z0-9-]+$`).MatchString(cr_url)
			if check_url_format {
				err := exec.Command("gcloud", "auth", "login").Run()
				if err == nil {
					gcrProjectID := regexp.MustCompile(`^[^/]+/(.+)$`).FindStringSubmatch(cr_url)[1]
					err := exec.Command("gcloud", "config", "set", "project", gcrProjectID)
					if err == nil {
						input = "gcloud auth configure-docker"
					}
				}
			}
			// gcloud auth login
			// gcloud config set project [PROJECT-ID]
			// gcloud auth configure-docker
       } else if strings.Contains(cr_url, "docker.io") {
		 check_url_format := regexp.MustCompile(`^docker\.io/[a-zA-Z0-9-]+$`).MatchString(cr_url)
		 if check_url_format { //check if docker login is not failed - skip this step
			  username := regexp.MustCompile(`^[^/]+/(.+)$`).FindStringSubmatch(cr_url)[1]
			  if os.Getenv("DOCKER_ACCESS_TOKEN") == "" {
                return diag.FromErr(fmt.Errorf("docker access token is not set. Please set it using 'export DOCKER_ACCESS_TOKEN=<your_token>'"))
            }
			  input = fmt.Sprintf("echo \"$DOCKER_ACCESS_TOKEN\" | docker login -u %s --password-stdin", username) // with username and token or password
			}
		 }
		 output,err := exec.Command("sh", "-c", input).CombinedOutput()

		 if err != nil {
			 return diag.FromErr(fmt.Errorf("you agagaga yeh %v \nAUTPUT!!!!!!! >>>> %s",err,string(output)))
		 }
	   }

		

	

	// Execute the Docker push command if provided
	// if docker_push != "" {
	// 	if strings.Contains(cr_url, "amazonaws.com") || strings.Contains(cr_url, "azurecr.io") {
	// 		fmt.Println("DO NUFIN!!!!")
	// 	} else {
	// 		cmd := exec.Command("sh", "-c", fmt.Sprintf("docker login --username %v --password %v", container_registry_url, cr_pass))
	// 		err := cmd.Run()
	// 		if err != nil {
	// 			// give them more attempts cos it can be wrong creds
	// 			return diag.FromErr(err)
	// 		}
	// 	}

	// 	output, err := exec.Command("sh", "-c", fmt.Sprintf("docker push %v", docker_push)).CombinedOutput()
	// 	if err != nil {
	// 		return diag.FromErr(fmt.Errorf("docker push failed with error: %v\nOutput: %s", err, string(output)))
	// 	}
	// }

	// Set the ID and timestamp for the resource
	d.SetId(fmt.Sprintf("%s-%s", strings.ToLower(build), strings.ToLower(test)))
	d.Set("timestamp", time.Now().Format(time.RFC3339))

	return diags

}

func resourceCICDRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceCICDUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceCICDCreate(ctx, d, m)
}

func resourceCICDDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
