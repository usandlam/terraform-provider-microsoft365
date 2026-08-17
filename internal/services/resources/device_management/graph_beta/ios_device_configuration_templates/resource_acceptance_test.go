package graphBetaIosDeviceConfigurationTemplates_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaIosDeviceConfigurationTemplates "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/ios_device_configuration_templates"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	resourceType = graphBetaIosDeviceConfigurationTemplates.ResourceName
	testResource = graphBetaIosDeviceConfigurationTemplates.IosDeviceConfigurationTemplatesTestResource{}
)

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

// General Device Restriction Tests
func TestAccResourceIosDeviceConfigurationTemplates_01_GeneralDeviceRestriction(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating general device restriction configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_general_device_restriction_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("iOS device configuration", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".general_restriction_example").ExistsInGraph(testResource),
					check.That(resourceType+".general_restriction_example").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".general_restriction_example").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-iOS-general-restriction-[a-z0-9]{8}$`)),
					check.That(resourceType+".general_restriction_example").Key("general_device_restriction.activation_lock_allow_when_supervised").HasValue("true"),
					check.That(resourceType+".general_restriction_example").Key("general_device_restriction.passcode_required").HasValue("true"),
					check.That(resourceType+".general_restriction_example").Key("assignments.#").HasValue("2"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing general device restriction configuration")
				},
				ResourceName:      resourceType + ".general_restriction_example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Trusted Root Certificate Tests
func TestAccResourceIosDeviceConfigurationTemplates_02_TrustedRootCertificate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating trusted root certificate configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_trusted_root_certificate_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("iOS device configuration", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".trusted_cert_example").ExistsInGraph(testResource),
					check.That(resourceType+".trusted_cert_example").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".trusted_cert_example").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-iOS-trusted-root-cert-[a-z0-9]{8}$`)),
					check.That(resourceType+".trusted_cert_example").Key("trusted_certificate.cert_file_name").HasValue("MicrosoftRootCertificateAuthority2011.cer"),
					check.That(resourceType+".trusted_cert_example").Key("trusted_certificate.trusted_root_certificate").Exists(),
					check.That(resourceType+".trusted_cert_example").Key("assignments.#").HasValue("2"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing trusted root certificate configuration")
				},
				ResourceName:      resourceType + ".trusted_cert_example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
