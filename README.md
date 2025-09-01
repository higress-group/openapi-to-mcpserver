# OpenAPI to MCP Server

A tool to convert OpenAPI specifications to MCP (Model Context Protocol) server configurations.

## Installation

```bash
go install github.com/higress-group/openapi-to-mcpserver/cmd/openapi-to-mcp@latest
```

## Usage

```bash
openapi-to-mcp --input path/to/openapi.json --output path/to/mcp-server.yaml
```

### Options

- `--input`: Path to the OpenAPI specification file (JSON or YAML) (required)
- `--output`: Path to the output MCP configuration file (YAML) (required)
- `--server-name`: Name of the MCP server (default: "openapi-server")
- `--tool-prefix`: Prefix for tool names (default: "")
- `--format`: Output format (yaml or json) (default: "yaml")
- `--validate`: Validate the OpenAPI specification (default: false)
- `--template`: Path to a template file to patch the output (default: "")

## Example

```bash
openapi-to-mcp --input petstore.json --output petstore-mcp.yaml --server-name petstore
```

### Converting OpenAPI to Higress REST-to-MCP Configuration

This tool can be used to convert an OpenAPI specification to a Higress REST-to-MCP configuration. Here's a complete example:

1. Start with an OpenAPI specification (petstore.json):

```json
{
  "openapi": "3.0.0",
  "info": {
    "version": "1.0.0",
    "title": "Swagger Petstore",
    "description": "A sample API that uses a petstore as an example to demonstrate features in the OpenAPI 3.0 specification"
  },
  "servers": [
    {
      "url": "http://petstore.swagger.io/v1"
    }
  ],
  "paths": {
    "/pets": {
      "get": {
        "summary": "List all pets",
        "operationId": "listPets",
        "parameters": [
          {
            "name": "limit",
            "in": "query",
            "description": "How many items to return at one time (max 100)",
            "required": false,
            "schema": {
              "type": "integer",
              "format": "int32"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "A paged array of pets",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "pets": {
                      "type": "array",
                      "items": {
                        "type": "object",
                        "properties": {
                          "id": {
                            "type": "integer",
                            "description": "Unique identifier for the pet"
                          },
                          "name": {
                            "type": "string",
                            "description": "Name of the pet"
                          },
                          "tag": {
                            "type": "string",
                            "description": "Tag of the pet"
                          }
                        }
                      }
                    },
                    "nextPage": {
                      "type": "string",
                      "description": "URL to get the next page of pets"
                    }
                  }
                }
              }
            }
          }
        }
      },
      "post": {
        "summary": "Create a pet",
        "operationId": "createPets",
        "requestBody": {
          "description": "Pet to add to the store",
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["name"],
                "properties": {
                  "name": {
                    "type": "string",
                    "description": "Name of the pet"
                  },
                  "tag": {
                    "type": "string",
                    "description": "Tag of the pet"
                  }
                }
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Null response"
          }
        }
      }
    },
    "/pets/{petId}": {
      "get": {
        "summary": "Info for a specific pet",
        "operationId": "showPetById",
        "parameters": [
          {
            "name": "petId",
            "in": "path",
            "required": true,
            "description": "The id of the pet to retrieve",
            "schema": {
              "type": "string"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Expected response to a valid request",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "id": {
                      "type": "integer",
                      "description": "Unique identifier for the pet"
                    },
                    "name": {
                      "type": "string",
                      "description": "Name of the pet"
                    },
                    "tag": {
                      "type": "string",
                      "description": "Tag of the pet"
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
```

2. Convert it to a Higress REST-to-MCP configuration:

```bash
openapi-to-mcp --input petstore.json --output petstore-mcp.yaml --server-name petstore
```

3. The resulting petstore-mcp.yaml file:

