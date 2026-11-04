# MANIFEST — EXIT-SEED-BUILD-EVIDENCE 证据包

切片：`EXIT-SEED-BUILD-EVIDENCE`（证据链补证类 `[SLICE]`；M 档；`docs/t027/CHECKLIST_v44.md` §四 层 3 第 7 项）。建包时点 2026-11-05（`HEAD` = `baa144d`，准则 v1.60）。

**包状态**：步骤 0 已按正式起跑消息完成（2026-11-05 起跑）；②a（Windows 侧 G-1 至 G-4 与 G-6/G-7 绑定）已落盘 `02-win-captures/`；②b（VM 侧 G-5 与真两代构建、G-6 Linux 工具链读数）已落盘 `02-linux-captures/`（2026-11-05，远程临时工作区已清理）；③ 证据链审计已落盘 `03-test-report.md`（末行 `TEST: DONE`，0 Critical / 0 High），其 6 项 Medium 已以 `claim-scope.md`（收窄断言）与 `anomalies.md`（编号化登记）处置，并由 `03-test-report-r2.md` 逐条判定闭合（6/6）、新提 6 项 Low（R2-1 至 R2-6，已就地修正；其中 R2-3 的计数更正以 `03-test-report-r2.md` 为准，`03-test-report.md` 按归档纪律保留原值不追溯改写）。⑤⑥⑦ 审查原文见 `04-reviews/`。⑧、⑨ 的产物与读数按轮追加归档。

**②b 的三角色提交映射**（与 T017 结构一致）：seed = `fc99f553c3066b24230a093d14bc17af0c6a398e`（构建 `prfrail`，作固定 seed 产物）；verifier = `18c9340c60deb30049554491f8905df51d0a3180`（构建 `prfrail-release`，即 T017 的 verifier baseline，与 candidate 不同）；candidate = `baa144df5898817e0ac7cd860928bb4a90cbd009`（构建 `prfrail`）。**偏差登记**：seed commit `fc99f55` **不含** `cmd/prfrail-release` 与 `internal/release`（`git ls-tree` 实测），故「用 seed 跑 selfhost」不可行；改为「seed 构建产物入 seed 根并被哈希锁定 + verifier 取自 T017 baseline」。

**②b 的模块获取偏差**：VM 无法访问 `proxy.golang.org` / `sum.golang.org`（IPv6 黑洞），模块经可达镜像 `goproxy.cn` 预取（`GOPROXY=https://goproxy.cn,direct`，`GOSUMDB=off`），随后**所有构建均在 `GOPROXY=off` 下完成**（构造期无网络依赖可证）；`go.sum` 哈希校验仍生效，变更见 `02b-00-env.txt` 的 note 行。

**瞬时工作区（②a）**：隔离构建与探测工作区位于 `tmp/EXIT-SEED-BUILD-EVIDENCE/work/`，捕获完成后即整目录删除——其中的隔离模块缓存含第三方 `.go` 源码，会污染 `gofmt -l .`（v44 §九 门禁坑第 12 条同类）。重放所需的逐条命令记录在 `02-win-captures/index.json` 的 `command` 字段。

**载荷集**：`MANIFEST.md` 与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（`MANIFEST` 归档条目集合 ≡ `SHA256SUMS` 条目集合 ≡ 包内载荷文件集合）。件数以 `SHA256SUMS.txt` 的哈希行为准，不设固定数。

## 步骤 0 的覆盖与纪律登记

