# Agent 入口

本文件是所有 AI Agent 的项目入口。不是百科，不是业务事实源。

**禁止**在本仓库另起一套文档系统、工作流、规则文件或工具专用说明书。规则只在 `docs/engineering/ENGINEERING.md`。

## 开始工作

按顺序读：

1. `README.md` — 项目是什么。
2. `.project/state.json` — 阶段和当前目标。
3. `docs/work/CURRENT.md` — 当前状态与 OPEN。

然后按 `.project/context-routes.json` 读对应文档。不要问用户文档在哪。

## 工作方式

- 当前工作是「主要目标」，不是锁死任务。用户插话就先做插话，再回主线。
- 不顺手重构无关问题。不把猜测写成正式事实。
- 卡住时停手：写观察到的失败，看最小证据，形成一个假设，只做一步验证。
- 发现产品/架构冲突时记 OPEN 或写 ADR，不静默改事实。

## 完工四问

1. 做了吗？
2. 测了吗？（被这次改动影响到的检查）
3. 文档需要更新吗？（改了正式事实才更新权威文档）
4. 违反 `ENGINEERING.md` 了吗？

## 事实归属

正式事实只有一个权威文档：归属见 `.project/ownership.json`，路径见 `.project/documents.json`。

`modules/` 是未启用空壳，不是权威。

无法确认归属时标记 OPEN 并询问，不要猜。
