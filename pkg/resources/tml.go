package resources

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/daniepett/thoughtspot-sdk-go"
	"github.com/daniepett/thoughtspot-sdk-go/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &TmlResource{}
	_ resource.ResourceWithConfigure = &TmlResource{}
	// _ resource.ResourceWithImportState = &TmlResource{}
)

func NewTmlResource() resource.Resource {
	return &TmlResource{}
}

type TmlResource struct {
	client *thoughtspot.Client
}

// orderResourceModel maps the resource schema data.
type TmlResourceModel struct {
	ID    types.String `tfsdk:"id"`
	Tml   types.String `tfsdk:"tml"`
	Guids types.List   `tfsdk:"guids"`
}

type TmlGuidModel struct {
	Original types.String `tfsdk:"original"`
	Computed types.String `tfsdk:"computed"`
}

func (o TmlGuidModel) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"original": types.StringType,
		"computed": types.StringType,
	}
}

// Metadata returns the resource type name.
func (r *TmlResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tml"
}

// Schema defines the schema for the resource.
func (r *TmlResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tml": schema.StringAttribute{
				Required: true,
			},
			"guids": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"original": schema.StringAttribute{
							Computed: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"computed": schema.StringAttribute{
							Computed: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
				Computed: true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *TmlResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func exportTml(ctx context.Context, client *thoughtspot.Client, id string, tml string, existingGuids []MetadataGuidModel) (*TmlResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	cr := models.ExportMetadataTMLRequest{
		Metadata: []models.ExportMetadataTypeInput{models.ExportMetadataTypeInput{
			Identifier: id,
		}},
		EdocFormat: "YAML",
		ExportOptions: models.ExportOptions{
			IncludeGuid: false,
		},
	}

	c, err := client.ExportMetadataTML(cr)
	if err != nil {
		diags.AddError(
			"Error Reading Metadata",
			"Could not read Metadata ID: "+err.Error(),
		)
		return nil, diags
	}

	if len(c) == 0 {
		return nil, diags
	}

	metadata := c[0]

	var guids []MetadataGuidModel

	if existingGuids != nil {
		guids = existingGuids
	} else {
		re := regexp.MustCompile(`guid: ([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})`)
		ogids := re.FindAllStringSubmatch(tml, -1)
		cgids := re.FindAllStringSubmatch(metadata.Edoc, -1)
		for j := range ogids {
			guid := MetadataGuidModel{
				Original: types.StringValue(ogids[j][1]),
				Computed: types.StringValue(cgids[j][1]),
			}
			guids = append(guids, guid)

		}
	}

	tmlExport := metadata.Edoc

	for _, guid := range guids {
		tmlExport = strings.Replace(tml, guid.Original.ValueString(), guid.Computed.ValueString(), 1)
	}

	lg, diag := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: MetadataGuidModel{}.attrTypes()}, guids)

	diags.Append(diag...)

	m := TmlResourceModel{
		ID:    types.StringValue(metadata.Info.Id),
		Tml:   types.StringValue(tmlExport),
		Guids: lg,
	}

	return &m, diags
}

// Create a new resource.
func (r *TmlResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan TmlResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cr := models.ImportMetadataTMLRequest{
		MetadataTmls: []string{plan.Tml.ValueString()},
		CreateNew:    true,
	}

	c, err := r.client.ImportMetadataTML(cr)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating data connection",
			"Could not create connection, unexpected error: "+err.Error(),
		)
		return
	}

	if len(c) == 0 {

	}

	if c[0].Response.Status.StatusCode == "ERROR" {
		resp.Diagnostics.AddError(
			"Error importing TML",
			"Could not import tml , unexpected error: "+c[0].Response.Status.ErrorMessage,
		)
		return
	}
	id := c[0].Response.Header["id_guid"].(string)

	ex, _ := exportTml(ctx, r.client, id, plan.Tml.ValueString(), nil)

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(id)
	plan.Tml = ex.Tml
	plan.Guids = ex.Guids

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *TmlResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state TmlResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	var guids []MetadataGuidModel
	diags = state.Guids.ElementsAs(ctx, &guids, false)
	resp.Diagnostics.Append(diags...)
	ex, _ := exportTml(ctx, r.client, state.ID.ValueString(), state.Tml.ValueString(), guids)

	if ex == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Guids = ex.Guids
	state.Tml = ex.Tml

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *TmlResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan TmlResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var guids []MetadataGuidModel
	diags = plan.Guids.ElementsAs(ctx, &guids, false)
	resp.Diagnostics.Append(diags...)
	tml := plan.Tml.ValueString()
	for _, guid := range guids {
		tml = strings.Replace(tml, guid.Original.ValueString(), guid.Computed.ValueString(), 1)

	}

	cr := models.ImportMetadataTMLRequest{
		MetadataTmls: []string{tml},
	}

	_, err := r.client.ImportMetadataTML(cr)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating data connection",
			"Could not create connection, unexpected error: "+err.Error(),
		)
		return

	}

	ex, diag := exportTml(ctx, r.client, plan.ID.ValueString(), tml, nil)
	resp.Diagnostics.Append(diag...)

	plan.Tml = ex.Tml
	plan.Guids = ex.Guids

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TmlResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state TmlResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cr := models.DeleteMetadataRequest{
		Metadata: []models.DeleteMetadataTypeInput{{Identifier: state.ID.ValueString()}},
	}

	err := r.client.DeleteMetadata(cr)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Metadata",
			"Could not Metadata, unexpected error: "+err.Error(),
		)
		return
	}
}

// func (r *TmlResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
// 	// Retrieve import ID and save to id attribute
// 	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
// }
