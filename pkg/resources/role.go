package resources

import (
	"context"
	"fmt"

	thoughtspot "github.com/daniepett/thoughtspot-sdk-go"
	"github.com/daniepett/thoughtspot-sdk-go/models"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &RoleResource{}
	_ resource.ResourceWithConfigure = &RoleResource{}
	// _ resource.ResourceWithImportState = &spaceResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewRoleResource() resource.Resource {
	return &RoleResource{}
}

// orderResource is the resource implementation.
type RoleResource struct {
	client *thoughtspot.Client
}

// orderResourceModel maps the resource schema data.
type RoleResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Privileges  types.List   `tfsdk:"privileges"`
}

// Role returns the resource type name.
func (r *RoleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

// Schema defines the schema for the resource.
func (r *RoleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"privileges": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(stringvalidator.OneOf([]string{"USERDATAUPLOADING",
						"DATADOWNLOADING",
						"DATAMANAGEMENT",
						"SHAREWITHALL",
						"JOBSCHEDULING",
						"A3ANALYSIS",
						"EXPERIMENTALFEATUREPRIVILEGE",
						"BYPASSRLS",
						"DISABLE_PINBOARD_CREATION",
						"DEVELOPER",
						"APPLICATION_ADMINISTRATION",
						"USER_ADMINISTRATION",
						"GROUP_ADMINISTRATION",
						"SYSTEM_INFO_ADMINISTRATION",
						"SYNCMANAGEMENT",
						"ORG_ADMINISTRATION",
						"ROLE_ADMINISTRATION",
						"AUTHENTICATION_ADMINISTRATION",
						"BILLING_INFO_ADMINISTRATION",
						"CONTROL_TRUSTED_AUTH",
						"TAGMANAGEMENT",
						"LIVEBOARD_VERIFIER",
						"CAN_MANAGE_CUSTOM_CALENDAR",
						"CAN_CREATE_OR_EDIT_CONNECTIONS",
						"CAN_MANAGE_WORKSHEET_VIEWS_TABLES",
						"CAN_MANAGE_VERSION_CONTROL",
						"THIRDPARTY_ANALYSIS",
						"CAN_CREATE_CATALOG",
						"ALLOW_NON_EMBED_FULL_APP_ACCESS",
						"CAN_ACCESS_ANALYST_STUDIO",
						"CAN_MANAGE_ANALYST_STUDIO",
						"PREVIEW_DOCUMENT_SEARCH",
						"CAN_SETUP_VERSION_CONTROL"}...)),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *RoleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*thoughtspot.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *thoughtspot.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Create a new resource.
func (r *RoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan RoleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	p := make([]string, 0, len(plan.Privileges.Elements()))
	_ = plan.Privileges.ElementsAs(ctx, &p, false)

	cr := models.CreateRoleRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Privileges:  p,
	}

	c, err := r.client.CreateRole(cr)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating user group",
			"Could not create user group, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(c.Id)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *RoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state RoleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cr := models.SearchRolesRequest{
		RoleIdentifiers: []string{state.ID.ValueString()},
	}

	c, err := r.client.SearchRoles(cr)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading User Group",
			"Could not read User Group ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	if len(c) == 0 {
		resp.State.RemoveResource(ctx)

		return
	}

	m := c[0]

	state.Name = types.StringValue(m.Name)
	state.Description = types.StringValue(m.Description)

	state.Privileges, _ = types.ListValueFrom(ctx, types.StringType, m.Privileges)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *RoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan RoleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	p := make([]string, 0, len(plan.Privileges.Elements()))
	_ = plan.Privileges.ElementsAs(ctx, &p, false)

	cr := models.UpdateRoleRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Privileges:  p,
	}

	_, err := r.client.UpdateRole(plan.ID.ValueString(), cr)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating User Group",
			"Could not update User Group, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *RoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state RoleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteRole(state.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting User Group",
			"Could not User Group, unexpected error: "+err.Error(),
		)
		return
	}
}

// func (r *RoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
// 	// Retrieve import ID and save to id attribute
// 	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
// }
