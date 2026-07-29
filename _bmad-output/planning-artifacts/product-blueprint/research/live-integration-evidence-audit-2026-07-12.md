---
status: PASS
audit_type: fobrain-current-live-integration-evidence
date: '2026-07-12'
scope: technical-integration-only
---

# FOBrain 当前 live 集成证据审计

## 结论

当前仓库已有真实 `provider_mode=live` 报告，足以直接复用为对应 capability 的技术集成证据，不需要在产品研究阶段重新证明这些接口存在。

| 项目 | 结果 |
| --- | ---: |
| 历史定义的只读工具总数 | 24 |
| live `passed` | 17 |
| live `blocked` | 6 |
| live `skipped` | 1 |
| connector live passed | 1 |
| sample discovery | passed |

通过项可标记为 `EXPERIMENT_RESULT / DIRECT / PARTIAL / INTEGRATION_AVAILABLE / OPEN / NOT_NORMATIVE`。`PARTIAL` 表示报告证明当前环境和调用链，但尚未证明目标用户任务、产品范围或跨客户适用性。

## 报告结果

| 报告 | provider | 状态 | 直接结论 |
| --- | --- | --- | --- |
| `batch-a-live-report.json` | live | passed | 2 个只读工具 + connector 通过 |
| `batch-b-live-report.json` | live | blocked | 6 个“我的／本部门范围”能力被空结果或 current-user scope 阻塞 |
| `batch-c-live-report.json` | live | passed | 5 个只读工具通过，1 个工单能力 skipped |
| `batch-d-live-report.json` | live | passed | 6 个参数化资产／漏洞查询工具通过 |
| `batch-e-live-report.json` | live | passed | 4 个详情／汇总／关联工具通过 |
| `sample-discovery-report.json` | live | passed | 可安全发现样本，报告不暴露真实样本值 |

## 验证命令

```text
jq -s '[.[].capability_results[]] ...' test-results/eino-workbench-fobrain-batch-{a,b,c,d,e}-live-report.json
git diff --check
```

实际聚合结果：`24 total / 17 passed / 6 blocked / 1 skipped / 1 connector passed`。

## 它证明什么

- 当前环境中指定 capability 的 provider 调用真实发生；
- passed 项产生稳定 `result_ref` 和 StructuredResult 路径；
- 报告记录了凭据、Authorization、raw payload 和样本值等脱敏检查；
- blocked／skipped 是可用的反证，不应被隐藏。

## 它不证明什么

- 用户最重要的实际任务；
- 任务频率、耗时、返工、错误或失败损失；
- FOBrain 原生操作是否已经足够；
- AI 是否改善任务结果；
- 哪些工具应该进入产品范围；
- 产品应是内嵌、独立 Workbench、通用平台或不做；
- 用户采用、经济买方和购买意愿。

因此后续只做**任务研究和任务回放**，不重复做“工具是否能调用”的广泛调研。
