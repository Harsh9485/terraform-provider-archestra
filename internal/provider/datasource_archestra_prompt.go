package provider

import (
	"context"
	"fmt"

	"github.com/archestra-ai/archestra/terraform-provider-archestra/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/oapi-codegen/runtime/types"
)

var _ datasource.DataSource = &ArchestraPromptDataSource{}

func NewArchestraPromptDataSource() datasource.DataSource {
	return &ArchestraPromptDataSource{}
}

type ArchestraPromptDataSource struct {
	client *client.ClientWithResponses
}

type ArchestraPromptDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Content     types.String `tfsdk:"content"`
	Tags        types.List   `tfsdk:"tags"`
	Visibility  types.String `tfsdk:"visibility"`
}

func (d *ArchestraPromptDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_archestra_prompt"
}

func (d *ArchestraPromptDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches an existing Archestra prompt by ID or name.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Prompt identifier",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the prompt",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the prompt",
				Computed:            true,
			},
			"content": schema.StringAttribute{
				MarkdownDescription: "The content of the prompt",
				Computed:            true,
			},
			"tags": schema.ListAttribute{
				MarkdownDescription: "List of tags for the prompt",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"visibility": schema.StringAttribute{
				MarkdownDescription: "Visibility of the prompt",
				Computed:            true,
			},
		},
	}
}

func (d *ArchestraPromptDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.ClientWithResponses)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.ClientWithResponses, got: %T", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *ArchestraPromptDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ArchestraPromptDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var prompt interface{} // Assume the type is the element from JSON200 slice

	if !data.ID.IsNull() {
		getResp, apiErr := d.client.GetPromptWithResponse(ctx, types.UUID(data.ID.ValueString()))
		if apiErr != nil {
			resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to read prompt, got error: %s", apiErr))
			return
		}
		if getResp.JSON200 == nil {
			resp.Diagnostics.AddError("Unexpected API Response", fmt.Sprintf("Expected 200 OK, got status %d", getResp.StatusCode()))
			return
		}
		prompt = getResp.JSON200
	} else if !data.Name.IsNull() {
		promptsResp, apiErr := d.client.GetPromptsWithResponse(ctx)
		if apiErr != nil {
			resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to read prompts, got error: %s", apiErr))
			return
		}
		if promptsResp.JSON200 == nil {
			resp.Diagnostics.AddError("Unexpected API Response", fmt.Sprintf("Expected 200 OK, got status %d", promptsResp.StatusCode()))
			return
		}
		for _, p := range *promptsResp.JSON200 {
			if p.Name == data.Name.ValueString() {
				prompt = p
				break
			}
		}
		if prompt == nil {
			resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Prompt with name '%s' not found", data.Name.ValueString()))
			return
		}
	} else {
		resp.Diagnostics.AddError("Invalid Configuration", "Either 'id' or 'name' must be provided")
		return
	}

	// Map to state (adjust based on available fields)
	if p, ok := prompt.(struct{ Id types.UUID; Name string; Content string }); ok { // Placeholder type
		data.ID = types.StringValue(p.Id.String())
		data.Name = types.StringValue(p.Name)
		data.Content = types.StringValue(p.Content)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
