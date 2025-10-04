package example

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type shopArticleResource struct {
	provider *exampleProvider
}

type shopArticleResourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

type shopArticleRESTModel struct {
	ID          int64  `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
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
func (s *shopArticleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan shopArticleResourceModel
	diag := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, err := json.Marshal(shopArticleRESTModel{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error during article Create", fmt.Sprintf("Json marshalling error: %v", err))
		return
	}

	url := fmt.Sprintf("%s/articles", s.provider.host)
	httpResp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		resp.Diagnostics.AddError("Error during article Create", fmt.Sprintf("HTTP Post Error with URL [%s]: %v", url, err))
		return
	}
	if httpResp.StatusCode >= 400 {
		respB, _ := io.ReadAll(httpResp.Body)
		resp.Diagnostics.AddError("Error during article Create", fmt.Sprintf("Backend returned bad status %d: %s", httpResp.StatusCode, respB))
		return
	}

	var respBody shopArticleRESTModel
	dec := json.NewDecoder(httpResp.Body)
	if err := dec.Decode(&respBody); err != nil {
		resp.Diagnostics.AddError("Error during article Create", fmt.Sprintf("Json unmarshalling error: %v", err))
		return
	}

	plan.ID = types.Int64Value(respBody.ID)
	diag = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diag...)
}

// Delete implements resource.Resource.
func (s *shopArticleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state shopArticleResourceModel
	diag := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/articles/%d", s.provider.host, state.ID.ValueInt64())
	httpReq, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		resp.Diagnostics.AddError("Error during article Delete", fmt.Sprintf("Error creating DELETE request for URL [%s]: %v", url, err))
		return
	}

	client := &http.Client{}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Error during article Delete", fmt.Sprintf("HTTP DELETE Error with URL [%s]: %v", url, err))
		return
	}
	if httpResp.StatusCode != 204 {
		resp.Diagnostics.AddError("Error during article Delete", fmt.Sprintf("HTTP DELETE Error with URL [%s] returned not the expected result code. Should be 204 but was %d", url, httpResp.StatusCode))
	}
}

// Read implements resource.Resource.
func (s *shopArticleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state shopArticleResourceModel
	diag := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/articles/%d", s.provider.host, state.ID.ValueInt64())
	httpResp, err := http.Get(url)
	if err != nil {
		resp.Diagnostics.AddError("Error during article Read", fmt.Sprintf("HTTP Get Error with URL [%s]: %v", url, err))
		return
	}

	var respBody shopArticleRESTModel
	dec := json.NewDecoder(httpResp.Body)
	if err := dec.Decode(&respBody); err != nil {
		resp.Diagnostics.AddError("Error during article Create", fmt.Sprintf("Json unmarshalling error: %v", err))
		return
	}

	state.Name = types.StringValue(respBody.Name)
	state.Description = types.StringValue(respBody.Description)
	diag = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diag...)
}

// Update implements resource.Resource.
func (s *shopArticleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state shopArticleResourceModel
	diag := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}
	diag = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create payload for update
	body, err := json.Marshal(shopArticleRESTModel{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error during article Create", fmt.Sprintf("Json marshalling error: %v", err))
		return
	}

	// perform PUT request to backend
	url := fmt.Sprintf("%s/articles/%d", s.provider.host, state.ID.ValueInt64())
	httpReq, err := http.NewRequest("PUT", url, bytes.NewReader(body))
	if err != nil {
		resp.Diagnostics.AddError("Error during article Update", fmt.Sprintf("Error creating PUT request for URL [%s]: %v", url, err))
		return
	}

	client := &http.Client{}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Error during article Update", fmt.Sprintf("HTTP PUT Error with URL [%s]: %v", url, err))
		return
	}
	if httpResp.StatusCode >= 400 {
		respB, _ := io.ReadAll(httpResp.Body)
		resp.Diagnostics.AddError("Error during article Update", fmt.Sprintf("Backend returned bad status %d: %s", httpResp.StatusCode, respB))
		return
	}

	// read response
	var respBody shopArticleRESTModel
	dec := json.NewDecoder(httpResp.Body)
	if err := dec.Decode(&respBody); err != nil {
		resp.Diagnostics.AddError("Error during article Create", fmt.Sprintf("Json unmarshalling error: %v", err))
		return
	}

	state.ID = types.Int64Value(respBody.ID)
	state.Name = types.StringValue(respBody.Name)
	state.Description = types.StringValue(respBody.Description)

	diag = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diag...)
}
