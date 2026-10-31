# P42-T027B-EVIDENCE-BIND ⑦ 终审归档（原样）

切片：`P42-T027B-EVIDENCE-BIND`。环节：⑦ 终审。载体：`independent-reviewer`；模型：`GPT-5.3-Codex (copilot)`。输出原样归档。

---

1. 编号: FR-01
严重度: Low
file:line: internal/adapters/agent_runner_replay_store.go#L221-L225, internal/adapters/evidence_bind_dir_windows_test.go#L106, internal/adapters/evidence_bind_dir_unix_test.go#L148
一句话问题: CloseBoundDirs 被定义为运行期应收口的生命周期动作，但当前仅在测试中调用，生产链路未接线。
建议修法: 在运行时拥有者的收尾路径增加一次性释放（例如 store 生命周期结束处显式调用），并在 ⑧a 记录为已落地或明确延后。

2. 编号: FR-02
严重度: Low
file:line: internal/adapters/bounddir_windows.go#L118-L136
一句话问题: Windows 侧二次复核已实现，但目录归属判断仍是路径前缀语义，reparse/别名边界属于残余风险。
建议修法: ⑧a 逐条登记该残余；后续片可补"对象标识级"校验（句柄对象身份或分层父目录身份校验），减少仅靠前缀判断的歧义面。

3. 编号: FR-03
严重度: Low
file:line: internal/adapters/bounddir_unix.go#L61, internal/adapters/agent_runner_replay_store.go#L260, internal/adapters/agent_runner_postflight.go#L319-L323
一句话问题: Unix 侧 O_NOFOLLOW 只约束叶子，requestId 目录作为中间组件若被替换为绝对目标符号链接，仍存在向量；当前主要由冻结摘要比对兜底为 fail-closed。
建议修法: 后续片改为两段式 dirfd 解析（先安全打开 requestId 目录，再相对打开目标文件），并保留现有摘要对账兜底。

4. 编号: FR-04
严重度: Info
file:line: internal/adapters/agent_runner_postflight.go#L645, internal/adapters/agent_runner_postflight.go#L367
一句话问题: 注释"compared against frozen manifest"与实际 usage 哈希比对语义不完全一致。
建议修法: 在 ⑧a 一并修正文案或登记为文案微瑕，避免误导后续审阅。

5. 编号: FR-05
严重度: Info
file:line: docs/validation/evidence/P42-T027B-EVIDENCE-BIND/05-prereview/review.md#L84
一句话问题: ⑤ 已登记"全量日志仅存一腿"，证据完备性存在轻微展示缺口。
建议修法: ⑧a 明确该事实及不影响结论的理由，或补齐第二腿日志索引/定位信息。

6. 编号: FR-06
严重度: Info
file:line: internal/evidence/boundread.go#L60-L88
一句话问题: 截断语义注释与调用方分类语义存在可读性压缩，容易被后续维护者误解为统一行为。
建议修法: 在 ⑧a 记录并后续统一措辞，明确"记录/PID拒绝、日志/事件可截断"的双类契约。

7. 编号: FR-07
严重度: Info
file:line: docs/validation/t027-replay-convergence.md#L38, docs/validation/evidence/P42-T027B-EVIDENCE-BIND/05-prereview/review.md#L37-L40, docs/S1_SUPPORT_MATRIX.md#L10-L13
一句话问题: 旧验证文档仍写 darwin 交叉通过，与本片 D7 结论和当前支持矩阵口径存在历史不一致。
建议修法: ⑧a 做"口径替代"登记，并在后续文档片修订旧条目。

可证伪性说明
1. FR-01 可被推翻的条件: 在非测试代码出现生产级 CloseBoundDirs 调用并有收尾路径证据；当前仅见测试调用。
2. FR-02 可被推翻的条件: Windows 侧改为对象身份级校验且覆盖 reparse 边界测试通过。
3. FR-03 可被推翻的条件: Unix 侧新增"中间组件不跟随符号链接"的两段式安全打开，并有对应攻击用例通过。
4. FR-04/FR-06 可被推翻的条件: 注释文本与实际判定分支一一对应，不再出现术语错配。
5. FR-05 可被推翻的条件: 证据包补齐第二腿或给出可复算索引。
6. FR-07 可被推翻的条件: 历史文档条目完成 supersede/修订，口径与当前 D7 一致。

交叉抽查（⑤/⑥ file:line，满足 ≥3）
1. ⑤ 的 N1 结论与代码一致: 05-prereview/review.md#L80 对应 internal/adapters/agent_runner_replay_store.go#L221-L225 与仅测试调用位置。
2. ⑤ 的 N2 结论与代码一致: 05-prereview/review.md#L81 对应 internal/adapters/bounddir_windows.go#L118-L136。
3. ⑤ 的 N3 结论与代码一致: 05-prereview/review.md#L82 对应 internal/adapters/bounddir_unix.go#L61 与 internal/adapters/agent_runner_replay_store.go#L260。
4. ⑥ 的"无 Medium+"与"Unix 端仍需 CI 实跑"口径一致: 06-scan/review.md#L24, 06-scan/review.md#L20。

链路一致性与可发布性判定
1. ①→②→⑤⑥ 总体一致，未发现新增 Medium+ 安全回归。
2. D1 已按"绑定 runs 基目录"落地并有冻结摘要兜底一致性证据。
3. D4 口径一致（本片不迁移 BuildReleaseInventory）。
4. D6 目前是登记边界而非闭环修复（见 FR-01）。
5. D7 目前是登记边界且存在历史文档待 supersede（见 FR-07）。
6. ⑤ r2 的 N1-N6 需要在 ⑧a 逐条落，不建议合并成"已知小问题"一条带过；至少给出逐项状态（接受/修复/后续片）与证据锚点。

FINAL REVIEW: PASS WITH FIXES
