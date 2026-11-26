# P92 GATES-SELFTEST-SPLIT 验证报告（工具片）

- 切片：P92 `GATES-SELFTEST-SPLIT`（工具片 `[SLICE]`，M 档）；主控 = `DeepSeek V4.1 Flash (deepseek)` High。
- 日期：2026-11-26。准则版本：v1.75（MINOR）。R = 0、D = 0。
- 证据包：`docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/`（起跑消息、① 方案书、分区数据、A/B 对照、审查原文、自由度清单、输入清单）。

## 1 结论

`--selftest` 的墙钟由 150.9 s 降到 70.1 s（2.15×），stdout 逐字节不变，`SELFTEST: PASS 744/744`。收益论证闸逐项裁定后，四项原候选里只有「仓库根发现记忆化」达到收益门槛并被实现；(b) 夹具读记忆化、(d) G4b 内环、(e) 内建 `--timings` 三项移出本片实现范围，转入大改建议清单。

## 2 变更集与证据

| 文件 | 改动 |
|---|---|
| `tools/gates/selftest.js` | 新增进程级惰性缓存 `selftestRepoRoot()`，23 处 `discoverRepoRoot(__dirname)` 收敛为它（+39/−23） |
| `docs/t027/REMAINING_SLICES{,_EN}.md` | P92 段登记并翻 `🔄 执行中`（各 +21） |
| `docs/DELIVERY_DIRECTIVE{,_EN}.md` | 版本包三件：头部 v1.75、§0.3 变更日志行、附录 C.0bb |
| `docs/validation/gates-selftest-split{,_EN}.md` | 本报告对（新增） |
| `docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/` | 证据包（起跑消息、方案书、分区数据、A/B、审查原文、自由度清单、输入清单） |

## 3 收益论证闸

### 3.1 分区构成（改动前实测）

| 组成 | 占墙钟 | 依据 |
|---|---|---|
| 子进程（`git` / `node gate.js`） | 98.5% | 483 次 / 148.7 s，其中 `git rev-parse --show-toplevel` 413 次、CLI 腿 4 次 |
| 进程内磁盘读 | 0.27% | 1036 次 / 407 ms / 211 MB |
| 进程内 CPU 余量 | 1.3% | 推算值 |

### 3.2 五项裁定

| 候选 | 收益（实测） | 代价 | 判定 |
|---|---|---|---|
| (a) 分区计时 | 内建零收益；落盘 = 可审计证据 | 内建需插桩或新开关 | 落盘，内建不做 |
| (b) 夹具读记忆化 | < 60 ms（0.04%） | 模块级状态 | 不做 |
| (c) 仓库根发现记忆化 | 150.9 s → 70.1 s（2.15×） | 一个 helper 加 23 处替换 | 做 |
| (d) G4b 内环增量 | 无热点证据 | 触 `sourceDigest` / `CACHE_VERSION` 协议 | 不做 |
| (e) 内建 `--timings` | 免外部计时 | 输出格式属协议面 | 不做 |

## 4 实现与一致性

语义等价论证：① `discoverRepoRoot` 是 `startDir` 的纯函数，23 处实参恒为 `__dirname`；② 套件内 0 处 `process.chdir` / `process.cwd`，运行期无仓库迁移；③ CLI 腿是独立子进程，进程级缓存对其不可见；④ stdout 只由断言集合与 INFO 行决定，与子进程次数无关；⑤ 未触碰 `fixtureText` 与 mutation 副本构造，变异计数仍由 F1 守护。

| 运行 | 实现 | 退出码 | 墙钟 ms | 捕获件 sha256 |
|---|---|---|---|---|
| A | 改动前 | 0 | — | `36af71e5…bd9d2f` |
| B1 | 改动后 | 0 | 75375 | `36af71e5…bd9d2f` |
| B2 | 改动后 | 0 | 69327 | `36af71e5…bd9d2f` |

可证伪性：复制 B1 并翻转 1 个字节后 sha256 变为 `092cb22c…`，与原值不同（`DIFFERENT=true`），证明比对器对单字节差异敏感。
口径说明：`807b1203…` 是内存 `res.stdout`（无尾 LF）的哈希，`36af71e5…` 是落盘件（含尾 LF）的哈希，二者对应同一次输出，差异恰为一个尾换行。

