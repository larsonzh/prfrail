# ProofRail 架构决策记录

[English](ADR_REGISTER_EN.md)

日期：2026-09-06。状态分为 RFC 已决定、提议待评审、待实验；本表不代表用户已批准 S1。每项批准时记录主体、日期、RFC 变更和测试证据，禁止模型代签。

| ID | 决策/方案 | 替代与理由 | 状态/完成条件 |
|---|---|---|---|
| ADR-001 | Go 模块化单体、CLI/TUI、single writer、无默认 Git 写入 | 不选脚本翻译/分布式系统；降低跨平台与恢复复杂度 | RFC §4/11/16 已决定；T001 确认依赖基线 |
| ADR-002 | 冻结 Schema 2020-12 工具链、canonical JSON/路径/错误/配置优先级 | 采用 `canonicalize` 4.0.0 复算 RFC 8785 向量；Ajv 8.20.0 独立于未来 Go 核心 | T002 已冻结 29 份结构 Schema；当前 82 个手写预期 fixture 与两条 canonical 向量全部通过 |
| ADR-003 | 事件/journal 与接受提交协议、Windows 原子替换/停机原语 | rename 只解决单文件；需证明 durable、失败回滚和跨文件发布读者语义 | T004 已实测：NTFS 开放 reader 阻止普通替换且目录 `Sync` 不可用，ext4 支持开放 reader 替换与目录 `Sync`；rename 后 receipt 前必须判 uncertain。S1 采用平台原语并补断电故障注入；不能用目录隔离替代 required OS 限制 |
| ADR-004 | SessionBridge v0.1.1 silent + 文件队列，核心自管持久幂等 | 不 fork 扩展、不用 visible/auto、无 GUI 兜底；内存缓存不等于 exactly-once | T004 已实测：Windows rename-only claim 可多赢家，固定 claim `O_EXCL`、requestId 去重和 generation fencing 为强制项；SessionBridge 40 项无模型测试通过，真实 LM API 仍属宿主能力门禁 |
| ADR-005 | bootstrap 首个 seed 由固定 commit、常规构建、独立 Go tests 和人工审计建立；N 构建 N+1，隔离重放 | 不接受候选自证或覆盖运行中 host | RFC §16.5 已定；T017 记录 seed commit/摘要、外部 oracle 和 rollback 路径 |
| ADR-006 | Windows 11 amd64 正式、Linux amd64 核心 CI；纯 Go 发布；固定工具链/依赖，S1 以 commit+CI+SHA256SUMS+SBOM/许可清单为发布证据 | RFC §16.2 提到 Win10，与 §9.6/14 的 Win11 范围不同；不建立 S1 离线密钥体系，签名/attestation 后续增强 | 2026-09-08 所有者批准简化；T017 落实基础发布证据 |
| ADR-007 | 最小 context、离线先行、按任务预算、双语 ID 追踪 | 不用昂贵模型反复全仓审读；不把英文译本做第二权威源 | 2026-09-06 所有者批准；S0 新增付费调用为 0，超出须另授权 |
| ADR-008 | RFC §19 六类产品最小闭环纳入 S1；仍用本地模块、不建商业后台 | 不采用只交执行引擎、隐式回写/部署或静默自更新；新增成本必须显式评审 | 2026-09-06 所有者批准范围；T002 已冻结记录，T003 已补样例，T019–T024 运行验证；未授权 S1 编码 |
| ADR-009 | S1 Windows amd64 使用手工解压的便携 ZIP；用户选择独立目录，以完整 `prfrail.exe` 路径运行，不修改 PATH/注册表/系统目录 | 不采用安装器、包管理器、自动更新或隐式版本切换；最小化安装副作用并支持目录级并存/回滚 | 2026-09-09 产品所有者明确批准；最终 ZIP 仍须绑定发行记录并复验 |
| ADR-010 | 产品版本采用 `vMAJOR.MINOR.PATCH`（首版候选 `v0.1.0`）；原型期新 minor 发布时，在 GitHub Release notes 与仓库支持矩阵同时公告前一 minor 的 EOL，初始窗口为发布后 90 个自然日 | patch 不启动新的 EOL 倒计时；每次 minor 发布前按稳定性、用户量和维护能力评审是否延长未来窗口。不得追溯缩短已公告窗口；例外须有带日期、理由和到期日的 ADR，经产品所有者明确批准后同步两处 | `PROPOSED / NOT APPROVED`；90 天适合当前早期原型，通知渠道与演进规则须在首版发布前获产品所有者明确批准 |

