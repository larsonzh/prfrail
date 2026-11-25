# ③ 测试报告 ③a 原文（P85 · T2 反例与变异）

### 1. A 组：修复既有 16 个测试（仅补装配输入，未改判定意图）
根因（②b-1 差异第 2 条）：夹具无 `PolicySnapshot` / 无冻结制品，且 grant 的 `PolicyHash` 是任意摘要，T2 后无法与磁盘制品对证。
统一修法（新增夹具助手）：先冻结策略体 → 用其摘要作为 grant 的 `PolicyHash`；装配类夹具补 `PolicySnapshot`；`p83BindReviewedPolicy` 改为「先冻结再绑定」；`TestP83T2a` 刻意让 newer/widened 两份 grant 携带同一冻结摘要。

### 2. B 组：新增反例（`internal/console/p85_t2_frozen_policy_tester_test.go`）
- CX1 `TestP85T2ReviewRefusesMissingPolicyArtifact`（反例·缺失，含正控制）
- CX2 `TestP85T2ReviewRefusesUnreadablePolicyArtifact`（3 子例：畸形 JSON / body 类型错 / 未知字段严格解码）
- CX3 `TestP85T2ReviewRefusesPolicyArtifactThatIsNotTheGrantedOne`（反例·合形不匹配）
- CX4 `TestP85T2PublishRefusesArtifactReplacedAfterReview`（替换 / 删除 2 子例）
- CX5 `TestP85T2AssemblyRefusesWithoutFrozenPolicy`（反例·装配期缺失）
每条均断言：不 approve / 不发布（零值）、证据完整（制品字节前后逐字节相同）、因果唯一性（CX4 拒绝后回写原制品 ⇒ completed；CX1/CX3 带正控制）。

### 3. C 组：变异检验（实测红证据）
基线 sha256（`runtime_agent.go`）= `BD88D96A…112E739F`
- 变异 1（删 `operatorApprovedPolicyHashForAttempt` 的制品比对）⇒ CX3 转红（`--- FAIL ... Outcome:approve PolicyHash:sha256:0965b7f4…`）
- 变异 2a（仅删 `verifyReviewPolicyBinding` 的一次磁盘重读）⇒ CX4 仍绿（第二处磁盘重读属双保险，比方案书更强，非缺陷）
- 变异 2b（连带把 `verifyFrozenPolicyArtifact` 改为不问磁盘结果）⇒ CX4 两子例转红（`Outcome:completed ... ErrorEvidence:[]`）
- 恢复：逆向编辑还原（未用 git checkout）⇒ sha256 回到 `BD88D96A…112E739F`，gofmt 空，测试复绿。

### 4. 验证输出
`ok github.com/larsonzh/prfrail/internal/console 4.756s`（16 既有 + 6 新增全 PASS）；全量剩 5 个 T3 相关失败（归 ③b），T2 范围全绿。

### 5. 边界
未跑 git add/commit/push；实现文件与变异前逐字节一致。不可测项：无。
