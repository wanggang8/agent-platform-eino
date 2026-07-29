package facts

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

const (
	// StructuredResultSchemaVersionV2 是 M1 唯一允许进入事实链的工具结果版本。
	StructuredResultSchemaVersionV2  = "tool.structured_result.v2"
	QueryResultSnapshotSchemaVersion = "eino_query_result_snapshot.v2"
	ProductFactsSchemaVersionV2      = "eino_product_facts.v2"
	NewVulnerabilityCapabilityID     = "capability.vulnerability.list_new"
	NewVulnerabilityQueryKind        = "new_vulnerabilities"
	CoverageCompleteSet              = "complete_set"
	FailedQueryMessage               = "无法读取漏洞事实。此次请求不是空结果。"
	ErrorCodeFixtureReadFailed       = "fixture_read_failed"
)

// QueryResultStatus 区分成功有结果和成功空结果；失败不创建 StructuredResult。
type QueryResultStatus string

const (
	QueryResultResolved QueryResultStatus = "resolved"
	QueryResultEmpty    QueryResultStatus = "empty"
)

// SnapshotItem 是 M1 冻结的最小安全行事实。
type SnapshotItem struct {
	ref          string
	displayLabel string
	discoveredAt time.Time
}

// NewSnapshotItem 在可信边界创建封闭行事实。
func NewSnapshotItem(ref, displayLabel string, discoveredAt time.Time) (SnapshotItem, error) {
	if !validOpaqueRef(ref, "item_") || strings.TrimSpace(displayLabel) == "" || discoveredAt.IsZero() {
		return SnapshotItem{}, errors.New("invalid snapshot item")
	}
	if ContainsUnsafeMaterial(displayLabel) || ContainsUnsafeMaterial(ref) {
		return SnapshotItem{}, ErrUnsafeFactMaterial
	}
	return SnapshotItem{ref: ref, displayLabel: displayLabel, discoveredAt: discoveredAt.UTC()}, nil
}

// MustSnapshotItem 仅供受版本控制的 fixture 与测试构造确定性行事实。
func MustSnapshotItem(ref, displayLabel string, discoveredAt time.Time) SnapshotItem {
	item, err := NewSnapshotItem(ref, displayLabel, discoveredAt)
	if err != nil {
		panic(err)
	}
	return item
}

func (item SnapshotItem) Ref() string             { return item.ref }
func (item SnapshotItem) DisplayLabel() string    { return item.displayLabel }
func (item SnapshotItem) DiscoveredAt() time.Time { return item.discoveredAt }

type structuredResultWire struct {
	SchemaVersion string                `json:"schema_version"`
	CapabilityID  string                `json:"capability_id"`
	Status        QueryResultStatus     `json:"status"`
	Query         structuredResultQuery `json:"query"`
	Data          structuredResultData  `json:"data"`
	Metadata      structuredResultMeta  `json:"metadata"`
}

type structuredResultQuery struct {
	Kind string `json:"kind"`
}

type structuredResultData struct {
	Summary string             `json:"summary"`
	Count   int                `json:"count"`
	Items   []snapshotItemWire `json:"items"`
}

type snapshotItemWire struct {
	SnapshotItemRef string `json:"snapshot_item_ref"`
	DisplayLabel    string `json:"display_label"`
	DiscoveredAt    string `json:"discovered_at"`
}

type structuredResultMeta struct {
	Safe       bool   `json:"safe"`
	Coverage   string `json:"coverage"`
	ObservedAt string `json:"observed_at"`
}

// StructuredResultInput 是 Safety Gate 批准后创建唯一事实材料的输入。
type StructuredResultInput struct {
	InternalResultRef string
	Status            QueryResultStatus
	Summary           string
	ObservedAt        time.Time
	Items             []SnapshotItem
}

// StructuredResult 按值持有 canonical safe JSON；内部引用永不进入该 JSON。
type StructuredResult struct {
	internalResultRef string
	status            QueryResultStatus
	summary           string
	observedAt        time.Time
	items             []SnapshotItem
	safeJSON          []byte
	digest            string
}