## 1. 发布与依赖决策提案

模块最低 Go 1.22 不因本机 Go 1.27.0 自动提升；发布使用实施时仍获安全维护的固定补丁版本，CI 同时验证最低兼容版本。若 TUI/依赖要求更高版本，先提 ADR/RFC 变更，不偷偷改 go.mod。纯 Go 核心优先标准库；TOML、Schema、canonical、TUI 采用成熟纯 Go 库，精确版本、许可证、漏洞和最低 Go 要求由 T001/T002 核对后锁定。本轮不安装项目依赖。

支持策略已由所有者于 2026-09-06 批准：仅维护当前稳定 minor 和前一个 minor 的安全修复，前一个 minor 保留 90 天；未知历史 run 保持只读、不自动迁移。ADR-010 提议采用三段式 `vMAJOR.MINOR.PATCH`，首版候选为 `v0.1.0`；patch 发布不重置 EOL，新 minor 发布时用 GitHub Release notes 与仓库支持矩阵同时公告前一 minor 的 EOL。90 天是早期原型的初始窗口，后续仅可通过发布前评审逐步延长未来窗口。例外须由带日期、理由和到期日的所有者批准 ADR 记录，且不得追溯缩短已公告窗口。该提议尚未获批，当前仍不是已发布 SLA。

发布信任策略由所有者于 2026-09-08 简化：S1 不建立离线 Ed25519 密钥体系，以固定 commit SHA、GitHub CI、SHA256SUMS、SBOM 和许可证清单作为基础发布证据。SHA256SUMS 只证明下载内容与清单一致，不证明发布者身份；签名或 keyless attestation 移至后续增强，不阻断 T017/T018 或 S1 exit。`signature-receipt` 能力保留供可选跨主机/后续签名场景使用。

## 2. 必须关闭的歧义

T002 已集中冻结空链、列表覆盖/合并、duration/size/预算单位与默认、id/attempt/requestId 类型、error code、路径/symlink/hardlink、签名覆盖范围、review policy 与 waiver 授权、票据阈值、hook runner 最小能力、write lease 续期/接管及 Windows 10 支持措辞。T003 独立样例/validator 与 T004 平台探针已完成；生产默认值仍须由对应 S1 实现、原生故障测试和供应链门禁关闭。

## 3. S0 Readiness 当前结论

结论：**NOT_READY**，不意味着本轮文档工作失败，也不授权跳过门禁。

| RFC §17 项 | 当前证据/缺口 | 关闭任务 |
|---|---|---|
| 1 P0 文档评审 | 所有者已于 2026-09-06 接受第三方只读审计建议并批准 S0 规格冻结/决策单；独立安全签署仍见第 3 项 | 已关闭（T001） |
| 2 Schema/黄金/validator | 29 份结构 Schema 各有正反例；当前 82 个 fixture 与两条 canonical 向量通过独立 oracle | 已关闭（T003） |
| 3 STRIDE | 风险/责任/AT 映射已写；T003 只读安全复核未发现 fail-open，运行期安全证据仍由后续 gate 提供 | 已关闭 S0 静态部分（T003） |
| 4 whois + 非 C | whois 保持映射样例；Go/JavaScript/Python 非 C fixture 已通过 | 已关闭（T003） |
| 5 技术探针 | Windows 11/NTFS、Ubuntu 24.04/ext4、TUI、进程树、文件队列及 SessionBridge 无模型合同已实测，限制见 `validation/s0-spikes.md` | 已关闭（T004） |
| 6 可追踪 backlog | 所有者已批准当前 S0/S1 范围和追踪计划 | 已关闭（T001） |
| 7 平台/Go/许可/周期/响应/发布证据 | 平台、支持期、逐依赖许可及简化发布证据策略已批准；精确依赖/渠道未锁定 | T017 |
| 8 用户批准 | 已批准 S0 规格冻结、S1 范围、零新增付费预算和 whois 影子原则；未授权 S1 编码、具体 whois 数据访问或发布 | 已关闭（T001）；执行授权另行取得 |
| 9 step/TUI/bootstrap | step/TUI 合同与 TUI 原型已验证；seed/oracle 身份仍未建立 | T017 |
| 10 多语言 schema | C+Go 与 C+JavaScript+Python 的所有权/依赖 fixture 已通过 | 已关闭（T003） |
| 11 handoff | 合法交接与只读 target 非法交接 fixture 已通过 | 已关闭（T003） |
| 12 documentation | 合法 impact/waiver、必需文档缺失及循环生成 fixture 已通过 | 已关闭（T003） |
| 13 产品与生命周期 | PC-01–PC-06 对应记录均有独立结构正反 fixture；运行验收仍属 S1 exit | 已关闭 S0 契约部分（T003） |

