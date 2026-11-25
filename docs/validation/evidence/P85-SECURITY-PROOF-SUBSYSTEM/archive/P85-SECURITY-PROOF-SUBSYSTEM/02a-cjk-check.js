// P85 02a: verify every CJK char in new CN contract text exists in HEAD *.md corpus.
const { execFileSync } = require('child_process');

const texts = [
  '对 `agent-runner-cli` 通道的可用性探测，探针可执行文件的安装来源必须可证明不可变：Windows 上以稳定绝对路径命名且其所在安装目录的 ACL 对本用户令牌关闭写入/删除；Linux 上其安装目录全路径各组件对本用户不可写，且可执行文件为 ELF 或对其 shebang 解释器同样满足不可变条件；其余平台 fail-closed。该证明以链配置目录派生路径下的 store-local 载体 `ai-install-source-proof.json` 持久化（写入方 = ProofRail 检查命令、读取方 = 检查命令与探针执行器，每次启动前重读并实时重验路径→文件身份→摘要；载体不得由调用方自选路径）。安装来源不可证明时，不得记录 `available`（以 `cli_unavailable` 分类拒绝）。本载体的 Windows ACL 与 Linux 模式/属主判定分别由平台专属实现承担，判定失败一律 fail-closed。',
  'secret 认证模式的 `agent-runner-cli` 探测在无法证明子进程禁写未脱敏落盘时，不得解析凭据、不得启动子进程（零启动拒发 BYOK）：证明依据 = 白名单 argv（无日志目的地）、环境白名单、仅继承标准流句柄、受限工作目录与有界内存输出，外加平台写限制判定（Sandbox 检查层）；任一项不可证即拒发，拒绝证据只用稳定闭集分类。host-auth 路径不适用本门禁。'
];

const files = execFileSync('git', ['ls-tree', '-r', 'HEAD', '--name-only'], { encoding: 'utf8', maxBuffer: 64 * 1024 * 1024 })
  .split('\n').map(s => s.trim()).filter(f => f.endsWith('.md'));

let corpus = '';
for (const f of files) {
  try { corpus += execFileSync('git', ['show', `HEAD:${f}`], { encoding: 'utf8', maxBuffer: 64 * 1024 * 1024 }); } catch (e) { }
}
const has = new Set(corpus);
console.log('corpus md files=' + files.length + ' chars=' + corpus.length);

const isCjk = c => c >= '\u4e00' && c <= '\u9fff';
const missing = new Map();
for (const t of texts) {
  for (const c of t) { if (isCjk(c) && !has.has(c)) missing.set(c, (missing.get(c) || 0) + 1); }
}
console.log('CJK total needed=' + [...new Set(texts.join(''))].filter(isCjk).length);
if (missing.size === 0) console.log('MISSING: none');
else console.log('MISSING: ' + [...missing.entries()].map(([c, n]) => `${c}(${n})`).join(' '));
