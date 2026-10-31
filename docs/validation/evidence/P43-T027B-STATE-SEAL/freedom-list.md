# P43-T027B-STATE-SEAL · 执行自由度清单（步骤 0 基线）

本片为 `[SLICE]` 实现类。以下为执行角色在片内可自由裁量的范围与不可动项。

## 可自由裁量
1. 新增文件的命名与拆分（`internal/chain/**`、`internal/tickets/**` 内），须与既有命名风格一致。
2. 新增测试文件的组织方式（表驱动 / 单例），须保持 ③ 先行且不得修改既有断言强度。
3. 内部辅助函数与类型的具体形状（不改变导出签名的前提下）。
4. 注释与技术说明的措辞（英文，遵循仓库 i18n 约定）。

## 不可动（硬边界）
1. 导出签名：`RunEvidenceDir`、`StepIntentPreparer`、`WorkspacePort`、`StepPort`、`Reviewer`、`Publisher`、`Acceptance` 的既有方法集不得删改；新增能力须以新增方法或新接口叠加。
2. 事件模型与契约：`docs/CONTRACTS.md` 定义的事件字段、闭集状态、`InputEvidence` 语义不得变更；若确需变更须走「协议先行」并在 ⑧ 阶段回写契约。
3. 不得改 `internal/adapters/**`、`internal/gates/**`、`internal/console/**`、schema、fixtures、`ci.yml`。
4. 不得降低任何既有测试断言强度，不得跳过测试（`t.Skip`）以转绿。
5. 不得提交 / 推送（停 ⑩ 等授权）；不得推 gitee。
6. 临时件只落 `tmp/P43-T027B-STATE-SEAL/`，片内只增不删。

## 判定冲突时的优先序
1. `docs/CONTRACTS.md`（协议）＞ 2. 本片 ① 架构方案（含修订注记）＞ 3. 既有测试断言 ＞ 4. 代码现状。
