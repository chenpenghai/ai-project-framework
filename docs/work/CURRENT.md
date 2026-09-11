# 当前工作状态

> 这是项目「现在在做什么」的唯一权威记录，不记录长期业务规则。
> 新产品从模版创建后，重写本节。

## 阶段

模版源：OS v2。

## 当前目标

加固文档执行：多余 md 校验失败、按路由读、完工必须对齐权威文档。不在本仓做业务功能。

## 已完成

- 用轻量 Current Work 模型替换旧 APF（插件 / 队列 / 扫描器 / Go 宿主）
- 删除旧代码：`.agents/`、`cmd/`、`internal/`、`prototype/`
- 核心权威文档 + registry + 校验脚本
- 领域文档改为 `modules/` 可选空壳
- 提交时自动生成 `template/`；GitHub Action 只校验
- README 写清挡住的问题；`template/` 不再带 `.github`
- OPEN 双写校验；产品/架构/当前工作槽位改为可填结构
- `docs/` 未登记 md 校验失败；入口按 always_read + 路由读；四问第 3 条改为未对齐即未完成

## 阻塞 / OPEN

产品未决事项的权威清单只在 `docs/product/PRODUCT.md`，此处不复制。

## 下一阶段

用户拷走 `template/` 开新产品。不要在本仓扩业务。
