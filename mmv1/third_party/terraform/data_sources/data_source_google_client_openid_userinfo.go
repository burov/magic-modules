package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleClientOpenIDUserinfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleClientOpenIDUserinfoRead,
		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleClientOpenIDUserinfoRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	// When UserProjectOverride and BillingProject is set for the provider, the header X-Goog-User-Project is set for the API requests.
	// But it causes an error when calling GetCurrentUserEmail. It makes sense to not set header X-Goog-User-Project by setting UserProjectOverride
	// to false when calling GetCurrentUserEmail, because it does not create a bill.
	origUserProjectOverride := config.UserProjectOverride
	config.UserProjectOverride = false
	email, err := GetCurrentUserEmail(config, userAgent)
	config.UserProjectOverride = origUserProjectOverride

	if err != nil {
		return err
	}
	d.SetId(email)
	if err := d.Set("email", email); err != nil {
		return fmt.Errorf("Error setting email: %s", err)
	}
	return nil
}
