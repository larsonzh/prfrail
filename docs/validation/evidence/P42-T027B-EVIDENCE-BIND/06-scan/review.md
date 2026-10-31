# P42-T027B-EVIDENCE-BIND ⑥ 独立扫描归档（原样）

切片：`P42-T027B-EVIDENCE-BIND`。环节：⑥ 独立扫描。载体：`independent-reviewer`；模型：`MAI-Code-1.1-Flash (copilot)`。第二轮（格式合规重发）输出原样归档；首轮输出因格式不合规未采纳（额度 2/3）。

---

## Section A — 总体结论
本次只读复核未发现 P42-T027B-EVIDENCE-BIND 的 Medium+ 缺陷，当前实现与设计一致，并已按 C3+N2 句柄绑定/目录绑定语义落地。

## Section B — 发现清单
- 无

## Section C — 逐条证据与推理
在 [internal/evidence/boundread.go](internal/evidence/boundread.go#L23-L152) 中，`OpenBoundFile` 只保留单个已打开的 `*os.File`，随后 `ReadAll`、`Hash`、`Identity` 都从同一句柄读取，不再按路径做第二次解析；这与“打开后再按照路径重读”的 TOCTOU 缺陷相反。 [internal/adapters/bounddir_windows.go](internal/adapters/bounddir_windows.go#L39-L196) 和 [internal/adapters/bounddir_unix.go](internal/adapters/bounddir_unix.go#L18-L99) 分别通过 Windows 的 `GetFinalPathNameByHandle` 复核与 Unix 的 `openat(..., O_NOFOLLOW)` 维持目录绑定，确保祖先目录被替换时读取仍落到原对象。 [internal/adapters/agent_runner_evidence.go](internal/adapters/agent_runner_evidence.go#L100-L201) 中的 `FrozenFacts` 也要求缺口统一返回 `ErrAgentRunnerEvidenceGap`，避免把“未完成”伪装成空摘要；其既定边界在 [internal/adapters/agent_runner_postflight.go](internal/adapters/agent_runner_postflight.go#L554-L613) 和 [internal/adapters/agent_runner_run.go](internal/adapters/agent_runner_run.go#L317-L346) 中被明确标注为本 slice 的 residual boundary，而非未修复漏洞。 [docs/validation/evidence/P42-T027B-EVIDENCE-BIND/01-arch-design/review.md](docs/validation/evidence/P42-T027B-EVIDENCE-BIND/01-arch-design/review.md#L31-L101) 的 D1/D4/F1 修订注记也保持一致：Q2 绑定对象收敛到基目录、释放外部不在本 slice 迁移、Windows 逐读复核已实现。

## Section D — 可证伪性审计
一个直接反例是：在 [internal/evidence/boundread.go](internal/evidence/boundread.go#L83-L152) 中删去 `Seek(0, io.SeekStart)` 和 `io.LimitReader(bound.file, readLimit(bound.bound))` 的同句柄读逻辑，改回 `os.ReadFile(path)`，或者删除 [internal/evidence/boundread_test.go](internal/evidence/boundread_test.go#L19-L92) 里 `TestBoundFileReadsOriginalAfterPathReplacement` 的断言；该变更会让“路径替换后仍读到旧文件”的安全判定失效，结论会从 PASS 翻成 FAIL。另一个反例是去掉 [internal/adapters/evidence_bind_dir_windows_test.go](internal/adapters/evidence_bind_dir_windows_test.go#L31-L88) / [internal/adapters/evidence_bind_dir_unix_test.go](internal/adapters/evidence_bind_dir_unix_test.go#L31-L147) 的祖先替换测试，或删除 [internal/adapters/bounddir_unix.go](internal/adapters/bounddir_unix.go#L43-L66) 的 `O_NOFOLLOW`/`openat` 绑定；那么目录绑定就无法证明“原对象仍被读取”，PASS 也会被推翻。

## Section E — 我无法验证的部分与什么会推翻 PASS
1. Unix 上的 `openat(..., O_NOFOLLOW)` 与 symlink 反制场景需要在 CI/Unix 环境中真实执行；本次只读复核无法替代该端到端验证，因此 [internal/adapters/evidence_bind_dir_unix_test.go](internal/adapters/evidence_bind_dir_unix_test.go#L98-L147) 的 Unix 分支仍是“需 CI 保留”的部分。
2. 本 slice 明确保留的 residual boundary 仍处于后续读路径 sweep 范围，例如 [internal/adapters/agent_runner_run.go](internal/adapters/agent_runner_run.go#L316-L346) 中的 `readAgentRunnerChildPIDs` 以及 [internal/adapters/agent_runner_replay_store.go](internal/adapters/agent_runner_replay_store.go#L220-L270) 的部分读入口；如果后续实现把这些位置也改回 `os.ReadFile`/`os.Stat` 后再读，或如果以上红测在 Windows/Unix 上任一失败，PASS 将被推翻。
3. 只要同句柄绑定、目录绑定和缺口门槛仍保持不变，且 Windows 的全量/聚焦检查结果维持为 0 FAIL，则本结论继续成立。

INDEPENDENT SCAN: PASS
