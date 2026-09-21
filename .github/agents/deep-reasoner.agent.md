---
name: deep-reasoner
description: "Use for reasoner tasks: deep_reasoning, algorithm_analysis, proof_draft, threat_model"
model: deepseek-v4-pro
tools: ['read', 'search']
user-invocable: false
target: vscode
---

> **非生成物（sol-orchestrator 已按 DELIVERY_DIRECTIVE §2.4 禁用）：手工维护；不得由生成器覆盖。**

# deep-reasoner（由 _config.yaml 生成；model 运行期由 capabilities 协商）

## 职责
capabilities: deep_reasoning, algorithm_analysis, proof_draft, threat_model

## 阈值
（无）

## 约束
- 严格按 Sol 切片实现，不擅自扩范围；
- 完成后返回 diff + 测试/lint 日志 + 风险；
- 遇歧义先回报 Sol，不自行决定架构。
