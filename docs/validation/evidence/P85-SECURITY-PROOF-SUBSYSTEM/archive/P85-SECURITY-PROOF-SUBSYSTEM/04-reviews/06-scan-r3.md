# ⑥ 独立扫描 r3 原文（P85 · MAI Code 2 Flash · 额度 3/3）

INDEPENDENT SCAN: FINDINGS

## Section A
结论：FINDINGS；本轮不可判断 PASS，因为存在 11 条未闭合项及平台证据缺口。

## Section B
- F01 — NOT-CLOSED：装配锚和磁盘重读已实现，见 agent_runner_policy.go#L124，但同一 OS 身份仍可改策略文件与锚；无执行令牌不可写存储边界。
- F02 — NOT-CLOSED：末组件已走 bound open，但读取路径未对目录再做 VerifyBoundDir；Linux OpenBoundDir 仍可经目录名解析到目标，见 agent_runner_policy.go#L593 与 bounddir_unix.go#L96。
- F03 — CLOSED：Windows 权利表按逐项拆分并逐权利探测，见 ai_install_proof_windows.go#L86；本机 Windows 定向测试 PASS。
- F04 — NOT-CLOSED：祖先链与 reparse 检查已落地，但实现明确保留 CreateProcessW 重新解析整条名称链窗口，见 process_windows.go#L510。
- F05 — NOT-CLOSED：载体事实已传到启动边界并逐项比较，但最终 guard 只持摘要/身份路径，未把载体 fileID 延伸到最终打开句柄，存在同摘要异 inode 的重开窗口，见 ai_probe_copilot_executor.go#L631。
- F06 — CLOSED：Linux 组件改为 Lstat，叶文件与摘要/身份同 bound handle，见 ai_install_proof_unix.go#L112；原生 Linux 执行仍在 Section E 记录为不可验证。
- F07 — CLOSED：env 形式拒绝、递归解释器和 ELF PT_INTERP 已实现，见 ai_install_proof_unix.go#L240；原生 Linux 复跑仍缺。
- F08 — NOT-CLOSED：有效 euid/gid 判定已加，但未实际核对 POSIX ACL、named ACL 或已打开写句柄；实现仍明确把已打开句柄移出威胁范围，见 ai_install_proof_unix.go#L158。
- F09 — CLOSED：策略制品和载体两处均使用 no-replace os.Link，无 Rename 替代分支，见 agent_runner_policy.go#L640。
- F10 — NOT-CLOSED：Unix 已同步目录，Windows 明确无目录句柄 sync，且实现承认该平台豁免，见 agent_runner_policy.go#L733。
- F11 — CLOSED：Schema 闭集、valid/invalid 夹具与契约措辞均同步；契约测试 4/4 PASS。
- F12 — CLOSED：本地拒绝证据只含 status+reason，profile/time 留在记录元数据，见 ai_preflight.go#L108。
- F13 — CLOSED：冻结模板在凭据解析前构造，并在 spawn 时逐项比较 argv、环境、绑定和工作目录，定向反例全 PASS。
- F14 — CLOSED：凭据键按平台归一，大小写歧义与认证模式不一致均拒绝，见 ai_probe_copilot_executor.go#L672。
- F15 — CLOSED：ai check 前置拒绝采用 cli_unavailable/ErrorClass 统一输出，见 cli.go#L867。
- F16 — NOT-CLOSED：共享测试仍包含合成 C:\... Windows 路径；在 Linux 上会因 filepath.IsAbs 变相失败，见 p85_t3_tester_test.go#L118。
- F17 — NOT-CLOSED：同 inode 改写为显式特权豁免 skip；真实部分 ACL、同 inode 改写、解释器漂移未形成原生双平台反例，见 p85_t3_linux_tester_test.go#L197。
- F18 — NOT-CLOSED：T2 mutation 证据完整，但 T8a 缺独立可控启动 seam；执行器反例仍可被宿主 T3 可证性或平台 verdict 覆盖，见 p85_t8a_tester_test.go#L550。
- F19 — CLOSED：task/step/attempt 使用受约束的无歧义键，磁盘复用比较走原字节，见 agent_runner_policy.go#L78。
- F20 — NOT-CLOSED：载体读写使用 bound dir/leaf，但 Unix 读取路径未在每次读取后校验目录句柄仍是派生目录；目录级重定向仍无独立证明，见 ai_install_proof.go#L130。
- F21 — CLOSED：Linux/Windows 专属实现与其他平台 fail-closed stub 已到位，契约已明确 Linux fd 3 例外，见 process_unsupported.go#L1。
- F22 — NOT-CLOSED：当前 CI 摘要是既有 scope=ci snapshot，不是 P85 双平台包级原生证据；未发现 P85 原生 CI 日志能证明本片包运行。

- New findings：无独立编号新发现；以上未闭合项均来自盲审原 22 条。

## Section C
证据链：工作区 HEAD e35e50f、门禁 index/tree 均 TOTAL_FAIL=0、SELFTEST 740/740、契约测试 4/4 PASS、go build ./... 与 go vet ./... 无输出成功。P85 定向 T2、T8a、Windows F04/F05 测试均 PASS；全量 Go 首次有 P84 C7 Windows lease Access is denied 瞬态，单独重跑后 PASS。

## Section D
可证伪点包括：同 UID 同时改策略文件和锚、目录替换为 junction/symlink、启动前替换成同摘要异 inode、Windows 创建进程前切换 reparse、named ACL 授权写、Windows crash 后重新读未持久目录、Linux 下运行共享 C:\ 载体测试、特权 fixture 运行同 inode 改写，以及 P85 CI 缺失时不代换 gate-ci。

## Section E
无法验证：Ubuntu VM 副本是旧/dirty 基线 dc785ed，无最新 P85 源文件，不是可比原生 Linux 证据；未取得 P85 双平台 CI 的包级 ok 行；独立 mutation 未以源码落盘运行，以上结论以静态差异、当前 Windows 实测和现有测试为准。
