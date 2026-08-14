package graphBetaIosDeviceConfigurationTemplates

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Object type helpers for null values / set construction

func AppListItemType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":          types.StringType,
			"app_id":        types.StringType,
			"app_store_url": types.StringType,
			"publisher":     types.StringType,
		},
	}
}

func GeneralDeviceRestrictionType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"activation_lock_allow_when_supervised": types.BoolType,

			"passcode_required":                                   types.BoolType,
			"passcode_required_type":                              types.StringType,
			"passcode_minimum_length":                             types.Int32Type,
			"passcode_minutes_of_inactivity_before_lock":          types.Int32Type,
			"passcode_minutes_of_inactivity_before_screen_timeout": types.Int32Type,
			"passcode_expiration_days":                            types.Int32Type,
			"passcode_previous_passcode_block_count":              types.Int32Type,
			"passcode_sign_in_failure_count_before_wipe":          types.Int32Type,
			"passcode_block_simple":                               types.BoolType,
			"passcode_block_fingerprint_unlock":                   types.BoolType,
			"passcode_minimum_character_set_count":                types.Int32Type,

			"app_store_blocked":                   types.BoolType,
			"app_store_block_in_app_purchases":    types.BoolType,
			"app_store_require_password":          types.BoolType,
			"app_store_block_automatic_downloads": types.BoolType,

			"compliant_app_list_type": types.StringType,
			"compliant_apps_list":     types.SetType{ElemType: AppListItemType()},

			"camera_blocked":                          types.BoolType,
			"screen_capture_blocked":                  types.BoolType,
			"siri_blocked":                            types.BoolType,
			"siri_blocked_when_locked":                types.BoolType,
			"airdrop_blocked":                         types.BoolType,
			"bluetooth_block_modification":            types.BoolType,
			"cellular_block_data_roaming":             types.BoolType,
			"cellular_block_voice_roaming":            types.BoolType,
			"cellular_block_personal_hotspot":         types.BoolType,
			"device_block_erase_content_and_settings": types.BoolType,
			"device_block_name_modification":          types.BoolType,
			"configuration_profile_block_changes":     types.BoolType,

			"icloud_block_backup":            types.BoolType,
			"icloud_block_document_sync":     types.BoolType,
			"icloud_block_photo_stream_sync": types.BoolType,
			"icloud_block_managed_apps_sync": types.BoolType,

			"safari_blocked":          types.BoolType,
			"safari_block_javascript": types.BoolType,
			"safari_block_popups":     types.BoolType,
		},
	}
}

func TrustedCertificateType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"cert_file_name":           types.StringType,
			"trusted_root_certificate": types.StringType,
		},
	}
}

// MapRemoteResourceStateToTerraform maps a Graph deviceConfiguration response
// back to the Terraform model, dispatching on the concrete Go type.
func MapRemoteResourceStateToTerraform(ctx context.Context, data *IosDeviceConfigurationTemplatesResourceModel, remoteResource graphmodels.DeviceConfigurationable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())

	switch config := remoteResource.(type) {
	case *graphmodels.IosGeneralDeviceConfiguration:
		mapIosGeneralDeviceConfiguration(ctx, data, config)
	case *graphmodels.IosTrustedRootCertificate:
		mapIosTrustedRootCertificate(ctx, data, config)
	default:
		tflog.Error(ctx, "Unknown device configuration type", map[string]any{
			"type": fmt.Sprintf("%T", config),
		})
	}

	assignments := remoteResource.GetAssignments()
	tflog.Debug(ctx, "Retrieved assignments from remote resource", map[string]any{
		"assignmentCount": len(assignments),
		"resourceId":      data.ID.ValueString(),
	})

	if len(assignments) == 0 {
		data.Assignments = types.SetNull(IosConfigurationTemplatesAssignmentType())
	} else {
		mapAssignmentsToTerraform(ctx, data, assignments)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}

