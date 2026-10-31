# MANIFEST — P44-T027B-GATE-READINESS 证据包

切片：`P44-T027B-GATE-READINESS`（P44，实现类 `[SLICE]`，C7+N4 门控就绪与一次性许可）。建包时点 2026-11-01（`HEAD` = `a87b078`，准则 v1.60）。

**包状态**：建包后按轮追加 01 / 05 / 06 / 07 四批载荷；⑧ 收口时归档 `07-final/review.md` 并以一次性脚本重算 `SHA256SUMS.txt`，本包载荷集至此闭合。

**载荷集**：`MANIFEST.md` 与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（件数以 `SHA256SUMS.txt` 的哈希行为准，不设固定数）。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 用户起跑消息原文（剥离附件清单；v33 §七 第 10 条「起跑消息即工件」） |
| `freedom-list.md` | 执行自由度清单（步骤 0 落盘，含角色载体铁律） |
| `01-arch-design/review.md` | ① 架构《C7+N4 门控就绪与一次性许可方案》（`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`）原文归档 + §八 四条修订注记（R-1～R-4） |
| `03-tests/mutation-evidence.txt` | ③ 变异验证原始输出（11 项变异 M1–M11，全部 RED；证明断言可证伪） |
| `05-prereview/review.md` | ⑤ 预审 r1 + r2 输出原样（`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`） |
| `06-scan/review.md` | ⑥ 独立扫描输出原样（Section A–E；`independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`） |
| `07-final/review.md` | ⑦ 终审 r3 输出原样（`independent-reviewer` + `GPT-5.3-Codex (copilot)`；r2 两项发现闭合核验 + 放行条件 5 条） |

## 过程偏差登记

| 编号 | 事实 | 影响面 | 为何不影响结论 | 补救 |
|---|---|---|---|---|
| D9 | ② 载体执行清理时误用整目录递归删除，令本片 `tmp/` 中间件（批次普查、④ 中间门禁读数等）丢失 | 仅 `tmp/` 下未入库的一次性文件；持久证据包与源码零影响 | `tmp/` 全程 gitignored，最终三作用域门禁读数在 ⑨ 重取；批次普查产出已被 ① 消费并归档 | 后续清理只允许删本载体子路径；本文件登记留痕 |
| D10 | ③ 首轮变异验证输出未落盘 | 可复算性降低 | 已在 `03-tests/mutation-evidence.txt` 复算重跑并落盘（11 项变异全 RED） | 已补救 |

## 未入库文件清单

无。

## 重建说明

载荷为 UTF-8 + BOM + LF 纯文本（`SHA256SUMS.txt` 无 BOM）；由 `tmp/P44-T027B-GATE-READINESS/` 一次性脚本生成，脚本按 §5.1 在 ⑧c 清除。
