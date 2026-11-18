# P81 ① 架构方案书 r1（`DIRECTIVE-MODEL-UPGRADE`）

> 载体 / 模型：`deep-reasoner` / `DeepSeek V4 Pro (deepseek)`（只读）。逐字转录自 ① 的最终消息。
> 编码：UTF-8 with BOM + LF。

## Q1 落点全表（HEAD=v1.69；CN 与 EN 行号逐点相同，已逐处实读）

**现役落点（要改，9 处 ×2 语言）**：

1. `docs/DELIVERY_DIRECTIVE.md:3` 头部
OLD: `版本：v1.69。日期：2026-11-15。`
NEW: `版本：v1.70。日期：2026-11-18。`
EN:3 OLD: `Version: v1.69. Date: 2026-11-15.` NEW: `Version: v1.70. Date: 2026-11-18.`

2. CN:302 §2.2 ⑥ 行（整行）
OLD: `| 独立扫描员 | ⑥ | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | `mai-code-1.1-flash` | 否 | 只读 |`
NEW: `| 独立扫描员 | ⑥ | `independent-reviewer` | `MAI Code 2 Flash (copilot)` | `mai-code-2-flash` | 否 | 只读 |`
EN:302 OLD: `| Independent scanner | ⑥ | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | `mai-code-1.1-flash` | No | Read-only |`
NEW: `| Independent scanner | ⑥ | `independent-reviewer` | `MAI Code 2 Flash (copilot)` | `mai-code-2-flash` | No | Read-only |`

3. CN:304 §2.2 ⑦ 深度档行（整行）
OLD: `| 独立终审员（深度档） | ⑦ | `independent-reviewer` | `GPT-6.1 Sol (copilot)` | `gpt-6.1-sol` | 否 | 只读 |`
NEW: `| 独立终审员（深度档） | ⑦ | `independent-reviewer` | `GPT-6.2 Sol (copilot)` | `gpt-6.2-sol` | 否 | 只读 |`
EN:304 OLD: `| Independent final reviewer (deep tier) | ⑦ | `independent-reviewer` | `GPT-6.1 Sol (copilot)` | `gpt-6.1-sol` | No | Read-only |`
NEW: 同构换 `GPT-6.2 Sol (copilot)` / `gpt-6.2-sol`。

4. CN:309 §2.2「⑦ 三档分档」句（行内仅 1 处）
OLD: `gpt-6.1-sol`（深度档）
NEW: `gpt-6.2-sol`（深度档）
EN:309 OLD: `gpt-6.1-sol` (deep tier) NEW: `gpt-6.2-sol` (deep tier)

5. CN:560 §5.2 ⑦ 分档表「深度档」行（行内仅 1 处）
OLD: `| 深度档 | `gpt-6.1-sol` |`
NEW: `| 深度档 | `gpt-6.2-sol` |`
EN:560 OLD: `| Deep tier | `gpt-6.1-sol` |` NEW: `| Deep tier | `gpt-6.2-sol` |`

6. CN:640 §6.3 G5-b 闭集句（行内仅此两枚子串；判据表述与比较规则逐字不动）
OLD: `mai-code-1.1-flash`、`gpt-6.1-sol`
NEW: `mai-code-2-flash`、`gpt-6.2-sol`
EN:640 同构：OLD: `mai-code-1.1-flash`, `gpt-6.1-sol` NEW: `mai-code-2-flash`, `gpt-6.2-sol`（该句用 `, ` 分隔，闭集其余 6 成员不动）

7. CN:784 §7.5 成本表 Sol 行（只换 modelId，单价不动）
OLD: `| `gpt-6.1-sol` | 200 | 1000 |`
NEW: `| `gpt-6.2-sol` | 200 | 1000 |`
EN:784 同（`| `gpt-6.1-sol` | 200 | 1000 |` → 换 modelId）

8. CN:865 §8.2 档位映射行
OLD: `| `gpt-6.1-sol` | **Extra High**（允许 **Max**） |`
NEW: `| `gpt-6.2-sol` | **Extra High**（允许 **Max**） |`
EN:865 OLD: `| `gpt-6.1-sol` | **Extra High** (**Max** allowed) |` → 换 modelId

