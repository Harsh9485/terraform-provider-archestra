package provider

import (
	"context"
	"fmt"

	"github.com/archestra-ai/archestra/terraform-provider-archestra/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/oapi-codegen/runtime/types"
)

var _ resource.Resource = &ArchestraPromptResource{}

func NewArchestraPromptResource() resource.Resource {
	return &ArchestraPromptResource{}
}

type ArchestraPromptResource struct {
	client *client.ClientWithResponses
}

type ArchestraPromptResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Content     types.String `tfsdk:"content"`
	Tags        types.List   `tfsdk:"tags"`
	Visibility  types.String `tfsdk:"visibility"`
	VersionID   types.String `tfsdk:"version_id"` // For rollback
}

func (r *ArchestraPromptResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_archestra_prompt"
}

func (r *ArchestraPromptResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an Archestra prompt, including creation, versioning, and rollback.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Prompt identifier",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the prompt",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the prompt",
				Optional:            true,
			},
			"content": schema.StringAttribute{
				MarkdownDescription: "The content of the prompt",
				Required:            true,
			},
			"tags": schema.ListAttribute{
				MarkdownDescription: "List of tags for the prompt",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"visibility": schema.StringAttribute{
				MarkdownDescription: "Visibility of the prompt (e.g., public, private)",
				Optional:            true,
			},
			"version_id": schema.StringAttribute{
				MarkdownDescription: "Version ID for rollback",
				Optional:            true,
			},
		},
	}
}

func (r *ArchestraPromptResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.ClientWithResponses)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.ClientWithResponses, got: %T", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *ArchestraPromptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ArchestraPromptResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare create request (assume only Name and Content are supported)
	createReq := client.CreatePromptJSONRequestBody{
		Name:    data.Name.ValueString(),
		Content: data.Content.ValueString(),
	}

	createResp, err := r.client.CreatePromptWithResponse(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to create prompt, got error: %s", err))
		return
	}

	if createResp.JSON201 == nil {
		resp.Diagnostics.AddError("Unexpected API Response", fmt.Sprintf("Expected 201 Created, got status %d", createResp.StatusCode()))
		return
	}

	// Set state
	data.ID = types.StringValue(createResp.JSON201.Id.String())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ArchestraPromptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ArchestraPromptResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	getResp, err := r.client.GetPromptWithResponse(ctx, types.UUID(data.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to read prompt, got error: %s", err))
		return
	}

	if getResp.JSON200 == nil {
		resp.Diagnostics.AddError("Unexpected API Response", fmt.Sprintf("Expected 200 OK, got status %d", getResp.StatusCode()))
		return
	}

	// Map response to state (assume only Name, Content, ID are available)
	prompt := getResp.JSON200
	data.Name = types.StringValue(prompt.Name)
	data.Content = types.StringValue(prompt.Content)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ArchestraPromptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ArchestraPromptResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Handle rollback if version_id is set
	if !data.VersionID.IsNull() {
		rollbackResp, err := r.client.RollbackPromptWithResponse(ctx, types.UUID(data.ID.ValueString()), client.RollbackPromptJSONRequestBody{
			VersionId: types.UUID(data.VersionID.ValueString()),
		})
		if err != nil {
			resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to rollback prompt, got error: %s", err))
			return
		}
		if rollbackResp.StatusCode() != 200 {
			resp.Diagnostics.AddError("Unexpected API Response", fmt.Sprintf("Expected 200 OK, got status %d", rollbackResp.StatusCode()))
			return
		}
	} else {
		// Regular update (assume only Name and Content are supported)
		updateReq := client.UpdatePromptJSONRequestBody{
			Name:    data.Name.ValueStringPointer(),
			Content: data.Content.ValueStringPointer(),
		}

		updateResp, err := r.client.UpdatePromptWithResponse(ctx, types.UUID(data.ID.ValueString()), updateReq)
		if err != nil {
			resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to update prompt, got error: %s", err))
			return
		}

		if updateResp.StatusCode() != 200 {
			resp.Diagnostics.AddError("Unexpected API Response", fmt.Sprintf("Expected 200 OK, got status %d", updateResp.StatusCode()))
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ArchestraPromptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ArchestraPromptResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteResp, err := r.client.DeletePromptWithResponse(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to delete prompt, got error: %s", err))
		return
	}

	if deleteResp.StatusCode() != 204 {
		resp.Diagnostics.AddError("Unexpected API Response", fmt.Sprintf("Expected 204 No Content, got status %d", deleteResp.StatusCode()))
		return
	}
}

// Helper functions (assume these are defined elsewhere or add them)
func convertTypesListToStringSlice(list types.List) []string {
	// Implementation to convert types.List to []string
}

func convertStringSliceToTypesList(slice []string) types.List {
	// Implementation to convert []string to types.List
}
