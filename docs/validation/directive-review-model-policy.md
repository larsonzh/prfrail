# P87 DIRECTIVE-REVIEW-MODEL-POLICY 验证报告（治理片）

- 切片：P87 `DIRECTIVE-REVIEW-MODEL-POLICY`（治理片 `[SLICE]`，M 档）；主控 = `DeepSeek V4.1 Flash (deepseek)` High。
- 日期：2026-11-24。准则版本：v1.72（MINOR）。ADR-048。R = 0、D = 0。
- 证据包：`docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/`（起跑消息、① 方案书、读数、审查原文、自由度清单、输入清单）。

## 1. 结论

条款落地完成，审查链与门禁全绿；**当前状态 = 待 ⑩ 授权提交**（授权后推送 origin、走 CI 双腿、⑩ 回执）。

## 2. 变更集与证据

- 准则对（`docs/DELIVERY_DIRECTIVE{,_EN}.md`）：§2.2 三行 Gemini 终审角色 + Sol 备选句；§5.2 四档选型矩阵（含平行备选档）+ 判级门槛共识 + 组合触发条款 + 适用域排除句 + 生效机制；新增 §5.4；§7.5 单价三行 + 三注记；§8.2 档位三行 + 取档授权；§6.3 G5-b 闭集 11 枚；头部 v1.72；§0.3 的 v1.72 变更日志行；附录 C.0ay。
- ADR 对：ADR-048（白名单扩列与选型矩阵，用户直接裁定）。
- 台账对：P87 段（步骤 0 登记 + 翻「执行中」）。
- 判据实现：`tools/gates/lib/criteria/g5b.js`（`CLOSED_SET` 11 枚，与 §2.2 / §6.3 逐字同序）+ `tools/gates/selftest.js`（T132a–c）+ 4 个夹具（3 合法 + 1 非法）。

## 3. 机械验证读数

- `gate.js --selftest`：`SELFTEST: PASS 744/744`（起跑基线 740/740，+4）。
- 三作用域（`--skip-base`）：tree / index / ci 均 `TOTAL_FAIL=0`（读数见 `03-verification/gate-*.txt`）。
- `go build ./...` / `go vet ./...` / `go test ./...`：均退出 0（15 包 ok）。
- 六形状（CN/EN 准则对）：行数 2540、加粗 6160、竖线 4287、空行 589、反引号 6838、标题行号逐一同位 —— 全等。
- `git diff --cached --numstat`：准则对、ADR 对、台账对三对逐对称。- CI：run \`45294006169\`（commit \`8b64ce3\`）双腿 success（windows-latest / ubuntu-latest）；其余三个 job 按 CI 定义 skipped。
## 4. 审查链与闭合

- ⑤ 预审（`DeepSeek V4 Pro`，额度 1 + 1）：r1 `PASS WITH FIXES`（6 Medium）→ 修复 → r2 收窄复审报 1 Medium（深度档缺 Neon 触发条件）+ 4 Low。**⑤ 额度用尽**，该 Medium 已按起跑消息 ①「目标」字段与原 §5.4 文本逐字落盘（`§5.2` 深度档行），其闭合交由 ⑥ 与 ⑦ 复核；3 条 Low 一并处置（生效机制尾句、§7.5 措辞、轻量档措辞）。
- ⑥ 独立扫描（`MAI Code 2 Flash`，额度 ≤ 3，用 2）：r1 `FINDINGS`（1 Medium：EN 轻量档措辞不等价；2 Low）→ 修复 → r2 `PASS`。
- ⑦ 盲审（`GPT-5.3-Codex` 常规档，另计 +1）：`FINDINGS`（1 Medium：备选额度口径与起跑消息双口径）→ 修复为「与深度档同口径、按切片档行取值」→ 终审判定闭合。
- ⑦ 终审（`GPT-5.3-Codex` 常规档，额度 1 + 3，用 1 终审）：`FINAL REVIEW: PASS`，逐条判定 ⑤⑥ 与盲审全部发现闭合，未发现新增 Medium+。

## 5. 计量

- ⑤ × 2（r1 / r2）；⑥ × 2（r1 / r2）；⑦ = 盲审 1 + 终审 1；批量辅助 `GLM-5.3-Flash` × 2（镜像语言侧检查、条款枚举核对）。
- 探针 0；未调用 Hydrogen；未附加 Neon 交叉复核（触发式口径，本片未触发）。
- 计量的单一真值源 = 会话调用回执；无额度例外。

## 6. 异常处理

- **编码归一误写（步骤 0）**：首次归一脚本以相对路径读取失败后误写，五个证据件被截为 3 字节（仅 BOM）。同轮重写并以「长度断言 > 1000 才写回」的脚本重做归一；无仓库外与 `tmp/` 根残留，不计 D。已登记于 `freedom-list.md` 与 `input-list.md`。
- **⑥ r1 覆盖面偏差**：r1 报告称 build/vet「0 行缺退出码」，已补 `EXIT=0` 行，r2 判定闭合。
- **⑤ 额度用尽后的闭合路径**：r2 的 1 Medium 未再追加 ⑤ 轮次（不自费加轮，§7.8.3），改由 ⑥ r2 与 ⑦ 终审复核闭合。
- 临时件 D 计数：0（`tmp/P87-DIRECTIVE-REVIEW-MODEL-POLICY/` 内脚本于 ⑧c 清理）。

## 7. 未闭环项与交接

- **Hydrogen 对照评审前置未完成**：轻量档（Gemini 4 Hydrogen）条款已落地、对照评审实测未完成，属 P87 已知前置缺口，待后续业务片按 §1.7.4 临时授权补足；补足后回填 P87 报告注记。该项**不得表述为已验证**（本片未调用 Hydrogen）。
- **⑦ 复审额度剩余 3**：本片终审一次通过，无复审。
- **CI 双腿已取证**：run \`45294006169\`（commit \`8b64ce3\`）双腿 success，已回填至 §3；后续 push 产生的 CI run 仅作观察记录。
- **v71 快照吸收**：本片条款与台账段由 v71 快照再吸收（本片不改 `CHECKLIST_v70{,_EN}`）。

## 8. ⑥ 输出转录（逐字）

```sixstage
### Section A
覆盖：r1 三条闭合路径、双语 §5.2/§7.5、修复摘要、原始 diff、六形状与门禁读数。