9. CN:867 §8.2「取档授权」句（行内 1 处）：OLD: `` `gpt-6.1-sol` 取 `Max` `` NEW: `` `gpt-6.2-sol` 取 `Max` ``；EN:867 OLD: `` with `gpt-6.1-sol` `` NEW: `` with `gpt-6.2-sol` ``

10. CN:1529 附录 B.6 ⑥ 载体行（整行）
OLD: `⑥ 载体：runSubagent(agentName='independent-reviewer', model='MAI-Code-1.1-Flash (copilot)')`
NEW: `⑥ 载体：runSubagent(agentName='independent-reviewer', model='MAI Code 2 Flash (copilot)')`
EN:1529 OLD: `⑥ Vehicle: runSubagent(agentName='independent-reviewer', model='MAI-Code-1.1-Flash (copilot)')` NEW: 换 `MAI Code 2 Flash (copilot)`

**历史行/冻结件（不改，§11.1 不追溯）**：CN/EN 同号 L65（§0.3 v1.18 日志行）、L1091/L1100/L1106/L1170（§12.4 OB-53/62/68/132 台账行，记录当时事实）、L1700/L1720/L1741（附录 C 既有台账块审查实录行）、L1802（C.0j「为何必需」）。均属历史登记，只增不溯改。

## Q2 实读项闭合

- §7.5 成本表（L781–786）只列 ⑦ 三模型，全准则无 ⑥ 单价行；⑥ 仅 §7.5 额度表「上限 3 次/片」与 §7.6 探针额度，无成本引用 ⇒ **MAI Code 2 Flash 无成本表引用需同步，确认**。
- Sol 行只换 modelId、200/1000 不动，确认。

## Q3 selftest 同步面

① 必同步夹具值：`tools/gates/testdata/fixtures/g5b-closed-set-legal-sol.txt:4`
OLD: `model: gpt-6.1-sol` NEW: `model: gpt-6.2-sol`（T41 的 PASS 语义不变）。**明确不改**：`g6-1-pending-done.txt`（G6-1 历史行**形态**夹具，其 MAI-Code-1.1-Flash 是历史转写，与闭集无关）。
② 需新增（建议全量，属本片「同步自测夹具」自由度）：`g5b-illegal-retired-sol.txt`（`model: gpt-6.1-sol` ⇒ FAIL）、`g5b-illegal-retired-mai.txt`（`model: mai-code-1.1-flash` ⇒ FAIL）、`g5b-closed-set-legal-mai.txt`（`model: mai-code-2-flash` ⇒ PASS）。计数：737 = 动态 `CHECKS.length`。算术：仅① ⇒ 737 不变；①+sol 退役（+1 check）⇒ 738；再+mai 退役 ⇒ 739；再+mai 合法 ⇒ **740**。
③ 断言硬编码：全 `tools/gates/**` grep 证实除 `g5b.js` 闭集与 sol 夹具外**无**任何 `gpt-6.1-sol` / `mai-code-1.1-flash` 字符串；`selftest.js` 的 T37/T42 变异用 `G5B_FIRST_MEMBER`/`G5B_LUNA_MEMBER`，不受影响。

## Q4 P80 留痕补落

判定 = **§12.4 表尾引用块**（不占 OB 号；与 GATES-EXT / DIRECTIVE-HAIKU-BUMP 两个 ⑦ 轮次例外块同构），插入位 = HAIKU-BUMP 块之后、「按比例抽查记录」块之前；附录 C.0aw 只做指针不复述。依据：§7.8.3「用户授权例外」行要求「授权原文 + 输入收窄声明落 §12.4 / 附录 C」；v63 裁定①「正式补落 §12.4 / 附录 C 并入下次准则修订片」。

## Q5 版本与日志

