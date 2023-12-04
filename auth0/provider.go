package auth0

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scastria/terraform-provider-auth0/client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AUTH0_DOMAIN", nil),
			},
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AUTH0_CLIENT_ID", nil),
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AUTH0_CLIENT_SECRET", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			//"konnect_plugin":                  resourcePlugin(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			//"konnect_consumer":      dataSourceConsumer(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	domain := d.Get("domain").(string)
	clientId := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)

	var diags diag.Diagnostics
	c, err := client.NewClient(ctx, domain, clientId, clientSecret)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return c, diags
}
