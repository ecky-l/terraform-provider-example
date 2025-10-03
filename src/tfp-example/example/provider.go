package example

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// This struct holds all the provider information. It is passed as receiver to every resource or data callback
type exampleProvider struct {
	host string
}

// Backing model for the schema.
// The types package contains specific data types for terraform schema resource data. They can be converted
// to Golang types and vice versa.
type exampleProviderModel struct {
	Host types.String `tfsdk:"host"`
}

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &exampleProvider{}
}

// Metadata returns the provider type name.
func (p *exampleProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "tfp_example"
}

func (p *exampleProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provider Configuration block to setup the example provider",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "HOST for Example API.",
				Optional:    true,
			},
		},
	}
}

// configure the provider's properties
func (provider *exampleProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var providerModel exampleProviderModel
	diags := req.Config.Get(ctx, &providerModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// check that the host was configured
	if providerModel.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"), "Unknown Example backend Host", "The \"host\" parameter in the provider config must not be empty")
	}

	provider.host = providerModel.Host.ValueString()

	// validate if the user has given a valid URL
	if _, err := url.ParseRequestURI(provider.host); err != nil {
		resp.Diagnostics.AddError("Invalid Host URL", fmt.Sprintf("The configured example host [%s] is not a valid URL", provider.host))
	}

	resp.DataSourceData = provider
	resp.ResourceData = provider
}

func (provider *exampleProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// Resources defines the resources implemented in the provider.
func (provider *exampleProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}
