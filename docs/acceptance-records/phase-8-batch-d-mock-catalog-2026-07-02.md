# Phase 8 Batch D Mock Catalog Acceptance Record

日期：2026-07-02

范围：Fobrain Batch D 六个参数化只读工具的 catalog、输入校验和 mock StructuredResult mapper。

## 验收矩阵

| Criterion | Source | Evidence | Verification method | Status | Notes/Risks |
| --- | --- | --- | --- | --- | --- |
| Batch D 六个工具进入 mock/参数化 provider catalog | `docs/fobrain-live-read-batch-plan.md`、`tool-matrix-24.json` | `internal/einoapp/providers/fobrain/batch_d_test.go` | `go test ./internal/einoapp/providers/fobrain -count=1` | Pass | 只有 client 实现 `ParameterizedQueryClient` 时才注册 Batch D；真实 HTTP client 不提前声明不可执行工具。 |
| mock 运行期 Batch D 有 policy context 且可调用 | `docs/capability-provider-contract.md` | `cmd/eino-workbench/main_test.go` | `go test ./cmd/eino-workbench -count=1` | Pass | policy context 从配置化 provider catalog 派生，不维护第二份工具清单。 |
| live HTTP 运行期不暴露 Batch D | `docs/fobrain-live-read-batch-plan.md` | `cmd/eino-workbench/main_test.go` | `go test ./cmd/eino-workbench -count=1` | Pass | HTTP mapper 未实现前，Batch D 不进入 live registry。 |
| owner / department / ip 必填字段与矩阵一致 | `docs/schemas/fobrain/tool_inputs.v1.schema.json` | Batch D catalog test、schema validation | `go test ./internal/einoapp/providers/fobrain -run BatchD -count=1`、`npm run eino-workbench:schema-test` | Pass | Go catalog 只覆盖轻量 `capabilities.JSONSchema` 字段映射；`additionalProperties` 和 min/max 仍由 docs schema 记录并由 schema-test 校验 schema 文件有效。 |
| Batch D 调用输出安全 StructuredResult candidate | `docs/facts-contract.md` | provider invoke test | `go test ./internal/einoapp/providers/fobrain -run BatchD -count=1` | Pass | 不保存 fobrain business payload 到 Product Facts。 |
| Provider 回填 raw title/query 不进入事实摘要 | `docs/06-security-and-projection.md`、`docs/facts-contract.md` | safe override regression test | `go test ./internal/einoapp/providers/fobrain -run BatchD -count=1` | Pass | StructuredResult 只使用本地元数据和已校验请求 query。 |
| 缺少必填参数不调用 client | `docs/capability-provider-contract.md` | missing argument test | `go test ./internal/einoapp/providers/fobrain -run BatchD -count=1` | Pass | 返回安全 provider error。 |

## 验证命令

```bash
go test ./internal/einoapp/providers/fobrain -count=1
go test ./cmd/eino-workbench -count=1
npm run eino-workbench:schema-test
```

## 剩余风险

- Batch D 尚未接入真实 HTTP client、真实分页、实体消歧恢复数据和 Workbench 视觉证据；真实 HTTP client 当前不会暴露 Batch D catalog。
- 不能据此声明 Fobrain 24 个只读工具恢复完成。
