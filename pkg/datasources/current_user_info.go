package datasources

import (
	"context"
	"fmt"

	thoughtspot "github.com/daniepett/thoughtspot-sdk-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &CurrentUserInfoDataSource{}
	_ datasource.DataSourceWithConfigure = &CurrentUserInfoDataSource{}
)

// NewSpacesDataSource is a helper function to simplify the provider implementation.
func NewCurrentUserInfoDataSource() datasource.DataSource {
	return &CurrentUserInfoDataSource{}
}

// spacesDataSource is the data source implementation.
type CurrentUserInfoDataSource struct {
	client *thoughtspot.Client
}

// spacesModel maps coffees schema data.
type CurrentUserInfoModel struct {
	Id types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`
	// Display name of the user.
	DisplayName types.String `tfsdk:"display_name"`
	// Visibility of the users. The `SHARABLE` property makes a user visible to other users and group, who can share objects with the user.
	Visibility types.String `tfsdk:"visibility"`
	// Unique identifier of author of the user.
	AuthorId types.String `tfsdk:"author_id"`
	// Defines whether the user can change their password.
	CanChangePassword types.Bool `tfsdk:"can_change_password"`
	// Defines whether the response has complete detail of the user.
	CompleteDetail types.Bool `tfsdk:"complete_detail"`
	// Creation time of the user in milliseconds.
	CreationTimeInMillis types.Float32 `tfsdk:"creation_time_in_millis"`
}

// Metadata returns the data source type name.
func (d *CurrentUserInfoDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_current_user_info"
}

// Schema defines the schema for the data source.
func (d *CurrentUserInfoDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"display_name": schema.StringAttribute{
							Computed: true,
						},
						"visibility": schema.StringAttribute{
							Computed: true,
						},
						"author_id": schema.StringAttribute{
							Computed: true,
						},
						"can_change_password": schema.BoolAttribute{
							Computed: true,
						},
						"complete_detail": schema.BoolAttribute{
							Computed: true,
						},
						"creation_time_in_millis": schema.Float32Attribute{
							Computed: true,
						},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *CurrentUserInfoDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state CurrentUserInfoModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	User, err := d.client.GetCurrentUserInfo()

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Qlik Cloud User",
			err.Error(),
		)
		return
	}

	state.Id = types.StringValue(User.Id)
	state.Name = types.StringValue(User.Name)
	state.DisplayName = types.StringValue(User.DisplayName)
	state.Visibility = types.StringValue(User.Visibility)
	state.AuthorId = types.StringValue(User.AuthorId)
	state.CanChangePassword = types.BoolValue(User.CanChangePassword)
	state.CompleteDetail = types.BoolValue(User.CompleteDetail)
	state.CreationTimeInMillis = types.Float32Value(User.CreationTimeInMillis)

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *CurrentUserInfoDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = client
}
