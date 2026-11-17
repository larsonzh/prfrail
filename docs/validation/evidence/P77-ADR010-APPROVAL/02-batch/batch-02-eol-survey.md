# P77 批量辅助 2 —— EOL / 安全维护窗口 / 版本维护政策 表述枚举（12 份）

**计数口径（事实性说明）**：
- `hits` = 给定关键词集（`EOL`、`End of Life`、`生命周期`、`维护窗口`、`安全维护`、`支持窗口`、`90 天`、`ADR-010`、`已批准`、`未获批`、`尚未获批`、`待审批`、`Release notes`、`公告`）的逐字命中行数；一行含多个关键词只计 1 行。
- 另列出两类**不计入 hits** 的条目（同样逐字摘录）：
  - **[B] 镜像条目**：与另一语言文件中关键词命中行内容对应的行（该行本身不含给定关键词集）。
  - **[C] 主题补充**：不含给定关键词、但内容直接涉及 EOL 通知绑定 / 支持起始日期 / 支持范围 / 支持矩阵生效批准状态的行（仅限此四类，未扩大）。
- 「版本与安全维护候选政策」整段按要求逐句列出（其中无关键词的句子标注 `[段内-无关键词]`）。

---

## 1. docs/S1_SUPPORT_MATRIX.md

命中：**是**（关键词命中行：9, 16, 18, 19, 21, 22, 23, 29, 31）

**「版本与安全维护候选政策」整段逐句（L16–L23）：**

- `docs/S1_SUPPORT_MATRIX.md:16` ｜命中词：`安全维护`
```text
## 版本与安全维护候选政策
```
- `docs/S1_SUPPORT_MATRIX.md:18` ｜命中词：`EOL`
```text
- 产品版本采用 `vMAJOR.MINOR.PATCH`，首版候选为 `v0.1.0`；patch 发布不启动新的 EOL 倒计时。
```
- `docs/S1_SUPPORT_MATRIX.md:19` ｜命中词：`维护窗口`、`90 天`
```text
- 当前稳定 minor 与前一个 minor 接收安全修复；前一个 minor 的候选维护窗口为 90 天。
```
- `docs/S1_SUPPORT_MATRIX.md:20` ｜[段内-无关键词]
```text
- 历史未知 run 只读保留，不自动迁移。
```
- `docs/S1_SUPPORT_MATRIX.md:21` ｜命中词：`待审批`、`ADR-010`、`Release notes`、`公告`、`EOL`
```text
- 待审批的 ADR-010 提议：新 minor 发布时在 GitHub Release notes 与本支持矩阵同时公告前一 minor 的 EOL 日期，日期为发布后第 90 个自然日。
```
- `docs/S1_SUPPORT_MATRIX.md:22` ｜命中词：`90 天`
```text
- 90 天是原型期初始窗口；每次 minor 发布前按稳定性、用户量和维护能力评审，仅可延长后续新公告的窗口。
```
- `docs/S1_SUPPORT_MATRIX.md:23` ｜命中词：`公告`（两处）、`尚未获批`（含子串 `未获批`）
```text
- 任何偏离须有带日期、理由和到期日的 ADR，经产品所有者明确批准后同步两处，且不得追溯缩短已公告窗口。该提议尚未获批，本页不能替代最终公告。
```

**其余命中：**

- `docs/S1_SUPPORT_MATRIX.md:9` ｜命中词：`生命周期`
```text
| Windows 11 amd64 | 正式支持候选 | 纯 Go `CGO_ENABLED=0` 便携 ZIP；手工解压、显式路径运行、不修改 PATH；最终 ZIP 须重跑安装生命周期；基础终端（逐行 CLI / `interactions tui` / `--progress`）保留；完整一体化 TUI 移出 S1、由 S2 承接 |
```
- `docs/S1_SUPPORT_MATRIX.md:29` ｜命中词：`EOL`
```text
> 切片同步（P55 `EXIT-SEED-BUILD-EVIDENCE`）：本片不适用（无平台与版本支持范围增减）；理由：补证只归档既有 seed 驱动构建的两平台捕获，平台清单、支持层级与 EOL 政策不变。
```
- `docs/S1_SUPPORT_MATRIX.md:31` ｜命中词：`EOL`
```text
> 切片同步（P57 `P5-P6-CLOSEOUT-ASSESSMENT`）：本片不适用（无平台清单、支持层级、EOL 政策或安装模型的增减）；理由：P57 只把 T018 剩余门整理为收口清单，未改变任何门的判定，本文件措辞逐字不动。
```

