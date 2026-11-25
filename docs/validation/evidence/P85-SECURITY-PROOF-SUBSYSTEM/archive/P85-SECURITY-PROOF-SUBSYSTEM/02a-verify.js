// P85 02a: verify encoding + locate inserted paragraphs on disk.
const fs = require('fs');
const marks = [
  ['CN-A', '对 `agent-runner-cli` 通道的可用性探测'],
  ['CN-B', 'secret 认证模式的 `agent-runner-cli` 探测'],
  ['EN-A', 'For availability probing of the `agent-runner-cli` channel'],
  ['EN-B', 'In secret authentication mode, probing `agent-runner-cli`'],
];
for (const f of ['docs/CONTRACTS.md', 'docs/CONTRACTS_EN.md']) {
  const b = fs.readFileSync(f);
  const t = b.toString('utf8').replace(/^\uFEFF/, '');
  const L = t.split('\n');
  console.log('== ' + f + ' BOM=' + (b[0] === 0xEF && b[1] === 0xBB && b[2] === 0xBF)
    + ' CR=' + (b.includes(13) ? 'Y' : 'N') + ' LF=' + (L.length - 1));
  for (const [name, m] of marks) {
    const i = L.findIndex(x => x.startsWith(m));
    if (i >= 0) {
      const l = L[i];
      console.log('  ' + name + ' L' + (i + 1) + ' len=' + l.length
        + ' head=' + JSON.stringify(l.slice(0, 50)) + ' ... tail=' + JSON.stringify(l.slice(-48)));
    }
  }
  // enumerate closed-set line containing byok
  const i2 = L.findIndex(x => x.includes('byok_confinement_unproven'));
  if (i2 >= 0) {
    const l = L[i2];
    const k = l.indexOf('`response_invalid`');
    console.log('  CLOSED-SET L' + (i2 + 1) + ' seg=' + JSON.stringify(l.slice(k - 6, k + 90)));
  }
}
