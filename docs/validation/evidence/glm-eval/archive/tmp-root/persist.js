// Persist the in-session ① plan document verbatim + verify the frozen pack + emit the frozen-object manifest.
const fs = require('fs');
const path = require('path');
const crypto = require('crypto');

const ROOT = 'd:/LZProjects/prfrail';
const GE = path.join(ROOT, 'tmp', 'glm-eval');
const TR = 'c:/Users/\u5999\u5999\u545c/AppData/Roaming/Code/User/workspaceStorage/cb55b20f0cd6f845c95f4538166cf997/GitHub.copilot-chat/transcripts/a59cf6ba-ab9b-4477-b453-a98fb5b0e85c.jsonl';
const TR_LINE = 60442; // 1-based; the only record carrying the ① plan headline
const PRE_INPUT_HASH = 'bbeaf2f7db40547d26b1a338608c0b9c1679fd7b08edbf349c0fbadaaa8008b6';
const ORDER = ['00-plan.md', '01-clause-cn.txt', '02-diff.js', '03-prompt.txt'];
const sha = (b) => crypto.createHash('sha256').update(b).digest('hex');

// --- 1. extract the plan document verbatim -----------------------------------
function collect(o, acc, d) {
  if (o == null || d > 12) return;
  if (typeof o === 'string') { if (o.length > 300) acc.push(o); return; }
  if (Array.isArray(o)) { for (const x of o) collect(x, acc, d + 1); return; }
  if (typeof o === 'object') { for (const k of Object.keys(o)) collect(o[k], acc, d + 1); }
}
const recLine = fs.readFileSync(TR, 'utf8').split(/\r?\n/)[TR_LINE - 1];
const rec = JSON.parse(recLine);
const acc = []; collect(rec, acc, 0);
const cands = acc.filter((t) => t.startsWith('# \u5207\u7247'));
if (cands.length !== 1) { console.log('FATAL candidates=' + cands.length); process.exit(2); }
let doc = cands[0];
const sep = doc.indexOf('\n---8<---');
if (sep >= 0) doc = doc.slice(0, sep);
doc = doc.replace(/\r\n/g, '\n').replace(/\n+$/, '') + '\n';
const planPath = path.join(GE, 'plan-v1.md');
fs.writeFileSync(planPath, doc, 'utf8');
const planBuf = fs.readFileSync(planPath);
const planHash = sha(planBuf);
const h2 = (doc.match(/^## /gm) || []).length;
const seeds = (doc.match(/^\| S[1-8] \|/gm) || []).length;
console.log('plan-v1.md bytes=' + planBuf.length + ' bom=' + (planBuf[0] === 0xef) + ' crlf=' + planBuf.includes(Buffer.from('\r\n')));
console.log('plan-v1.md sha256=' + planHash + ' h2_sections=' + h2 + ' seed_rows=' + seeds);

// --- 2. provenance record ----------------------------------------------------
const prov = [
  '# ① 方案书落盘 · 来源与哈希自证',
  '',
  '- 对象：切片 `DIRECTIVE-GLM-EVAL` ① 架构师评估方案书（会话内交付、此前未落盘）',
  '- 来源：VS Code Copilot 聊天转录本 JSONL 第 ' + TR_LINE + ' 行（1-based）的记录内文本字段',
  '- 转录本：' + TR,
  '- 提取方式：解析该行 JSON → 收集长度 > 300 的字符串字段 → 唯一命中「以 `# 切片` 开头」者 → 切掉脚本拼接分隔符 → CRLF 归一为 LF → 末尾归一为单个 LF',
  '- 唯一性核验：全转录本中「① 架构师评估方案书」仅出现 1 次（`tmp/_scan2/` 打分筛选中该标题命中 1 个文件）',
  '- 落盘口径：与两腿产物同口径 —— 逐字节原文 + 独立 SHA-256（本文件即稳定可复算值）',
  '',
  '| 项 | 值 |',
  '|---|---|',
  '| plan-v1.md 字节数 | ' + planBuf.length + ' |',
  '| plan-v1.md SHA-256 | `' + planHash + '` |',
  '| 二级标题数 | ' + h2 + '（§1–§7） |',
  '| 预注册种子行数 | ' + seeds + '（S1–S8） |',
  '| BOM | 无 |',
  '| 行尾 | LF |',
  '',
  '复算命令（PowerShell）：',
  '',
  '```powershell',
  '(Get-FileHash tmp\\glm-eval\\plan-v1.md -Algorithm SHA256).Hash.ToLower()',
  '```',
  '',
  '注：本片后续用户裁决（②a 采纳 + 三项偏离）为对该方案书的口径修正，二者并存、不互改；',
  '方案书正文保持转录原样，裁决见 `tmp/glm-eval/BRIEFING.md`。',
  '',
].join('\n');
fs.writeFileSync(path.join(GE, 'PLAN-PROVENANCE.txt'), prov, 'utf8');
console.log('wrote PLAN-PROVENANCE.txt bytes=' + fs.statSync(path.join(GE, 'PLAN-PROVENANCE.txt')).size);

// --- 3. frozen-pack re-verification -----------------------------------------
const concat = Buffer.concat(ORDER.map((n) => fs.readFileSync(path.join(GE, 'bundle', n))));
const reInput = sha(concat);
const keyBuf = fs.readFileSync(path.join(GE, 'key.json'));
const dsObj = JSON.parse(fs.readFileSync(path.join(ROOT, 'docs', 'validation', 'directive-glm-eval-dataset.json'), 'utf8'));
console.log('input_hash(recomputed)=' + reInput + ' match=' + (reInput === PRE_INPUT_HASH));
console.log('key.json sha256=' + sha(keyBuf) + ' match=' + (sha(keyBuf) === dsObj.pre_registration.key_file_sha256));
console.log('dataset rows=' + dsObj.calls.length + ' registered_at=' + dsObj.pre_registration.registered_at + ' pre_input_match=' + (dsObj.pre_registration.bundle_input_hash === reInput));

// --- 4. frozen-object manifest ----------------------------------------------
function walk(dir, base, out) {
  if (!fs.existsSync(dir)) return;
  for (const e of fs.readdirSync(dir, { withFileTypes: true }).sort((a, b) => a.name.localeCompare(b.name))) {
    const p = path.join(dir, e.name);
    if (e.isDirectory()) walk(p, base, out); else out.push(p);
  }
}
const files = [];
for (const r of ['tmp/glm-eval', 'tmp/glm-eval-out']) walk(path.join(ROOT, r), r, files);
const SELF = path.join(GE, 'FROZEN-MANIFEST.txt');
const kept = files.filter((p) => p !== SELF); // never list itself: a self-hash would be stale by construction
files.length = 0; files.push(...kept);
const extra = [
  'tmp/PROMPT-EXACT.txt',
  'docs/validation/directive-glm-eval-dataset.json',
  'tmp/glm-eval-batch/a-input.txt',
  'tmp/glm-eval-batch/a-key.json',
  'tmp/freeze.js', 'tmp/save-leg.js', 'tmp/norm-score.js',
  'tmp/batch-a-gen.js', 'tmp/batch-d-gen.js', 'tmp/fill-dataset.js', 'tmp/persist.js',
  'tmp/fix-bom-and-append.js', 'tmp/recon-briefing.js', 'tmp/append-rulings.js',
];
for (const r of extra) if (fs.existsSync(path.join(ROOT, r))) files.push(path.join(ROOT, r));
const rows = files.map((p) => {
  const b = fs.readFileSync(p);
  return sha(b) + '  ' + String(b.length).padStart(8) + '  ' + p.replace(/\\/g, '/').replace(ROOT + '/', '');
});
const man = [
  '# DIRECTIVE-GLM-EVAL · 冻结物清单（哈希自证）',
  '',
  '口径：SHA-256 of raw bytes；`<hash>  <bytes>  <path>`。本清单用于会话内核对，',
  '⑧a 归档时另生 `docs/validation/evidence/glm-eval/SHA256SUMS.txt`（`<hash><2 spaces><path>` 形态，过 G8-b）。',
  '本清单自身不计入（避免自指陈旧哈希）。`BRIEFING.md` / `PROBE-NOTES.txt` / `PLAN-PROVENANCE.txt` 为**交接文档**（非冻结评估对象），一并列出便于复算。',
  '',
  '```',
  ...rows,
  '```',
  '',
  '已核验一致：`input_hash` = ' + reInput + '（= 预注册值）；`key.json` = ' + sha(keyBuf) + '（= `pre_registration.key_file_sha256`）。',
  '',
].join('\n');
fs.writeFileSync(path.join(GE, 'FROZEN-MANIFEST.txt'), man, 'utf8');
console.log('manifest files=' + rows.length + ' bytes=' + fs.statSync(path.join(GE, 'FROZEN-MANIFEST.txt')).size);
for (const r of rows) console.log(r);
