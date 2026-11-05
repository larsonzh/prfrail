# ⑥ 独立扫描原文 — `G4B-CORPUS-DECOUPLE`

性质：`independent-reviewer` + `MAI-Code-1.1-Flash (copilot)` 的**最终消息逐字转录**。审查时点 2026-11-05，只读；未向该环提供任何其它审查方的清单。

---

## Section A
INDEPENDENT SCAN: PASS
根据只读核对的实现和断言，没有发现违反这两条硬边界的情况：G4b 语料枚举确实被改成对象库读取，并且缓存键/版本逻辑仍然绑定要件，而且自测对关键回归有较强变异敏感性。

## Section B
(no findings)

## Section C
1) 对象库枚举真实生效，且不再依赖索引。读取主控落盘的 tools/gates/lib/ctx.js:536-551，listAllMarkdown 的实现是 git(['ls-tree', '-r', '-z', '--name-only', ref])，随后用 '\0' 分割，并用 /\.md$/i 过滤；这正是“对象库 + NUL 分隔 + JS 过滤”的组合，且注释明确指出不能用 ls-files 或 *.md pathspec，原因是会让 G4b 语料收缩并制造假阳性。也就是说，问题不只是报告里说了，而是实现上确实走了 tree/object store ，而不是 index。读取主控落盘的 tools/gates/lib/ctx.js:540-551 还特别说明了 -z 的必要性：非 ASCII 路径会被 quotePath 包裹，newline 分割会把路径弄坏，故 NUL 分隔是 load-bearing。

2) 缓存键仍然绑定解析 sha / whitelist 文本摘要 / 判据源码摘要 / 格式版本，并且 fail-closed 在缓存前。读取主控落盘的 tools/gates/lib/criteria/g4b.js:43-89 和 99-170：CACHE_VERSION 是 v2，keyString() 返回 [CACHE_VERSION, rs, kind, source, sha256(text), sourceDigest()]，其中 rs 来自 resolvedSha(ctx)，text 来自 whitelist.text，sourceDigest() 哈希了 WHITELIST_REL、cjkChars.toString() 与 runBody()；同时 run() 在缓存前就先读取 whitelist，并在 whitelist.error 时抛 UsageError，避免缓存掩盖坏白名单。validValue() 也要求 v === 2 且 corpus 是字符串数组，和缓存写入时的 { v: 2, corpus: [...] } 保持一致。缓存命中只读 corpus，不写 verdict，且任何 get 失败/读坏/写坏都 fail-open 回全量计算，这也和硬边界一致。

3) 新增断言具备变异敏感性，且覆盖了关键回归面。读取主控落盘的 tools/gates/selftest.js:6579-6663，T114/T117 明确锁定 object-store + -z + .md 过滤；T115 锁定“staged markdown deletion 后 cold/warm 相等且 verdict 不变”；T109/T110/T111 测试解析 sha、whitelist 摘要、sourceDigest 变化时必须 miss，而只有同一 key/同一判据源码/同一白名单才可 hit；T112 锁定 fail-open 与 key 形状，并确认 g4bCacheDisabled / no-sha 时都退回全量计算；T116 锁定缓存存储的 v 值必须等于 CACHE_VERSION。检索 tools/gates 中 listAllMarkdown 的唯一消费者也只有 G4b 自身和 selftest，没有发现其它判据被这次变动连带地依赖到这个枚举函数；因此“其它判据被波及”这点在当前仓库代码中未发现。

## Section D
- 改回 ls-files / 改回 newline split：T114、T117、T115 会红；T116 不会红，因为它只看 v 值和 corpus 结构，不管文件枚举实现。
- 删掉 -z 与 '\0' 分割：T114、T117 会红；T115 也会红；T116 不会红。
- 删除 resolved sha 段：T109 会红；T110/T111 不会红；T116 不会红。
- 删除 whitelist 摘要段：T110 会红；T109/T111 不会红。
- 删除 sourceDigest 段：T111 会红；T109/T110 不会红。
- 删除 validValue 里的 v===2 或软编码 v 绑定：T112 fail-open 2 会红；T116 会红；T108/T109/T110/T111 一般不会因这条单独变红。
- 删除 cacheUsable() 对 no-sha 的保护：T112 的“git 无法解析 sha 时静默全量计算、无缓存记录”会红；其余 T109-T111 不会特定红。
- 删除 whitelist.error 的 fail-closed：我没有看到当下断言覆盖“坏白名单错误应在缓存之前失败”的案例；这是目前没有覆盖的不变量。另一个没有覆盖的不变量是“另一个 consumer 不会再调用 listAllMarkdown”，但代码搜索显示当前 tools/gates 内没有其它消费者，故这条未覆盖风险已被消除。

## Section E
我无法验证的部分：我只读取了任务明确列出的主控落盘文件和实现文件；没有读取别的审查方清单，也没有执行 git/selftest；因此我不能对“其它审查团队是否还看到了别的命题”作独立确认。另一个无法验证点是：我未读取 diff.patch 或其它审查包内容，所以我不声明其余结论。

会推翻我的 PASS 的条件：如果在实际仓库里还有一条不在当前读到的实现中、但在未落盘分支/未读文件里调用 listAllMarkdown；如果后续新增了第二个 listAllMarkdown consumer；如果实际对象库语义不同于这里的 HEAD/object-store 表述；如果白名单 loader 在某个代码路径中返回 error 后仍被缓存绕过；或者如果该实现和当前断言的对象库语义存在未覆盖的环境差异。
INDEPENDENT SCAN: PASS
