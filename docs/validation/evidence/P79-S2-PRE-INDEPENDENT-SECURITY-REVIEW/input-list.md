# 评审输入清单（P79；逐件申报，文件名 + 行数 + 字节数）

生成方式：脚本按磁盘实测生成（`tmp/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/gen-input-list.js`），非手工记数。
口径：行数 = 按 LF 切分的行数（与 `gate.js` 的文本口径一致）；字节数 = 磁盘原始字节数（含 BOM）。
用途：⑤ / ⑥ / ⑦ 的审查输入包逐件断言来源（§5.1 的时序说明，OB-93）。

## A. 安全基线与需求

| 序号 | 文件 | 行数 | 字节数 | 存在 |
|---|---|---|---|---|
| 1 | `docs/SECURITY.md` | 76 | 10861 | 是 |
| 2 | `docs/SECURITY_EN.md` | 76 | 12146 | 是 |
| 3 | `docs/PRODUCT_REQUIREMENTS.md` | 132 | 16667 | 是 |
| 4 | `docs/PRODUCT_REQUIREMENTS_EN.md` | 130 | 18395 | 是 |
| 5 | `docs/ADR_REGISTER.md` | 117 | 65499 | 是 |
| 6 | `docs/ADR_REGISTER_EN.md` | 118 | 77732 | 是 |
| 7 | `docs/DEV_PLAN.md` | 265 | 86644 | 是 |
| 8 | `docs/DEV_PLAN_EN.md` | 266 | 99231 | 是 |
| 9 | `docs/TEST_STRATEGY.md` | 90 | 13560 | 是 |
| 10 | `docs/TEST_STRATEGY_EN.md` | 90 | 14808 | 是 |

## B. 执行准则与快照

| 序号 | 文件 | 行数 | 字节数 | 存在 |
|---|---|---|---|---|
| 1 | `docs/DELIVERY_DIRECTIVE.md` | 2474 | 506266 | 是 |
| 2 | `docs/DELIVERY_DIRECTIVE_EN.md` | 2474 | 617890 | 是 |
| 3 | `docs/t027/CHECKLIST_v65.md` | 221 | 28992 | 是 |
| 4 | `docs/t027/CHECKLIST_v65_EN.md` | 221 | 33671 | 是 |
| 5 | `docs/t027/REMAINING_SLICES.md` | 2448 | 381800 | 是 |
| 6 | `docs/t027/REMAINING_SLICES_EN.md` | 2435 | 428686 | 是 |

## C. 门禁与判据来源（硬安全门禁面证据）

| 序号 | 文件 | 行数 | 字节数 | 存在 |
|---|---|---|---|---|
| 1 | `docs/validation/sg6-gatelevel-consolidate.md` | 263 | 21612 | 是 |
| 2 | `docs/validation/s1-exit.md` | 77 | 10221 | 是 |
| 3 | `docs/validation/s1-exit_EN.md` | 77 | 11792 | 是 |
| 4 | `tools/gates/gate.js` | 295 | 11126 | 是 |
| 5 | `tools/gates/lib/criteria/g1.js` | 259 | 11994 | 是 |
| 6 | `tools/gates/lib/criteria/g5a.js` | 67 | 2240 | 是 |
| 7 | `tools/gates/lib/criteria/g7.js` | 1060 | 46219 | 是 |
| 8 | `tools/gates/lib/criteria/g9.js` | 205 | 8737 | 是 |

## D. 产品侧硬安全门禁实现

| 序号 | 文件 | 行数 | 字节数 | 存在 |
|---|---|---|---|---|
| 1 | `internal/gates/policy.go` | 263 | 8782 | 是 |
| 2 | `internal/gates/runner.go` | 399 | 14238 | 是 |
| 3 | `internal/gates/readiness.go` | 238 | 9321 | 是 |
| 4 | `internal/gates/probe_permit.go` | 319 | 12926 | 是 |
| 5 | `internal/guard/lease.go` | 331 | 11563 | 是 |
| 6 | `internal/guard/sandbox.go` | 115 | 4081 | 是 |
| 7 | `internal/guard/sandbox_unsupported.go` | 20 | 574 | 是 |
| 8 | `internal/chain/postflight.go` | 464 | 19294 | 是 |
| 9 | `internal/evidence/verify.go` | 120 | 4188 | 是 |
| 10 | `internal/evidence/disposition.go` | 717 | 26896 | 是 |
| 11 | `internal/evidence/event.go` | 271 | 9993 | 是 |

## E. 秘密与凭据面实现

| 序号 | 文件 | 行数 | 字节数 | 存在 |
|---|---|---|---|---|
| 1 | `internal/adapters/ai_secret_windows.go` | 165 | 4890 | 是 |
| 2 | `internal/adapters/ai_provider.go` | 122 | 3876 | 是 |
| 3 | `internal/adapters/ai_probe_copilot.go` | 222 | 8529 | 是 |
| 4 | `internal/console/ai.go` | 259 | 10466 | 是 |
| 5 | `docs/OPERATIONS.md` | 148 | 18308 | 是 |

合计：申报文件 40 件、字节合计 2704714（由本脚本实测，未手工记数）。