**[C] 主题补充（不计入 hits）：**

- `docs/S1_SUPPORT_MATRIX.md:5` ｜（支持矩阵自身生效批准状态）
```text
日期：2026-09-09。状态：`CANDIDATE`，待正式发行载体和 S1 exit 批准后生效；当前不是已发布 SLA。
```
- `docs/S1_SUPPORT_MATRIX.md:27` ｜[B]（EN L27 关键词命中行的 CN 镜像；含支持起始日期与最终发行说明批准）
```text
候选发行要求完整 commit SHA、GitHub CI、`SHA256SUMS`、SPDX SBOM 和许可证清单。校验和证明内容与清单一致，不认证发布者身份；S1 不宣称签名或 attestation。安装模型见 [便携 ZIP 安装指南](INSTALLATION.md)；具体文件名、下载位置和支持起始日期须与最终发行说明一起批准。
```

`docs/S1_SUPPORT_MATRIX.md: hits=9`

---

## 2. docs/S1_SUPPORT_MATRIX_EN.md

命中：**是**（关键词命中行：18, 21, 27, 29, 31）

**「Candidate Version and Security-Maintenance Policy」整段逐句（L16–L23）：**

- `docs/S1_SUPPORT_MATRIX_EN.md:16` ｜[B]（CN L16 关键词命中行的 EN 镜像）
```text
## Candidate Version and Security-Maintenance Policy
```
- `docs/S1_SUPPORT_MATRIX_EN.md:18` ｜命中词：`EOL`
```text
- Product versions use `vMAJOR.MINOR.PATCH`, beginning with candidate `v0.1.0`; patch releases do not start a new EOL countdown.
```
- `docs/S1_SUPPORT_MATRIX_EN.md:19` ｜[B]（CN L19 镜像）
```text
- The current stable minor and previous minor receive security fixes; the candidate maintenance window for the previous minor is 90 days.
```
- `docs/S1_SUPPORT_MATRIX_EN.md:20` ｜[段内-无关键词]
```text
- Unknown historical runs remain read-only and are not automatically migrated.
```
- `docs/S1_SUPPORT_MATRIX_EN.md:21` ｜命中词：`EOL`、`ADR-010`、`Release notes`
```text
- Proposed ADR-010: when a new minor is released, announce the previous minor's EOL date in both the GitHub Release notes and this support matrix, with EOL on the 90th calendar day after release.
```
- `docs/S1_SUPPORT_MATRIX_EN.md:22` ｜[B]（CN L22 镜像）
```text
- Ninety days is the initial prototype-phase window; before each minor release, review stability, adoption, and maintenance capacity, and only lengthen future newly announced windows.
```
- `docs/S1_SUPPORT_MATRIX_EN.md:23` ｜[B]（CN L23 镜像）
```text
- Deviations require a dated ADR with rationale and expiry, explicit product-owner approval, and synchronized publication in both channels, and must not retroactively shorten an announced window. This proposal is not approved, and this page does not replace the final notice.
```

**其余命中：**

- `docs/S1_SUPPORT_MATRIX_EN.md:27` ｜命中词：`Release notes`
```text
A candidate release requires a full commit SHA, GitHub CI, `SHA256SUMS`, an SPDX SBOM, and a license manifest. Checksums prove agreement with a manifest, not publisher identity; S1 does not claim signing or attestation. See the [portable ZIP installation guide](INSTALLATION_EN.md); filenames, download locations, and the support start date must be approved with the final release notes.
```
- `docs/S1_SUPPORT_MATRIX_EN.md:29` ｜命中词：`EOL`
```text
> Slice sync (P55 `EXIT-SEED-BUILD-EVIDENCE`): not applicable (no platform or version support-range change); reason: the supplement only archives the two-platform capture of the existing seed-driven build, so the platform list, support tiers and EOL policy stay unchanged.
```
- `docs/S1_SUPPORT_MATRIX_EN.md:31` ｜命中词：`EOL`
```text
> Slice sync (P57 `P5-P6-CLOSEOUT-ASSESSMENT`): not applicable (no change to the platform list, support tiers, EOL policy or install model); reason: P57 only turns the remaining T018 gates into a close-out inventory and changes no gate verdict, so this file's wording stays byte-identical.
```

