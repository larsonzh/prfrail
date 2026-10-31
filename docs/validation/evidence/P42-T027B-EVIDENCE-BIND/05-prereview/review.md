# P42-T027B-EVIDENCE-BIND ⑤ 预审归档（两轮原样）

切片：`P42-T027B-EVIDENCE-BIND`。环节：⑤ 预审。载体：`deep-reasoner`；模型：`DeepSeek V4 Pro (deepseek)`。两轮输出原样归档。

---

## 第 1 轮输出（原样）

# P42-T027B-EVIDENCE-BIND ⑤ 预审（只读）

载体：`deep-reasoner`；模型：DeepSeek V4 Pro。全程未修改任何文件、无 git 写操作。输入：① `review.md`、`impl.diff`（664 行全读）、`status.txt`、实现 11 个新文件 + 5 个测试文件、改动后生产代码全文核对。

## 逐条结论

### 1. 同句柄语义（C3 替换清单逐点验真）— CLOSED
- `BoundFile` 单次 open：`internal/evidence/boundread.go:47-65`（openBoundFileHandle）→ `NewBoundFile`（:75-110）f.Stat 身份、目录/非 regular/超限 open 期拒绝；`ReadAll` bound+1 截断（:116-134）、`Hash` 仅整读成功（:141-160，`written > bound` 拒绝）、`Identity` 永不重 stat（:166-172）；`readLimit` MaxInt64 防溢（:189-194）。
- `readBoundedAgentRunnerArtifact` 的 `os.ReadFile(path)` 分支确已删除（`agent_runner_evidence.go:419-438` 全走 `file.ReadAll`）；`readAgentRunnerPIDReport` stat→read 双解析已消除（:382-400）。逐点 grep 核对：替换清单 9 处全部落地，`filequeue.go` 生产代码零 `os.Stat` 残留；残留的 `os.ReadFile/os.Stat/os.Lstat` 均为清单外边界点（launch/terminal/mirror/slot/root-marker 校验），且各有 P42 boundary 注释登记，无未登记的「stat 后另开读」或「open 后按路径重读」形状（反例甲门禁通过）。
- 缺失保留 `ErrNotExist`：`boundopen_windows.go:24-30` 回传 `*os.PathError`，`TestBoundFileMissingPreservesNotExist` 钉住。

### 2. D1（改绑 `runs` 基目录）— CLOSED WITH NOTE
- (a) 不变式 4：对经 `readBoundEvidence` 的读取（pre/post/diff）保住——Windows 用句柄派生真名跟随对象，Unix `openat` 相对 dirfd 解析；`TestStoreReadBoundEvidenceSurvivesAncestorReplacement` 在 Windows 实测通过。但按字面「后续 I/O 全部经句柄」未满足：见 F3。
- (b) 同层替换残余面：`runs/<requestId>` 在绑定基目录内仍可按名替换；残余受控成立：唯一生产调用方 = `readRecordedArtifact`，且 postflight 对 storedPost/storedDiff 做冻结摘要比对、Leg 1/2 现场重捕重算——替换字节必然 fail-closed，无法通过。
- (c) 应回写设计并登记后续片（读路径全面句柄化 sweep 已列设计 §六）。

### 3. 目录原语正确性 — CLOSED WITH NOTE
- Windows：`CreateFile + FILE_FLAG_BACKUP_SEMANTICS + share RWD`、`GetFinalPathNameByHandle` 派生真名、`GetFileInformationByHandle` 复核 DIRECTORY；Unix：`openat + O_NOFOLLOW + O_CLOEXEC`、fstat 判目录、dev+ino 复核。
- `cleanBoundRelativePath` 逃逸拒绝完备。
- 句柄生命周期：`CloseBoundDirs` 锁外关闭、幂等；缺口见 F1 与 F2。

### 4. 语义漂移与边界（D2/D3/D4/D5/D8）— CLOSED WITH NOTE
- D2：纯新增，零既有符号变更；seam 签名逐字冻结。✓
- D3：保留 Lstat 类型门保住 symlink 叶拒绝语义；hash 校验说法不准确（见 F4），残余窗口极小。
- D4：`BuildReleaseInventory` 未迁移系设计清单项收窄（依赖方向不支持复用）；须回写设计 + 登记。
- D5：新上界 64MiB/1MiB，fail-closed 方向正确。
- D8：记录类拒绝语义在调用方保持；日志 truncated 语义与既有测试一致。

