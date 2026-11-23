# L4 修后 diff 摘录（2026-11-23 · 审查环后修改 · 依据 D-2）

> **如实标注**：以下为⑦审查环结束后、按 Carbon High 定向小读 Low-4 建议改动的纯注释/表述级修改；
> 属 P26 裁定 D-2 约束范围（审查环结束后按审查方建议改措辞须有下游实质审查）。
> 处置：⑧c/⑩ 回执附本摘录 + 顾问独立实质核验 + 用户终审（用户 2026-11-23 停点 6 追加裁定，同轮授权例外）。

## 1. `internal/console/ai_install_proof.go`（复用分支，约 :125）

- 修改前：
  `// Windows boundary (F04/F10, registered): the hard link gives ATOMIC VISIBILITY ...`
- 修改后：
  `// Windows durability boundary (independent S1 registration, user ruling of 2026-11-23, not an F04 exemption): the hard link gives ATOMIC VISIBILITY ...`
- 所依据的代码事实：
  - 所在分支 = 载体重入复用分支：`aiInstallSourceConfirmBoundDirectory`（同绑定句柄内 read/confirm/sync 三合一，:118-:133）；
  - 该注释主体是 **F10 耐久性边界**（Windows 无法经 Go sync 目录句柄，无断电耐久），不是 CreateTemp/Link 发布路径；
  - 发布路径 S1 边界已独立登记于 `internal/console/agent_runner_policy_publish_other.go:16` 与
    `internal/console/ai_install_proof_publish_other.go:15`（"it is not an F04 exemption"）；
  - 同文件 :258 已为「not an F04 exemption」措辞。

## 2. `internal/evidence/acl_linux.go` 文件头（约 :19-25）

- 修改前：
  `The path form remains only for the nominal entry, which holds no descriptor at all and re-lstats the name around the read.`
- 修改后：
  `The path form remains for the two call sites that hold no descriptor: the nominal write-domain entry, which re-lstats the name after the read, and the install-source ancestor walk, which lstats each component once before the read (see AccessACLPath).`
- 所依据的代码事实：
  - 名义入口 = `agent_runner_policy_write_guard_unix.go` `managedAgentPolicyDomainAfterProbe` 的 `!hasHandle` 分支
    （按名 Getxattr **之后**二次 `os.Lstat` 全字段比对，约 :191-:205）；
  - 祖先链 = `internal/adapters/ai_install_proof_unix.go` 约 :168-:178：每组件**单次前置** `os.Lstat` +
    拒 symlink + `installSourceDenyAccessACLPath`，**无**读后二次 Lstat。

## 3. `internal/evidence/acl_linux.go` `AccessACLPath` 注释（约 :95-107）

- 修改前：
  `... used only where no descriptor exists (the nominal write-domain entry and the ancestor components of the install-source walk); its callers re-lstat the name around the read so a name that moved is refused instead of described (P85 R7).`
- 修改后（双要点）：
  ① nominal write-domain entry：读后二次 re-lstat 并与 before 观测比对，移动即拒；读取本身仍按名解析 =
  已登记内部 API 边界；② install-source ancestor walk：每组件**单次前置** lstat 并拒链接组件，
  身份/mode 事实先于 ACL 读，该路径**无**读后 re-lstat。
- 所依据的代码事实：同上第 2 项两处调用点。

## 4. `internal/adapters/ai_install_proof_unsupported.go`（约 :10-13）

- 修改前：
  `installSourceInterpreter is declared by the Windows and Linux halves only:`
- 修改后：
  `installSourceInterpreter is declared by the LINUX half only:`
- 所依据的代码事实：
  - `internal/adapters/ai_install_proof_unix.go` 首行 build tag = `//go:build linux`（仅 Linux 编译）；
  - 全仓 `installSourceInterpreter` 定义仅 `ai_install_proof_unix.go:313` 一处（R7b 处置后行号；顾问核验实测）；Windows 无 shebang
    解释器链实现，无该函数声明。

## 修改后机械验证
- `gofmt -l` 零输出；`go build ./...` exit=0；门禁 index/tree 双作用域 TOTAL_FAIL=0
  （2026-11-23 主控实测，修后即跑）。
