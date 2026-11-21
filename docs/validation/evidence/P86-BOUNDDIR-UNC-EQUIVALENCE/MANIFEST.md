# MANIFEST — P86-BOUNDDIR-UNC-EQUIVALENCE 证据包

切片：`BOUNDDIR-UNC-EQUIVALENCE`（P86，实现类 `[SLICE]`，S 档）。**定稿**（2026-11-22；结论 `CLOSED`，见报告对 `docs/validation/bounddir-unc-equivalence{,_EN}.md`；⑧c 前只增不删）。

**载荷集**：本清单与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（磁盘 / 本清单 / `SHA256SUMS.txt`）。
**重算口径**：每次新增载荷后按磁盘派生重建 `SHA256SUMS.txt`，路径相对本目录。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 起跑消息逐字原文（§13.1.10） |
| `01-arch/01-arch-design.md` | ① 架构方案书（契约判定、§3 冲突裁定、§10 Rev 2） |
| `03-verification/build.txt` | go build ./... 读数 |
| `03-verification/gate-ci.txt` | gate.js --all --scope=ci 读数 |
| `03-verification/gate-index.txt` | gate.js --all --scope=index 读数（提交前） |
| `03-verification/gate-selftest.txt` | gate.js --selftest 读数（740/740） |
| `03-verification/gate-tree.txt` | gate.js --all --scope=tree 读数 |
| `03-verification/mutation-branch-removed.txt` | 删去 P86 分支后的红测日志 |
| `03-verification/mutation-record.txt` | 变异检验四步哈希与退出码记录 |
| `03-verification/probe-clean-windows.txt` | filepath.Clean / VolumeName 事实探针输出 |
| `03-verification/test-all-windows.txt` | go test ./... -count=1 全量读数 |
| `03-verification/test-p86-verbose.txt` | TestP86 定向 -v 读数（生产入口本机 SKIP） |
| `03-verification/test-post-restore.txt` | 逐字节恢复后的回归日志 |
| `03-verification/vet.txt` | go vet ./... 读数 |
| `04-reviews/05-pre-review-r1.md` | ⑤ 预审 r1 逐字 |
| `04-reviews/05-pre-review-r2.md` | ⑤ 预审 r2 逐字 |
| `04-reviews/06-scan-r1.md` | ⑥ 扫描 r1 逐字 |
| `04-reviews/06-scan-r2.md` | ⑥ 扫描 r2 逐字 |
| `04-reviews/07-blind-review.md` | ⑦ 盲审逐字 |
| `04-reviews/07-final-review-r1.md` | ⑦ 终审 r1 逐字 |
| `freedom-list.md` | 执行自由度清单（步骤 0 骨架；⑧b 定稿） |
| `input-list.md` | 输入清单（步骤 0 落盘） |

## 未入库文件清单

（无。`tmp/P86-BOUNDDIR-UNC-EQUIVALENCE/` 下的过程件不入库，⑧c 清除。）
