package graphBetaIosDeviceConfigurationTemplates

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// IosConfigurationTemplatesAssignmentType returns the object type for the
// assignments Set elements.
func IosConfigurationTemplatesAssignmentType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":        types.StringType,
			"group_id":    types.StringType,
			"filter_id":   types.StringType,
			"filter_type": types.StringType,
		},
	}
}

func mapAssignmentsToTerraform(ctx context.Context, data *IosDeviceConfigurationTemplatesResourceModel, assignments []graphmodels.DeviceConfigurationAssignmentable) {
	if len(assignments) == 0 {
		data.Assignments = types.SetNull(IosConfigurationTemplatesAssignmentType())
		return
	}

	assignmentValues := []attr.Value{}

	for i, assignment := range assignments {
		target := assignment.GetTarget()
		if target == nil {
			tflog.Warn(ctx, "Assignment target is nil, skipping", map[string]any{"index": i})
			continue
		}

		odataType := target.GetOdataType()
		if odataType == nil {
			tflog.Warn(ctx, "Assignment target OData type is nil, skipping", map[string]any{"index": i})
			continue
		}

		assignmentObj := map[string]attr.Value{
			"type":        types.StringNull(),
			"group_id":    types.StringNull(),
			"filter_id":   types.StringNull(),
			"filter_type": types.StringNull(),
		}

		switch *odataType {
		case "#microsoft.graph.allDevicesAssignmentTarget":
			assignmentObj["type"] = types.StringValue("allDevicesAssignmentTarget")

		case "#microsoft.graph.allLicensedUsersAssignmentTarget":
			assignmentObj["type"] = types.StringValue("allLicensedUsersAssignmentTarget")

		case "#microsoft.graph.groupAssignmentTarget":
			assignmentObj["type"] = types.StringValue("groupAssignmentTarget")
			if groupTarget, ok := target.(graphmodels.GroupAssignmentTargetable); ok {
				if groupId := groupTarget.GetGroupId(); groupId != nil && *groupId != "" {
					assignmentObj["group_id"] = convert.GraphToFrameworkString(groupId)
				}
			}

		case "#microsoft.graph.exclusionGroupAssignmentTarget":
			assignmentObj["type"] = types.StringValue("exclusionGroupAssignmentTarget")
			if groupTarget, ok := target.(graphmodels.ExclusionGroupAssignmentTargetable); ok {
				if groupId := groupTarget.GetGroupId(); groupId != nil && *groupId != "" {
					assignmentObj["group_id"] = convert.GraphToFrameworkString(groupId)
				}
			}

		default:
			tflog.Warn(ctx, "Unknown target type encountered", map[string]any{"index": i, "targetType": *odataType})
		}

		filterID := target.GetDeviceAndAppManagementAssignmentFilterId()
		if filterID != nil && *filterID != "" && *filterID != "00000000-0000-0000-0000-000000000000" {
			assignmentObj["filter_id"] = convert.GraphToFrameworkString(filterID)
		} else {
			assignmentObj["filter_id"] = types.StringValue("00000000-0000-0000-0000-000000000000")
		}

		filterType := target.GetDeviceAndAppManagementAssignmentFilterType()
		if filterType != nil {
			switch *filterType {
			case graphmodels.INCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE:
				assignmentObj["filter_type"] = types.StringValue("include")
			case graphmodels.EXCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE:
				assignmentObj["filter_type"] = types.StringValue("exclude")
			case graphmodels.NONE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE:
				assignmentObj["filter_type"] = types.StringValue("none")
			default:
				assignmentObj["filter_type"] = types.StringValue("none")
			}
		} else {
			assignmentObj["filter_type"] = types.StringValue("none")
		}

		objValue, diags := types.ObjectValue(IosConfigurationTemplatesAssignmentType().(types.ObjectType).AttrTypes, assignmentObj)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to create assignment object value", map[string]any{"index": i, "errors": diags.Errors()})
			continue
		}
		assignmentValues = append(assignmentValues, objValue)
	}

	if len(assignmentValues) == 0 {
		data.Assignments = types.SetNull(IosConfigurationTemplatesAssignmentType())
		return
	}

	setVal, diags := types.SetValue(IosConfigurationTemplatesAssignmentType(), assignmentValues)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to create assignments set", map[string]any{"errors": diags.Errors()})
		data.Assignments = types.SetNull(IosConfigurationTemplatesAssignmentType())
		return
	}
	data.Assignments = setVal
}