// NewStructuredResult 规范化排序并计算不可变内容摘要。
func NewStructuredResult(input StructuredResultInput) (StructuredResult, error) {
	if !validOpaqueRef(input.InternalResultRef, "result_ref:") || strings.TrimSpace(input.Summary) == "" || input.ObservedAt.IsZero() {
		return StructuredResult{}, errors.New("invalid structured result")
	}
	if ContainsUnsafeMaterial(input.InternalResultRef) || ContainsUnsafeMaterial(input.Summary) {
		return StructuredResult{}, ErrUnsafeFactMaterial
	}
	items := append([]SnapshotItem(nil), input.Items...)
	for _, item := range items {
		if _, err := NewSnapshotItem(item.ref, item.displayLabel, item.discoveredAt); err != nil {
			return StructuredResult{}, err
		}
	}
	sort.Slice(items, func(i, j int) bool {
		if !items[i].discoveredAt.Equal(items[j].discoveredAt) {
			return items[i].discoveredAt.After(items[j].discoveredAt)
		}
		return items[i].ref < items[j].ref
	})
	if (input.Status == QueryResultResolved && len(items) == 0) || (input.Status == QueryResultEmpty && len(items) != 0) {
		return StructuredResult{}, errors.New("structured result status/count mismatch")
	}
	if input.Status != QueryResultResolved && input.Status != QueryResultEmpty {
		return StructuredResult{}, errors.New("invalid structured result status")
	}
	wireItems := make([]snapshotItemWire, len(items))
	for index, item := range items {
		wireItems[index] = snapshotItemWire{SnapshotItemRef: item.ref, DisplayLabel: item.displayLabel, DiscoveredAt: item.discoveredAt.Format(time.RFC3339Nano)}
	}
	wire := structuredResultWire{
		SchemaVersion: StructuredResultSchemaVersionV2,
		CapabilityID:  NewVulnerabilityCapabilityID,
		Status:        input.Status,
		Query:         structuredResultQuery{Kind: NewVulnerabilityQueryKind},
		Data:          structuredResultData{Summary: input.Summary, Count: len(items), Items: wireItems},
		Metadata:      structuredResultMeta{Safe: true, Coverage: CoverageCompleteSet, ObservedAt: input.ObservedAt.UTC().Format(time.RFC3339Nano)},
	}
	safeJSON, err := json.Marshal(wire)
	if err != nil {
		return StructuredResult{}, fmt.Errorf("marshal structured result: %w", err)
	}
	sum := sha256.Sum256(safeJSON)
	return StructuredResult{
		internalResultRef: input.InternalResultRef,
		status:            input.Status, summary: input.Summary, observedAt: input.ObservedAt.UTC(),
		items: items, safeJSON: safeJSON, digest: hex.EncodeToString(sum[:]),
	}, nil
}

func (result StructuredResult) SchemaVersion() string     { return StructuredResultSchemaVersionV2 }
func (result StructuredResult) InternalResultRef() string { return result.internalResultRef }
func (result StructuredResult) Status() QueryResultStatus { return result.status }
func (result StructuredResult) Summary() string           { return result.summary }
func (result StructuredResult) Count() int                { return len(result.items) }
func (result StructuredResult) ObservedAt() time.Time     { return result.observedAt }
func (result StructuredResult) ContentDigest() string     { return result.digest }
func (result StructuredResult) ByteSize() int             { return len(result.safeJSON) }
func (result StructuredResult) SafeJSON() []byte          { return append([]byte(nil), result.safeJSON...) }
func (result StructuredResult) Items() []SnapshotItem {
	return append([]SnapshotItem(nil), result.items...)
}

// FreshnessPolicy 是版本化快照有效期规则，TTL 由配置注入。
type FreshnessPolicy struct {
	version string
	ttl     time.Duration
}

func NewFreshnessPolicy(version string, ttl time.Duration) (FreshnessPolicy, error) {
	if strings.TrimSpace(version) == "" || ttl <= 0 {
		return FreshnessPolicy{}, errors.New("invalid freshness policy")
	}
	return FreshnessPolicy{version: version, ttl: ttl}, nil
}

func (policy FreshnessPolicy) Version() string                  { return policy.version }
func (policy FreshnessPolicy) ExpiresAt(at time.Time) time.Time { return at.UTC().Add(policy.ttl) }

// Clock 使快照时间在测试和生产装配中均可控。
type Clock interface{ Now() time.Time }

type QueryResultSnapshotInput struct {
	SnapshotID, WorkspaceID, ConversationID, ActorID, RunID, ToolCallID string
	Result                                                              StructuredResult
	Clock                                                               Clock
	Freshness                                                           FreshnessPolicy
}

// QueryResultSnapshot 冻结查询作用域、来源和完整安全行事实。
type QueryResultSnapshot struct {
	snapshotID, workspaceID, conversationID, actorID, runID, toolCallID string
	capturedAt, expiresAt                                               time.Time
	freshnessVersion                                                    string
	items                                                               []SnapshotItem
}