```yaml
server:
  name: petstore
tools:
  - name: showPetById
    description: Info for a specific pet
    args:
      - name: petId
        description: The id of the pet to retrieve
        type: string
        required: true
        position: path
    requestTemplate:
      url: /pets/{petId}
      method: GET
    responseTemplate:
      prependBody: |
        # API Response Information

        Below is the response from an API call. To help you understand the data, I've provided:

        1. A detailed description of all fields in the response structure
        2. The complete API response

        ## Response Structure

        > Content-Type: application/json

        - **id**: Unique identifier for the pet (Type: integer)
        - **name**: Name of the pet (Type: string)
        - **tag**: Tag of the pet (Type: string)

        ## Original Response

  - name: createPets
    description: Create a pet
    args:
      - name: name
        description: Name of the pet
        type: string
        required: true
        position: body
      - name: tag
        description: Tag of the pet
        type: string
        position: body
    requestTemplate:
      url: /pets
      method: POST
      headers:
        - key: Content-Type
          value: application/json
    responseTemplate: {}

  - name: listPets
    description: List all pets
    args:
      - name: limit
        description: How many items to return at one time (max 100)
        type: integer
        position: query
    requestTemplate:
      url: /pets
      method: GET
    responseTemplate:
      prependBody: |
        # API Response Information

        Below is the response from an API call. To help you understand the data, I've provided:

        1. A detailed description of all fields in the response structure
        2. The complete API response

        ## Response Structure

        > Content-Type: application/json

        - **pets**:  (Type: array)
          - **pets[].id**: Unique identifier for the pet (Type: integer)
          - **pets[].name**: Name of the pet (Type: string)
          - **pets[].tag**: Tag of the pet (Type: string)
        - **nextPage**: URL to get the next page of pets (Type: string)

        ## Original Response
```

4. This configuration can be used with Higress by adding it to your Higress gateway configuration.

Note how the tool automatically sets the `position` field for each parameter based on its location in the OpenAPI specification:

- The `petId` parameter is set to `position: path` because it's defined as `in: path` in the OpenAPI spec
- The `limit` parameter is set to `position: query` because it's defined as `in: query` in the OpenAPI spec
- The request body properties (`name` and `tag`) are set to `position: body`

The MCP server will automatically handle these parameters in the correct location when making API requests.

