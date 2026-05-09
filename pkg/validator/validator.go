package validator

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// MCPCompatMode defines how to handle incompatible operations
type MCPCompatMode int

const (
	// MCPCompatOff disables MCP compatibility checking
	MCPCompatOff MCPCompatMode = iota
	// MCPCompatStrict reports errors and stops on incompatible operations
	MCPCompatStrict
	// MCPCompatWarn skips incompatible operations with warnings
	MCPCompatWarn
)

// ParseMCPCompatMode parses a string into MCPCompatMode
func ParseMCPCompatMode(s string) (MCPCompatMode, bool) {
	switch s {
	case "strict":
		return MCPCompatStrict, true
	case "warn":
		return MCPCompatWarn, true
	case "":
		return MCPCompatOff, true
	default:
		return MCPCompatOff, false
	}
}

// MCPCompatIssue represents a single compatibility issue found during validation
type MCPCompatIssue struct {
	Path        string // API path, e.g., "/pets/{id}/upload"
	Method      string // HTTP method, e.g., "POST"
	OperationID string // Operation ID if available
	Reason      string // Human-readable description of the issue
}

func (i MCPCompatIssue) String() string {
	op := i.OperationID
	if op == "" {
		op = fmt.Sprintf("%s %s", strings.ToUpper(i.Method), i.Path)
	}
	return fmt.Sprintf("[MCP incompatible] %s: %s", op, i.Reason)
}

// MCPCompatResult holds the validation result for the entire document
type MCPCompatResult struct {
	Issues            []MCPCompatIssue
	IncompatiblePaths map[string]map[string]bool // path -> method -> true
}

// IsCompatible returns true if a specific path+method combination is compatible
func (r *MCPCompatResult) IsCompatible(path, method string) bool {
	if r == nil || r.IncompatiblePaths == nil {
		return true
	}
	methods, exists := r.IncompatiblePaths[path]
	if !exists {
		return true
	}
	return !methods[strings.ToLower(method)]
}

// unsupportedContentTypes lists content types that cannot be represented as JSON input in MCP
var unsupportedContentTypes = []string{
	"application/octet-stream",
	"text/event-stream",
	"application/xml",
	"text/xml",
	"image/",
	"audio/",
	"video/",
}

// ValidateDocument checks all operations in an OpenAPI document for MCP compatibility.
// Returns the validation result containing all issues found.
func ValidateDocument(doc *openapi3.T) *MCPCompatResult {
	result := &MCPCompatResult{
		Issues:            []MCPCompatIssue{},
		IncompatiblePaths: make(map[string]map[string]bool),
	}

	if doc == nil || doc.Paths == nil {
		return result
	}

	for path, pathItem := range doc.Paths {
		operations := getOperations(pathItem)
		for method, operation := range operations {
			issues := checkOperation(path, method, operation)
			if len(issues) > 0 {
				result.Issues = append(result.Issues, issues...)
				if result.IncompatiblePaths[path] == nil {
					result.IncompatiblePaths[path] = make(map[string]bool)
				}
				result.IncompatiblePaths[path][method] = true
			}
		}
	}

	return result
}

// getOperations returns a map of HTTP method to operation
func getOperations(pathItem *openapi3.PathItem) map[string]*openapi3.Operation {
	operations := make(map[string]*openapi3.Operation)
	if pathItem.Get != nil {
		operations["get"] = pathItem.Get
	}
	if pathItem.Post != nil {
		operations["post"] = pathItem.Post
	}
	if pathItem.Put != nil {
		operations["put"] = pathItem.Put
	}
	if pathItem.Delete != nil {
		operations["delete"] = pathItem.Delete
	}
	if pathItem.Options != nil {
		operations["options"] = pathItem.Options
	}
	if pathItem.Head != nil {
		operations["head"] = pathItem.Head
	}
	if pathItem.Patch != nil {
		operations["patch"] = pathItem.Patch
	}
	if pathItem.Trace != nil {
		operations["trace"] = pathItem.Trace
	}
	return operations
}

