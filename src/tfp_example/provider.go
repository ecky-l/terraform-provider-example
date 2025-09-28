package tfp_example

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type exampleProvider struct {
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
		Description: "Provider Configuration block to setup the galapagos provider",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "HOST for Example API. May also be provided via EXAMPLE_HOST environment variable.",
				Optional:    true,
			},
		},
	}
}

func (provider *exampleProvider) Configure(_ context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

func (provider *exampleProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// Resources defines the resources implemented in the provider.
func (provider *exampleProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}
