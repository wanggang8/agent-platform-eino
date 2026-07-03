import type { DisplayField, JsonValue, StructuredItem, StructuredResult } from "../../../contracts/generated";

type StructuredResultViewProps = {
  readonly result: StructuredResult;
};

// StructuredResultView 是工具结果的唯一展示入口，避免前端直接解释 raw provider payload。
export function StructuredResultView({ result }: StructuredResultViewProps) {
  if (result.schema_version === "tool.structured_result.v1") {
    return (
      <div className="structured-result" data-testid="structured-result-view">
        <p>{result.data.summary ?? "结构化结果已生成。"}</p>
        <FactGrid facts={result.data.facts ?? []} />
      </div>
    );
  }

  return (
    <div className="structured-result" data-testid="structured-result-view">
      <div className="structured-heading">
        <strong>{result.data.title ?? displayTypeLabel(result.display_type)}</strong>
        <span>{result.data.summary ?? statusLabel(result.status)}</span>
      </div>
      {result.data.metrics && result.data.metrics.length > 0 ? (
        <div className="metric-strip">
          {result.data.metrics.map((metric) => (
            <div className={`metric ${metric.tone ?? "neutral"}`} key={metric.label}>
              <span>{metric.label}</span>
              <strong>{formatDisplayValue(metric.value)}</strong>
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

// ResultTable 渲染 fobrain 结构化集合结果，列定义来自 contract。
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
          {columns.map((column) => <span role="cell" key={column.key}>{formatDisplayValue(readValue(item, column.key), column.type)}</span>)}
        </div>
      ))}
    </div>
  );
}

// FactGrid 展示安全事实键值对。
function FactGrid({ facts }: { readonly facts: readonly DisplayField[] }) {
  if (facts.length === 0) return null;
  return (
    <dl className="fact-grid">
      {facts.map((fact) => (
        <div key={fact.label}>
          <dt>{fact.label}</dt>
          <dd>{formatDisplayValue(fact.value, fact.type)}</dd>
        </div>
      ))}
    </dl>
  );
}

// readValue 只从 StructuredItem 安全字段中取值。
function readValue(item: StructuredItem, key: string): JsonValue {
  const value = (item as Readonly<Record<string, JsonValue | undefined>>)[key];
  return value ?? "";
}

// formatDisplayValue 只输出产品化摘要，禁止把嵌套对象作为 JSON 字符串展示给用户。
function formatDisplayValue(value: JsonValue, fieldType?: DisplayField["type"]): string {
  if (value === null) return "-";
  if (Array.isArray(value)) return value.map((item) => formatDisplayValue(item, fieldType)).join("、");
  if (typeof value === "object") return "已整理";
  if (typeof value === "string") return normalizeDisplayCode(value, fieldType);
  return String(value);
}

// normalizeDisplayCode 将 provider/status 安全枚举映射为中文展示，避免产品界面出现原始英文状态码。
function normalizeDisplayCode(value: string, fieldType?: DisplayField["type"]): string {
  const normalized = value.toLowerCase();
  const labels: Record<string, string> = {
    resolved: "已解析",
    succeeded: "已完成",
    completed: "已完成",
    available: "可用",
    unavailable: "不可用",
    bound: "已绑定",
    unbound: "未绑定",
    configured: "已配置",
    missing: "缺少配置",
    workspace: "当前工作区",
    "ws-demo": "演示工作区",
    internal: "内网",
    external: "外网",
    high: "高危",
    medium: "中危",
    low: "低危",
    critical: "严重"
  };
  if (labels[normalized]) return labels[normalized];
  if (fieldType === "status" || fieldType === "severity" || fieldType === "badge" || fieldType === "code") return value.replace(/[A-Za-z_:-]+/g, "已整理");
  return value;
}

// displayTypeLabel 将结构化展示类型映射为中文标签。
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

// statusLabel 将结构化状态映射为中文标签。
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
  return labels[status] ?? "未知";
}
