---
page_title: "microsoft365_graph_beta_device_management_ios_device_configuration_templates Resource - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Manages iOS/iPadOS configuration templates in Microsoft Intune. This resource creates device configurations for iOS/iPadOS devices including general device restrictions (with Activation Lock) and trusted root certificates.
---

# microsoft365_graph_beta_device_management_ios_device_configuration_templates (Resource)

Manages iOS/iPadOS configuration templates in Microsoft Intune. This resource creates device configurations for iOS/iPadOS devices including general device restrictions (with Activation Lock) and trusted root certificates.

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

**Optional:**
- `None` `[N/A]`

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

    passcode_required                                    = true
    passcode_required_type                               = "alphanumeric"
    passcode_minimum_length                              = 6
    passcode_minutes_of_inactivity_before_lock           = 2
    passcode_minutes_of_inactivity_before_screen_timeout = 5
    passcode_expiration_days                             = 90
    passcode_previous_passcode_block_count               = 5
    passcode_sign_in_failure_count_before_wipe           = 10
    passcode_block_simple                                = true
    passcode_block_fingerprint_unlock                    = false
    passcode_minimum_character_set_count                 = 1

    app_store_blocked                   = false
    app_store_block_in_app_purchases    = true
    app_store_require_password          = true
    app_store_block_automatic_downloads = false

    camera_blocked                          = false
    screen_capture_blocked                  = false
    siri_blocked                            = false
    siri_blocked_when_locked                = true
    airdrop_blocked                         = false
    bluetooth_block_modification            = false
    cellular_block_data_roaming             = true
    cellular_block_voice_roaming            = true
    cellular_block_personal_hotspot         = false
    device_block_erase_content_and_settings = true
    device_block_name_modification          = true
    configuration_profile_block_changes     = true

    icloud_block_backup            = false
    icloud_block_document_sync     = false
    icloud_block_photo_stream_sync = false
    icloud_block_managed_apps_sync = false

    safari_blocked          = false
    safari_block_javascript = false
    safari_block_popups     = true
  }

  role_scope_tag_ids = ["0"]

  assignments = [
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000001"
      filter_id   = "00000000-0000-0000-0000-000000000002"
      filter_type = "include"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000003"
    }
  ]

  timeouts = {
    create = "3m"
    read   = "3m"
    update = "3m"
    delete = "3m"
  }
}

# Example 2: iOS Trusted Root Certificate
resource "microsoft365_graph_beta_device_management_ios_device_configuration_templates" "trusted_certificate_example" {
  display_name = "iOS trusted root CA"
  description  = "Deploys the corporate root CA to iOS/iPadOS devices."

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

# Example 3: iOS General Device Restriction with a compliant apps allow-list
resource "microsoft365_graph_beta_device_management_ios_device_configuration_templates" "compliant_apps_example" {
  display_name = "iOS compliant apps allow-list"

  general_device_restriction = {
    compliant_app_list_type = "appsInListCompliant"
    compliant_apps_list = [
      {
        name      = "Microsoft Outlook"
        app_id    = "com.microsoft.Office.Outlook"
        publisher = "Microsoft Corporation"
      },
      {
        name      = "Microsoft Teams"
        app_id    = "com.microsoft.skype.teams"
        publisher = "Microsoft Corporation"
      }
    ]
  }

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000001"
    }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) The display name for the iOS configuration template.

### Optional

