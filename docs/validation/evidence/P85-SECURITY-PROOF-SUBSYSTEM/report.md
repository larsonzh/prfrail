# P85 `SECURITY-PROOF-SUBSYSTEM` 切片验证报告

- 切片：P85（实现类 `[SLICE]`，L 档，协议先行）；目标 = 闭合 P83 §8 顺延的 F3/F6/F8（T2 冻结策略制品 / T3 不可变安装来源 / T8a BYOK 证明载体）。
- 结论（代码复审口径）：**T2 / T3 / T8a 全部 CLOSED**；⑨ 证据债 F17/F22 已于 2026-11-23 凭双腿 CI 证据销账（见 §12；缺证不得关片口径由此闭合）。
- 候选：暂存工作树（`HEAD = e35e50f`，R=7 批终态，全文件精确暂存、porcelain 无 MM/AM/??）。
- 日期：2026-11-23。主控 = DeepSeek V4 Pro Max（用户同轮授权，起跑后固定）。

## §1 异常登记区

| # | 事项 | 处置 |
|---|---|---|
| E-1 | R 计数达 7（§9.2 R>3 暂停） | 每轮停点履行人工审核 + 用户同轮授权续跑；OB-142 已登记（不得作先例） |
| E-2 | ⑦ 额度 1+3 用尽后追加四读（Carbon/Neon/Sol/Codex） | OB-141（Codex 第四读）、Sol 第五读按 §7.8.3 用户授权例外留痕；三读试水系治理片试水授权；均不得作先例 |
| E-3 | R=7 后定向小读（Carbon High，仅读 R=7 diff+测试） | OB-143（§1.7.4 临时授权，不得作先例） |
| E-4 | L4 注释措辞修改发生在审查环结束后 | 依据 P26 裁定 D-2：附 diff 摘录（`04-reviews/09-l4-diff-excerpt.md`）+ 顾问实质核验 + 用户终审（⑩ 回执登记） |
| E-5 | Linux 目录槽夹具初版最小 ACL 与 mode 等价、内核不落盘，VM 首跑该用例失败 | 夹具改扩展 ACL（命名用户 u:4242 + mask）后双平台全绿；性质=测试夹具修正，非生产代码回退 |
| E-6 | Windows 发布 `CreateTemp`/`Link` 按路径写入残余 | 用户终裁：登记 **S1 独立边界**（见 §8），不挂 F04 豁免、不改代码路径 |
| E-7 | gopls unused 跨 GOOS 变体诊断 12 符号（修复后重扫） | 三类处置完成，三视图清零（§10/§12） |
| E-8 | 主提交 CI windows 腿 2 例环境性失败（卷根争用 + 8.3 短名拼写；本机/VM 均绿） | **R=8** 测试侧修复（生产零改动）：卷根探针争用重试+skip、父目录命名降级对象身份比对；提交 `6480bb1`（用户 2026-11-23 授权）；复跑 CI 45208001981 双腿绿（§12） |

## §2 变更概要

