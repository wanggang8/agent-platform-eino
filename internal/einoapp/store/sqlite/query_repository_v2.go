package sqlite

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"agent-platform-eino/internal/einoapp/facts"
)

const (
	// QuerySchemaEpoch 是 Greenfield M1 唯一支持的数据库 epoch。
	QuerySchemaEpoch     = 2
	minimumSQLiteVersion = "3.51.3"
)

var (
	ErrUnsupportedSchemaEpoch = errors.New("unsupported_schema_epoch")
	ErrRepositoryNotReady     = errors.New("query repository not ready")
	ErrCorruptFacts           = errors.New("corrupt query facts")
)

type QueryWritePoint string

const (
	WriteStructuredResult QueryWritePoint = "structured_result"
	WriteQuerySnapshot    QueryWritePoint = "query_snapshot"
	WriteAggregate        QueryWritePoint = "aggregate"
	WriteFactEvent        QueryWritePoint = "fact_event"
)

// QueryOptions 固定连接池、锁等待和测试故障注入；生产配置必须显式传入。
type QueryOptions struct {
	BusyTimeout   time.Duration
	PoolSize      int
	FaultInjector func(QueryWritePoint) error
}

// ReadinessReport 是不包含连接串或凭据的安全就绪结果。
type ReadinessReport struct {
	Ready         bool
	Reason        string
	SchemaEpoch   int
	SQLiteVersion string
	JournalMode   string
}

// ConnectionInvariant 记录一个真实 sql.Conn 的连接级 PRAGMA。
type ConnectionInvariant struct {
	ForeignKeys bool
	BusyTimeout time.Duration
}

// QueryRepository 是 M1 唯一 SQLite Product Facts 实现。
type QueryRepository struct {
	db        *sql.DB
	options   QueryOptions
	readiness ReadinessReport
}

var _ facts.QueryRepository = (*QueryRepository)(nil)