**[C] 主题补充（不计入 hits）：**

- `docs/S1_SUPPORT_MATRIX_EN.md:5` ｜（支持矩阵自身生效批准状态）
```text
Date: 2026-09-09. Status: `CANDIDATE`, effective only after the formal release carrier and S1 exit are approved. This is not a published SLA.
```
- `docs/S1_SUPPORT_MATRIX_EN.md:9` ｜[B]（CN L9 关键词命中行的 EN 镜像）
```text
| Windows 11 amd64 | Production-support candidate | Pure-Go `CGO_ENABLED=0` portable ZIP; manual extraction, explicit-path invocation, no PATH modification; rerun the lifecycle against the final ZIP; the basic terminal (line-oriented CLI / `interactions tui` / `--progress`) stays; the full unified TUI moves out of S1 and is taken up by S2 |
```

`docs/S1_SUPPORT_MATRIX_EN.md: hits=5`

---

## 3. docs/INSTALLATION_PLAN.md

命中：**是**（关键词命中行：28, 32, 38；三条均位于文件开头 40 行内）

- `docs/INSTALLATION_PLAN.md:28` ｜命中词：`生命周期`
```text
| 卸载与数据保留 | 手工删除版本目录；run/store 默认保留，数据删除单独授权 | T024 生命周期记录和显式删除授权 | S1 模型与库级规则已冻结 |
```
- `docs/INSTALLATION_PLAN.md:32` ｜命中词：`EOL`
```text
S1 不指定默认安装目录，不发布一键安装脚本，不修改 PATH，不声明自动更新，也不要求管理员权限。state/store 默认位置、正式下载入口与 EOL 通知仍待冻结。
```
- `docs/INSTALLATION_PLAN.md:38` ｜命中词：`生命周期`
```text
该结果关闭“候选 artifact 是否能执行生命周期”的工程疑问，并支撑 2026-09-09 获批的手工便携 ZIP、不改 PATH 模型。正式下载入口和最终 ZIP 尚未发布，仍须用最终载体重跑。
```

**[B] 镜像条目（不计入 hits）：**

- `docs/INSTALLATION_PLAN.md:63` ｜（EN L63 关键词命中行的 CN 镜像，含 `发行说明` 与 `获批`）
```text
T016 已冻结 CLI、配置和当前源码运行方式；T024 已冻结库级备份、升级前阻断、停用和数据保留规则；T017 的候选 CI artifact、固定 commit、SHA256SUMS、SBOM/许可证门禁与双平台可信矩阵已完成。签名/attestation 是后续增强。S1 安装模型已于 2026-09-09 获批；最终 ZIP 复验和 S1 exit 获批前，本文仍不构成生产发布声明。[操作手册](OPERATIONS.md) 继续区分当前可执行与未来设计；候选发布边界见 [支持矩阵](S1_SUPPORT_MATRIX.md) 与 [发行说明](S1_RELEASE_NOTES.md)。
```

`docs/INSTALLATION_PLAN.md: hits=3`

---

## 4. docs/INSTALLATION_PLAN_EN.md

命中：**是**（关键词命中行：32, 63）

