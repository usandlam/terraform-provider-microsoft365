package graphBetaIosDeviceConfigurationTemplates

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ModifyPlan validates that exactly one configuration type is specified.
func (r *IosDeviceConfigurationTemplatesResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan IosDeviceConfigurationTemplatesResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	configCount := 0
	if !plan.GeneralDeviceRestriction.IsNull() && !plan.GeneralDeviceRestriction.IsUnknown() {
		configCount++
	}
	if !plan.TrustedCertificate.IsNull() && !plan.TrustedCertificate.IsUnknown() {
		configCount++
	}

	switch {
	case configCount == 0:
		resp.Diagnostics.AddError(
			"Missing Configuration",
			"Exactly one of general_device_restriction or trusted_certificate must be specified.",
		)
	case configCount > 1:
		resp.Diagnostics.AddError(
			"Multiple Configurations",
			"Only one of general_device_restriction or trusted_certificate may be specified.",
		)
	}
}