// OpenQuery 打开全新 epoch；旧 schema 直接拒绝，不迁移、不删库、不 dual write。
func OpenQuery(ctx context.Context, path string, options QueryOptions) (*QueryRepository, error) {
	if strings.TrimSpace(path) == "" || options.BusyTimeout <= 0 || options.PoolSize < 2 {
		return nil, errors.New("invalid query repository options")
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(options.PoolSize)
	db.SetMaxIdleConns(options.PoolSize)
	db.SetConnMaxIdleTime(0)
	db.SetConnMaxLifetime(0)
	fail := func(err error) (*QueryRepository, error) {
		_ = db.Close()
		return nil, err
	}
	if err := db.PingContext(ctx); err != nil {
		return fail(err)
	}
	fresh, epoch, err := inspectQueryEpoch(ctx, db)
	if err != nil {
		return fail(err)
	}
	if !fresh && epoch != QuerySchemaEpoch {
		return fail(ErrUnsupportedSchemaEpoch)
	}
	journalMode, err := enableWAL(ctx, db)
	if err != nil {
		return fail(err)
	}
	if journalMode != "wal" {
		return fail(fmt.Errorf("%w: journal_mode=%s", ErrRepositoryNotReady, journalMode))
	}
	if err := configureConnections(ctx, db, options); err != nil {
		return fail(err)
	}
	if fresh {
		if err := createQuerySchema(ctx, db); err != nil {
			return fail(err)
		}
	}
	version, err := querySQLiteVersion(ctx, db)
	if err != nil {
		return fail(err)
	}
	if !SQLiteVersionAtLeast(version, minimumSQLiteVersion) {
		return fail(fmt.Errorf("%w: sqlite_version=%s", ErrRepositoryNotReady, version))
	}
	if err := verifyDatabaseHealth(ctx, db); err != nil {
		return fail(err)
	}
	repository := &QueryRepository{
		db:        db,
		options:   options,
		readiness: ReadinessReport{Ready: true, SchemaEpoch: QuerySchemaEpoch, SQLiteVersion: version, JournalMode: journalMode},
	}
	return repository, nil
}

func (repo *QueryRepository) Close() error { return repo.db.Close() }

func (repo *QueryRepository) Readiness(_ context.Context) ReadinessReport { return repo.readiness }

// ConnectionInvariants 同时占用多个真实连接，避免反复验证同一连接。
func (repo *QueryRepository) ConnectionInvariants(ctx context.Context, count int) ([]ConnectionInvariant, error) {
	if count < 1 || count > repo.options.PoolSize {
		return nil, errors.New("invalid connection invariant count")
	}
	connections := make([]*sql.Conn, 0, count)
	defer func() {
		for _, connection := range connections {
			_ = connection.Close()
		}
	}()
	for range count {
		connection, err := repo.db.Conn(ctx)
		if err != nil {
			return nil, err
		}
		connections = append(connections, connection)
	}
	values := make([]ConnectionInvariant, 0, count)
	for _, connection := range connections {
		var foreignKeys int
		var busyMilliseconds int64
		if err := connection.QueryRowContext(ctx, `PRAGMA foreign_keys`).Scan(&foreignKeys); err != nil {
			return nil, err
		}
		if err := connection.QueryRowContext(ctx, `PRAGMA busy_timeout`).Scan(&busyMilliseconds); err != nil {
			return nil, err
		}
		values = append(values, ConnectionInvariant{ForeignKeys: foreignKeys == 1, BusyTimeout: time.Duration(busyMilliseconds) * time.Millisecond})
	}
	return values, nil
}

func (repo *QueryRepository) CommitQuery(ctx context.Context, aggregate facts.QueryFacts) error {
	if !repo.readiness.Ready {
		return ErrRepositoryNotReady
	}
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if aggregate.Status() == facts.QueryFactsSucceeded {
		if err := repo.inject(WriteStructuredResult); err != nil {
			return err
		}
		result := aggregate.Result()
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO query_structured_results(run_id, internal_result_ref, schema_version, status, summary, observed_at, safe_json, content_digest, byte_size)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			aggregate.RunID(), result.InternalResultRef(), result.SchemaVersion(), string(result.Status()), result.Summary(), formatQueryTime(result.ObservedAt()), result.SafeJSON(), result.ContentDigest(), result.ByteSize()); err != nil {
			return err
		}

		if err := repo.inject(WriteQuerySnapshot); err != nil {
			return err
		}
		snapshot := aggregate.Snapshot()
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO query_result_snapshots(run_id, snapshot_id, workspace_id, conversation_id, actor_id, tool_call_id, query_kind, captured_at, expires_at, freshness_version, coverage)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			aggregate.RunID(), snapshot.SnapshotID(), snapshot.WorkspaceID(), snapshot.ConversationID(), snapshot.ActorID(), snapshot.ToolCallID(), snapshot.QueryKind(), formatQueryTime(snapshot.CapturedAt()), formatQueryTime(snapshot.ExpiresAt()), snapshot.FreshnessVersion(), snapshot.Coverage()); err != nil {
			return err
		}
		for index, item := range snapshot.Items() {
			if _, err := tx.ExecContext(ctx, `
				INSERT INTO query_snapshot_items(run_id, ordinal, snapshot_item_ref, display_label, discovered_at)
				VALUES (?, ?, ?, ?, ?)`, aggregate.RunID(), index, item.Ref(), item.DisplayLabel(), formatQueryTime(item.DiscoveredAt())); err != nil {
				return err
			}
		}
	}

	if err := repo.inject(WriteAggregate); err != nil {
		return err
	}
	safeError := aggregate.SafeError()
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO query_aggregates(run_id, workspace_id, conversation_id, actor_id, query_sequence, status, created_at, safe_error_code, safe_error_message)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		aggregate.RunID(), aggregate.WorkspaceID(), aggregate.ConversationID(), aggregate.ActorID(), aggregate.QuerySequence(), string(aggregate.Status()), formatQueryTime(aggregate.CreatedAt()), safeError.Code, safeError.Message); err != nil {
		return err
	}

	if err := repo.inject(WriteFactEvent); err != nil {
		return err
	}
	for _, event := range aggregate.Events() {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO query_fact_events(run_id, event_id, sequence, event_type, created_at)
			VALUES (?, ?, ?, ?, ?)`, aggregate.RunID(), event.EventID(), event.Sequence(), string(event.Type()), formatQueryTime(event.CreatedAt())); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (repo *QueryRepository) GetQuery(ctx context.Context, workspaceID, runID string) (facts.QueryFacts, error) {
	row := repo.db.QueryRowContext(ctx, `
		SELECT conversation_id, actor_id, query_sequence, status, created_at, safe_error_code, safe_error_message
		FROM query_aggregates WHERE workspace_id = ? AND run_id = ?`, workspaceID, runID)
	var conversationID, actorID, status, createdText, errorCode, errorMessage string
	var sequence int64
	if err := row.Scan(&conversationID, &actorID, &sequence, &status, &createdText, &errorCode, &errorMessage); errors.Is(err, sql.ErrNoRows) {
		return facts.QueryFacts{}, facts.ErrNotFound
	} else if err != nil {
		return facts.QueryFacts{}, err
	}
	createdAt, err := parseQueryTime(createdText)
	if err != nil {
		return facts.QueryFacts{}, ErrCorruptFacts
	}
	if facts.QueryFactsStatus(status) == facts.QueryFactsFailed {
		aggregate, err := facts.NewFailedQueryFacts(facts.FailedQueryFactsInput{
			WorkspaceID: workspaceID, ConversationID: conversationID, ActorID: actorID, RunID: runID,
			QuerySequence: sequence, CreatedAt: createdAt, Error: facts.SafeQueryError{Code: errorCode, Message: errorMessage},
		})
		if err != nil {
			return facts.QueryFacts{}, ErrCorruptFacts
		}
		if err := repo.verifyEvents(ctx, aggregate); err != nil {
			return facts.QueryFacts{}, err
		}
		return aggregate, nil
	}
	if facts.QueryFactsStatus(status) != facts.QueryFactsSucceeded {
		return facts.QueryFacts{}, ErrCorruptFacts
	}
	return repo.loadSucceededQuery(ctx, workspaceID, conversationID, actorID, runID, sequence, createdAt)
}

func (repo *QueryRepository) LatestQuery(ctx context.Context, workspaceID string) (facts.QueryFacts, error) {
	var runID string
	err := repo.db.QueryRowContext(ctx, `
		SELECT run_id FROM query_aggregates WHERE workspace_id = ?
		ORDER BY query_sequence DESC, created_at DESC, run_id ASC LIMIT 1`, workspaceID).Scan(&runID)
	if errors.Is(err, sql.ErrNoRows) {
		return facts.QueryFacts{}, facts.ErrNotFound
	}
	if err != nil {
		return facts.QueryFacts{}, err
	}
	return repo.GetQuery(ctx, workspaceID, runID)
}

func (repo *QueryRepository) loadSucceededQuery(ctx context.Context, workspaceID, conversationID, actorID, runID string, sequence int64, createdAt time.Time) (facts.QueryFacts, error) {
	var internalRef, schemaVersion, resultStatus, summary, observedText, digest string
	var safeJSON []byte
	var byteSize int
	err := repo.db.QueryRowContext(ctx, `
		SELECT internal_result_ref, schema_version, status, summary, observed_at, safe_json, content_digest, byte_size
		FROM query_structured_results WHERE run_id = ?`, runID).Scan(&internalRef, &schemaVersion, &resultStatus, &summary, &observedText, &safeJSON, &digest, &byteSize)
	if err != nil || schemaVersion != facts.StructuredResultSchemaVersionV2 {
		return facts.QueryFacts{}, ErrCorruptFacts
	}
	observedAt, err := parseQueryTime(observedText)
	if err != nil {
		return facts.QueryFacts{}, ErrCorruptFacts
	}

	var snapshotID, snapshotWorkspace, snapshotConversation, snapshotActor, toolCallID, queryKind, capturedText, expiresText, freshnessVersion, coverage string
	err = repo.db.QueryRowContext(ctx, `
		SELECT snapshot_id, workspace_id, conversation_id, actor_id, tool_call_id, query_kind, captured_at, expires_at, freshness_version, coverage
		FROM query_result_snapshots WHERE run_id = ?`, runID).Scan(&snapshotID, &snapshotWorkspace, &snapshotConversation, &snapshotActor, &toolCallID, &queryKind, &capturedText, &expiresText, &freshnessVersion, &coverage)
	if err != nil || snapshotWorkspace != workspaceID || snapshotConversation != conversationID || snapshotActor != actorID || queryKind != facts.NewVulnerabilityQueryKind || coverage != facts.CoverageCompleteSet {
		return facts.QueryFacts{}, ErrCorruptFacts
	}
	items, err := repo.loadSnapshotItems(ctx, runID)
	if err != nil {
		return facts.QueryFacts{}, err
	}
	result, err := facts.NewStructuredResult(facts.StructuredResultInput{InternalResultRef: internalRef, Status: facts.QueryResultStatus(resultStatus), Summary: summary, ObservedAt: observedAt, Items: items})
	if err != nil || result.ContentDigest() != digest || result.ByteSize() != byteSize || !bytes.Equal(result.SafeJSON(), safeJSON) {
		return facts.QueryFacts{}, ErrCorruptFacts
	}
	capturedAt, err := parseQueryTime(capturedText)
	if err != nil {
		return facts.QueryFacts{}, ErrCorruptFacts
	}
	expiresAt, err := parseQueryTime(expiresText)
	if err != nil || !expiresAt.After(capturedAt) {
		return facts.QueryFacts{}, ErrCorruptFacts
	}
	policy, err := facts.NewFreshnessPolicy(freshnessVersion, expiresAt.Sub(capturedAt))
	if err != nil {
		return facts.QueryFacts{}, ErrCorruptFacts
	}
	snapshot, err := facts.NewQueryResultSnapshot(facts.QueryResultSnapshotInput{
		SnapshotID: snapshotID, WorkspaceID: workspaceID, ConversationID: conversationID, ActorID: actorID,
		RunID: runID, ToolCallID: toolCallID, Result: result, Clock: persistedClock{at: capturedAt}, Freshness: policy,
	})
	if err != nil {
		return facts.QueryFacts{}, ErrCorruptFacts
	}
	aggregate, err := facts.NewSucceededQueryFacts(facts.SucceededQueryFactsInput{
		WorkspaceID: workspaceID, ConversationID: conversationID, ActorID: actorID, RunID: runID,
		QuerySequence: sequence, CreatedAt: createdAt, Result: result, Snapshot: snapshot,
	})
	if err != nil {
		return facts.QueryFacts{}, ErrCorruptFacts
	}
	if err := repo.verifyEvents(ctx, aggregate); err != nil {
		return facts.QueryFacts{}, err
	}
	return aggregate, nil
}

func (repo *QueryRepository) loadSnapshotItems(ctx context.Context, runID string) ([]facts.SnapshotItem, error) {
	rows, err := repo.db.QueryContext(ctx, `
		SELECT snapshot_item_ref, display_label, discovered_at
		FROM query_snapshot_items WHERE run_id = ? ORDER BY ordinal ASC`, runID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []facts.SnapshotItem
	for rows.Next() {
		var ref, label, discoveredText string
		if err := rows.Scan(&ref, &label, &discoveredText); err != nil {
			return nil, err
		}
		discoveredAt, err := parseQueryTime(discoveredText)
		if err != nil {
			return nil, ErrCorruptFacts
		}
		item, err := facts.NewSnapshotItem(ref, label, discoveredAt)
		if err != nil {
			return nil, ErrCorruptFacts
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (repo *QueryRepository) verifyEvents(ctx context.Context, aggregate facts.QueryFacts) error {
	rows, err := repo.db.QueryContext(ctx, `SELECT event_id, sequence, event_type, created_at FROM query_fact_events WHERE run_id = ? ORDER BY sequence`, aggregate.RunID())
	if err != nil {
		return err
	}
	defer rows.Close()
	expected := aggregate.Events()
	index := 0
	for rows.Next() {
		if index >= len(expected) {
			return ErrCorruptFacts
		}
		var eventID, eventType, createdText string
		var sequence int64
		if err := rows.Scan(&eventID, &sequence, &eventType, &createdText); err != nil {
			return err
		}
		createdAt, err := parseQueryTime(createdText)
		if err != nil || eventID != expected[index].EventID() || sequence != expected[index].Sequence() || eventType != string(expected[index].Type()) || !createdAt.Equal(expected[index].CreatedAt()) {
			return ErrCorruptFacts
		}
		index++
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if index != len(expected) {
		return ErrCorruptFacts
	}
	return nil
}

func (repo *QueryRepository) inject(point QueryWritePoint) error {
	if repo.options.FaultInjector == nil {
		return nil
	}
	return repo.options.FaultInjector(point)
}

func inspectQueryEpoch(ctx context.Context, db *sql.DB) (fresh bool, epoch int, err error) {
	var schemaMeta, legacy int
	if err := db.QueryRowContext(ctx, `SELECT count(*) FROM sqlite_master WHERE type='table' AND name='query_schema_meta'`).Scan(&schemaMeta); err != nil {
		return false, 0, err
	}
	if err := db.QueryRowContext(ctx, `SELECT count(*) FROM sqlite_master WHERE type='table' AND name='schema_versions'`).Scan(&legacy); err != nil {
		return false, 0, err
	}
	if schemaMeta == 0 {
		if legacy != 0 {
			return false, 0, nil
		}
		var userTables int
		if err := db.QueryRowContext(ctx, `SELECT count(*) FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'`).Scan(&userTables); err != nil {
			return false, 0, err
		}
		return userTables == 0, 0, nil
	}
	if err := db.QueryRowContext(ctx, `SELECT epoch FROM query_schema_meta WHERE singleton = 1`).Scan(&epoch); err != nil {
		return false, 0, err
	}
	return false, epoch, nil
}

