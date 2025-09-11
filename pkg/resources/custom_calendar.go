package resources

import (
	"context"
	"fmt"

	thoughtspot "github.com/daniepett/thoughtspot-sdk-go"
	"github.com/daniepett/thoughtspot-sdk-go/models"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
	_ resource.Resource              = &CustomCalendarResource{}
	_ resource.ResourceWithConfigure = &CustomCalendarResource{}
	// _ resource.ResourceWithImportState = &spaceResource{}
)

func NewCustomCalendarResource() resource.Resource {
	return &CustomCalendarResource{}
}

type CustomCalendarResource struct {
	client *thoughtspot.Client
}

type CustomCalendarResourceModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	ExistingTable     types.Bool   `tfsdk:"existing_table"`
	TableReference    types.Object `tfsdk:"table_reference"`
	StartDate         types.String `tfsdk:"start_date"`
	EndDate           types.String `tfsdk:"end_date"`
	CalendarType      types.String `tfsdk:"calendar_type"`
	MonthOffset       types.String `tfsdk:"month_offset"`
	StartDayOfWeek    types.String `tfsdk:"start_day_of_week"`
	QuarterNamePrefix types.String `tfsdk:"quarter_name_prefix"`
	YearNamePrefix    types.String `tfsdk:"year_name_prefix"`
}

type CustomCalendarTableReferenceModel struct {
	ConnectionIdentifier types.String `tfsdk:"connection_identifier"`
	DatabaseName         types.String `tfsdk:"database_name"`
	SchemaName           types.String `tfsdk:"schema_name"`
	TableName            types.String `tfsdk:"table_name"`
}

// CustomCalendar returns the resource type name.
func (r *CustomCalendarResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_calendar"
}

