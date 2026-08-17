resource "random_string" "general_restriction_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_groups_group" "general_restriction_group_1" {
  display_name     = "acc-test-ios-general-restriction-group-1-${random_string.general_restriction_suffix.result}"
  mail_nickname    = "acc-test-ios-restriction-1-${random_string.general_restriction_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "general_restriction_group_2" {
  display_name     = "acc-test-ios-general-restriction-group-2-${random_string.general_restriction_suffix.result}"
  mail_nickname    = "acc-test-ios-restriction-2-${random_string.general_restriction_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_device_management_ios_device_configuration_templates" "general_restriction_example" {
  display_name = "acc-test-iOS-general-restriction-${random_string.general_restriction_suffix.result}"
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

    safari_block_popups = true
  }

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.general_restriction_group_1.id
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.general_restriction_group_2.id
    }
  ]

  depends_on = [
    microsoft365_graph_beta_groups_group.general_restriction_group_1,
    microsoft365_graph_beta_groups_group.general_restriction_group_2,
  ]

  timeouts = {
    create = "3m"
    read   = "3m"
    update = "3m"
    delete = "3m"
  }
}
