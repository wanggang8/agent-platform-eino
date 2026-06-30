# Workbench Visual Acceptance Task 15

## Scope

- Canonical target: `docs/design/assets/workbench-target-ui-v1-2026-06-26.png`.
- Real region screenshots: 150 from `test-results/workbench-real-model-tool-1782607638250-*`.
- Real full-shell screenshots: 25 captured under `test-results/workbench-ui-rewrite/15-visual-acceptance/full-shell/`.
- Functional/data source: `test-results/workbench-real-model-tool-report.json`.

## Data Correctness Gate

- Scenarios: `25`.
- Real report passed/failed: `25/0`.
- Screenshots in report: `150`.
- Full-shell screenshots captured: `25`.
- All scenario verdicts passed: `True`.
- All final states completed: `True`.
- All DOM/fresh/tool/answer/auxiliary/schema/leak gates passed: `True`.

## Visual Metric Summary

| Region | Count | Avg color distance | Max color distance | Avg density | Avg contrast | Statuses |
| --- | ---: | ---: | ---: | ---: | ---: | --- |
| `audit` | 25 | 5.97 | 6.4 | 0.0416 | 22.29 | `{"metric_ok": 25}` |
| `evidence` | 25 | 6.27 | 26.65 | 0.0808 | 26.98 | `{"metric_ok": 25}` |
| `fresh-main-chat` | 25 | 2.85 | 18.6 | 0.0669 | 23.65 | `{"metric_ok": 25}` |
| `full-shell` | 25 | 5.4 | 10.49 | 0.0812 | 24.92 | `{"metric_ok": 25}` |
| `internal-details` | 25 | 8.16 | 8.77 | 0.0307 | 15.92 | `{"metric_ok": 25}` |
| `main-chat` | 25 | 2.85 | 18.6 | 0.0669 | 23.65 | `{"metric_ok": 25}` |
| `process` | 25 | 7.76 | 8.12 | 0.0349 | 18.64 | `{"metric_ok": 25}` |

## Metric-Flagged Items For Human Review

These are not automatic failures. They identify screenshots whose aggregate color/density/contrast differ enough from the target crop to require visual inspection.

| Scenario | Region | Status | Color distance | Density | Contrast | Path |
| --- | --- | --- | ---: | ---: | ---: | --- |
| - | - | - | - | - | - | - |

## Contact Sheets

- `test-results/workbench-ui-rewrite/15-visual-acceptance/region-contact-sheets/audit.png`
- `test-results/workbench-ui-rewrite/15-visual-acceptance/region-contact-sheets/evidence.png`
- `test-results/workbench-ui-rewrite/15-visual-acceptance/region-contact-sheets/fresh-main-chat.png`
- `test-results/workbench-ui-rewrite/15-visual-acceptance/region-contact-sheets/full-shell.png`
- `test-results/workbench-ui-rewrite/15-visual-acceptance/region-contact-sheets/internal-details.png`
- `test-results/workbench-ui-rewrite/15-visual-acceptance/region-contact-sheets/main-chat.png`
- `test-results/workbench-ui-rewrite/15-visual-acceptance/region-contact-sheets/process.png`

## Required Manual Review Checklist

- Left sidebar: brand, new chat, search/filter, conversation rows, selected row border, footer settings.
- Header: breadcrumb/title, chips, action buttons, spacing, no wrapping.
- Main conversation: user bubble, assistant text, single tool card, multi-tool/workflow cards, status icons, final answer actions.
- Inspector: tabs, evidence table, detail panels, lower panels, audit/runtime/structured-result tabs.
- Composer: mode controls, input placeholder, tool/capability/attachment buttons, send button.
- Typography/color/border/spacing/texture: compare against the canonical target and contact sheets.
- Data correctness: verify visible claims match `structured_result` and the real report gates.

## Current Conclusion

Task 15 has generated the full evidence package and metric triage. Human visual review of the generated contact sheets is still required before marking this visual acceptance goal complete.
