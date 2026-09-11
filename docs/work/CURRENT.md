# 当前工作状态

> 这是项目「现在在做什么」的唯一权威记录，不记录长期业务规则。
> 新产品从模版创建后，重写本节。

## 阶段

模版源：OS v2。

## 当前目标

把文档系统做到用户拷走 `template/` 就能用：收口生成物、补 OPEN 双写校验、把槽位写成可填结构。不在本仓做业务功能。

## 已完成

- 用轻量 Current Work 模型替换旧 APF（插件 / 队列 / 扫描器 / Go 宿主）
- 删除旧代码：`.agents/`、`cmd/`、`internal/`、`prototype/`
- 核心权威文档 + registry + 校验脚本
- 领域文档改为 `modules/` 可选空壳
- 提交时自动生成 `template/`；GitHub Action 只校验
- README 写清挡住的问题；`template/` 不再带 `.github`
- OPEN 双写校验；产品/架构/当前工作槽位改为可填结构

## 阻塞 / OPEN

产品未决事项的权威清单只在 `docs/product/PRODUCT.md`，此处不复制。

## 下一阶段

用户拷走 `template/` 开新产品。不要在本仓扩业务。