- `docs/INSTALLATION_PLAN_EN.md:32` ｜命中词：`EOL`
```text
S1 defines no default install directory, publishes no one-click installer, does not modify PATH, claims no automatic updates, and requires no administrator privileges. The default state/store location, formal download entry point, and EOL notification channel remain open.
```
- `docs/INSTALLATION_PLAN_EN.md:63` ｜命中词：`Release notes`
```text
T016 has frozen CLI, configuration, and current source-run behavior; T024 has frozen library-level backup, pre-upgrade blocking, retirement, and data-retention rules. T017 candidate CI artifacts, fixed-commit, SHA256SUMS, SBOM/license gates, and the dual-platform trust matrix are complete. Signing/attestation is a later enhancement. The S1 installation model was approved on 2026-09-09. This document still does not claim a production release until the final ZIP rerun and S1 exit are approved; [Operations](OPERATIONS_EN.md) continues to distinguish executable behavior from future design. See the candidate [support matrix](S1_SUPPORT_MATRIX_EN.md) and [release notes](S1_RELEASE_NOTES_EN.md) for release boundaries.
```

**[B] 镜像条目（不计入 hits）：**

- `docs/INSTALLATION_PLAN_EN.md:28` ｜（CN L28 关键词命中行的 EN 镜像）
```text
| Uninstall/data retention | Manually remove the version directory; retain run/store by default and separately authorize data deletion | T024 lifecycle record and explicit deletion authorization | S1 model and library rules frozen |
```
- `docs/INSTALLATION_PLAN_EN.md:38` ｜（CN L38 关键词命中行的 EN 镜像）
```text
This closes the engineering question of whether candidate artifacts can execute the lifecycle and supports the manually extracted portable-ZIP, no-PATH model approved on 2026-09-09. The formal download entry point and final ZIP do not yet exist, so the final carrier must rerun the drill.
```

`docs/INSTALLATION_PLAN_EN.md: hits=2`

---

## 5. docs/PRODUCT_REQUIREMENTS.md

命中：**是**（关键词命中行：41, 105, 107, 123）

- `docs/PRODUCT_REQUIREMENTS.md:41` ｜命中词：`生命周期`
```text
本文件中的角色、场景和 REQ 必须不依赖 whois 背景即可理解和验收。whois 只提供历史问题证据与未来影子夹具；其 A/B、D/V、Step47 或 start-file 行为不构成 ProofRail 隐含需求。正常流、失败流、恢复、人工交接、交付和生命周期统一由 [业务流程](BUSINESS_WORKFLOWS.md) 串联；安装发行的未知项由 [安装部署规划](INSTALLATION_PLAN.md) 管理，不允许实现者自行选择默认值。
```
- `docs/PRODUCT_REQUIREMENTS.md:105` ｜命中词：`已批准`（§7 标题）
```text
## 7. 已批准边界、待落实事项与后续调研
```
- `docs/PRODUCT_REQUIREMENTS.md:107` ｜命中词：`EOL`、`90 天`（§7 内「EOL 通知」相关句）
```text
2026-09-06，产品所有者批准 S0 规格冻结及 ADR §4：S0 新增付费调用为 0，Windows 11 amd64 正式/Linux amd64 核心 CI，当前与前一 minor（前一 minor 90 天），逐依赖许可审计，以及未来 whois 脱敏只读影子验证原则。2026-09-08，所有者进一步批准 S1 简化发布信任：固定 commit SHA、GitHub CI、SHA256SUMS、SBOM 和许可证清单为必需证据，不建立离线密钥体系，签名/attestation 后续增强。精确依赖、EOL 通知、安全评审者和具体 whois 输入仍待后续任务落实，不因已有实现而自动取得付费调用、数据访问或发布授权。
```
- `docs/PRODUCT_REQUIREMENTS.md:123` ｜命中词：`待审批`
```text
| PC-03 | 运行期间授权撤销、审批积压缺乏产品流程 | REQ-005/021/023 | 哈希绑定运行授权、撤销/过期阻断、待审批队列；S2 脱敏通知/RBAC | AT-18 |
```

`docs/PRODUCT_REQUIREMENTS.md: hits=4`

---

## 6. docs/PRODUCT_REQUIREMENTS_EN.md

命中：**是**（关键词命中行：105）

