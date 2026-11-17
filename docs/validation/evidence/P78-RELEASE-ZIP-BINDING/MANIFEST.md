# MANIFEST — P78-RELEASE-ZIP-BINDING 证据包

切片：`RELEASE-ZIP-BINDING`（P78，v0.1.0 Release 创建 + 最终 ZIP 构建 + 发行记录绑定 + 安装生命周期复验，实现类 `[SLICE]`，M-L 档）。建包时点 2026-11-17（`HEAD` = `c935713`，准则 v1.69）。

**包状态**：步骤 0 建包（三件载荷：`00-briefing.md` / `input-list.md` / `freedom-list.md`）；后续按轮追加 `01-arch/` / `02-batch/` / `04-reviews/` / `06-gates/` 载荷；⑧c 收口时归档定稿件并重算 `SHA256SUMS.txt`。

**载荷集**：本清单与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等。两者不在下表内（避免与载荷集混计）。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 用户起跑消息原文（§13.1.10；2026-11-17 P78 起跑，M-L 实现类） |
| `01-arch/design-r1.md` | ① 阶段方案书原文（`deep-reasoner` + `DeepSeek V4 Pro`；Q1–Q6） |
| `02-batch/batch-01-release-pipeline.md` | 批量辅助 1 事实枚举（`GLM-5.3-Flash`）：ci.yml job/inputs 与发布链路 |
| `02-batch/batch-02-release-docs.md` | 批量辅助 2 事实枚举（`GLM-5.3-Flash`）：八组文档的发布相关表述 |
| `04-reviews/final-r1.md` | （后续轮次追加载荷） |
| `04-reviews/final-r2.md` | （后续轮次追加载荷） |
| `04-reviews/prereview-r1.md` | （后续轮次追加载荷） |
| `04-reviews/prereview-r2.md` | （后续轮次追加载荷） |
| `04-reviews/scan-r1.md` | （后续轮次追加载荷） |
| `06-gates/base-go.txt` | （后续轮次追加载荷） |
| `06-gates/gate-ci.txt` | （后续轮次追加载荷） |
| `06-gates/gate-index.txt` | （后续轮次追加载荷） |
| `06-gates/gate-selftest.txt` | （后续轮次追加载荷） |
| `06-gates/gate-tree.txt` | （后续轮次追加载荷） |
| `freedom-list.md` | 执行自由度清单（步骤 0 落盘初稿；⑥ 之后、⑦ 之前定稿，§1.7.3） |
| `input-list.md` | 本片输入清单（步骤 0 落盘，§13.1.1） |

## 后续追加计划

- 按 ① / 批量辅助 / ⑤ / ⑥ / ⑦ / ⑨ 各轮追加；⑧c 只重算 `SHA256SUMS.txt` 并核验五条读数。

## 编码

全部 `.md` 为 UTF-8 with BOM + LF（仓库约定）；`SHA256SUMS.txt` 为无 BOM + LF。
