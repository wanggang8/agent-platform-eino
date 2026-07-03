package execution

import (
	"context"
	"errors"
	"time"
)

// ErrCheckpointMissing 表示 pending 绑定的内部 checkpoint 不存在或不可恢复。
var ErrCheckpointMissing = errors.New("checkpoint missing")

// CheckpointResolver 只解析 Product Facts 中的安全 checkpoint_ref 到内部 checkpoint。
// 返回的 checkpoint id 只能继续传给 Eino Runner，不能进入产品输出。
type CheckpointResolver interface {
	ResolveCheckpointID(ctx context.Context, checkpointRef string, runID string, pendingID string) (checkpointID string, existed bool, err error)
}

// ApprovalCheckpointStore 保存 approval 中断点和安全引用映射。
// Set/Get 使用内部 checkpoint id；BindCheckpointRef 只接收可进入 Product Facts 的安全 checkpoint_ref。
type ApprovalCheckpointStore interface {
	CheckpointResolver
	Set(ctx context.Context, checkpointID string, checkpoint []byte) error
	Get(ctx context.Context, checkpointID string) ([]byte, bool, error)
	BindCheckpointRef(ctx context.Context, checkpointRef string, checkpointID string, runID string, pendingID string, createdAt time.Time) error
}
