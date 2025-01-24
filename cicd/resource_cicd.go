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
        CustomizeDiff: customizeDiff, // Add CustomizeDiff function

        Schema: map[string]*schema.Schema{
            "step_1": {
                Type:     schema.TypeString,
                Required: true,
            },
            "step_2": {
                Type:     schema.TypeString,
                Optional: true,
            },
            "path": {
                Type:     schema.TypeString,
                Required: true,
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

func executeCommand(command string) error {
    fmt.Printf("Executing command: %s\n", command)
    cmd := exec.Command("sh", "-c", command)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    if err != nil {
        return fmt.Errorf("command failed with error: %v", err)
    }
    return nil
}

func resourceCICDCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics

    step1 := d.Get("step_1").(string)
    step2 := d.Get("step_2").(string)
    path := d.Get("path").(string)

    // Change directory
    if path != "" {
        fmt.Printf("Changing directory to: %s\n", path)
        err := os.Chdir(path)
        if err != nil {
            fmt.Printf("Failed to change directory: %v\n", err)
            return diag.FromErr(fmt.Errorf("failed to change directory: %v", err))
        }
    }

    // Execute step 1
    if step1 != "" {
        err := executeCommand(step1)
        if err != nil {
            return diag.FromErr(err)
        }
    }

    // Execute step 2
    if step2 != "" {
        err := executeCommand(step2)
        if err != nil {
            return diag.FromErr(err)
        }
    }

    // Set the ID and timestamp for the resource
    d.SetId(fmt.Sprintf("%s-%s", strings.ToLower(step1), strings.ToLower(step2)))
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