// C1 closure probe: re-run the v1.12 review's 34-item reconciliation mechanically.
// Each entry is [review item #, CN needle, EN needle]; both must be present.
const fs = require('fs');

const cn = fs.readFileSync('docs/DELIVERY_DIRECTIVE.md', 'utf8');
const en = fs.readFileSync('docs/DELIVERY_DIRECTIVE_EN.md', 'utf8');

const ITEMS = [
  [1, '**术语消歧**：本节『门禁判据』', '**Terminology disambiguation**'],
  [2, '裁决注记（2026-09-23，v1.12）', '**Ruling note (2026-09-23, v1.12)**'],
  [3, '| 版本 | 日期 | 变更 | 受影响的在飞切片 |', '| Version | Date | Change | Slices in flight (affected) |'],
  [4, '| 分级 | [ROUTINE] / [SLICE] / [PILOT]', '| Tier | [ROUTINE] / [SLICE] / [PILOT]'],
  [5, '生产代码语义边界，2026-09-23 用户裁决', 'production-code semantic boundary, user ruling 2026-09-23'],
  [6, '**计数单位**：「前 20 个字符」', '**Counting unit**'],
  [7, '可判分界 = 是否落盘', 'decidable dividing line'],
  [8, '不得以 §1.7.1 为准', 'outrank the execution freedom of §1.7.1'],
  [9, '| `modelId`（调用串，唯一权威） | UI 显示名（仅供人工核对） |', '| `modelId` (call string, sole authority) | UI display name (for manual cross-check only) |'],
  [10, "model='deepseek-v4-pro'", "model='deepseek-v4-pro'"],
  [11, 'UI 显示名不得作为调用串写入任务包或模板', 'a UI display name must not be used as a call string'],
  [12, '**黑名单含 GPT-5.x 全系**', '**The blacklist covers the whole GPT-5.x family**'],
  [13, '**闭集枚举**', '**closed-set enumeration**'],
  [14, '修改权归实现者', 'the implementer owns them'],
  [15, '**豁免（2026-09-23）**：已归档、带 MANIFEST/SHA256SUMS', '**Exemption (2026-09-23)**: archived **frozen evidence packs**'],
  [16, '| 高风险（返工≥2 或审查矛盾或原生反复失败） | **必** | 需 | 需 | **必** | **必** | **必** |', '| High risk (rework ≥2, or contradictory reviews, or repeated native failures) | **required** | required | required | **required** | **required** | **required** |'],
  [17, '⑧b 重写后若 diff 涉及', "if the ⑧b rewrite's diff touches"],
  [18, '**聚焦门禁定义（2026-09-23）**', '**Focused-gate definition (2026-09-23)**'],
  [19, '停止态切片（§1.6）不参与本判据', "does not take part in this criterion's"],
  [20, '**首轮不附任何 ⑤/⑥ 清单**', '**The first round attaches no ⑤/⑥ checklist**'],
  [21, '### 7.2.1 触发', '### 7.2.1 Trigger'],
  [22, '⑥ 独立扫描层的重发按 §7.2 成本上限', 'resends of the ⑥ independent-scan layer follow the §7.2 cost cap'],
  [23, '**计数口径（2026-09-23）**', '**Counting rule (2026-09-23)**'],
  [24, '**分母定义（2026-09-23 用户裁决，固定）**', '**Denominator definition (user ruling 2026-09-23, fixed)**'],
  [25, '**复评条款（2026-09-23 用户裁决）**', '**Re-evaluation clause (user ruling 2026-09-23)**'],
  [26, '**已处置**（v1.1）；⑦ 复审确认**未单独执行**', 'the ⑦ re-review confirmation **was not run separately**'],
  [27, '| OB-28 |', '| OB-28 |'],
  [28, '**已落地**（', '| **Landed**'],
  [29, 'G1–G7 七类机械检查全绿', 'all seven categories of mechanical checks G1–G7 green'],
  [30, '**闭集枚举表**（⑦ 终审', '**closed-set enumeration table** of §4.1'],
  [31, '切片：[切片 ID] [描述]', 'SLICE:[slice ID] [description]'],
  [32, '§2.4、§2.6、§3.1', '§2.4, §2.6, §3.1'],
  [33, '**更正注记（2026-09-23）**', '**Correction note (2026-09-23)**'],
  [34, '### C.0h', '### C.0h'],
];

let missing = 0;
for (const [n, a, b] of ITEMS) {
  const okCn = cn.includes(a);
  const okEn = en.includes(b);
  if (!okCn || !okEn) {
    console.log(`MISSING #${n} cn=${okCn} en=${okEn} :: ${a}`);
    missing += 1;
  }
}
console.log(`[c1-recheck] items=${ITEMS.length} missing=${missing} result=${missing === 0 ? 'C1 CLOSED' : 'STILL OPEN'}`);
