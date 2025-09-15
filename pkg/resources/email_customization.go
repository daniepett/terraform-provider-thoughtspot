package resources

import (
	"context"
	"fmt"

	thoughtspot "github.com/daniepett/thoughtspot-sdk-go"
	"github.com/daniepett/thoughtspot-sdk-go/models"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &EmailCustomizationResource{}
	_ resource.ResourceWithConfigure = &EmailCustomizationResource{}
	// _ resource.ResourceWithImportState = &spaceResource{}
)

func NewEmailCustomizationResource() resource.Resource {
	return &EmailCustomizationResource{}
}

type EmailCustomizationResource struct {
	client *thoughtspot.Client
}

type EmailCustomizationResourceModel struct {
	ID                           types.String `tfsdk:"id"`
	OrgIdentifier                types.String `tfsdk:"org_identifier"`
	CtaButtonBgColor             types.String `tfsdk:"cta_button_bg_color"`
	CtaTextFontColor             types.String `tfsdk:"cta_text_font_color"`
	PrimaryBgColor               types.String `tfsdk:"primary_bg_color"`
	HomeURL                      types.String `tfsdk:"home_url"`
	LogoURL                      types.String `tfsdk:"logo_url"`
	FontFamily                   types.String `tfsdk:"font_family"`
	ProductName                  types.String `tfsdk:"product_name"`
	FooterAddress                types.String `tfsdk:"footer_address"`
	FooterPhone                  types.String `tfsdk:"footer_phone"`
	ReplacementValueForLiveboard types.String `tfsdk:"replacement_value_for_liveboard"`
	ReplacementValueForAnswer    types.String `tfsdk:"replacement_value_for_answer"`
	ReplacementValueForSpotIQ    types.String `tfsdk:"replacement_value_for_spot_iq"`
	HideFooterAddress            types.Bool   `tfsdk:"hide_footer_address"`
	HideFooterPhone              types.Bool   `tfsdk:"hide_footer_phone"`
	HideManageNotification       types.Bool   `tfsdk:"hide_manage_notification"`
	HideMobileAppNudge           types.Bool   `tfsdk:"hide_mobile_app_nudge"`
	HidePrivacyPolicy            types.Bool   `tfsdk:"hide_privacy_policy"`
	HideProductName              types.Bool   `tfsdk:"hide_product_name"`
	HideTSVocabularyDefinitions  types.Bool   `tfsdk:"hide_ts_vocabulary_definitions"`
	HideNotificationStatus       types.Bool   `tfsdk:"hide_notification_status"`
	HideErrorMessage             types.Bool   `tfsdk:"hide_error_message"`
	HideUnsubscribeLink          types.Bool   `tfsdk:"hide_unsubscribe_link"`
	HideModifyAlert              types.Bool   `tfsdk:"hide_modify_alert"`
	ValidateCustomization        types.Bool   `tfsdk:"validate_customization"`
}

// EmailCustomization returns the resource type name.
func (r *EmailCustomizationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_email_customization"
}

