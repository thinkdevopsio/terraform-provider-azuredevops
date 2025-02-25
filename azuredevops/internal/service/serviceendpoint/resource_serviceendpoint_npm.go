package serviceendpoint

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceServiceEndpointNpm schema and implementation for npm service endpoint resource
func ResourceServiceEndpointNpm() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointNpm, expandServiceEndpointNpm)

	r.Schema["url"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		Description:  "Url for the npm registry",
	}

	r.Schema["access_token"] = &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		Sensitive:        true,
		DiffSuppressFunc: tfhelper.DiffFuncSuppressSecretChanged,
		ValidateFunc:     validation.StringIsNotWhiteSpace,
		Description:      "The access token for npm registry",
	}
	// Add a spot in the schema to store the token secretly
	stSecretHashKey, stSecretHashSchema := tfhelper.GenerateSecreteMemoSchema("access_token")
	r.Schema[stSecretHashKey] = stSecretHashSchema

	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointNpm(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *uuid.UUID, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"apitoken": d.Get("access_token").(string),
		},
		Scheme: converter.String("Token"),
	}
	serviceEndpoint.Type = converter.String("externalnpmregistry")
	serviceEndpoint.Url = converter.String(d.Get("url").(string))
	return serviceEndpoint, projectID, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointNpm(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *uuid.UUID) {
	doBaseFlattening(d, serviceEndpoint, projectID)

	tfhelper.HelpFlattenSecret(d, "access_token")

	d.Set("url", *serviceEndpoint.Url)
	d.Set("access_token", (*serviceEndpoint.Authorization.Parameters)["apitoken"])
}
