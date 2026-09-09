# 可选领域模块

这里的文件**不是权威**。启用步骤：

1. 拷到 `docs/<name>/`（保持文件名）。
2. 在 `.project/documents.json` 增加一条，`id` 用下面的 id。
3. 若有运行时事实，在 `.project/ownership.json` 增加 fact，`knowledge_owner` 等于该 id。
4. 需要时在 `.project/context-routes.json` 给相关 task 的 `optional` 加上该 id。
5. 跑 `node scripts/check-project.mjs`。

| id | 源 | 何时启用 |
|---|---|---|
| design | `modules/design/DESIGN.md` | 有视觉/交互规则 |
| data-model | `modules/data/DATA-MODEL.md` | 有持久化模型 |
| security | `modules/security/SECURITY.md` | 有身份、密钥、隐私 |
| payment | `modules/payment/PAYMENT.md` | 有订单/支付 |
| events | `modules/analytics/EVENTS.md` | 有埋点/指标 |
| report-engine | `modules/report/REPORT-ENGINE.md` | 有评分/报告 |

不要提前启用。空的权威文档比没有更有害。
