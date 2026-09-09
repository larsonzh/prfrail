# S1 Windows 便携包安装演练

[English](s1-install-drill_EN.md)

日期：2026-09-09，更新于 2026-09-10。结论：`PASS AS CANDIDATE DRILL`。本报告证明已归档的 Windows amd64 可信 artifact 可在显式临时目录完成离线核验、首次运行、升级、回滚和保留数据的卸载；它不授权 GitHub Release。产品所有者随后选择手工解压便携 ZIP、用户自选独立目录、显式路径运行且不修改 PATH 的 S1 模型。

## 输入

1. 旧包：GitHub Actions run `34277671704` 的 `prfrail-Windows-1e7af676e8e84028c7ecfef1ccde91213738c20c`。
2. 新包：Node 24 Actions run `34325758548` 的 `prfrail-Windows-603e8dc03e9fb74ce9c53243d38e704fd8469c6a`。
3. 两包均含 `prfrail.exe`、`SHA256SUMS`、`sbom.spdx.json` 和 `licenses.json`；使用当前受审源码构建的 `prfrail-release verify` 分别按对应完整 candidate SHA 核验，退出码均为 0。
4. 演练根目录仅为仓库忽略的 `tmp/t018-install-drill/`，不写用户 PATH、注册表、系统目录、源树或 Git。

## 执行结果

| 步骤 | 结果 | 证据 |
|---|---|---|
| 首次安装 | PASS | 旧包复制到独立 `installed/current`；`version --json` 返回 `0.1.0+1e7af676e8e84028c7ecfef1ccde91213738c20c` |
| 升级 | PASS | 新包先复制到 staging，再将 current 保留为 rollback；新版返回 `0.1.0+603e8dc03e9fb74ce9c53243d38e704fd8469c6a` |
| 回滚 | PASS | 删除新版 current 后恢复旧目录；旧版再次返回原完整版本 |
| 卸载 | PASS | 删除 installed 二进制目录；独立 `data/retained-evidence.txt` 仍存在且内容为 `retain-me` |
| 安装教程复验 | PASS | 将 run `34325758548` 的四文件 artifact 压为 ZIP 后按 `INSTALLATION.md` 从干净目录解压；逐项 SHA-256、两份 JSON、`version`、Go 样例 `validate/preview` 均通过，进程及用户 PATH 前后不变 |
| `gh` 下载路径 | PASS | `gh run download 34325758548 --repo larsonzh/prfrail --name prfrail-Windows-603e8dc03e9fb74ce9c53243d38e704fd8469c6a` 成功，并自动展开为且仅为四个预期文件；临时目录已删除 |

旧包 `prfrail.exe` 的 SHA-256 为 `5b30f2b3f0de23ce662d6c75a27395cc59d5087678016faf3daebbda7ffc2db1`；新包为 `31335030152a417757990fd36e64f882d891a14112e0a2b865e1fc7a4280bfc8`。校验和只证明内容与清单一致，不认证发布者身份。

## 最终候选复验（嵌套目录变更前）

2026-09-10，使用可信 workflow run [`34358467083`](https://github.com/larsonzh/prfrail/actions/runs/34358467083) 的 Windows artifact `prfrail-Windows-bcc235496e0bece825047e95f3165a36585a59c0` 完成第二轮生命周期及五步快速上手复验。该 artifact ID 为 `10106867086`，大小 `2,398,240` 字节，未过期，GitHub archive digest 为 `sha256:10d902164194a8b1c6d1121adba5363bb69f829ac3516f7b891c6f5f41d08e06`。

| 检查 | 结果 | 证据 |
|---|---|---|
| 升级 | PASS | 旧版 `0.1.0+603e8dc03e9fb74ce9c53243d38e704fd8469c6a` 升级为 `0.1.0+bcc235496e0bece825047e95f3165a36585a59c0`；旧/新 EXE SHA-256 分别为 `31335030152a417757990fd36e64f882d891a14112e0a2b865e1fc7a4280bfc8` 和 `c77a5780abbc381a7a696c24d3f428286fa1ebe9d77f21f06bcb6ac2ce41d692` |
| 回滚 | PASS | 恢复后 `version --json` 再次返回旧版完整版本 |
| 安装副作用 | PASS | 进程 PATH、用户 PATH 和 ProofRail 进程 ID 在演练前后均未改变 |
| 五步快速上手 | PASS | `init`、`validate`、`preview`、noop-only `run`、`report` 全部通过；run 为 `chainState: COMPLETED`、`events: 11`，报告为 `taskState: PASSED=1`、`stepState: NOOP_RECORDED=1` |
| 清理 | PASS | 临时演练目录已删除，Git 工作树保持干净 |

该 run 的 verified artifact 仍是归档根目录直接包含四个文件的旧布局。2026-09-10 新增的 `prfrail/` 单一顶层目录尚未由该 run 验证，因此以上结果是最终候选内容与生命周期证据，不是新 ZIP 布局的最终载体验收。

## 边界

1. 这是已有 CI artifact 的候选便携安装演练，不是正式发行包或最终用户教程。
2. S1 不设默认安装目录、不修改 PATH、不提供隐式版本切换；下载入口和最终文件名仍待发行批准。
3. 含 `prfrail/` 顶层目录的新正式 ZIP 仍需重跑同一流程，并将版本、commit、run、文件名、支持矩阵和发行说明绑定到该载体。
4. 临时 artifact 与演练目录在证据记录完成后删除，不纳入版本控制。