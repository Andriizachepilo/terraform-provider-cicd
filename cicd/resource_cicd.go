package cicd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
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
			"container_registry_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
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
			return fmt.Errorf("failed to change directory: %v", err),"."
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
	container_registry_url := d.Get("container_registry_url").(string)
	cr_pass := d.Get("container_registry_password").(string)
	docker_push := d.Get("docker_push").(string)

	// feedback := func(dia diag.Diagnostics, processName string, output string) diag.Diagnostics {
	// 	return append(dia, diag.Diagnostic{
	// 		Severity: diag.Warning,
	// 		Summary:  fmt.Sprintf("Step %v completed!\n%v", processName, output),
	// 	})
	// }
	feedback := func(processName, output string) diag.Diagnostics {
		return diag.Diagnostics{
			{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Step %v completed!\n%v", processName, output),
			},
		}
	}

	// if working_directory != "" {
	// 	err := os.Chdir(working_directory)
	// 	if err != nil {
	// 		return diag.FromErr(fmt.Errorf("failed to change directory: %v", err))
	// 	}
	// }

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
	
	// if build != "" {
	// 	err, outp := executeCommand(build, working_directory)
	// 	if err != nil {
	// 		return diag.FromErr(err)
	// 	} else {
	// 	    return feedback(diags, "build", outp)
	// 	}
	// }

	// if test != "" {
	// 	err, outp := executeCommand(test, working_directory)
	// 	if err != nil {
	// 		return diag.FromErr(err)
	// 	} else {
	// 		return feedback(diags, "test", outp)
	// 	}
	// }

	// // // Execute the build_and_test command if provided
	// if build_and_test != "" {
	// 	err, outp := executeCommand(build_and_test, working_directory)
	// 	if err != nil {
	// 		return diag.FromErr(err)
	// 	} else {
	// 		return feedback(diags, "build_and_test", outp)
	// 	}
	// }

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

	// Execute the Docker push command if provided
	if docker_push != "" {
		if strings.Contains(container_registry_url, "amazonaws.com") || strings.Contains(container_registry_url, "azurecr.io") {
			fmt.Println("DO NUFIN!!!!")
		} else {
			cmd := exec.Command("sh", "-c", fmt.Sprintf("docker login --username %v --password %v", container_registry_url, cr_pass))
			err := cmd.Run()
			if err != nil {
				// give them more attempts cos it can be wrong creds
				return diag.FromErr(err)
			}
		}

		output, err := exec.Command("sh", "-c", fmt.Sprintf("docker push %v", docker_push)).CombinedOutput()
		if err != nil {
			return diag.FromErr(fmt.Errorf("docker push failed with error: %v\nOutput: %s", err, string(output)))
		}
	}

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
