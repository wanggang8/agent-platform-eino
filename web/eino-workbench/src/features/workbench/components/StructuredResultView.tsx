import type { JsonValue, StructuredItem, StructuredResult } from "../../../contracts/generated";

type StructuredResultViewProps = {
  readonly result: StructuredResult;
};

export function StructuredResultView({ result }: StructuredResultViewProps) {
  if (result.schema_version === "tool.structured_result.v1") {
    return (
      <div className="structured-result">
        <p>{result.data.summary ?? "结构化结果已生成。"}</p>
        <FactGrid facts={result.data.facts ?? []} />
      </div>
    );
  }

  return (
    <div className="structured-result">
      <div className="structured-heading">
        <strong>{result.data.title ?? displayTypeLabel(result.display_type)}</strong>
        <span>{result.data.summary ?? statusLabel(result.status)}</span>
      </div>
      {result.data.metrics && result.data.metrics.length > 0 ? (
        <div className="metric-strip">
          {result.data.metrics.map((metric) => (
            <div className={`metric ${metric.tone ?? "neutral"}`} key={metric.label}>
              <span>{metric.label}</span>
              <strong>{formatValue(metric.value)}</strong>
              {metric.sub_label ? <small>{metric.sub_label}</small> : null}
            </div>
          ))}
        </div>
      ) : null}
      {result.data.columns && result.data.items ? <ResultTable columns={result.data.columns} items={result.data.items} /> : null}
      {result.data.facts ? <FactGrid facts={result.data.facts} /> : null}
    </div>
  );
}

function ResultTable({
  columns,
  items
}: {
  readonly columns: NonNullable<Extract<StructuredResult, { schema_version: "fobrain.tool_result.v2" }>["data"]["columns"]>;
  readonly items: readonly StructuredItem[];
}) {
  return (
    <div className="result-table" role="table">
      <div className="result-row is-header" role="row">
        {columns.map((column) => <span role="columnheader" key={column.key}>{column.label}</span>)}
      </div>
      {items.map((item, index) => (
        <div className="result-row" role="row" key={item.entity_ref ?? item.row_ref ?? String(index)}>
          {columns.map((column) => <span role="cell" key={column.key}>{formatValue(readValue(item, column.key))}</span>)}
        </div>
      ))}
    </div>
  );
}

function FactGrid({ facts }: { readonly facts: readonly { readonly label: string; readonly value: JsonValue }[] }) {
  if (facts.length === 0) return null;
  return (
    <dl className="fact-grid">
      {facts.map((fact) => (
        <div key={fact.label}>
          <dt>{fact.label}</dt>
          <dd>{formatValue(fact.value)}</dd>
        </div>
      ))}
    </dl>
  );
}

function readValue(item: StructuredItem, key: string): JsonValue {
  const value = (item as Readonly<Record<string, JsonValue | undefined>>)[key];
  return value ?? "";
}

function formatValue(value: JsonValue): string {
  if (value === null) return "-";
  if (Array.isArray(value)) return value.map(formatValue).join(", ");
  if (typeof value === "object") return JSON.stringify(value);
  return String(value);
}

function displayTypeLabel(displayType: string) {
  const labels: Record<string, string> = {
    entity_collection: "对象列表",
    entity_detail: "对象详情",
    metrics_summary: "指标摘要",
    operation_result: "操作结果",
    connector_status: "连接器状态",
    entity_resolution: "实体消歧"
  };
  return labels[displayType] ?? "结构化结果";
}

function statusLabel(status: string) {
  const labels: Record<string, string> = {
    resolved: "已解析",
    waiting: "等待中",
    pending_approval: "等待审批",
    empty: "无结果",
    failed: "失败",
    partial: "部分结果",
    not_found: "未找到"
  };
  return labels[status] ?? status;
}
