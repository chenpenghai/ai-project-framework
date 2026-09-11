# 项目操作系统模版

`os_version: 2`

> **新项目不要用仓库根目录。只复制 `template/` 文件夹。**

## 怎么复制

1. **下载 zip（不用 git）**  
   打开这个链接会只打包 `template/`：[下载 template 文件夹](https://download-directory.github.io/?url=https://github.com/chenpenghai/ai-project-framework/tree/main/template)

2. **一条命令**

   ```text
   npx degit chenpenghai/ai-project-framework/template my-project
   ```

3. **已经 clone 了整仓**  
   把 `template/` 拷到别处，在新目录里执行 `git init`。

复制后先填这四份：`docs/product/PRODUCT.md`、`docs/architecture/ARCHITECTURE.md`、`docs/work/CURRENT.md`、`.project/state.json`。然后跑 `node scripts/check-project.mjs`。

给 AI 的第一句话：

> 这是从 ai-project-framework 的 template/ 创建的项目。禁止重建文档系统、禁止发明 queue / Work Unit / 每工具一份规则。只改这四份：`docs/product/PRODUCT.md`、`docs/architecture/ARCHITECTURE.md`、`docs/work/CURRENT.md`、`.project/state.json`。需要设计/数据/安全/支付/分析/报告时，从 `modules/` 拷到 `docs/` 并登记进 `.project/documents.json`。改完跑 `node scripts/check-project.mjs`。

浏览目录：[GitHub 上的 template/](https://github.com/chenpenghai/ai-project-framework/tree/main/template)

## 本仓库自己

本仓的产品就是这套 OS。根目录是咱们的源，不要往这里塞业务项目。

提交时会自动生成 `template/` 并一起提交。GitHub Action 只校验，不生成。

## 目录

```text
AGENTS.md                         Agent 入口（指针，不是百科）
docs/engineering/ENGINEERING.md   规则唯一权威
docs/product/PRODUCT.md           产品事实（本仓：OS 模版）
docs/architecture/ARCHITECTURE.md 架构现状（只写已落地的）
docs/work/CURRENT.md              现在在做什么
docs/decisions/                   ADR（不改写历史）
.project/                         机器可读 registry
modules/                          未启用的领域空壳
scripts/check-project.mjs         registry 校验
scripts/build-template.mjs        生成本地 template/
template/                         用户拷走的干净实例
```

核心永远在 `docs/` 且已登记。`modules/` 不是权威。
