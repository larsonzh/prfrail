# FIX-GAP-77-SYMLINK-ESCAPE 执行自由度清单

| 项 | 内容 |
|---|---|
| 切片 | `FIX-GAP-77-SYMLINK-ESCAPE`（M 档 · 底层安全修复类 · `[SLICE]`） |
| 定位 | v42 §四 层 2 第 5 项；修复 GAP-77（Linux 平台 `bounddir_unix.go:64` 的 `openat(O_NOFOLLOW)` 仅护末组件，中间符号链接逃逸）；依赖 P52 闭环 |
| 可自行决定的自由度 | ① 修复方案二选一（(A) 逐级循环 `openat` + `O_NOFOLLOW`；(B) `O_PATH` 打开中间目录 + `/proc/self/fd` 相对打开），由 ① 架构评估后定稿；② T51-1 `t.Skip` 移除的时机与红测步骤编排（TDD 红测先行）；③ 界内多层目录读取回归用例的增补方式；④ 文档回写的措辞与落点小节（`SECURITY{,_EN}` SEC-02、`TEST_STRATEGY{,_EN}` AT-02、`ARCHITECTURE{,_EN}` §5.1、`DEV_PLAN{,_EN}`）；⑤ 批量辅助（GLM）派单的具体检索问题 |
| 必须升级的事项 | 改变 `bounddir_unix.go` 任何接口签名（起跑消息明令禁止）；修改 `CONTRACTS{,_EN}` 任何段落；改动 Windows 实现（`bounddir_windows.go`，本片限 Unix 平台）；改变 fail-closed 语义或放宽拒绝条件；误杀界内多层目录的正常读取；对既有测试的删除或语义弱化 |
| 载体铁律 | ① ⑤ `deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（只读）；② `prfrail-implementer`、③ `prfrail-tester`、⑧ `prfrail-documenter` + `DeepSeek V4.1 Flash (deepseek)`（可写）；⑥ `independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`；⑦ `independent-reviewer` + `GPT-5.3-Codex (copilot)`；批量辅助 `deep-reasoner` + `GLM-5.3-Flash (glm)`；每次 `runSubagent` 必须同时传 `agentName` 与 `model`；严格串行禁止并行 |
| 主控禁止 | 禁止代写 ②③ 的生产代码与测试；主控只做编排、落盘、门禁与回执 |
| 纪律 V 拆步 | 先设计逐级解析方案（① 架构定稿）→ ③ 红测先行 → ② 只改一个函数 → 验证 → 再谈下一步；严禁打包大改；超时后先查 `git status` 半截产物再决定回退 |
| 红测先行红线 | 移除 `TestBoundDirMiddleSymlinkEscapesBound` 的 `t.Skip` 必须发生在生产代码修改**之前**，并在 Ubuntu VM（`bash -lc`）上跑出真实红（读出界外字节）并保存日志；若红测未按 TDD 时序执行即属偏离 |
| 临时件 | 一律 `tmp/FIX-GAP-77-SYMLINK-ESCAPE/`，只增不删，⑧c 清理；测字符数用内联脚本，严禁在 `$env:TEMP` 等系统临时目录生成文件 |
| 禁止 | 推 gitee；在 ⑩ 之前 commit/push（⑩ 停点等同一轮授权） |
| 证据包 | `docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/`（`00-briefing.md` 为起跑消息原文；`.md` 一律 UTF-8 + BOM + LF） |
| 环境事实 | 本机 Windows + Go（直连，无代理）；远程 Ubuntu VM `10.0.0.199`（go1.22.12 @ `~/go1.22` 已就位，须 `bash -lc` 加载环境）承担 Linux 红绿验证与 `-race`；本机 `go test ./...` + 三作用域门禁为本地验收基线 |
| 主控模型 | 本片为底层并发/句柄类安全修复切片 ⇒ 主控提升为 `DeepSeek V4 Pro`（v42「纪律 P」+ 起跑消息第 3 条明示）；降级须按 OB-37 元规则独立登记 |
| ① 双打裁定（临时授权留痕） | 用户 2026-11-03 裁定 ① 架构采用双打模式：第一棒 `deep-reasoner` + `DeepSeek V4 Pro (deepseek)` 出 A/B 对比方案；第二棒 `independent-reviewer` + `GPT-6.1 Sol (copilot)` 底层安全抗压审查（**§1.7.4 临时授权，白名单外，一次性，不得作为先例**）；红队发现缝隙可直接打回 V4 Pro 修改。审查红线：接口签名不得改变；必须设计红测先行（跑红后才能 ② 进场）；若 `O_PATH` 方案依赖 `/proc/self/fd` 导致无 `/proc` 的极简容器环境不可用，架构必须设计降级为逐级 `openat` |
| ① 验收裁定（用户 2026-11-04，原文见 `01-arch/rulings.md`） | (a) 威胁边界**包含**外部 rename ⇒ `openat2` 主路径为硬要求；<5.6 内核等不支持环境**必须 fail-closed 拒绝**（返回 `ENOSYS` 包装错误），不允许降级到有竞态窗口的 walk 路径（`allowUnsafeFallback` 配置本片不启用）；(b) 挂载导入属攻击面 ⇒ 必须保留 `NO_XDEV`；(c) 不引入深度帽（删除 r1 的 64）；(d) 文档同步（SECURITY/TEST_STRATEGY/ARCHITECTURE/DEV_PLAN）纳入本片 ⑧ 派单；(e) 签认 `openat2` 超出二选一清单、批准为主路径 |
