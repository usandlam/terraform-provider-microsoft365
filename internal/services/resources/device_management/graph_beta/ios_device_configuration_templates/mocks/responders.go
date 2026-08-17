// Package mocks provides HTTP responders for testing the iOS device
// configuration templates resource. It registers in-memory CRUD responders
// for /deviceManagement/deviceConfigurations plus supporting dependency
// endpoints (groups, roleScopeTags, assignmentFilters) so unit and
// acceptance tests can be written without hitting Graph.
package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	deviceConfigurations map[string]map[string]any
	assignments          map[string][]any
}

func init() {
	mockState.deviceConfigurations = make(map[string]map[string]any)
	mockState.assignments = make(map[string][]any)
	mocks.GlobalRegistry.Register("ios_device_configuration_templates", &IosDeviceConfigurationTemplatesMock{})
}

// IosDeviceConfigurationTemplatesMock registers httpmock responders for the
// deviceConfigurations endpoints. Handles create/read/patch/delete and the
// /assign sub-path for both iosGeneralDeviceConfiguration and
// iosTrustedRootCertificate variants.
type IosDeviceConfigurationTemplatesMock struct{}

var _ mocks.MockRegistrar = (*IosDeviceConfigurationTemplatesMock)(nil)

func (m *IosDeviceConfigurationTemplatesMock) RegisterMocks() {
	mockState.Lock()
	mockState.deviceConfigurations = make(map[string]map[string]any)
	mockState.assignments = make(map[string][]any)
	mockState.Unlock()

	m.registerDependencyMocks()

	baseURL := "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations"

	httpmock.RegisterResponder("POST", baseURL, func(req *http.Request) (*http.Response, error) {
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"invalid body"}}`), nil
		}

		odataType, _ := body["@odata.type"].(string)
		if odataType != "#microsoft.graph.iosGeneralDeviceConfiguration" &&
			odataType != "#microsoft.graph.iosTrustedRootCertificate" {
			return httpmock.NewStringResponse(400, fmt.Sprintf(`{"error":{"code":"BadRequest","message":"unsupported @odata.type %s"}}`, odataType)), nil
		}

		id := uuid.NewString()
		body["id"] = id

		mockState.Lock()
		mockState.deviceConfigurations[id] = body
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, body)
	})

	httpmock.RegisterRegexpResponder("GET",
		mustRegex(`^https://graph\.microsoft\.com/beta/deviceManagement/deviceConfigurations/([^/?]+)(\?.*)?$`),
		func(req *http.Request) (*http.Response, error) {
			id := extractID(req.URL.Path)
			mockState.Lock()
			defer mockState.Unlock()
			cfg, ok := mockState.deviceConfigurations[id]
			if !ok {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"not found"}}`), nil
			}
			resp := make(map[string]any, len(cfg)+1)
			for k, v := range cfg {
				resp[k] = v
			}
			assignments := mockState.assignments[id]
			if assignments == nil {
				assignments = []any{}
			}
			resp["assignments"] = assignments
			return httpmock.NewJsonResponse(200, resp)
		},
	)

	httpmock.RegisterRegexpResponder("PATCH",
		mustRegex(`^https://graph\.microsoft\.com/beta/deviceManagement/deviceConfigurations/([^/?]+)$`),
		func(req *http.Request) (*http.Response, error) {
			id := extractID(req.URL.Path)
			var body map[string]any
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"invalid body"}}`), nil
			}
			mockState.Lock()
			defer mockState.Unlock()
			cfg, ok := mockState.deviceConfigurations[id]
			if !ok {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"not found"}}`), nil
			}
			for k, v := range body {
				cfg[k] = v
			}
			cfg["id"] = id
			mockState.deviceConfigurations[id] = cfg
			return httpmock.NewJsonResponse(200, cfg)
		},
	)

	httpmock.RegisterRegexpResponder("DELETE",
		mustRegex(`^https://graph\.microsoft\.com/beta/deviceManagement/deviceConfigurations/([^/?]+)$`),
		func(req *http.Request) (*http.Response, error) {
			id := extractID(req.URL.Path)
			mockState.Lock()
			defer mockState.Unlock()
			delete(mockState.deviceConfigurations, id)
			delete(mockState.assignments, id)
			return httpmock.NewStringResponse(204, ""), nil
		},
	)

	httpmock.RegisterRegexpResponder("POST",
		mustRegex(`^https://graph\.microsoft\.com/beta/deviceManagement/deviceConfigurations/([^/?]+)/assign$`),
		func(req *http.Request) (*http.Response, error) {
			id := extractIDFromAssignPath(req.URL.Path)
			var body map[string]any
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"invalid body"}}`), nil
			}
			assignments, _ := body["assignments"].([]any)
			mockState.Lock()
			mockState.assignments[id] = assignments
			mockState.Unlock()
			return httpmock.NewJsonResponse(200, map[string]any{"value": assignments})
		},
	)
}

// RegisterErrorMocks registers responders that always return errors, for
// tests that exercise failure paths.
func (m *IosDeviceConfigurationTemplatesMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.deviceConfigurations = make(map[string]map[string]any)
	mockState.assignments = make(map[string][]any)
	mockState.Unlock()

	m.registerDependencyMocks()

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations",
		httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Error creating iOS device configuration template"}}`))
	httpmock.RegisterRegexpResponder("GET",
		mustRegex(`^https://graph\.microsoft\.com/beta/deviceManagement/deviceConfigurations/[^/?]+(\?.*)?$`),
		httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"not found"}}`))
	httpmock.RegisterRegexpResponder("PATCH",
		mustRegex(`^https://graph\.microsoft\.com/beta/deviceManagement/deviceConfigurations/[^/?]+$`),
		httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Error updating iOS device configuration template"}}`))
	httpmock.RegisterRegexpResponder("DELETE",
		mustRegex(`^https://graph\.microsoft\.com/beta/deviceManagement/deviceConfigurations/[^/?]+$`),
		httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Error deleting iOS device configuration template"}}`))
	httpmock.RegisterRegexpResponder("POST",
		mustRegex(`^https://graph\.microsoft\.com/beta/deviceManagement/deviceConfigurations/[^/?]+/assign$`),
		httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Error assigning iOS device configuration template"}}`))
}

// CleanupMockState resets the in-memory mock state between tests.
func (m *IosDeviceConfigurationTemplatesMock) CleanupMockState() {
	mockState.Lock()
	mockState.deviceConfigurations = make(map[string]map[string]any)
	mockState.assignments = make(map[string][]any)
	mockState.Unlock()
}

// registerDependencyMocks registers stub responders for the groups,
// roleScopeTags, and assignmentFilters endpoints referenced by assignment
// blocks. These do not affect the resource-under-test but keep the Graph
// SDK from erroring on incidental lookups.
func (m *IosDeviceConfigurationTemplatesMock) registerDependencyMocks() {
	httpmock.RegisterRegexpResponder("GET",
		mustRegex(`^https://graph\.microsoft\.com/beta/deviceManagement/roleScopeTags/([^/]+)$`),
		func(req *http.Request) (*http.Response, error) {
			tagID := lastSegment(req.URL.Path)
			return httpmock.NewJsonResponse(200, map[string]any{
				"@odata.type": "#microsoft.graph.roleScopeTag",
				"id":          tagID,
				"displayName": fmt.Sprintf("Role Scope Tag %s", tagID),
				"description": "Test role scope tag",
			})
		},
	)

	httpmock.RegisterRegexpResponder("GET",
		mustRegex(`^https://graph\.microsoft\.com/beta/groups/([^/]+)$`),
		func(req *http.Request) (*http.Response, error) {
			groupID := lastSegment(req.URL.Path)
			return httpmock.NewJsonResponse(200, map[string]any{
				"@odata.type":     "#microsoft.graph.group",
				"id":              groupID,
				"displayName":     fmt.Sprintf("Test Group %s", groupID),
				"description":     "Test group for iOS device configuration",
				"groupTypes":      []string{},
				"securityEnabled": true,
			})
		},
	)

	httpmock.RegisterRegexpResponder("GET",
		mustRegex(`^https://graph\.microsoft\.com/beta/deviceManagement/assignmentFilters/([^/]+)$`),
		func(req *http.Request) (*http.Response, error) {
			filterID := lastSegment(req.URL.Path)
			return httpmock.NewJsonResponse(200, map[string]any{
				"@odata.type":   "#microsoft.graph.deviceAndAppManagementAssignmentFilter",
				"id":            filterID,
				"displayName":   fmt.Sprintf("Test Assignment Filter %s", filterID),
				"description":   "Test assignment filter",
				"platform":      "iOS",
				"rule":          "(device.deviceOwnership -eq \"Corporate\")",
				"roleScopeTags": []string{"0"},
			})
		},
	)

	httpmock.RegisterRegexpResponder("GET",
		mustRegex(`^https://graph\.microsoft\.com/beta/deviceManagement/deviceConfigurations/([^/]+)/assignments$`),
		func(req *http.Request) (*http.Response, error) {
			configID := secondToLastSegment(req.URL.Path)
			mockState.Lock()
			assignments := mockState.assignments[configID]
			mockState.Unlock()
			if assignments == nil {
				assignments = []any{}
			}
			return httpmock.NewJsonResponse(200, map[string]any{
				"@odata.context": fmt.Sprintf("https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceConfigurations('%s')/assignments", configID),
				"value":          assignments,
			})
		},
	)
}

func mustRegex(pattern string) *regexp.Regexp {
	return regexp.MustCompile(pattern)
}

func extractID(path string) string {
	return lastSegment(path)
}

func extractIDFromAssignPath(path string) string {
	trimmed := strings.TrimSuffix(path, "/assign")
	return lastSegment(trimmed)
}

func lastSegment(path string) string {
	parts := strings.Split(strings.TrimSuffix(path, "/"), "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

func secondToLastSegment(path string) string {
	parts := strings.Split(strings.TrimSuffix(path, "/"), "/")
	if len(parts) < 2 {
		return ""
	}
	return parts[len(parts)-2]
}
