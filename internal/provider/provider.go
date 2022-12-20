package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	dd "github.com/doximity/defect-dojo-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DefectDojoProvider satisfies the tfsdk.Provider interface and usually is included
// with all Resource and DataSource implementations.
type DefectDojoProvider struct {
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

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &DefectDojoProvider{}

// providerData can be used to store data from the Terraform configuration.
type providerData struct {
	BaseUrl  types.String `tfsdk:"base_url"`
	ApiKey   types.String `tfsdk:"api_key"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func newClient(ctx context.Context, url string, token string, user string, pass string) (*dd.ClientWithResponses, error) {
	// if we have no api key but we do have a user and password explicitly set,
	// we can use the API to get the token dynamically instead.
	if token == "" && user != "" && pass != "" {
		tokenclient, err := dd.NewClientWithResponses(url)
		if err != nil {
			return nil, fmt.Errorf("Error instantiating the client. This should never happen.")
		}

		tokenResponse, err := tokenclient.ApiTokenAuthCreateWithResponse(ctx, dd.ApiTokenAuthCreateJSONRequestBody{
			Username: user,
			Password: pass,
		})

		if err != nil {
			return nil, fmt.Errorf("Network error retrieving the token via the API: %s", err)
		}

		if tokenResponse.StatusCode() == 200 {
			token = tokenResponse.JSON200.Token
		} else {
			return nil, fmt.Errorf("Error retrieving the api token via the API. Unexpected response code: %d", tokenResponse.StatusCode())
		}
	}

	if token == "" {
		return nil, fmt.Errorf("Could not determine the api key for the defectdojo service. No api_key value provided and no DEFECTDOJO_APIKEY environment variable.")
	}

	tokenProvider, err := securityprovider.NewSecurityProviderApiKey("header", "Authorization", fmt.Sprintf("Token %s", token))
	if err != nil {
		return nil, fmt.Errorf("Error instantiating the security provider. This should never happen.")
	}

	client, err := dd.NewClientWithResponses(url, dd.WithRequestEditorFn(tokenProvider.Intercept))
	if err != nil {
		return nil, fmt.Errorf("Error instantiating the client. This should never happen.")
	}

	return client, nil
}

func (p *DefectDojoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data providerData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if p.configured {
		resp.DataSourceData = p.client
		resp.ResourceData = p.client
		return
	}

	var (
		url   string
		token string
		user  string
		pass  string
	)

	if data.BaseUrl.IsNull() {
		url = os.Getenv("DEFECTDOJO_BASEURL")
	} else {
		url = data.BaseUrl.ValueString()
	}

	if url == "" {
		resp.Diagnostics.AddError(
			"Unable to configure provider",
			"Could not determine the url of the defectdojo service. No base_url value provided and no DEFECTDOJO_BASEURL environment variable.",
		)
		return
	}

	if !data.ApiKey.IsNull() {
		token = data.ApiKey.ValueString()
	} else {
		token = os.Getenv("DEFECTDOJO_APIKEY")
	}

	if !data.Username.IsNull() {
		user = data.Username.ValueString()
	} else {
		user = os.Getenv("DEFECTDOJO_USERNAME")
	}

	if !data.Password.IsNull() {
		pass = data.Password.ValueString()
	} else {
		pass = os.Getenv("DEFECTDOJO_PASSWORD")
	}

	client, err := newClient(ctx, url, token, user, pass)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to configure provider",
			fmt.Sprintf("%s", err),
		)
		p.configured = false
		p.client = nil
		return
	}

	p.client = client
	p.configured = true

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *DefectDojoProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "defectdojo"
	resp.Version = p.version
}

func (p *DefectDojoProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProductResource,
		NewProductTypeResource,
		NewJiraProductConfigurationResource,
	}
}

func (p *DefectDojoProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewProductDataSource,
		NewProductTypeDataSource,
	}

}

func (p *DefectDojoProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Base URL of the defectdojo installation",
				Optional:            true,
				Required:            false,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "The API Key used to authenticate to defectdojo",
				Optional:            true,
				Required:            false,
				Sensitive:           true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The username used to authenticate to defectdojo. Has no effect if api_key is set.",
				Optional:            true,
				Required:            false,
				Sensitive:           true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password used to authenticate to defectdojo. Has no effect if api_key is set.",
				Optional:            true,
				Required:            false,
				Sensitive:           true,
			},
		},
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &DefectDojoProvider{
			version: version,
		}
	}
}
