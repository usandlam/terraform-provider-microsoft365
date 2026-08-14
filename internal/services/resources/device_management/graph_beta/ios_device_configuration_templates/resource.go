package graphBetaIosDeviceConfigurationTemplates

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_management_ios_device_configuration_templates"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &IosDeviceConfigurationTemplatesResource{}
	_ resource.ResourceWithConfigure   = &IosDeviceConfigurationTemplatesResource{}
	_ resource.ResourceWithImportState = &IosDeviceConfigurationTemplatesResource{}
	_ resource.ResourceWithModifyPlan  = &IosDeviceConfigurationTemplatesResource{}
	_ resource.ResourceWithIdentity    = &IosDeviceConfigurationTemplatesResource{}
)

func NewIosDeviceConfigurationTemplatesResource() resource.Resource {
	return &IosDeviceConfigurationTemplatesResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/deviceConfigurations",
	}
}

type IosDeviceConfigurationTemplatesResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *IosDeviceConfigurationTemplatesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *IosDeviceConfigurationTemplatesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *IosDeviceConfigurationTemplatesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *IosDeviceConfigurationTemplatesResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *IosDeviceConfigurationTemplatesResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages iOS/iPadOS configuration templates in Microsoft Intune. " +
			"This resource creates device configurations for iOS/iPadOS devices including " +
			"general device restrictions (with Activation Lock) and trusted root certificates.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the iOS configuration template.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The display name for the iOS configuration template.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Optional description of the resource. Maximum length is 1500 characters.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1500),
				},
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Set of scope tag IDs for this template profile.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			"general_device_restriction": schema.SingleNestedAttribute{
				Optional: true,
				MarkdownDescription: "iOS/iPadOS general device restriction template " +
					"(iosGeneralDeviceConfiguration). Covers passcode, App Store, common device " +
					"restrictions, iCloud/backup, Safari, and the Activation Lock allow-when-supervised toggle. " +
					"This is a curated subset of the full Graph resource surface; additional properties can be " +
					"added over time without breaking the schema.",
				Attributes: map[string]schema.Attribute{
					"activation_lock_allow_when_supervised": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to allow activation lock when the device is in supervised mode. Enables Activation Lock on supervised iOS/iPadOS devices; use the `bypass_activation_lock` device action to remove it later.",
					},
					// Passcode
					"passcode_required": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to require a passcode.",
					},
					"passcode_required_type": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Type of passcode that is required. Possible values: deviceDefault, alphanumeric, numeric.",
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "alphanumeric", "numeric"),
						},
					},
					"passcode_minimum_length": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "Minimum length of passcode. Valid values 4 to 14.",
						Validators: []validator.Int32{
							int32validator.Between(4, 14),
						},
					},
					"passcode_minutes_of_inactivity_before_lock": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "Minutes of inactivity before a passcode is required.",
					},
					"passcode_minutes_of_inactivity_before_screen_timeout": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "Minutes of inactivity before the screen times out.",
					},
					"passcode_expiration_days": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "Number of days before the passcode expires. Valid values 1 to 65535.",
						Validators: []validator.Int32{
							int32validator.Between(1, 65535),
						},
					},
					"passcode_previous_passcode_block_count": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "Number of previous passcodes to block. Valid values 1 to 24.",
						Validators: []validator.Int32{
							int32validator.Between(1, 24),
						},
					},
					"passcode_sign_in_failure_count_before_wipe": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "Number of sign in failures allowed before wiping the device. Valid values 2 to 11.",
						Validators: []validator.Int32{
							int32validator.Between(2, 11),
						},
					},
					"passcode_block_simple": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block simple passcodes.",
					},
					"passcode_block_fingerprint_unlock": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block fingerprint unlock.",
					},
					"passcode_minimum_character_set_count": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "Number of character sets a passcode must contain. Valid values 0 to 4.",
						Validators: []validator.Int32{
							int32validator.Between(0, 4),
						},
					},
					// App Store
					"app_store_blocked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block the user from using the App Store. Requires a supervised device for iOS 13+.",
					},
					"app_store_block_in_app_purchases": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block the user from making in-app purchases.",
					},
					"app_store_require_password": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to require a password when using the App Store.",
					},
					"app_store_block_automatic_downloads": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block the automatic downloading of apps purchased on other devices (supervised, iOS 9.0+).",
					},
					// App lists
					"compliant_app_list_type": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "List type that is in the compliant apps list. Possible values: none, appsInListCompliant, appsNotInListCompliant.",
						Validators: []validator.String{
							stringvalidator.OneOf("none", "appsInListCompliant", "appsNotInListCompliant"),
						},
					},
					"compliant_apps_list": schema.SetNestedAttribute{
						Optional:            true,
						MarkdownDescription: "List of apps in the compliance list (allow or block list, controlled by compliant_app_list_type).",
						NestedObject: schema.NestedAttributeObject{
							Attributes: appListItemAttributes(),
						},
					},
					// Common restrictions
					"camera_blocked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block the user from accessing the camera. Requires a supervised device for iOS 13+.",
					},
					"screen_capture_blocked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block the user from taking screen captures.",
					},
					"siri_blocked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block the user from using Siri.",
					},
					"siri_blocked_when_locked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block the user from using Siri when locked.",
					},
					"airdrop_blocked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to allow AirDrop when the device is in supervised mode.",
					},
					"bluetooth_block_modification": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to allow modification of Bluetooth settings when the device is in supervised mode (iOS 10.0+).",
					},
					"cellular_block_data_roaming": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block data roaming.",
					},
					"cellular_block_voice_roaming": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block voice roaming.",
					},
					"cellular_block_personal_hotspot": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block Personal Hotspot.",
					},
					"device_block_erase_content_and_settings": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to allow the use of the 'Erase all content and settings' option when supervised.",
					},
					"device_block_name_modification": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to allow device name modification (supervised, iOS 9.0+).",
					},
					"configuration_profile_block_changes": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block the user from installing configuration profiles and certificates interactively (supervised).",
					},
					// iCloud
					"icloud_block_backup": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block iCloud backup. Requires a supervised device for iOS 13+.",
					},
					"icloud_block_document_sync": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block iCloud document sync. Requires a supervised device for iOS 13+.",
					},
					"icloud_block_photo_stream_sync": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block iCloud Photo Stream sync.",
					},
					"icloud_block_managed_apps_sync": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block managed apps from using cloud sync.",
					},
					// Safari
					"safari_blocked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block the user from using Safari. Requires a supervised device for iOS 13+.",
					},
					"safari_block_javascript": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block JavaScript in Safari.",
					},
					"safari_block_popups": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block popups in Safari.",
					},
				},
			},
			"trusted_certificate": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "iOS/iPadOS trusted root certificate configuration (iosTrustedRootCertificate).",
				Attributes: map[string]schema.Attribute{
					"cert_file_name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The file name of the certificate file (.cer file) as displayed in the Intune admin center.",
					},
					"trusted_root_certificate": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The base64-encoded trusted root certificate content. Typically supplied via `filebase64(\"my-root-cert.cer\")`.",
					},
				},
			},
			"assignments": commonschemagraphbeta.DeviceConfigurationWithAllGroupAssignmentsAndFilterSchema(),
			"timeouts":    commonschema.ResourceTimeouts(ctx),
		},
	}
}

// appListItemAttributes returns the schema attributes for an appListItem entry.
func appListItemAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The display name of the app.",
		},
		"app_id": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The app identifier (bundle identifier for iOS).",
		},
		"app_store_url": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The App Store URL of the app.",
		},
		"publisher": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The publisher of the app.",
		},
	}
}
