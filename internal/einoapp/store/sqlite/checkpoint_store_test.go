package sqlite_test

import (
	"bytes"
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/store/sqlite"
)

func TestCheckPointStorePersistsCheckpointAndSafeRefAcrossRestart(t *testing.T) {
	// CheckPointStore 是 Eino 内部恢复材料，Product Facts 只能引用安全 checkpoint_ref。
	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "checkpoint.db")
	now := time.Unix(700, 0).UTC()

	store, err := sqlite.OpenCheckpointStore(ctx, dbPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Set(ctx, "eino-internal-checkpoint-1", []byte("serialized-eino-checkpoint")); err != nil {
		t.Fatal(err)
	}
	if err := store.BindCheckpointRef(ctx, "checkpoint_ref:safe-1", "eino-internal-checkpoint-1", "run-1", "pending-1", now); err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}

	reopened, err := sqlite.OpenCheckpointStore(ctx, dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()

	payload, existed, err := reopened.Get(ctx, "eino-internal-checkpoint-1")
	if err != nil {
		t.Fatal(err)
	}
	if !existed || !bytes.Equal(payload, []byte("serialized-eino-checkpoint")) {
		t.Fatalf("checkpoint payload existed=%v payload=%q", existed, payload)
	}
	checkpointID, existed, err := reopened.ResolveCheckpointID(ctx, "checkpoint_ref:safe-1", "run-1", "pending-1")
	if err != nil {
		t.Fatal(err)
	}
	if !existed || checkpointID != "eino-internal-checkpoint-1" {
		t.Fatalf("resolved checkpoint id = %q existed=%v", checkpointID, existed)
	}
}

func TestCheckPointStoreRejectsUnsafeSafeRefAndReportsMissing(t *testing.T) {
	// raw checkpoint 标识不能伪装成 Product Facts 中的 checkpoint_ref。
	ctx := context.Background()
	store, err := sqlite.OpenCheckpointStore(ctx, filepath.Join(t.TempDir(), "checkpoint.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	for _, testCase := range []struct {
		name          string
		checkpointRef string
		checkpointID  string
	}{
		{name: "explicit raw marker", checkpointRef: "checkpoint-raw-eino-id", checkpointID: "eino-internal-checkpoint-1"},
		{name: "same as internal id", checkpointRef: "cp1", checkpointID: "cp1"},
		{name: "missing safe prefix", checkpointRef: "550e8400-e29b-41d4-a716-446655440000", checkpointID: "550e8400-e29b-41d4-a716-446655440000"},
		{name: "outer whitespace", checkpointRef: " checkpoint_ref:space ", checkpointID: "eino-internal-checkpoint-1"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			err := store.BindCheckpointRef(ctx, testCase.checkpointRef, testCase.checkpointID, "run-1", "pending-1", time.Time{})
			if !errors.Is(err, facts.ErrUnsafeFactMaterial) {
				t.Fatalf("unsafe checkpoint ref err = %v, want ErrUnsafeFactMaterial", err)
			}
		})
	}

	checkpointID, existed, err := store.ResolveCheckpointID(ctx, "checkpoint_ref:missing", "run-1", "pending-1")
	if err != nil {
		t.Fatal(err)
	}
	if existed || checkpointID != "" {
		t.Fatalf("missing checkpoint resolved id=%q existed=%v", checkpointID, existed)
	}
}

func TestCheckPointStoreBindingIsImmutable(t *testing.T) {
	// checkpoint_ref 首次绑定后不可静默重绑，避免旧 pending 被动变成 checkpoint_missing。
	ctx := context.Background()
	store, err := sqlite.OpenCheckpointStore(ctx, filepath.Join(t.TempDir(), "checkpoint.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	if err := store.Set(ctx, "eino-internal-checkpoint-1", []byte("serialized")); err != nil {
		t.Fatal(err)
	}
	if err := store.Set(ctx, "eino-internal-checkpoint-2", []byte("serialized-2")); err != nil {
		t.Fatal(err)
	}
	if err := store.BindCheckpointRef(ctx, "checkpoint_ref:immutable", "eino-internal-checkpoint-1", "run-1", "pending-1", time.Time{}); err != nil {
		t.Fatal(err)
	}
	if err := store.BindCheckpointRef(ctx, "checkpoint_ref:immutable", "eino-internal-checkpoint-1", "run-1", "pending-1", time.Time{}); err != nil {
		t.Fatalf("same checkpoint binding should be idempotent: %v", err)
	}
	if err := store.BindCheckpointRef(ctx, "checkpoint_ref:immutable", "eino-internal-checkpoint-2", "run-1", "pending-1", time.Time{}); !errors.Is(err, facts.ErrIdempotencyConflict) {
		t.Fatalf("rebinding checkpoint ref err = %v, want ErrIdempotencyConflict", err)
	}
}

func TestCheckPointStoreResolveRequiresRunAndPendingBinding(t *testing.T) {
	// 同一个安全 checkpoint_ref 只能给绑定的 run/pending 使用，避免恢复到其他 pending 的 checkpoint。
	ctx := context.Background()
	store, err := sqlite.OpenCheckpointStore(ctx, filepath.Join(t.TempDir(), "checkpoint.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	if err := store.Set(ctx, "eino-internal-checkpoint-1", []byte("serialized")); err != nil {
		t.Fatal(err)
	}
	if err := store.BindCheckpointRef(ctx, "checkpoint_ref:safe-1", "eino-internal-checkpoint-1", "run-1", "pending-1", time.Time{}); err != nil {
		t.Fatal(err)
	}
	if checkpointID, existed, err := store.ResolveCheckpointID(ctx, "checkpoint_ref:safe-1", "run-2", "pending-2"); err != nil || existed || checkpointID != "" {
		t.Fatalf("misbound checkpoint resolved id=%q existed=%v err=%v", checkpointID, existed, err)
	}
}
