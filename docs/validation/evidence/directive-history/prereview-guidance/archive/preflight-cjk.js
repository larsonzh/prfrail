// G4b preflight (step 0): find CJK chars in the planned new text that are absent
// from the HEAD *.md corpus and the cjk-newwords whitelist. No writes.
const { execFileSync } = require('child_process');
const fs = require('fs');
const sh = (a) => execFileSync('git', a, { encoding: 'utf8', maxBuffer: 1 << 28 });

const files = sh(['ls-files', '*.md']).split('\n').filter(Boolean);
let corpus = '';
for (const f of files) {
  try { corpus += sh(['show', 'HEAD:' + f]); } catch (e) { /* ignore */ }
}
const isCjk = (ch) => {
  const u = ch.codePointAt(0);
  return (u >= 0x3400 && u <= 0x4dbf) || (u >= 0x4e00 && u <= 0x9fff) || (u >= 0xf900 && u <= 0xfaff);
};
const setOf = (s) => new Set(Array.from(s).filter(isCjk));
const corpusSet = setOf(corpus);
const wlRaw = fs.existsSync('tools/gates/cjk-newwords.txt')
  ? fs.readFileSync('tools/gates/cjk-newwords.txt', 'utf8') : '';
const wlSet = setOf(wlRaw);

const draft = [
  // --- step 0 ledger rows ---
  '**试点继续顺延第 7 次（用户直接裁定）**：本片（`DIRECTIVE-PREREVIEW-GUIDANCE`）为纯文档准则治理片 ⇒ 不满足 §12.1 的试点选择标准；按 §7.8.6 与 §12.1 的配套②升级用户裁决，顺延次数以 §12.4 为唯一真值源（本次记为第 7 次）。四要素：裁决人 = 用户；日期 = 2026-10-13（本片步骤 0）；理由 = 纯文档准则治理片不满足 §12.1 的试点选择标准；依据 = 用户同轮裁定原文。**按 OB-37 元规则独立留痕**，如实标注「超出 OB-11 字面标准、系用户直接裁定」，且**不得作为先例**。',
  '**已裁决**（用户直接裁定：继续执行、不作为试点）',
  '**n=2 追加（2026-10-13）**：V4 Pro 引导版（15 min / 成本约 ¥3.07）与 Flash 引导版（19 min / 成本约 ¥1.09）两臂（同一 vendor 家族的两个变体）均抓到核心缺陷 ⇒ 引导增益跨变体同向；引导文本与读数见用户实验报告（未入库）。',
  '**已处置**（切片 `DIRECTIVE-PREREVIEW-GUIDANCE` 起跑，2026-10-13）',
  // --- D-1 five items + closing sentence ---
  '对 diff 中提及但未展开的工件（装置脚本、测试文件、证据工件、报告附件等），必须实际读取关键文件本身，不得仅凭报告叙述或推测下结论。',
  '对 diff 或报告涉及的关键数值、时间戳、状态序列断言，必须独立复算至少 2 处：涉及原始工件时重派生；涉及跨文件引用时定位引用点核对。',
  '若发现报告或设计文档的叙述与证据工件之间存在任何矛盾（哪怕很小），必须作为独立发现列出并判定严重度——不得以「可能只是笔误」放过。',
  '反例分析须给出至少一条「若删除某个守卫 / 改回某个逻辑，哪条断言或哪份证据会失败」的可证伪判据；若无，须明确说明原因。',
  '结论必须基于审查方自己实际执行的核验；不得引用报告自身对自身的自述作为证据（自证不可信）。',
  '以上为最低核验动作，不得以此替代自主判断；超出最低要求继续深挖亦不在此限。',
  // --- D-2 severity boundary ---
  '报告核心因果被自身证据推翻；协议属性与设计不符但 fail-closed 已避免错误结论；二者不因定级差异削弱实质发现。',
  // --- D-3 / D-4 prohibition + pointer ---
  '改文件；以测试代码或文档变更本身作为判断实现正确性的依据（读取上述工件以核验报告声明不在此限；保持独立视角）',
  '最低核验动作：见 §4.2.1（不在此重复正文）。',
].join('\n');

const draftSet = setOf(draft);
const unknown = Array.from(draftSet).filter((c) => !corpusSet.has(c) && !wlSet.has(c));
console.log('mdFiles=' + files.length + ' corpusCjk=' + corpusSet.size + ' whitelist=' + wlSet.size + ' draftCjk=' + draftSet.size);
console.log('unknown=' + JSON.stringify(unknown));
console.log('RESULT=' + (unknown.length === 0 ? 'PASS' : 'FINDINGS'));
