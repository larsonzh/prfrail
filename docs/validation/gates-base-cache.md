# P96 GATES-BASE-CACHE 验证报告（跨作用域 base 结论缓存）

切片：`GATES-BASE-CACHE`（P96）。日期：2026-11-28。规模：M。类型：工具片（条款先行）。

## 1. 结论

- 结论：**完成**。缓存只覆盖 `build` / `vet` / `test` 三条 §6.1 结论；`gofmt` 恒真跑、不参与缓存。
- 边界：判据判定逻辑零改动；绝不缓存任何 G1–G9 判据结论；冷热报告逐字节一致。
- 版本：准则 **v1.78 → v1.79 → v1.80 → v1.81**（同片四段修订）。

## 2. 收益论证（实测）

| 场景 | 冷遍 | 热遍 | 说明 |
|---|---|---|---|
| `--all --scope=tree` | 105705 ms | 11484 ms | 命中一条目；整报告 sha256 相等 |

- 收益：单切片多作用域连跑可省约 **262 s**（保守口径）；成本为热遍 **+1.5–3 s**（`gofmt` 真跑与键计算）。
- 判定：收益成立；过度失效方向 fail-safe（只会多冷跑）。

## 3. 实现

- 键：`base-v2` | HEAD | porcelain `-z` 摘要 | 未跟踪内容映射 | 已改跟踪内容映射 | 忽略 `.go` 映射 | `go version` | `go env` 全量摘要（排除 `GOGCCFLAGS`）。
- `GOWORK` 非空且非 `off` ⇒ 本遍放弃缓存（仓外模块源码不可入键）。
- 只缓存 **PASS**；FAIL 永不落盘；label 与 detail 按位精确形态；存储值含 `key` 绑定。
- 任何异常 ⇒ fail-open 全量；命中路径无 hit/miss 标记。

## 4. 验收证据

| 项 | 读数 | 判定 |
|---|---|---|
| 冷热 A/B | `FULL_IDENTICAL=True` | 通过 |
| 三反例 | tracked / untracked / toolchain 均冷跑且换键 | 通过 |
| fail-open | 翻转 1 字节 ⇒ 回全量并自愈 | 通过 |
| GOWORK 非空 | `computeKey = null` | 通过 |
| `SELFTEST` | 806/806 | 通过 |
| tree / index | `TOTAL_FAIL=0` | 通过 |
| 冷热两遍 × 三作用域 | tree 125038→26392 ms、index 124263→22862 ms、ci 116902→21920 ms，`TOTAL_FAIL=0` | 通过 |
| CI 双腿 run `46053254410` | `SELFTEST: PASS 806/806` ×2、`GATE REPORT scope=ci changed=17` ×2、`TOTAL_FAIL=0` ×2 | 通过 |
| 主提交 | `81a6794` | 落地 |

## 5. 审查链

- ① `DeepSeek V4 Pro`：`DESIGN: DONE`；⑤ 预审同模型 `PRE-REVIEW: PASS`。
- ③ `DeepSeek V4.1 Flash`：r1 报 F-1、r2 报 F-2，r3 `TEST: PASS`（冷热与变异由主控自执行，③ 为事后复核）。
- ⑥ `MAI Code 2 Flash`：r1 `FINDINGS`、r2 `FINDINGS`、r3 `INDEPENDENT SCAN: PASS`。
- ⑦ 深度档 `GPT-6.2 Sol` Max：r1–r4 `FINDINGS`；Neon 交叉复核 4 项 HIGH；额度用尽后经用户裁定收敛。
- 收窄复审（额度外，OB-37）：v1.80 轮 1 HIGH（`GO111MODULE` / `go.work`）；v1.81 轮 2 HIGH（过滤过宽 / 仓外模块），均已按用户裁决处置。

## 6. 异常处理

