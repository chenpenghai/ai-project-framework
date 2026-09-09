# ADR-0001 轻量项目操作系统

- Status: accepted
- Date: 2026-09-09

## 背景

每个新项目都让 AI 从零设计文档系统和规则，成本高、结果还不一致。曾经试过两套更重的方案：

1. **APF 插件层**（本仓库旧形态）：宿主 hook、状态机循环、扫描器/图。和「在仓库里放一套可复制的规则」不是同一个产品。
2. **iqtest 的 Work Unit 队列**：`queue.json`、七状态、`rules.json` 与 ENGINEERING 双重定义。单人 + AI 会被插话打乱，机制被绕过。

需要一套能直接复制进新仓库的轻量 OS。

## 决策

采用 **工具无关、仓库内、可复制** 的轻量操作系统（`os_version: 2`）：

- `AGENTS.md` 唯一 Agent 入口，不存业务事实。
- `ENGINEERING.md` 唯一规则正文。
- `.project/` 只保留 state / documents / ownership / context-routes。
- 领域文档按需从 `modules/` 启用。
- 工作循环：读状态 → 干活（允许插话）→ 改事实则同步文档 → 四问完工。
- 用 `scripts/check-project.mjs` 校验 registry，不靠模型自觉发现撞号和死链。

## 原因

规则必须在仓库里，换 Agent 才不会丢。轻量循环匹配实际插话节奏。校验脚本能抓住 iqtest 已经出现的损坏：ADR 撞号、owner 对不上文档 id、OPEN 清单双写。

## 后果

正面：新项目从模版创建，不再设计 OS。  
代价：没有自动强制「改了代码必改文档」；靠四问和校验脚本。OS 升级靠 `os_version`，旧项目可以不跟。

## 拒绝的方案

1. 继续做宿主插件：解决的是 IDE 时机，不是新项目冷启动。
2. 保留 Work Unit / queue / rules.json：已被实践否定。
3. 每个 AI 工具一份完整规则：必然分叉。
4. 把 React/Fastify 等 stack 打进 OS：下一个非该栈项目会整套推倒。
5. 预铺支付/报告等全部领域文档：未用到的权威文件会变成假事实。
