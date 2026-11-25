// P85 02a: assert exact anchors exist exactly once before editing.
const fs = require('fs');
const cn = fs.readFileSync('docs/CONTRACTS.md', 'utf8').replace(/^\uFEFF/, '');
const en = fs.readFileSync('docs/CONTRACTS_EN.md', 'utf8').replace(/^\uFEFF/, '');

const checks = [
  ['CN-A', cn, '`response_invalid`、`timeout` 和 `unknown`。某通道不可用不得自动切换 provider、模型或其他通道；非 AI 功能可继续运行，但对应 AI 入口保持禁用并显示结构化原因。\n\n`windows-credential:<target>` 引用由当前用户的 Windows Credential Manager generic credential 实现，持久性为 local-machine；'],
  ['CN-B', cn, '非 Windows 平台在各自 SecretStore backend 完成前明确 fail-closed，不共享或猜测 Windows 凭据。\n\n链配置可选顶层 `ai`'],
  ['EN-A', en, '`response_invalid`, `timeout`, and `unknown`. An unavailable channel cannot automatically switch provider, model, or channel. Non-AI features may continue, while the affected AI entry point remains disabled and displays its structured reason.\n\nA `windows-credential:<target>` reference uses a generic credential'],
  ['EN-B', en, 'Non-Windows platforms fail closed until their own SecretStore backend exists; they do not share or infer Windows credentials.\n\nThe optional top-level `ai` section'],
];

for (const [name, src, anchor] of checks) {
  const n = src.split(anchor).length - 1;
  console.log(name + ' occurrences=' + n + (n === 1 ? ' OK' : ' FAIL'));
}
// G6-4 precheck on new CN text: digits{2,} followed by optional spaces then 行
const newCn = ['对 `agent-runner-cli` 通道', 'secret 认证模式的 `agent-runner-cli` 探测'];
for (const t of newCn) console.log('G6-4 pattern in new text: ' + /\d{2,}\s*行/.test(t));
