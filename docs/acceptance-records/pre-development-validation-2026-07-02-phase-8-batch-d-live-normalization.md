# Pre-development Validation: Phase 8 Batch D Live Normalization

日期：2026-07-02

范围：Fobrain Batch D 参数化查询接入真实 HTTP client，并以当前真实接口返回重新归一化安全字段。

## 校验结论

| Item | Evidence | Result |
| --- | --- | --- |
| Phase 定位 | `docs/fobrain-live-read-batch-plan.md` | 属于 Phase 8 Batch D live read mapper 切片，不声明 24 个只读工具完成。 |
| 用户决策 | 用户确认“以当前真实接口返回为准重新归一化” | 字段映射以 live 探测为事实源，旧项目只用于定位候选路径。 |
| live 配置 | `configs/eino-workbench.local.yaml` 脱敏检查 | Fobrain live、workspace token、`authorization` 参数和私有证书跳过均已配置；未记录 token。 |
| 旧项目约束 | `/Users/vick/Desktop/project/ai-agent` 只读参考 | 未修改旧项目，未复制旧 runtime、旧执行链路或 Workbench 接口。 |
| 安全边界 | `docs/facts-contract.md`、`docs/06-security-and-projection.md` | raw provider payload 只在 provider 边界内解析，Product Facts 仍只接收 StructuredResult。 |

## 开工边界

- 允许：实现 `HTTPClient` 的 Batch D `ParameterizedQueryClient`、当前真实路径优先、旧 `/api/v1` 路径兼容 fallback、测试和文档。
- 不允许：把 raw payload、token、credential ref 写入 Product Facts、报告、audit 或 replay。
- 后续：完整 live smoke 报告、Workbench 视觉证据和 24/24 最终恢复需独立门禁。
