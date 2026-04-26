# AI Changes

## 2026-04-15 16:07 CST
- 修复 `nullable oneOf + $ref` 组合 schema 的类型推断，避免 `GenerateSql` / `CreateReport` 等工具生成空类型或空对象 schema。
- 涉及文件：`pkg/converter/schema_helpers.go`、`pkg/converter/converter.go`、`pkg/converter/converter_test.go`。
- `unwrapSingleComposedSchema` 现在支持从 `null + 一个真实 schema` 的组合中提取有效 schema，并保留原有类型推断路径。
- 新增回归测试覆盖 `GenerateSql` 风格的 request/response nullable composed refs，校验 `aggregation`、`promql_condition`、`search_event_type`、`metadata` 均能推断出合法类型。
- 兼容性修正：保留外层 `description`，并继续维持标量数组项的空描述输出，避免既有 golden 文件回归。
- 全量测试通过，重新生成 OpenObserve MCP 产物后 `type: ""` 数量为 0。

## 2026-04-15 16:04 CST
- 修复请求参数与 requestBody 嵌套 schema 的类型推断，避免生成 `type: ""` 的非法 MCP schema。
- 涉及文件：`pkg/converter/converter.go`、`pkg/converter/schema_helpers.go`、`pkg/converter/converter_test.go`。
- 统一 `convertParameters` / `convertRequestBody` 使用 `effectiveSchemaType` 推断对象、数组和 enum-only 字段类型。
- 数组 items 仅在类型非空时输出 `type`，并保留兼容性所需的空 `description` 字段。
- 新增回归测试覆盖 `createPipeline` 风格的嵌套 request body，验证不会再生成空类型。
- 使用真实 OpenObserve OpenAPI 文档复验，`createPipeline.nodes.items.type=object`，生成 YAML 中 `type: ""` 数量为 0。