func NewQueryResultSnapshot(input QueryResultSnapshotInput) (QueryResultSnapshot, error) {
	if input.Clock == nil || input.Freshness.version == "" || input.Result.safeJSON == nil {
		return QueryResultSnapshot{}, errors.New("invalid query result snapshot dependencies")
	}
	for _, value := range []string{input.SnapshotID, input.WorkspaceID, input.ConversationID, input.ActorID, input.RunID, input.ToolCallID} {
		if strings.TrimSpace(value) == "" || ContainsUnsafeMaterial(value) {
			return QueryResultSnapshot{}, errors.New("invalid query result snapshot scope")
		}
	}
	capturedAt := input.Clock.Now().UTC()
	if capturedAt.IsZero() {
		return QueryResultSnapshot{}, errors.New("invalid snapshot clock")
	}
	return QueryResultSnapshot{
		snapshotID: input.SnapshotID, workspaceID: input.WorkspaceID, conversationID: input.ConversationID,
		actorID: input.ActorID, runID: input.RunID, toolCallID: input.ToolCallID,
		capturedAt: capturedAt, expiresAt: input.Freshness.ExpiresAt(capturedAt), freshnessVersion: input.Freshness.Version(),
		items: input.Result.Items(),
	}, nil
}

func (snapshot QueryResultSnapshot) SnapshotID() string       { return snapshot.snapshotID }
func (snapshot QueryResultSnapshot) WorkspaceID() string      { return snapshot.workspaceID }
func (snapshot QueryResultSnapshot) ConversationID() string   { return snapshot.conversationID }
func (snapshot QueryResultSnapshot) ActorID() string          { return snapshot.actorID }
func (snapshot QueryResultSnapshot) RunID() string            { return snapshot.runID }
func (snapshot QueryResultSnapshot) ToolCallID() string       { return snapshot.toolCallID }
func (snapshot QueryResultSnapshot) QueryKind() string        { return NewVulnerabilityQueryKind }
func (snapshot QueryResultSnapshot) CapturedAt() time.Time    { return snapshot.capturedAt }
func (snapshot QueryResultSnapshot) ExpiresAt() time.Time     { return snapshot.expiresAt }
func (snapshot QueryResultSnapshot) FreshnessVersion() string { return snapshot.freshnessVersion }
func (snapshot QueryResultSnapshot) Coverage() string         { return CoverageCompleteSet }
func (snapshot QueryResultSnapshot) Items() []SnapshotItem {
	return append([]SnapshotItem(nil), snapshot.items...)
}

type QueryFactsStatus string

const (
	QueryFactsSucceeded QueryFactsStatus = "succeeded"
	QueryFactsFailed    QueryFactsStatus = "failed"
)

type FactEventType string

const (
	FactEventQuerySucceeded FactEventType = "query.succeeded"
	FactEventQueryFailed    FactEventType = "query.failed"
)

// FactEvent 是最小持久序列事实，不等同于完整 Event Sourcing。
type FactEvent struct {
	eventID   string
	sequence  int64
	typeName  FactEventType
	createdAt time.Time
}

func (event FactEvent) EventID() string      { return event.eventID }
func (event FactEvent) Sequence() int64      { return event.sequence }
func (event FactEvent) Type() FactEventType  { return event.typeName }
func (event FactEvent) CreatedAt() time.Time { return event.createdAt }

type SafeQueryError struct{ Code, Message string }

type SucceededQueryFactsInput struct {
	WorkspaceID, ConversationID, ActorID, RunID string
	QuerySequence                               int64
	CreatedAt                                   time.Time
	Result                                      StructuredResult
	Snapshot                                    QueryResultSnapshot
}

type FailedQueryFactsInput struct {
	WorkspaceID, ConversationID, ActorID, RunID string
	QuerySequence                               int64
	CreatedAt                                   time.Time
	Error                                       SafeQueryError
}

// QueryFacts 是 M1 当前聚合；结果、快照、状态和事件必须原子提交。
type QueryFacts struct {
	workspaceID, conversationID, actorID, runID string
	querySequence                               int64
	status                                      QueryFactsStatus
	createdAt                                   time.Time
	result                                      StructuredResult
	snapshot                                    QueryResultSnapshot
	hasResult, hasSnapshot                      bool
	safeError                                   SafeQueryError
	events                                      []FactEvent
}