### 5. D7（darwin 交叉编译失效）— CLOSED WITH NOTE
- 实测双工具链 `GOOS=darwin` 均失败于 `bounddir_unix.go:61: undefined: syscall.Openat`；`GOOS=linux` 通过。
- 支持矩阵：macOS = 不支持（正式支持候选 Windows amd64、核心 CI Linux amd64）；CI 矩阵仅 ubuntu/windows。判定可接受。
- 需登记：① 本片起 darwin 交叉编译不再成立；② `docs/validation/t027-replay-convergence.md:38` 曾把「Darwin/amd64 交叉编译通过」登记为验证步骤，该条目失效，须在报告中标注。

### 6. 红测 I1–I5 覆盖与不可测项 — CLOSED WITH NOTE
- I1/I2/I3/I4/I5 覆盖如实现；8 条新测试 Windows 侧全部 PASS。
- 不可测项登记缺口：unix 测试本机未跑（CI 未跑）；unix 侧 CloseBoundDirs 证明力待 CI。

### 7. 371 PASS 真实性抽查 — CLOSED
- 逐字复现 371 PASS / 1 SKIP / 0 FAIL；全量 15 包 ok / 0 FAIL。附加验证：go1.22.12 下双绿。

## 新发现

- F1（Medium）：设计反例乙规定的逐读复核未实现（`bounddir_windows.go:78-100`）。实际后果被冻结摘要比对兜底，属纵深防御缺失；须补 ~10 行复核或显式登记省略理由并回写设计。
- F2（Low）：`bindEvidenceDir` 注释与 `:370` 附近过度声称；须修正。
- F3（Low）：不变式 4 按字面未兑现（hashLogs/fileHash/readPIDs 仍按路径解析）；须登记或整改。
- F4（Low）：D3 注释高估验证强度；须改写。
- F5（Info）：unix 红测未跑、CloseBoundDirs 未接 run 收尾。
- F6（Info）：双工具链差异已实测消除。

PRE-REVIEW: PASS WITH FIXES

---

## 第 2 轮输出（原样）

# P42-T027B-EVIDENCE-BIND ⑤ 预审第 2 轮（只读）— F 项清零确认

只读声明：本环未修改任何文件、未执行任何 git 写操作；全部结论基于工作区直读 + 包内快照/日志机械核验。

### F 项逐条核对

- F1 — CLOSED WITH NOTE：`openBoundDirFile` 派生名→`CreateFile`→`confirmOpenedArtifact`（L127-146）；文件句柄最终路径必须位于目录句柄当次 `currentName()` 之下（`isUnderDirectory` L148-162），不符即 CloseHandle + 拒绝——fail-closed 成立；handle→handle 比较未重走调用方路径；`stripExtendedPathPrefix` 处理 `\\?\`/`\\?\UNC\`。窗口分析成立；Unix openat 原子解析无需后置复核。回归：I3 PASS；全量含 adapters 49.55s。NOTE：风险 1（reparse 解析边界）与风险 2（仅前缀比较）判定可接受，但须在 ⑧a 报告逐条落。
- F2 — CLOSED：全部注释与实现逐条吻合。
- F3 — CLOSED：三处边界登记如实（微瑕 Info：`:643` 措辞 frozen manifest → 实为 Facts.UsageHash）。
- F4 — CLOSED：无残留过度声称（Info 微瑕：`OpenBoundFileTruncating` 注释措辞）。
- 设计回写与包自洽 — CLOSED：修订注记与实现一致、声明冲突以本节为准；review.md 15,727 字节 BOM+LF；SHA256SUMS 两条与磁盘一致；MANIFEST 2 件载荷三方相等；编码合规。
- 回归抽查：focused 667 PASS / 1 SKIP / 0 FAIL；full1 15 ok / 0 FAIL（tmp 仅存档一腿，N5）。
- D1/D3/D4/D6/D7 登记口径：D1 ✓、D3 ✓、D4 ✓；D6/D7 基本可用，⑧a 须落：① CloseBoundDirs 无生产调用方 + 注释矛盾（N1）；② reparse/前缀比较残余（N2）；③ Unix 中间组件符号链接向量（N3）；④ ×2 第二腿日志（N5）。

### 新发现

- N1（Low）`agent_runner_replay_store.go:220-224`：注释自相矛盾且生产无调用方；⑧a 登记 + 措辞统一或后续片接线。
- N2（Low）`bounddir_windows.go:127-162`：reparse 解析边界与仅前缀比较未显式登记；⑧a 须逐条落。
- N3（Low）Unix `runs/<requestId>` 绝对目标符号链接向量；设计原文该句属过度声称，修订注记优先条款可覆盖；⑧a 应精确登记向量。
- N4（Info）`agent_runner_postflight.go:643` 措辞。
- N5（Info）tmp 仅存档全量一腿。
- N6（Info）`boundread.go` 注释措辞。

PRE-REVIEW: PASS
