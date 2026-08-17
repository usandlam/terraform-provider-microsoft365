package graphBetaIosDeviceConfigurationTemplates

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// IosDeviceConfigurationTemplatesResourceModel describes the resource data model.
type IosDeviceConfigurationTemplatesResourceModel struct {
	ID              types.String `tfsdk:"id"`
	DisplayName     types.String `tfsdk:"display_name"`
	Description     types.String `tfsdk:"description"`
	RoleScopeTagIds types.Set    `tfsdk:"role_scope_tag_ids"`
	// Nested configuration blocks (mutually exclusive)
	GeneralDeviceRestriction types.Object   `tfsdk:"general_device_restriction"`
	TrustedCertificate       types.Object   `tfsdk:"trusted_certificate"`
	Assignments              types.Set      `tfsdk:"assignments"`
	Timeouts                 timeouts.Value `tfsdk:"timeouts"`
}

// GeneralDeviceRestrictionResourceModel describes a curated subset of
// iosGeneralDeviceConfiguration properties (device restrictions template).
// Only high-value fields are exposed in this first pass; remaining properties
// can be added incrementally without breaking the schema.
type GeneralDeviceRestrictionResourceModel struct {
	// Activation Lock
	ActivationLockAllowWhenSupervised types.Bool `tfsdk:"activation_lock_allow_when_supervised"`

	// Passcode
	PasscodeRequired                               types.Bool   `tfsdk:"passcode_required"`
	PasscodeRequiredType                           types.String `tfsdk:"passcode_required_type"`
	PasscodeMinimumLength                          types.Int32  `tfsdk:"passcode_minimum_length"`
	PasscodeMinutesOfInactivityBeforeLock          types.Int32  `tfsdk:"passcode_minutes_of_inactivity_before_lock"`
	PasscodeMinutesOfInactivityBeforeScreenTimeout types.Int32  `tfsdk:"passcode_minutes_of_inactivity_before_screen_timeout"`
	PasscodeExpirationDays                         types.Int32  `tfsdk:"passcode_expiration_days"`
	PasscodePreviousPasscodeBlockCount             types.Int32  `tfsdk:"passcode_previous_passcode_block_count"`
	PasscodeSignInFailureCountBeforeWipe           types.Int32  `tfsdk:"passcode_sign_in_failure_count_before_wipe"`
	PasscodeBlockSimple                            types.Bool   `tfsdk:"passcode_block_simple"`
	PasscodeBlockFingerprintUnlock                 types.Bool   `tfsdk:"passcode_block_fingerprint_unlock"`
	PasscodeMinimumCharacterSetCount               types.Int32  `tfsdk:"passcode_minimum_character_set_count"`

	// App Store
	AppStoreBlocked                 types.Bool `tfsdk:"app_store_blocked"`
	AppStoreBlockInAppPurchases     types.Bool `tfsdk:"app_store_block_in_app_purchases"`
	AppStoreRequirePassword         types.Bool `tfsdk:"app_store_require_password"`
	AppStoreBlockAutomaticDownloads types.Bool `tfsdk:"app_store_block_automatic_downloads"`

	// App lists
	CompliantAppListType types.String `tfsdk:"compliant_app_list_type"`
	CompliantAppsList    types.Set    `tfsdk:"compliant_apps_list"`

	// Common device restrictions
	CameraBlocked                      types.Bool `tfsdk:"camera_blocked"`
	ScreenCaptureBlocked               types.Bool `tfsdk:"screen_capture_blocked"`
	SiriBlocked                        types.Bool `tfsdk:"siri_blocked"`
	SiriBlockedWhenLocked              types.Bool `tfsdk:"siri_blocked_when_locked"`
	AirDropBlocked                     types.Bool `tfsdk:"airdrop_blocked"`
	BluetoothBlockModification         types.Bool `tfsdk:"bluetooth_block_modification"`
	CellularBlockDataRoaming           types.Bool `tfsdk:"cellular_block_data_roaming"`
	CellularBlockVoiceRoaming          types.Bool `tfsdk:"cellular_block_voice_roaming"`
	CellularBlockPersonalHotspot       types.Bool `tfsdk:"cellular_block_personal_hotspot"`
	DeviceBlockEraseContentAndSettings types.Bool `tfsdk:"device_block_erase_content_and_settings"`
	DeviceBlockNameModification        types.Bool `tfsdk:"device_block_name_modification"`
	ConfigurationProfileBlockChanges   types.Bool `tfsdk:"configuration_profile_block_changes"`

	// iCloud / backup
	ICloudBlockBackup          types.Bool `tfsdk:"icloud_block_backup"`
	ICloudBlockDocumentSync    types.Bool `tfsdk:"icloud_block_document_sync"`
	ICloudBlockPhotoStreamSync types.Bool `tfsdk:"icloud_block_photo_stream_sync"`
	ICloudBlockManagedAppsSync types.Bool `tfsdk:"icloud_block_managed_apps_sync"`

	// Safari
	SafariBlocked         types.Bool `tfsdk:"safari_blocked"`
	SafariBlockJavaScript types.Bool `tfsdk:"safari_block_javascript"`
	SafariBlockPopups     types.Bool `tfsdk:"safari_block_popups"`
}

// TrustedCertificateResourceModel describes iosTrustedRootCertificate.
type TrustedCertificateResourceModel struct {
	CertFileName           types.String `tfsdk:"cert_file_name"`
	TrustedRootCertificate types.String `tfsdk:"trusted_root_certificate"`
}

// AppListItemResourceModel describes an appListItem entry, reused across
// compliant_apps_list / apps_visibility_list / apps_single_app_mode_list.
type AppListItemResourceModel struct {
	Name        types.String `tfsdk:"name"`
	AppId       types.String `tfsdk:"app_id"`
	AppStoreUrl types.String `tfsdk:"app_store_url"`
	Publisher   types.String `tfsdk:"publisher"`
}