For more information about using this configuration with Higress REST-to-MCP, please refer to the [Higress REST-to-MCP documentation](https://higress.cn/en/ai/mcp-quick-start/#configuring-rest-api-mcp-server).

## Features

- Converts OpenAPI paths to MCP tools
- Supports both JSON and YAML OpenAPI specifications
- Generates MCP configuration with server and tool definitions
- Preserves parameter descriptions and types
- Automatically sets parameter positions based on OpenAPI parameter locations
- Handles path, query, header, cookie, and body parameters
- Generates response templates with field descriptions and improved formatting for LLM understanding
- **NEW**: Generates output schemas for MCP Protocol 2025-06-18 support
- **NEW**: Supports structured content responses for enhanced MCP compatibility
- Optional validation of OpenAPI specifications (disabled by default)
- Supports template-based patching of the generated configuration

## Template-Based Patching

You can use the `--template` flag to provide a YAML file that will be used to patch the generated configuration. This is useful for adding common headers, authentication, or other customizations to all tools in the configuration.

Example template file:

```yaml
server:
  name: "openapi-server"
  config:
    apiKey: ""

tools:
  outputSchema:
    type: object
    description: "Default output schema for all tools"
    properties:
      success:
        type: boolean
        description: "Indicates if the operation was successful"
      message:
        type: string
        description: "Response message"
  requestTemplate:
    headers:
      - key: Authorization
        value: "APPCODE {{.config.apiKey}}"
      - key: X-Ca-Nonce
        value: "{{uuidv4}}"
```

When applied, this template will:

1. Add an `apiKey` field to the server config
2. Add a default output schema to all tools (enabling MCP 2025-06-18 support)
3. Add the specified headers to all tools in the configuration

Usage:

```bash
openapi-to-mcp --input api-spec.json --output mcp-server.yaml --server-name my-server --template template.yaml
```

The template values like `{{.config.apiKey}}` or `"{{uuidv4}}"` are not processed by the tool but are preserved in the output for use by the MCP server at runtime.

## Security Scheme Conversion

The tool now supports the conversion of security schemes defined in your OpenAPI specification.

### Server-Level Security Schemes

Security schemes defined in the `components.securitySchemes` section of your OpenAPI document are converted into a list under `server.securitySchemes` in the generated MCP configuration.

**Example OpenAPI Snippet (`components.securitySchemes`):**

```json
{
  "components": {
    "securitySchemes": {
      "BasicAuth": {
        "type": "http",
        "scheme": "basic"
      },
      "ApiKeyAuth": {
        "type": "apiKey",
        "in": "header",
        "name": "X-API-KEY"
      }
    }
  }
}
```

**Corresponding MCP YAML Output (`server.securitySchemes`):**

```yaml
server:
  name: your-server-name
  securitySchemes:
    - id: ApiKeyAuth # Note: Schemes are sorted by ID in the output
      type: apiKey
      in: header
      name: X-API-KEY
    - id: BasicAuth
      type: http
      scheme: basic
  # ... other server config ...
```

The `defaultCredential` field within a security scheme is an MCP-specific extension and is not derived from the OpenAPI specification. You can set it using the `--template` feature if needed.

### Tool-Level Security Requirements

Security requirements defined at the operation level in your OpenAPI document (using the `security` keyword) are converted into a list under `requestTemplate.security` for the corresponding tool. Each entry in this list will reference the `id` of a security scheme defined in `server.securitySchemes`.

**Example OpenAPI Snippet (Operation with security):**

```json
{
  "paths": {
    "/protected_resource": {
      "get": {
        "summary": "Access a protected resource",
        "operationId": "getProtectedResource",
        "security": [
          {
            "ApiKeyAuth": []
          }
        ],
        "responses": {
          "200": { "description": "Success" }
        }
      }
    }
  }
}
```

**Corresponding MCP YAML Output (Tool's `requestTemplate.security`):**

```yaml
tools:
  - name: getProtectedResource
    description: Access a protected resource
    # ... args ...
    requestTemplate:
      url: /protected_resource # Actual URL depends on your server config in OpenAPI
      method: GET
      security:
        id: ApiKeyAuth
    # ... responseTemplate ...
```

### Template Overrides for Security

You can use the `--template` option to:

- Add new security schemes to `server.securitySchemes`.
- Override existing security schemes (e.g., to add `defaultCredential`).
- Override or set `security` requirements for all tools via the `tools.requestTemplate.security` path in your template file.

If the template defines `server.securitySchemes` or `tools.requestTemplate.security`, these will replace any schemes/requirements derived from the OpenAPI specification.

## MCP Protocol 2025-06-18 Support

This tool now supports the latest MCP Protocol version 2025-06-18, which introduces two key features:

### 1. Output Schema

Tools can now define schemas for their structured outputs. The tool automatically generates output schemas based on your OpenAPI response definitions.

**Example generated output schema:**

```yaml
tools:
  - name: getUser
    description: Get user information
    # ... other fields ...
    outputSchema:
      type: object
      description: User information
      properties:
        id:
          type: integer
          description: User ID
        name:
          type: string
          description: User name
        email:
          type: string
          description: User email
      required:
        - id
        - name
```

### 2. Structured Content

For MCP 2025-06-18 clients, responses include both traditional text content and structured JSON data. This allows MCP clients to process responses in a more structured way.

**MCP 2025-06-18 is automatically enabled when output schemas are present.**

### Backward Compatibility

The tool maintains full backward compatibility:

- Older MCP clients continue to work without changes
- New features are automatically enabled when output schemas are generated
- Existing functionality remains unchanged

### Using with Higress

When using the generated configuration with Higress:

1. The tool automatically generates output schemas for all endpoints with response definitions
2. MCP clients that support 2025-06-18 will automatically use the structured content
3. MCP clients that don't support 2025-06-18 will continue to work with the traditional text responses
4. No additional configuration is required - the presence of output schemas enables the new features

For more information about MCP Protocol 2025-06-18, refer to the [MCP specification](https://modelcontextprotocol.io/specification).

## Contributing

We welcome contributions from the community! Please see our [Contributing Guide](CONTRIBUTING.md) for information on how to get started.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