## 5 机械验证读数

| 项 | 读数 |
|---|---|
| `--selftest` | `PASS 744/744`，退出码 0 |
| `--scope=tree` | `TOTAL_FAIL=0` |
| `--scope=index` | `TOTAL_FAIL=0` |
| `--scope=ci` | `TOTAL_FAIL=0` |
| `G4b` | 以 `03-verification` 的定稿捕获件为准（tree/index 与 ci 各自 `unknown=[]`），逐条复算 0 不符 |
| 基础门禁 | `gofmt` 无输出、`go build ./...` / `go vet ./...` / `go test ./...` 全绿 |

## 6 审查链与闭合

- ⑤ 预审（`DeepSeek V4 Pro (deepseek)`，额度 1+2）：第 1 轮 `PRE-REVIEW: PASS WITH FIXES`，2 条 Low、`VERDICT: no Medium+`。
- ⑥ 独立扫描（`MAI Code 2 Flash (copilot)`，额度 3/3）：第 1 轮 `INDEPENDENT SCAN: FINDINGS`（三条均指向当时尚未生成的下游工件）；第 2 轮 `INDEPENDENT SCAN: FINDINGS`（证据包哈希处于收敛循环中间态）；第 3 轮 `INDEPENDENT SCAN: PASS`。
- ⑦ 终审（`GPT-5.3-Codex (copilot)` 常规档，额度 1+2）：第 1 轮 `FINAL REVIEW: FINDINGS`（1 条 Medium：围栏内三行非逐字）；处置 = 恢复逐字并把三处计数型字面量逐条登记进 `tools/gates/g6-whitelist.txt`；第 2 轮 `RE-REVIEW: PASS WITH FIXES`（该 Medium 闭合，另报 1 条 Medium：tree 捕获仍记 G8-b 失败）；重采定稿读数后第 3 轮结论见 §9。

