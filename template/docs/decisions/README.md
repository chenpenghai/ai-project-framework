# 决策记录

本文只说明 ADR 的格式和编号规则，不写产品或架构事实。

重大产品/架构/OS 决定写在这里的 `NNNN-slug.md`。

## 规则

- 文件名：`NNNN-short-slug.md`，`NNNN` 四位数字，仓库内唯一。
- 状态：`draft` | `accepted` | `superseded`。
- `accepted` 之后不改写正文。要改决定，新建 ADR，并写明 `Supersedes: docs/decisions/NNNN-old.md`（写路径，不写「ADR-0001」这种会撞号的称呼）。
- 被替代的旧文件把状态改为 `superseded`，并指向新文件路径。
- 接受后的「当前系统是什么」写进对应权威文档（产品/架构/工程），ADR 只保留「为什么」。

## 模板

```markdown
# ADR-NNNN 标题

- Status: accepted
- Date: YYYY-MM-DD
- Supersedes: （无则删此行）

## 背景
## 决策
## 原因
## 后果
## 拒绝的方案
```