- `docs/PRODUCT_REQUIREMENTS_EN.md:105` ｜命中词：`EOL`（§7 内「EOL notice」相关句，含 90-day 窗口表述）
```text
On 2026-09-06, the product owner approved the S0 specification freeze and ADR section 4: zero new paid S0 calls; Windows 11 amd64 production/Linux amd64 core CI; current and previous minor support with a 90-day previous-minor window; per-dependency license review; and future redacted read-only whois shadow validation. On 2026-09-08, the owner simplified S1 release trust: fixed commit SHA, GitHub CI, SHA256SUMS, SBOM and license manifest are required evidence; no offline key system is maintained, and signing/attestation is deferred. Exact dependencies, EOL notice, security reviewer, and concrete whois input remain later-task details; existing implementation grants no automatic authority for paid calls, data access, or release.
```

**[B] 镜像条目（不计入 hits）：**

- `docs/PRODUCT_REQUIREMENTS_EN.md:39` ｜（CN L41 关键词命中行的 EN 镜像）
```text
The roles, scenarios, and REQs in this document must be understandable and testable without whois knowledge. whois supplies historical problem evidence and a future shadow fixture only; its A/B, D/V, Step47, or start-file behavior is not an implicit ProofRail requirement. [Business workflows](BUSINESS_WORKFLOWS_EN.md) connect the normal, failure, recovery, handoff, delivery, and lifecycle paths. [Installation planning](INSTALLATION_PLAN_EN.md) owns unresolved distribution decisions so implementers cannot invent defaults.
```
- `docs/PRODUCT_REQUIREMENTS_EN.md:103` ｜（CN L105 §7 标题的 EN 镜像）
```text
## 7. Approved Boundaries, Open Details and Further Research
```
- `docs/PRODUCT_REQUIREMENTS_EN.md:121` ｜（CN L123 关键词命中行的 EN 镜像）
```text
| PC-03 | Revocation and pending approvals lack an operational flow | REQ-005/021/023 | Hash-bound run authorization, revocation/expiry blocking and approval inbox; S2 redacted notifications/RBAC | AT-18 |
```

`docs/PRODUCT_REQUIREMENTS_EN.md: hits=1`

---

## 7. docs/S1_RELEASE_NOTES.md

命中：**是**（关键词命中行：23）

**「发布前剩余门禁」段逐行（L20–L25）：**

- `docs/S1_RELEASE_NOTES.md:20` ｜[段标题-无关键词]
```text
## 发布前剩余门禁
```
- `docs/S1_RELEASE_NOTES.md:22` ｜[段内-无关键词]
```text
1. 对最终 ZIP 复跑核验、首次运行、升级、回滚和卸载。
```
- `docs/S1_RELEASE_NOTES.md:23` ｜命中词：`EOL` ｜**EOL 相关句**
```text
2. 将版本、artifact 文件名、下载位置、EOL 通知方式和支持日期绑定到本说明及支持矩阵。
```
- `docs/S1_RELEASE_NOTES.md:24` ｜[C]（发行批准门禁）
```text
3. 获得独立评审结论和产品所有者明确 S1 exit/发行批准。
```
- `docs/S1_RELEASE_NOTES.md:25` ｜[段内-无关键词]
```text
4. 完成 T027 并通过 AT-23；T028 已完成、AT-24 条件闭合；S1 仍因 T027/AT-23 阻断。
```

**其余条目：**

- `docs/S1_RELEASE_NOTES.md:1` ｜[B]（EN L1 关键词命中行的 CN 镜像）
```text
# ProofRail S1 候选发行说明
```
- `docs/S1_RELEASE_NOTES.md:5` ｜[C]（支持起始日期与批准人冻结状态）
```text
日期：2026-09-09。状态：`DRAFT / NOT RELEASED`。版本、下载链接、正式 artifact 文件名、支持起始日期和批准人均待 S1 exit 冻结；本文件不授权 tag、GitHub Release、发布或安装。
```
- `docs/S1_RELEASE_NOTES.md:18` ｜[C]（支持范围以支持矩阵为准）
```text
支持范围以 [S1 候选支持矩阵](S1_SUPPORT_MATRIX.md) 为准。S1 安装模型是[手工解压 Windows 便携 ZIP](INSTALLATION.md)、显式路径运行且不修改 PATH。当前仍无正式发行包，界面是逐行 CLI。正式 AgentRunner/CLI Agent 仍未交付；受监督 visible 黑箱候选的产品侧已由 T028 实现（AT-24 条件闭合），但真实 CLI 端到端依赖 T027-B、尚不可运行发布；该路径即使完成也只保证归还后的产物验收，不保证工具、网络、费用或外部副作用，并须用户显式确认风险和监督。S1 也不提供签名/attestation、Linux 正式支持、自动更新、Web 控制台、完整 TUI 或任意 executable hook。
```

