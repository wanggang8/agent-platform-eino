package execution

import (
	"context"
	"errors"
)

// ErrCheckpointMissing 表示 pending 绑定的内部 checkpoint 不存在或不可恢复。
var ErrCheckpointMissing = errors.New("checkpoint missing")

// CheckpointResolver 只解析 Product Facts 中的安全 checkpoint_ref 到内部 checkpoint。
// 返回的 checkpoint id 只能继续传给 Eino Runner，不能进入产品输出。
type CheckpointResolver interface {
	ResolveCheckpointID(ctx context.Context, checkpointRef string, runID string, pendingID string) (checkpointID string, existed bool, err error)
}
