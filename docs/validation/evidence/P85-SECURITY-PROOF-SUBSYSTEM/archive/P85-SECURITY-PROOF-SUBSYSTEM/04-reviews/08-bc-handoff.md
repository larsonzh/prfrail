# P85 · B/C 类移交文案（⑦ 终审输入件之一；⑥ r4 核如实性）

## B 类（四条 · 威胁模型边界 + 设计层已登记残余，交 ⑦ 终审逐条裁定）

- **F01（同 OS 身份可改制品+锚）**：S1 威胁模型显式边界 = `SECURITY.md:11`
  「相同 OS 用户下目录分离不是恶意进程沙箱」。本片防线 = 装配冻结锚 +
  四等式对账（`Digest(body)==制品 policyHash==recorded==grant.PolicyHash==锚`）+
  no-replace 硬链接 + 发布/晋升前磁盘重读 + 路径不相交断言；防的是
  「装配后意外重写与漂移对账」，同 OS 身份恶意进程的主动写入不在 S1 承诺。
  实现注释已逐字登记（`agent_runner_policy.go` 锚类型注释块）。
- **F04（CreateProcessW 重解析窗口）**：防线 = 全链逐级判定（卷根身份锚定 +
  非卷根逐权利 + 每级 reparse 拒绝）+ 叶级 `openApprovedExecutable` 持有
  deny-write/deny-delete 句柄跨 verify→spawn + pre-spawn 实时重验。残余 =
  内核 `CreateProcessW` 最终名称解析瞬间（句柄锁定的是叶对象，整条名称链的
  最终解析由内核完成）；已由 ① 方案书 Rev2 §2.2 与实现注释登记。
- **F08（POSIX named ACL / 已打开写句柄）**：Linux 判定口径 = 逐组件
  `Lstat` 拒链接 + 模式位按 euid 有效权限 + 属主非本用户/root；named ACL
  与「已打开写句柄」属 S1 威胁模型外（同 OS 身份进程），未纳入判定；已由
  ① 方案书 Rev2 §3.1 与实现注释登记。
- **F10（Windows 目录 sync 平台豁免）**：Go 在 Windows 上无法 sync 目录句柄
  （`syscall.Fsync` 对目录返回错误），沿 applier 原子写先例：rename/link 成功
  即持久点；Linux 侧已补父目录 + 已存在目录的 sync。豁免已由实现注释登记。

## C 类（六条 · ⑨ 原生验证与证据计划，随输入包移交）

- **F06/F07（Linux 原生复跑）**：Ubuntu VM（10.0.0.199，larson，go1.22.12）
  同步本片工作树后原生执行：`p85_t3_linux_tester_test.go`（不可写目录 /
  属主 / shebang 解释器三反例 + 非 ELF 反例 + 真实 ELF 正例 + PT_INTERP
  断言）+ `go test -count=1 ./...` 全量；读数入证据包 `03-verification/`。
- **F16（Linux 共享测试含 `C:\` 路径）**：VM 原生跑时验证共享测试文件在
  Linux 腿的实况（跳过或通过），并按其结果决定「宿主绝对路径化」或
  「EXEMPT 登记」；读数入证据包。
- **F17（同 inode 改写特权夹具）**：维持 EXEMPT-REGISTERED 四要素登记
  （root 属主夹具需 chown、VM 无免密 sudo），承接测试两枚已列；⑦ 终审复核
  登记如实性。
- **F18（T8a 独立启动 seam）**：变异红留痕（checker 门禁 / executor 断言块
  / override 键循环三处）已由 ③ 报告归档；VM/CI 复跑 T8a 用例时输出留痕；
  独立 seam 缺口的登记随报告 §8。
- **F22（CI 双平台包级证据）**：提交后 CI 双腿（windows/ubuntu）包级
  `ok <pkg>` 行 + Gates 步硬证据，入报告「提交与 CI」节（§11.2 体例）；
  gate-ci 本地读数只作观察。
- **⑥ r3 的 Section E 观察（VM 旧基线 `dc785ed`）**：⑨ 前必须 `git fetch`
  并把 VM 工作树同步到本片 HEAD（`e35e50f` + 本片未提交改动经 `git diff
  --binary` 或 rsync 同步），旧基线读数不可比，不得复用。
