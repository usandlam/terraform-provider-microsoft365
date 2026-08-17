resource "random_string" "trusted_cert_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_groups_group" "trusted_cert_group_1" {
  display_name     = "acc-test-ios-trusted-cert-group-1-${random_string.trusted_cert_suffix.result}"
  mail_nickname    = "acc-test-ios-trusted-cert-1-${random_string.trusted_cert_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "trusted_cert_group_2" {
  display_name     = "acc-test-ios-trusted-cert-group-2-${random_string.trusted_cert_suffix.result}"
  mail_nickname    = "acc-test-ios-trusted-cert-2-${random_string.trusted_cert_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_device_management_ios_device_configuration_templates" "trusted_cert_example" {
  display_name = "acc-test-iOS-trusted-root-cert-${random_string.trusted_cert_suffix.result}"
  description  = "Install company root certificate for secure connections"

  trusted_certificate = {
    cert_file_name           = "MicrosoftRootCertificateAuthority2011.cer"
    trusted_root_certificate = filebase64("tests/terraform/acceptance/MicrosoftRootCertificateAuthority2011.cer")
  }

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.trusted_cert_group_1.id
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.trusted_cert_group_2.id
    }
  ]

  depends_on = [
    microsoft365_graph_beta_groups_group.trusted_cert_group_1,
    microsoft365_graph_beta_groups_group.trusted_cert_group_2,
  ]

  timeouts = {
    create = "50s"
    read   = "5m"
    update = "30m"
    delete = "30m"
  }
}
