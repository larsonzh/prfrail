# MANIFEST — CONTRACTS-S7-DURABILITY 证据包

切片：`CONTRACTS-S7-DURABILITY`（契约与实现类 `[SLICE]`；L 档；v43 §四 层 3 第 6 项，R-02 裁定已批准）。建包时点 2026-11-04（`HEAD` = `2b8ac3e`，准则 v1.60）。

**包状态**：步骤 0 骨架落盘；① 架构（双打：V4 Pro + Sol 红队）、③ 红基线、② 实施、⑤、⑥、⑦、⑨ 门禁读数按轮追加归档。

**载荷集**：`MANIFEST.md` 与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（件数以 `SHA256SUMS.txt` 的哈希行为准，不设固定数）。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 用户起跑消息原文（含 R-02 裁定、① 双打红队授权、模型绑定、纪律 V 拆步要求、红测先行红线与文档同步要求） |
| `01-arch/design-r2.md` | ① 架构第二轮设计原文（V4 Pro 只读载体，逐字归档；Q1 = B 计划冻结后预建，三个 High 逐条闭合，七条边界映射，预建全集持久化结论） |
| `01-arch/design-r2b.md` | ① 架构第三轮修订原文（V4 Pro 只读载体，逐字归档；按用户裁定处置第二次红队 1 High + 2 Medium + 3 Low） |
| `01-arch/design-r2c.md` | ① 架构微修原文（V4 Pro 只读载体，逐字归档；激活语义改 store 级 + A2 激活语义句） |
| `01-arch/design.md` | ① 架构设计原文（V4 Pro 只读载体，逐字归档；含 §7 修订文本、C6 函数级落地设计、零影响证明、核心声明论证与三判据对应） |
| `01-arch/red-team-sol-review-r2.md` | 第二次 Sol 红队审查原文（design-r2 收窄复核；§1.7.4 临时授权 ADR-029；逐字归档） |
| `01-arch/red-team-sol-review.md` | ① 双打第二棒：GPT-6.1 Sol 红队抗压审查原文（§1.7.4 临时授权，逐字归档） |
| `01-arch/ruling-2026-11-04.md` | 用户裁定原文（2026-11-04 目标边界裁定 + §7 修订边界 + 临时停点；逐字归档） |
| `01-arch/ruling-2026-11-04-q1.md` | 用户 Q1 裁定原文（2026-11-04 Q1 = B 计划冻结后预建 + B 实现边界七条 + 第二次红队预授权；逐字归档） |
| `05-prereview/review-r1.md` | ⑤ 预审 r1 原文（PRE-REVIEW: PASS WITH FIXES，0 High / 0 Medium / 4 Low；逐字归档） |
| `06-scan/review.md` | ⑥ 独立扫描原文（INDEPENDENT SCAN: PASS，A–E 形状；逐字归档） |
| `07-final/review.md` | ⑦ 终审原文（Codex `FINAL REVIEW: PASS`，0C/0H/0M/0L；等价性论证 §7.2.1 口径通过）+ ⑦ 复审 #1/3 原文（②e 收窄输入，Codex `FINAL REVIEW: PASS`，0 发现；两轮逐字归档） |
| `08-gates/cross-platform-linux.txt` | ④ 跨平台编译证据（GOOS=linux build/vet，exit 0/0） |
| `08-gates/gate-readings.txt` | ⑨ 门禁读数（三作用域 TOTAL_FAIL=0 + SELFTEST PASS 697/697） |
| `08-gates/green-c6slot-tests.txt` | ③ 绿侧证据（`go test -tags c6slot -count=1 ./internal/adapters/` 通过） |
| `08-gates/red-baseline-c6slot-compile.txt` | ③ 红侧证据（HEAD=2b8ac3e + 最终测试文件的编译级红复现；原始 ③a 捕获存于 tmp/ 已随 ⑧c 清除，⑧b 以确定性复现归档，red-ness 相同） |
| `freedom-list.md` | 执行自由度清单（步骤 0 落盘，含载体铁律、必须升级事项、接口签名冻结、红测先行红线与拆步约束） |
| `input-list.md` | 输入包清单（OB-93 口径，逐件字节数与前置核对要点） |

后续环节（①/③/②/⑤/⑥/⑦/⑨）的产物按轮追加到上表并同步重算 `SHA256SUMS.txt`（步骤 0 只登记当前已落盘载荷）。
