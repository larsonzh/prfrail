# T017 Verifier Baseline 独立主审存档

[English](t017-baseline-review_EN.md)

日期：2026-09-08。结论：`PASS FOR BASELINE REVIEW`。未发现可复现的 `BLOCKER` 或 `HIGH`；允许将当前受审实现作为 verifier baseline 的候选提交，但本结论不是 T017 完成、发布授权、发布者身份认证或 GitHub CI 成功证明。

## 评审来源与边界

1. 评审对象是基于 HEAD `fc99f553c3066b24230a093d14bc17af0c6a398e` 的 T017 未提交工作树，包括 workflow、发布 verifier、工作流契约/变异测试、fixture 和相关文档；它在提交前不是不可变对象。
2. 评审由用户发起的 GPT-5.3 Codex 严格只读主审完成，口径要求对抗性复核、findings 优先、禁止改文件/联网/提交/推送/发布，并明确区分本地验证与真实 GitHub CI。
3. 本页归档会话结论及本地命令证据，不是签名 receipt 或可独立认证的模型身份记录。baseline 落库前仍须由维护者核对实际 diff 与精确文件清单。
4. 状态保持 `IMPLEMENTED_AWAITING_CI`；只有较早 verifier baseline commit 验证不同 candidate commit，且首次 GitHub Windows/Linux 原生矩阵通过并归档 run 证据后，才可关闭 T017。

另见结论相同的 [DeepSeek V4 Pro 独立主审存档](t017-deepseek-baseline-review.md)。两份记录相互印证，但不构成签名或 attestation。

## Findings

1. `MEDIUM`：完整脚本 SHA-256 是内容漂移门禁，不是 PowerShell 语义证明。若 workflow 与测试中的期望摘要在同一变更中被恶意同步修改，本地测试可以重新变绿；控制措施是独立评审 baseline，禁止从待验 workflow 自动刷新摘要。后续可考虑将摘要基线外置为受保护只读工件。
2. `LOW`：生产代码离线静态检查覆盖直接 `net` import 和已列举下载器/联网 Git 字面量，但命令词表不是完备的 OS 网络沙箱。当前执行器不声明网络隔离能力，`deny/loopback/allowlist` 均在执行前 fail-close。
3. `LOW`：本地工具链为 Go 1.27.0，而 CI 固定 Go 1.22.12；本地全绿不能替代目标 Windows/Linux runner 结果。

## 本地验证证据

以下检查由主审在 `d:\LZProjects\prfrail` 只读或仅产生常规 Go 构建缓存的条件下执行，退出码均为 0：

1. `CGO_ENABLED=0 go build ./...`
2. `CGO_ENABLED=0 go vet ./...`
3. `CGO_ENABLED=0 go test -count=1 ./...`
4. 聚焦 workflow 契约与变异测试：通过；共 540 个变异子测试。
5. `git diff --check`
6. 关键文件编辑器诊断：无错误。
7. 改动文件编码/行尾检查：通过；Markdown 为 UTF-8 BOM + LF，Go/YAML/JSON 为 UTF-8 无 BOM + LF。
8. 生产 Go 文件静态抽查：直接 `net` import 为 0；`os/exec` 仅位于已允许执行边界，未发现常见下载器或联网 Git 子命令。
9. `tmp/` 仅保留 `.gitkeep`；`sessbridge` 未修改。

## Baseline 落库门禁

1. 提交前复核当前工作树相对 `fc99f553c3066b24230a093d14bc17af0c6a398e` 的完整 diff，确认除本存档与链接外没有主审后代码漂移。
2. 使用精确文件列表暂存并建立单独 verifier baseline commit；不得把未审改动、临时文件或候选提交混入 baseline。
3. baseline commit 推送 `origin/main` 后，再建立内容不同且 SHA 不同的 candidate commit。
4. 从 `main` 手工 dispatch workflow，记录 verifier SHA、candidate SHA、run URL、平台结果和最终 artifact/SHA256SUMS 摘要。
5. 任一平台、输入 pin 或证据校验失败均保持阻断；不得将本页的 PASS 当作发布授权。

## 未执行事项

本次主审未执行 `git commit`、`git push`、tag、签名、attestation、GitHub Release 或 workflow dispatch，也未重跑真实 runner 攻击演练。