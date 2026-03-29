package parser

import "testing"

func TestParseOpenAPI31TypeArray(t *testing.T) {
	t.Parallel()

	spec := []byte(`{
		"openapi": "3.1.0",
		"info": {
			"title": "Nullable field API",
			"version": "1.0.0"
		},
		"paths": {
			"/users": {
				"post": {
					"operationId": "createUser",
					"requestBody": {
						"required": true,
						"content": {
							"application/json": {
								"schema": {
									"type": "object",
									"properties": {
										"nickname": {
											"type": ["string", "null"],
											"description": "Optional nickname"
										}
									}
								}
							}
						}
					},
					"responses": {
						"200": {
							"description": "ok"
						}
					}
				}
			}
		}
	}`)

	parser := NewParser()
	if err := parser.Parse(spec); err != nil {
		t.Fatalf("Parse() returned error for OpenAPI 3.1 type array: %v", err)
	}
}