- `assignments` (Attributes Set) Assignments for the device configuration. Each assignment specifies the target group and schedule for script execution. Supports group filters. (see [below for nested schema](#nestedatt--assignments))
- `description` (String) Optional description of the resource. Maximum length is 1500 characters.
- `general_device_restriction` (Attributes) iOS/iPadOS general device restriction template (iosGeneralDeviceConfiguration). Covers passcode, App Store, common device restrictions, iCloud/backup, Safari, and the Activation Lock allow-when-supervised toggle. This is a curated subset of the full Graph resource surface; additional properties can be added over time without breaking the schema. (see [below for nested schema](#nestedatt--general_device_restriction))
- `role_scope_tag_ids` (Set of String) Set of scope tag IDs for this template profile.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `trusted_certificate` (Attributes) iOS/iPadOS trusted root certificate configuration (iosTrustedRootCertificate). (see [below for nested schema](#nestedatt--trusted_certificate))

### Read-Only

- `id` (String) The unique identifier for the iOS configuration template.

<a id="nestedatt--assignments"></a>
### Nested Schema for `assignments`

Required:

- `type` (String) Type of assignment target. Must be one of: 'allDevicesAssignmentTarget', 'allLicensedUsersAssignmentTarget', 'groupAssignmentTarget', 'exclusionGroupAssignmentTarget'.

Optional:

- `filter_id` (String) ID of the filter to apply to the assignment. Required when filter_type is 'include' or 'exclude'. Should be omitted when filter_type is 'none'.
- `filter_type` (String) Type of filter to apply. Must be one of: 'include', 'exclude', or 'none'.
- `group_id` (String) The Entra ID group ID to include or exclude in the assignment. Required when type is 'groupAssignmentTarget' or 'exclusionGroupAssignmentTarget'.


<a id="nestedatt--general_device_restriction"></a>
### Nested Schema for `general_device_restriction`

Optional:

- `activation_lock_allow_when_supervised` (Boolean) Indicates whether or not to allow activation lock when the device is in supervised mode. Enables Activation Lock on supervised iOS/iPadOS devices; use the `bypass_activation_lock` device action to remove it later.
- `airdrop_blocked` (Boolean) Indicates whether or not to allow AirDrop when the device is in supervised mode.
- `app_store_block_automatic_downloads` (Boolean) Indicates whether or not to block the automatic downloading of apps purchased on other devices (supervised, iOS 9.0+).
- `app_store_block_in_app_purchases` (Boolean) Indicates whether or not to block the user from making in-app purchases.
- `app_store_blocked` (Boolean) Indicates whether or not to block the user from using the App Store. Requires a supervised device for iOS 13+.
- `app_store_require_password` (Boolean) Indicates whether or not to require a password when using the App Store.
- `bluetooth_block_modification` (Boolean) Indicates whether or not to allow modification of Bluetooth settings when the device is in supervised mode (iOS 10.0+).
- `camera_blocked` (Boolean) Indicates whether or not to block the user from accessing the camera. Requires a supervised device for iOS 13+.
- `cellular_block_data_roaming` (Boolean) Indicates whether or not to block data roaming.
- `cellular_block_personal_hotspot` (Boolean) Indicates whether or not to block Personal Hotspot.
- `cellular_block_voice_roaming` (Boolean) Indicates whether or not to block voice roaming.
- `compliant_app_list_type` (String) List type that is in the compliant apps list. Possible values: none, appsInListCompliant, appsNotInListCompliant.
- `compliant_apps_list` (Attributes Set) List of apps in the compliance list (allow or block list, controlled by compliant_app_list_type). (see [below for nested schema](#nestedatt--general_device_restriction--compliant_apps_list))
- `configuration_profile_block_changes` (Boolean) Indicates whether or not to block the user from installing configuration profiles and certificates interactively (supervised).
- `device_block_erase_content_and_settings` (Boolean) Indicates whether or not to allow the use of the 'Erase all content and settings' option when supervised.
- `device_block_name_modification` (Boolean) Indicates whether or not to allow device name modification (supervised, iOS 9.0+).
- `icloud_block_backup` (Boolean) Indicates whether or not to block iCloud backup. Requires a supervised device for iOS 13+.
- `icloud_block_document_sync` (Boolean) Indicates whether or not to block iCloud document sync. Requires a supervised device for iOS 13+.
- `icloud_block_managed_apps_sync` (Boolean) Indicates whether or not to block managed apps from using cloud sync.
- `icloud_block_photo_stream_sync` (Boolean) Indicates whether or not to block iCloud Photo Stream sync.
- `passcode_block_fingerprint_unlock` (Boolean) Indicates whether or not to block fingerprint unlock.
- `passcode_block_simple` (Boolean) Indicates whether or not to block simple passcodes.
- `passcode_expiration_days` (Number) Number of days before the passcode expires. Valid values 1 to 65535.
- `passcode_minimum_character_set_count` (Number) Number of character sets a passcode must contain. Valid values 0 to 4.
- `passcode_minimum_length` (Number) Minimum length of passcode. Valid values 4 to 14.
- `passcode_minutes_of_inactivity_before_lock` (Number) Minutes of inactivity before a passcode is required.
- `passcode_minutes_of_inactivity_before_screen_timeout` (Number) Minutes of inactivity before the screen times out.
- `passcode_previous_passcode_block_count` (Number) Number of previous passcodes to block. Valid values 1 to 24.
- `passcode_required` (Boolean) Indicates whether or not to require a passcode.
- `passcode_required_type` (String) Type of passcode that is required. Possible values: deviceDefault, alphanumeric, numeric.
- `passcode_sign_in_failure_count_before_wipe` (Number) Number of sign in failures allowed before wiping the device. Valid values 2 to 11.
- `safari_block_javascript` (Boolean) Indicates whether or not to block JavaScript in Safari.
- `safari_block_popups` (Boolean) Indicates whether or not to block popups in Safari.
- `safari_blocked` (Boolean) Indicates whether or not to block the user from using Safari. Requires a supervised device for iOS 13+.
- `screen_capture_blocked` (Boolean) Indicates whether or not to block the user from taking screen captures.
- `siri_blocked` (Boolean) Indicates whether or not to block the user from using Siri.
- `siri_blocked_when_locked` (Boolean) Indicates whether or not to block the user from using Siri when locked.

<a id="nestedatt--general_device_restriction--compliant_apps_list"></a>
### Nested Schema for `general_device_restriction.compliant_apps_list`

Required:

- `name` (String) The display name of the app.

Optional:

- `app_id` (String) The app identifier (bundle identifier for iOS).
- `app_store_url` (String) The App Store URL of the app.
- `publisher` (String) The publisher of the app.



<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--trusted_certificate"></a>
### Nested Schema for `trusted_certificate`

Required:

- `cert_file_name` (String) The file name of the certificate file (.cer file) as displayed in the Intune admin center.
- `trusted_root_certificate` (String) The base64-encoded trusted root certificate content. Typically supplied via `filebase64("my-root-cert.cer")`.

## Important Notes

- **Mutually exclusive template types**: Exactly one of `general_device_restriction` or `trusted_certificate` must be specified per resource instance.
- **Activation Lock**: Enabling `activation_lock_allow_when_supervised` in `general_device_restriction` is how Activation Lock is turned on in Intune. Use the `microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock` device action to bypass an active Activation Lock at runtime.
- **Assignment Required**: Policies must be assigned to device or user groups to be deployed.

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash

# Example 1: Import an iOS general device restriction template
echo "terraform import microsoft365_graph_beta_device_management_ios_device_configuration_templates.general_restriction_example \"12345678-1234-1234-1234-123456789012\""

# Example 2: Import an iOS trusted root certificate template
echo "terraform import microsoft365_graph_beta_device_management_ios_device_configuration_templates.trusted_certificate_example \"87654321-4321-4321-4321-210987654321\""
```
