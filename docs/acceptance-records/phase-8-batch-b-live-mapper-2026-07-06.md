# Phase 8 Batch B HTTP Live Mapper

日期：2026-07-06  
方案：Batch B “我的范围”六个只读工具 HTTP live mapper focused tests  
状态：Passed for local mapper tests

## Scope

- `tool.fobrain.my_assets`
- `tool.fobrain.my_department_assets`
- `tool.fobrain.my_vulnerabilities`
- `tool.fobrain.my_department_vulnerabilities`
- `tool.fobrain.my_business_systems`
- `tool.fobrain.my_important_business_systems`

## Verification

```bash
go test ./internal/einoapp/providers/fobrain -run 'HTTPClientMyScope' -count=1
go test ./internal/einoapp/providers/fobrain -count=1
```

## Evidence

- HTTP client 已实现 `MyScopeClient`，Batch B 进入真实 provider 可选能力边界。
- 资产/漏洞“我的范围”从 `/api/v1/user` 当前用户派生 `oper_info.name`、`person_info.name`。
- 本部门资产/漏洞从当前用户部门派生 `business_department.name.keyword`、`person_department.name.keyword`。
- 资产、漏洞、业务系统均按当前 `/api/...` 优先，404/405 时 fallback 到 `/api/v1/...`。
- 业务系统 owner scope 使用旧源码确认的 `person_base.name` 条件；用户输入 keyword 只作为 `business_name` 额外收窄条件，避免放宽“我的范围”。
- 重要业务系统按 Fobrain 源码字段 `assets_attribute.important_types` 通过 `search_condition` 过滤 1/2 两档重要性。
- 业务错误、认证失败和 raw provider payload 不进入 Product Facts 或调用方错误文本。

## Explicit Non-Claims

- 不声明 Batch B 真实环境 live pass。
- 不声明 Batch B 分页批量验收完成。
- 不声明 Batch B Workbench 视觉通过。
- 不声明 Batch B Action API 同源 smoke 通过。
- 不声明 Fobrain 24/24 只读恢复完成。

## Next Gate

下一步应补 `fobrain-batch-b` live smoke/report schema 和 server smoke 场景，使用真实配置验证六个工具的分页、空态和脱敏报告。
