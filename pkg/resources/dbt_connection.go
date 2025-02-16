package resources

import (
	"context"
	"fmt"
	"time"

	thoughtspot "github.com/daniepett/thoughtspot-sdk-go"
	"github.com/daniepett/thoughtspot-sdk-go/models"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
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
	_ resource.Resource              = &DbtConnectionResource{}
	_ resource.ResourceWithConfigure = &DbtConnectionResource{}
	// _ resource.ResourceWithImportState = &spaceResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewDbtConnectionResource() resource.Resource {
	return &DbtConnectionResource{}
}

// orderResource is the resource implementation.
type DbtConnectionResource struct {
	client *thoughtspot.Client
}

// orderResourceModel maps the resource schema data.
type DbtConnectionResourceModel struct {
	ID             types.String   `tfsdk:"id"`
	ConnectionName types.String   `tfsdk:"connection_name"`
	DatabaseName   types.String   `tfsdk:"database_name"`
	ImportType     types.String   `tfsdk:"import_type"`
	AccessToken    types.String   `tfsdk:"access_token"`
	DbtUrl         types.String   `tfsdk:"dbt_url"`
	AccountId      types.String   `tfsdk:"account_id"`
	ProjectId      types.String   `tfsdk:"project_id"`
	DbtEnvId       types.String   `tfsdk:"dbt_env_id"`
	ProjectName    types.String   `tfsdk:"project_name"`
	FileContent    types.String   `tfsdk:"file_content"`
	Timeouts       timeouts.Value `tfsdk:"timeouts"`
}

// Metadata returns the resource type name.
func (r *DbtConnectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dbt_connection"
}

// Schema defines the schema for the resource.
func (r *DbtConnectionResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"connection_name": schema.StringAttribute{
				Required: true,
			},
			"database_name": schema.StringAttribute{
				Required: true,
			},
			"import_type": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"DBT_CLOUD",
						"ZIP_FILE"}...),
				},
			},
			"access_token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"dbt_url": schema.StringAttribute{
				Optional: true,
			},
			"account_id": schema.StringAttribute{
				Optional: true,
			},
			"project_id": schema.StringAttribute{
				Optional: true,
			},
			"dbt_env_id": schema.StringAttribute{
				Optional: true,
			},
			"project_name": schema.StringAttribute{
				Optional: true,
			},
			"file_content": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
			}),
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *DbtConnectionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *DbtConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan DbtConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	cr := models.DbtConnectionRequest{
		ConnectionName: plan.ConnectionName.ValueString(),
		DatabaseName:   plan.DatabaseName.ValueString(),
		ImportType:     plan.ImportType.ValueString(),
		AccessToken:    plan.AccessToken.ValueString(),
		DbtUrl:         plan.DbtUrl.ValueString(),
		AccountId:      plan.AccountId.ValueString(),
		ProjectId:      plan.ProjectId.ValueString(),
		DbtEnvId:       plan.DbtEnvId.ValueString(),
		ProjectName:    plan.ProjectName.ValueString(),
		FileContent:    plan.FileContent.ValueString(),
	}

	createTimeout, diags := plan.Timeouts.Create(ctx, 2*time.Minute)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	_, err := r.client.CreateDbtConnection(cr)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating dbt Connection",
			"Could not create Dbt Connection, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	// plan.ID = types.StringValue(c.Id)
	// plan.DataWarehouseType = types.StringValue(dtype)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *DbtConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state DbtConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.SearchDbtConnection()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading dbt Connections",
			"Could not read Database DbtConnection ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// if len(c) == 0 {
	// 	resp.State.RemoveResource(ctx)

	// 	return
	// }

	// conn := c[0]

	// state.Name = types.StringValue(conn.Name)
	// state.Description = types.StringValue(conn.Description)
	// state.DataWarehouseType = types.StringValue(conn.DataWarehouseType)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DbtConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan DbtConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// var config map[string]interface{}

	// if !plan.Snowflake.IsNull() {
	// 	var cs DbtConnectionSnowflakeModel

	// 	diags = plan.Snowflake.As(ctx, &cs, basetypes.ObjectAsOptions{})
	// 	resp.Diagnostics.Append(diags...)

	// 	config = map[string]interface{}{
	// 		"authenticationType": cs.AuthenticationType.ValueString(),
	// 		"accountName":        cs.AccountName.ValueString(),
	// 		"password":           cs.Password.ValueString(),
	// 		// "PrivateKey":        cs.PrivateKey.ValueString(),
	// 		"role":             cs.Role.ValueString(),
	// 		"warehouse":        cs.Warehouse.ValueString(),
	// 		"database":         cs.Database.ValueString(),
	// 		"client_id":        cs.OauthClientId.ValueString(),
	// 		"client_secret":    cs.OauthClientSecret.ValueString(),
	// 		"scope":            cs.Scope.ValueString(),
	// 		"auth_url":         cs.AuthUrl.ValueString(),
	// 		"access_token_url": cs.AccessTokenUrl.ValueString(),
	// 	}

	// 	if !cs.User.IsNull() {
	// 		config["user"] = cs.User.ValueString()
	// 	}

	// }

	// edl := []map[string]interface{}{}

	// for _, source := range plan.ExternalDatabases {
	// 	ed := map[string]interface{}{
	// 		"name": source.Name.ValueString(),
	// 	}
	// 	edl = append(edl, ed)
	// }

	// cr := models.UpdateDbtConnectionRequest{
	// 	Name:        plan.Name.ValueString(),
	// 	Description: plan.Description.ValueString(),
	// 	DataWarehouseConfig: map[string]interface{}{
	// 		"configuration":     config,
	// 		"externalDatabases": edl,
	// 	},
	// 	Validate: plan.Validate.ValueBool(),
	// }

	// // Create new space
	// err := r.client.UpdateDbtConnection(plan.ID.ValueString(), cr)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error updating Database DbtConnection",
	// 		"Could not update Database DbtConnection, unexpected error: "+err.Error(),
	// 	)
	// 	return
	// }

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DbtConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state DbtConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDbtConnection(state.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Database DbtConnection",
			"Could not delete Database DbtConnection, unexpected error: "+err.Error(),
		)
		return
	}
}

// func (r *DbtConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
// 	// Retrieve import ID and save to id attribute
// 	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
// }