- **T2（F3）冻结策略制品**：装配前冻结 + 按 run/task/step/attempt 重读 + 发布前复核；`persistAgentRunnerFrozenPolicyProven` 生产装配不注入 proof（同名无证明包装已删）；四槽摘要 `:2` 固定 `[策略目录|父目录|制品|锚]`、空槽 `∅`、拒绝路径保持四槽（R=7）。
- **T3（F6）不可变安装来源**：`openInstallSourceLeaf` 改 `OpenBoundDirVerified`（同柄分类/摘要、错误携带 directoryID）；载体发布前后根链重绑定对证；Windows `verifyInstallSourceImmutability` 保留 leaf、探测前后 `FileIDAt == leaf.FileID()`（R=7）；shebang 256 字节窗口 + ELF PT_INTERP 不变。
- **T8a（F8）BYOK 证明载体**：生产装配不注入 Sandbox，双平台默认 `VerdictUnknown` ⇒ 凭据解析前零启动拒发，分类 `byok_confinement_unproven`；S1 保持 BYOK 禁用、不宣称 OS 写限制已获证明。
- **契约**：`docs/CONTRACTS{,_EN}.md` :214（byok_confinement_unproven 闭集）、:216（T3 安装来源 + Windows 根链）、:220（BYOK 零启动门 + fd3）、:231（ManagedWriteProofShape 创建/复用）、:233（域 :2 四槽）。
- **准则台账**：`docs/DELIVERY_DIRECTIVE{,_EN}.md` OB-140/141/142/143（四要素、不得作先例）。
- **测试**：R4 期 M13–M17、R6 期 slot_order/sandwich/windows_forensics/windows_publish_recheck、R7 期 `p85_r7_role_suffix` / `p85_r7_refusal_slot` / `p85_r7_open_bounddir_verified` + 目录槽夹具（扩展 ACL 判别 + 闭端核验腿）共 20+ 用例。
- **IDE 告警清零**：死代码 9 删（含连带 `openAgentRunnerPolicyFile`）；3 组符号归位（`managedAgentPolicyWriteProofDomain`→`linux||windows` 文件；`boundRootName`/`errBoundRootArtifact`→`!darwin` 文件（保 freebsd 编译）；unsupported 删 `installSourceInterpreter`；`p85R4SwapPolicyDirectory`→linux 测试文件）；零 `//nolint`。

## §3 执行流水线（步骤 / 角色 / 状态 / 返工）

① 设计 Rev1–Rev4（V4 Pro Max）→ ②a 契约先行（CONTRACTS 五处双语）→ ②b 实现（Flash High）→ ③ 测试 →
④ 门禁 → ⑤ 预审 r1/r2 → ⑥ 扫描 r3/r4 → ⑦ 盲审（22 条）→ ⑦ 终审（Sol Max）→ 复审 r1/r2/r3 →
R=4/R=5/R=6 修复 → 四读同输入（Carbon/Neon/Sol/Codex）→ **R=7 合并修复批**（OB-142）→
定向小读（Carbon High，OB-143）→ ⑧a 初稿 → ⑨ 原生验证 → ⑧b 定稿 → ⑧c 收尾 → ⑩ 停点。
**返工计数 R = 7**（r3 后 R=4/5/6 三组 + R=7 合并批；每轮用户授权，OB-140/142 留痕）。

## §4 环境与档位

- 本机 Windows 11 amd64（非提升 shell）；Ubuntu VM 10.0.0.199（larson 非 root，go1.22.12，`GOPROXY=goproxy.cn`）。
- 主控/① = DeepSeek V4 Pro Max；②③ = DeepSeek V4.1 Flash High；⑤ = DeepSeek V4 Pro；
  ⑥ = MAI Code 2 Flash（+1 轮 OB-140）；⑦ = GPT-6.2 Sol Max（1+3）+ 四读（§5）；
  定向小读 = Gemini 4 Carbon High（OB-143）。探针 = 0。

## §5 审查结论摘要（⑤ ⑥ ⑦）

- ⑤ 预审 r1 FINDINGS → r2 PASS；⑥ 扫描 r3/r4；⑦ 盲审 FINDINGS（22 条）→ 终审 FINDINGS（3H+2M）→
  复审 r1/r2/r3 FINDINGS → R=4/R=5/R=6 → **四读同输入矩阵** → R=7 → 定向小读。
- **四读矩阵（同输入 = R6 批；按「判级分歧」记录，n=1 单输入观测，不作模型能力排序）**：

| 读 | 结论 | 组1/组2/组3 | T2/T3/T8a |
|---|---|---|---|
| Gemini 4 Carbon | PASS WITH FIXES | C/C/C | C/C/C |
| Gemini 4 Neon | FINDINGS（4M+2L） | N/C/C | N/N/C |
| GPT-6.2 Sol | FINDINGS（1M+3L） | N/N/C | N/N/C |
| GPT-5.3-Codex（盲） | PASS | C/C/C | C/C/C |

  判级分歧行（P87 数据点）：组 1 Carbon/Codex CLOSED vs Neon/Sol NOT-CLOSED；组 2 Carbon/Neon/Codex
  CLOSED vs Sol NOT-CLOSED；T2/T3 Carbon/Codex CLOSED vs Neon/Sol NOT-CLOSED。核心观察 = 同输入下
  PASS 与 FINDINGS 并存（判级哲学差异：fail-closed 兜底 vs 契约完备性）；判级标准共识与派单策略
  由 P87 治理片裁定。P87 原始数据已归档 `docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/P87-four-read-data/`。
