# Agent 入口

本文件是所有 AI Agent 的项目入口。不是百科，不是业务事实源。

**禁止**在本仓库另起一套文档系统、工作流、规则文件或工具专用说明书。规则只在 `docs/engineering/ENGINEERING.md`。

## 开始工作

按顺序读：

1. `README.md` — 项目是什么（给人看的入口）。
2. `.project/state.json` — 阶段和当前目标。
3. `docs/work/CURRENT.md` — 当前状态与 OPEN。

然后读 `.project/context-routes.json` 的 `always_read`，再读当前任务的 `required`（需要时 `optional`）。不要问用户文档在哪。不要为了解项目而列出或读完 `docs/`。

新建权威文档前先查 `.project/documents.json`：id 已在就改那份；没有对应项就先问或先登记。没有 owner 就标 OPEN，不建文件。

若当前仓库是刚从本模版创建的**新产品**：先填 `PRODUCT.md` / `ARCHITECTURE.md` / `CURRENT.md` / `state.json`，不要设计新的操作系统。

## 工作方式

- 当前工作是「主要目标」，不是锁死任务。用户插话就先做插话，再回主线。
- 不顺手重构无关问题。不把猜测写成正式事实。
- 发现产品/架构冲突时记 OPEN 或写 ADR，不静默改事实。

## 完工四问

1. 做了吗？
2. 测了吗？（被这次改动影响到的检查；本仓至少 `node scripts/check-project.mjs`）
3. 改了正式事实却没改对应权威文档 = 未完成。过程中该改就改。
4. 违反 `ENGINEERING.md` 了吗？

## 事实归属

正式事实只有一个权威文档：归属见 `.project/ownership.json`，路径见 `.project/documents.json`。未登记的文件不是权威。

`modules/` 是未启用空壳，不是权威。

无法确认归属时标记 OPEN 并询问，不要猜。
