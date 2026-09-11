# 项目操作系统模版

`os_version: 2`

给新项目用的 **Agent 工作系统**：入口、规则、文档归属、当前状态、可选领域空壳。

不是业务脚手架。不带 React / Fastify / 数据库。不绑任何 AI 工具或模型。

## 新项目怎么用

1. 用本仓库 GitHub Template 创建，或 clone 后改 remote。
2. 把下面这段作为给 AI 的第一句话：

> 这是从 ai-project-framework 创建的项目。禁止重建文档系统、禁止发明 queue / Work Unit / 每工具一份规则。只改这四份：`docs/product/PRODUCT.md`、`docs/architecture/ARCHITECTURE.md`、`docs/work/CURRENT.md`、`.project/state.json`。需要设计/数据/安全/支付/分析/报告时，从 `modules/` 拷到 `docs/` 并登记进 `.project/documents.json`。改完跑 `node scripts/check-project.mjs`。

3. 跑 `node scripts/check-project.mjs`。
4. 宿主若强制要 `CLAUDE.md` / Copilot 入口：只保留一行指针到 `AGENTS.md`，不要复制规则正文。

## 本仓库自己

本仓的产品就是这套 OS。改 OS 时走同一套循环。不要往这里塞业务项目。

## 目录

```text
AGENTS.md                         Agent 入口（指针，不是百科）
docs/engineering/ENGINEERING.md   规则唯一权威
docs/product/PRODUCT.md           产品事实
docs/architecture/ARCHITECTURE.md 架构现状（只写已落地的）
docs/work/CURRENT.md              现在在做什么
docs/decisions/                   ADR（不改写历史）
.project/                         机器可读 registry
modules/                          未启用的领域空壳
scripts/check-project.mjs         registry 校验
```

核心永远在 `docs/` 且已登记。`modules/` 不是权威。