未覆盖见 Section E。

### Section B
| # | SEVERITY | 位置 | 发现 | 建议 |
|---|---|---|---|---|
本轮无残余发现行，也未发现新增 Medium+。

### Section C
1. **Medium — 已闭合**：CN 的轻量档已写为“定稿复核与定向小读”，EN 同位写为“finalization reviews and targeted small reads”，不再等同于 `final-review passes`；见 [CN](docs/DELIVERY_DIRECTIVE.md#L574) 与 [EN](docs/DELIVERY_DIRECTIVE_EN.md#L574)。
2. **Low — 已闭合**：§7.5 已明确 Sol 1M 档输入为上表两倍、输出为上表 1.5 倍；见 [CN](docs/DELIVERY_DIRECTIVE.md#L830) 与 [EN](docs/DELIVERY_DIRECTIVE_EN.md#L830)。
3. **Low — 已闭合**：`build.txt` 与 `vet.txt` 均落盘 `EXIT=0`；见 [build.txt](docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/build.txt#L2) 与 [vet.txt](docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/vet.txt#L2)。

### Section D
未发现新问题：

- 双语 §5.2 和 §7.5 的修订口径一致；[shape-six.txt](docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/shape-six.txt#L1) 六形状全等。
- G5-b 三枚 Gemini 在 CN、EN 与实现闭集一致；见 [CN](docs/DELIVERY_DIRECTIVE.md#L680)、[EN](docs/DELIVERY_DIRECTIVE_EN.md#L680)、[g5b.js](tools/gates/lib/criteria/g5b.js#L34)。
- [gate-tree.txt](docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/gate-tree.txt#L48)、[gate-index.txt](docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/gate-index.txt#L48)、[gate-ci.txt](docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/gate-ci.txt#L45) 均为 `TOTAL_FAIL=0`；[gate-selftest.txt](docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/gate-selftest.txt#L749) 为 `SELFTEST: PASS 744/744`。

### Section E
无法验证：落盘材料之外的 GitHub live CI、真实外部模型调用、Operator UI 实时定价，以及 `build.txt`/`vet.txt` 缺失的完整构建与 vet 原始日志。

INDEPENDENT SCAN: PASS
```

## 9. 工件围栏

```artifacts
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-review-model-policy.md	repo
docs/validation/directive-review-model-policy_EN.md	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/00-briefing.md	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/01-arch/01-arch-design.md	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/build.txt	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/gate-ci.txt	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/gate-index.txt	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/gate-selftest.txt	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/gate-tree.txt	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/numstat.txt	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/shape-six.txt	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/test-all.txt	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/vet.txt	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/04-reviews/05-pre-review-r1.md	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/04-reviews/05-pre-review-r2.md	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/04-reviews/06-scan-r1.md	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/04-reviews/06-scan-r2.md	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/04-reviews/07-blind-review.md	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/04-reviews/07-final-review-r1.md	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/MANIFEST.md	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/SHA256SUMS.txt	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/freedom-list.md	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/input-list.md	repo
tools/gates/lib/criteria/g5b.js	repo
tools/gates/selftest.js	repo
tools/gates/testdata/fixtures/g5b-closed-set-legal-carbon.txt	repo
tools/gates/testdata/fixtures/g5b-closed-set-legal-hydrogen.txt	repo
tools/gates/testdata/fixtures/g5b-closed-set-legal-neon.txt	repo
tools/gates/testdata/fixtures/g5b-illegal-gemini-unknown.txt	repo
```

SLICE: READY_FOR_REVIEW
