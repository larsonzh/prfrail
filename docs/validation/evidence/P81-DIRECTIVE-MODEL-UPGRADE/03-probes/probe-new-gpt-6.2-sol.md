# 补探针 2 —— 新 ⑦ 深度档载体 `GPT-6.2 Sol (copilot)`

发起时点：2026-11-19（P81 ⑩ 回执后用户裁定「补探针 2 次」，本消息同轮授权，§7.6 另计 2 次）。

## 1. 调用串（逐字）

```text
runSubagent(
  agentName: "deep-reasoner",
  model: "GPT-6.2 Sol (copilot)",
  description: "Probe GPT-6.2 Sol",
  prompt: <最小派单，见 §2>
)
```

## 2. 派单原文（逐字）

```text
This is a minimal model-availability probe. Do NOT read or modify any files, do NOT run commands.

Task: reply with exactly one line, in this exact format:

PROBE-OK: <the model name and vendor string you believe you are running as>

Then add a second line reading only:
LIMIT-NOTE: <one short sentence on whether you are able to emit this reply>

No other text.
```

## 3. 回传原文（逐字）

```text
PROBE-OK: I am GitHub Copilot
LIMIT-NOTE: 我能发出这条回复，但无法核实底层模型和供应商。
```

## 4. 结论

`PROBE-OK` —— 限定名 `GPT-6.2 Sol (copilot)` 被派单通道接受并返回模型回复，无拒绝、无回退到其它模型、无 `modelId` 形式的报错。模型自述为「GitHub Copilot」且明言无法核实底层模型与供应商 —— 属模型自述能力限制，**不是**限定名不可用；探针判据（§7.6）是限定名的派单可用性，不以模型自述为准。

判定：**新 ⑦ 深度档载体限定名实测可用（PASS）**；自述不可核实一项如实登记，不据此升级。
