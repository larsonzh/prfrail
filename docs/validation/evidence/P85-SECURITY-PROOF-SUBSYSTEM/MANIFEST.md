# MANIFEST — P85-SECURITY-PROOF-SUBSYSTEM 证据包

切片：`SECURITY-PROOF-SUBSYSTEM`（P85，实现类 `[SLICE]`，L 档）。**进行中**（2026-11-22 起跑；⑧c 前只增不删）。

**载荷集**：本清单与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（磁盘 / 本清单 / `SHA256SUMS.txt`）。
**重算口径**：新增载荷后按磁盘派生重建 `SHA256SUMS.txt`，路径相对本目录。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 起跑消息逐字原文（§13.1.10） |
| `freedom-list.md` | 执行自由度清单（F-01~F-12） |
| `input-list.md` | 输入清单（步骤 0 落盘） |
| `P87-four-read-data/07-three-read-comparison.md` | 四读同输入对照表（2026-11-23，R6 批同输入；P87 治理片原始数据） |
| `P87-four-read-data/07-carbon-main-read.md` | Gemini 4 Carbon 主读原文 |
| `P87-four-read-data/07-neon-cross-read.md` | Gemini 4 Neon 交叉复核原文 |
| `P87-four-read-data/07-sol-baseline-read.md` | GPT-6.2 Sol 基准读原文 |
| `P87-four-read-data/07-codex-baseline-read.md` | GPT-5.3-Codex 法定基线盲读原文 |
| `P87-four-read-data/07-re-review-r3.md` | r3 复审原文（四读共同输入之一） |
| `P87-four-read-data/diff.patch` | R6 批全量 diff（四读共同输入之一，776KB） |
| `P87-four-read-data/numstat.txt` | R6 批 numstat 读数（四读共同输入之一） |
| `P87-four-read-data/gate-index.txt` | R6 批 gate-index 读数（四读共同输入之一） |
| `P87-four-read-data/gate-tree.txt` | R6 批 gate-tree 读数（四读共同输入之一） |
| `P87-four-read-data/gate-ci.txt` | R6 批 gate-ci 读数（四读共同输入之一） |
| `P87-four-read-data/selftest.txt` | R6 批 SELFTEST 读数（四读共同输入之一） |
| `P87-four-read-data/vm-r6-legs.txt` | R6 批 Ubuntu VM 腿读数（四读共同输入之一） |
| `P87-four-read-data/test-r6.txt` | R6 批 Windows 测试读数（四读共同输入之一） |
| `P87-four-read-data/mutation-r6.txt` | R6 批变异读数（四读共同输入之一） |
| `report.md` | P85 切片验证报告（⑧a 初稿 = ⑧b 定稿；含 §8 S1 边界独立登记段与自由度清单回写） |
| `07-carbon-r7-targeted-read.md` | R=7 定向小读原文（Gemini 4 Carbon High，⑧c 输入，OB-143） |
| `09-l4-diff-excerpt.md` | L4 修后 diff 摘录（审查环后修改、依据 P26 裁定 D-2；供顾问实质核验） |
| `local-full-log.txt` | Windows 本机全量 build/test/vet 日志与退出码（⑨ 证据） |
| `vm-full-log.txt` | Ubuntu VM 非 root 全量 build/test/vet 日志与退出码（⑨ 证据） |
| `local-r7-legs.txt` | Windows 非提升 R7 新用例逐腿 PASS/SKIP 读数（⑨ 正对照） |
| `vm-r7-legs.txt` | Ubuntu VM 非 root R7 新用例逐腿 PASS/SKIP 读数（⑨ 正对照） |
| `mutation-r7-executed.txt` | MUT-4A/4B + MUT-5 红→SHA-256 逐字节复原→绿证据（⑨） |
| `mutation-r6.txt` | 变异记录（含 R7 补记 MUT-4/MUT-5 执行态） |
| `old-wording-grep.txt` | ⑧a 自查③ 旧口径 grep 命令与结果 |
