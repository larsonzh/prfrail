# 证轨 P38 报告：T027-B 代码级重构施工图（准则 v1.59）

切片 `P38-T027B-TACTICAL-BLUEPRINT`（评估类 `[SLICE]`）；执行日 2026-10-30 至 2026-10-31；准则 v1.59；两轮推演由只读载体执行、主控逐字转录；R = 0、D = 0。

## 1 执行摘要

- **调用事实**：第 1 轮（Sol）按用户裁定拆为两批执行（批次 A = C1–C4，批次 B = C5–C7），第 2 轮（Astra）为红队验收；可用正文 3 份，另有 1 次首派被判作废。
- **转录与校验**：三份正文均逐字转录，`sha256` 源件与落盘件一致；批次 A 源 10,052 字符、批次 B 源 9,808 字符、第 2 轮源 2,362 字符。
- **批量辅助**：实派 2 次（输入包统计、证据包一致性复算），均为 `deep-reasoner` + `GLM-5.3-Flash (glm)`，逐次机械复核通过。
- **核心结论**：施工图覆盖 C1–C7，逐项给出文件与函数级改法、新增测试与断言点；Astra 判 `VERDICT: REVISE`——3 项 BLOCKED（C1 / C4 / C6）、4 项 PARTIAL（C2 / C3 / C5 / C7），并给出新漏洞 N1–N4。

## 2 两轮回顾

### 2.1 第 1 轮 · 施工图（Sol，分批）

- **批次 A（C1–C4）**：C1 以每进程唯一 reaper 与有界停机替换无界回收等待；C2 以 `ProcessSpec.OnSpawn` 与预登记 starting 所有者消除「已 spawn 却按未启动重试」的路径；C3 把限长读取与 PID 读取绑定同一已开对象；C4 以「原件封存先于 intent」重排终局提交顺序。
- **批次 B（C5–C7）**：C5 以执行上下文摘要 `K=(run,task,step,attempt,requestId,requestHash)` 选择不可变上下文并共用 `outcomeOnce`；C6 以持久包装 `DurableCostLedger` 提供 `Reserve`、`RequireOutstandingReservation`、`Settle`；C7 以 `EvaluateT027BReadiness` 分离探针、真实验收与生产发布，消除 AT-23 自依赖。
- **跨项约束**：两批各给出一组不变量（所有权常驻、核心单写、锁不嵌套、封存至结算至 closure 的顺序）。

### 2.2 第 2 轮 · 红队验收（Astra）

- **判定**：C1 / C4 / C6 为 BLOCKED；C2 / C3 / C5 / C7 为 PARTIAL，分别残留宿主崩溃窗口、目录替换窗口、attempt 写死与探针额度可重复消费。
- **新漏洞**：N1 为 Job 归属建立前的宿主崩溃窗口；N2 为目录 helper 返回后的祖先替换；N3 为 `Reopen` 写死 attempt；N4 为 `ProbePermit` 额度可重复取得。
- **过度设计**：未认定必然过度设计；建议把 `ArtifactCommitter` 与 `DurableCostLedger` 限定在本期固定流程与已证明后端，不追加通用事务框架。
- **契约阻塞**：`CONTRACT-BLOCK`——Windows 平台发布耐久升级须先修订 `docs/CONTRACTS.md` §7 并经独立评审。

## 3 输入包与计量

- **输入包**：条目 59（3 件文档段 + 56 个 `.go`），行数 14,714，字节 648,707（约 634 KB）；逐件明细见 `01-sol-blueprint/input-list.md`。
- **第 2 轮输入收窄**：仅两份施工图正文，不附代码包（用户指令，节省推理时间）；按 v20 规则在本报告声明「收窄」。
- **输出计量**：批次 A 源字节 13,442、批次 B 源字节 14,583、第 2 轮源字节 4,242，均为无 BOM 的文本长度。

## 4 待用户裁定

| # | 项 | 现状 | 阻塞谁 |
|---|---|---|---|
| 1 | Astra 判 REVISE 的 N1–N4 是否在本片内补一轮修订 | 待裁 | 11 月实施 |
| 2 | `CONTRACT-BLOCK` 是否立项为契约修订片 | 待裁 | C4 与 C6 落地 |
| 3 | C2 / C3 / C5 / C7 的 PARTIAL 是否可接受为 11 月起点 | 待裁 | T027-B 实施 |

## 5 流水线裁剪与工序声明

- **裁剪依据**：本片为评估类 `[SLICE]`，无实现 diff、无实现审查对象。
- **批量辅助**：实派 2 次，场景为步骤 0 输入包统计与证据包一致性复算；未派排序、择优或任何结论性工作。
- **声明口径**：下方两条为 G7-f 的跳过通道；本报告不携带 `sixstage` 围栏。

- ⑥ 跳过：本片为评估类切片，无实现 diff。
- ⑦ 不适用：本片无实现审查对象；第 2 轮 Astra 仅作红队审查，不产出 ⑦ 判定。

## 6 制品登记

变更集为 15 件：准则对（头部升 v1.59、§0.3 日志行、§12.4 的 OB-131）、台账对（P38 的台账行与其切片块）、报告对、证据包 9 件。

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/t027b-tactical-blueprint.md	repo
docs/validation/t027b-tactical-blueprint_EN.md	repo
docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/00-doc-inventory.md	repo
docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/freedom-list.md	repo
docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/MANIFEST.md	repo
docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/SHA256SUMS.txt	repo
docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/01-sol-blueprint/input-list.md	repo
docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/01-sol-blueprint/review.md	repo
docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/01-sol-blueprint/review-batch-b.md	repo
docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/02-astra-redteam/input-list.md	repo
docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/02-astra-redteam/review.md	repo
```

## 7 异常处理 / 事实包

1. **档位**：用户侧设定 1M / Max；主控不可观测，按事实登记。transcript 偏移作为「调用已发生」的证据：批次 A 派单约 L101607，批次 B 派单约 L101855，第 2 轮派单约 L102004，记录时 transcript 总行数 102,123；该偏移**不**主张证明档位。
2. **输出长度阻力**：第 1 轮首派连续两次触发自重启与截断，原始件字节 42,509 与 15,428 已留证但未作底稿；改为分批后批次 A 仍含一次自重启，主控只取第一遍完整段，第二遍作废。
3. **止损**：按用户裁定在批次 A 首次完整返回后即停止重试，未再消耗额度；第 2 轮输入按裁定收窄为两份施工图。
4. **无其它异常**：无越权调用、无临时件漏报、无契约与代码改动。
