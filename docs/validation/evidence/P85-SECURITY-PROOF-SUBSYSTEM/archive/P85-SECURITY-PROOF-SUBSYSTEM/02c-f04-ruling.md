# P85 · F04 卷根排除边界裁定登记稿（报告 §8 引用，逐字备查）

- 裁决人 = 用户；日期 = 2026-11-23；理由 = 顾问疑点（盘根默认 ACL 对
  Authenticated Users 授予「仅此文件夹的创建子目录」权利，若全链逐权利判定含
  卷根且计创建权利，任何安装路径在盘根组件即判可变更 ⇒ 系统性误拒，与 F08
  修正过严误拒的精神相反）；依据 = 本机最小实测 + ① single-round 快速裁定。
- **B 裁定（① 裁决）**：全链判定维持纳入卷根级（不排除卷根本身——排除会漏判
  「持盘根 DELETE_CHILD / WRITE_DAC 换掉首级祖先」的替换路径）；卷根以 P84
  身份原语锚定（`evidence.BoundDir.IsFixedVolumeRoot`，FileIndexLow==5）。
- **用户边界（2026-11-23 裁定，四笔）**：
  1. 排除范围 = 仅「真卷根」的 `FILE_ADD_FILE` / `FILE_ADD_SUBDIRECTORY` 两权利
     （卷根目录项系统保留、不可 swap，建子项不影响既有路径组件的名称解析）；
  2. 非卷根祖先目录保持全权利判定（rev2 七项，含 `FILE_DELETE_CHILD` /
     `FILE_ADD_FILE`——它们直接决定既有目录项能否被删换）；
  3. 卷根身份以 P84 原语锚定（FileIndexLow==5）；**身份锚定失败 ⇒ 回全权利
     判定（七项）fail-closed，不得按卷根排除**；
  4. 报告 §8 登记：B 裁定 + icacls 实测证据 + 排除理由；如实注明本裁定系
     顾问疑点触发、① single-round 裁决。
- **实测证据（2026-11-23，本机 Medium 非提升令牌，C# 内同步 errno）**：
  - `C:\` 的 `FILE_ADD_SUBDIRECTORY` = GRANTED（`icacls C:\` 显示
    `NT AUTHORITY\Authenticated Users:(AD)`——append-data/add-subdirectory、
    仅此文件夹、无继承）；`FILE_ADD_FILE` = ACCESS_DENIED；
    `GENERIC_WRITE` = ACCESS_DENIED；
  - `C:\Windows\System32` / `C:\Program Files`：全部 ACCESS_DENIED（Users 仅 RX）。
  - 教训：跨 P/Invoke 层读 errno 会被污染（PowerShell 版本假 GRANTED），
    必须在同一 C#/Go 调用内同步读取；Go 实现全部经 syscall 直读。
- **实现落点**：`internal/adapters/ai_install_proof_windows.go`——目录表恢复
  rev2 七项；新增 `installSourceVolumeRootMutationRights`（五项）；卷根分支
  「锚定成功 ⇒ 五项表；锚定失败 ⇒ 七项表」；objectStates 记录
  `anchor:fixed-volume` / `anchor:unproven` 标记。
- **③ 测试义务**：正例「卷根仅授 `FILE_ADD_SUBDIRECTORY` ⇒ 判定通过」钉住
  排除不回归；反例「非卷根祖先授 `FILE_ADD_FILE` ⇒ 拒」；反例「锚定失败回
  全表（授 ADD_FILE ⇒ 拒）」。
