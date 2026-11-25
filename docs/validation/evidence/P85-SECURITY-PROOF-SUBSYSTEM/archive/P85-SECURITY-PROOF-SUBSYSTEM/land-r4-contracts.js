// P85 ②-R4: append-only contract landing (CONTRACTS.md / CONTRACTS_EN.md).
const fs = require('fs');
const child = require('child_process');

const CN = 'docs/CONTRACTS.md';
const EN = 'docs/CONTRACTS_EN.md';

const cnT3Append = 'Windows 上该安装来源路径的根名称链（卷根以下各组件）必须逐组件以「打开为自身」分类，junction、symlink 与其他 reparse 标记一律拒绝；UNC、扩展与相对路径同样拒绝。该判据只约束取证绑定的名称链，与 spawn 名的既有表述互不替代。';
const cnP1 = '冻结策略制品在装配期封存时，其摘要锚必须以 `managedWriteProofShape` 记录 `ManagedWriteProofShape`（闭集仅 `creation` 与 `reuse`）；重入复验按该 Shape 重建同域并逐字节对证，字段缺失或取值未知一律拒绝（fail-closed），既有锚不迁移。';
const cnP2 = '受管身份写保护观测的摘要域为 `proofrail:agent-runner-policy-write-guard:2\\n`（域版本 `:2`）；观察向量固定为四槽 `[策略目录 | 父目录 | 制品 | 锚]`，缺失的槽以规范空标记 `∅` 占位，槽序与槽数不随对象是否存在而变化，跨版本摘要不得比较为相等。';

const enT3Append = ' On Windows, the root name chain of that install-source path (every component below the volume root) must be classified component by component as "opened as itself", and a junction, symlink or any other reparse marker is refused; UNC, extended and relative paths are likewise refused. This criterion binds only the name chain of the evidence binding; it is not interchangeable with any existing statement about the spawn name.';
const enP1 = 'When the frozen policy artifact is sealed at assembly, its digest anchor must record `ManagedWriteProofShape` in `managedWriteProofShape` (closed set: `creation` and `reuse` only); re-entry re-verification rebuilds the same domain from that shape and compares byte for byte, refusing a missing or unknown value (fail-closed), and existing anchors are not migrated.';
const enP2 = 'The digest domain of the managed-identity write-protection observation is `proofrail:agent-runner-policy-write-guard:2\\n` (domain version `:2`); the observed vector is always the four fixed slots `[policy directory | parent directory | artifact | anchor]`, an absent slot is written as the canonical empty marker `∅`, and the slot order and count do not vary with which objects exist, so digests from different domain versions are never compared as equal.';

// Anchors: exact, unique paragraph prefixes taken from disk (CN216 / EN210, CN230 / EN224).
const cnT3Anchor = '对 `agent-runner-cli` 通道的可用性探测，探针可执行文件的安装来源必须可证明不可变：';
const enT3Anchor = 'For availability probing of the `agent-runner-cli` channel, the install source of the probe executable must be provably immutable:';
const cnPAnchor = 'CLI Agent 只在为当前 task/attempt 物化的可丢弃 run-workspace 中工作，不得写源目录、store、策略、门禁实现或接受记录。';
const enPAnchor = 'The CLI Agent works only inside the disposable run-workspace materialized for the current task/attempt.';

function load(path) {
  const raw = fs.readFileSync(path, 'utf8');
  const text = raw.replace(/^\uFEFF+/, '');
  if (text.includes('\r')) throw new Error(path + ': CR found');
  return text;
}

function locate(lines, anchor, label) {
  const hits = [];
  lines.forEach((l, i) => { if (l.startsWith(anchor)) hits.push(i); });
  if (hits.length !== 1) throw new Error(label + ': anchor hits=' + hits.length);
  return hits[0];
}

function land(path, t3Anchor, t3Append, pAnchor, p1, p2) {
  const text = load(path);
  const lines = text.split('\n');
  const iT3 = locate(lines, t3Anchor, path + ' T3');
  const iP = locate(lines, pAnchor, path + ' P');
  const oldT3 = lines[iT3];
  lines[iT3] = oldT3 + t3Append;
  if (!lines[iT3].startsWith(oldT3)) throw new Error(path + ': T3 prefix broken');
  // Insert P1, blank, P2 immediately after the anchor paragraph.
  lines.splice(iP + 1, 0, p1, '', p2);
  const out = '\uFEFF' + lines.join('\n');
  fs.writeFileSync(path, out, { encoding: 'utf8' });
  return { t3Line: iT3 + 1, p1Line: iP + 2, p2Line: iP + 4, oldT3Len: oldT3.length, newT3Len: lines[iT3].length };
}

// New-CJK preflight against the HEAD .md corpus (G4b discipline), plus the
// registered new-word whitelist. Any character outside both aborts the landing.
function cjkSet(s) {
  const o = new Set();
  for (const ch of s) {
    const c = ch.codePointAt(0);
    if ((c >= 0x4e00 && c <= 0x9fff) || (c >= 0x3400 && c <= 0x4dbf) || (c >= 0xf900 && c <= 0xfaff)) o.add(ch);
  }
  return o;
}
const corpusText = child.execSync('git grep -h -I -e . HEAD -- "*.md"', { encoding: 'utf8', maxBuffer: 512 * 1024 * 1024 });
const corpus = cjkSet(corpusText);
if (fs.existsSync('tools/gates/cjk-newwords.txt')) {
  for (const c of cjkSet(fs.readFileSync('tools/gates/cjk-newwords.txt', 'utf8'))) corpus.add(c);
}
const added = [cnT3Append, cnP1, cnP2, enT3Append, enP1, enP2];
const mine = new Set();
for (const s of added) for (const c of cjkSet(s)) mine.add(c);
const unknown = [...mine].filter(c => !corpus.has(c));
console.log('NEW_CJK_COUNT=' + mine.size);
if (unknown.length) {
  console.log('UNKNOWN=' + JSON.stringify(unknown));
  process.exit(2);
}
console.log('UNKNOWN=[]');

const cn = land(CN, cnT3Anchor, cnT3Append, cnPAnchor, cnP1, cnP2);
const en = land(EN, enT3Anchor, enT3Append, enPAnchor, enP1, enP2);
console.log(JSON.stringify({ cn, en }, null, 2));
