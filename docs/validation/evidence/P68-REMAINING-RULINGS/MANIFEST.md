# MANIFEST — P68 REMAINING-RULINGS 证据包

切片：`P68 REMAINING-RULINGS`（层 2 治理侧 · 分阶段策略封板确认 + 四项 R 裁定建议 · 评估类 `[SLICE]`，M 档）
起跑日期：2026-11-12；主控 = `DeepSeek V4.1 Flash (deepseek)`（默认，未登记 ADR）

## 归档条目清单

| 路径 | 原始字节 sha256 | 归档字节 sha256 | 说明 |
|---|---|---|---|
| `00-briefing.md` | `4db0e79d58cf53bc7b70524a47c67473dd57a5c1d187ac5105c752a7b7a23a38` | `4db0e79d58cf53bc7b70524a47c67473dd57a5c1d187ac5105c752a7b7a23a38` | 起跑消息逐字落盘（步骤 0） |
| `01-arch/design-r1.md` | `86149733de93bda0ac3bf23689a56174f52771011d0a09e545d255a6da336854` | `86149733de93bda0ac3bf23689a56174f52771011d0a09e545d255a6da336854` | ① 五件评估（A 封板确认 + B–E 四项 R 建议；结论行 DESIGN: DONE） |
| `01-arch/r07-refined.md` | `d597aaa82ad465d751c2ac2441c19bd8d8c4cfb9c531afaec67de8436be512da` | `d597aaa82ad465d751c2ac2441c19bd8d8c4cfb9c531afaec67de8436be512da` | R-07 完整建议（技术前提 A/B/C/D 核实 + 三方案 + 四组合 + 推荐） |
| `02-batch/glm-batch1.txt` | `3f06fb001ccd58e7a2d4a67781c3f12f808dfc0232445ee5d860b13cebd55727` | `3f06fb001ccd58e7a2d4a67781c3f12f808dfc0232445ee5d860b13cebd55727` | 批量辅助 ① 四项 R 现状锚点枚举（实得 2418 字符） |
| `02-batch/glm-batch2.txt` | `eda69afc54df01c938f6b680c090256891f6d798bd2851a791ef3e4d3ac3d63d` | `eda69afc54df01c938f6b680c090256891f6d798bd2851a791ef3e4d3ac3d63d` | 批量辅助 ② 候选方案与影响维度枚举（实得 1428 字符） |
| `03-test/consistency-r1.md` | `6d2d2d2c51af24fecb0cd882b24b661a4952069dc2a75b92fdd3fb06e3cddb5c` | `6d2d2d2c51af24fecb0cd882b24b661a4952069dc2a75b92fdd3fb06e3cddb5c` | ③ 一致性检查（34 锚点机检；1 处 Low 行号偏差修正） |
| `04-reviews/05-prereview-r1.md` | `d2d40957533a35a43df0d23d060d56a8319eb92e419e049d9896ab16e3192c9e` | `d2d40957533a35a43df0d23d060d56a8319eb92e419e049d9896ab16e3192c9e` | ⑤ 预审 r1 逐字归档（PASS；1 项 Low 已采纳整改） |
| `04-reviews/07-final-r1.md` | `2a2d174b6db6dd0f18a0f7129715bf376833bd2baaf50974f8a5c0e7c1d158c9` | `2a2d174b6db6dd0f18a0f7129715bf376833bd2baaf50974f8a5c0e7c1d158c9` | ⑦ 终审 r1 归档（FINAL REVIEW: PASS；0 findings） |
| `06-gates/gates-r1.txt` | `bb556b8ac48dbae2441f2c1994c32c44c08de940b2b472dc5899e24175cb31b9` | `bb556b8ac48dbae2441f2c1994c32c44c08de940b2b472dc5899e24175cb31b9` | ④ 门禁原始读数（selftest + tree/index/ci，TOTAL_FAIL=0） |
| `06-gates/gates-r2.txt` | `6304625c1c7990f3bcf88f83ce52ff56635d2d42dfec1057e9f1142827265498` | `6304625c1c7990f3bcf88f83ce52ff56635d2d42dfec1057e9f1142827265498` | ⑨ 终态门禁原始读数（SELFTEST 737/737 + 三作用域 TOTAL_FAIL=0） |
| `06-gates/step0-selftest-r1.txt` | `3a5d1a7c18d224991cb1e4c5130542fe3db5dca535d225d6b4b0811876b04fd7` | `3a5d1a7c18d224991cb1e4c5130542fe3db5dca535d225d6b4b0811876b04fd7` | 步骤 0 SELFTEST 锚点原始行（737/737） |
| `freedom-list.md` | `b659107bf7bc667f2fb90a51b647d5e62983ea125164cd8fe3773c544625443b` | `b659107bf7bc667f2fb90a51b647d5e62983ea125164cd8fe3773c544625443b` | 执行自由度清单（§5.1：A 自主 / B 停升级 / C 硬禁止） |
| `input-list.md` | `0116b8f4f76370470570e67993903cf9f35d4722b2de15aa022801e07fcebb3b` | `0116b8f4f76370470570e67993903cf9f35d4722b2de15aa022801e07fcebb3b` | 输入文档清单（步骤 0 前置核对，file:line 级） |

## 编码约束

- 全部 `.md` / `.txt`：UTF-8 **with BOM** + **LF**。
- 载荷任何改动后须重算 `SHA256SUMS.txt` 并同步本清单状态列。