1. **覆盖**：`00-briefing.md` 的前一版为无触发词抢跑期自撰草稿（转录的是更早一轮「核锚点并看队列」的消息）。按用户 2026-11-05 裁定，已用正式起跑消息**整份覆盖**，并同步重算 `SHA256SUMS.txt`；草稿内容不得被引用为本片事实来源。
2. **纪律登记**：抢跑事件记 `D = 1`（违反 §1.6 / §5.0；**经用户强制打断纠正**，非主控自行发现；处置 = 就地转正，未造成污染）。登记落点 = `docs/validation/exit-seed-build-evidence{,_EN}.md` 的「异常与兜底记录」章，本包 `freedom-list.md` 第三节先行留痕。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 用户消息逐字原文（含主控指令「纪律登记」与「覆盖警告」，以及正式起跑消息本体；围栏块承载，围栏外加性质说明与覆盖声明） |
| `input-list.md` | 输入包清单 + 缺件断言实测 + 前置核对记录（AT-13 口径、`s1-selfhost.md` 历史读数、初步缺口方向） |
| `freedom-list.md` | 执行自由度清单（步骤 0 起草；⑥ 后 ⑦ 前定稿）+ 不派发环节 + 纪律登记 + 升级边界 |
| `02-win-captures/02a-00-env.txt` | Windows 侧环境读数（`go version` / `go env` / `git rev-parse` / node / os） |
| `02-win-captures/02a-01-go-mod-download.txt` | 隔离模块缓存下的 `go mod download` 原始输出 |
| `02-win-captures/02a-02-build-candidate.txt` | 候选二进制 clean-room 构建（独立 GOCACHE/GOMODCACHE） |
| `02-win-captures/02a-03-build-release-tool.txt` | 发布工具二进制 clean-room 构建 |
| `02-win-captures/02a-04-hash-tree-seed.txt` | seed 树基线哈希 |
| `02-win-captures/02a-05-probe-actual.txt` | 候选进程产生 actual oracle（G-2 实际侧） |
| `02-win-captures/02a-06-selfhost-success.txt` | `selfhost` 成功路径结果 JSON（SeedHash / CandidateHash / OracleHash） |
| `02-win-captures/02a-07-neg-roots-overlap.txt` | 反例：seed 与 candidate 根重叠 ⇒ fail-close |
| `02-win-captures/02a-08-neg-oracle-ownership.txt` | 反例：expected oracle 归属候选根 ⇒ fail-close |
| `02-win-captures/02a-09-neg-seed-changed.txt` | 反例：seed 哈希不匹配 ⇒ fail-close |
| `02-win-captures/02a-10-neg-oracle-drift.txt` | 反例：actual oracle 漂移 ⇒ fail-close（G-4） |
| `02-win-captures/02a-11-seed-intact.txt` | 失败后 seed 哈希前后对照（证明失败不损 seed） |
| `02-win-captures/02a-12-candidate-listing.txt` | 候选根与 probe 工作根清单（expected 不在候选根） |
| `02-win-captures/expected.json` | expected oracle 原件（外部基准，来自 `testdata/selfhost/`） |
| `02-win-captures/actual.json` | actual oracle 原件（由候选进程产生） |
| `02-win-captures/bindings.json` | HEAD / seed 哈希 / 工具链 / 五枚夹具哈希 / 准则版本绑定（G-7） |
| `02-win-captures/index.json` | 机械索引：逐件命令、退出码、判定、字节数、CRLF 归一计数、sha256 |
| `02-linux-captures/02b-00-env.txt` | VM 环境读数（go1.22.12、内核、os-release、隔离缓存、代理口径） |
| `02-linux-captures/02b-01-clone.txt` | 部分克隆（`--filter=blob:none`） |
| `02-linux-captures/02b-02-worktree-seed.txt` | seed commit 工作树 |
| `02-linux-captures/02b-03-worktree-verifier.txt` | verifier commit 工作树 |
| `02-linux-captures/02b-04-worktree-candidate.txt` | candidate commit 工作树 |
| `02-linux-captures/02b-05-module-preload-seed.txt` | seed 工作树模块预取 |
| `02-linux-captures/02b-06-module-preload-verifier.txt` | verifier 工作树模块预取 |
| `02-linux-captures/02b-07-module-preload-candidate.txt` | candidate 工作树模块预取 |
| `02-linux-captures/02b-08-build-seed-binary.txt` | seed 二进制构建（GOPROXY=off） |
| `02-linux-captures/02b-09-seed-binary-version.txt` | seed 二进制可执行性（`version --json`） |
| `02-linux-captures/02b-10-build-verifier-tool.txt` | verifier 工具构建（GOPROXY=off） |
| `02-linux-captures/02b-11-build-candidate-binary.txt` | candidate 二进制构建（GOPROXY=off） |
| `02-linux-captures/02b-12-seed-tree-hash.txt` | seed 树基线哈希 |
| `02-linux-captures/02b-13-probe-actual.txt` | candidate 产生 actual oracle |
| `02-linux-captures/02b-14-selfhost-success.txt` | 两代 selfhost 验证成功（SeedHash / CandidateHash / OracleHash） |
| `02-linux-captures/02b-15-neg-roots-overlap.txt` | 反例：根重叠 ⇒ fail-close |
| `02-linux-captures/02b-16-neg-oracle-ownership.txt` | 反例：expected 归属候选根 ⇒ fail-close |
| `02-linux-captures/02b-17-neg-seed-changed.txt` | 反例：seed 哈希不匹配 ⇒ fail-close |
| `02-linux-captures/02b-18-neg-oracle-drift.txt` | 反例：actual 漂移 ⇒ fail-close |
| `02-linux-captures/02b-19-seed-intact.txt` | 失败后 seed 哈希前后对照 |
| `02-linux-captures/02b-20-candidate-listing.txt` | 候选根与 probe 根清单（expected 不在候选根） |
| `02-linux-captures/02b-21-bindings.json` | 三提交 / 工具链 / 内核 / 三二进制 / 五夹具哈希绑定 |
| `02-linux-captures/02b-22-verdicts.txt` | 远端逐件判定原始件（id / exit / verdict） |
| `02-linux-captures/index.json` | Linux 侧机械索引（逐件命令、退出码、判定、字节数、sha256） |
| `03-test-report.md` | ③ 校验测试报告（哈希复核 / 双侧交叉核对 / 失败路径断言 / 发现清单 / 末行 `TEST: DONE`） |
| `03-test-report-r2.md` | ③ 第 2 轮复审报告（6 项 Medium 闭合判定 / 新发现 R2-1 至 R2-6 / 末行 `TEST: DONE`） |
| `anomalies.md` | 异常与偏离编号台账（`D = 1`、`A-1` 至 `A-7`；含 D = 1 措辞矛盾的权威表述） |
| `claim-scope.md` | 证据能力边界：本片成立 S-1 至 S-7 / 不 assert N-1 至 N-7 / 未覆盖面 |
| `replay.md` | 逐条重放指令（Windows 与 Linux，含反例与漂移输入件重建） |
| `04-reviews/05-pre-review-r1.md` | ⑤ 预审第 1 轮原文（`PRE-REVIEW: PASS WITH FIXES`；1 Medium + 3 Low） |
| `04-reviews/05-pre-review-r2.md` | ⑤ 预审第 2 轮原文（整改后复审；`PASS WITH FIXES`；N-1/N-2 两项 Low 已修） |
| `04-reviews/06-independent-scan.md` | ⑥ 独立扫描原文（A–E 五段；Section B = `(no findings)`；`INDEPENDENT SCAN: PASS`） |
| `04-reviews/07-final-review-r1.md` | ⑦ 独立终审第 1 轮原文（1 Low；`FINAL REVIEW: PASS WITH FIXES`） |

**`index.json` 口径（v2）**：两份索引均为**派生元数据**——`bytes` / `sha256` 以磁盘实况重算，捕获期原始读数保留在 `captureMeta.rawBytes` / `captureMeta.crlfNormalized`。此变更闭合 `03-test-report.md` 的 F-03（三处滞后哈希 + `rawBytes` 同名异义）。

后续环节的产物按轮追加到上表并同步重算 `SHA256SUMS.txt`。

## 编码形态

本包载荷一律 UTF-8 with BOM + LF（`.md`），`SHA256SUMS.txt` 为 UTF-8 without BOM + LF。
