// Reconstruct BRIEFING.md from the transcript's own tool-call payloads and locate any byte-level drift.
const fs = require('fs');
const crypto = require('crypto');
const TR = 'c:/Users/\u5999\u5999\u545c/AppData/Roaming/Code/User/workspaceStorage/cb55b20f0cd6f845c95f4538166cf997/GitHub.copilot-chat/transcripts/a59cf6ba-ab9b-4477-b453-a98fb5b0e85c.jsonl';
const BR = 'd:/LZProjects/prfrail/tmp/glm-eval/BRIEFING.md';
const sha = (s) => crypto.createHash('sha256').update(typeof s === 'string' ? Buffer.from(s, 'utf8') : s).digest('hex');

const steps = [];
function walk(o, fn, d) {
  if (o == null || d > 14) return;
  if (Array.isArray(o)) { o.forEach((x) => walk(x, fn, d + 1)); return; }
  if (typeof o === 'object') { fn(o); Object.keys(o).forEach((k) => walk(o[k], fn, d + 1)); }
}
fs.readFileSync(TR, 'utf8').split(/\r?\n/).forEach((ln, i) => {
  if (!ln.includes('BRIEFING.md')) return;
  let obj; try { obj = JSON.parse(ln); } catch (e) { return; }
  walk(obj, (o) => {
    if (typeof o.filePath !== 'string' || !o.filePath.endsWith('BRIEFING.md')) return;
    let repl = null;
    if (Array.isArray(o.replacements)) repl = o.replacements.filter((r) => String(r.filePath || '').endsWith('BRIEFING.md'));
    else if (o.oldString) repl = [{ oldString: o.oldString, newString: o.newString }];
    steps.push({ line: i + 1, content: o.content || null, repl });
  }, 0);
});
console.log('payload records=' + steps.length);
steps.forEach((s, i) => console.log('  #' + i + ' L' + s.line + (s.content ? ' CREATE bytes=' + Buffer.byteLength(s.content) : ' EDIT n=' + (s.repl ? s.repl.length : 0)) + (s.content || s.repl ? '' : '  (no payload: read/lookup)')));

let cur = null;
for (const s of steps) {
  if (s.content) { cur = s.content.replace(/\r\n/g, '\n'); continue; }
  if (!cur || !s.repl) continue;
  for (const r of s.repl) {
    if (!cur.includes(r.oldString)) { console.log('MISS L' + s.line + ' old=' + JSON.stringify(r.oldString.slice(0, 40))); continue; }
    cur = cur.replace(r.oldString, r.newString);
  }
}
const now = fs.readFileSync(BR, 'utf8').split('\n## \u8ffd\u52a0')[0] + '\n';
console.log('reconstructed bytes=' + Buffer.byteLength(cur) + ' sha=' + sha(cur));
console.log('current base  bytes=' + Buffer.byteLength(now) + ' sha=' + sha(now));
const a = Buffer.from(cur, 'utf8'), b = Buffer.from(now, 'utf8');
console.log('equal=' + a.equals(b) + ' (reported_end = 14687 / d8463ef5..., measured_now = 14690 / 64775589...)');
if (!a.equals(b)) {
  let i = 0; while (i < Math.min(a.length, b.length) && a[i] === b[i]) i++;
  console.log('first diff at byte ' + i + '  reconstructed=' + JSON.stringify(a.slice(Math.max(0, i - 60), i + 60).toString('utf8')));
  console.log('                       measured     =' + JSON.stringify(b.slice(Math.max(0, i - 60), i + 60).toString('utf8')));
}