- **R=7 合并批**：Neon 4M+2L + Sol 3L + 用户终裁 S1 边界（§8）；Codex 新发现 = 0。
- **定向小读（⑧c 输入）**：`TARGETED REVIEW: PASS WITH FIXES` —— 9 项全 CLOSED、0 Medium+、4 Low 不阻断
  （L1 p85_t3_windows_f04 传 nil leaf 后不再走祖先拒绝分支，生产由 VerifyAIInstallSource 覆盖；
  L2 名义入口大写后缀 role=-1（仅非生产入口）；L3 aclErr 丢 probeErr 未 Join（仅诊断）；L4 注释措辞（已修，§1 E-4））。
  L1/L2/L3 用户裁定接受不修、如实记入本报告。

## §6 ⑧a 自查（四条，有/无 + 证据路径）

| # | 项 | 结论 | 证据路径 |
|---|---|---|---|
| 1 | 成本表机械审计 | 有 | §14 成本表；credits 用户登记口径（Sol 67.9/14m50s、Codex 129.9/26m58s、四读合计 ≈958） |
| 2 | 「已处置」file:line | 有 | `docs/DELIVERY_DIRECTIVE.md` L1184–L1187（OB-140/141/142/143）；CONTRACTS :214/:216/:220/:231/:233 |
| 3 | 旧口径 grep | 有 | 命令与结果 = `05-input/old-wording-grep.txt`（仅存「not an F04 exemption」新口径句与 `bounddir_windows.go` 三处 F04 合法登记——spawn 最终名/artifact 打开最终名/目录 flush，均非发布路径） |
| 4 | ⑥/⑦ 每轮 Info/Low 全落报告 | 有 | 见 §5：Carbon 3L + Neon 2L + Sol 3L + 定向小读 4L 全部列出；⑥ Info 见 R4–R6 记录 |

## §7 A-NN 编号完整性

- 盲审 22 条 F01–F22 逐条处置留痕（⑨ 前全落）；R 轮次 MUT-1..MUT-5；OB-140/141/142/143；DR 无新增；编号连续无缺口。

## §8 Windows 发布路径 S1 边界独立登记（用户 2026-11-23 终裁）

- **边界事实**：Windows 下 `os.CreateTemp` + `os.Link` 无法相对目录句柄执行，发布写入按用户态路径解析；
  前后根链重绑定对证（`OpenBoundDirVerified` 两次）在写前/写后执行，但前后对证之间被插链/替换时
  链接已发生、后置拒绝不能撤回。no-replace 硬链接提供原子可见性，**不提供断电耐久**。
- **登记口径**：本条为 **S1 用户态路径发布边界**（用户 2026-11-23 终裁、独立登记），**不是** F04 豁免；
  代码注释 5 处已改登本口径（`agent_runner_policy_publish_other.go:14-20`、`ai_install_proof_publish_other.go:13-19`、
  `ai_install_proof.go:126/258` 等），`bounddir_windows.go` 三处既有 F04 提及（spawn 最终名 / artifact 打开最终名 /
  目录 flush）为**其它已登记残余**，不属本边界。
- **后续**：S2 若引入句柄相对创建/链接或转交由受管执行者写入，本边界重新评估。

## §9 审查输入包摘要（含盲审隔离证明）

- ⑦ 各轮输入：冻结 diff（diff.patch 776KB）+ 双语契约 + 双平台读数；**盲审首轮未附任何审查清单**（隔离成立）；
  四读同输入 = R6 批（r3 三条 + 三组修复 + 契约槽位句 + 双平台读数），Codex 为盲读（禁读前三腿报告，已声明）。
- R=7 定向小读输入 = 仅 `r7-fix.diff`（22 文件 78KB）+ R7 起跑消息 + 对应测试，无全量重审。

## §10 可证伪性（机械检查 + 变异检验）

