# 证轨 P37 报告：T027-B 生产级架构设计（准则 v1.58）

切片 `P37-T027B-BREAKTHROUGH`（评估类 `[SLICE]`）；执行日 2026-10-30；准则 v1.58；两轮均由只读载体 `independent-reviewer` 执行、主控逐字转录；R = 0、D = 0。

## 1 执行摘要

- **调用事实**：两轮连续执行、不分批。第 1 轮 `GPT-6.1 Sol (copilot)`，第 2 轮 `GPT-6 Astra (copilot)`（表外授权见 §12.4 的 OB-130）。总调用 2 次，未超「≤ 3 次」上限。
- **转录与校验**：两份 `review.md` 均逐字转录，`sha256` 与原消息归一形态一致（仅 BOM 与换行归一差异）；第 1 轮行数 111、字节 10723，第 2 轮行数 68、字节 9277。
- **批量辅助**：实派 2 次，均为 `deep-reasoner` + `GLM-5.3-Flash (glm)`；场景 ①/②（输入包统计与数字核对）与场景 ②（两轮转录数字核对），逐次机械复核通过。
- **核心判断**：第 1 轮给出「G0–G4 门序 + 模块划分 + 状态机 + 关键接口 + 资源约束」的完整设计；第 2 轮指出三处致命交接缺口（启动所有权、证据耐久、门序前置），并给出最小重构与 11 月最小边界。**两者合起来才是可落地的设计基线**。

## 2 两轮回顾

### 2.1 第 1 轮 · 架构设计（Sol）

- **门序**：G0 立项（OB-55 归属与冻结项）→ G1 耐久（按平台 × 前提举证）→ G2 边界（受限候选可运行）→ G3 装配（零模型三任务）→ G4 真实闭环。门与门之间以证据绑定，不以日期替代条件。
- **模块**：`chain` 独有调度与接受决定；`guard` 持有受管进程与停机；adapter 只转换外部事实；SessionBridge 不进入主执行链。
- **状态机**：completed 仅允许 step PASSED；未知一律落在 uncertain / paused，禁止以 completed 表达未知。
- **接口**：复用 admission / dispatch / postflight 原签名，新增 `Prepare` / `WaitDispatched` / `Decode` 三个内部接口，均要求同绑定幂等、失败保留底层分类。
- **资源**：固定候选版本、单一 watchdog、持久预算与 unknown 占用、凭据最小注入、网络默认拒绝。

### 2.2 第 2 轮 · 批判与重构（Astra）

- **三处致命**：① 启动所有权交接（spawn 后身份发布失败可能被当作「未启动」，回收失败又释放槽位）；② 完成记录耐久不等于其引用证据原件耐久；③ G1 隐含本轮边界外的生产前置、G2 又提前需要真实调用，原门序不能直接当作 11 月承诺。
- **最小重构 6 条**：所有权与短期互斥分离、Prepare 结果改为不可复用受管资源、安全决定绑定打开对象、增加证据提交屏障、按 request 摘要选择不可变上下文、持久预留与结算确认接到调用与发布边界。
- **降级 3 条**：拆分「开发顺序」与「生产准入顺序」；把通用 Wait 重入收窄为终局事务恢复；把完整 session resume 降为后续能力。
- **11 月最小边界**：一个 Windows 前提 P、一个候选、一任务一 Agent 步骤 + 声明 hook；先用零模型替身验证装配与串行接受；真实三任务闭环另设条件验收，替身不得冒充真实证据。

## 3 输入包与计量

- **第 1 轮输入**：62 件（文档段 6 件、Go 源 56 件），行数 14923、字节 712662；其中 `internal/adapters` 38 件 / 395919 字节、`internal/chain` 18 件 / 184724 字节。逐件清单见 `01-sol-design/input-list.md`。
- **第 2 轮输入**：第 1 轮设计书全文 + 同一原始输入包（逐件重复投递，避免上下文漂移）。
- **模型调用**：2 次（上限 3 次）。**批量辅助**：2 次，按 §2.2 不计入该上限。
- **口径**：整文件按原始字节（含 BOM）；行区间按切片后以 LF 重建计字节，不含原文件 BOM。

## 4 待用户裁定

1. OB-55 是否归入独立 T027-B 实施片，以及实施片边界。
2. 冻结哪个候选与 Windows 前提 P；是否批准同边界耐久与真实探针取证。
3. AT-23 三单元映射与签署责任人；首验是否覆盖 session resume。
4. 资源、费用与人工介入上限，以及试点受阻后的停止或延期条件。
5. 是否接受「先交限定机制与装配片，真实三任务闭环另设条件验收」。

## 5 流水线裁剪与工序声明

- ⑥ 跳过：本片为评估类切片，无实现 diff。
- ⑦ 不适用：本片无实现审查对象；第 2 轮 Astra 仅作跨模型架构批判，不产出 ⑦ 判定。
- ⑧a 转录：两轮由主控逐字转录，`sha256` 与原消息归一形态一致。
- ⑨：编码门禁 + 三作用域 + `SELFTEST` 见 ⑩ 回执读数。
- ⑧c：临时目录 `tmp/P37-T027B-BREAKTHROUGH/` 清理并做三项零残留回读。
- ⑩：停点，等 commit / push 授权。

## 6 制品登记

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/t027b-breakthrough.md	repo
docs/validation/t027b-breakthrough_EN.md	repo
docs/validation/evidence/P37-T027B-BREAKTHROUGH/00-doc-inventory.md	repo
docs/validation/evidence/P37-T027B-BREAKTHROUGH/freedom-list.md	repo
docs/validation/evidence/P37-T027B-BREAKTHROUGH/MANIFEST.md	repo
docs/validation/evidence/P37-T027B-BREAKTHROUGH/SHA256SUMS.txt	repo
docs/validation/evidence/P37-T027B-BREAKTHROUGH/01-sol-design/input-list.md	repo
docs/validation/evidence/P37-T027B-BREAKTHROUGH/01-sol-design/review.md	repo
docs/validation/evidence/P37-T027B-BREAKTHROUGH/02-astra-critique/input-list.md	repo
docs/validation/evidence/P37-T027B-BREAKTHROUGH/02-astra-critique/review.md	repo
```

## 7 异常处理 / 事实包

1. **档位**：用户侧设定 1M / Max；主控不可观测，按事实登记。transcript 偏移作为「调用已发生」的证据：第 1 轮派单约 L98894–L98898，第 2 轮派单约 L99179–L99183，记录时 transcript 总行数 99502；该偏移**不**主张证明档位。
2. **批量辅助**：实派 2 次（场景 ①/②、场景 ②），逐次机械复核通过；未派排序、Top-N 择优与任何结论性工作。
3. **⑨ 期**：门禁红一律在授权边界内改写后重跑，过程与结论见 ⑩ 回执；未绕过任何判据。
4. **无其它异常**：无越权调用、无临时件漏报、无契约与代码改动。
