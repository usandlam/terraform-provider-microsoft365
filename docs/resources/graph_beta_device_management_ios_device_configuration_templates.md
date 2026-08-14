---
page_title: "microsoft365_graph_beta_device_management_ios_device_configuration_templates Resource - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Manages iOS/iPadOS configuration templates in Microsoft Intune. This resource creates device configurations for iOS/iPadOS devices including general device restrictions (with Activation Lock) and trusted root certificates.
---

# microsoft365_graph_beta_device_management_ios_device_configuration_templates (Resource)

Manages iOS/iPadOS configuration templates in Microsoft Intune. Exposes two mutually exclusive template types under a single resource:

- `general_device_restriction` — an `iosGeneralDeviceConfiguration` device restrictions profile. Includes `activation_lock_allow_when_supervised`, which is how Activation Lock is enabled on supervised iOS/iPadOS devices in Intune. Use the `microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock` action to bypass an active Activation Lock at runtime.
- `trusted_certificate` — an `iosTrustedRootCertificate` profile that deploys a base64-encoded root CA to devices.

Inline `assignments` are supported using the same shape as the macOS templates resource (all-devices, all-licensed-users, group, exclusion-group), each with optional assignment filter.

## Microsoft Documentation

- [iosGeneralDeviceConfiguration resource type](https://learn.microsoft.com/graph/api/resources/intune-deviceconfig-iosgeneraldeviceconfiguration?view=graph-rest-beta)
- [iosTrustedRootCertificate resource type](https://learn.microsoft.com/graph/api/resources/intune-deviceconfig-iostrustedrootcertificate?view=graph-rest-beta)
- [Apple device restriction settings in Microsoft Intune](https://learn.microsoft.com/intune/device-configuration/templates/ref-device-restrictions-apple)
- [Disable Activation Lock (iOS)](https://learn.microsoft.com/intune/device-management/actions/disable-activation-lock)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `DeviceManagementConfiguration.Read.All`
- `DeviceManagementConfiguration.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.28.0-alpha | Experimental | Initial release |

## Example Usage

```terraform
# Terraform resource configuration for Microsoft 365 Graph Beta Device Management iOS/iPadOS Configuration Templates

# Example 1: iOS General Device Restriction (with Activation Lock)
resource "microsoft365_graph_beta_device_management_ios_device_configuration_templates" "general_restriction_example" {
  display_name = "iOS device restrictions with Activation Lock"
  description  = "Enables Activation Lock on supervised devices and enforces a passcode policy."

  general_device_restriction = {
    activation_lock_allow_when_supervised = true

    passcode_required          = true
    passcode_required_type     = "alphanumeric"
    passcode_minimum_length    = 6
    passcode_block_simple      = true
    cellular_block_data_roaming = true
    safari_block_popups         = true
  }

  assignments = [
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000001"
      filter_id   = "00000000-0000-0000-0000-000000000002"
      filter_type = "include"
    }
  ]
}

# Example 2: iOS Trusted Root Certificate
resource "microsoft365_graph_beta_device_management_ios_device_configuration_templates" "trusted_certificate_example" {
  display_name = "iOS trusted root CA"

  trusted_certificate = {
    cert_file_name           = "corp-root-ca.cer"
    trusted_root_certificate = filebase64("${path.module}/corp-root-ca.cer")
  }

  assignments = [
    {
      type = "allDevicesAssignmentTarget"
    }
  ]
}
```

## Schema

### Required

- `display_name` (String) The display name for the iOS configuration template.

### Optional

- `assignments` (Attributes Set) Inline assignments. Each element takes `type` (one of `allDevicesAssignmentTarget`, `allLicensedUsersAssignmentTarget`, `groupAssignmentTarget`, `exclusionGroupAssignmentTarget`), optional `group_id`, and optional `filter_id`/`filter_type`.
- `description` (String)
- `general_device_restriction` (Attributes) iOS general device restriction template. See fields below.
- `role_scope_tag_ids` (Set of String) Defaults to `["0"]`.
- `timeouts` (Attributes)
- `trusted_certificate` (Attributes) iOS trusted root certificate. Fields: `cert_file_name`, `trusted_root_certificate` (base64-encoded).

Exactly one of `general_device_restriction` or `trusted_certificate` must be set.

### `general_device_restriction` fields

All fields optional. Grouped for readability:

- **Activation Lock:** `activation_lock_allow_when_supervised`
- **Passcode:** `passcode_required`, `passcode_required_type` (`deviceDefault` / `alphanumeric` / `numeric`), `passcode_minimum_length` (4–14), `passcode_minutes_of_inactivity_before_lock`, `passcode_minutes_of_inactivity_before_screen_timeout`, `passcode_expiration_days` (1–65535), `passcode_previous_passcode_block_count` (1–24), `passcode_sign_in_failure_count_before_wipe` (2–11), `passcode_block_simple`, `passcode_block_fingerprint_unlock`, `passcode_minimum_character_set_count` (0–4)
- **App Store:** `app_store_blocked`, `app_store_block_in_app_purchases`, `app_store_require_password`, `app_store_block_automatic_downloads`
- **Compliant apps:** `compliant_app_list_type` (`none` / `appsInListCompliant` / `appsNotInListCompliant`), `compliant_apps_list` (list of `{ name, app_id, app_store_url, publisher }`)
- **Device restrictions:** `camera_blocked`, `screen_capture_blocked`, `siri_blocked`, `siri_blocked_when_locked`, `airdrop_blocked`, `bluetooth_block_modification`, `cellular_block_data_roaming`, `cellular_block_voice_roaming`, `cellular_block_personal_hotspot`, `device_block_erase_content_and_settings`, `device_block_name_modification`, `configuration_profile_block_changes`
- **iCloud / backup:** `icloud_block_backup`, `icloud_block_document_sync`, `icloud_block_photo_stream_sync`, `icloud_block_managed_apps_sync`
- **Safari:** `safari_blocked`, `safari_block_javascript`, `safari_block_popups`

The Graph `iosGeneralDeviceConfiguration` resource exposes many additional properties (media content ratings, wallpaper, network usage rules, single sign-on payloads, kiosk-mode detail, etc.). Those can be added incrementally to this resource without breaking the schema.

### Read-Only

- `id` (String)

## Import

```shell
terraform import microsoft365_graph_beta_device_management_ios_device_configuration_templates.example <device-configuration-id>
```
