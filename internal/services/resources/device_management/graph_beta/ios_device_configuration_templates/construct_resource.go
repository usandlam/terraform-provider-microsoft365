package graphBetaIosDeviceConfigurationTemplates

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource dispatches to the concrete iOS device configuration
// constructor based on which nested block is set on the Terraform model.
func constructResource(ctx context.Context, data *IosDeviceConfigurationTemplatesResourceModel) (graphmodels.DeviceConfigurationable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	var requestBody graphmodels.DeviceConfigurationable

	switch {
	case !data.GeneralDeviceRestriction.IsNull() && !data.GeneralDeviceRestriction.IsUnknown():
		requestBody = constructIosGeneralDeviceConfiguration(ctx, data)
	case !data.TrustedCertificate.IsNull() && !data.TrustedCertificate.IsUnknown():
		requestBody = constructIosTrustedRootCertificate(ctx, data)
	default:
		return nil, fmt.Errorf("no configuration type specified")
	}

	if requestBody == nil {
		return nil, fmt.Errorf("failed to construct configuration")
	}

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))
	return requestBody, nil
}

// constructIosGeneralDeviceConfiguration builds an IosGeneralDeviceConfiguration.
func constructIosGeneralDeviceConfiguration(ctx context.Context, data *IosDeviceConfigurationTemplatesResourceModel) graphmodels.DeviceConfigurationable {
	tflog.Debug(ctx, "Constructing IosGeneralDeviceConfiguration")

	cfg := graphmodels.NewIosGeneralDeviceConfiguration()

	var g GeneralDeviceRestrictionResourceModel
	diags := data.GeneralDeviceRestriction.As(ctx, &g, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		tflog.Error(ctx, "Failed to extract general_device_restriction data")
		return nil
	}

	// Activation Lock
	convert.FrameworkToGraphBool(g.ActivationLockAllowWhenSupervised, cfg.SetActivationLockAllowWhenSupervised)

	// Passcode
	convert.FrameworkToGraphBool(g.PasscodeRequired, cfg.SetPasscodeRequired)
	if err := convert.FrameworkToGraphEnum(g.PasscodeRequiredType, graphmodels.ParseRequiredPasswordType, cfg.SetPasscodeRequiredType); err != nil {
		tflog.Error(ctx, "Failed to set passcode_required_type", map[string]any{"error": err.Error()})
		return nil
	}
	convert.FrameworkToGraphInt32(g.PasscodeMinimumLength, cfg.SetPasscodeMinimumLength)
	convert.FrameworkToGraphInt32(g.PasscodeMinutesOfInactivityBeforeLock, cfg.SetPasscodeMinutesOfInactivityBeforeLock)
	convert.FrameworkToGraphInt32(g.PasscodeMinutesOfInactivityBeforeScreenTimeout, cfg.SetPasscodeMinutesOfInactivityBeforeScreenTimeout)
	convert.FrameworkToGraphInt32(g.PasscodeExpirationDays, cfg.SetPasscodeExpirationDays)
	convert.FrameworkToGraphInt32(g.PasscodePreviousPasscodeBlockCount, cfg.SetPasscodePreviousPasscodeBlockCount)
	convert.FrameworkToGraphInt32(g.PasscodeSignInFailureCountBeforeWipe, cfg.SetPasscodeSignInFailureCountBeforeWipe)
	convert.FrameworkToGraphBool(g.PasscodeBlockSimple, cfg.SetPasscodeBlockSimple)
	convert.FrameworkToGraphBool(g.PasscodeBlockFingerprintUnlock, cfg.SetPasscodeBlockFingerprintUnlock)
	convert.FrameworkToGraphInt32(g.PasscodeMinimumCharacterSetCount, cfg.SetPasscodeMinimumCharacterSetCount)

	// App Store
	convert.FrameworkToGraphBool(g.AppStoreBlocked, cfg.SetAppStoreBlocked)
	convert.FrameworkToGraphBool(g.AppStoreBlockInAppPurchases, cfg.SetAppStoreBlockInAppPurchases)
	convert.FrameworkToGraphBool(g.AppStoreRequirePassword, cfg.SetAppStoreRequirePassword)
	convert.FrameworkToGraphBool(g.AppStoreBlockAutomaticDownloads, cfg.SetAppStoreBlockAutomaticDownloads)

	// Compliant apps list
	if err := convert.FrameworkToGraphEnum(g.CompliantAppListType, graphmodels.ParseAppListType, cfg.SetCompliantAppListType); err != nil {
		tflog.Error(ctx, "Failed to set compliant_app_list_type", map[string]any{"error": err.Error()})
		return nil
	}
	if err := setAppListItems(ctx, g.CompliantAppsList, cfg.SetCompliantAppsList); err != nil {
		tflog.Error(ctx, "Failed to set compliant_apps_list", map[string]any{"error": err.Error()})
		return nil
	}

	// Common device restrictions
	convert.FrameworkToGraphBool(g.CameraBlocked, cfg.SetCameraBlocked)
	convert.FrameworkToGraphBool(g.ScreenCaptureBlocked, cfg.SetScreenCaptureBlocked)
	convert.FrameworkToGraphBool(g.SiriBlocked, cfg.SetSiriBlocked)
	convert.FrameworkToGraphBool(g.SiriBlockedWhenLocked, cfg.SetSiriBlockedWhenLocked)
	convert.FrameworkToGraphBool(g.AirDropBlocked, cfg.SetAirDropBlocked)
	convert.FrameworkToGraphBool(g.BluetoothBlockModification, cfg.SetBluetoothBlockModification)
	convert.FrameworkToGraphBool(g.CellularBlockDataRoaming, cfg.SetCellularBlockDataRoaming)
	convert.FrameworkToGraphBool(g.CellularBlockVoiceRoaming, cfg.SetCellularBlockVoiceRoaming)
	convert.FrameworkToGraphBool(g.CellularBlockPersonalHotspot, cfg.SetCellularBlockPersonalHotspot)
	convert.FrameworkToGraphBool(g.DeviceBlockEraseContentAndSettings, cfg.SetDeviceBlockEraseContentAndSettings)
	convert.FrameworkToGraphBool(g.DeviceBlockNameModification, cfg.SetDeviceBlockNameModification)
	convert.FrameworkToGraphBool(g.ConfigurationProfileBlockChanges, cfg.SetConfigurationProfileBlockChanges)

	// iCloud
	convert.FrameworkToGraphBool(g.ICloudBlockBackup, cfg.SetICloudBlockBackup)
	convert.FrameworkToGraphBool(g.ICloudBlockDocumentSync, cfg.SetICloudBlockDocumentSync)
	convert.FrameworkToGraphBool(g.ICloudBlockPhotoStreamSync, cfg.SetICloudBlockPhotoStreamSync)
	convert.FrameworkToGraphBool(g.ICloudBlockManagedAppsSync, cfg.SetICloudBlockManagedAppsSync)

	// Safari
	convert.FrameworkToGraphBool(g.SafariBlocked, cfg.SetSafariBlocked)
	convert.FrameworkToGraphBool(g.SafariBlockJavaScript, cfg.SetSafariBlockJavaScript)
	convert.FrameworkToGraphBool(g.SafariBlockPopups, cfg.SetSafariBlockPopups)

	return cfg
}

