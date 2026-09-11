# 系统架构

> Knowledge Owner：architecture。本文是架构事实的唯一权威来源。
> 只写当前仓库里已经存在的结构。未落地选型写 OPEN 或 ADR，不写进本文当现状。

## 当前结构

已落地的是工作循环与文档归属，不是业务系统。

```text
AGENTS.md + ENGINEERING.md     工作循环
.project/                      状态、文档登记、归属、路由
docs/                          已启用的权威知识
modules/                       未启用空壳（非权威）
scripts/check-project.mjs      registry 校验
```

业务目录尚未落地。

## 边界

产品事实写在 `docs/product/PRODUCT.md`。未决定的技术选型标 OPEN。
