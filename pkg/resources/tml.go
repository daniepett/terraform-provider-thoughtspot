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
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
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
	ID          types.String `tfsdk:"id"`
	Tml         types.String `tfsdk:"tml"`
	Guids       types.List   `tfsdk:"guids"`
	UseObjectId types.Bool   `tfsdk:"use_object_id"`
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

func requiresReplaceIfGuidChanged() planmodifier.String {
	return requiresReplaceIfGuidChangedModifier{}
}

type requiresReplaceIfGuidChangedModifier struct{}

func (m requiresReplaceIfGuidChangedModifier) Description(ctx context.Context) string {
	return "Forces replacement if the guid changes"
}

func (m requiresReplaceIfGuidChangedModifier) MarkdownDescription(ctx context.Context) string {
	return "Forces replacement if the guid changes"
}

func (m requiresReplaceIfGuidChangedModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.ConfigValue.IsNull() || req.StateValue.IsNull() {
		return
	}

	if req.ConfigValue.Equal(req.StateValue) {
		return
	}

	// If we don't have a guid in the state or config, then don't evaluate
	if !strings.HasPrefix(req.StateValue.String(), "guid:") || !strings.HasPrefix(req.ConfigValue.String(), "guid:") {
		return
	}

	re := regexp.MustCompile(`guid: ([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})`)
	currentGuids := re.FindAllStringSubmatch(req.StateValue.String(), -1)
	newGuids := re.FindAllStringSubmatch(req.ConfigValue.String(), -1)
	// Checks the first guid value hasn't changed
	if currentGuids[0][1] != newGuids[0][1] {
		resp.RequiresReplace = true
	}
}

func (r *TmlResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config TmlResourceModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.UseObjectId.ValueBool() {
		tml := config.Tml.ValueString()

		// Check for GUIDs in the TML string
		guidRegex := regexp.MustCompile(`^guid: ([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})`)
		if guidRegex.MatchString(tml) {
			resp.Diagnostics.AddAttributeError(
				path.Root("tml"),
				"GUIDs Not Allowed When Using Object ID",
				"When 'use_object_id' is set to true, the 'tml' attribute must not contain any GUIDs.",
			)
		}
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
				PlanModifiers: []planmodifier.String{
					requiresReplaceIfGuidChanged(),
				},
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
			"use_object_id": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Flag to use object id and not guid mapping in TML import",
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
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

func exportTml(ctx context.Context, client *thoughtspot.Client, id string, tml string, existingGuids []MetadataGuidModel, useObjectId bool) (*TmlResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	cr := models.ExportMetadataTMLRequest{
		Metadata: []models.ExportMetadataTypeInput{models.ExportMetadataTypeInput{
			Identifier: id,
		}},
		EdocFormat: "YAML",
		ExportOptions: models.ExportOptions{
			IncludeGuid:  !useObjectId,
			IncludeObjId: useObjectId,
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

	if c[0].Info.Status.StatusCode == "ERROR" {
		diags.AddError(
			"Error reading TML",
			"Could not read tml , unexpected error: "+c[0].Info.Status.ErrorMessage,
		)
		return nil, diags
	}

	if len(c) == 0 {
		return nil, diags
	}

	metadata := c[0]

	tmlExport := metadata.Edoc
	var guids []MetadataGuidModel

	re := regexp.MustCompile(`guid: ([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})`)
	ogids := re.FindAllStringSubmatch(tml, -1)
	cgids := re.FindAllStringSubmatch(metadata.Edoc, -1)
	if (len(ogids) == 0 || len(cgids) == 0) && !useObjectId {
		diags.AddError(
			"Could not extract guids from TML",
			"No guids found for Metadata ID: "+id,
		)
		return nil, diags
	}
	if existingGuids != nil && len(ogids) != len(cgids) {
		guids = existingGuids
	} else {
		for j := range ogids {
			guid := MetadataGuidModel{
				Original: types.StringValue(ogids[j][1]),
				Computed: types.StringValue(cgids[j][1]),
			}
			guids = append(guids, guid)

		}
	}

	if len(guids) > 0 {
		for _, guid := range guids {
			tmlExport = strings.Replace(tml, guid.Computed.ValueString(), guid.Original.ValueString(), 1)
		}
	} else {
		// Computed attribute can't be nil
		guid := MetadataGuidModel{
			Original: types.StringValue(id),
			Computed: types.StringValue(id),
		}
		guids = append(guids, guid)

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
		ImportPolicy: "ALL_OR_NONE",
		CreateNew:    false,
	}

	c, err := r.client.ImportMetadataTML(cr)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing TML",
			"Could not import tml , unexpected error: "+err.Error(),
		)
		return
	}

	if c[0].Response.Status.StatusCode == "ERROR" {
		resp.Diagnostics.AddError(
			"Error importing TML",
			"Could not import tml , unexpected error: "+c[0].Response.Status.ErrorMessage,
		)
		return
	}
	id := c[0].Response.Header.IdGuid

	ex, diags := exportTml(ctx, r.client, id, plan.Tml.ValueString(), nil, plan.UseObjectId.ValueBool())
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(id)
	// Commented out while TS fixes viz_guid implementations
	// plan.Tml = ex.Tml
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

	// Ensure guids is not nil
	if guids == nil {
		guids = []MetadataGuidModel{}
	}
	ex, diags := exportTml(ctx, r.client, state.ID.ValueString(), state.Tml.ValueString(), guids, state.UseObjectId.ValueBool())
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if ex == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Tml = ex.Tml
	state.Guids = ex.Guids

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

	tml := plan.Tml.ValueString()

	var guids []MetadataGuidModel
	diags = plan.Guids.ElementsAs(ctx, &guids, false)
	resp.Diagnostics.Append(diags...)
	for _, guid := range guids {
		tml = strings.Replace(tml, guid.Original.ValueString(), guid.Computed.ValueString(), 1)

	}

	cr := models.ImportMetadataTMLRequest{
		MetadataTmls: []string{tml},
		ImportPolicy: "ALL_OR_NONE",
	}

	c, err := r.client.ImportMetadataTML(cr)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing TML",
			"Could not import tml , unexpected error: "+err.Error(),
		)
		return

	}

	if c[0].Response.Status.StatusCode == "ERROR" {
		resp.Diagnostics.AddError(
			"Error importing TML",
			"Could not import tml , unexpected error: "+c[0].Response.Status.ErrorMessage,
		)
		return
	}

	// ex, diag := exportTml(ctx, r.client, plan.ID.ValueString(), plan.Tml.ValueString(), nil)
	// resp.Diagnostics.Append(diag...)

	// plan.Tml = ex.Tml
	// plan.Guids = ex.Guids

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
