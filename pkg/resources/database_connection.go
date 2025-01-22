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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &DatabaseConnectionResource{}
	_ resource.ResourceWithConfigure = &DatabaseConnectionResource{}
	// _ resource.ResourceWithImportState = &spaceResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewDatabaseConnectionResource() resource.Resource {
	return &DatabaseConnectionResource{}
}

// orderResource is the resource implementation.
type DatabaseConnectionResource struct {
	client *thoughtspot.Client
}

// orderResourceModel maps the resource schema data.
type DatabaseConnectionResourceModel struct {
	ID                  types.String                                       `tfsdk:"id"`
	Name                types.String                                       `tfsdk:"name"`
	Description         types.String                                       `tfsdk:"description"`
	DataWarehouseType   types.String                                       `tfsdk:"data_warehouse_type"`
	DataWarehouseConfig DatabaseConnectionDataWarehouseConfigResourceModel `tfsdk:"data_warehouse_config"`
	Validate            types.Bool                                         `tfsdk:"validate"`
}

type DatabaseConnectionDataWarehouseConfigResourceModel struct {
	Configuration     DatabaseConnectionDataWarehouseConfigConfigurationResourceModel      `tfsdk:"configuration"`
	ExternalDatabases []DatabaseConnectionDataWarehouseConfigExternalDatabaseResourceModel `tfsdk:"external_databases"`
}

type DatabaseConnectionDataWarehouseConfigConfigurationResourceModel struct {
	AccountName       types.String `tfsdk:"account_name"`
	User              types.String `tfsdk:"user"`
	Password          types.String `tfsdk:"password"`
	PrivateKey        types.String `tfsdk:"private_key"`
	Role              types.String `tfsdk:"role"`
	Warehouse         types.String `tfsdk:"warehouse"`
	Database          types.String `tfsdk:"database"`
	OauthClientId     types.String `tfsdk:"oauth_client_id"`
	OauthClientSecret types.String `tfsdk:"oauth_client_secret"`
	Scope             types.String `tfsdk:"scope"`
	AuthUrl           types.String `tfsdk:"auth_url"`
	AccessTokenUrl    types.String `tfsdk:"access_token_url"`
}

type DatabaseConnectionDataWarehouseConfigExternalDatabaseResourceModel struct {
	Name types.String `tfsdk:"name"`
}

// Metadata returns the resource type name.
func (r *DatabaseConnectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_connection"
}

