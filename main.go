package main

import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
    "github.com/Andriizachepilo/terraform-provider-cicd/cicd"
)

func main() {
    plugin.Serve(&plugin.ServeOpts{
        ProviderFunc: func() *schema.Provider {
            return &schema.Provider{
                ResourcesMap: map[string]*schema.Resource{
                    "cicd_example": cicd.ResourceCICD(),
                },
            }
        },
    })
}