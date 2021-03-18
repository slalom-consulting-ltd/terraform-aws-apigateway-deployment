package apigwdeployment

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"profile": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("api_gw_aws_profile", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"apigwdeployment_methods": resourceAwsApiGatewayDeployment(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	profile := d.Get("profile").(string)

	var diags diag.Diagnostics

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(profile))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "no default config found: ~/.aws/config",
		})
		return nil, diags
	}

	return cfg, diags
}