// Schema defines the schema for the resource.
func (r *DatabaseConnectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"data_warehouse_type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"SNOWFLAKE"}...),
				},
			},
			"data_warehouse_config": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"configuration": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"account_name": schema.StringAttribute{
								Optional: true,
							},
							"user": schema.StringAttribute{
								Optional: true,
							},
							"password": schema.StringAttribute{
								Optional:  true,
								Sensitive: true,
							},
							"private_key": schema.StringAttribute{
								Optional:  true,
								Sensitive: true,
							},
							"role": schema.StringAttribute{
								Optional: true,
							},
							"warehouse": schema.StringAttribute{
								Optional: true,
							},
							"database": schema.StringAttribute{
								Optional: true,
							},
							"oauth_client_id": schema.StringAttribute{
								Optional: true,
							},
							"oauth_client_secret": schema.StringAttribute{
								Optional:  true,
								Sensitive: true,
							},
							"scope": schema.StringAttribute{
								Optional: true,
							},
							"auth_url": schema.StringAttribute{
								Optional: true,
							},
							"access_token_url": schema.StringAttribute{
								Optional: true,
							},
						},
					},
					"external_databases": schema.ListNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required: true,
								},
							},
						},
					},
				},
			},
			"validate": schema.BoolAttribute{
				Required: true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *DatabaseConnectionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *DatabaseConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan DatabaseConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cr := models.CreateConnectionRequest{
		Name:                plan.Name.ValueString(),
		Description:         plan.Description.ValueString(),
		DataWarehouseType:   plan.DataWarehouseType.ValueString(),
		DataWarehouseConfig: map[string]interface{}{},
		Validate:            plan.Validate.ValueBool(),
	}

	cr.DataWarehouseConfig["configuration"] = map[string]interface{}{
		"account_name":        plan.DataWarehouseConfig.Configuration.AccountName.ValueString(),
		"user":                plan.DataWarehouseConfig.Configuration.User.ValueString(),
		"password":            plan.DataWarehouseConfig.Configuration.Password.ValueString(),
		"private_key":         plan.DataWarehouseConfig.Configuration.PrivateKey.ValueString(),
		"role":                plan.DataWarehouseConfig.Configuration.Role.ValueString(),
		"warehouse":           plan.DataWarehouseConfig.Configuration.Warehouse.ValueString(),
		"database":            plan.DataWarehouseConfig.Configuration.Database.ValueString(),
		"oauth_client_id":     plan.DataWarehouseConfig.Configuration.OauthClientId.ValueString(),
		"oauth_client_secret": plan.DataWarehouseConfig.Configuration.OauthClientSecret.ValueString(),
		"scope":               plan.DataWarehouseConfig.Configuration.Scope.ValueString(),
		"auth_url":            plan.DataWarehouseConfig.Configuration.AuthUrl.ValueString(),
		"access_token_url":    plan.DataWarehouseConfig.Configuration.AccessTokenUrl.ValueString(),
	}

	edl := []map[string]interface{}{}

	for _, source := range plan.DataWarehouseConfig.ExternalDatabases {
		ed := map[string]interface{}{
			"name": source.Name.ValueString(),
		}
		edl = append(edl, ed)
	}
	cr.DataWarehouseConfig["externalDatabases"] = edl

	// Create new space
	c, err := r.client.CreateConnection(cr)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating data connection",
			"Could not create connection, unexpected error: "+err.Error(),
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
func (r *DatabaseConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state DatabaseConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cr := models.SearchConnectionRequest{
		Connections: []models.ConnectionInput{
			{
				Identifier: state.ID.ValueString(),
			}},
	}

	c, err := r.client.SearchConnection(cr)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Database Connection",
			"Could not read Database Connection ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	if len(c) == 0 {
		resp.State.RemoveResource(ctx)

		return
	}

	conn := c[0]

	state.Name = types.StringValue(conn.Name)
	state.Description = types.StringValue(conn.Description)
	state.DataWarehouseType = types.StringValue(conn.DataWarehouseType)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DatabaseConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan DatabaseConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cr := models.UpdateConnectionRequest{
		Name:                plan.Name.ValueString(),
		Description:         plan.Description.ValueString(),
		DataWarehouseConfig: map[string]interface{}{},
		Validate:            plan.Validate.ValueBool(),
	}

	cr.DataWarehouseConfig["configuration"] = map[string]interface{}{
		"account_name":        plan.DataWarehouseConfig.Configuration.AccountName.ValueString(),
		"user":                plan.DataWarehouseConfig.Configuration.User.ValueString(),
		"password":            plan.DataWarehouseConfig.Configuration.Password.ValueString(),
		"private_key":         plan.DataWarehouseConfig.Configuration.PrivateKey.ValueString(),
		"role":                plan.DataWarehouseConfig.Configuration.Role.ValueString(),
		"warehouse":           plan.DataWarehouseConfig.Configuration.Warehouse.ValueString(),
		"database":            plan.DataWarehouseConfig.Configuration.Database.ValueString(),
		"oauth_client_id":     plan.DataWarehouseConfig.Configuration.OauthClientId.ValueString(),
		"oauth_client_secret": plan.DataWarehouseConfig.Configuration.OauthClientSecret.ValueString(),
		"scope":               plan.DataWarehouseConfig.Configuration.Scope.ValueString(),
		"auth_url":            plan.DataWarehouseConfig.Configuration.AuthUrl.ValueString(),
		"access_token_url":    plan.DataWarehouseConfig.Configuration.AccessTokenUrl.ValueString(),
	}

	edl := []map[string]interface{}{}

	for _, source := range plan.DataWarehouseConfig.ExternalDatabases {
		ed := map[string]interface{}{
			"name": source.Name.ValueString(),
		}
		edl = append(edl, ed)
	}
	cr.DataWarehouseConfig["externalDatabases"] = edl

	// Create new space
	err := r.client.UpdateConnection(plan.ID.ValueString(), cr)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Database Connection",
			"Could not update Database Connection, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DatabaseConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state DatabaseConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteConnection(state.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Database Connection",
			"Could not delete Database Connection, unexpected error: "+err.Error(),
		)
		return
	}
}

// func (r *DatabaseConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
// 	// Retrieve import ID and save to id attribute
// 	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
// }
