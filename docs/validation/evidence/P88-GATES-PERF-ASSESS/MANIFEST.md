# MANIFEST — P88-GATES-PERF-ASSESS 证据包

切片：`GATES-PERF-ASSESS`（P88，工具片 `[SLICE]`，S 档）。**按磁盘派生重建**；⑧c 前只增不删。

**载荷集**：本清单与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（磁盘 / 本清单 / `SHA256SUMS.txt`）。
**重算口径**：新增载荷后按磁盘派生重建 `SHA256SUMS.txt`，路径相对本目录且用 POSIX 分隔符。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 起跑消息逐字原文（§13.1.10） |
| `01-arch/01-arch-design.md` | ① 架构方案书（影响判定 / 优化落点 / A-B 验证设计 / 风险 / 大改建议） |
| `02-impl/02-notes.md` | ② 实现说明（改动清单、fail-open 判定表、不改变输出的论证） |
| `03-verification/ab-report.md` | ③ A/B 逐字一致性 + 耗时对照报告 |
| `03-verification/ab-sha256.txt` | ③ A/B 各格 stdout 的 sha256 清单 |
| `03-verification/baseline.tsv` | 优化前耗时构成基线（逐项 exit + elapsed ms） |
| `03-verification/build.txt` | go build ./... 读数 |
| `03-verification/gate-ci.txt` | gate.js --all --scope=ci 读数 |
| `03-verification/gate-index.txt` | gate.js --all --scope=index 读数 |
| `03-verification/gate-selftest.txt` | gate.js --selftest 读数（744/744） |
| `03-verification/gate-tree.txt` | gate.js --all --scope=tree 读数 |
| `03-verification/gofmt.txt` | gofmt -l . 读数（空输出 = 通过） |
| `03-verification/test.txt` | go test ./... 包级 ok/FAIL 读数 |
| `03-verification/timings.md` | 优化前后耗时对照汇总 |
| `03-verification/vet.txt` | go vet ./... 读数 |
| `04-reviews/05-pre-review-r1.md` | ⑤ 预审 r1 逐字 |
| `04-reviews/06-scan-r1.md` | ⑥ 独立扫描 r1 逐字（报 1 HIGH） |
| `04-reviews/06-scan-r2.md` | ⑥ 独立扫描 r2 逐字（补齐配对证据后 PASS） |
| `04-reviews/07-final-review-r1.md` | ⑦ 终审 r1 逐字（FINAL REVIEW: PASS，0 发现） |
| `05-batch-assist/glm-mirror-check.md` | 批量辅助：报告对双语镜像核对逐字（MIRROR: PASS） |
| `freedom-list.md` | 执行自由度清单（F-1..F-10） |
| `input-list.md` | 输入清单与步骤 0 锚点读数 |
