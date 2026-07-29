# ADR：提高 modernc SQLite 的 WAL 安全基线

状态：Accepted

日期：2026-07-20

## 背景

Story 1.2 / M1 必须启用 SQLite WAL，并通过统一连接策略验证 `journal_mode=wal`、每连接外键、busy timeout、schema epoch 和 readiness。当前 `go.mod` 与架构版本表固定 `modernc.org/sqlite v1.34.5`；该模块在主要目标平台内含 SQLite 3.46.0。

SQLite 官方于 2026-04-13 更新 WAL 文档，确认 SQLite 3.7.0～3.51.2 存在低概率但可能造成数据库损坏的 WAL-reset 竞态，修复版本为 3.51.3 及以上（另有部分旧稳定线回移版本）。`modernc.org/sqlite` 官方 changelog 记录 `v1.46.2` 升级到 SQLite 3.51.3。

继续以 v1.34.5 实现多连接 WAL 会让 Story 1.2 的持久化与重启恢复基线建立在已知缺陷版本上，不符合 fail-closed 目标。

## 决策

- Story 1.2 开始实现时，把 `modernc.org/sqlite` 的精确目标版本从 `v1.34.5` 提高到 `v1.46.2`。
- 不使用浮动版本；`go.mod`、`go.sum`、架构版本表和 canonical 工具链证据必须记录同一精确版本。
- 启动与测试必须读取 `sqlite_version()`，并拒绝低于 3.51.3 的运行时进入 ready 状态。
- SQLite 仍保持纯 Go、单实例、本地文件和 WAL 方案；本 ADR 不改变数据库选型、schema epoch、repository 边界或事务模型。
- 本次只做安全必需升级，不顺带采用 v1.46.2 之后的功能，也不做无关依赖升级。

## 备选方案

- 保留 v1.34.5 并限制为单连接：无法满足已批准的统一连接／pool 策略，也把安全性依赖于容易漂移的运行参数，拒绝。
- 关闭 WAL：与 AD-18、AD-20、AD-25 及 Story 1.2 验收要求冲突，拒绝。
- 直接采用当前最新 modernc SQLite：包含更多与 M1 无关的变化，扩大回归面；在没有额外需求时不采用。
- 更换 SQLite driver：会改变 CGO、工具链和部署假设，超出本次修复范围，拒绝。

## 影响

- Story 1.2 的第一项实现任务必须更新模块并在 canonical 环境运行 Go、race、SQLite restart 与完整基线验证。
- Story 1.1 的历史 `G-TOOLCHAIN=PASS` 仍证明工具链本身；它不能替代升级后 Story 1.2 的新 clean GitHub Actions 证据。
- 若 v1.46.2 在固定 Go 1.26.5 或目标 linux/amd64 环境无法通过现有门禁，Story 1.2 保持阻断，必须新建 ADR 重新选择，而不是降级回受影响版本。

## 官方依据

- SQLite WAL：<https://www.sqlite.org/wal.html>（2026-07-20 访问；WAL-reset 缺陷与修复版本）
- modernc SQLite changelog：<https://gitlab.com/cznic/sqlite/-/blob/master/CHANGELOG.md>（2026-07-20 访问；v1.46.2 对应 SQLite 3.51.3）

## 关联验收

```bash
go list -m modernc.org/sqlite
go test ./internal/einoapp/store/sqlite -count=1
go test -race ./internal/einoapp/store/sqlite ./internal/einoapp/execution -count=1
go test ./...
bash scripts/run_toolchain_baseline.sh
```
