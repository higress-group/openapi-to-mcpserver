package converter

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/higress-group/openapi-to-mcpserver/pkg/models"
	"github.com/higress-group/openapi-to-mcpserver/pkg/parser"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestEndToEndConversion(t *testing.T) {
	// Test cases
	testCases := []struct {
		name           string
		inputFile      string
		expectedOutput string
		serverName     string
		templatePath   string
	}{
		{
			name:           "Petstore API",
			inputFile:      "../../test/petstore.json",
			expectedOutput: "../../test/expected-petstore-mcp.yaml",
			serverName:     "petstore",
		},
		{
			name:           "Path Parameters API",
			inputFile:      "../../test/path-params.json",
			expectedOutput: "../../test/expected-path-params-mcp.yaml",
			serverName:     "path-params-api",
		},
		{
			name:           "Header Parameters API",
			inputFile:      "../../test/header-params.json",
			expectedOutput: "../../test/expected-header-params-mcp.yaml",
			serverName:     "header-params-api",
		},
		{
			name:           "Cookie Parameters API",
			inputFile:      "../../test/cookie-params.json",
			expectedOutput: "../../test/expected-cookie-params-mcp.yaml",
			serverName:     "cookie-params-api",
		},
		{
			name:           "Request Body Types API",
			inputFile:      "../../test/request-body-types.json",
			expectedOutput: "../../test/expected-request-body-types-mcp.yaml",
			serverName:     "request-body-types-api",
		},
		{
			name:           "Petstore API with Template",
			inputFile:      "../../test/petstore.json",
			expectedOutput: "../../test/expected-petstore-template-mcp.yaml",
			serverName:     "petstore",
			templatePath:   "../../test/template.yaml",
		},
		{
			name:           "Security Schemes API",
			inputFile:      "../../test/security-test.json",
			expectedOutput: "../../test/expected-security-test-mcp.yaml",
			serverName:     "openapi-server", // Matches the default or can be specified if different
		},
		{
			name:           "Tools Args array of object",
			inputFile:      "../../test/tools-args-array-of-object.json",
			expectedOutput: "../../test/expected-tools-args-array-of-object-mcp.yaml",
			serverName:     "openapi-server",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new parser
			p := parser.NewParser()

			// Parse the OpenAPI specification
			err := p.ParseFile(tc.inputFile)
			assert.NoError(t, err)

			// Create a new converter
			c := NewConverter(p, models.ConvertOptions{
				ServerName:   tc.serverName,
				TemplatePath: tc.templatePath,
			})

			// Convert the OpenAPI specification to an MCP configuration
			config, err := c.Convert()
			assert.NoError(t, err)

			// Marshal the MCP configuration to YAML
			var buffer bytes.Buffer
			encoder := yaml.NewEncoder(&buffer)
			encoder.SetIndent(2)

			if err := encoder.Encode(config); err != nil {
				fmt.Printf("Error encoding YAML: %v\n", err)
				return
			}
			actualYAML := buffer.Bytes()
			assert.NoError(t, err)

			// If the expected output file doesn't exist, write the actual output to it
			if _, err := os.Stat(tc.expectedOutput); os.IsNotExist(err) {
				err = os.WriteFile(tc.expectedOutput, actualYAML, 0o644)
				assert.NoError(t, err)
				t.Logf("Created expected output file: %s", tc.expectedOutput)
			}

			// Read the expected output
			expectedYAML, err := os.ReadFile(tc.expectedOutput)
			assert.NoError(t, err)

			// Compare the actual and expected output
			assert.Equal(t, string(expectedYAML), string(actualYAML))
		})
	}
}
