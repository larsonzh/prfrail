# MANIFEST — P87-DIRECTIVE-REVIEW-MODEL-POLICY 证据包

切片：`DIRECTIVE-REVIEW-MODEL-POLICY`（P87，治理片 `[SLICE]`，M 档）。**按磁盘派生重建**（⑧c 前只增不删；每次新增载荷后由脚本重建）。

**载荷集**：本清单与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（磁盘 / 本清单 / `SHA256SUMS.txt`）。
**重算口径**：新增载荷后按磁盘派生重建 `SHA256SUMS.txt`，路径相对本目录且用 POSIX 分隔符；P85 四读数据不复制进本包，按路径引用。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 起跑消息逐字原文（§13.1.10） |
| `01-arch/01-arch-design.md` | ① 架构方案书（契约判定、ADR-048 形态、风险枚举、主控就地裁定表） |
| `03-verification/build.txt` | go build ./... 读数 |
| `03-verification/gate-ci.txt` | gate.js --all --scope=ci --skip-base 读数 |
| `03-verification/gate-index.txt` | gate.js --all --scope=index --skip-base 读数 |
| `03-verification/gate-selftest.txt` | gate.js --selftest 读数（744/744） |
| `03-verification/gate-tree.txt` | gate.js --all --scope=tree --skip-base 读数 |
| `03-verification/numstat.txt` | git diff --cached --numstat 双语对称读数 |
| `03-verification/shape-six.txt` | G3(a.2) 六形状比对读数 |
| `03-verification/test-all.txt` | go test ./... 包级 ok/FAIL 汇总 |
| `03-verification/vet.txt` | go vet ./... 读数 |
| `04-reviews/05-pre-review-r1.md` | ⑤ 预审 r1 逐字 |
| `04-reviews/05-pre-review-r2.md` | （未登记描述） |
| `04-reviews/06-scan-r1.md` | ⑥ 独立扫描 r1 逐字 |
| `04-reviews/06-scan-r2.md` | （未登记描述） |
| `04-reviews/07-blind-review.md` | ⑦ 盲审逐字 |
| `04-reviews/07-final-review-r1.md` | ⑦ 终审 r1 逐字 |
| `freedom-list.md` | 执行自由度清单 |
| `input-list.md` | 输入清单与步骤 0 读数 |

## 未入库文件清单

无
