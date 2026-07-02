# Phase 8 Batch D Live Normalization Acceptance Record

日期：2026-07-02

范围：Fobrain Batch D 六个参数化只读工具的真实 HTTP client mapper 与当前接口字段归一化。

## Live 探测结论

- 当前真实用户接口：`/api/user` 可用，`/api/v1/user` 返回 404；Batch A fallback 已覆盖。
- 当前真实资产列表接口：`/api/asset`、`/api/internal_asset`、`/api/external_ip_asset` 可用，`/api/v1/asset` 返回 404。
- 当前真实漏洞列表接口：`/api/threat_center` 可用，`/api/v1/threat_center` 返回 404。
- 资产列表包装：`{code,data,message}`，`data` 包含 `items/page/per_page/total`。
- 漏洞列表包装：`{code,data,message}`，`data` 包含 `items/page/per_page/total`。
- 探测输出未记录 token、raw payload 或完整业务值。

## 验收矩阵

| Criterion | Source | Evidence | Verification method | Status | Notes/Risks |
| --- | --- | --- | --- | --- | --- |
| HTTP client 实现 Batch D 参数化只读接口 | `docs/fobrain-live-read-batch-plan.md` | `internal/einoapp/providers/fobrain/http_client_batch_d.go` | `go test ./internal/einoapp/providers/fobrain -count=1` | Pass | `HTTPClient` 现在实现 `ParameterizedQueryClient`。 |
| 当前真实路径优先，旧 `/api/v1` 兼容 fallback | live 探测、旧项目只读路径参考 | `TestHTTPClientParameterizedAssetByIPUsesCurrentAPIPath`、`TestHTTPClientParameterizedQueryFallsBackToLegacyV1Path` | `go test ./internal/einoapp/providers/fobrain -run HTTPClientParameterized -count=1` | Pass | 当前环境优先 `/api/asset` 与 `/api/threat_center`。 |
| 六个 Batch D 工具都有 HTTP 参数构造覆盖 | `docs/fixtures/fobrain/tool-matrix-24.json` | IP 两项端到端映射测试，owner/department 四项 search_condition 测试 | `go test ./internal/einoapp/providers/fobrain -count=1` | Pass | live 环境当前未提供稳定负责人/部门样本，单元测试固定参数契约。 |
| raw provider payload 不进入 Product Facts | `docs/facts-contract.md` | StructuredResult mapper 仍只使用本地元数据和已校验 query；HTTP mapper 输出安全 `QueryResultItem` | `go test ./internal/einoapp/providers/fobrain -count=1` | Pass | Product Facts 未保存业务表格 payload。 |
| Batch A live 未回归 | Batch A 验收门禁 | `test-results/eino-workbench-fobrain-batch-a-live-report.json` | `bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-batch-a --config configs/eino-workbench.local.yaml` | Pass | 私有证书配置继续生效。 |

## 验证命令

```bash
go test ./internal/einoapp/providers/fobrain -count=1
go test ./cmd/eino-workbench -count=1
go test ./...
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-batch-a --config configs/eino-workbench.local.yaml
```

## 剩余风险

- 尚未生成 Batch D 专用 live report schema 和 smoke 脚本；本记录引用的是脱敏 live 探测和 HTTP client 单元测试。
- 尚未补齐 Workbench Batch D 视觉证据。
- 尚不能声明 Fobrain 24 个只读工具恢复完成。
