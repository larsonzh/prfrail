// Append the user's three rulings (hash authority / encoding caliber / BOM event destination) to BRIEFING.md.
const fs = require('fs');
const crypto = require('crypto');
const BR = 'd:/LZProjects/prfrail/tmp/glm-eval/BRIEFING.md';
const MARK = '\n## \u8ffd\u52a0\uff082026-10-05 \u8865 2\uff09';
const PRIOR = '580f0d9cb7d2f9a4ece41e894250084b8f374084dd84075f4306eca9c0231eea';
const sha = (b) => crypto.createHash('sha256').update(b).digest('hex');

const cur = fs.readFileSync(BR);
const bom = cur[0] === 0xef && cur[1] === 0xbb && cur[2] === 0xbf;
const txt = (bom ? cur.slice(3) : cur).toString('utf8');
const base = Buffer.from(txt.split(MARK)[0], 'utf8'); // strip any prior 补-2 block => idempotent
console.log('on-disk bytes=' + cur.length + ' bom=' + bom + ' base_bytes=' + base.length + ' base_sha=' + sha(base));
console.log('base == user-designated authority 580f0d9c = ' + sha(base).startsWith('580f0d9c'));

const head = [
  '',
  '## 追加（2026-10-05 补 2）·用户裁定（哈希权威 / 编码口径 / BOM 事件去向）',
  '',
  '- **哈希权威（裁定三）**：本文件的验收值以 `FROZEN-MANIFEST.txt` 所列**完整哈希**为**唯一权威**；本文件内的 `append_sha256`（自排除口径，剔除末二行）**只作补充校验、非权威**。复核命令：`node tmp/persist.js`（重算清单并比对全部哈希），无需自排除口径。',
  '- **编码口径（裁定二）**：`tmp/` 下**过程件**保持**无 BOM + LF**（本文件、`plan-v1.md`、`PROBE-NOTES.txt`、`GLM-USAGE-READING.txt` 等）。依据：§6.3 的 G4a 作用域 = 本次**待提交变更集**，而 `tmp/` 被 gitignore、**不在集内** ⇒ G4a 不适用；无 BOM 不影响 Node / 编辑器读取。若未来这些交接件**归档**进 `docs/validation/evidence/**`，届时在**归档执行片**内按 §1.5.1「归档后类型归一」（`.md` → with BOM + LF）统一转换并产生新哈希 —— **属预期流程、非偏差**。',
  '- **BOM 事件去向（裁定一）**：外部保存定性确认（`base_sha = d8463ef5…` 精确一致 ⇒ 内容逐字节未变；已复原为无 BOM + LF；取证链 = `tmp/recon-briefing.js` + `tmp/fix-bom-and-append.js`）。⇒ **⑧a 报告须写入「异常与兜底记录」节**（本片新增必载项）。',
  '- 其余确认（数据集 5 行累计口径、成本节必载项 4 条、冻结包逐字节未变、门禁全绿、授权边界）见上一节，用户已逐项确认。',
  '- **主哈希变更（追加导致，非重写）**：本文件权威哈希 `' + PRIOR + '`（16,973 B）→ `__NEWHASH__`（同链登记：`d8463ef5…` → `580f0d9c…` → 本值；口径：本文件**末二行** —— 本行与 `append_sha256` 行 —— 不计入哈希计算，自指不可内嵌；权威值始终以 `FROZEN-MANIFEST.txt` 为准）',
];
const hashedBlock = head.slice(0, -1).join('\n') + '\n';
const apSha = sha(Buffer.from(hashedBlock, 'utf8')).slice(0, 16);
const newHash = sha(Buffer.concat([base, Buffer.from(hashedBlock, 'utf8')]));
const finalText = (head.join('\n') + '\n' + 'append_sha256=__APSHA__\n')
  .replace('__NEWHASH__', newHash).replace('__APSHA__', apSha);
fs.writeFileSync(BR, Buffer.concat([base, Buffer.from(finalText, 'utf8')]));
const out = fs.readFileSync(BR);
console.log('BRIEFING.md: bytes=' + out.length + ' bom=' + (out[0] === 0xef) + ' crlf=' + out.includes(Buffer.from('\r\n')) + ' sha=' + sha(out));
console.log('append_sha256=' + apSha + '  new_hash_self_excluded=' + newHash);
console.log('verify self-excluded = ' + sha(Buffer.from(out.toString('utf8').split('\n').slice(0, -3).join('\n') + '\n', 'utf8')));
