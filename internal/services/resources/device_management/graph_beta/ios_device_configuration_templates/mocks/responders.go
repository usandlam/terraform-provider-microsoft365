// Package mocks provides HTTP responders for testing the iOS device
// configuration templates resource. This scaffold registers minimal in-memory
// state and CRUD responders for `/deviceManagement/deviceConfigurations`
// endpoints so acceptance tests can be written without hitting Graph.
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

	baseURL := "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations"

	// POST /deviceManagement/deviceConfigurations
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

	// GET/PATCH/DELETE and /assign for a specific id
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
			// Attach assignments if expand=assignments was requested (we always attach; caller can ignore).
			resp := make(map[string]any, len(cfg)+1)
			for k, v := range cfg {
				resp[k] = v
			}
			resp["assignments"] = mockState.assignments[id]
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

func mustRegex(pattern string) *regexp.Regexp {
	return regexp.MustCompile(pattern)
}

func extractID(path string) string {
	// path is /beta/deviceManagement/deviceConfigurations/{id}
	parts := strings.Split(strings.TrimSuffix(path, "/"), "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

func extractIDFromAssignPath(path string) string {
	// path is /beta/deviceManagement/deviceConfigurations/{id}/assign
	trimmed := strings.TrimSuffix(path, "/assign")
	parts := strings.Split(strings.TrimSuffix(trimmed, "/"), "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}
