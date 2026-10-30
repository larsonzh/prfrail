## 判定
| 项 | 判定 | 理由（≤2 行） |
|---|---|---|
| C1 | BLOCKED | 唯一 Wait、独立停机截止及外部调用不持 run 锁，覆盖重复回收、无界回收等待和持锁重入问题；超时不提前释放所有权。[证据](docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/01-sol-blueprint/review.md#L2) |
| C2 | PARTIAL | 已 spawn 故障保留所有者、禁止自动抢槽，能阻止错误重试；但 spawn 后才加入 Job，宿主崩溃窗口仍可能产生孤儿，见 N1。[证据](docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/01-sol-blueprint/review.md#L10) |
| C3 | PARTIAL | 单文件读取与 PID 查询绑定同一对象有效；目录 helper 返回即关闭句柄，后续路径操作的安全仍依赖未展开的全路径 ACL 边界，见 N2。[证据](docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/01-sol-blueprint/review.md#L24) |
| C4 | BLOCKED | 原件闭包耐久先于 intent，恢复与加载均验证原件；缺口不得补造、覆盖赢家或继续结算，堵住“收据完整但证据不成立”的放行路径。[证据](docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/01-sol-blueprint/review.md#L32) |
| C5 | PARTIAL | 共享 outcome、原 intent 恢复及冻结候选复用能防重复采集与重入漂移；但 `Reopen(..., 1, parent)` 写死 attempt，与执行身份包含 attempt 的设计矛盾，见 N3。[证据](docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/01-sol-blueprint/review-batch-b.md#L5) |
| C6 | BLOCKED | 就生产预算路径而言，固定 scope、稳定跨进程锁、锁内重读、耐久确认及 unknown 保留占用，覆盖快照双写、确认丢失和重启释放额度问题。[证据](docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/01-sol-blueprint/review-batch-b.md#L21) |
| C7 | PARTIAL | 区分 F、真实验收和生产发布，能消除 AT-23 自依赖；但独立 F 只规定 Permit 校验，未规定授权额度的消费与并发累计，见 N4。[证据](docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/01-sol-blueprint/review-batch-b.md#L41) |

## 新漏洞
- `N1` 宿主在 `cmd.Start` 成功后、加入 Job 前崩溃 → 内存 registry 随宿主消失，未入 Job 的子进程可能成为孤儿，即使仍处于挂起态；保留槽位只能阻止再启动 → 明确创建时即建立 Job 归属，或提供可证明的恢复回收路径；增加该窗口的宿主崩溃测试，未闭环前保持 Windows `unproven`。
- `N2` `RunEvidenceDir` 返回后，候选仍有任一祖先目录的替换权限；仅禁止替换请求目录本身不足以排除此条件 → 后续按路径打开可能转向另一对象或根外，前面的句柄核验不能约束此次 I/O → 将目录绑定延续至实际 I/O，或证明 ACL 隔离覆盖全部相关路径组件；测试必须包含 helper 返回后的替换。
- `N3` `attempt>1`，尤其同时存在 attempt=1 的旧工作区 → 固定传入 `1` 会选错重开身份；即使后续校验拒绝，也会导致合法恢复失败，不能认定证据与执行上下文一致 → 使用持久执行上下文中的真实 attempt，并增加同任务不同 attempt 的工作区、候选和证据隔离测试。
- `N4` 同一 `ProbePermit` 被重复或并发用于 F，而 `Require` 仅校验授权参数 → 每次调用都可重新获得 `maximumRequests` 额度，绕开生产 dispatcher 后缺少累计预算约束 → 明确 Permit 一次性消费或授权窗口共享计数，调用前原子预留；失败、超时与重启不得重置已消费额度。

## 过度设计
- 未发现仅凭这两份施工图即可认定的必然过度设计，也没有足够排期依据断言无法在 11 月落地 → 保留必要的所有权、封存和持久预算职责，将 `ArtifactCommitter`、`DurableCostLedger` 限定为本期固定流程与已证明后端；不追加通用事务框架，Windows 耐久升级继续留在独立契约修订范围。

## 结论
- VERDICT: REVISE，启动归属、目录绑定、attempt 传递及探针额度消费仍有未闭合反例；以上为施工图静态判定，[拟新增测试及 lint 均未执行](docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/01-sol-blueprint/review-batch-b.md#L59)。
