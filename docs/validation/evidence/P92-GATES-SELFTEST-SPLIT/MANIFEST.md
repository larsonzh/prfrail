# MANIFEST — P92-GATES-SELFTEST-SPLIT 证据包

切片：`GATES-SELFTEST-SPLIT`（P92，工具片 `[SLICE]`，M 档）。2026-11-26 起跑。

**载荷集**：本清单与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（磁盘 / 本清单 / `SHA256SUMS.txt`）。
**重算口径**：新增载荷后按磁盘派生重建 `SHA256SUMS.txt`，路径相对本目录（POSIX 分隔）；哈希在编码归一之后计算。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 起跑消息逐字原文与四要素（§13.1.10） |
| `01-design.md` | ① 影响判定 + 收益论证闸 + 实现范围裁定（方案书） |
| `02-probe-data.md` | 分区计时与子进程归因实测数据（改动前 / 记忆化变体 / 落地后复测） |
| `03-verification/ab-hashes.txt` | A/B/B2 sha256 与可证伪翻转对照 |
| `03-verification/base-gates.txt` | 基础门禁：gofmt / go build / go vet / go test 原始读数 |
| `03-verification/gate-ci-observational.txt` | ⑨ 本地 ci 作用域观察读数（推送前非权威） |
| `03-verification/gate-index-final.txt` | ⑨ 定稿态 index 作用域原始读数 |
| `03-verification/gate-tree-final.txt` | ⑨ 定稿态 tree 作用域原始读数 |
| `03-verification/selftest-A.txt` | 改动前基线捕获（A 侧，与 B1/B2 同 sha256） |
| `03-verification/selftest-B1.txt` | 改动后第 1 次捕获（B1） |
| `03-verification/selftest-B2.txt` | 改动后第 2 次捕获（B2） |
| `03-verification/selftest-final.txt` | 载荷 |
| `04-reviews/05-pre-review-r1.md` | ⑤ 预审 r1 原文 + 主控处置表 |
| `04-reviews/06-scan-r1.md` | ⑥ 独立扫描 r1 原文（FINDINGS，指向当时未生成的下游工件） |
| `04-reviews/06-scan-r2.md` | ⚠️ 独立扫描 r2 原文（FINDINGS：证据包哈希在收敛循环中的中间态漂移） |
| `04-reviews/06-scan-r3.md` | ⚠️ 独立扫描 r3 原文（收敛后重跑） |
| `04-reviews/07-final-r1.md` | ⑦ 终审 r1 原文（FINDINGS：围栏三行非逐字） |
| `04-reviews/07-final-r2.md` | ⑦ 复审 r2 原文（前 Medium 闭合；新报 tree 捕获中间态） |
| `04-reviews/07-final-r4.md` | ⑦ 收窄确认读 r4 原文（额度例外，1 次）+ 主控处置 |
| `freedom-list.md` | 执行自由度清单 F-1..F-7（§1.7.3） |
| `input-list.md` | 输入清单（权威输入 / 改动候选面 / 实测件 / 步骤 0 读数） |