// Schema defines the schema for the resource.
func (r *CustomCalendarResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the custom calendar.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"existing_table": schema.BoolAttribute{
				Required:    true,
				Description: "Defines the creation method",
			},
			"start_date": schema.StringAttribute{
				Optional:    true,
				Description: "Start date for the calendar in MM/dd/yyyy format.",
			},
			"end_date": schema.StringAttribute{
				Optional:    true,
				Description: "End date for the calendar in MM/dd/yyyy format.",
			},
			"calendar_type": schema.StringAttribute{
				Optional:    true,
				Description: "Type of the calendar.",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"MONTH_OFFSET",
						"FOUR_FOUR_FIVE",
						"FOUR_FIVE_FOUR",
						"FIVE_FOUR_FOUR"}...),
				},
			},
			"month_offset": schema.StringAttribute{
				Optional:    true,
				Description: "Specify the month in which the fiscal or custom calendar year should start. For example, if you set month_offset to \"April\", the custom calendar will treat \"April\" as the first month of the year, and the related attributes such as quarters and start date will be based on this offset. The default value is January, which represents the standard calendar year (January to December).",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"January",
						"February",
						"March",
						"April",
						"May",
						"June",
						"July",
						"August",
						"September",
						"October",
						"November",
						"December"}...),
				},
			},
			"start_day_of_week": schema.StringAttribute{
				Optional:    true,
				Description: "Specify the starting day of the week.",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"Sunday",
						"Monday",
						"Tuesday",
						"Wednesday",
						"Thursday",
						"Friday",
						"Saturday"}...),
				},
			},
			"quarter_name_prefix": schema.StringAttribute{
				Optional:    true,
				Description: "Prefix to add before the quarter.",
			},
			"year_name_prefix": schema.StringAttribute{
				Optional:    true,
				Description: "Prefix to add before the year.",
			},
		},
		Blocks: map[string]schema.Block{
			"table_reference": schema.SingleNestedBlock{
				Validators: []validator.Object{
					objectvalidator.IsRequired(),
				},
				Attributes: map[string]schema.Attribute{
					"connection_identifier": schema.StringAttribute{
						Optional:    true,
						Description: "Unique ID or name of the connection.",
					},
					"database_name": schema.StringAttribute{
						Optional:    true,
						Description: "Name of the database.",
					},
					"schema_name": schema.StringAttribute{
						Optional:    true,
						Description: "Name of the schema.",
					},
					"table_name": schema.StringAttribute{
						Optional:    true,
						Description: "Name of the table. Table names may be case-sensitive depending on the database system.",
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *CustomCalendarResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *CustomCalendarResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan CustomCalendarResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var cm string
	if plan.ExistingTable.ValueBool() {
		cm = "FROM_EXISTING_TABLE"
	} else {
		cm = "FROM_INPUT_PARAMS"
	}

	var tr CustomCalendarTableReferenceModel

	diags = plan.TableReference.As(ctx, &tr, basetypes.ObjectAsOptions{})
	resp.Diagnostics.Append(diags...)

	cr := models.CreateCustomCalendarRequest{
		Name:           plan.Name.ValueString(),
		CreationMethod: cm,
		TableReference: models.CustomCalendarTableReference{
			ConnectionIdentifier: tr.ConnectionIdentifier.ValueString(),
			DatabaseName:         tr.DatabaseName.ValueString(),
			SchemaName:           tr.SchemaName.ValueString(),
			TableName:            tr.TableName.ValueString(),
		},
		StartDate:         plan.StartDate.ValueString(),
		EndDate:           plan.EndDate.ValueString(),
		CalendarType:      plan.CalendarType.ValueString(),
		MonthOffset:       plan.MonthOffset.ValueString(),
		StartDayOfWeek:    plan.StartDayOfWeek.ValueString(),
		QuarterNamePrefix: plan.QuarterNamePrefix.ValueString(),
		YearNamePrefix:    plan.YearNamePrefix.ValueString(),
	}

	c, err := r.client.CreateCalendar(cr)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating user group",
			"Could not create user group, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(c.CalendarId)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *CustomCalendarResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state CustomCalendarResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cr := models.SearchCustomCalendarsRequest{
		NamePattern: state.Name.ValueString(),
		RecordSize:  "1",
	}

	c, err := r.client.SearchCalendars(cr)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Custom Calendar",
			"Could not read Custom Calendars "+state.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	if len(c) == 0 {
		resp.State.RemoveResource(ctx)
		return
	}

	cal := c[0]

	state.Name = types.StringValue(cal.CalendarName)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *CustomCalendarResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan CustomCalendarResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var cm string
	if plan.ExistingTable.ValueBool() {
		cm = "FROM_EXISTING_TABLE"
	} else {
		cm = "FROM_INPUT_PARAMS"
	}

	var tr CustomCalendarTableReferenceModel

	diags = plan.TableReference.As(ctx, &tr, basetypes.ObjectAsOptions{})
	resp.Diagnostics.Append(diags...)

	cr := models.UpdateCustomCalendarRequest{
		UpdateMethod: cm,
		TableReference: models.CustomCalendarTableReference{
			ConnectionIdentifier: tr.ConnectionIdentifier.ValueString(),
			DatabaseName:         tr.DatabaseName.ValueString(),
			SchemaName:           tr.SchemaName.ValueString(),
			TableName:            tr.TableName.ValueString(),
		},
		StartDate:         plan.StartDate.ValueString(),
		EndDate:           plan.EndDate.ValueString(),
		CalendarType:      plan.CalendarType.ValueString(),
		MonthOffset:       plan.MonthOffset.ValueString(),
		StartDayOfWeek:    plan.StartDayOfWeek.ValueString(),
		QuarterNamePrefix: plan.QuarterNamePrefix.ValueString(),
		YearNamePrefix:    plan.YearNamePrefix.ValueString(),
	}

	err := r.client.UpdateCalendar(plan.ID.ValueString(), cr)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Custom Calendar",
			"Could not update Custom Calendar, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *CustomCalendarResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state CustomCalendarResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCalendar(state.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Custom Calendar",
			"Could not delete custom calendar, unexpected error: "+err.Error(),
		)
		return
	}
}

// func (r *CustomCalendarResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
// 	// Retrieve import ID and save to id attribute
// 	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
// }
