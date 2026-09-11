# 工程规则

> Knowledge Owner：engineering。本文是工程规范与工作循环的唯一权威来源。
> 流程正文只写在这里。`AGENTS.md` 只保留入口指针。

`os_version: 2`

## 不变量

1. 代码、测试、文档描述同一个系统事实。
2. 同一规则只有一个正式实现位置。
3. 未决定的内容标 OPEN，不得被代码或文档偷偷固定。
4. 权威文档只写现状；变更历史进 `CURRENT.md` / ADR；实现细节留在代码里。
5. 未落地的技术选型不得写成 Accepted 现状。写 ADR 草稿或 OPEN。
6. 密钥与私有凭证不进客户端，不进文档。
7. 声称完成前，被这次改动影响的检查必须通过。

## 文档规则

- 每种正式事实一个 Knowledge Owner，见 `.project/ownership.json`。
- 路径与登记见 `.project/documents.json`。未登记的文件不是权威。
- `documents.json` 的 `id` 是主键。`ownership.json` 的 `knowledge_owner` 必须是某个 `id`。
- `documents[].authority` 里的每个名字必须是 `ownership.json` 的 fact `id`。
- `modules/` 下的空壳未拷贝、未登记前不是权威。
- README、注释、聊天记录不得形成第二套规则。
- 重大决定进 `docs/decisions/`。Accepted 的 ADR 不改写，由新 ADR 按**文件路径** supersede。
- ADR 编号 `NNNN-slug.md`，NNNN 全局唯一，禁止两个 `0001`。

改了正式事实，同步对应权威文档再继续。纯样式/文案不必。

## 工作循环

```text
读状态 → 做当前目标（允许插话）→ 改变事实则同步文档 → 完工四问
```

- `.project/state.json` 记录当前主要目标。同一时间一个主目标。
- 插话是正常工作方式。先完成插话，再回主线。
- 阶段外的小问题可以直接修；影响大的记入 `CURRENT.md` 的 OPEN，不扩散范围。

开工三步、完工四问见 `AGENTS.md`。不要在这里再复制一份。

## 禁止

- 新建第二套 OS、AGENTS、`.project` 规则、Work Unit、任务队列、七状态生命周期。
- `.project/queue.json`、`.project/rules.json`。
- 为 Claude / Codex / Cursor / Copilot / Gemini 各写一份完整规则。宿主若强制要入口文件，只放一行指针到 `AGENTS.md`。
- 把模版空壳或假数据当成产品事实。
- 在 `AGENTS.md` 复制本文正文。

## 校验

```text
node scripts/check-project.mjs
```

改 registry 或权威文档路径后必须通过。
