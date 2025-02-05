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

func executeCommand(command string, workDir string) error {
	if workDir == "" {
		return fmt.Errorf("working diretory is not specified")
	} else {
		fmt.Printf("Executing command: %s\n", command)
		cmd := exec.Command("sh", "-c", command)
		output, err := cmd.CombinedOutput()
		fmt.Printf("Command output: %s\n", string(output))
		if err != nil {
			return fmt.Errorf("command failed with error: %v\nOutput: %s", err, string(output))
		}
		return nil
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

	// Change to the specified working directory if provided
	if working_directory != "" {
		fmt.Printf("Changing directory to: %s\n", working_directory)
		err := os.Chdir(working_directory)
		if err != nil {
			fmt.Printf("Failed to change directory: %v\n", err)
			return diag.FromErr(fmt.Errorf("failed to change directory: %v", err))
		}
	}

	// Execute the build command if provided
	if build != "" {
		err := executeCommand(build, working_directory)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Execute the test command if provided
	if test != "" {
		err := executeCommand(test, working_directory)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Execute the build_and_test command if provided
	if build_and_test != "" {
		err := executeCommand(build_and_test, working_directory)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Execute the Docker build command if provided
	if docker_build != "" {
		if dockerfile_dir != "" {
			fmt.Printf("Changing directory to: %s\n", dockerfile_dir)
			err := os.Chdir(dockerfile_dir)
			if err != nil {
				fmt.Printf("Failed to change directory: %v\n", err)
				return diag.FromErr(fmt.Errorf("failed to change directory: %v", err))
			}
		}
		cmd := exec.Command("sh", "-c", fmt.Sprintf("docker build -t %v .", docker_build))
		output, err := cmd.CombinedOutput()
		fmt.Printf("Command output: %s\n", string(output))
		if err != nil {
			if strings.Contains(string(output), "Is the docker daemon running?") {
				fmt.Println("Docker is not running!")
			} else if strings.Contains(string(output), "Not found") {
				fmt.Println("Docker is not installed, please install Docker and try again")
			} else {
				return diag.FromErr(fmt.Errorf("docker build failed with error: %v\nOutput: %s", err, string(output)))
			}
		} else {
			fmt.Println("Image was built successfully!")
		}
	}

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

		cmd := exec.Command("sh", "-c", fmt.Sprintf("docker push %v", docker_push))
		output, err := cmd.CombinedOutput()
		fmt.Printf("Command output: %s\n", string(output))
		if err != nil {
			return diag.FromErr(fmt.Errorf("docker push failed with error: %v\nOutput: %s", err, string(output)))
		} else {
			fmt.Println("Images were pushed successfully")
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
