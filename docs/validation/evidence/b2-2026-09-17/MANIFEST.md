# B2 evidence bundle manifest (2026-09-17)

All artifacts are sanitized copies: the local user name is replaced by <user>; no token value
appears anywhere (verified by an automated scan for gho_/ghp_/github_pat_/40-hex patterns).
The machine-ops/ copies are the archived snapshot of the maintained scripts now living in
tools/agent-probe/appcontainer-b2/ (same bytes at archive time); scratch output of those scripts
goes to the git-ignored tmp/b2/ tree, never next to the versioned scripts.
Every text artifact is stored as UTF-8 with LF: repository encoding convention applies, so the
scripts (*.ps1) carry a BOM while logs, results and JSON do not. Exception:
measuring/probe10-dir-listing.raw.txt keeps the raw OEM code-page bytes captured from cmd; its
decoded counterpart is measuring/probe10-dir-listing.txt. The live scratch copies under tmp/
are normalized identically (normalize-encoding.ps1 -Root tmp\b2), so a live copy and its archived
counterpart differ only by the user-name sanitization. Hashes below cover the stored bytes.

- boundary\boundary-battery.txt  1087 bytes  sha256 d7103cf26122767851b22844b11f474e5b9acc0bc81bbf2a3d9e285c4a80257b
- container-exec\argv.json  4538 bytes  sha256 6cf843d80b15993d1084b9cdf6035ee6ba364bcb1f7a71ba1b2c596dfe6395ba
- container-exec\exec-check-pre-subst-stderr.txt  2075 bytes  sha256 45e742bb24deb87d6077140b38bf3a9351670876eb9454b204d9132e47988774
- container-exec\result.json  1706 bytes  sha256 308e6a9fb8b979a256b6ee71913cd3f98dde7acaf86cec4ecb0dc57758c11dc9
- container-exec\stdout.txt  70 bytes  sha256 29c4198d94948906858abafe036110547bae7dddc92e2805b2e0a9a2b33ac7ba
- machine-ops\ac-lib.ps1  13872 bytes  sha256 fd13e7197488160a7c7b5e1523504bc3eae608f5540ab06722d8bd6988719a21
- machine-ops\apply-b2-loopback.ps1  5426 bytes  sha256 8b418578168330a0aac2f1526a6697e15f89d5494af89cd03e76e92ecd85fe94
- machine-ops\assemble-evidence.ps1  12892 bytes  sha256 672da48c1e08b9005b79f50c7da25e18c9956ccabf8271d64673ff49290cebb9
- machine-ops\b2-loopback-lib.ps1  13818 bytes  sha256 69fa3af9a4bb448df1e46bb7016a5dd9e19ca763eb48bcaac92500ed10b4bfb1
- machine-ops\boundary-battery.ps1  7015 bytes  sha256 66c0ef577612bf8e5fdde32ac2f05e9462bc6fe6029ffe509f5acadfa696c856
- machine-ops\container-candidate.ps1  17076 bytes  sha256 3f6ab7df16c64a9ed46120b5c4a285dceac0804d438537ba52b085776508edc9
- machine-ops\container-modelcall.ps1  12519 bytes  sha256 8ecedfc4436a6dfcb0eb6455657dc924688db0bbcff9c2107dd56f8be9e6415e
- machine-ops\container-nodecall.ps1  10549 bytes  sha256 d017ee189eff034591dacc63e5a7f3e998e69a09be0316ce12eda8c6a46c635b
- machine-ops\Invoke-B2LoopbackRehearsal.ps1  4526 bytes  sha256 cc455d46dfd677866fc6b756caef5e9c1dd0c3b977754be01d1c7daecd1fc009
- machine-ops\node-client.js  6853 bytes  sha256 f91f9e508f2488be8bf9c0931821901bc5a34f7433327f28387532ba7f25b7d2
- machine-ops\normalize-encoding.ps1  7214 bytes  sha256 3bd3205184a80315acf69739d16b5660704596ae14e246e06b4f544416864b56
- machine-ops\path-probe.ps1  6150 bytes  sha256 601f701573e4602749e247a5c0ea5487b46fbd7dc6abc30850ba775dc6ad2797
- machine-ops\remove-b2-container-profile.ps1  1479 bytes  sha256 4ba9dbc5e17007a23c9eaa24e34efbdbefb68b4580de59a364ba38620066e6fd
- machine-ops\revert-b2-loopback.ps1  5323 bytes  sha256 f3dd1a9493add7f6fdb2e46fd86c61e0350d7092b86f0e577648a28fbe980b49
- machine-ops\Test-B2Tooling.ps1  6408 bytes  sha256 c6e7ab7532e472cbc63a770c5f73287dd378f5d796b4d05015d913ee53c0f0d1
- machine-ops\verify-b2-loopback.ps1  2520 bytes  sha256 c15578ae3eb74ee7789580cc0432f36b1094db62060fc0fb3caf83e86a0bde1c
- measuring\probe10-dir-listing.raw.txt  357 bytes  sha256 d03e1798378eedff6f8932d7fc92e7d882d53fa15edd9a47ab2d0fd41cc89ae6
- measuring\probe10-dir-listing.txt  364 bytes  sha256 47b7e3ad0000f9e84b36d82ef88850bcad47a32abeea1c09de15ff3999668476
- measuring\probe13-path-spelling.txt  152 bytes  sha256 8dbb8848a1760fdc7052c31c33f42a69d97078372e7e1b8c996ece795e477ddb
- measuring\profile-detection.txt  284 bytes  sha256 81216f67aa56a32e8f4468ec96ad80ef92a9c7594ec3e31e4ae703717fcad7aa
- measuring\profile-detection-probe.ps1  1498 bytes  sha256 a30152c690339ae6e3d0aac87da06720c89fd47c6b9d05d7f6613c484ee48057
- measuring\ws-child-console.txt  0 bytes  sha256 e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
- measuring\ws-result.txt  370 bytes  sha256 560f5cfbfba7951a79eb47f6967588d651ca1f61fa31be568b30ffa5d0ada884
- proxy-decisions\allow-run.jsonl  2102 bytes  sha256 97740fd630dcba592a6aade5c23fa3a91ec8561f6f65453c531fac8390108847
- proxy-decisions\allow-run-result.json  1256 bytes  sha256 1c25cd175039e72a1c24d3181766bfcfdc8535060566741ce82dbbc7104b8c5b
- proxy-decisions\allow-run-stderr.txt  117 bytes  sha256 11b3375effa193ec49dd01d4ea2620f237d44ec305abe299525b1b1671c78583
- proxy-decisions\allow-run-usage.json  381 bytes  sha256 7dc303dd60add54709202128b2c8a5205b549e62e317feda89e4b26340493c38
- proxy-decisions\candidate-session-log.txt  6240 bytes  sha256 03da8663a20b5cba852de393244e5210037bccbf9150f5692b930558d3e48345
- proxy-decisions\deny-run.jsonl  1962 bytes  sha256 bf8a65dae4119517104f203d143ed3fd2f0ad71c72a3875c1047e9bf518ef1b4
- proxy-decisions\deny-run-result.json  1584 bytes  sha256 1051cb7623814cd35f4f6eaa0b98601cb36bf6319af562fd02fee8ce0f2e043a
- restorability\rehearsal.txt  4755 bytes  sha256 ff457d14a9804c25d485f8e28f94e5dbe1866bd140c9c5435c9f58e5ebd76229
- restorability\verify-current-state.txt  370 bytes  sha256 51803177624f1c4825463f5e08fc6d46da454353afd67350647793443301fb7f
- token-exchange\node-out.txt  826 bytes  sha256 ef60519db57c585f031924496412f50c06130d5aeed095a4f435a0adf0bca368
- token-exchange\node-result.txt  826 bytes  sha256 ef60519db57c585f031924496412f50c06130d5aeed095a4f435a0adf0bca368
- token-exchange\proxy.jsonl  2326 bytes  sha256 073d6d5eea6b2ebf5893732843140d5155d4bfd9fabaa293646923938fb4e340
- token-exchange\token-err.txt  13 bytes  sha256 e5073381f828efcb033b8e6feee3df7bbc4a0182b68809ff10d950b81572856f
