# input-list — `GOVERNANCE-BASELINE-VERIFY`（P47）输入包清单

口径：OB-93（逐件字节数，不复制入包者仅登记路径与字节数）。采集时点 2026-11-02，`HEAD` = `0e46bd0`。

## 准则与快照

| 路径 | 字节 | 用途 |
|---|---|---|
| `docs/DELIVERY_DIRECTIVE.md` | 446616 | §12.4 台账、§1.7.4 独立登记、附录 C.0–C.0x 治理记录 |
| `docs/t027/CHECKLIST_v36.md` | 21439 | 跨会话快照（锚点、四层队列、纪律 S / 评估类停点纪律） |
| `docs/t027/REMAINING_SLICES.md` | 271187 | 切片台账（P32/P33/P34 状态行与详情块） |

## 断点定位

| 路径 | 字节 | 用途 |
|---|---|---|
| `docs/validation/strategic-review-2026-11-02.md` | 13312 | P46 报告：断点 A/B 与七项「未闭合（证据不足）」 |
| `docs/validation/strategic-review-2026-10-30.md` | 5606 | P35 报告：七项发现的**原始定义**（P3/P5/G4/G5/G7/M6/S7） |

## 历史记录（G4/G5/G7 指向的 P32/P33/P34）

| 路径 | 字节 | 用途 |
|---|---|---|
| `docs/validation/directive-d-mechanical-guard.md` | 28757 | P32 报告（R 计数争议对象） |
| `docs/validation/directive-v21-obs-sweep.md` | 29137 | P33 报告（独立登记争议对象） |
| `docs/validation/directive-tmp-self-heal.md` | 21067 | P34 报告（特例排除域争议对象） |

## 权威正文（P3 指向的语义分叉面）

| 路径 | 字节 | 用途 |
|---|---|---|
| `docs/PRODUCT_REQUIREMENTS.md` | 16087 | REQ-023 / AT-24 `unknown` 语义 |
| `docs/TEST_STRATEGY.md` | 11214 | 接受判据（与 REQ-023 的分叉面） |

## 未复制入包

`.yml`、Go 源码与 CI 工件不复制；如需引用仅登记路径与行号。
