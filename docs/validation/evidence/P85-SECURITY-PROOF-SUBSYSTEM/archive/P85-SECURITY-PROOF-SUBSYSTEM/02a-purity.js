// P85 02a: byte-level purity check — current tree must equal HEAD + exactly the intended edits.
const { execFileSync } = require('child_process');
const fs = require('fs');

const headLines = f => execFileSync('git', ['show', 'HEAD:' + f], { encoding: 'utf8', maxBuffer: 64 * 1024 * 1024 })
  .replace(/^\uFEFF/, '').split('\n');
const curLines = f => fs.readFileSync(f, 'utf8').replace(/^\uFEFF/, '').split('\n');

const PARA_A_CN = '对 `agent-runner-cli` 通道的可用性探测，探针可执行文件的安装来源必须可证明不可变：Windows 上以稳定绝对路径命名且其所在安装目录的 ACL 对本用户令牌关闭写入/删除；Linux 上其安装目录全路径各组件对本用户不可写，且可执行文件为 ELF 或对其 shebang 解释器同样满足不可变条件；其余平台 fail-closed。该证明以链配置目录派生路径下的 store-local 载体 `ai-install-source-proof.json` 持久化（写入方 = ProofRail 检查命令、读取方 = 检查命令与探针执行器，每次启动前重读并实时重验路径→文件身份→摘要；载体不得由调用方自选路径）。安装来源不可证明时，不得记录 `available`（以 `cli_unavailable` 分类拒绝）。本载体的 Windows ACL 与 Linux 模式/属主判定分别由平台专属实现承担，判定失败一律 fail-closed。';
const PARA_B_CN = 'secret 认证模式的 `agent-runner-cli` 探测在无法证明子进程禁写未脱敏落盘时，不得解析凭据、不得启动子进程（零启动拒发 BYOK）：证明依据 = 白名单 argv（无日志目的地）、环境白名单、仅继承标准流句柄、受限工作目录与有界内存输出，外加平台写限制判定（Sandbox 检查层）；任一项不可证即拒发，拒绝证据只用稳定闭集分类。host-auth 路径不适用本门禁。';
const PARA_A_EN = 'For availability probing of the `agent-runner-cli` channel, the install source of the probe executable must be provably immutable: on Windows it is named by a stable absolute path and the ACL of its install directory denies write and delete to the current user\'s token; on Linux every component of its install directory path is not writable by the current user, and the executable is either ELF or its shebang interpreter satisfies the same immutability condition; all other platforms fail closed. The proof is persisted as the store-local `ai-install-source-proof.json` artifact under a path derived from the chain configuration directory (written by the ProofRail check command, read by the check command and the probe executor, re-read and re-verified before every start for path, file identity, and digest; callers cannot choose its path). If the install source cannot be proven, `available` must not be recorded and the verdict is refused as `cli_unavailable`. Windows ACL and Linux mode/owner evaluation are owned by platform-specific implementations, and any failed evaluation fails closed.';
const PARA_B_EN = 'In secret authentication mode, probing `agent-runner-cli` must not resolve credentials and must not start the child process when it cannot be proven that the child cannot write unredacted output to disk (zero-start BYOK refusal): the proof basis is a whitelisted argv with no log destination, a whitelisted environment, inheritance of standard stream handles only, a confined working directory with bounded in-memory output, plus platform write-restriction evaluation in the sandbox check layer; if any item cannot be proven the verdict is refused, and refusal evidence uses only stable closed-set classifications. The host-auth path is not subject to this gate.';

function build(f, preflightIdx, credIdx, enumOld, enumNew, paraA, paraB) {
  const L = headLines(f);
  const out = [];
  for (let i = 0; i < L.length; i++) {
    let line = L[i];
    if (i === preflightIdx) {
      if (!line.includes(enumOld)) throw new Error(f + ': enum anchor missing at ' + (i + 1));
      out.push(line.replace(enumOld, enumNew));
      out.push('', paraA);
    } else if (i === credIdx) {
      out.push(line, '', paraB);
    } else {
      out.push(line);
    }
  }
  return out;
}

function cmp(name, expected, actual) {
  if (expected.length !== actual.length) { console.log(name + ' LEN FAIL ' + expected.length + ' vs ' + actual.length); return; }
  let bad = 0;
  for (let i = 0; i < expected.length; i++) if (expected[i] !== actual[i]) { if (bad < 3) console.log(name + ' DIFF L' + (i + 1)); bad++; }
  console.log(name + ' ' + (bad === 0 ? 'EXACT-MATCH OK (lines=' + actual.length + ')' : bad + ' differing lines'));
}

cmp('CN', build('docs/CONTRACTS.md', 213, 215, '`timeout` 和 `unknown`', '`timeout`、`byok_confinement_unproven` 和 `unknown`', PARA_A_CN, PARA_B_CN), curLines('docs/CONTRACTS.md'));
cmp('EN', build('docs/CONTRACTS_EN.md', 207, 209, '`timeout`, and `unknown`', '`timeout`, `byok_confinement_unproven`, and `unknown`', PARA_A_EN, PARA_B_EN), curLines('docs/CONTRACTS_EN.md'));