func mapIosGeneralDeviceConfiguration(ctx context.Context, data *IosDeviceConfigurationTemplatesResourceModel, config *graphmodels.IosGeneralDeviceConfiguration) {
	tflog.Debug(ctx, "Mapping IosGeneralDeviceConfiguration")

	compliantAppsSet := mapAppListItemsToSet(ctx, config.GetCompliantAppsList())

	model := GeneralDeviceRestrictionResourceModel{
		ActivationLockAllowWhenSupervised: convert.GraphToFrameworkBool(config.GetActivationLockAllowWhenSupervised()),

		PasscodeRequired:                              convert.GraphToFrameworkBool(config.GetPasscodeRequired()),
		PasscodeRequiredType:                          convert.GraphToFrameworkEnum(config.GetPasscodeRequiredType()),
		PasscodeMinimumLength:                         convert.GraphToFrameworkInt32(config.GetPasscodeMinimumLength()),
		PasscodeMinutesOfInactivityBeforeLock:         convert.GraphToFrameworkInt32(config.GetPasscodeMinutesOfInactivityBeforeLock()),
		PasscodeMinutesOfInactivityBeforeScreenTimeout: convert.GraphToFrameworkInt32(config.GetPasscodeMinutesOfInactivityBeforeScreenTimeout()),
		PasscodeExpirationDays:                        convert.GraphToFrameworkInt32(config.GetPasscodeExpirationDays()),
		PasscodePreviousPasscodeBlockCount:            convert.GraphToFrameworkInt32(config.GetPasscodePreviousPasscodeBlockCount()),
		PasscodeSignInFailureCountBeforeWipe:          convert.GraphToFrameworkInt32(config.GetPasscodeSignInFailureCountBeforeWipe()),
		PasscodeBlockSimple:                           convert.GraphToFrameworkBool(config.GetPasscodeBlockSimple()),
		PasscodeBlockFingerprintUnlock:                convert.GraphToFrameworkBool(config.GetPasscodeBlockFingerprintUnlock()),
		PasscodeMinimumCharacterSetCount:              convert.GraphToFrameworkInt32(config.GetPasscodeMinimumCharacterSetCount()),

		AppStoreBlocked:                 convert.GraphToFrameworkBool(config.GetAppStoreBlocked()),
		AppStoreBlockInAppPurchases:     convert.GraphToFrameworkBool(config.GetAppStoreBlockInAppPurchases()),
		AppStoreRequirePassword:         convert.GraphToFrameworkBool(config.GetAppStoreRequirePassword()),
		AppStoreBlockAutomaticDownloads: convert.GraphToFrameworkBool(config.GetAppStoreBlockAutomaticDownloads()),

		CompliantAppListType: convert.GraphToFrameworkEnum(config.GetCompliantAppListType()),
		CompliantAppsList:    compliantAppsSet,

		CameraBlocked:                      convert.GraphToFrameworkBool(config.GetCameraBlocked()),
		ScreenCaptureBlocked:               convert.GraphToFrameworkBool(config.GetScreenCaptureBlocked()),
		SiriBlocked:                        convert.GraphToFrameworkBool(config.GetSiriBlocked()),
		SiriBlockedWhenLocked:              convert.GraphToFrameworkBool(config.GetSiriBlockedWhenLocked()),
		AirDropBlocked:                     convert.GraphToFrameworkBool(config.GetAirDropBlocked()),
		BluetoothBlockModification:         convert.GraphToFrameworkBool(config.GetBluetoothBlockModification()),
		CellularBlockDataRoaming:           convert.GraphToFrameworkBool(config.GetCellularBlockDataRoaming()),
		CellularBlockVoiceRoaming:          convert.GraphToFrameworkBool(config.GetCellularBlockVoiceRoaming()),
		CellularBlockPersonalHotspot:       convert.GraphToFrameworkBool(config.GetCellularBlockPersonalHotspot()),
		DeviceBlockEraseContentAndSettings: convert.GraphToFrameworkBool(config.GetDeviceBlockEraseContentAndSettings()),
		DeviceBlockNameModification:        convert.GraphToFrameworkBool(config.GetDeviceBlockNameModification()),
		ConfigurationProfileBlockChanges:   convert.GraphToFrameworkBool(config.GetConfigurationProfileBlockChanges()),

		ICloudBlockBackup:          convert.GraphToFrameworkBool(config.GetICloudBlockBackup()),
		ICloudBlockDocumentSync:    convert.GraphToFrameworkBool(config.GetICloudBlockDocumentSync()),
		ICloudBlockPhotoStreamSync: convert.GraphToFrameworkBool(config.GetICloudBlockPhotoStreamSync()),
		ICloudBlockManagedAppsSync: convert.GraphToFrameworkBool(config.GetICloudBlockManagedAppsSync()),

		SafariBlocked:         convert.GraphToFrameworkBool(config.GetSafariBlocked()),
		SafariBlockJavaScript: convert.GraphToFrameworkBool(config.GetSafariBlockJavaScript()),
		SafariBlockPopups:     convert.GraphToFrameworkBool(config.GetSafariBlockPopups()),
	}

	objectValue, diags := types.ObjectValueFrom(ctx, GeneralDeviceRestrictionType().AttrTypes, model)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to create general_device_restriction object", map[string]any{"errors": diags.Errors()})
		return
	}
	data.GeneralDeviceRestriction = objectValue
	data.TrustedCertificate = types.ObjectNull(TrustedCertificateType().AttrTypes)
}

func mapIosTrustedRootCertificate(ctx context.Context, data *IosDeviceConfigurationTemplatesResourceModel, config *graphmodels.IosTrustedRootCertificate) {
	tflog.Debug(ctx, "Mapping IosTrustedRootCertificate")

	var trustedRootCert types.String
	if certBytes := config.GetTrustedRootCertificate(); certBytes != nil {
		trustedRootCert = types.StringValue(base64.StdEncoding.EncodeToString(certBytes))
	} else {
		trustedRootCert = types.StringNull()
	}

	certModel := TrustedCertificateResourceModel{
		CertFileName:           convert.GraphToFrameworkString(config.GetCertFileName()),
		TrustedRootCertificate: trustedRootCert,
	}

	objectValue, diags := types.ObjectValueFrom(ctx, TrustedCertificateType().AttrTypes, certModel)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to create trusted_certificate object", map[string]any{"errors": diags.Errors()})
		return
	}
	data.TrustedCertificate = objectValue
	data.GeneralDeviceRestriction = types.ObjectNull(GeneralDeviceRestrictionType().AttrTypes)
}

// mapAppListItemsToSet converts a slice of AppListItemable to a TF Set of AppListItem objects.
func mapAppListItemsToSet(ctx context.Context, apps []graphmodels.AppListItemable) types.Set {
	if len(apps) == 0 {
		return types.SetNull(AppListItemType())
	}

	appModels := make([]AppListItemResourceModel, 0, len(apps))
	for _, app := range apps {
		if app == nil {
			continue
		}
		appModels = append(appModels, AppListItemResourceModel{
			Name:        convert.GraphToFrameworkString(app.GetName()),
			AppId:       convert.GraphToFrameworkString(app.GetAppId()),
			AppStoreUrl: convert.GraphToFrameworkString(app.GetAppStoreUrl()),
			Publisher:   convert.GraphToFrameworkString(app.GetPublisher()),
		})
	}

	setValue, diags := types.SetValueFrom(ctx, AppListItemType(), appModels)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to create app list items set", map[string]any{"errors": diags.Errors()})
		return types.SetNull(AppListItemType())
	}
	return setValue
}
