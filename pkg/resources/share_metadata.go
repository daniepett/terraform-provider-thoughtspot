package resources

import (
	"context"
	"fmt"

	thoughtspot "github.com/daniepett/thoughtspot-sdk-go"
	"github.com/daniepett/thoughtspot-sdk-go/models"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &ShareMetadataResource{}
	_ resource.ResourceWithConfigure = &ShareMetadataResource{}
	// _ resource.ResourceWithImportState = &spaceResource{}
)

func NewShareMetadataResource() resource.Resource {
	return &ShareMetadataResource{}
}

type ShareMetadataResource struct {
	client *thoughtspot.Client
}

type ShareMetadataResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	MetadataType        types.String `tfsdk:"metadata_type"`
	MetadataIdentifier  types.String `tfsdk:"metadata_identifier"`
	PrincipalType       types.String `tfsdk:"principal_type"`
	PrincipalIdentifier types.String `tfsdk:"principal_identifier"`
	ShareMode           types.String `tfsdk:"share_mode"`
	Discoverable        types.Bool   `tfsdk:"discoverable"`
}

// ShareMetadata returns the resource type name.
func (r *ShareMetadataResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_share_metadata"
}

// Schema defines the schema for the resource.
func (r *ShareMetadataResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"metadata_type": schema.StringAttribute{
				Required:    true,
				Description: "Type of metadata. Required if identifier in metadata_identifier is a name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"LIVEBOARD",
						"ANSWER",
						"LOGICAL_TABLE",
						"LOGICAL_COLUMN",
						"CONNECTION"}...),
				},
			},
			"metadata_identifier": schema.StringAttribute{
				Required:    true,
				Description: "Unique ID or name of metadata object.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"principal_type": schema.StringAttribute{
				Required:    true,
				Description: "Principal type. Accepts `USER`, `USER_GROUP`",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"USER",
						"USER_GROUP"}...),
				},
			},
			"principal_identifier": schema.StringAttribute{
				Required:    true,
				Description: "Unique ID or name of the principal object such as a user or group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"share_mode": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Type of access to the shared object. Accepts `READ_ONLY`, `MODIFY`, `NO_ACCESS`",
				Default:     stringdefault.StaticString("READ_ONLY"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"READ_ONLY",
						"MODIFY",
						"NO_ACCESS"}...),
				},
			},
			"discoverable": schema.BoolAttribute{
				Required:    true,
				Description: "Flag to make the object discoverable.",
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *ShareMetadataResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *ShareMetadataResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ShareMetadataResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cr := models.ShareMetadataRequest{
		MetadataType:        plan.MetadataType.ValueString(),
		MetadataIdentifiers: []string{plan.MetadataIdentifier.ValueString()},
		Permissions: []models.SharePermissionsInput{models.SharePermissionsInput{
			Principal: models.PrincipalsInput{
				Identifier: plan.PrincipalIdentifier.ValueString(),
				Type:       plan.PrincipalType.ValueString(),
			},
			ShareMode: plan.ShareMode.ValueString(),
		},
		},
		HasLenientDiscoverability: plan.Discoverable.ValueBool(),
	}

	err := r.client.ShareMetadata(cr)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating user group",
			"Could not create user group, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(plan.MetadataIdentifier.ValueString() + "|" + plan.PrincipalType.ValueString())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *ShareMetadataResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state ShareMetadataResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cr := models.FetchPermissionsOnMetadataRequest{
		Metadata:   []models.PermissionsMetadataTypeInput{models.PermissionsMetadataTypeInput{Identifier: state.MetadataIdentifier.ValueString()}},
		Principals: []models.PrincipalsInput{models.PrincipalsInput{Identifier: state.PrincipalIdentifier.ValueString(), Type: state.PrincipalType.ValueString()}},
	}

	c, err := r.client.FetchPermissionsOnMetadata(cr)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading User Group",
			"Could not read User Group ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	if len(c.MetadataPermissionDetails) == 0 {
		resp.State.RemoveResource(ctx)

		return
	}

	p := c.MetadataPermissionDetails[0].PrincipalPermissionInfo[0].PrincipalPermissions[0].Permission

	state.ShareMode = types.StringValue(p)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ShareMetadataResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan ShareMetadataResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cr := models.ShareMetadataRequest{
		MetadataType:        plan.MetadataType.ValueString(),
		MetadataIdentifiers: []string{plan.MetadataIdentifier.ValueString()},
		Permissions: []models.SharePermissionsInput{models.SharePermissionsInput{
			Principal: models.PrincipalsInput{
				Identifier: plan.PrincipalIdentifier.ValueString(),
				Type:       plan.PrincipalType.ValueString(),
			},
			ShareMode: plan.ShareMode.ValueString(),
		},
		},
		HasLenientDiscoverability: plan.Discoverable.ValueBool(),
	}

	err := r.client.ShareMetadata(cr)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating user group",
			"Could not create user group, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ShareMetadataResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state ShareMetadataResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cr := models.ShareMetadataRequest{
		MetadataType:        state.MetadataType.ValueString(),
		MetadataIdentifiers: []string{state.MetadataIdentifier.ValueString()},
		Permissions: []models.SharePermissionsInput{models.SharePermissionsInput{
			Principal: models.PrincipalsInput{
				Identifier: state.PrincipalIdentifier.ValueString(),
				Type:       state.PrincipalType.ValueString(),
			},
			ShareMode: "NO_ACCESS",
		},
		},
		HasLenientDiscoverability: state.Discoverable.ValueBool(),
	}

	err := r.client.ShareMetadata(cr)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating user group",
			"Could not create user group, unexpected error: "+err.Error(),
		)
		return
	}

}

// func (r *ShareMetadataResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
// 	// Retrieve import ID and save to id attribute
// 	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
// }
