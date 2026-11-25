// P85 ②-R4 strict byte equality of landed content + pre-change shape deltas + G6-4 scan.
const fs = require('fs');
const child = require('child_process');

function loadText(p) {
  return fs.readFileSync(p, 'utf8').replace(/^\uFEFF+/, '');
}
function shape(text) {
  const L = text.split('\n');
  const count = re => L.reduce((n, s) => n + (s.match(re) || []).length, 0);
  return { lines: L.length, blank: L.filter(s => s === '').length, bold: count(/\*\*/g), bars: count(/\|/g), bt: count(/`/g) };
}
const cnT3Append = 'Windows 上该安装来源路径的根名称链（卷根以下各组件）必须逐组件以「打开为自身」分类，junction、symlink 与其他 reparse 标记一律拒绝；UNC、扩展与相对路径同样拒绝。该判据只约束取证绑定的名称链，与 spawn 名的既有表述互不替代。';
const cnP1 = '冻结策略制品在装配期封存时，其摘要锚必须以 `managedWriteProofShape` 记录 `ManagedWriteProofShape`（闭集仅 `creation` 与 `reuse`）；重入复验按该 Shape 重建同域并逐字节对证，字段缺失或取值未知一律拒绝（fail-closed），既有锚不迁移。';
const cnP2 = '受管身份写保护观测的摘要域为 `proofrail:agent-runner-policy-write-guard:2\\n`（域版本 `:2`）；观察向量固定为四槽 `[策略目录 | 父目录 | 制品 | 锚]`，缺失的槽以规范空标记 `∅` 占位，槽序与槽数不随对象是否存在而变化，跨版本摘要不得比较为相等。';
const enT3Append = ' On Windows, the root name chain of that install-source path (every component below the volume root) must be classified component by component as "opened as itself", and a junction, symlink or any other reparse marker is refused; UNC, extended and relative paths are likewise refused. This criterion binds only the name chain of the evidence binding; it is not interchangeable with any existing statement about the spawn name.';
const enP1 = 'When the frozen policy artifact is sealed at assembly, its digest anchor must record `ManagedWriteProofShape` in `managedWriteProofShape` (closed set: `creation` and `reuse` only); re-entry re-verification rebuilds the same domain from that shape and compares byte for byte, refusing a missing or unknown value (fail-closed), and existing anchors are not migrated.';
const enP2 = 'The digest domain of the managed-identity write-protection observation is `proofrail:agent-runner-policy-write-guard:2\\n` (domain version `:2`); the observed vector is always the four fixed slots `[policy directory | parent directory | artifact | anchor]`, an absent slot is written as the canonical empty marker `∅`, and the slot order and count do not vary with which objects exist, so digests from different domain versions are never compared as equal.';

const idxCN = child.execSync('git show :docs/CONTRACTS.md', { encoding: 'utf8' }).replace(/^\uFEFF+/, '');
const idxEN = child.execSync('git show :docs/CONTRACTS_EN.md', { encoding: 'utf8' }).replace(/^\uFEFF+/, '');
const wtCN = loadText('docs/CONTRACTS.md');
const wtEN = loadText('docs/CONTRACTS_EN.md');
const iCN = idxCN.split('\n');
const iEN = idxEN.split('\n');
const wCN = wtCN.split('\n');
const wEN = wtEN.split('\n');

const checks = {
  cnT3Exact: wCN[215] === iCN[215] + cnT3Append,
  enT3Exact: wEN[209] === iEN[209] + enT3Append,
  cnP1Exact: wCN[230] === cnP1,
  cnBlank: wCN[231] === '',
  cnP2Exact: wCN[232] === cnP2,
  enP1Exact: wEN[224] === enP1,
  enBlank: wEN[225] === '',
  enP2Exact: wEN[226] === enP2,
  cnAllElseEqual: (() => {
    const a = idxCN.split('\n'); a[215] = a[215] + cnT3Append; a.splice(230, 0, cnP1, '', cnP2);
    return a.join('\n') === wtCN;
  })(),
  enAllElseEqual: (() => {
    const a = idxEN.split('\n'); a[209] = a[209] + enT3Append; a.splice(224, 0, enP1, '', enP2);
    return a.join('\n') === wtEN;
  })(),
};
const pre = { cn: shape(idxCN), en: shape(idxEN) };
const post = { cn: shape(wtCN), en: shape(wtEN) };
const delta = (k, f) => post[f][k] - pre[f][k];
const deltaOfDelta = k => (post.en[k] - post.cn[k]) - (pre.en[k] - pre.cn[k]);

// G6-4: no two-or-more-digit number immediately followed by an optional space and the character 行.
const addedLines = child.execSync('git diff -U0 -- docs/CONTRACTS.md docs/CONTRACTS_EN.md', { encoding: 'utf8' })
  .split('\n').filter(l => l.startsWith('+') && !l.startsWith('+++')).map(l => l.slice(1));
const g64 = addedLines.filter(l => /[0-9]{2,}\s*行/.test(l));
console.log(JSON.stringify({
  checks,
  pre, post,
  myDeltas: {
    cn: { lines: delta('lines', 'cn'), blank: delta('blank', 'cn'), bold: delta('bold', 'cn'), bars: delta('bars', 'cn'), bt: delta('bt', 'cn') },
    en: { lines: delta('lines', 'en'), blank: delta('blank', 'en'), bold: delta('bold', 'en'), bars: delta('bars', 'en'), bt: delta('bt', 'en') },
  },
  pairDeltaBefore: { lines: pre.en.lines - pre.cn.lines, blank: pre.en.blank - pre.cn.blank, bold: pre.en.bold - pre.cn.bold, bars: pre.en.bars - pre.cn.bars, bt: pre.en.bt - pre.cn.bt },
  pairDeltaAfter: { lines: post.en.lines - post.cn.lines, blank: post.en.blank - post.cn.blank, bold: post.en.bold - post.cn.bold, bars: post.en.bars - post.cn.bars, bt: post.en.bt - post.cn.bt },
  pairDeltaUnchanged: ['lines', 'blank', 'bold', 'bars', 'bt'].every(k => deltaOfDelta(k) === 0),
  addedLineCount: addedLines.length,
  g64Hits: g64,
}, null, 2));