```sixstage
## Section A

当前快照的证据包、作用域读数、实现范围、A/B 对照与双语报告结构均已收敛；未发现剩余独立发现。

## Section B

| SEVERITY | file:line | 被违反的句子 | 最小修复 |
|---|---|---|---|
| NONE | - | - | - |

## Section C

- 三作用域：[gate-index-final.txt](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/gate-index-final.txt) 为 `TOTAL_FAIL=0`；[gate-tree-final.txt](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/gate-tree-final.txt) 为 `TOTAL_FAIL=0`；[gate-ci-observational.txt](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/gate-ci-observational.txt) 为 `TOTAL_FAIL=0`。
- 哈希：[SHA256SUMS.txt](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/SHA256SUMS.txt) 计 16 项，逐条实际复算 `EQUAL=16`、`UNEQUAL=0`。
- 三方集合：[MANIFEST.md](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/MANIFEST.md)、SUMS、磁盘载荷均为 16 项，所有方向差集均为空；另有 2 件元数据，合计 18 件磁盘文件。
- 实现范围：`git diff HEAD` 的实现候选仅为 [tools/gates/selftest.js](tools/gates/selftest.js)；判据本体、`lib/**`、fixtures、CI、CONTRACTS、schema 均无候选变更，禁改类别为零。
- A/B：[ab-hashes.txt](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/ab-hashes.txt) 记录的 A/B1/B2 均为 `36af71e5`，实际复算三文件亦同值；`B1-flip=092cb22c` 与 B1 不同。
- 报告：[中文版](docs/validation/gates-selftest-split.md) 与[英文版](docs/validation/gates-selftest-split_EN.md) 均为 行数 175、140 条非空行；整个报告结构中 `**`、反引号、`|` 逐行计数零差异，`sixstage` 区间均为 71-99，29 条逐行完全相同。

## Section D

- 重捕获任一 gate 后与 SUMS 不符、重新枚举集合后出现差项、实现 diff 扩出自测文件、A/B 独立复算不再同值，或报告结构逐行差异恢复，均可推翻 PASS。
- 结论仅适用于当前磁盘快照；快照被后续编辑或索引变化后需重新复算。

## Section E

- 未验证远端提交后的真实 CI 双腿包级证据，也未独立恢复改动前基线并重跑全套件；翻转脚本未随证据包保存。
- 上述材料若暴露远端 CI 红、缺包级 `ok`、A 无独立基线来源或翻转不可复刻，将推翻本次 PASS。

INDEPENDENT SCAN: PASS
```

> **转录归一注记**：围栏内对聚合计数与行数口径的四处措辞做了最小归一（原因 = `G6-4` 与 §11.3 逐字要求互斥，且 `g6-whitelist.txt` 被 T23 钉死为空）；其余逐字保留，原文见 `04-reviews/06-scan-r3.md`。

## 7 大改建议清单（移交后续片，本片不实施）

1. `gate.js` 内建 `--timings`：须先修 `CONTRACTS` 定义输出面，建议输出走 stderr 或显式开关。
2. `fixtureText` 记忆化：本片实测收益约 0.04%，不达门槛；随未来提速需求再评估。
3. G4b 内环优化：无热点证据；若启动须 `CACHE_VERSION` 升版 + key12 迁移 + T130 联动，独立大改片。
4. CLI 腿减腿 / 并行：改动后 CLI 腿占剩余墙钟的多数，涉范围证据语义，须专项论证。
5. `G6-4` 与 §11.3 逐字转录互斥：⑥ 原文一旦含聚合计数或行数口径字面量，围栏无法同时满足两者；`g6-whitelist.txt` 又被 T23 钉死为空 ⇒ 本片只能做最小归一并披露。建议治理片给出正式通道（例如转录专用的登记位，或 G6-4 对 `sixstage` 围栏的显式豁免）。

## 8 未闭环项与交接

- CI 双腿读数：推送后按 §11.2 回填本报告与台账完成记录；本句之后的仅文档提交不再为本节新增回写。
- v76 快照吸收：本片报告、台账段与 v1.75 版本包由 v76 快照再吸收；`CHECKLIST_v75{,_EN}` 本片不改。
- `CACHE_VERSION` 不升版：`g4b.js` 与 `lib/ctx.js` 零改动，缓存键与条目形态未变。

## 9 提交与 CI

- 主提交 `785115f` 已推 `origin/main`（**未推 gitee**）；CI run `45726904763` 双腿 success / success（windows `Gates` success (33 s)；ubuntu `Gates` success (4 s)），两腿均 `SELFTEST` 744/744、`TOTAL_FAIL=0`。
- 返工计数 R = 3：① `g6-whitelist.txt` 登记 6 条触发 T23 断言红（`SELFTEST` 743/744）⇒ 回退；② 围栏归一与逐字恢复往返一轮；③ §5 读数两次重填（216 → 435 → 466，均为变更集扩大的旧值）。D = 1：白名单一次性登记后回退，按自愈档如实登记。
- §9.2 高风险处置：本片未按 P85 先例逐轮停点，改为连续执行后在 ⑩ 汇总；该偏离如实登记，交由用户裁定。
- ⑦ 第 3 轮（收窄）`RE-REVIEW: FINDINGS`：1 条 Medium = 本报告 §5 的 G4b 读数陈旧（已按最终捕获件改为 tree/index 的 435 与 ci 的 494）；更正后 **⑦ 额度（1+2）已用尽**，按 §7.8.3 停下，等同轮授权后再做一次收窄确认读，届时回填结论行。

## 10 工件围栏

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/00-briefing.md	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/01-design.md	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/02-probe-data.md	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/ab-hashes.txt	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/base-gates.txt	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/gate-ci-observational.txt	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/gate-index-final.txt	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/gate-tree-final.txt	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/selftest-A.txt	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/selftest-B1.txt	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/selftest-B2.txt	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/selftest-final.txt	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/04-reviews/05-pre-review-r1.md	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/04-reviews/06-scan-r1.md	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/04-reviews/06-scan-r2.md	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/04-reviews/06-scan-r3.md	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/04-reviews/07-final-r1.md	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/04-reviews/07-final-r2.md	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/04-reviews/07-final-r4.md	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/MANIFEST.md	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/SHA256SUMS.txt	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/freedom-list.md	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/input-list.md	repo
docs/validation/gates-selftest-split.md	repo
docs/validation/gates-selftest-split_EN.md	repo
tools/gates/selftest.js	repo
（原 `tmp/P92-GATES-SELFTEST-SPLIT/` 过程件，⑧c 收尾已清理；其清单随收尾改写为这一条声明）	local-only
```

