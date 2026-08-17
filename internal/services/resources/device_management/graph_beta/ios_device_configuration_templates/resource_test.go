package graphBetaIosDeviceConfigurationTemplates_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	iosConfigMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/ios_device_configuration_templates/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *iosConfigMocks.IosDeviceConfigurationTemplatesMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	configMock := &iosConfigMocks.IosDeviceConfigurationTemplatesMock{}
	configMock.RegisterMocks()
	return mockClient, configMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *iosConfigMocks.IosDeviceConfigurationTemplatesMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	configMock := &iosConfigMocks.IosDeviceConfigurationTemplatesMock{}
	configMock.RegisterErrorMocks()
	return mockClient, configMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func TestUnitResourceIosDeviceConfigurationTemplates_01_GeneralDeviceRestriction(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, configMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer configMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_general_device_restriction_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".general_restriction_example").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".general_restriction_example").Key("display_name").HasValue("unit-test-iOS-general-device-restriction-example"),
					check.That(resourceType+".general_restriction_example").Key("description").HasValue("Enables Activation Lock on supervised devices and enforces a passcode policy."),
					check.That(resourceType+".general_restriction_example").Key("general_device_restriction.activation_lock_allow_when_supervised").HasValue("true"),
					check.That(resourceType+".general_restriction_example").Key("general_device_restriction.passcode_required").HasValue("true"),
					check.That(resourceType+".general_restriction_example").Key("general_device_restriction.passcode_required_type").HasValue("alphanumeric"),
					check.That(resourceType+".general_restriction_example").Key("general_device_restriction.passcode_minimum_length").HasValue("6"),
					check.That(resourceType+".general_restriction_example").Key("general_device_restriction.compliant_app_list_type").HasValue("appsInListCompliant"),
					check.That(resourceType+".general_restriction_example").Key("general_device_restriction.compliant_apps_list.#").HasValue("2"),
					check.That(resourceType+".general_restriction_example").Key("general_device_restriction.safari_block_popups").HasValue("true"),
					check.That(resourceType+".general_restriction_example").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".general_restriction_example").Key("assignments.#").HasValue("4"),
				),
			},
		},
	})
}

func TestUnitResourceIosDeviceConfigurationTemplates_02_TrustedRootCertificate(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, configMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer configMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_trusted_root_certificate_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".trusted_cert_example").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".trusted_cert_example").Key("display_name").HasValue("unit-test-iOS-trusted-root-certificate-example"),
					check.That(resourceType+".trusted_cert_example").Key("description").HasValue("Install company root certificate for secure connections"),
					check.That(resourceType+".trusted_cert_example").Key("trusted_certificate.cert_file_name").HasValue("MicrosoftRootCertificateAuthority2011.cer"),
					check.That(resourceType+".trusted_cert_example").Key("trusted_certificate.trusted_root_certificate").Exists(),
					check.That(resourceType+".trusted_cert_example").Key("assignments.#").HasValue("4"),
				),
			},
		},
	})
}

func TestUnitResourceIosDeviceConfigurationTemplates_03_CreateWithError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, configMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer configMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_general_device_restriction_maximal.tf"),
				ExpectError: regexp.MustCompile(`(?s)Bad Request|400`),
			},
		},
	})
}

func TestUnitResourceIosDeviceConfigurationTemplates_04_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, configMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer configMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_general_device_restriction_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".general_restriction_example").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".general_restriction_example").Key("display_name").HasValue("unit-test-iOS-general-device-restriction-example"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_general_device_restriction_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".general_restriction_example").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".general_restriction_example").Key("display_name").HasValue("unit-test-iOS-general-device-restriction-example"),
				),
			},
		},
	})
}

func TestUnitResourceIosDeviceConfigurationTemplates_05_ImportState(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, configMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer configMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_trusted_root_certificate_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".trusted_cert_example").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				),
			},
			{
				ResourceName:      resourceType + ".trusted_cert_example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
