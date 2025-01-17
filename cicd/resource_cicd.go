package cicd

import (
    "context"
    "fmt"
    "os/exec"

    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceCICD() *schema.Resource {
    return &schema.Resource{
        CreateContext: resourceCICDCreate,
        ReadContext:   resourceCICDRead,
        UpdateContext: resourceCICDUpdate,
        DeleteContext: resourceCICDDelete,

        Schema: map[string]*schema.Schema{
            "programming_language": {
                Type:     schema.TypeString,
                Required: true,
            },
            "build_tool": {
                Type:     schema.TypeString,
                Required: true,
            },
            "docker_tag": {
                Type:     schema.TypeString,
                Required: true,
            },
            "registry_url": {
                Type:     schema.TypeString,
                Required: true,
            },
            "registry_credentials": {
                Type:      schema.TypeString,
                Required:  true,
                Sensitive: true,
            },
            "skip_build": {
                Type:     schema.TypeBool,
                Optional: true,
                Default:  false,
            },
            "skip_test": {
                Type:     schema.TypeBool,
                Optional: true,
                Default:  false,
            },
            "skip_build_docker_image": {
                Type:     schema.TypeBool,
                Optional: true,
                Default:  false,
            },
            "skip_push_docker_image": {
                Type:     schema.TypeBool,
                Optional: true,
                Default:  false,
            },
            "custom_build_command": {
                Type:     schema.TypeString,
                Optional: true,
            },
            "custom_test_command": {
                Type:     schema.TypeString,
                Optional: true,
            },
            "environment_variables": {
                Type:     schema.TypeMap,
                Optional: true,
                Elem:     &schema.Schema{Type: schema.TypeString},
            },
        },
    }
}

func resourceCICDCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics

    programmingLanguage := d.Get("programming_language").(string)
    buildTool := d.Get("build_tool").(string)
    dockerTag := d.Get("docker_tag").(string)
    registryURL := d.Get("registry_url").(string)
    registryCredentials := d.Get("registry_credentials").(string)
    skipBuild := d.Get("skip_build").(bool)
    skipTest := d.Get("skip_test").(bool)
    skipBuildDockerImage := d.Get("skip_build_docker_image").(bool)
    skipPushDockerImage := d.Get("skip_push_docker_image").(bool)
    customBuildCommand := d.Get("custom_build_command").(string)
    customTestCommand := d.Get("custom_test_command").(string)
    environmentVariables := d.Get("environment_variables").(map[string]interface{})

    // Process environment variables
    for key, value := range environmentVariables {
        err := exec.Command("sh", "-c", fmt.Sprintf(`export %s=%s`, key, value)).Run()
        if err != nil {
            return diag.FromErr(err)
        }
    }

    // Implement the build, test, and docker steps
    if !skipBuild {
        if customBuildCommand != "" {
            err := exec.Command("sh", "-c", customBuildCommand).Run()
            if err != nil {
                return diag.FromErr(err)
            }
        } else {
            err := exec.Command("sh", "-c", fmt.Sprintf("%s %s", buildTool, "build")).Run()
            if err != nil {
                return diag.FromErr(err)
            }
        }
    }

    if !skipTest {
        if customTestCommand != "" {
            err := exec.Command("sh", "-c", customTestCommand).Run()
            if err != nil {
                return diag.FromErr(err)
            }
        } else {
            err := exec.Command("sh", "-c", fmt.Sprintf("%s %s", buildTool, "test")).Run()
            if err != nil {
                return diag.FromErr(err)
            }
        }
    }

    if !skipBuildDockerImage {
        err := exec.Command("sh", "-c", fmt.Sprintf("docker build -t %s .", dockerTag)).Run()
        if err != nil {
            return diag.FromErr(err)
        }
    }

    if !skipPushDockerImage {
        err := exec.Command("sh", "-c", fmt.Sprintf("echo %s | docker login %s --username %s --password-stdin", registryCredentials, registryURL, registryCredentials)).Run()
        if err != nil {
            return diag.FromErr(err)
        }
        err = exec.Command("sh", "-c", fmt.Sprintf("docker push %s", dockerTag)).Run()
        if err != nil {
            return diag.FromErr(err)
        }
    }

    // Set the ID for the resource
    d.SetId(fmt.Sprintf("%s-%s", programmingLanguage, dockerTag))

    return diags
}

func resourceCICDRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    // Read logic here
    return nil
}

func resourceCICDUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    // Update logic here
    return nil
}

func resourceCICDDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    // Delete logic here
    return nil
}