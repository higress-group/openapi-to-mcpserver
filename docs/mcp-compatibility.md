# MCP Compatibility Guide

MCP (Model Context Protocol) tools communicate via JSON-RPC. All tool inputs must be JSON-serializable, and tool calls follow a synchronous request-response model. This means certain OpenAPI operations cannot be directly converted to MCP tools.

The `--mcp-compat` flag detects these incompatibilities before conversion.

## Detected Patterns

| Incompatible Pattern | Reason |
|---|---|
| `application/octet-stream` request body | Binary stream cannot be JSON-serialized |
| `multipart/form-data` request body | Multipart encoding has no JSON equivalent |
| `text/event-stream` request body | Streaming input not supported |
| `application/xml` / `text/xml` request body | XML cannot be passed as JSON tool arguments |
| `image/*`, `audio/*`, `video/*` request body | Media types cannot be JSON-serialized as input |
| `format: "binary"` fields in request body | Binary data cannot be transmitted as MCP tool input |
| `text/event-stream` in response | MCP tool calls are synchronous; streaming responses are incompatible |

> **Note**: MCP tool *results* (outputs) support image and audio via base64 encoding. The check only flags issues on the **input** side, except for streaming responses which fundamentally conflict with MCP's synchronous model.

## Modes

### Strict Mode (`--mcp-compat strict`)

Prints all incompatible operations and exits with a non-zero status code. No output file is generated.

```
$ openapi-to-mcp --input api.json --output mcp.yaml --mcp-compat strict

Error: OpenAPI specification contains MCP-incompatible operations:
  [MCP incompatible] uploadFile: request body content type "application/octet-stream" is not supported by MCP (only JSON-serializable inputs are allowed)
  [MCP incompatible] streamEvents: response (200) uses streaming content type "text/event-stream", MCP tool calls are synchronous request-response
  [MCP incompatible] uploadAvatar: request body content type "multipart/form-data" is not supported by MCP (multipart form data cannot be represented as JSON tool arguments)
```

**When to use**: CI pipelines, pre-commit checks, or when you want to ensure the entire spec is MCP-compatible.

### Warn Mode (`--mcp-compat warn`)

Prints warnings for incompatible operations, skips them, and converts the remaining compatible operations.

```
$ openapi-to-mcp --input api.json --output mcp.yaml --mcp-compat warn

Warning: [MCP incompatible] uploadFile: request body content type "application/octet-stream" is not supported by MCP (only JSON-serializable inputs are allowed)
Warning: [MCP incompatible] streamEvents: response (200) uses streaming content type "text/event-stream", MCP tool calls are synchronous request-response
Warning: [MCP incompatible] uploadAvatar: request body content type "multipart/form-data" is not supported by MCP (multipart form data cannot be represented as JSON tool arguments)

Successfully converted OpenAPI specification to MCP configuration: mcp.yaml
```

The output file will only contain tools from compatible operations.

**When to use**: When your API has a mix of compatible and incompatible endpoints, and you want to convert what's possible.

## How to Fix Incompatible Operations

If you need the flagged operations available as MCP tools:

1. **File uploads** (`multipart/form-data`, `application/octet-stream`): Redesign the API to accept a URL or base64-encoded string in a JSON body instead of raw binary.
2. **Streaming responses** (`text/event-stream`): Provide an alternative polling endpoint that returns the current state as a JSON object.
3. **XML request bodies**: Add an `application/json` content type alternative to the operation.

## Default Behavior

When `--mcp-compat` is not specified, no compatibility check is performed and all operations are converted regardless of compatibility. This preserves backward compatibility with existing workflows.
