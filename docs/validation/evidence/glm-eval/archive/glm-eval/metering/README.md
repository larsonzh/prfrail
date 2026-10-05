# DIRECTIVE-GLM-EVAL · 计量记录规程（metering）

本目录用于保存每一次付费调用的**调用前 / 调用后**读数原始文件与说明。
本片（②实现者）不发起任何调用；调用由主控在执行步 2 / 3 使用。

## 1. 逐步规程

对每一次调用（`call_no` 递增），按顺序执行：

1. **调用前读**：记录当前用量读数，写入 `{seq}-pre.{ext}`。
   - `{seq}` = 两位序号，如 `01`、`02`（对应数据集 `call_no`）。
   - `{ext}` 按可得形态取 `txt` / `json` / `png`（截图）。
2. **调用**：只用 `03-prompt.txt` 作为任务包正文（bundle 四文件同为输入）。
3. **调用后读**：立即记录用量读数，写入 `{seq}-post.{ext}`。
4. **10 分钟复核**：在不早于 10 分钟后再次读取用量，写入 `{seq}-recheck.{ext}`，
   用于确认区间计量（interval-sum）是否稳定、是否出现滞后入账。

## 2. 文件命名

| 文件 | 含义 |
|---|---|
| `{seq}-pre.*` | 调用前读数 |
| `{seq}-post.*` | 调用后读数 |
| `{seq}-recheck.*` | 调用后 ≥10 分钟复核读数 |

原始读数**只增不改**；补记一律新开文件并在数据集的 `notes` 中说明。

## 3. metering_mode 判定

- `exact-delta`：读数提供精确的差分（可直接相减得到本次消耗）。
- `interval-sum`：只能按时间区间累加，需 `pre` / `post` / `recheck` 三读交叉核对。
- `unavailable`：无法取得任何可信读数（须在 `notes` 说明原因）。
- `user-provided-cumulative`：逐次读数不可得，改记**用户提供的累计读数**（`measured_tokens` 记累计口径字符串，
  **不拆分**到单次；须在 `notes` 说明读数来源与未剥离项）。本片对 GLM 四行即取此口径
  （见 `tmp/glm-eval/GLM-USAGE-READING.txt` 的更正节）。

## 4. 与数据集的关系

读数文件是本目录的原始证据；汇总值填入
`docs/validation/directive-glm-eval-dataset.json` 的 `calls[]` 对应行的
`measured_tokens` / `metering_mode` / `duration_s` 列。
