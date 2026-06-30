# Workbench Full Scenario Acceptance

## Scope

- Replays every real conversation from the source real-model/tool report.
- Captures full shell, sidebar, header, main chat, composer, all inspector tabs, evidence row detail, and narrative-link result.
- Validates visible text safety, old UI removal, tab presence, narrative linkage, row detail interaction, terminal tool states, overflow candidates, and implementation boundary checks.

## Summary

- Source report: `test-results/workbench-ui-rewrite/16-product-narrative/workbench-real-model-tool-report-1782607638250.json`
- Base URL: `http://127.0.0.1:8080`
- Scenarios: `25`
- Passed / failed: `25/0`
- Screenshots: `275`
- Static boundary verdict: `passed`

## Scenario Matrix

| Scenario | Screenshots | Verdict | Failed checks |
| --- | ---: | --- | --- |
| `current_user_context` | 11 | passed | - |
| `my_permissions` | 11 | passed | - |
| `assets_by_owner` | 11 | passed | - |
| `vulnerabilities_by_owner` | 11 | passed | - |
| `department_assets` | 11 | passed | - |
| `department_vulnerabilities` | 11 | passed | - |
| `ip_assets` | 11 | passed | - |
| `ip_vulnerabilities` | 11 | passed | - |
| `asset_detail` | 11 | passed | - |
| `vulnerability_detail` | 11 | passed | - |
| `business_risk_summary` | 11 | passed | - |
| `business_list` | 11 | passed | - |
| `external_high_risk_assets` | 11 | passed | - |
| `threat_relevance_list` | 11 | passed | - |
| `vulnerability_status_summary` | 11 | passed | - |
| `pending_tickets` | 11 | passed | - |
| `ip_stats` | 11 | passed | - |
| `vul_stats` | 11 | passed | - |
| `my_assets` | 11 | passed | - |
| `my_department_assets` | 11 | passed | - |
| `my_vulnerabilities` | 11 | passed | - |
| `my_department_vulnerabilities` | 11 | passed | - |
| `my_business_systems` | 11 | passed | - |
| `my_important_business_systems` | 11 | passed | - |
| `fobrain_connector_status` | 11 | passed | - |

## Screenshot Region Coverage

| Region | Count |
| --- | ---: |
| `full-shell` | 25 |
| `sidebar` | 25 |
| `header` | 25 |
| `main-chat` | 25 |
| `composer` | 25 |
| `narrative-link-result` | 25 |
| `inspector-evidence` | 25 |
| `inspector-structured` | 25 |
| `inspector-runtime` | 25 |
| `inspector-audit` | 25 |
| `evidence-row-detail` | 25 |

## Static Boundary Findings

- None.

## Artifacts

- JSON report: `test-results/workbench-ui-rewrite/17-full-scenario-acceptance/full-scenario-acceptance.json`
- Screenshots: `test-results/workbench-ui-rewrite/17-full-scenario-acceptance/screenshots/`
