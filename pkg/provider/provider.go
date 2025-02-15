package provider

import (
	"context"
	"os"

	"terraform-provider-thoughtspot/pkg/datasources"
	"terraform-provider-thoughtspot/pkg/resources"

	thoughtspot "github.com/daniepett/thoughtspot-sdk-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &thoughtspotProvider{}
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &thoughtspotProvider{
			version: version,
		}
	}
}

type thoughtspotProvider struct {
	version string
}

type thoughtspotProviderModel struct {
	Host          types.String `tfsdk:"host"`
	Username      types.String `tfsdk:"username"`
	Password      types.String `tfsdk:"password"`
	OrgIdentifier types.String `tfsdk:"org_identifier"`
}

func (p *thoughtspotProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "thoughtspot"
	resp.Version = p.version
}

func (p *thoughtspotProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Required: true,
			},
			"username": schema.StringAttribute{
				Required: true,
			},
			"password": schema.StringAttribute{
				Required: true,
			},
			"org_identifier": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (p *thoughtspotProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config thoughtspotProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("THOUGHTSPOT_HOST")
	username := os.Getenv("THOUGHTSPOT_USERNAME")
	password := os.Getenv("THOUGHTSPOT_PASSWORD")
	org_identifier := os.Getenv("THOUGHTSPOT_ORG_IDENTIFIER")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}
	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	if !config.OrgIdentifier.IsNull() {
		org_identifier = config.OrgIdentifier.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Host",
			"The provider cannot create the ThoughtSpot API client as there is a missing or empty value for the ThoughtSpot Host. "+
				"Set the host value in the configuration or use the THOUGHTSPOT_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing Username",
			"The provider cannot create the ThoughtSpot API client as there is a missing or empty value for the ThoughtSpot Username. "+
				"Set the username value in the configuration or use the THOUGHTSPOT_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Password",
			"The provider cannot create the ThoughtSpot API client as there is a missing or empty value for the ThoughtSpot Password. "+
				"Set the password value in the configuration or use the THOUGHTSPOT_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if org_identifier == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("org_identifier"),
			"Missing OrgIdentifier",
			"The provider cannot create the ThoughtSpot API client as there is a missing or empty value for the ThoughtSpot OrgIdentifier. "+
				"Set the org_identifier value in the configuration or use the THOUGHTSPOT_ORG_IDENTIFIER environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new Thoughtspot Client using the configuration values
	client, err := thoughtspot.NewClient(&host, &username, &password, &org_identifier)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Thoughtspot Client",
			"An unexpected error occurred when creating the Thoughtspot Client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Thoughtspot Client Error: "+err.Error(),
		)
		return
	}

	// Make the Qlik client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *thoughtspotProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewCurrentUserInfoDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *thoughtspotProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewUserGroupResource,
		resources.NewMetadataResource,
		resources.NewConnectionResource,
		resources.NewRoleResource,
		resources.NewTmlResource,
	}
}