`docs/S1_RELEASE_NOTES.md: hits=1`

---

## 8. docs/S1_RELEASE_NOTES_EN.md

命中：**是**（关键词命中行：1, 23）

- `docs/S1_RELEASE_NOTES_EN.md:1` ｜命中词：`Release notes`（标题）
```text
# ProofRail S1 Candidate Release Notes
```

**「Remaining Pre-Release Gates」段逐行（L20–L25）：**

- `docs/S1_RELEASE_NOTES_EN.md:20` ｜[B]（CN L20 段标题的 EN 镜像）
```text
## Remaining Pre-Release Gates
```
- `docs/S1_RELEASE_NOTES_EN.md:22` ｜[段内-无关键词]
```text
1. Rerun verification, first run, upgrade, rollback, and uninstall against the final ZIP.
```
- `docs/S1_RELEASE_NOTES_EN.md:23` ｜命中词：`EOL` ｜**EOL 相关句**
```text
2. Bind version, artifact filenames, download location, EOL notification channel, and support date to these notes and the support matrix.
```
- `docs/S1_RELEASE_NOTES_EN.md:24` ｜[C]（release approval 门禁）
```text
3. Obtain the independent-review verdict and explicit product-owner S1 exit/release approval.
```
- `docs/S1_RELEASE_NOTES_EN.md:25` ｜[段内-无关键词]
```text
4. Complete T027 and pass AT-23; T028 is complete with AT-24 conditionally closed; S1 remains blocked by T027/AT-23.
```

**其余条目：**

- `docs/S1_RELEASE_NOTES_EN.md:5` ｜[C]（support start date 与 approver 冻结状态）
```text
Date: 2026-09-09. Status: `DRAFT / NOT RELEASED`. Version, download links, formal artifact filenames, support start date, and approver remain pending S1 exit. This file does not authorize a tag, GitHub Release, publication, or installation.
```
- `docs/S1_RELEASE_NOTES_EN.md:18` ｜[C]（support matrix 支持范围引用）
```text
See the [S1 candidate support matrix](S1_SUPPORT_MATRIX_EN.md). The S1 installation model is a [manually extracted Windows portable ZIP](INSTALLATION_EN.md), invoked by explicit path without modifying PATH. No formal release package exists and the current interface is line-oriented CLI. Formal AgentRunner/CLI Agent are still undelivered; the supervised visible black-box candidate's product side is implemented by T028 (AT-24 conditionally closed) but its real CLI end-to-end depends on T027-B and it is not yet a runnable release. Even when implemented, the latter warrants post-return artifact acceptance only, not tools, network, cost or external effects, and requires explicit user risk acknowledgment and supervision. S1 also excludes signing/attestation, production Linux support, automatic updates, a Web console, the full TUI, and arbitrary executable hooks.
```

`docs/S1_RELEASE_NOTES_EN.md: hits=2`

---

## 9. docs/INSTALLATION.md

命中：**否**（给定关键词集 0 命中；已对全文件独立检索）

**[B] 镜像条目（不计入 hits）：**

- `docs/INSTALLATION.md:206` ｜（EN L206 关键词命中行的 CN 镜像，含 `发行说明`）
```text
正式发行前，以上命令须用最终 ZIP 再执行一次，并在 [S1 exit 报告](validation/s1-exit.md)、[支持矩阵](S1_SUPPORT_MATRIX.md)和[发行说明](S1_RELEASE_NOTES.md)中绑定最终版本、commit、run 和文件名。
```

`docs/INSTALLATION.md: hits=0`

---

## 10. docs/INSTALLATION_EN.md

命中：**是**（关键词命中行：206）

