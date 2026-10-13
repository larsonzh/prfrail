// Slice DIRECTIVE-PREREVIEW-GUIDANCE : assertion-based landing script.
// CN authoritative + _EN strict mirror : 6 + 6 frozen edits, nothing else.
// Every anchor must match exactly once; files must be BOM + LF; verify after write.
'use strict';

const fs = require('fs');
const path = require('path');

const ROOT = path.resolve(__dirname, '..', '..');
const CN = path.join(ROOT, 'docs', 'DELIVERY_DIRECTIVE.md');
const EN = path.join(ROOT, 'docs', 'DELIVERY_DIRECTIVE_EN.md');

// ---------- helpers ----------
function count(hay, needle) { return hay.split(needle).length - 1; }
function shape(s) {
  const lines = s.split('\n');
  return {
    lines: lines.length,
    bold: (s.match(/\*\*/g) || []).length,
    bars: (s.match(/\|/g) || []).length,
    blank: lines.filter((x) => x === '').length,
  };
}
function assertOnce(name, hay, needle) {
  const n = count(hay, needle);
  if (n !== 1) throw new Error(`ANCHOR FAIL [${name}] occurrences=${n} (expected 1)`);
}
function assertClean(name, s) {
  if (!s.startsWith('\uFEFF')) throw new Error(`BOM FAIL [${name}] file does not start with U+FEFF`);
  if (s.indexOf('\r') !== -1) throw new Error(`EOL FAIL [${name}] file contains CR`);
}

// ---------- frozen edits (verbatim from the slice definition) ----------
const CN_EDITS = [
  {
    id: 'CN-1 §4.2 pre-reviewer prohibited cell',
    old: '| 预审员 | 方案书、实现 diff、契约原句 | 《预审报告》：总体结论 + 偏离项 + 反例分析（并发/崩溃/边界） | 改文件；看测试代码（保持独立视角） |',
    new: '| 预审员 | 方案书、实现 diff、契约原句 | 《预审报告》：总体结论 + 偏离项 + 反例分析（并发/崩溃/边界） | **改文件**；其余见 §4.2.1 / §4.2.2 |',
  },
  {
    id: 'CN-2 §4.2.1 + §4.2.2 block insert',
    old:
      '| 独立终审员（⑦） | 最终 diff + 契约原句（**首轮不附任何清单**） | 安全 / 架构一致性 / 完整性 / 测试反例四段 + 结论行 | 改文件；首轮参考任何中间结论 |' +
      '\n\n### 4.3 严重度定义',
    new:
      '| 独立终审员（⑦） | 最终 diff + 契约原句（**首轮不附任何清单**） | 安全 / 架构一致性 / 完整性 / 测试反例四段 + 结论行 | 改文件；首轮参考任何中间结论 |' +
      '\n\n**4.2.1 预审员的最低核验动作（五条）**：\n' +
      '\n' +
      '1. 对 diff 中提及但未展开的工件（装置脚本、测试文件、证据工件、报告附件等），必须实际读取关键文件本身，不得仅凭报告叙述或推测下结论。\n' +
      '2. 对 diff 或报告涉及的关键数值、时间戳、状态序列断言，必须独立复算至少 2 处：涉及原始工件时重派生；涉及跨文件引用时定位引用点核对。\n' +
      '3. 若发现报告或设计文档的叙述与证据工件之间存在任何矛盾（哪怕很小），必须作为独立发现列出并判定严重度——不得以「可能只是笔误」放过。\n' +
      '4. 反例分析须给出至少一条「若删除某个守卫 / 改回某个逻辑，哪条断言或哪份证据会失败」的可证伪判据；若无，须明确说明原因。\n' +
      '5. 结论必须基于审查方自己实际执行的核验；不得引用报告自身对自身的自述作为证据（自证不可信）。\n' +
      '\n' +
      '以上为最低核验动作，不得以此替代自主判断；超出最低要求继续深挖亦不在此限。\n' +
      '\n' +
      '**4.2.2 预审员的禁止项（权威措辞）**：改文件；以测试代码或文档变更本身作为判断实现正确性的依据（读取上述工件以核验报告声明不在此限；保持独立视角）。\n' +
      '\n### 4.3 严重度定义',
  },
  {
    id: 'CN-3 §4.3 Critical meaning',
    old: '| Critical | 数据丢失、安全边界失效、静默错误结论 | 必须修复并重跑该审查环节 |',
    new: '| Critical | 数据丢失、安全边界失效、静默错误结论；且报告核心因果被自身证据推翻者属本级别 | 必须修复并重跑该审查环节 |',
  },
  {
    id: 'CN-4 §4.3 High meaning',
    old: '| High | 功能失效、崩溃、契约冲突 | 同上 |',
    new: '| High | 功能失效、崩溃、契约冲突；协议属性与设计不符、但 fail-closed 已避免错误结论者属本级别 | 同上 |',
  },
  {
    id: 'CN-5 §4.3 severity-boundary note insert',
    old: '| Low | 风格、命名、注释、可读性 | 记录在案，可延后 |\n\n---',
    new:
      '| Low | 风格、命名、注释、可读性 | 记录在案，可延后 |\n' +
      '\n' +
      '> **定级边界注记**：同一实质发现不因定级差异而削弱其效力——定级决定的是处置动作（重跑该环 vs 定向复核），不决定发现是否成立。\n' +
      '\n---',
  },
  {
    id: 'CN-6 appendix B.5 forbidden line',
    old: '禁止：改文件；查看测试代码或文档变更（保持独立视角）。',
    new: '禁止：改文件；以测试代码或文档变更本身作为判断实现正确性的依据（读取上述工件以核验报告声明不在此限；保持独立视角）。最低核验动作：见 §4.2.1（不在此重复正文）。',
  },
];