func createQuerySchema(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	statements := []string{
		`CREATE TABLE query_schema_meta (singleton INTEGER PRIMARY KEY CHECK(singleton = 1), epoch INTEGER NOT NULL, created_at TEXT NOT NULL)`,
		`CREATE TABLE query_aggregates (
			run_id TEXT PRIMARY KEY, workspace_id TEXT NOT NULL, conversation_id TEXT NOT NULL, actor_id TEXT NOT NULL,
			query_sequence INTEGER NOT NULL CHECK(query_sequence >= 1), status TEXT NOT NULL CHECK(status IN ('succeeded','failed')),
			created_at TEXT NOT NULL, safe_error_code TEXT NOT NULL DEFAULT '', safe_error_message TEXT NOT NULL DEFAULT '',
			UNIQUE(workspace_id, query_sequence)
		)`,
		`CREATE INDEX idx_query_aggregates_latest ON query_aggregates(workspace_id, query_sequence DESC, created_at DESC)`,
		`CREATE TABLE query_structured_results (
			run_id TEXT PRIMARY KEY, internal_result_ref TEXT NOT NULL UNIQUE, schema_version TEXT NOT NULL, status TEXT NOT NULL,
			summary TEXT NOT NULL, observed_at TEXT NOT NULL, safe_json BLOB NOT NULL, content_digest TEXT NOT NULL, byte_size INTEGER NOT NULL,
			FOREIGN KEY(run_id) REFERENCES query_aggregates(run_id) DEFERRABLE INITIALLY DEFERRED
		)`,
		`CREATE TABLE query_result_snapshots (
			run_id TEXT PRIMARY KEY, snapshot_id TEXT NOT NULL UNIQUE, workspace_id TEXT NOT NULL, conversation_id TEXT NOT NULL, actor_id TEXT NOT NULL,
			tool_call_id TEXT NOT NULL, query_kind TEXT NOT NULL, captured_at TEXT NOT NULL, expires_at TEXT NOT NULL,
			freshness_version TEXT NOT NULL, coverage TEXT NOT NULL,
			FOREIGN KEY(run_id) REFERENCES query_aggregates(run_id) DEFERRABLE INITIALLY DEFERRED
		)`,
		`CREATE TABLE query_snapshot_items (
			run_id TEXT NOT NULL, ordinal INTEGER NOT NULL, snapshot_item_ref TEXT NOT NULL, display_label TEXT NOT NULL, discovered_at TEXT NOT NULL,
			PRIMARY KEY(run_id, ordinal), UNIQUE(run_id, snapshot_item_ref),
			FOREIGN KEY(run_id) REFERENCES query_result_snapshots(run_id) DEFERRABLE INITIALLY DEFERRED
		)`,
		`CREATE TABLE query_fact_events (
			run_id TEXT NOT NULL, event_id TEXT NOT NULL UNIQUE, sequence INTEGER NOT NULL CHECK(sequence >= 1), event_type TEXT NOT NULL, created_at TEXT NOT NULL,
			PRIMARY KEY(run_id, sequence), FOREIGN KEY(run_id) REFERENCES query_aggregates(run_id) DEFERRABLE INITIALLY DEFERRED
		)`,
	}
	for _, statement := range statements {
		if _, err := tx.ExecContext(ctx, statement); err != nil {
			return err
		}
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO query_schema_meta(singleton, epoch, created_at) VALUES (1, ?, ?)`, QuerySchemaEpoch, formatQueryTime(time.Now().UTC())); err != nil {
		return err
	}
	return tx.Commit()
}

func enableWAL(ctx context.Context, db *sql.DB) (string, error) {
	var mode string
	if err := db.QueryRowContext(ctx, `PRAGMA journal_mode=WAL`).Scan(&mode); err != nil {
		return "", err
	}
	return strings.ToLower(mode), nil
}

func configureConnections(ctx context.Context, db *sql.DB, options QueryOptions) error {
	connections := make([]*sql.Conn, 0, options.PoolSize)
	defer func() {
		for _, connection := range connections {
			_ = connection.Close()
		}
	}()
	for range options.PoolSize {
		connection, err := db.Conn(ctx)
		if err != nil {
			return err
		}
		connections = append(connections, connection)
	}
	milliseconds := options.BusyTimeout.Milliseconds()
	for _, connection := range connections {
		if _, err := connection.ExecContext(ctx, `PRAGMA foreign_keys=ON`); err != nil {
			return err
		}
		if _, err := connection.ExecContext(ctx, fmt.Sprintf("PRAGMA busy_timeout=%d", milliseconds)); err != nil {
			return err
		}
		var foreignKeys int
		var configuredBusy int64
		if err := connection.QueryRowContext(ctx, `PRAGMA foreign_keys`).Scan(&foreignKeys); err != nil {
			return err
		}
		if err := connection.QueryRowContext(ctx, `PRAGMA busy_timeout`).Scan(&configuredBusy); err != nil {
			return err
		}
		if foreignKeys != 1 || configuredBusy != milliseconds {
			return ErrRepositoryNotReady
		}
	}
	return nil
}

func verifyDatabaseHealth(ctx context.Context, db *sql.DB) error {
	var integrity string
	if err := db.QueryRowContext(ctx, `PRAGMA integrity_check`).Scan(&integrity); err != nil || integrity != "ok" {
		return fmt.Errorf("%w: integrity_check=%s", ErrRepositoryNotReady, integrity)
	}
	rows, err := db.QueryContext(ctx, `PRAGMA foreign_key_check`)
	if err != nil {
		return err
	}
	if rows.Next() {
		_ = rows.Close()
		return fmt.Errorf("%w: foreign_key_check failed", ErrRepositoryNotReady)
	}
	if err := rows.Close(); err != nil {
		return err
	}
	// TEMP 表属于单一连接；建表和回滚探针必须固定在同一个 sql.Conn 上。
	connection, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	defer connection.Close()
	if _, err := connection.ExecContext(ctx, `CREATE TEMP TABLE IF NOT EXISTS query_readiness_probe(value INTEGER NOT NULL)`); err != nil {
		return fmt.Errorf("%w: write probe: %v", ErrRepositoryNotReady, err)
	}
	tx, err := connection.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO query_readiness_probe(value) VALUES (1)`); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%w: write probe: %v", ErrRepositoryNotReady, err)
	}
	return tx.Rollback()
}

