# Terraform resource configuration for Microsoft 365 Graph Beta Device Management iOS/iPadOS Configuration Templates

# Example 1: iOS General Device Restriction (with Activation Lock)
resource "microsoft365_graph_beta_device_management_ios_device_configuration_templates" "general_restriction_example" {
  display_name = "iOS device restrictions with Activation Lock"
  description  = "Enables Activation Lock on supervised devices and enforces a passcode policy."

  general_device_restriction = {
    activation_lock_allow_when_supervised = true

    passcode_required                                     = true
    passcode_required_type                                = "alphanumeric"
    passcode_minimum_length                               = 6
    passcode_minutes_of_inactivity_before_lock            = 2
    passcode_minutes_of_inactivity_before_screen_timeout  = 5
    passcode_expiration_days                              = 90
    passcode_previous_passcode_block_count                = 5
    passcode_sign_in_failure_count_before_wipe            = 10
    passcode_block_simple                                 = true
    passcode_block_fingerprint_unlock                     = false
    passcode_minimum_character_set_count                  = 1

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