func NewSucceededQueryFacts(input SucceededQueryFactsInput) (QueryFacts, error) {
	if err := validateQueryIdentity(input.WorkspaceID, input.ConversationID, input.ActorID, input.RunID, input.QuerySequence, input.CreatedAt); err != nil {
		return QueryFacts{}, err
	}
	if input.Result.safeJSON == nil || input.Snapshot.snapshotID == "" || input.Snapshot.workspaceID != input.WorkspaceID || input.Snapshot.runID != input.RunID {
		return QueryFacts{}, errors.New("query facts result/snapshot mismatch")
	}
	createdAt := input.CreatedAt.UTC()
	return QueryFacts{
		workspaceID: input.WorkspaceID, conversationID: input.ConversationID, actorID: input.ActorID, runID: input.RunID,
		querySequence: input.QuerySequence, status: QueryFactsSucceeded, createdAt: createdAt,
		result: input.Result, snapshot: input.Snapshot, hasResult: true, hasSnapshot: true,
		events: []FactEvent{{eventID: input.RunID + ":1", sequence: 1, typeName: FactEventQuerySucceeded, createdAt: createdAt}},
	}, nil
}

func NewFailedQueryFacts(input FailedQueryFactsInput) (QueryFacts, error) {
	if err := validateQueryIdentity(input.WorkspaceID, input.ConversationID, input.ActorID, input.RunID, input.QuerySequence, input.CreatedAt); err != nil {
		return QueryFacts{}, err
	}
	if input.Error.Code != ErrorCodeFixtureReadFailed || input.Error.Message != FailedQueryMessage || ContainsUnsafeMaterial(input.Error.Message) {
		return QueryFacts{}, ErrUnsafeFactMaterial
	}
	createdAt := input.CreatedAt.UTC()
	return QueryFacts{
		workspaceID: input.WorkspaceID, conversationID: input.ConversationID, actorID: input.ActorID, runID: input.RunID,
		querySequence: input.QuerySequence, status: QueryFactsFailed, createdAt: createdAt, safeError: input.Error,
		events: []FactEvent{{eventID: input.RunID + ":1", sequence: 1, typeName: FactEventQueryFailed, createdAt: createdAt}},
	}, nil
}

func (facts QueryFacts) WorkspaceID() string           { return facts.workspaceID }
func (facts QueryFacts) ConversationID() string        { return facts.conversationID }
func (facts QueryFacts) ActorID() string               { return facts.actorID }
func (facts QueryFacts) RunID() string                 { return facts.runID }
func (facts QueryFacts) QuerySequence() int64          { return facts.querySequence }
func (facts QueryFacts) Status() QueryFactsStatus      { return facts.status }
func (facts QueryFacts) CreatedAt() time.Time          { return facts.createdAt }
func (facts QueryFacts) HasStructuredResult() bool     { return facts.hasResult }
func (facts QueryFacts) HasSnapshot() bool             { return facts.hasSnapshot }
func (facts QueryFacts) Result() StructuredResult      { return cloneStructuredResult(facts.result) }
func (facts QueryFacts) Snapshot() QueryResultSnapshot { return cloneQuerySnapshot(facts.snapshot) }
func (facts QueryFacts) SafeError() SafeQueryError     { return facts.safeError }
func (facts QueryFacts) Events() []FactEvent           { return append([]FactEvent(nil), facts.events...) }

func cloneStructuredResult(result StructuredResult) StructuredResult {
	result.items = append([]SnapshotItem(nil), result.items...)
	result.safeJSON = append([]byte(nil), result.safeJSON...)
	return result
}

func cloneQuerySnapshot(snapshot QueryResultSnapshot) QueryResultSnapshot {
	snapshot.items = append([]SnapshotItem(nil), snapshot.items...)
	return snapshot
}

func cloneQueryFacts(value QueryFacts) QueryFacts {
	value.result = cloneStructuredResult(value.result)
	value.snapshot = cloneQuerySnapshot(value.snapshot)
	value.events = append([]FactEvent(nil), value.events...)
	return value
}

func validateQueryIdentity(workspaceID, conversationID, actorID, runID string, sequence int64, createdAt time.Time) error {
	for _, value := range []string{workspaceID, conversationID, actorID, runID} {
		if strings.TrimSpace(value) == "" || ContainsUnsafeMaterial(value) {
			return errors.New("invalid query identity")
		}
	}
	if sequence < 1 || createdAt.IsZero() {
		return errors.New("invalid query sequence or time")
	}
	return nil
}

func validOpaqueRef(value, prefix string) bool {
	if strings.TrimSpace(value) != value || !strings.HasPrefix(value, prefix) || len(value) <= len(prefix) {
		return false
	}
	for _, char := range value[len(prefix):] {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '_' || char == '-' {
			continue
		}
		return false
	}
	return true
}
