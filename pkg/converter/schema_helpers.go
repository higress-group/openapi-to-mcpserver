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
	schema = unwrapSingleComposedSchema(schema)
	if schema == nil {
		return ""
	}

	if typ := schemaType(schema); typ != "" {
		return typ
	}

	if len(schema.Enum) > 0 {
		return openapi3.TypeString
	}

	if len(schema.Properties) > 0 {
		return openapi3.TypeObject
	}

	if schema.Items != nil {
		return openapi3.TypeArray
	}

	return ""
}

func unwrapSingleComposedSchema(schema *openapi3.Schema) *openapi3.Schema {
	if schema == nil {
		return nil
	}

	if schemaType(schema) != "" || len(schema.Enum) > 0 || len(schema.Properties) > 0 || schema.Items != nil {
		return schema
	}

	if ref := singleNonNullSchemaRef(schema.OneOf); ref != nil {
		return unwrapSingleComposedSchema(ref.Value)
	}

	if ref := singleNonNullSchemaRef(schema.AnyOf); ref != nil {
		return unwrapSingleComposedSchema(ref.Value)
	}

	if ref := singleNonNullSchemaRef(schema.AllOf); ref != nil {
		return unwrapSingleComposedSchema(ref.Value)
	}

	return schema
}

func singleNonNullSchemaRef(refs openapi3.SchemaRefs) *openapi3.SchemaRef {
	if len(refs) == 0 {
		return nil
	}

	var candidate *openapi3.SchemaRef
	for _, ref := range refs {
		if ref == nil || ref.Value == nil {
			return nil
		}

		if isNullOnlySchema(ref.Value) {
			continue
		}

		if candidate != nil {
			return nil
		}

		candidate = ref
	}

	return candidate
}

func isNullOnlySchema(schema *openapi3.Schema) bool {
	if schema == nil || schema.Type == nil {
		return false
	}

	types := schema.Type.Slice()
	return len(types) == 1 && types[0] == openAPITypeNull
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}

	return ""
}
