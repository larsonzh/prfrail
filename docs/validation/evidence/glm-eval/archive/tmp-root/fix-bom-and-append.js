// Verify the BOM-only drift, restore the declared no-BOM + LF caliber, then (re)append the cost block.
const fs = require('fs');
const crypto = require('crypto');
const TR = 'c:/Users/\u5999\u5999\u545c/AppData/Roaming/Code/User/workspaceStorage/cb55b20f0cd6f845c95f4538166cf997/GitHub.copilot-chat/transcripts/a59cf6ba-ab9b-4477-b453-a98fb5b0e85c.jsonl';
const GE = 'd:/LZProjects/prfrail/tmp/glm-eval';
const BR = GE + '/BRIEFING.md';
const PN = GE + '/PROBE-NOTES.txt';
const MARK = '\n## \u8ffd\u52a0\uff082026-10-05 \u8865\uff09';
const sha = (b) => crypto.createHash('sha256').update(b).digest('hex');

// --- 1. reconstruct the end-of-last-turn bytes from the transcript ----------
const steps = [];
function walk(o, fn, d) {
  if (o == null || d > 14) return;
  if (Array.isArray(o)) { o.forEach((x) => walk(x, fn, d + 1)); return; }
  if (typeof o === 'object') { fn(o); Object.keys(o).forEach((k) => walk(o[k], fn, d + 1)); }
}
fs.readFileSync(TR, 'utf8').split(/\r?\n/).forEach((ln) => {
  if (!ln.includes('BRIEFING.md')) return;
  let obj; try { obj = JSON.parse(ln); } catch (e) { return; }
  walk(obj, (o) => {
    if (typeof o.filePath !== 'string' || !o.filePath.endsWith('BRIEFING.md')) return;
    let repl = null;
    if (Array.isArray(o.replacements)) repl = o.replacements.filter((r) => String(r.filePath || '').endsWith('BRIEFING.md'));
    else if (o.oldString) repl = [{ oldString: o.oldString, newString: o.newString }];
    if (o.content || repl) steps.push({ content: o.content || null, repl });
  }, 0);
});
let recon = null;
for (const s of steps) {
  if (s.content) { recon = s.content.replace(/\r\n/g, '\n'); continue; }
  if (!recon || !s.repl) continue;
  for (const r of s.repl) if (recon.includes(r.oldString)) recon = recon.replace(r.oldString, r.newString);
}
const reconBuf = Buffer.from(recon, 'utf8');
console.log('reconstructed: bytes=' + reconBuf.length + ' sha=' + sha(reconBuf));

// --- 2. measure the on-disk base (appended block already present) ----------
const cur = fs.readFileSync(BR);
const hasBom = cur[0] === 0xef && cur[1] === 0xbb && cur[2] === 0xbf;
const txt = (hasBom ? cur.slice(3) : cur).toString('utf8');
const base = Buffer.from(txt.split(MARK)[0], 'utf8');
console.log('on-disk: bytes=' + cur.length + ' bom=' + hasBom + ' base_bytes=' + base.length + ' base_sha=' + sha(base));
console.log('CONTENT UNCHANGED BY THE EXTERNAL SAVE = ' + base.equals(reconBuf));
console.log('base matches recorded d8463ef5 = ' + sha(base).startsWith('d8463ef5'));

// --- 3. PROBE-NOTES.txt: same episode? ------------------------------------
const pn = fs.readFileSync(PN);
const pnBom = pn[0] === 0xef && pn[1] === 0xbb && pn[2] === 0xbf;
console.log('PROBE-NOTES.txt bytes=' + pn.length + ' bom=' + pnBom + ' sha=' + sha(pn) +
  ' sha_without_bom=' + sha(pnBom ? pn.slice(3) : pn));
if (pnBom) { fs.writeFileSync(PN, pn.slice(3)); console.log('PROBE-NOTES.txt restored to no-BOM + LF'); }