允许下一步：评审本包并执行 S0 规格化和独立技术探针。进入 S1 必须逐项提交证据与显式用户批准，不能由文档存在性自动放行。

## 4. T001 所有者决策单

状态：`APPROVED_FOR_S0_SPEC_FREEZE`。2026-09-06，产品所有者（当前仓库用户）基于第三方只读审计明确批准 S0 规格冻结及本决策单内容；T001 完成，允许进入 T002。该记录不是独立安全签署，也不授权 S1 生产编码、付费调用、具体 whois 数据访问、Git 操作或发布。

| 决策 | 批准内容 | 当前约束/待落实 | 批准记录 |
|---|---|---|---|
| S1 范围 | R1–R22 加 PC-01–PC-06；S2/S3 增强不前移 | 仅批准范围进入规格冻结，未授权 S1 生产编码 | 2026-09-06，产品所有者批准 |
| S0 预算 | 只使用当前已授权会话和本地工具；新增付费调用为 0 | 不调用收费模型/服务；自动安装依赖仍需任务内明确依据 | 2026-09-06，产品所有者批准 |
| 产品所有者 | 当前仓库用户作为最终范围/发布批准者 | S1 编码、Git、发布仍须各自明确授权 | 2026-09-06，产品所有者确认 |
| 安全评审 | 指定一名未生产该候选的独立评审者；同一人兼任须披露独立性不足 | 评审者尚待指定，§17.3 保持 NOT_READY | 原则于 2026-09-06 批准；人员/结论待记录 |
| 平台/支持 | Windows 11 amd64 正式；Linux amd64 核心 CI；当前与前一 minor，前一 minor 90 天 | 首版前冻结 EOL 通知方式和例外 | 2026-09-06，产品所有者批准 |
| 发布信任 | S1 使用固定 commit SHA、GitHub CI、SHA256SUMS、SBOM 和许可证清单；不建立离线密钥体系 | SHA256SUMS 不得表述为身份认证；签名/attestation 后续增强 | 2026-09-08，产品所有者批准简化策略 |
| 安装模型 | Windows amd64 便携 ZIP 手工解压到用户自选独立目录；显式路径运行，不修改 PATH/注册表/系统目录 | 最终 ZIP 仍需绑定发行记录并复验；不授权发布 | 2026-09-09，产品所有者批准 |
| whois 影子 | 原则允许未来使用脱敏只读固定夹具 | 每次执行仍需输入路径、脱敏范围和时段授权；当前不得读取资产 | 2026-09-06，产品所有者原则批准 |
| 依赖许可 | 优先标准库；每个 TOML/Schema/canonical/TUI 候选逐项核许可证、最低 Go 与维护状态 | T002 记录名称/版本/来源/许可证/摘要/日期后方可引入 | 2026-09-06，产品所有者批准策略 |

本次明确批准的是“批准 S0 规格冻结”。“批准 S1 编码”仍须另行给出，二者不能互相推断。本次批准不自动批准 Git commit/push/release，也不批准读取具体 whois 资产或新增付费调用。