- **变异矩阵（全部红 → 逐字节复原 → 绿）**：MUT-1 槽序退回路径序（R6）；MUT-2 Windows 三明治（R6）；
  MUT-3 载体发布复核回退 FileIDAt（R6）；**MUT-4A** fd→按名 ACL（红 = `TestP85R7LinuxDirectorySlotReadsTheACLFromItsDescriptor`，
  VM 非 root）；**MUT-4B** 删除闭名对证（裁定指定；红 = `TestP85R7LinuxDirectorySlotClosingComparisonRefusesAMovedName`，
  VM 非 root）；**MUT-5** `OpenBoundDirVerified`→`OpenBoundDir` 跟链（红 = `TestP85R7OpenInstallSourceLeafRefusesALinkedChain`，
  Windows 非提升）。复原均 SHA-256 逐字节一致（`05-input/mutation-r7-executed.txt`、`mutation-r6.txt`）。
- **反例**：缺失/畸形/合形错配载体拒绝；跨角色同名认领拒绝；Role/后缀矛盾与大小写拒绝；junction/symlink 链拒绑。

## §11 门禁结果

- `node tools/gates/gate.js --all --skip-base --scope=index|tree|ci`：**TOTAL_FAIL=0 × 3**；`--selftest` **PASS 740/740**。
- 三变体交叉构建 `GOOS=windows/linux/darwin GOARCH=amd64 go build ./...` 全绿；三视图 `go vet` 全绿。
- 三视图 `gopls check -severity=hint`（internal/**/*.go 411 文件）过滤 `unused` = **0 行 × 3**。
- 暂存卫生：porcelain 全部 `M/A` 精确暂存、MM/AM/?? = 0；未提交（等 ⑩ 授权）。

## §12 ⑨ 原生验证结果（双端）

- **Windows 本机（非提升）**：`go build ./...` exit=0；`go test ./...` 15 包 ok / FAIL=0，exit=0；
  `go vet ./...` exit=0。逐腿（R7 新用例）：`OpenInstallSourceLeafRefusesALinkedChain` PASS（junction 腿）、
  `UnobservableSlotKeepsItsPositionInTheDigest` PASS、`RoleSuffixContradictionIsRefused` 8/8 PASS；无 SKIP。
  原始日志 = `05-input/local-full-log.txt`、`local-r7-legs.txt`。
- **Ubuntu VM（larson 非 root，go1.22.12）**：`go build ./...` exit=0；`go test ./...` 15 包 ok / FAIL=0，exit=0；
  `go vet ./...` exit=0。逐腿：目录槽三腿（`ReadsTheACLFromItsDescriptor` / `RefusesAReplacedName` /
  `ClosingComparisonRefusesAMovedName`）PASS；`OpenInstallSourceLeafRefusesALinkedChain` PASS（symlink 腿）；
  `UnobservableSlotKeepsItsPositionInTheDigest` PASS；`RoleSuffixContradictionIsRefused` 8/8 PASS。
  原始日志 = `05-input/vm-full-log.txt`、`vm-r7-legs.txt`。
- **F17 特权夹具（四要素如实登记）**：裁决人=用户；日期=2026-11-23；适用=全片；理由=F17 豁免只限同 inode
  特权腿——本片未执行需 chown/SeTakeOwnership 的特权夹具，相关正对照以非提升/非 root 原生腿替代并逐腿记录；
  依据=⑦ 终审口径 + 用户裁定。**未执行即如实标注，不记闭合**。
- **CI 双腿（已销账）**：主提交 `cc7d90f` 首跑 CI 45201921016：ubuntu ✓、windows 腿 2 例环境性失败（§1 E-8）→
  R=8 测试侧修复（`6480bb1`）→ 复跑 CI **45208001981 双腿绿**：windows 4m55s / ubuntu 3m25s；
  Gates 步 **TOTAL_FAIL=0 × 2**；`adapters`/`console`/`guard` 包级 `ok` 行双腿全有（另 ubuntu Race 三包 ok）；
  Test 步 `--- FAIL` = 0。原始证据 = `05-input/ci-45208001981-evidence.txt`（已归档本包）。
  **F17/F22 销账**：双腿 Gates + 包级 ok 硬证据齐备 ⇒ 两项证据债按 ⑨ 口径销账；缺证不得关片条件闭合。

