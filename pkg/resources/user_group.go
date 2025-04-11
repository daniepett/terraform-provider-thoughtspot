package resources

import (
	"context"
	"fmt"

	thoughtspot "github.com/daniepett/thoughtspot-sdk-go"
	"github.com/daniepett/thoughtspot-sdk-go/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &UserGroupResource{}
	_ resource.ResourceWithConfigure = &UserGroupResource{}
	// _ resource.ResourceWithImportState = &spaceResource{}
)

func NewUserGroupResource() resource.Resource {
	return &UserGroupResource{}
}

type UserGroupResource struct {
	client *thoughtspot.Client
}

type UserGroupResourceModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	DisplayName       types.String `tfsdk:"display_name"`
	DefaultLiveboards types.List   `tfsdk:"default_liveboards"`
	Description       types.String `tfsdk:"description"`
	Privileges        types.List   `tfsdk:"privileges"`
	SubGroups         types.List   `tfsdk:"sub_groups"`
	Type              types.String `tfsdk:"type"`
	Users             types.List   `tfsdk:"users"`
	Visibility        types.String `tfsdk:"visibility"`
	Roles             types.List   `tfsdk:"roles"`
	RbacEnabled       types.Bool   `tfsdk:"rbac_enabled"`
}

// UserGroup returns the resource type name.
func (r *UserGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_group"
}