// checkOperation performs all compatibility checks on a single operation
func checkOperation(path, method string, operation *openapi3.Operation) []MCPCompatIssue {
	var issues []MCPCompatIssue

	operationID := ""
	if operation.OperationID != "" {
		operationID = operation.OperationID
	}

	// Check request body content types and schemas
	if operation.RequestBody != nil && operation.RequestBody.Value != nil {
		requestBody := operation.RequestBody.Value
		for contentType, mediaType := range requestBody.Content {
			// Check for unsupported content types
			if issue := checkContentType(path, method, operationID, contentType); issue != nil {
				issues = append(issues, *issue)
				continue
			}

			// Check for binary/stream formats in schema
			if mediaType.Schema != nil && mediaType.Schema.Value != nil {
				schemaIssues := checkSchemaForBinary(path, method, operationID, contentType, mediaType.Schema.Value)
				issues = append(issues, schemaIssues...)
			}
		}
	}

	// Check response for streaming (text/event-stream)
	if operation.Responses != nil {
		for code, responseRef := range operation.Responses {
			if responseRef == nil || responseRef.Value == nil {
				continue
			}
			for contentType := range responseRef.Value.Content {
				if strings.Contains(contentType, "text/event-stream") {
					issues = append(issues, MCPCompatIssue{
						Path:        path,
						Method:      method,
						OperationID: operationID,
						Reason:      fmt.Sprintf("response (%s) uses streaming content type %q, MCP tool calls are synchronous request-response", code, contentType),
					})
				}
			}
		}
	}

	return issues
}

// checkContentType checks if a request body content type is supported by MCP
func checkContentType(path, method, operationID, contentType string) *MCPCompatIssue {
	ct := strings.ToLower(contentType)

	for _, unsupported := range unsupportedContentTypes {
		if strings.Contains(ct, unsupported) {
			return &MCPCompatIssue{
				Path:        path,
				Method:      method,
				OperationID: operationID,
				Reason:      fmt.Sprintf("request body content type %q is not supported by MCP (only JSON-serializable inputs are allowed)", contentType),
			}
		}
	}

	// multipart/form-data is generally problematic for MCP
	if strings.Contains(ct, "multipart/form-data") {
		return &MCPCompatIssue{
			Path:        path,
			Method:      method,
			OperationID: operationID,
			Reason:      fmt.Sprintf("request body content type %q is not supported by MCP (multipart form data cannot be represented as JSON tool arguments)", contentType),
		}
	}

	return nil
}

// checkSchemaForBinary checks a schema for binary/stream format fields
func checkSchemaForBinary(path, method, operationID, contentType string, schema *openapi3.Schema) []MCPCompatIssue {
	var issues []MCPCompatIssue

	// Check top-level schema format
	if schema.Format == "binary" {
		issues = append(issues, MCPCompatIssue{
			Path:        path,
			Method:      method,
			OperationID: operationID,
			Reason:      fmt.Sprintf("request body (%s) uses format \"binary\" which cannot be transmitted as MCP tool input", contentType),
		})
		return issues
	}

	// Check object properties for binary fields
	if schema.Type == "object" {
		for propName, propRef := range schema.Properties {
			if propRef.Value == nil {
				continue
			}
			if propRef.Value.Format == "binary" {
				issues = append(issues, MCPCompatIssue{
					Path:        path,
					Method:      method,
					OperationID: operationID,
					Reason:      fmt.Sprintf("request body field %q uses format \"binary\" which cannot be transmitted as MCP tool input", propName),
				})
			}
			// Check array items with binary format
			if propRef.Value.Type == "array" && propRef.Value.Items != nil && propRef.Value.Items.Value != nil {
				if propRef.Value.Items.Value.Format == "binary" {
					issues = append(issues, MCPCompatIssue{
						Path:        path,
						Method:      method,
						OperationID: operationID,
						Reason:      fmt.Sprintf("request body field %q is an array of binary items which cannot be transmitted as MCP tool input", propName),
					})
				}
			}
		}
	}

	// Check array items
	if schema.Type == "array" && schema.Items != nil && schema.Items.Value != nil {
		if schema.Items.Value.Format == "binary" {
			issues = append(issues, MCPCompatIssue{
				Path:        path,
				Method:      method,
				OperationID: operationID,
				Reason:      fmt.Sprintf("request body (%s) is an array of binary items which cannot be transmitted as MCP tool input", contentType),
			})
		}
	}

	return issues
}