const EN_EDITS = [
  {
    id: 'EN-1 §4.2 pre-reviewer prohibited cell',
    old: '| Pre-reviewer | proposal, implementation diff, verbatim contract text | Pre-review Report: overall conclusion + deviation items + counterexample analysis (concurrency/crash/boundary) | modify files; look at test code (keep an independent perspective) |',
    new: '| Pre-reviewer | proposal, implementation diff, verbatim contract text | Pre-review Report: overall conclusion + deviation items + counterexample analysis (concurrency/crash/boundary) | **modify files**; for the rest see §4.2.1 / §4.2.2 |',
  },
  {
    id: 'EN-2 §4.2.1 + §4.2.2 block insert',
    old:
      '| Independent final reviewer (⑦) | final diff + verbatim contract text (**no list attached in the first round**) | four sections: security / architecture consistency / completeness / test counterexamples + conclusion line | modify files; in the first round reference any intermediate conclusion |' +
      '\n\n### 4.3 Severity Definitions',
    new:
      '| Independent final reviewer (⑦) | final diff + verbatim contract text (**no list attached in the first round**) | four sections: security / architecture consistency / completeness / test counterexamples + conclusion line | modify files; in the first round reference any intermediate conclusion |' +
      '\n\n**4.2.1 Minimum verification actions of the pre-reviewer (five items)**:\n' +
      '\n' +
      '1. For artifacts mentioned but not expanded in the diff (harness scripts, test files, evidence artifacts, report attachments, etc.), the key files themselves must actually be read; a conclusion must not rest on the report narrative or speculation alone.\n' +
      '2. For key numbers, timestamps and state-sequence assertions involved in the diff or the report, at least 2 must be independently recomputed: re-derive when raw artifacts are involved; locate the reference point and check it when cross-file references are involved.\n' +
      '3. Any contradiction (however small) found between the narrative of the report or design document and the evidence artifacts must be listed as an independent finding with a severity - it must not be let go as "probably just a typo".\n' +
      '4. The counterexample analysis must give at least one falsifiable criterion of the form "if guard X is deleted / logic Y is reverted, which assertion or piece of evidence would fail"; if there is none, the reason must be stated explicitly.\n' +
      '5. The conclusion must rest on verification the reviewer actually performed; the report\'s own self-description must not be cited as evidence (self-attestation is not trustworthy).\n' +
      '\n' +
      'The above are minimum verification actions; they must not replace independent judgment; going beyond the minimum to dig deeper is likewise not restricted.\n' +
      '\n' +
      '**4.2.2 Prohibited items of the pre-reviewer (authoritative wording)**: editing files; using test code or documentation changes themselves as the basis for judging implementation correctness (reading the above artifacts to verify report claims is not restricted by this; keep an independent perspective).\n' +
      '\n### 4.3 Severity Definitions',
  },
  {
    id: 'EN-3 §4.3 Critical meaning',
    old: '| Critical | data loss, security boundary failure, silently wrong conclusion | must fix and re-run that review stage |',
    new: '| Critical | data loss, security boundary failure, silently wrong conclusion; and a report whose core causation is overturned by its own evidence belongs to this level | must fix and re-run that review stage |',
  },
  {
    id: 'EN-4 §4.3 High meaning',
    old: '| High | functional failure, crash, contract conflict | same as above |',
    new: '| High | functional failure, crash, contract conflict; a protocol property that does not match the design but fail-closed already prevented a wrong conclusion belongs to this level | same as above |',
  },
  {
    id: 'EN-5 §4.3 severity-boundary note insert',
    old: '| Low | style, naming, comments, readability | recorded on the ledger, may be deferred |\n\n---',
    new:
      '| Low | style, naming, comments, readability | recorded on the ledger, may be deferred |\n' +
      '\n' +
      '> **Severity-boundary note**: the same substantive finding is not weakened by a difference in severity rating - the rating determines the handling action (re-run that stage vs a targeted re-check), not whether the finding stands.\n' +
      '\n---',
  },
  {
    id: 'EN-6 appendix B.5 forbidden line',
    old: 'Forbidden: editing files; inspecting test code or documentation changes (preserve an independent perspective).',
    new: 'Forbidden: editing files; using test code or documentation changes themselves as the basis for judging implementation correctness (reading the above artifacts to verify report claims is not restricted by this; keep an independent perspective). Minimum verification actions: see §4.2.1 (not repeated here in full).',
  },
];

