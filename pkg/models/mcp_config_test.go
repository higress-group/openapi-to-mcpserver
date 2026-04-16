package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestServerConfigPassthroughAuthHeader(t *testing.T) {
	testCases := []struct {
		name          string
		config        ServerConfig
		expectedYAML  string
		expectedValue *bool
	}{
		{
			name: "passthroughAuthHeader set to true",
			config: ServerConfig{
				Name:                  "test-server",
				Type:                  "mcp-proxy",
				Transport:             "sse",
				McpServerURL:          "http://backend.example.com/mcp",
				PassthroughAuthHeader: ptrBool(true),
			},
			expectedYAML: `name: test-server
type: mcp-proxy
transport: sse
mcpServerURL: http://backend.example.com/mcp
passthroughAuthHeader: true
`,
			expectedValue: ptrBool(true),
		},
		{
			name: "passthroughAuthHeader set to false",
			config: ServerConfig{
				Name:                  "test-server",
				Type:                  "mcp-proxy",
				Transport:             "sse",
				McpServerURL:          "http://backend.example.com/mcp",
				PassthroughAuthHeader: ptrBool(false),
			},
			expectedYAML: `name: test-server
type: mcp-proxy
transport: sse
mcpServerURL: http://backend.example.com/mcp
passthroughAuthHeader: false
`,
			expectedValue: ptrBool(false),
		},
		{
			name: "passthroughAuthHeader nil (not set)",
			config: ServerConfig{
				Name:                  "test-server",
				Type:                  "mcp-proxy",
				Transport:             "sse",
				McpServerURL:          "http://backend.example.com/mcp",
				PassthroughAuthHeader: nil,
			},
			expectedYAML: `name: test-server
type: mcp-proxy
transport: sse
mcpServerURL: http://backend.example.com/mcp
`,
			expectedValue: nil,
		},
		{
			name: "passthroughAuthHeader with other fields",
			config: ServerConfig{
				Name:                      "test-server",
				Type:                      "mcp-proxy",
				Transport:                 "sse",
				McpServerURL:              "http://backend.example.com/mcp",
				Timeout:                   ptrInt(60000),
				DefaultDownstreamSecurity: &ToolSecurityRequirement{ID: "client-auth"},
				DefaultUpstreamSecurity:   &ToolSecurityRequirement{ID: "backend-auth", Passthrough: true},
				PassthroughAuthHeader:     ptrBool(true),
			},
			expectedYAML: `name: test-server
type: mcp-proxy
transport: sse
mcpServerURL: http://backend.example.com/mcp
timeout: 60000
defaultUpstreamSecurity:
    id: backend-auth
    passthrough: true
defaultDownstreamSecurity:
    id: client-auth
passthroughAuthHeader: true
`,
			expectedValue: ptrBool(true),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test YAML serialization
			yamlBytes, err := yaml.Marshal(tc.config)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedYAML, string(yamlBytes))

			// Test YAML deserialization
			var unmarshaledConfig ServerConfig
			err = yaml.Unmarshal(yamlBytes, &unmarshaledConfig)
			assert.NoError(t, err)

			// Verify PassthroughAuthHeader value
			if tc.expectedValue == nil {
				assert.Nil(t, unmarshaledConfig.PassthroughAuthHeader)
			} else {
				assert.NotNil(t, unmarshaledConfig.PassthroughAuthHeader)
				assert.Equal(t, *tc.expectedValue, *unmarshaledConfig.PassthroughAuthHeader)
			}
		})
	}
}

