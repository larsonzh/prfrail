# P49 RACE-LEG-EXPAND 起跑消息（原文落盘）

> 落盘纪律：v38 §七 第 9 条「起跑消息即工件」。以下为用户起跑消息原文，未改写（附件 `CHECKLIST_v38.md` 为仓库文件，不复制入包）。

---

切片：P49 RACE-LEG-EXPAND CI Race 腿扩包（M 档 · 基础设施类 [SLICE]）
=====================================================================

【〇·零 起跑前自检与上下文锚定】
1. 请先完整阅读 `docs/t027/CHECKLIST_v38.md`（跨会话快照，注意以实际仓库 HEAD 为准）。
2. 核对第一节的 5 个锚点（HEAD 应为 `81dd49f`，若生成了本消息提交则顺延）。
3. **主控模型声明**：本会话主控模型维持 **`DeepSeek V4 Pro`**。
4. **记忆唤醒（红线复述）**：请在回执中复述以下五条规则，证明你已掌握：
   - 「评估/实现类切片 ⑨/⑧c 全绿后必停 ⑩」：未获同轮授权绝不自行 `commit`/`push`。
   - 「纪律 S（会话长度与压缩纪律）」：上下文达 100K 或出现压缩提示，立即准备结片。
   - 「纪律 T/U」：防字面量误杀（禁写量词形态、禁逐字引用触发串）。
   - **v38 §九 P48 门禁坑**：特别注意 `tmp/` 下的源码副本必须用 `.txt` 后缀，防 `go build` 污染；写 `.md` 必须显式补 BOM。
5. 若锚点 ≥2 项不符，立即停止并报仓库实况。

【〇 本片起跑前的裁定】
1. **R-03 裁定**：批准 `RACE-LEG-EXPAND` 独立立项，作为层 1 治理收尾。
2. **目标**：把 CI 的 Race 腿（`CGO_ENABLED=1 go test -race`）从当前的 `internal/adapters` 与 `internal/chain`，**扩包到 `internal/console`、`internal/gates`、`internal/guard`、`internal/tickets`**。
3. **风险定级**：基础设施类。涉及 `.github/workflows/ci.yml` 的修改。绝对不能手滑删除现有的 CI 步骤。

【一、授权边界（流水线角色与模型）】
- ① 架构：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（评估 cgo 在 windows/ubuntu 两腿的构建影响与时长）。
- ② 实现：`prfrail-implementer` + `DeepSeek V4.1 Flash (deepseek)`（仅改 `ci.yml`）。
- ③ 测试：`prfrail-tester` + `DeepSeek V4.1 Flash (deepseek)`（本地跑通扩包后的 race 测试）。
- ⑤ 预审：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（额度 1+1）。
- ⑥ 独立扫描：`independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`（额度 3）。
- ⑦ 终审：`independent-reviewer` + `GPT-5.3-Codex (copilot)`（额度 1+2）。
- ⑧ 文档：`prfrail-documenter` + `DeepSeek V4.1 Flash (deepseek)`。
- ⑨ / ⑧c / ⑩：由主控（V4 Pro）执行。

【二、派单与执行硬性要求】
- 每次 `runSubagent` 必须显式同时传 `agentName` 与 `model`。
- **严格串行**：禁止并行调用子代理。
- 主控不得代写任何 ②③ 代码或测试。

【三、批量辅助】
实派 `deep-reasoner` + `GLM-5.3-Flash (glm)` ≥ 2 次（场景：当前 `ci.yml` 结构检索、各包 Race 耗时评估），严格串行。

【四、临时脚本与中间产物】
一律落盘 `tmp/RACE-LEG-EXPAND/`；只增不删；⑧c 统一清理。

【五、步骤 0 待办】
1. 建目录 `tmp/RACE-LEG-EXPAND/` 与 `docs/validation/evidence/RACE-LEG-EXPAND/`。
2. 起跑消息逐字落盘为 `00-briefing.md`（BOM+LF）。
3. 在 `REMAINING_SLICES{,_EN}.md` 登记台账（状态设为 `🔄 进行中`，`R/D=0`，临时件落点）。
4. 落盘 `freedom-list.md`、`MANIFEST.md` 与 `SHA256SUMS.txt`。
5. `artifacts` 围栏声明 `00-briefing.md`。
6. 核对前置阅读：`docs/validation/p45-e2e-zero-model.md`（其中提及了 Race 腿的未覆盖部分）。

【六、执行流水线与核心要求】
1. **① 架构**：检查当前 `.github/workflows/ci.yml` 中 `Race` 腿的配置。评估 Windows 环境下的 cgo 依赖，输出《CI 扩包方案》（明确要加的目录，以及 `ci.yml` 的精确修改内容）。
2. **③ 测试先行（TDD 强制）**：**必须先在本地未修改的 `ci.yml` 结构上，写一个脚本或简单断言证明“当前 CI 未覆盖 `internal/console` 等包的 Race”。**（例如：本地执行排除 `--race` 的测试或检查 `ci.yml` 对应的参数行）。
3. **② 实现**：修改 `ci.yml`，将 `--race` 的目标包扩展为 `./internal/adapters/... ./internal/chain/... ./internal/console/... ./internal/gates/... ./internal/guard/... ./internal/tickets/...`。如果因环境或超时问题无法全部扩包，必须在报告中明确说明并降级。
4. **④ 门禁**：`node tools/gates/gate.js --all --scope=tree` 全绿；`actionlint`（如果有）或 YAML 语法检查通过。
5. **⑤⑥⑦ 审查**：重点关注 `ci.yml` 的修改是否安全、是否会引发超时、`Race` 腿的 cgo 参数是否正确。

【七、执行节奏与红线】
1. ④ 全绿直接派 ⑤；⑤/⑥/⑦ 报 Medium+ 原位重跑该环。
2. **⑩ 是唯一常规停点**。⑨/⑧c 全绿后，**必须停在 ⑩ 等授权，严禁自行 commit/push**。
3. 严格遵守 v38 §九 的六条门禁坑（防 `.go` 污染、补 BOM 等）。

【八、文档同步更新（v38 第 15 条铁律）】
本片修改了 CI 架构，在 ⑧ 阶段，**必须同步更新**：
1. `docs/ARCHITECTURE.md`（及 `_EN`）：在 §3 数据流或 §4 存储中，若有提及测试设施，补充说明现在所有核心包已全部纳入 Race 腿覆盖。
2. `docs/DEV_PLAN.md`（及 `_EN`）：检查是否有“CI 覆盖不全”的过时描述并更新。
3. 台账回写：完成记录中体现文档更新状态。

【九、⑩ 回执要求】
1. 派单表（`agentName` 与 `model` 逐字机械读取）。
2. 红绿证据（修改前的 `ci.yml` 片段证明未包含，修改后的 YAML 文件展示）。
3. 三作用域门禁与 `SELFTEST` 读数。
4. **文档同步更新**的 `git diff --stat`。
5. 最终事实包，含未闭环项。

现在开始执行 P49 步骤 0。