§0.3 新增行（CN，插于 v1.69 行后）；受影响的在飞切片 = 本片；无其他在飞片需回写。
附录 C 新块标题：`### C.0aw P81 DIRECTIVE-MODEL-UPGRADE（v1.70）`；表格项仿 C.0av：切片 / 规模、触发、目标、交付物、范围、审查（执行片回填）、版本 v1.70（G1-b）、附注（P80 授权原文落点 = §12.4 表尾引用块，本块只做指针）。
§0.2 MINOR 依据：① 角色集三元组与流水线结构零变更 ⇒ 非 MAJOR；② 只换两处模型限定名与标识、数值与判据逻辑不动 ⇒ 属「规则增补」；③ 同型先例 v1.18（闭集成员增补，MINOR）、v1.26（GLM 扩列，MINOR）；④ 裁决注记「自 v2.0 起严格执行」下无歧义。

## Q6 三方闭集一致

替换后闭集（8 成员，顺序不动）：`deepseek-v4-pro`、`deepseek-flash`、`mai-code-2-flash`、`gpt-5.3-codex`、`gpt-6.2-sol`、`gpt-6-luna`、`glm-5.3`、`glm-5.3-flash`。
自证脚本（§12 机械读数，仿 haiku-bump 先例）：

```js
const fs=require('fs'),g5b=require('./tools/gates/lib/criteria/g5b');
const want=[...g5b.closedSet].sort().join(',');
for(const f of ['docs/DELIVERY_DIRECTIVE.md','docs/DELIVERY_DIRECTIVE_EN.md']){
 const t=fs.readFileSync(f,'utf8');
 const en=f.includes('_EN');
 const a=[...t.matchAll(/\| `((?:deepseek|mai|gpt|glm)[a-z0-9.\-]+)` \| (?:否|No) \|/g)].map(m=>m[1]);
 const m=en?t.match(/closed set = ([^(\n]+)/):t.match(/闭集 = ([^（(\n]+)/);
 const b=m[1].replace(/`/g,'').split(en?', ':'、').map(s=>s.trim());
 console.log(f,[...new Set(a)].sort().join(',')===b.sort().join(',')&&b.sort().join(',')===want?'OK':'MISMATCH');
}
```

另：G5-c 已动态读 §2.2 限定名列（改表即改判定集合）；CI `--selftest` 含 T41 载荷断言与 T37/T42 变异。

## Q7 文档同步三层

第一层（必改）：两准则文件（Q1 十点 + §0.3 行 + §12.4 引用块 + C.0aw + 头部）；`tools/gates/lib/criteria/g5b.js`（CLOSED_SET 两成员）；`g5b-closed-set-legal-sol.txt`；新增 3 夹具与 `selftest.js` 的 FIXTURES/check（若采纳）；`REMAINING_SLICES{,_EN}` P81 行的 ⑧b 完成记录（正文行不改——新旧名并存是替换定义的固有形态）。
第二层（视实读）：`.github/agents/sol-orchestrator.agent.md:74` 是 **family 字段匹配注记**（`oswe-vscode-modelD`），不硬编码 modelId ⇒ 默认不改，仅当新模型 family 变化时同步；其 `_config.*` / `_discovered_models*` **未被 git 跟踪**（`ls-files` 为空，生成器产物被忽略）⇒ 不属本片。
第三层（明确不改，理由）：`testdata/agent-runner/capability-probes/github-copilot-cli-windows-evidence.json`（历史探针证据，§11.1）；`ADR_REGISTER{,_EN}` ADR-019（历史登记；本片不新增 ADR）；`DEV_PLAN{,_EN}` 的 P40 试点历史记录；`docs/validation/*{,_EN}.md` 与 `docs/t027/CHECKLIST_v2x–v63`（历史报告与冻结快照）；`.vscode/settings.json`（未跟踪，操作员环境）；`internal/**` 产品代码。**`LEGACY_FLASH_MODEL_IDS` 不需要加旧 `mai-code-1.1-flash`**：该名单针对 DeepSeek Flash 旧 ID，旧 MAI 退役后由 §2.2 闭集即拒；若代码侧加 ID 属改产品代码 ⇒ 扩范围，不落地。

## 待升级

- 新增 3 份夹具判为 §1.7.1「已授权包内夹具同步」的执行自由度、非扩范围；若 ⑤ / ⑦ 持异见按 §1.7.2 升级。仅 ① 为必改底线。
- 未发现用户裁定与准则现状冲突。

DESIGN: DONE