// Schema defines the schema for the resource.
func (r *UserGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"display_name": schema.StringAttribute{
				Required: true,
			},
			"default_liveboards": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default: listdefault.StaticValue(types.ListValueMust(
					types.StringType,
					[]attr.Value{},
				)),
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"privileges": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default: listdefault.StaticValue(types.ListValueMust(
					types.StringType,
					[]attr.Value{},
				)),
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"sub_groups": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default: listdefault.StaticValue(types.ListValueMust(
					types.StringType,
					[]attr.Value{},
				)),
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Optional: true,
			},
			"users": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "List of user names to add to the user group, if not defined Terraform will not manage user assignment to the group",
				Optional:    true,
			},
			"visibility": schema.StringAttribute{
				Optional: true,
			},
			"roles": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default: listdefault.StaticValue(types.ListValueMust(
					types.StringType,
					[]attr.Value{},
				)),
			},
			"rbac_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *UserGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *UserGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan UserGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dli := make([]string, 0, len(plan.DefaultLiveboards.Elements()))
	diags = plan.DefaultLiveboards.ElementsAs(ctx, &dli, false)
	resp.Diagnostics.Append(diags...)

	sgi := make([]string, 0, len(plan.SubGroups.Elements()))
	diags = plan.SubGroups.ElementsAs(ctx, &sgi, false)
	resp.Diagnostics.Append(diags...)

	var ui []string
	if !plan.Users.IsNull() {
		ui = make([]string, 0, len(plan.Users.Elements()))
		diags = plan.Users.ElementsAs(ctx, &ui, false)
		resp.Diagnostics.Append(diags...)
	} else {
		ui = nil
	}

	var ri []string
	var p []string
	if plan.RbacEnabled.ValueBool() {
		ri = make([]string, 0, len(plan.Roles.Elements()))
		diags = plan.Roles.ElementsAs(ctx, &ri, false)
		resp.Diagnostics.Append(diags...)
		p = nil
	} else {
		p = make([]string, 0, len(plan.Privileges.Elements()))
		diags = plan.Privileges.ElementsAs(ctx, &p, false)
		resp.Diagnostics.Append(diags...)

	}

	cr := models.CreateUserGroupRequest{
		Name:                        plan.Name.ValueString(),
		DisplayName:                 plan.DisplayName.ValueString(),
		DefaultLiveboardIdentifiers: dli,
		Description:                 plan.Description.ValueString(),
		Privileges:                  p,
		SubGroupIdentifiers:         sgi,
		Type:                        plan.Type.ValueString(),
		UserIdentifiers:             ui,
		Visibility:                  plan.Visibility.ValueString(),
		RoleIdentifiers:             ri,
	}

	c, err := r.client.CreateUserGroup(cr)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating user group",
			"Could not create user group, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(c.Id)
	// plan.Name = types.StringValue(m["name"].(string))
	// plan.Type = types.StringValue(m["UserGroup_type"].(string))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *UserGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state UserGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cr := models.SearchUserGroupsRequest{
		GroupIdentifier: state.ID.ValueString(),
	}

	c, err := r.client.SearchUserGroups(cr)
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

	users := make([]string, len(m.Users))
	if !state.Users.IsNull() {
		for i := range m.Users {
			users[i] = m.Users[i].Name
		}
	} else {
		users = nil
	}

	sg := make([]string, len(m.SubGroups))
	for i := range m.SubGroups {
		sg[i] = m.SubGroups[i].Name
	}

	dl := make([]string, len(m.DefaultLiveboards))
	for i := range m.DefaultLiveboards {
		dl[i] = m.DefaultLiveboards[i].Id
	}

	state.Name = types.StringValue(m.Name)
	state.DisplayName = types.StringValue(m.DisplayName)
	state.Description = types.StringValue(m.Description)
	state.Type = types.StringValue(m.Type)
	state.Users, _ = types.ListValueFrom(ctx, types.StringType, users)
	state.SubGroups, _ = types.ListValueFrom(ctx, types.StringType, sg)
	state.DefaultLiveboards, _ = types.ListValueFrom(ctx, types.StringType, dl)
	state.Visibility = types.StringValue(m.Visibility)

	if state.RbacEnabled.ValueBool() {
		roles := make([]string, len(m.Roles))
		for i := range m.Roles {
			roles[i] = m.Roles[i].Id
		}
		state.Roles, diags = types.ListValueFrom(ctx, types.StringType, roles)
		resp.Diagnostics.Append(diags...)
		state.Privileges = types.ListValueMust(
			types.StringType,
			[]attr.Value{},
		)
	} else {
		state.Roles = types.ListValueMust(
			types.StringType,
			[]attr.Value{},
		)
		state.Privileges, diags = types.ListValueFrom(ctx, types.StringType, m.Privileges)
		resp.Diagnostics.Append(diags...)
	}

	fmt.Print("This is state", state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *UserGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan UserGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dli := make([]string, 0, len(plan.DefaultLiveboards.Elements()))
	diags = plan.DefaultLiveboards.ElementsAs(ctx, &dli, false)
	resp.Diagnostics.Append(diags...)

	sgi := make([]string, 0, len(plan.SubGroups.Elements()))
	diags = plan.SubGroups.ElementsAs(ctx, &sgi, false)
	resp.Diagnostics.Append(diags...)

	var ui []string
	if !plan.Users.IsNull() {
		ui = make([]string, 0, len(plan.Users.Elements()))
		diags = plan.Users.ElementsAs(ctx, &ui, false)
		resp.Diagnostics.Append(diags...)
	} else {
		ui = nil
	}

	var ri []string
	var p []string
	if plan.RbacEnabled.ValueBool() {
		ri = make([]string, 0, len(plan.Roles.Elements()))
		diags = plan.Roles.ElementsAs(ctx, &ri, false)
		resp.Diagnostics.Append(diags...)
		p = nil
	} else {
		p = make([]string, 0, len(plan.Privileges.Elements()))
		diags = plan.Privileges.ElementsAs(ctx, &p, false)
		resp.Diagnostics.Append(diags...)
		ri = nil
	}

	cr := models.UpdateUserGroupRequest{
		Name:                        plan.Name.ValueString(),
		DisplayName:                 plan.DisplayName.ValueString(),
		DefaultLiveboardIdentifiers: dli,
		Description:                 plan.Description.ValueString(),
		Privileges:                  p,
		SubGroupIdentifiers:         sgi,
		Type:                        plan.Type.ValueString(),
		UserIdentifiers:             ui,
		Visibility:                  plan.Visibility.ValueString(),
		RoleIdentifiers:             ri,
	}

	err := r.client.UpdateUserGroup(plan.ID.ValueString(), cr)
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

func (r *UserGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state UserGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteUserGroup(state.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting User Group",
			"Could not User Group, unexpected error: "+err.Error(),
		)
		return
	}
}

// func (r *UserGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
// 	// Retrieve import ID and save to id attribute
// 	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
// }