// ---------- load + pre-write assertions ----------
let cn = fs.readFileSync(CN, 'utf8');
let en = fs.readFileSync(EN, 'utf8');
assertClean('CN', cn);
assertClean('EN', en);

const cnBefore = shape(cn);
const enBefore = shape(en);

// assert every anchor is unique on the untouched source, then apply
for (const e of CN_EDITS) assertOnce(e.id, cn, e.old);
for (const e of EN_EDITS) assertOnce(e.id, en, e.old);

for (const e of CN_EDITS) cn = cn.split(e.old).join(e.new);
for (const e of EN_EDITS) en = en.split(e.old).join(e.new);

// post-transform well-formedness
assertClean('CN-after', cn);
assertClean('EN-after', en);

// ---------- write ----------
fs.writeFileSync(CN, cn, 'utf8');
fs.writeFileSync(EN, en, 'utf8');

// ---------- verify from disk ----------
const cn2 = fs.readFileSync(CN, 'utf8');
const en2 = fs.readFileSync(EN, 'utf8');
assertClean('CN-disk', cn2);
assertClean('EN-disk', en2);

const report = [];
for (const e of CN_EDITS) {
  const oldLeft = count(cn2, e.old);
  const newSeen = count(cn2, e.new);
  const ok = oldLeft === 0 && newSeen === 1;
  report.push({ id: e.id, oldLeft, newSeen, status: ok ? 'LANDED' : 'NOT-LANDED' });
  if (!ok) throw new Error(`VERIFY FAIL [${e.id}] oldLeft=${oldLeft} newSeen=${newSeen}`);
}
for (const e of EN_EDITS) {
  const oldLeft = count(en2, e.old);
  const newSeen = count(en2, e.new);
  const ok = oldLeft === 0 && newSeen === 1;
  report.push({ id: e.id, oldLeft, newSeen, status: ok ? 'LANDED' : 'NOT-LANDED' });
  if (!ok) throw new Error(`VERIFY FAIL [${e.id}] oldLeft=${oldLeft} newSeen=${newSeen}`);
}

const cnAfter = shape(cn2);
const enAfter = shape(en2);
const d = (a, b, k) => b[k] - a[k];
const cnDelta = { lines: d(cnBefore, cnAfter, 'lines'), bold: d(cnBefore, cnAfter, 'bold'), bars: d(cnBefore, cnAfter, 'bars'), blank: d(cnBefore, cnAfter, 'blank') };
const enDelta = { lines: d(enBefore, enAfter, 'lines'), bold: d(enBefore, enAfter, 'bold'), bars: d(enBefore, enAfter, 'bars'), blank: d(enBefore, enAfter, 'blank') };

// parity assertion : CN delta === EN delta (item by item) and === frozen expectation
const expect = { lines: 14, bold: 8, bars: 0, blank: 5 };
for (const k of ['lines', 'bold', 'bars', 'blank']) {
  if (cnDelta[k] !== enDelta[k]) throw new Error(`PARITY FAIL ${k}: CN=${cnDelta[k]} EN=${enDelta[k]}`);
  if (cnDelta[k] !== expect[k]) throw new Error(`EXPECT FAIL ${k}: got=${cnDelta[k]} want=${expect[k]}`);
}
// full-file parity (leftover shape difference must not pre-exist unnoticed)
for (const k of ['lines', 'bold', 'bars', 'blank']) {
  if (cnAfter[k] !== enAfter[k]) throw new Error(`FULL-FILE PARITY FAIL ${k}: CN=${cnAfter[k]} EN=${enAfter[k]}`);
}

console.log('=== per-edit coverage ===');
for (const r of report) console.log(`${r.status}\t${r.id}\toldLeft=${r.oldLeft}\tnewSeen=${r.newSeen}`);
console.log('\n=== CN shape ===');
console.log('before', JSON.stringify(cnBefore));
console.log('after ', JSON.stringify(cnAfter));
console.log('delta ', JSON.stringify(cnDelta));
console.log('\n=== EN shape ===');
console.log('before', JSON.stringify(enBefore));
console.log('after ', JSON.stringify(enAfter));
console.log('delta ', JSON.stringify(enDelta));
console.log('\n=== parity ===');
console.log('CN delta === EN delta === {lines:14,bold:8,bars:0,blank:5}: true');
console.log('full-file CN === EN (lines/bold/bars/blank): true');
console.log('\nLAND: DONE');