// constructIosTrustedRootCertificate builds an IosTrustedRootCertificate.
func constructIosTrustedRootCertificate(ctx context.Context, data *IosDeviceConfigurationTemplatesResourceModel) graphmodels.DeviceConfigurationable {
	tflog.Debug(ctx, "Constructing IosTrustedRootCertificate")

	certConfig := graphmodels.NewIosTrustedRootCertificate()

	var certData TrustedCertificateResourceModel
	diags := data.TrustedCertificate.As(ctx, &certData, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		tflog.Error(ctx, "Failed to extract trusted_certificate data")
		return nil
	}

	convert.FrameworkToGraphString(certData.CertFileName, certConfig.SetCertFileName)

	if !certData.TrustedRootCertificate.IsNull() && !certData.TrustedRootCertificate.IsUnknown() {
		certBase64 := certData.TrustedRootCertificate.ValueString()
		certBytes, err := base64.StdEncoding.DecodeString(certBase64)
		if err != nil {
			tflog.Error(ctx, "Failed to decode base64 certificate data", map[string]any{"error": err.Error()})
			return nil
		}
		certConfig.SetTrustedRootCertificate(certBytes)
	}

	return certConfig
}

// setAppListItems reads a Terraform Set of AppListItemResourceModel entries
// and passes them to the given Graph SDK setter.
func setAppListItems(ctx context.Context, appsSet types.Set, setter func([]graphmodels.AppListItemable)) error {
	if appsSet.IsNull() || appsSet.IsUnknown() {
		return nil
	}

	var appModels []AppListItemResourceModel
	diags := appsSet.ElementsAs(ctx, &appModels, false)
	if diags.HasError() {
		return fmt.Errorf("failed to extract app list items: %v", diags.Errors())
	}

	apps := make([]graphmodels.AppListItemable, 0, len(appModels))
	for _, appModel := range appModels {
		app := graphmodels.NewAppListItem()
		convert.FrameworkToGraphString(appModel.Name, app.SetName)
		convert.FrameworkToGraphString(appModel.AppId, app.SetAppId)
		convert.FrameworkToGraphString(appModel.AppStoreUrl, app.SetAppStoreUrl)
		convert.FrameworkToGraphString(appModel.Publisher, app.SetPublisher)
		apps = append(apps, app)
	}
	setter(apps)
	return nil
}