- 自主发现 High：`runBaseGates` 的 `emit` 误把 detail 写入 label ⇒ 缓存永不写入且报告 detail 退化；④ 前就地修复，不单计 R。
- 偏离登记：冷热 A/B 与三反例变异由主控自执行、③ 事后复核（用户已背书）；独立性影响由 ⑥⑦ 与 ⑨ 兜底。
- HIGH-1：`GOGCCFLAGS` 过滤过宽会误删 `CGO_CFLAGS` 行 ⇒ 改为精确变量名过滤。
- HIGH-2：`go.work` 仓外模块源码不可入键 ⇒ 方案 (b) 保守禁用。
- 实现期发现：`go env` 含随机 `-ffile-prefix-map` ⇒ 排除 `GOGCCFLAGS`。

## 7. 已知限制

- 本地手工改写缓存条目属威胁模型外（§5.1 运行时缓存，与 `g4b` 同模型）。
- 无关环境变量变化导致多余冷跑（正确性优先，已在条款注明）。

## 8. 版本包

| 版本 | 内容 |
|---|---|
| v1.78 | §6.1 首落 base 结论缓存条款 |
| v1.79 | ⑥ 整改：键段补列、`GOTOOLCHAIN`、`key` 绑定 |
| v1.80 | `gofmt` 恒真跑、缓存收敛为三条 |
| v1.81 | `go env` 全量摘要 + `GOWORK` 保守禁用 |

## 9. 工件

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
tools/gates/gate.js	repo
tools/gates/selftest.js	repo
tools/gates/lib/baseCache.js	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/gates-base-cache.md	repo
docs/validation/gates-base-cache_EN.md	repo
docs/validation/evidence/P96-GATES-BASE-CACHE/00-briefing.md	repo
docs/validation/evidence/P96-GATES-BASE-CACHE/01-design.md	repo
docs/validation/evidence/P96-GATES-BASE-CACHE/input-list.md	repo
docs/validation/evidence/P96-GATES-BASE-CACHE/freedom-list.md	repo
docs/validation/evidence/P96-GATES-BASE-CACHE/MANIFEST.md	repo
docs/validation/evidence/P96-GATES-BASE-CACHE/SHA256SUMS.txt	repo
docs/validation/evidence/P96-GATES-BASE-CACHE/04-reviews/03-tester-review.md	repo
docs/validation/evidence/P96-GATES-BASE-CACHE/04-reviews/06-scan-r1.md	repo
docs/validation/evidence/P96-GATES-BASE-CACHE/04-reviews/06-scan-r2.md	repo
docs/validation/evidence/P96-GATES-BASE-CACHE/04-reviews/06-scan-r3.md	repo
docs/validation/evidence/P96-GATES-BASE-CACHE/04-reviews/07-deep-review-log.md	repo
docs/validation/evidence/P96-GATES-BASE-CACHE/05-batch/01-mirror-check.md	repo
docs/validation/evidence/P96-GATES-BASE-CACHE/05-batch/02-number-check.md	repo
```

## 10. ⑥ 独立扫描结论（逐字）

```sixstage
## Section A
`INDEPENDENT SCAN: PASS`

## Section B
`none`

## Section C
- M-1：**CLOSED（条款内）**。条目仅缓存 PASS；残缺、非法 detail、非法 key 均不回放，异常回全量。
- M-2：**CLOSED**。删除路径仍计入 porcelain 总预算，且状态摘要本身进入键。
- M-3：**CLOSED**。`GOTOOLCHAIN` 已在实现变量集与条款键段中。
- LOW：**CLOSED**。条款已同步列出 ignored `.go` 与 gofmt 相关表述；运行时缓存例外固定读数已同步为 2。

## Section D
可证伪点：任一四条形态可伪造记录、FAIL 被写入、201 删除未触发 null、`GOTOOLCHAIN` 变更未失效，即可反推对应项未闭合。

## Section E
未独立重跑；结论依据 r3 门禁读数、源码、差异与条款。
INDEPENDENT SCAN: PASS
```