func querySQLiteVersion(ctx context.Context, db *sql.DB) (string, error) {
	var version string
	err := db.QueryRowContext(ctx, `SELECT sqlite_version()`).Scan(&version)
	return version, err
}

// SQLiteVersionAtLeast 按数值元组比较；任一 malformed 输入都 fail closed。
func SQLiteVersionAtLeast(actual, minimum string) bool {
	actualParts, ok := parseSQLiteVersion(actual)
	if !ok {
		return false
	}
	minimumParts, ok := parseSQLiteVersion(minimum)
	if !ok {
		return false
	}
	for index := range actualParts {
		if actualParts[index] > minimumParts[index] {
			return true
		}
		if actualParts[index] < minimumParts[index] {
			return false
		}
	}
	return true
}

func parseSQLiteVersion(value string) ([3]int, bool) {
	var parsed [3]int
	parts := strings.Split(value, ".")
	if len(parts) != 3 {
		return parsed, false
	}
	for index, part := range parts {
		number, err := strconv.Atoi(part)
		if err != nil || number < 0 {
			return parsed, false
		}
		parsed[index] = number
	}
	return parsed, true
}

type persistedClock struct{ at time.Time }

func (clock persistedClock) Now() time.Time { return clock.at }

func formatQueryTime(value time.Time) string { return value.UTC().Format(time.RFC3339Nano) }

func parseQueryTime(value string) (time.Time, error) { return time.Parse(time.RFC3339Nano, value) }
