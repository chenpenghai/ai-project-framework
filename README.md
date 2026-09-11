# 空白项目模板

和 AI 开新项目，它几乎每次都会自己发明一套文档、规则，有时再加一个任务队列。换个仓库不一样，换个工具又重写一遍。

这份模板把工作系统一次性放进仓库：一个入口、一份规则、事实有归属、用脚本检查有没有写坏。不是 React / 后端 / 数据库脚手架。

换 Claude Code、Codex、Grok、CodeBuddy、ZCode、OpenCode、Cursor 都能读同一套。Gemini CLI 除外。

它具体挡住这些事：

- 禁止再设计第二套文档系统和入口文件
- 禁止给每个 AI 工具各写一份完整规则，宿主入口只许指向 `AGENTS.md`
- 没决定的事标 OPEN，不许写进代码或文档当成事实
- 支付、报告等文档用到再启用；空的权威文件比没有更有害
- 插话先做，不做那种会被插话绕开的任务队列
- 校验脚本检查文档死链、归属对不上、ADR 撞号，不靠模型自觉

## 怎么用

只需要 [下载 template 文件夹](https://download-directory.github.io/?url=https://github.com/chenpenghai/ai-project-framework/tree/main/template)，让 AI 在里面干活。别的不用管。
