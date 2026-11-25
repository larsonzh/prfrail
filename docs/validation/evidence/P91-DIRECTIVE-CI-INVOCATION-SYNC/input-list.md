# input-list — P91 `DIRECTIVE-CI-INVOCATION-SYNC`

## 步骤 0 锚点实读（2026-11-26）

| 项 | 起跑消息所列锚点 | 实读 | 结论 |
|---|---|---|---|
| `git rev-parse --short HEAD` | `254e98f` | `254e98f` | 一致 |
| `git rev-parse --short origin/main` | `254e98f` | `254e98f` | 一致 |
| `git status --porcelain` | 空 | 空 | 一致 |
| 准则版本（CN 头行） | v1.73 | `v1.73`（日期 2026-11-25） | 一致 |
| 准则对行数（步骤 0 自测项，起跑消息未列） | — | 2557 / 2557，严格同位 | 补充读数 |
| SELFTEST | 744/744 | 见 `03-verification/selftest-step0.txt` | 一致 |
| ADR / OB / DR / 试点顺延 | 48 / 144 / 18 / 25 | 48 / 144 / 18 / 25 | 一致 |
| `tmp` 顶层三件 | `.gitkeep` / `gates-cache` / `find-quotes.js` | 同 | 一致 |

## 本片读取的权威输入

| 文件 | 用途 |
|---|---|
| `docs/DELIVERY_DIRECTIVE.md:3` | 头部版本行（版本包三件之一，G1-b） |
| `docs/DELIVERY_DIRECTIVE.md:43-120` | §0.3 变更日志区间与最新行（v1.73） |
| `docs/DELIVERY_DIRECTIVE.md:675` / `_EN.md:675` | 待同步的「脚本化状态」句（唯一改动正文点） |
| `docs/DELIVERY_DIRECTIVE.md:1220-1235` | §12.4 OB 台账尾部（OB-144 之后追加 OB-145） |
| `docs/DELIVERY_DIRECTIVE.md:2476-2491` | 附录 C 的 C.0az 块与 C.1 起点（新块 C.0ba 插入点） |
| `docs/validation/gates-ci-perf.md:64` / `_EN.md:64` | K4 拟改文本唯一出处 |
| `.github/workflows/ci.yml:80-110` | `Gates` 步与 `Format` 步实况（P90 落地态） |
| `docs/t027/CHECKLIST_v74.md:15-57` | v74 快照锚点与 OB-145 候选位 |
| `docs/t027/REMAINING_SLICES.md:2414-2452` | P90 台账段体例（P91 段照此结构） |

## 审查输入包

| 环节 | 输入包 |
|---|---|
| ① | 本文件的锚点表 + K4 原文 + `ci.yml:96-108` 实况 + 附录 C / §0.3 结构位置 |
| ⑤ / ⑥ / ⑦ | `05-review-pack/`（CN/EN diff、numstat、门禁原始读数、自由度清单、台账与报告） |
