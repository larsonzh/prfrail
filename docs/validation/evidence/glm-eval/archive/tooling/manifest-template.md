# 冻结证据包 · 切片 `DIRECTIVE-GLM-EVAL`

| 项 | 值 |
|---|---|
| 包根 | `docs/validation/evidence/glm-eval/` |
| 切片 | `DIRECTIVE-GLM-EVAL`（评估片；**只评估、不入白名单**） |
| 归档日期 | 2026-10-05 |
| 载荷范围 | `archive/glm-eval/**`（会话工作目录全量：bundle / 两腿输出 / 评分与污染脚本 / 自测器与夹具 / 计量规程 / 交接文档 / 方案书与勘误）；`archive/glm-eval-out/**`（两腿原始输出 + 自测报告时间点快照）；`archive/glm-eval-batch/**`（批量腿 a 的输入与钥匙）；`archive/tmp-root/**`（承重脚本与 `PROMPT-EXACT.txt`）；`archive/tooling/**`（本包的可复现工具：打包器与 MANIFEST 模板） |
| 条目数 | {{COUNT}} |
| 哈希清单 | `SHA256SUMS.txt`（`<hash><2 spaces><path>` 形态，POSIX 分隔符） |
| 归档口径 | 逐字节复制，**未做任何编码或行尾转换**：载荷原始即为「非 `.ps1` ⇒ 无 BOM」且一律 LF，已满足 §1.5.1「归档后类型归一」。故每条**原始字节 sha256 = 归档字节 sha256**（见下「后缀与编码映射」节） |
| 自洽核验 | 本包须过 `node tools/gates/gate.js --check=G8-b --scope=tree`（归档条目集合 ≡ `SHA256SUMS.txt` 条目集合 ≡ 包内载荷集合，且逐条重算 sha256 一致） |

## 归档条目清单

| 归档路径 | 字节 | sha256 | 归档理由 |
|---|---|---|---|
{{ENTRIES}}

## 后缀与编码映射

| 原始路径 | 归档路径 | 原始字节 sha256 | 归档字节 sha256 | 归一 |
|---|---|---|---|---|
{{ENCMAP}}

## 未入库文件清单

无。本包不引用任何未入库文件；全部载荷均已随包入库。

## 排除清单（有意不归档）

- 无。本包覆盖 `tmp/glm-eval/**`、`tmp/glm-eval-out/**`、`tmp/glm-eval-batch/**` 与 `tmp/` 根承重脚本的全部条目（来源清单 = `tmp/glm-eval/FROZEN-MANIFEST.txt` 的 72 条 + 该清单自身）。
- 一次性中间物（`tmp/_g4bchk.js`、`tmp/_g4b.txt`、`tmp/_persist2.log`、`tmp/_append-*.txt`、`tmp/_manifest-*.md`）属 §1.5.1 第 ③ 类，**用完即清**，不入包。

## 重建说明

- **不适用**：本包全部为**原始字节归档**，无「重建归档」条目（§1.5.1 第 ④ 类）。
- 唯一被覆写过的载荷是 `archive/glm-eval/selftest/selftest-report.txt`：其**脚手架期原始字节不可复现**（已作为不可逆取证损失登记于报告与 `TEST-REPORT` §12）；包内存放的是整改后重跑版本，并另存时间点快照 `archive/glm-eval-out/selftest-report-20261005.txt`（与前者同哈希）。

## 对应引用点

- 报告：`docs/validation/directive-glm-eval{,_EN}.md`（§3.2 bundle 表 / §3.6 冻结自证 / §6 可证伪性 / §9 门禁结果）。
- 评估数据集：`docs/validation/directive-glm-eval-dataset.json`。
- 自由度清单：`docs/validation/evidence/DIRECTIVE-GLM-EVAL-freedom-list.md`（证据根级 `*-freedom-list.md` 形态）。
- 台账：`docs/t027/REMAINING_SLICES{,_EN}.md` 的 `DIRECTIVE-GLM-EVAL` 完成行。
