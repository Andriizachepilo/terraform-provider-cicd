package cicd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

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
            "step_1": {
                Type:     schema.TypeString,
                Required: true,
            },
            "step_2": {
                Type:     schema.TypeString,
                Required: false,
                Optional: true,
            },
            "path": {
                Type:     schema.TypeString,
                Required: true,
            },
        },
    }
}

func resourceCICDCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics

    step_1 := d.Get("step_1").(string)
    step_2 := d.Get("step_2").(string)
    path := d.Get("path").(string)
    
    if path != "" {
        fmt.Println("Changing directory to:", path)
        // Attempt to change directory
        err := os.Chdir(path)
        if err != nil {
            return diag.FromErr(fmt.Errorf("failed to change directory: %v", err))
        }
    } else {
        // If path is not specified, output a message
        fmt.Println("Path has not been specified.")
        // Here, you can return an error if you want to stop execution when path is empty
        err := exec.Command("sh", "-c", "echo 'Path has not been specified.'").Run()
        if err != nil {
            return diag.FromErr(fmt.Errorf("failed to execute command for empty path: %v", err))
        }
    }

     if step_1 != "" {
        command := strings.ToLower(step_1) //split and check if build is succesful, check for dependencies and install ?
        err := exec.Command("sh", "-c", command).Run()
        if err != nil {
            return diag.FromErr(err)
        }
     } else {
        err := exec.Command("sh", "-c", fmt.Sprintf("%s", "There are not steps specified")).Run()
        if err != nil {
            return diag.FromErr(err)
        }
     }

    if step_2 != "" {
        command := strings.ToLower(step_2)
        err := exec.Command("sh", "-c", command).Run()
        if err != nil {
            return diag.FromErr(err)
        }
     } else {
        err := exec.Command("sh","-c", fmt.Sprintf( "%s", "There are no steps specified")).Run()
        if err != nil {
            return diag.FromErr(err)
        }
     }



    // Set the ID for the resource
    d.SetId(fmt.Sprintf("%s-%s", strings.ToLower(step_1), strings.ToLower(step_2)))


    return diags
}





func resourceCICDRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    return nil
}

func resourceCICDUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    return nil
}

func resourceCICDDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    return nil
}