- `docs/INSTALLATION_EN.md:206` ｜命中词：`Release notes`
```text
Before formal release, rerun these commands against the final ZIP and bind its version, commit, run, and filename in the [S1 exit report](validation/s1-exit_EN.md), [support matrix](S1_SUPPORT_MATRIX_EN.md), and [release notes](S1_RELEASE_NOTES_EN.md).
```

`docs/INSTALLATION_EN.md: hits=1`

---

## 11. docs/OPERATIONS.md

命中：**是**（关键词命中行：139, 146）

- `docs/OPERATIONS.md:139` ｜命中词：`待审批`
```text
| PC-03 审批/撤销 | 查看待审批队列与候选摘要；撤销后确认无新投递及在途停机证据 | 不能停机则维持暂停；不通过 waive 放宽硬门禁 |
```
- `docs/OPERATIONS.md:146` ｜命中词：`生命周期`
```text
历史证据含秘密：隔离、停止外发、撤销凭据，再由人决定审计保留与删除冲突；不改旧哈希伪造完整链。删除不能保证 SSD/外部副本安全擦除。报告默认本地无遥测；提交支持材料只提供复现所需脱敏最小包。T020 已完成库级导出能力，T021 已实现授权/撤销与审批队列（approvals），T022 已实现副作用分类/观察与恢复诊断（库级），T023 已实现成本账本与本地成本报告（cost report），T024 已实现库级完整备份、新 store 核验、生命周期处置与发布清单；T017 发行物与安装命令仍待实现，完整产品流程当前仍不能作为可执行命令教程。
```

`docs/OPERATIONS.md: hits=2`

---

## 12. docs/OPERATIONS_EN.md

命中：**否**（给定关键词集 0 命中；已对全文件独立检索，EOL 邻近表述使用 `lifecycle` 等英文词，不在给定关键词集内）

**[B] 镜像条目（不计入 hits）：**

- `docs/OPERATIONS_EN.md:139` ｜（CN L139 关键词命中行的 EN 镜像）
```text
| PC-03 Approval/revocation | Inspect inbox/candidate hash; verify blocked new dispatch and in-flight stop evidence after revocation | Remain paused if stopping fails; no hard-gate waiver |
```
- `docs/OPERATIONS_EN.md:146` ｜（CN L146 关键词命中行的 EN 镜像，含 `lifecycle disposition`）
```text
For historical secrets, quarantine/block export/revoke credentials, then obtain a human retention/deletion decision; never rewrite old hashes as intact history. Deletion does not guarantee SSD/external-copy erasure. Reports stay local without default telemetry; support submissions contain only minimal redacted reproduction data. T020 now provides library-level export capability; T021 adds authorization/revocation with an approvals inbox, T022 adds side-effect classification/observation and recovery diagnosis (library level), T023 adds the cost ledger with local cost reporting (`cost report`), and T024 adds library-level complete backup, new-store verification, lifecycle disposition, and release inventory. T017 artifacts and installation commands remain pending, so the full product flow is still not an executable command tutorial today.
```

`docs/OPERATIONS_EN.md: hits=0`

---

## CN/EN 命中数对照表（hits = 给定关键词集逐字命中行数）

| # | 文件 | hits |
|---|---|---|
| 1 | docs/S1_SUPPORT_MATRIX.md | 9 |
| 2 | docs/S1_SUPPORT_MATRIX_EN.md | 5 |
| 3 | docs/INSTALLATION_PLAN.md | 3 |
| 4 | docs/INSTALLATION_PLAN_EN.md | 2 |
| 5 | docs/PRODUCT_REQUIREMENTS.md | 4 |
| 6 | docs/PRODUCT_REQUIREMENTS_EN.md | 1 |
| 7 | docs/S1_RELEASE_NOTES.md | 1 |
| 8 | docs/S1_RELEASE_NOTES_EN.md | 2 |
| 9 | docs/INSTALLATION.md | 0 |
| 10 | docs/INSTALLATION_EN.md | 1 |
| 11 | docs/OPERATIONS.md | 2 |
| 12 | docs/OPERATIONS_EN.md | 0 |

关键词集命中合计：30 行（CN 21 / EN 9）。

BATCH-02: DONE