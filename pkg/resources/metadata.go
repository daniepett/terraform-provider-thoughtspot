package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/daniepett/thoughtspot-sdk-go"
	"github.com/daniepett/thoughtspot-sdk-go/models"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &MetadataResource{}
	_ resource.ResourceWithConfigure = &MetadataResource{}
	// _ resource.ResourceWithImportState = &spaceResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewMetadataResource() resource.Resource {
	return &MetadataResource{}
}

// orderResource is the resource implementation.
type MetadataResource struct {
	client *thoughtspot.Client
}

// orderResourceModel maps the resource schema data.
type MetadataResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Tml         types.String `tfsdk:"tml"`
	TmlComputed types.String `tfsdk:"tml_computed"`
}

// Metadata returns the resource type name.
func (r *MetadataResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_metadata"
}

// Schema defines the schema for the resource.
func (r *MetadataResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"type": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tml_computed": schema.StringAttribute{
				Computed: true,
			},
			"tml": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *MetadataResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*thoughtspot.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *qlikcloud.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Create a new resource.
func (r *MetadataResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan MetadataResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cr := models.ImportMetadataTMLRequest{
		MetadataTmls: []string{plan.Tml.ValueString()},
	}

	c, err := r.client.ImportMetadataTML(cr)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating data connection",
			"Could not create connection, unexpected error: "+err.Error(),
		)
		return
	}

	m := c[0].Response.Header
	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(m["id_guid"].(string))
	plan.Name = types.StringValue(m["name"].(string))
	plan.Type = types.StringValue(m["metadata_type"].(string))

	s := []string{"guid: ", plan.ID.ValueString(), "\n", plan.Tml.ValueString()}

	plan.TmlComputed = types.StringValue(strings.Join(s, ""))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *MetadataResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state MetadataResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cr := models.ExportMetadataTMLRequest{
		Metadata: []models.ExportMetadataTypeInput{
			{
				Type:       state.Type.ValueString(),
				Identifier: state.ID.ValueString(),
			}},
		EdocFormat: "YAML",
		ExportOptions: models.ExportOptions{
			IncludeGuid: false,
		},
	}

	c, err := r.client.ExportMetadataTML(cr)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Metadata",
			"Could not read Metadata ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	if len(c) == 0 {
		resp.State.RemoveResource(ctx)

		return
	}

	m := c[0]

	tml := strings.Join(strings.Split(m.Edoc, "\n")[1:], "\n")

	state.Tml = types.StringValue(tml)
	state.TmlComputed = types.StringValue(m.Edoc)
	state.Name = types.StringValue(m.Info.Name)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *MetadataResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan MetadataResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	s := []string{"guid: ", plan.ID.ValueString(), "\n", plan.Tml.ValueString()}
	plan.TmlComputed = types.StringValue(strings.Join(s, ""))
	cr := models.ImportMetadataTMLRequest{
		MetadataTmls: []string{plan.TmlComputed.ValueString()},
	}

	c, err := r.client.ImportMetadataTML(cr)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating data connection",
			"Could not create connection, unexpected error: "+err.Error(),
		)
		return
	}

	fmt.Print("This is the response from update", c)

	m := c[0].Response.Header
	// Map response body to schema and populate Computed attribute values
	plan.Name = types.StringValue(m["name"].(string))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *MetadataResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state MetadataResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cr := models.DeleteMetadataRequest{
		Metadata: []models.DeleteMetadataTypeInput{
			{
				Type:       state.Type.ValueString(),
				Identifier: state.ID.ValueString(),
			}},
	}

	err := r.client.DeleteMetadata(cr)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting metadata",
			"Could not metadata, unexpected error: "+err.Error(),
		)
		return
	}
}

// func (r *MetadataResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
// 	// Retrieve import ID and save to id attribute
// 	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
// }
