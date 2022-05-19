package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	dd "github.com/doximity/defect-dojo-client-go"
)

// provider satisfies the tfsdk.Provider interface and usually is included
// with all Resource and DataSource implementations.
type provider struct {
	// client can contain the upstream provider SDK or HTTP client used to
	// communicate with the upstream service. Resource and DataSource
	// implementations can then make calls using this client.
	client *dd.ClientWithResponses

	// configured is set to true at the end of the Configure method.
	// This can be used in Resource and DataSource implementations to verify
	// that the provider was previously configured.
	configured bool

	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// providerData can be used to store data from the Terraform configuration.
type providerData struct {
	BaseUrl  types.String `tfsdk:"base_url"`
	ApiKey   types.String `tfsdk:"api_key"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	var data providerData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if p.configured {
		return
	}

	var (
		url   string
		token string
		user  string
		pass  string
	)

	if data.BaseUrl.Null {
		url = os.Getenv("DEFECTDOJO_BASEURL")
	} else {
		url = data.BaseUrl.Value
	}

	if url == "" {
		resp.Diagnostics.AddError(
			"Unable to configure provider",
			"Could not determine the url of the defectdojo service. No base_url value provided and no DEFECTDOJO_BASEURL environment variable.",
		)
		return
	}

	if !data.ApiKey.Null {
		token = data.ApiKey.Value
	} else {
		token = os.Getenv("DEFECTDOJO_APIKEY")
	}

	if !data.Username.Null {
		user = data.Username.Value
	} else {
		user = os.Getenv("DEFECTDOJO_USERNAME")
	}

	if !data.Password.Null {
		pass = data.Password.Value
	} else {
		pass = os.Getenv("DEFECTDOJO_PASSWORD")
	}

	// if we have no api key but we do have a user and password explicitly set,
	// we can use the API to get the token dynamically instead.
	if token == "" && user != "" && pass != "" {
		tokenclient, err := dd.NewClientWithResponses(url)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to configure provider",
				"Error instantiating the client. This should never happen.",
			)
			return
		}

		tokenResponse, err := tokenclient.ApiTokenAuthCreateWithResponse(ctx, dd.ApiTokenAuthCreateJSONRequestBody{
			Username: user,
			Password: pass,
		})

		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to configure provider",
				fmt.Sprintf("Network error retrieving the token via the API: %s", err),
			)
			return
		}

		if tokenResponse.StatusCode() == 200 {
			token = tokenResponse.JSON200.Token
		} else {
			resp.Diagnostics.AddError(
				"Unable to configure provider",
				fmt.Sprintf("Error retrieving the api token via the API. Unexpected response code: %d", tokenResponse.StatusCode()),
			)
			return
		}
	}

	if token == "" {
		resp.Diagnostics.AddError(
			"Unable to configure provider",
			"Could not determine the api key for the defectdojo service. No api_key value provided and no DEFECTDOJO_APIKEY environment variable.",
		)
		return
	}

	tokenProvider, err := securityprovider.NewSecurityProviderApiKey("header", "Authorization", fmt.Sprintf("Token %s", token))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to configure provider",
			"Error instantiating the security provider. This should never happen.",
		)
		return
	}

	client, err := dd.NewClientWithResponses(url, dd.WithRequestEditorFn(tokenProvider.Intercept))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to configure provider",
			"Error instantiating the client. This should never happen.",
		)
		return
	}

	p.client = client

	p.configured = true
}

func (p *provider) GetResources(ctx context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"defectdojo_product":      productResourceType{},
		"defectdojo_product_type": productTypeResourceType{},
	}, nil
}

func (p *provider) GetDataSources(ctx context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		"defectdojo_product":      productDataSourceType{},
		"defectdojo_product_type": productTypeDataSourceType{},
	}, nil
}

func (p *provider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"base_url": {
				MarkdownDescription: "Base URL of the defectdojo installation",
				Optional:            true,
				Required:            false,
				Type:                types.StringType,
			},
			"api_key": {
				MarkdownDescription: "The API Key used to authenticate to defectdojo",
				Optional:            true,
				Required:            false,
				Type:                types.StringType,
				Sensitive:           true,
			},
			"username": {
				MarkdownDescription: "The username used to authenticate to defectdojo. Has no effect if api_key is set.",
				Optional:            true,
				Required:            false,
				Type:                types.StringType,
				Sensitive:           true,
			},
			"password": {
				MarkdownDescription: "The password used to authenticate to defectdojo. Has no effect if api_key is set.",
				Optional:            true,
				Required:            false,
				Type:                types.StringType,
				Sensitive:           true,
			},
		},
	}, nil
}

func New(version string) func() tfsdk.Provider {
	return func() tfsdk.Provider {
		return &provider{
			version: version,
		}
	}
}

// convertProviderType is a helper function for NewResource and NewDataSource
// implementations to associate the concrete provider type. Alternatively,
// this helper can be skipped and the provider type can be directly type
// asserted (e.g. provider: in.(*provider)), however using this can prevent
// potential panics.
func convertProviderType(in tfsdk.Provider) (provider, diag.Diagnostics) {
	var diags diag.Diagnostics

	p, ok := in.(*provider)

	if !ok {
		diags.AddError(
			"Unexpected Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received. This is always a bug in the provider code and should be reported to the provider developers.", p),
		)
		return provider{}, diags
	}

	if p == nil {
		diags.AddError(
			"Unexpected Provider Instance Type",
			"While creating the data source or resource, an unexpected empty provider instance was received. This is always a bug in the provider code and should be reported to the provider developers.",
		)
		return provider{}, diags
	}

	return *p, diags
}