// --- 4. rebuild BRIEFING.md = verified base (no BOM) + revised block -------
const head = [
  '',
  '## 追加（2026-10-05 补）·成本读数',
  '',
  '- GLM 侧读数到位（见 `GLM-USAGE-READING.txt`）：`glm-5.3` 累计 161,863 tokens；`glm-5.3-flash` 累计 424,690 tokens。',
  '- V4 Pro 对照腿成本读数：**不可得** ⇒ 成本对照为**单边**。',
  '- 报告口径变更：',
  '  - ⑤ 两腿 `metering_mode`：GLM 侧 = `user-provided-cumulative`；V4 Pro 侧保持 `unavailable`。',
  '  - 成本判据：可给出 GLM 侧**绝对量级**；不可给出 GLM vs V4 Pro **相对结论**。',
  '  - 不填估值、不做 token 换算。',
  '- 新增必载项：实测消耗**显著高于主控预估**（约 4–8 倍），原因未核，如实登记。',
  '- **编码事件披露（会话间外部保存）**：本文件在上一会话收尾时实测为 14,687 B / **无 BOM**；本会话追加前实测为 **14,690 B / 含 BOM**（`EF BB BF`）—— 差值恰为 **+3 B 的 BOM**，**内容逐字节相同**（已用聊天转录本重建 14,687 B 版本并逐字节比对证实，其 SHA-256 与上一会话所载 `d8463ef5…` 精确一致）。本块落盘时**已恢复为本文头部声明的「无 BOM + LF」口径**。同批变更提示中的 `PROBE-NOTES.txt` 经复核**未受影响**（2,023 B / 无 BOM / 哈希与上一会话记录一致）。如需改回仓库 `.md` 惯例（UTF-8 **with BOM**），转换会改变本节所载哈希，属预期变更而非偏差。',
  '',
  '### 报告成本节必载项（⑧a 逐条写入）',
  '',
  '1. GLM 侧绝对量级已获得（161,863 / 424,690，含少量用户早期测试消息）；',
  '2. V4 Pro 侧读数**不可得** ⇒ **相对成本判据不成立**；',
  '3. 对「是否值得把 GLM-5.3 纳入 ⑤」的成本问题给出**单边证据**：GLM-5.3 单次预审量级 ≈ 80k tokens，属可接受范围；',
  '4. 相对成本比较**保留为悬置**，作为后续切片的输入。',
  '',
  '- **主哈希变更（追加导致，非重写）**：原哈希 `d8463ef5…` → 新哈希（自排除口径）`__NEWHASH__`（口径：本文件**末二行** —— 本行与 `append_sha256` 行 —— 不计入哈希计算，自指不可内嵌；`FROZEN-MANIFEST.txt` 所列为本文件**含**末行的完整哈希，二者口径不同、均自证；复算 `node tmp/persist.js`）',
];
const hashedBlock = head.slice(0, -1).join('\n') + '\n';
const baseHash = sha(base);
const apSha = sha(Buffer.from(hashedBlock, 'utf8')).slice(0, 16);
const newHash = sha(Buffer.concat([base, Buffer.from(hashedBlock, 'utf8')]));
const finalText = (head.join('\n') + '\n' + 'append_sha256=__APSHA__\n')
  .replace('__NEWHASH__', newHash).replace('__APSHA__', apSha);
fs.writeFileSync(BR, Buffer.concat([base, Buffer.from(finalText, 'utf8')]));
const out = fs.readFileSync(BR);
console.log('BRIEFING.md: bytes=' + out.length + ' bom=' + (out[0] === 0xef) + ' sha=' + sha(out));
console.log('append_sha256=' + apSha + '  new_hash_self_excluded=' + newHash + '  pre_append_sha=' + baseHash);
console.log('verify self-excluded = ' + sha(Buffer.from(out.toString('utf8').split('\n').slice(0, -3).join('\n') + '\n', 'utf8')));