// Schema defines the schema for the resource.
func (r *EmailCustomizationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"org_identifier": schema.StringAttribute{
				Optional:    true,
				Description: "Unique ID or name of org",
			},
			"cta_button_bg_color": schema.StringAttribute{
				Optional:    true,
				Description: "Background color for call-to-action button in hex format",
			},
			"cta_text_font_color": schema.StringAttribute{
				Optional:    true,
				Description: "Text color for call-to-action button in hex format",
			}, "primary_bg_color": schema.StringAttribute{
				Optional:    true,
				Description: "Primary background color in hex format",
			},
			"home_url": schema.StringAttribute{
				Optional: true, Description: "Home page URL (HTTP/HTTPS only)",
			},
			"logo_url": schema.StringAttribute{
				Optional:    true,
				Description: "Logo image URL (HTTP/HTTPS only)",
			},
			"font_family": schema.StringAttribute{
				Optional:    true,
				Description: "Font family for email content (e.g., Arial, sans-serif)",
			},
			"product_name": schema.StringAttribute{Optional: true,
				Description: "Product name to display",
			},
			"footer_address": schema.StringAttribute{
				Optional:    true,
				Description: "Footer address text",
			},
			"footer_phone": schema.StringAttribute{
				Optional:    true,
				Description: "Footer phone number",
			}, "replacement_value_for_liveboard": schema.StringAttribute{
				Optional:    true,
				Description: "Replacement value for Liveboard",
			},
			"replacement_value_for_answer": schema.StringAttribute{
				Optional: true, Description: "Replacement value for Answer",
			},
			"replacement_value_for_spot_iq": schema.StringAttribute{
				Optional:    true,
				Description: "Replacement value for SpotIQ",
			},
			"hide_footer_address": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to hide footer address",
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				}},
			"hide_footer_phone": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to hide footer phone number",
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				}},
			"hide_manage_notification": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to hide manage notification link",
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				}},
			"hide_mobile_app_nudge": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to hide mobile app nudge",
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"hide_privacy_policy": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to hide privacy policy link",
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"hide_product_name": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to hide product name",
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"hide_ts_vocabulary_definitions": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to hide ThoughtSpot vocabulary definitions",
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			}, "hide_notification_status": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to hide notification status",
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"hide_error_message": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to hide error message",
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"hide_unsubscribe_link": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to hide unsubscribe link",
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"hide_modify_alert": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to hide modify alert",
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"validate_customization": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to send validation email to the logged in user",
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *EmailCustomizationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *EmailCustomizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan EmailCustomizationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tp := models.TemplatePropertiesInputCreate{
		CtaButtonBgColor:             plan.CtaButtonBgColor.ValueString(),
		CtaTextFontColor:             plan.CtaTextFontColor.ValueString(),
		PrimaryBgColor:               plan.PrimaryBgColor.ValueString(),
		HomeURL:                      plan.HomeURL.ValueString(),
		LogoURL:                      plan.LogoURL.ValueString(),
		FontFamily:                   plan.FontFamily.ValueString(),
		ProductName:                  plan.ProductName.ValueString(),
		FooterAddress:                plan.FooterAddress.ValueString(),
		FooterPhone:                  plan.FooterPhone.ValueString(),
		ReplacementValueForLiveboard: plan.ReplacementValueForLiveboard.ValueString(),
		ReplacementValueForAnswer:    plan.ReplacementValueForAnswer.ValueString(),
		ReplacementValueForSpotIQ:    plan.ReplacementValueForSpotIQ.ValueString(),
		HideFooterAddress:            plan.HideFooterAddress.ValueBool(),
		HideFooterPhone:              plan.HideFooterPhone.ValueBool(),
		HideManageNotification:       plan.HideManageNotification.ValueBool(),
		HideMobileAppNudge:           plan.HideMobileAppNudge.ValueBool(),
		HidePrivacyPolicy:            plan.HidePrivacyPolicy.ValueBool(),
		HideProductName:              plan.HideProductName.ValueBool(),
		HideTSVocabularyDefinitions:  plan.HideTSVocabularyDefinitions.ValueBool(),
		HideNotificationStatus:       plan.HideNotificationStatus.ValueBool(),
		HideErrorMessage:             plan.HideErrorMessage.ValueBool(),
		HideUnsubscribeLink:          plan.HideUnsubscribeLink.ValueBool(),
		HideModifyAlert:              plan.HideModifyAlert.ValueBool(),
	}

	cr := models.CreateEmailCustomizationRequest{
		OrgIdentifier:      plan.OrgIdentifier.ValueString(),
		TemplateProperties: tp,
	}

	c, err := r.client.CreateEmailCustomization(cr)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating email customization",
			"Could not create email customization, unexpected error: "+err.Error(),
		)
		return
	}

	if plan.ValidateCustomization.ValueBool() {
		err := r.client.ValidateEmailCustomization()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error validation email customization",
				"Could not validate email customization, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(fmt.Sprintf("%d", c.Org.Id))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *EmailCustomizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state EmailCustomizationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cr := models.CustomizationEmailSearchRequest{
		OrgIdentifiers: []string{state.ID.ValueString()},
	}

	c, err := r.client.SearchEmailCustomization(cr)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Email Customization",
			"Could not read Email Customization"+err.Error(),
		)
		return
	}

	if len(c) == 0 {
		resp.State.RemoveResource(ctx)
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *EmailCustomizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan EmailCustomizationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tp := models.TemplatePropertiesInputCreate{
		CtaButtonBgColor:             plan.CtaButtonBgColor.ValueString(),
		CtaTextFontColor:             plan.CtaTextFontColor.ValueString(),
		PrimaryBgColor:               plan.PrimaryBgColor.ValueString(),
		HomeURL:                      plan.HomeURL.ValueString(),
		LogoURL:                      plan.LogoURL.ValueString(),
		FontFamily:                   plan.FontFamily.ValueString(),
		ProductName:                  plan.ProductName.ValueString(),
		FooterAddress:                plan.FooterAddress.ValueString(),
		FooterPhone:                  plan.FooterPhone.ValueString(),
		ReplacementValueForLiveboard: plan.ReplacementValueForLiveboard.ValueString(),
		ReplacementValueForAnswer:    plan.ReplacementValueForAnswer.ValueString(),
		ReplacementValueForSpotIQ:    plan.ReplacementValueForSpotIQ.ValueString(),
		HideFooterAddress:            plan.HideFooterAddress.ValueBool(),
		HideFooterPhone:              plan.HideFooterPhone.ValueBool(),
		HideManageNotification:       plan.HideManageNotification.ValueBool(),
		HideMobileAppNudge:           plan.HideMobileAppNudge.ValueBool(),
		HidePrivacyPolicy:            plan.HidePrivacyPolicy.ValueBool(),
		HideProductName:              plan.HideProductName.ValueBool(),
		HideTSVocabularyDefinitions:  plan.HideTSVocabularyDefinitions.ValueBool(),
		HideNotificationStatus:       plan.HideNotificationStatus.ValueBool(),
		HideErrorMessage:             plan.HideErrorMessage.ValueBool(),
		HideUnsubscribeLink:          plan.HideUnsubscribeLink.ValueBool(),
		HideModifyAlert:              plan.HideModifyAlert.ValueBool(),
	}

	cr := models.UpdateEmailCustomizationRequest{
		OrgIdentifier:      plan.ID.ValueString(),
		TemplateProperties: tp,
	}

	err := r.client.UpdateEmailCustomization(cr)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating email customization",
			"Could not update email customization, unexpected error: "+err.Error(),
		)
		return
	}

	if plan.ValidateCustomization.ValueBool() {
		err := r.client.ValidateEmailCustomization()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error validation email customization",
				"Could not validate email customization, unexpected error: "+err.Error(),
			)
			return
		}
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *EmailCustomizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state EmailCustomizationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cr := models.CustomizationEmailDeleteRequest{
		OrgIdentifiers: []string{state.ID.ValueString()},
	}

	err := r.client.DeleteOrgEmailCustomization(cr)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting email customization",
			"Could not delete email customization, unexpected error: "+err.Error(),
		)
		return
	}
}

// func (r *EmailCustomizationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
// 	// Retrieve import ID and save to id attribute
// 	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
// }
