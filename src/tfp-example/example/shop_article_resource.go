package example

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type shopArticleResource struct {
	provider *exampleProvider
}

type shopArticleModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

// NewApplicationResource is a helper function to simplify the provider implementation.
func NewShopArticleResource() resource.Resource {
	return &shopArticleResource{}
}

func (s *shopArticleResource) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	s.provider = request.ProviderData.(*exampleProvider)
}

// Metadata implements resource.Resource.
func (s *shopArticleResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_shop_article"
}

// Schema implements resource.Resource.
func (s *shopArticleResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "An example shop article resource, just for demonstration purposes.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Computed or literal import ID for the shop article (UUID)",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Shop article name",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Shop article description",
				Required:    true,
			},
		},
	}
}

// Create implements resource.Resource.
func (s *shopArticleResource) Create(context.Context, resource.CreateRequest, *resource.CreateResponse) {
	panic("unimplemented")
}

// Delete implements resource.Resource.
func (s *shopArticleResource) Delete(context.Context, resource.DeleteRequest, *resource.DeleteResponse) {
	panic("unimplemented")
}

// Read implements resource.Resource.
func (s *shopArticleResource) Read(context.Context, resource.ReadRequest, *resource.ReadResponse) {
	panic("unimplemented")
}

// Update implements resource.Resource.
func (s *shopArticleResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {
	panic("unimplemented")
}
