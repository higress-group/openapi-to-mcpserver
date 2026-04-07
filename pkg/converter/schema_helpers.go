package converter

import "github.com/getkin/kin-openapi/openapi3"

const openAPITypeNull = "null"

func schemaType(schema *openapi3.Schema) string {
	if schema == nil || schema.Type == nil {
		return ""
	}

	types := schema.Type.Slice()
	for _, typ := range types {
		if typ != openAPITypeNull {
			return typ
		}
	}

	if len(types) == 0 {
		return ""
	}

	return types[0]
}

func schemaRefType(schemaRef *openapi3.SchemaRef) string {
	if schemaRef == nil || schemaRef.Value == nil {
		return ""
	}

	return effectiveSchemaType(schemaRef.Value)
}

func effectiveSchemaType(schema *openapi3.Schema) string {
	if schema == nil {
		return ""
	}

	if typ := schemaType(schema); typ != "" {
		return typ
	}

	if len(schema.Properties) > 0 {
		return openapi3.TypeObject
	}

	if schema.Items != nil {
		return openapi3.TypeArray
	}

	return ""
}