func TestMCPConfigPassthroughAuthHeader(t *testing.T) {
	testCases := []struct {
		name         string
		config       MCPConfig
		expectedYAML string
	}{
		{
			name: "MCPConfig with passthroughAuthHeader true",
			config: MCPConfig{
				Server: ServerConfig{
					Name:                  "proxy-server",
					Type:                  "mcp-proxy",
					Transport:             "sse",
					McpServerURL:          "http://localhost:8080/mcp",
					PassthroughAuthHeader: ptrBool(true),
				},
				Tools: []Tool{},
			},
			expectedYAML: `server:
    name: proxy-server
    type: mcp-proxy
    transport: sse
    mcpServerURL: http://localhost:8080/mcp
    passthroughAuthHeader: true
`,
		},
		{
			name: "MCPConfig with passthroughAuthHeader false",
			config: MCPConfig{
				Server: ServerConfig{
					Name:                  "proxy-server",
					Type:                  "mcp-proxy",
					Transport:             "sse",
					McpServerURL:          "http://localhost:8080/mcp",
					PassthroughAuthHeader: ptrBool(false),
				},
				Tools: []Tool{},
			},
			expectedYAML: `server:
    name: proxy-server
    type: mcp-proxy
    transport: sse
    mcpServerURL: http://localhost:8080/mcp
    passthroughAuthHeader: false
`,
		},
		{
			name: "MCPConfig without passthroughAuthHeader",
			config: MCPConfig{
				Server: ServerConfig{
					Name:         "proxy-server",
					Type:         "mcp-proxy",
					Transport:    "sse",
					McpServerURL: "http://localhost:8080/mcp",
				},
				Tools: []Tool{},
			},
			expectedYAML: `server:
    name: proxy-server
    type: mcp-proxy
    transport: sse
    mcpServerURL: http://localhost:8080/mcp
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test YAML serialization
			yamlBytes, err := yaml.Marshal(tc.config)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedYAML, string(yamlBytes))

			// Test YAML deserialization
			var unmarshaledConfig MCPConfig
			err = yaml.Unmarshal(yamlBytes, &unmarshaledConfig)
			assert.NoError(t, err)

			// Verify the config can be properly unmarshaled
			assert.Equal(t, tc.config.Server.Name, unmarshaledConfig.Server.Name)
			assert.Equal(t, tc.config.Server.Type, unmarshaledConfig.Server.Type)
			assert.Equal(t, tc.config.Server.Transport, unmarshaledConfig.Server.Transport)
			assert.Equal(t, tc.config.Server.McpServerURL, unmarshaledConfig.Server.McpServerURL)

			if tc.config.Server.PassthroughAuthHeader != nil {
				assert.NotNil(t, unmarshaledConfig.Server.PassthroughAuthHeader)
				assert.Equal(t, *tc.config.Server.PassthroughAuthHeader, *unmarshaledConfig.Server.PassthroughAuthHeader)
			} else {
				assert.Nil(t, unmarshaledConfig.Server.PassthroughAuthHeader)
			}
		})
	}
}

func TestPassthroughAuthHeaderParseFromYAML(t *testing.T) {
	testCases := []struct {
		name        string
		yamlInput   string
		expectedVal *bool
	}{
		{
			name:        "parse true value",
			yamlInput:   "passthroughAuthHeader: true",
			expectedVal: ptrBool(true),
		},
		{
			name:        "parse false value",
			yamlInput:   "passthroughAuthHeader: false",
			expectedVal: ptrBool(false),
		},
		{
			name:        "parse nil (field not present)",
			yamlInput:   "name: test-server",
			expectedVal: nil,
		},
		{
			name: "parse with other server config fields",
			yamlInput: `name: test-server
type: mcp-proxy
transport: sse
mcpServerURL: http://backend/mcp
timeout: 30000
passthroughAuthHeader: true
`,
			expectedVal: ptrBool(true),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var config ServerConfig
			err := yaml.Unmarshal([]byte(tc.yamlInput), &config)
			assert.NoError(t, err)

			if tc.expectedVal == nil {
				assert.Nil(t, config.PassthroughAuthHeader)
			} else {
				assert.NotNil(t, config.PassthroughAuthHeader)
				assert.Equal(t, *tc.expectedVal, *config.PassthroughAuthHeader)
			}
		})
	}
}

// TestPassthroughAuthHeaderBackwardCompatibility tests backward compatibility:
// When the field is not present (nil), it should be treated as false (default behavior).
// This ensures old MCP configurations without this field continue to work correctly.
func TestPassthroughAuthHeaderBackwardCompatibility(t *testing.T) {
	testCases := []struct {
		name            string
		yamlInput       string
		fieldIsNil      bool
		semanticDefault bool // The default behavior when nil (should be false)
		description     string
	}{
		{
			name: "old config without passthroughAuthHeader - should default to false",
			yamlInput: `name: legacy-server
type: mcp-proxy
transport: sse
mcpServerURL: http://backend:8080/mcp
timeout: 60000
`,
			fieldIsNil:      true,
			semanticDefault: false,
			description:     "Legacy configs should default to NOT passing through auth header",
		},
		{
			name: "old config with only basic fields - should default to false",
			yamlInput: `name: minimal-server
mcpServerURL: http://localhost/mcp
`,
			fieldIsNil:      true,
			semanticDefault: false,
			description:     "Minimal configs should default to NOT passing through auth header",
		},
		{
			name: "new config explicitly set to false",
			yamlInput: `name: explicit-false-server
passthroughAuthHeader: false
`,
			fieldIsNil:      false,
			semanticDefault: false,
			description:     "Explicit false should be respected",
		},
		{
			name: "new config explicitly set to true",
			yamlInput: `name: explicit-true-server
passthroughAuthHeader: true
`,
			fieldIsNil:      false,
			semanticDefault: true,
			description:     "Explicit true should be respected",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var config ServerConfig
			err := yaml.Unmarshal([]byte(tc.yamlInput), &config)
			assert.NoError(t, err, "YAML parsing should succeed")

			// Verify the field is nil/not nil as expected
			if tc.fieldIsNil {
				assert.Nil(t, config.PassthroughAuthHeader, tc.description)
			} else {
				assert.NotNil(t, config.PassthroughAuthHeader, tc.description)
			}

			// Verify semantic behavior: nil should be treated as false
			actualBehavior := getPassthroughAuthHeaderBehavior(config)
			assert.Equal(t, tc.semanticDefault, actualBehavior, tc.description)
		})
	}
}

// getPassthroughAuthHeaderBehavior returns the semantic behavior of PassthroughAuthHeader.
// When the field is nil (not set), it defaults to false (NOT passing through auth header).
// This helper function demonstrates how the field should be used in actual code.
func getPassthroughAuthHeaderBehavior(config ServerConfig) bool {
	if config.PassthroughAuthHeader == nil {
		return false // Default: do NOT pass through auth header
	}
	return *config.PassthroughAuthHeader
}

// Helper functions
func ptrBool(v bool) *bool {
	return &v
}

func ptrInt(v int) *int {
	return &v
}
