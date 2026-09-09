# 系统架构

> Knowledge Owner：architecture。本文是架构事实的唯一权威来源。
> 只写当前仓库里已经存在的结构。未落地选型写 OPEN 或 ADR，不写进本文当现状。
> 新产品从模版创建后，重写为该产品的真实结构。

## 当前结构

```text
AGENTS.md + ENGINEERING.md     工作循环
        ↓
.project/                      状态、文档登记、归属、路由
        ↓
docs/                          已启用的权威知识
modules/                       未启用空壳（非权威）
scripts/check-project.mjs      registry 校验
```

## 边界

- **OS**：开工、归属、对齐、ADR、校验。跨项目复制。
- **产品实例**：`PRODUCT.md`、`ARCHITECTURE.md`、`CURRENT.md`、`state.json`、按需模块。每个项目自己填。
- **Stack**：语言、框架、目录（apps/packages）不属于 OS。需要时另建 stack 模版，不要写进本文件冒充现状。

## 机器可读契约

| 文件 | 职责 |
|---|---|
| `.project/state.json` | 阶段、当前目标、`os_version` |
| `.project/documents.json` | 权威文档路径；`id` 是主键 |
| `.project/ownership.json` | `knowledge_owner` 必须是 `documents.json` 的 `id` |
| `.project/context-routes.json` | 按任务读的最小文档集 |

`ownership.json` 的 `knowledge_owner` 与 `documents.json` 的 `id` 必须能对上。对不上就是 registry 损坏。
