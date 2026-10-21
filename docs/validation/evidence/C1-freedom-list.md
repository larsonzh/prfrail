# 执行自由度清单 · 切片 `C1`

> 归属：`docs/validation/t027-c1-linux-native{,_EN}.md` 的 §9 引用本清单为权威，本清单不改变报告含义。
> 口径：F = 本片自行作出的选择；M = 已裁定、不得重开的判定；D = 偏离（异常）登记。每行给出 `file:line` 或偏移证据。

## F 选择

| # | 选择 | 证据 |
|---|---|---|
| F1 | 验证驱动用 bash 写在 `tools/validation/c1-linux/`，不在 PowerShell 里套 `ssh` | `tools/validation/c1-linux/run-c1-faces.sh:1-60` |
| F2 | 四个面各作主动收窄（机制存在 / 进程树与出口 / 实现层 proven 路径 / f1–f4 + CLI 腿） | ① 设计；报告 §5 |
| F3 | f5 排除但保留 `face_needs_cmd_exe` 兜底分支，让「跳过」可记录而非静默 | `run-c1-faces.sh:234-241`、`:259-265` |
| F4 | 调用矩阵与 `run-faces.ps1` 逐组对齐，不新增 `f4-pos-r2` | `run-c1-faces.sh:397-417` 对 `run-faces.ps1:248-285` |
| F5 | VM 侧输出改到 `~/prfrail-c1/c1-out`，脚本默认值仍留 `$C1_REPO/tmp/c1/` | `run-c1-faces.sh:122-124`；③ 运行命令 |
| F6 | 脚本用一次 `scp` 送 VM 并核对 sha256，不用 heredoc 传输 | ③ 报告第 1 节 |
| F7 | ④ 除判据外还跑三条 base 腿（`go build` / `go vet` / `go test`） | `tmp/C1/review-pack/03-gate4-tree.txt` |
| F8 | 三个 Medium 统一取「登记已知限制 + 收窄读」，不动脚本 | 报告 §4、§8 |
| F9 | 证据包按 `docs/validation/evidence/C1-2026-10-21/` 归档，正文引用其相对路径 | `docs/validation/evidence/C1-2026-10-21/MANIFEST.md` |
| F10 | 包内剔除 run root、seed、emitted inputs 与二进制，只留可复核读数 | 同 MANIFEST 的排除段 |

## M 判定（不得重开）

| # | 判定 | 依据 |
|---|---|---|
| M1 | 本片只覆盖 Linux AgentRunner native validation，不含 Linux SecretStore 与 SessionBridge visible | 用户起跑指令 |
| M2 | f5 排除、`f4` 的 pos 腿不得补 r2 重放；验证片不得改被测 harness | 用户 2026-10-21 裁定 |
| M3 | ④ 的 `G5(d)` 属结构性预期红，⑤ 起派不因此阻塞；⑨ 为权威读数 | 用户 2026-10-21 裁定 |
| M4 | 临时件违规第 9 次按「自报 + 自删 + 零残留」判据继续，且作为子项合并计数 | 用户 2026-10-21 裁定 |
| M5 | 三个 Medium 均以「已知限制 + 收窄读」闭合，不改脚本以保全 ③ 证据哈希 | 用户 2026-10-21 裁定 |
| M6 | KL-C1-7 升格为复用前置门槛，下一片复用前必须先修 | 用户 2026-10-21 裁定 |

## D 偏离

| # | 偏离 | 影响 | 证据 |
|---|---|---|---|
| D1 | ② 接线自测在 `%TEMP%` 留下 `c1-wire*` 目录与 `c1-inputs.json` | 放置规则第 9 次，子项合并计数；文件已删 | 报告 §3；`tmp/C1/violation-9.md` |
| D2 | 主控首轮「零残留」误报（`-First 5` 截断 + 过滤过窄） | 方法缺陷，已自我纠正并写入核验纪律 | 报告 §3；`tmp/C1/violation-9.md` |
| D3 | ⑤ r1 Medium #1：脚本未强制绝对路径，② 自报与源码不符 | 登记 KL-C1-1；⑤ r2 判闭合 | 报告 §8 |
| D4 | ⑥ r1 Medium #1：复用输出目录可致旧 verdict 被当本轮结果 | 登记 KL-C1-7；⑥ r2 判闭合 | 报告 §8 |
| D5 | ⑦ r1 Medium #1：断网归因过强；r2 判不闭合，修订后 r3 判闭合 | 登记 KL-C1-8；两次返工 | 报告 §4、§8 |
| D6 | ④ 的 `G5(d)` 暂存未声明红 | 结构性预期红，非缺陷；⑨ 复核 | 报告 §7 |
| D7 | 主控在 PowerShell 内联 `node -e` 写临时代码时，反引号被当作转义符，产出过一份错字临时件 | 在 `tmp/C1/` 内删除重写，未入库 | 报告 §3 的自我纠错节 |

- 本清单不新增任何未在报告中登记的自由度；F9/F10 只描述证据包的归档口径，不构成规格。
