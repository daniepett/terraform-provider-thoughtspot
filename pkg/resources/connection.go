package resources

import (
	"context"
	"fmt"

	thoughtspot "github.com/daniepett/thoughtspot-sdk-go"
	"github.com/daniepett/thoughtspot-sdk-go/models"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &ConnectionResource{}
	_ resource.ResourceWithConfigure = &ConnectionResource{}
	// _ resource.ResourceWithImportState = &spaceResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewConnectionResource() resource.Resource {
	return &ConnectionResource{}
}

// orderResource is the resource implementation.
type ConnectionResource struct {
	client *thoughtspot.Client
}

// orderResourceModel maps the resource schema data.
type ConnectionResourceModel struct {
	ID                types.String                                                 `tfsdk:"id"`
	Name              types.String                                                 `tfsdk:"name"`
	Description       types.String                                                 `tfsdk:"description"`
	DataWarehouseType types.String                                                 `tfsdk:"data_warehouse_type"`
	ExternalDatabases []ConnectionDataWarehouseConfigExternalDatabaseResourceModel `tfsdk:"external_databases"`
	Validate          types.Bool                                                   `tfsdk:"validate"`
	Snowflake         types.Object                                                 `tfsdk:"snowflake"`
	Redshift          types.Object                                                 `tfsdk:"redshift"`
}

type ConnectionSnowflakeModel struct {
	AuthenticationType types.String `tfsdk:"authentication_type"`
	AccountName        types.String `tfsdk:"account_name"`
	User               types.String `tfsdk:"user"`
	Password           types.String `tfsdk:"password"`
	PrivateKey         types.String `tfsdk:"private_key"`
	Role               types.String `tfsdk:"role"`
	Warehouse          types.String `tfsdk:"warehouse"`
	Database           types.String `tfsdk:"database"`
	OauthClientId      types.String `tfsdk:"oauth_client_id"`
	OauthClientSecret  types.String `tfsdk:"oauth_client_secret"`
	Scope              types.String `tfsdk:"scope"`
	AuthUrl            types.String `tfsdk:"auth_url"`
	AccessTokenUrl     types.String `tfsdk:"access_token_url"`
}

func (o ConnectionSnowflakeModel) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"authentication_type": types.StringType,
		"account_name":        types.StringType,
		"user":                types.StringType,
		"password":            types.StringType,
		"private_key":         types.StringType,
		"role":                types.StringType,
		"warehouse":           types.StringType,
		"database":            types.StringType,
		"oauth_client_id":     types.StringType,
		"oauth_client_secret": types.StringType,
		"scope":               types.StringType,
		"auth_url":            types.StringType,
		"access_token_url":    types.StringType,
	}
}

type ConnectionDataWarehouseConfigExternalDatabaseResourceModel struct {
	Name types.String `tfsdk:"name"`
}

// Metadata returns the resource type name.
func (r *ConnectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connection"
}

// Schema defines the schema for the resource.
func (r *ConnectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				Computed: true,
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
			"validate": schema.BoolAttribute{
				Required: true,
			},
		},
		Blocks: map[string]schema.Block{
			"snowflake": schema.SingleNestedBlock{
				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("snowflake"),
						path.MatchRoot("redshift"),
					}...),
				},
				Attributes: map[string]schema.Attribute{
					"authentication_type": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf([]string{
								"SERVICE_ACCOUNT",
								"OAUTH",
								"EXTOAUTH",
								"KEY_PAIR",
								"OAUTH_WITH_PKCE",
								"EXTOAUTH_WITH_PKCE"}...),
						},
					},
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
			"redshift": schema.SingleNestedBlock{
				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("snowflake"),
						path.MatchRoot("redshift"),
					}...),
				},
				Attributes: map[string]schema.Attribute{
					"account_name": schema.StringAttribute{
						Optional: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *ConnectionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *ConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dtype string
	var config map[string]interface{}

	if !plan.Snowflake.IsNull() {
		dtype = "SNOWFLAKE"
		var cs ConnectionSnowflakeModel

		diags = plan.Snowflake.As(ctx, &cs, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)

		config = map[string]interface{}{
			"authenticationType": cs.AuthenticationType.ValueString(),
			"accountName":        cs.AccountName.ValueString(),
			// "password":           cs.Password.ValueString(),
			// "privateKey":        cs.PrivateKey.ValueString(),
			// "role":          cs.Role.ValueString(),
			"warehouse":     cs.Warehouse.ValueString(),
			"database":      cs.Database.ValueString(),
			"client_id":     cs.OauthClientId.ValueString(),
			"client_secret": cs.OauthClientSecret.ValueString(),
			// "scope":             cs.Scope.ValueString(),
			// "auth_url":           cs.AuthUrl.ValueString(),
			// "accesstoken_url":    cs.AccessTokenUrl.ValueString(),
		}

		if !cs.User.IsNull() {
			config["user"] = cs.User.ValueString()
		}

	}

	edl := []map[string]interface{}{}

	for _, source := range plan.ExternalDatabases {
		ed := map[string]interface{}{
			"name": source.Name.ValueString(),
		}
		edl = append(edl, ed)
	}

	cr := models.CreateConnectionRequest{
		Name:              plan.Name.ValueString(),
		Description:       plan.Description.ValueString(),
		DataWarehouseType: dtype,
		DataWarehouseConfig: map[string]interface{}{
			"configuration":     config,
			"externalDatabases": edl,
		},
		Validate: plan.Validate.ValueBool(),
	}

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
	plan.DataWarehouseType = types.StringValue(dtype)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *ConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state ConnectionResourceModel
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

func (r *ConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan ConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config map[string]interface{}

	if !plan.Snowflake.IsNull() {
		var cs ConnectionSnowflakeModel

		diags = plan.Snowflake.As(ctx, &cs, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)

		config = map[string]interface{}{
			"authenticationType": cs.AuthenticationType.ValueString(),
			"accountName":        cs.AccountName.ValueString(),
			"password":           cs.Password.ValueString(),
			// "PrivateKey":        cs.PrivateKey.ValueString(),
			"role":             cs.Role.ValueString(),
			"warehouse":        cs.Warehouse.ValueString(),
			"database":         cs.Database.ValueString(),
			"client_id":        cs.OauthClientId.ValueString(),
			"client_secret":    cs.OauthClientSecret.ValueString(),
			"scope":            cs.Scope.ValueString(),
			"auth_url":         cs.AuthUrl.ValueString(),
			"access_token_url": cs.AccessTokenUrl.ValueString(),
		}

		if !cs.User.IsNull() {
			config["user"] = cs.User.ValueString()
		}

	}

	edl := []map[string]interface{}{}

	for _, source := range plan.ExternalDatabases {
		ed := map[string]interface{}{
			"name": source.Name.ValueString(),
		}
		edl = append(edl, ed)
	}

	cr := models.UpdateConnectionRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		DataWarehouseConfig: map[string]interface{}{
			"configuration":     config,
			"externalDatabases": edl,
		},
		Validate: plan.Validate.ValueBool(),
	}

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

func (r *ConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state ConnectionResourceModel
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

// func (r *ConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
// 	// Retrieve import ID and save to id attribute
// 	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
// }
