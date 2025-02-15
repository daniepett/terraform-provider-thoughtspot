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
	_ resource.Resource              = &MetadataResource{}
	_ resource.ResourceWithConfigure = &MetadataResource{}
	// _ resource.ResourceWithImportState = &MetadataResource{}
)

func NewMetadataResource() resource.Resource {
	return &MetadataResource{}
}

type MetadataResource struct {
	client *thoughtspot.Client
}

// orderResourceModel maps the resource schema data.
type MetadataResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Metadata     types.List   `tfsdk:"metadata"`
	ImportPolicy types.String `tfsdk:"import_policy"`
}

type MetadataGuidModel struct {
	Original types.String `tfsdk:"original"`
	Computed types.String `tfsdk:"computed"`
}

type MetadataExportModel struct {
	ID    types.String `tfsdk:"id"`
	Tml   types.String `tfsdk:"tml"`
	Guids types.List   `tfsdk:"guids"`
}

func (o MetadataGuidModel) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"original": types.StringType,
		"computed": types.StringType,
	}
}

func (o MetadataExportModel) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":    types.StringType,
		"tml":   types.StringType,
		"guids": types.ListType{ElemType: types.ObjectType{AttrTypes: MetadataGuidModel{}.attrTypes()}},
	}
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
			"import_policy": schema.StringAttribute{
				Optional: true,
			},
		},
		Blocks: map[string]schema.Block{
			"metadata": schema.ListNestedBlock{
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedBlockObject{
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
				},
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

func exportTmlsMetadata(ctx context.Context, client *thoughtspot.Client, ids []string, tmls []string) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	var emi []models.ExportMetadataTypeInput
	for _, element := range ids {
		emi = append(emi, models.ExportMetadataTypeInput{
			Identifier: element,
		})
	}

	cr := models.ExportMetadataTMLRequest{
		Metadata:   emi,
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
		return types.ListNull(types.ObjectType{AttrTypes: MetadataExportModel{}.attrTypes()}), diags
	}

	if len(c) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: MetadataExportModel{}.attrTypes()}), diags
	}

	var mems []MetadataExportModel
	for i := range c {
		re := regexp.MustCompile(`guid: ([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})`)
		ogids := re.FindAllStringSubmatch(tmls[i], -1)
		cgids := re.FindAllStringSubmatch(c[i].Edoc, -1)
		var guids []MetadataGuidModel

		tmlExport := c[i].Edoc
		for j := range ogids {
			guid := MetadataGuidModel{
				Original: types.StringValue(ogids[j][1]),
				Computed: types.StringValue(cgids[j][1]),
			}
			tmlExport = strings.Replace(tmls[i], guid.Computed.ValueString(), guid.Original.ValueString(), 1)
			guids = append(guids, guid)

		}
		lg, diag := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: MetadataGuidModel{}.attrTypes()}, guids)

		diags.Append(diag...)

		mems = append(mems, MetadataExportModel{
			ID:    types.StringValue(c[i].Info.Id),
			Tml:   types.StringValue(tmlExport),
			Guids: lg,
		})

	}

	m, diag := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: MetadataExportModel{}.attrTypes()}, mems)

	diags.Append(diag...)

	return m, diags
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

	var tmls []string
	var metadata []MetadataExportModel
	diags = plan.Metadata.ElementsAs(ctx, &metadata, false)
	resp.Diagnostics.Append(diags...)

	for _, t := range metadata {
		tmls = append(tmls, t.Tml.ValueString())
	}

	cr := models.ImportMetadataTMLRequest{
		MetadataTmls: tmls,
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

	var ids []string
	for i := range c {
		ids = append(ids, c[i].Response.Header["id_guid"].(string))
	}

	ex, _ := exportTmlsMetadata(ctx, r.client, ids, tmls)

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(ids[0])

	plan.Metadata = ex

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

	var tmls []string
	var ids []string

	var metadata []MetadataExportModel
	diags = state.Metadata.ElementsAs(ctx, &metadata, false)
	for _, t := range metadata {
		tmls = append(tmls, t.Tml.ValueString())
		ids = append(ids, t.ID.ValueString())
	}

	ex, _ := exportTmlsMetadata(ctx, r.client, ids, tmls)

	state.Metadata = ex

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

	var formattedTmls []string
	var tmls []string
	var metadata []MetadataExportModel
	var originalIds []string
	diags = plan.Metadata.ElementsAs(ctx, &metadata, false)
	resp.Diagnostics.Append(diags...)

	for _, t := range metadata {
		var guids []MetadataGuidModel
		diags = t.Guids.ElementsAs(ctx, &guids, false)
		resp.Diagnostics.Append(diags...)
		tml := t.Tml.ValueString()
		tmls = append(tmls, tml)
		for _, guid := range guids {
			tml = strings.Replace(tml, guid.Original.ValueString(), guid.Computed.ValueString(), 1)

		}
		formattedTmls = append(formattedTmls, tml)
		originalIds = append(originalIds, t.ID.ValueString())
	}

	cr := models.ImportMetadataTMLRequest{
		MetadataTmls: formattedTmls,
	}

	c, err := r.client.ImportMetadataTML(cr)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating data connection",
			"Could not create connection, unexpected error: "+err.Error(),
		)
		return
	}

	var ids []string
	for _, r := range c {
		ids = append(ids, r.Response.Header["id_guid"].(string))
	}

	ex, diag := exportTmlsMetadata(ctx, r.client, ids, tmls)

	resp.Diagnostics.Append(diag...)
	// fmt.Print("This is the read", ex)
	// Map response body to schema and populate Computed attribute values
	var newMetadata []MetadataExportModel
	diags = ex.ElementsAs(ctx, &newMetadata, false)
	resp.Diagnostics.Append(diags...)

	for _, item := range newMetadata {
		fmt.Print("This is the new tml", item.Tml)
	}

	plan.ID = types.StringValue(ids[0])
	plan.Metadata = ex

	// var deletedIds []models.DeleteMetadataTypeInput

	// for _, val1 := range originalIds {
	// 	found := false

	// 	for _, val2 := range ids {
	// 		if val1 == val2 {
	// 			found = true
	// 			break
	// 		}
	// 	}
	// 	if !found {
	// 		fmt.Print("Not found", val1)
	// 		deletedIds = append(deletedIds, models.DeleteMetadataTypeInput{
	// 			Identifier: val1,
	// 		})
	// 	}
	// }

	// fmt.Print("these are deleted", deletedIds)

	// if len(deletedIds) > 0 {
	// 	cr := models.DeleteMetadataRequest{
	// 		Metadata: deletedIds,
	// 	}
	// 	err := r.client.DeleteMetadata(cr)

	// 	if err != nil {
	// 		resp.Diagnostics.AddError(
	// 			"Error deleting Metadata",
	// 			"Could not Metadata, unexpected error: "+err.Error(),
	// 		)
	// 		return
	// 	}
	// }

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
		Metadata: []models.DeleteMetadataTypeInput{},
	}

	var metadata []MetadataExportModel

	diags = state.Metadata.ElementsAs(ctx, &metadata, false)
	resp.Diagnostics.Append(diags...)

	for _, t := range metadata {
		cr.Metadata = append(cr.Metadata, models.DeleteMetadataTypeInput{Identifier: t.ID.ValueString()})
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

// func (r *MetadataResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
// 	// Retrieve import ID and save to id attribute
// 	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
// }