## §13 ⑧c 收尾清洁记录

- 全仓 `gofmt -l` 零输出；三视图 gopls unused 清零；编码门禁：变更文件 BOM/EOL 全检
  （md=BOM+LF、go/json/txt=noBOM+LF）；`tmp/` 根仅 `.gitkeep` + 既有切片目录；未跟踪文件 0（全部精确暂存）。
- 证据包自洽：P85 证据目录 MANIFEST 与 SHA256SUMS 同步（§9 载荷；冻结包不参与逐文件编码改写）。

## §14 成本与计量（用户登记口径）

- ⑦ 法定额度：终审 1 + 复审 3（F-01 用户裁定）+ 盲审 1（F-03）——共 5 次 Sol Max 调用。
- 追加：四读（Carbon / Neon / Sol 第 5 读 / Codex）+ 定向小读（Carbon High，OB-143）。
- credits：Sol 基准读 67.9 / 14m50s；Codex 129.9 / 26m58s；四读合计 ≈ 958（其余读数与 ② 两单 Flash credits 由用户登记）。
- 探针 = 0；⑥ = MAI Code 2 Flash（额度 3 + OB-140 追加 1）。

## §15 明确未执行事项

| # | 事项 | 原因与口径 |
|---|---|---|
| 1 | F17 同 inode 特权夹具 | 非提升/非 root 无法构造；四要素登记（§12），不记闭合 |
| 2 | P85 双腿 CI（Gates + 包级 ok） | 需提交后触发；⑩ 授权后补齐，缺证不得关片 |
| 3 | L1/L2/L3（定向小读 Low） | 用户裁定接受不修，如实记录（§5） |
| 4 | Windows 发布路径句柄相对创建/链接 | 用户终裁：S1 独立边界登记（§8），不修代码 |

## §16 先例禁用核对

- OB-140（⑥ +1）/ OB-141（Codex 第四读）/ OB-142（R=7 批）/ OB-143（定向小读）均「不得作为先例」；
  F-01 ⑦ 额度、F-04 主控模型、四读试水、D-2 追加裁定均为用户逐项授权，不作先例；
  本报告四读矩阵仅按判级分歧记录（n=1），无能力排序表述。

## §17 下一步建议

1. ⑩ 授权提交（`git add` 精确暂存已就绪；提交信息 `<type>: <summary>` 英文）→ 推 origin（仅 origin，禁止 gitee）。
2. 提交后触发 CI 双腿 → 归档 Gates 步与 adapters/console/guard 包级 `ok` 行 → F17/F22 证据债销账。
3. 台账 `REMAINING_SLICES{,_EN}.md` P85 段与 SECURITY 注记在 CI 证据落齐后按既有流程回填；P87 治理片排片沿用本片四读数据。

## 附 A 执行自由度清单（⑧b 定稿回写，与 `docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/freedom-list.md` 一致）

| 编号 | 自由度 | 终判 |
|---|---|---|
| F-01 | ⑦ 额度 1+3（用户裁定） | in-bounds（OB 留痕） |
| F-02 | ⑦ = GPT-6.2 Sol Max 深度档 | in-bounds |
| F-03 | 盲审必抽 1 | in-bounds |
| F-04 | 主控/① 统一 V4 Pro Max | in-bounds（用户授权） |
| F-05 | ②③ Flash High 不中途改档 | in-bounds |
| F-06 | ⑤ V4 Pro 1+1 | in-bounds |
| F-07 | ⑥ ≤3 + OB-140 | in-bounds（追加已登记） |
| F-08 | GLM-5.3-Flash ≥2 串行 | in-bounds |
| F-09 | 协议先行 | in-bounds |
| F-10 | ⑨ 含 Ubuntu VM + 双平台 CI | in-bounds（CI 待提交后） |
| F-11 | 探针 0 | in-bounds |
| F-12 | 不推 gitee；英文提交信息；⑩ 停点 | in-bounds |
| F-13（追加） | 四读 + 定向小读（模型名单扩列类临时授权） | in-bounds（OB-141/143 + 试水授权，不得作先例） |
