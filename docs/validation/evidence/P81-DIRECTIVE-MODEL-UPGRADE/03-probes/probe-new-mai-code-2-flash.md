# 补探针 1 —— 新 ⑥ 载体 `MAI Code 2 Flash (copilot)`

发起时点：2026-11-19（P81 ⑩ 回执后用户裁定「补探针 2 次」，本消息同轮授权，§7.6 另计 2 次）。

## 1. 调用串（逐字）

```text
runSubagent(
  agentName: "deep-reasoner",
  model: "MAI Code 2 Flash (copilot)",
  description: "Probe MAI Code 2 Flash",
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
PROBE-OK: MAI Code 2 Flash (Microsoft)
LIMIT-NOTE: I am able to emit this reply.
```

## 4. 结论

`PROBE-OK` —— 限定名 `MAI Code 2 Flash (copilot)` 被派单通道接受并返回模型回复，无拒绝、无回退到其它模型、无 `modelId` 形式的报错。自述供应商串为 `Microsoft`（与 `(copilot)` 限定名不冲突；限定名的供应商段是本仓库的派单语法，不是模型自述）。

判定：**新 ⑥ 载体限定名实测可用（PASS）**。
