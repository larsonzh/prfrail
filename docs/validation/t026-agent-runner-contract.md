﻿# T026 AgentRunner 契约与能力验证报告

[English](t026-agent-runner-contract_EN.md)

日期：2026-09-10。结论：`CONTRACT COMPLETE / RUNTIME PROBE BLOCKED BY AUTHENTICATION / AT-23 NOT PASS`。

## 已完成范围

1. 四份 Schema 冻结 request、capability、event 和 completion wire；契约门禁现覆盖 34 份 Schema、14 个 fixture 文件承载的 96 个用例及 2 条 canonical 向量。
2. `internal/adapters/agent_runner*.go` 实现 RFC 8785 域分隔摘要、严格解码、排序集合、create/resume 约束、request/event/completion 绑定，以及可从持久记录重建的幂等/冲突索引。
3. capability 固定 12 项矩阵；仅全项带运行证据且为 verified 才可 compatible。event 同时锁定 eventId、session/sequence 和 request/session，禁止同一 request 漂移到不同 session。completion 同时锁定 completionId 与 requestId，且 completed 只证明外部执行结束、进程树停止、日志/用量完整和输出 manifest 捕获，不代表 task PASS。固定 capability 测试还会读取实际证据文件并核对其字节 SHA-256，防止记录引用不存在或过时的证据。

## 候选探针

- 候选：GitHub Copilot CLI 1.0.83（Windows x64，构建提交 `e6a98f1`）。原生 `copilot.exe` 为 144,796,448 字节，SHA-256 为 `d3f3bb7b8bbf68357ad29f514a179d09f76135483d8bfb643131b8600f671ee2`；平台包元数据摘要为 `d6d07e0943462974bba6380170ca6c6d147e6e609b2cc5ceff165f8f16faa5ed`。
- 绕过 VS Code bootstrap 和 npm shim 直接执行原生二进制，`--version` 返回 1.0.83。`--help` 仅用于设计探针，不作为能力通过证据。
- 用户已明确授权一次最小、可能计费的模型调用。首个启动脚本因 Windows PowerShell 5.1 不支持 `Get-Date -AsUTC` 而在 Copilot 启动前失败，不计为授权调用；修正后的调用使用隔离 workspace、JSONL stream、30 AI credits 软上限、仅 `shell` 工具、精确 `shell(pwd)` 自动批准、HTTP/HTTPS deny、禁用 built-in MCP、系统临时目录、自定义指令、远程控制/导出及自动更新。
- 修正后的 CLI 进程未产生 JSONL、stderr 或 `usage.json`。468 字节的专用进程日志（SHA-256 `7d1b6cd0a3f696a9233b51071fff75c94542e33b2a0c1c697ef2008841ac5129`）显示认证编排拒绝了当前 GitHub 凭据，原因为 classic PAT 不受支持；日志不包含凭据值。进程在该错误后被中断，事后无匹配探针进程。没有模型响应或用量证据，现有产物不能证明是否到达模型端或计费，因此没有在“一次调用”授权下重试。
- 用户再次授权后，第二次预检仅在 Copilot 子进程中移除 `GH_TOKEN`，不修改用户环境或配置，以尝试回落到本地 OAuth 凭据。CLI 明确返回“未找到认证信息”并以 1 退出；381 字节的 `usage.json`（SHA-256 `f9f90cadcbcb1d9ba7176d520df7585abc79ff3ef751937204c366c6765b0c98`）记录 0 次用户请求、0 输入/输出 token、0 API 时长、0 premium request cost 且无 model metrics。因此第二次授权也没有发生模型调用或计费。
- 第二次预检显式使用 `--model auto`。Copilot CLI 不继承 VS Code 聊天输入框选择的模型；可用 `--model <name>` 为单次运行指定 CLI 支持的模型，或用 `--model auto` 让 CLI 自主路由。未显式传参时使用 CLI 自身的 `model` 配置或默认值。
- 非交互、cwd、事件/完整日志、session create/resume、取消、进程树停止、工具/网络/权限、用量及无人值守确认共 12 项仍记为 unknown，整体 blocked。机器证据及带 hash 的 blocked capability record 位于 `testdata/agent-runner/capability-probes/`。

## 门禁结果

1. `go build ./...`：通过。
2. `go vet ./...`：通过。
3. `go test -count=1 ./...`：通过，13 个含测试包及 2 个无测试 command 包。
4. `node tools/contracts/contracts.test.js`：通过，2/2；34 份 Schema、96 个用例和 2 条 canonical 向量通过。
5. `git diff --check`：通过。

## 解阻条件

若要重试，需先在用户自己的终端完成 CLI 1.0.83 支持的 OAuth 登录，或配置具备 `Copilot Requests` 权限的 fine-grained PAT，再取得一次新的最小、可能计费模型调用授权。后续应先用可观察的前台输出确认 CLI 已进入非交互模式，再在隔离 workspace 中逐项运行 12 项探针并更新 capability record。全部 required 项 verified 后才可关闭 T026 并开始 T027；当前不得宣称 AT-23 PASS。