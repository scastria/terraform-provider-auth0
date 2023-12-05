package auth0

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scastria/terraform-provider-auth0/auth0/client"
	"net/http"
	"net/url"
	"regexp"
)

func dataSourceEmailUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEmailUsersRead,
		Schema: map[string]*schema.Schema{
			"email": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`), "must be a valid email address"),
			},
			"user_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceEmailUsersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	email := d.Get("email").(string)
	requestQuery := url.Values{
		client.Email: []string{email},
	}
	body, err := c.HttpRequest(ctx, http.MethodGet, client.EmailUsersPath, requestQuery, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVals := []client.EmailUser{}
	err = json.NewDecoder(body).Decode(&retVals)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	numEmailUsers := len(retVals)
	if numEmailUsers < 1 {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("No user exists with that email"))
	}
	d.Set("email", email)
	userIds := []string{}
	for _, value := range retVals {
		userIds = append(userIds, value.UserId)
	}
	d.Set("user_ids", userIds)
	d.SetId(email)
	return diags
}
