package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v3"
)

// Parser represents an OpenAPI parser
type Parser struct {
	doc              *openapi3.T
	ValidateDocument bool
}

// NewParser creates a new OpenAPI parser
func NewParser() *Parser {
	return &Parser{
		ValidateDocument: false, // Default to no validation
	}
}

// SetValidation sets whether to validate the OpenAPI document
func (p *Parser) SetValidation(validate bool) {
	p.ValidateDocument = validate
}

// ParseFile parses an OpenAPI document from a file
func (p *Parser) ParseFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read OpenAPI file: %w", err)
	}

	return p.Parse(data)
}

// Parse parses an OpenAPI document from bytes
func (p *Parser) Parse(data []byte) error {
	loader := openapi3.NewLoader()

	normalizedData, err := normalizeOpenAPI31Document(data)
	if err != nil {
		return fmt.Errorf("failed to normalize OpenAPI document: %w", err)
	}

	// Try to parse as JSON first
	var doc *openapi3.T

	// Parse the document (loader can handle both JSON and YAML)
	doc, err = loader.LoadFromData(normalizedData)

	if err != nil {
		return fmt.Errorf("failed to parse OpenAPI document: %w", err)
	}

	// Validate the document if validation is enabled
	if p.ValidateDocument {
		err = doc.Validate(context.Background())
		if err != nil {
			return fmt.Errorf("invalid OpenAPI document: %w", err)
		}
	}

	p.doc = doc
	return nil
}

// GetDocument returns the parsed OpenAPI document
func (p *Parser) GetDocument() *openapi3.T {
	return p.doc
}

// GetPaths returns all paths in the OpenAPI document
func (p *Parser) GetPaths() map[string]*openapi3.PathItem {
	if p.doc == nil || p.doc.Paths == nil {
		return nil
	}
	return p.doc.Paths.Map()
}

// GetServers returns all servers in the OpenAPI document
func (p *Parser) GetServers() []*openapi3.Server {
	if p.doc == nil {
		return nil
	}
	return p.doc.Servers
}

// GetInfo returns the info section of the OpenAPI document
func (p *Parser) GetInfo() *openapi3.Info {
	if p.doc == nil {
		return nil
	}
	return p.doc.Info
}

// isJSON checks if the data is in JSON format
func isJSON(data []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(data, &js) == nil
}

var schemaObjectKeywords = map[string]struct{}{
	"additionalProperties":  {},
	"contains":              {},
	"else":                  {},
	"if":                    {},
	"items":                 {},
	"not":                   {},
	"propertyNames":         {},
	"then":                  {},
	"unevaluatedItems":      {},
	"unevaluatedProperties": {},
}

var schemaMapKeywords = map[string]struct{}{
	"$defs":             {},
	"dependentSchemas":  {},
	"patternProperties": {},
	"properties":        {},
}

var schemaArrayKeywords = map[string]struct{}{
	"allOf":       {},
	"anyOf":       {},
	"oneOf":       {},
	"prefixItems": {},
}

func normalizeOpenAPI31Document(data []byte) ([]byte, error) {
	var raw any
	var err error

	if isJSON(data) {
		err = json.Unmarshal(data, &raw)
	} else {
		err = yaml.Unmarshal(data, &raw)
	}
	if err != nil {
		return nil, err
	}

	root, ok := raw.(map[string]any)
	if !ok {
		return data, nil
	}

	version, _ := root["openapi"].(string)
	if !strings.HasPrefix(version, "3.1") {
		return data, nil
	}

	normalized, ok := normalizeDocumentValue(root).(map[string]any)
	if !ok {
		return data, nil
	}

	return json.Marshal(normalized)
}

func normalizeDocumentValue(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		normalized := make(map[string]any, len(typed))
		for key, child := range typed {
			switch {
			case hasKey(schemaObjectKeywords, key):
				normalized[key] = normalizePossibleSchema(child)
			case hasKey(schemaMapKeywords, key):
				normalized[key] = normalizeSchemaMap(child)
			case hasKey(schemaArrayKeywords, key):
				normalized[key] = normalizeSchemaArray(child)
			default:
				normalized[key] = normalizeDocumentValue(child)
			}
		}
		return normalized
	case []any:
		normalized := make([]any, len(typed))
		for i, child := range typed {
			normalized[i] = normalizeDocumentValue(child)
		}
		return normalized
	default:
		return value
	}
}

func normalizePossibleSchema(value any) any {
	switch typed := value.(type) {
	case bool:
		if typed {
			return map[string]any{}
		}
		return map[string]any{
			"not": map[string]any{},
		}
	default:
		return normalizeDocumentValue(value)
	}
}

func normalizeSchemaMap(value any) any {
	typed, ok := value.(map[string]any)
	if !ok {
		return normalizeDocumentValue(value)
	}

	normalized := make(map[string]any, len(typed))
	for key, child := range typed {
		normalized[key] = normalizePossibleSchema(child)
	}
	return normalized
}

func normalizeSchemaArray(value any) any {
	typed, ok := value.([]any)
	if !ok {
		return normalizeDocumentValue(value)
	}

	normalized := make([]any, len(typed))
	for i, child := range typed {
		normalized[i] = normalizePossibleSchema(child)
	}
	return normalized
}

func hasKey(values map[string]struct{}, key string) bool {
	_, ok := values[key]
	return ok
}

// GetOperationID generates an operation ID if one is not provided
func (p *Parser) GetOperationID(path string, method string, operation *openapi3.Operation) string {
	if operation.OperationID != "" {
		return operation.OperationID
	}

	// Generate an operation ID based on the path and method
	pathName := strings.ReplaceAll(path, "/", "_")
	pathName = strings.ReplaceAll(pathName, "{", "")
	pathName = strings.ReplaceAll(pathName, "}", "")
	return fmt.Sprintf("%s%s", strings.ToLower(method), pathName)
}
