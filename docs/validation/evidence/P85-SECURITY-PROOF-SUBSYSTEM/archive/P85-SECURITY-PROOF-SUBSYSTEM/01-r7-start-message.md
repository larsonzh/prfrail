# P85 R=7 合并修复批起跑消息（2026-11-23 · §13.1.10 工件）

- 切片：P85 `SECURITY-PROOF-SUBSYSTEM`（prfrail）；批次 = 合并修复批 **R=7**（§9.2「R>3 暂停」由停点 5
  履行人工审核；用户 2026-11-23 同轮授权续跑；OB-142 已登记、不得作先例）。
- 修复集（融合表 9 项，Codex 新发现 = 0，无待定项）：
  ① Linux 目录槽 ACL fd 化（`BoundDir.AccessACL()` + 目录槽 fd 路径 + 不实登记注释改写 + 窗口方向更正；
  :155 拒绝分支逻辑随 fd 化重写，#8 自然覆盖）；② Windows DACL 指纹移入闭名对证 + 闭端根链复查；
  ③ Windows leaf 身份对证（探测前后 `FileIDAt == leaf.FileID()`，unsupported 变体签名同步）；
  ④ 测试/变异补缺（目录槽夹具判别性、OpenBoundDirVerified 变异、mutation-r6 补 2 项）；
  ⑤ Role/后缀矛盾拒绝；⑥ 非生产名义入口二次 Lstat；⑦ 拒绝路径固定四槽、空槽 ∅；
  ⑧ Windows 发布 `CreateTemp/Link` 残余按用户终裁登记 S1 独立边界（注释措辞不挂 F04 豁免）；
  ⑨ P87 数据归档（本批同做，已完成）。
- 载体：② = `prfrail-implementer`（`DeepSeek V4.1 Flash (deepseek)`）；④ = 主控机械验证；
  终验定向小读 = `Gemini 4 Carbon (copilot)` High（仅读 R=7 diff + 对应测试；§1.7.4 临时授权 +
  OB 独立留痕、不得作先例，派读前登记）。
- P87 数据归档目录声明：`docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/P87-four-read-data/`
  （四读对照表 + 四份模型报告 + r3 原文 + R6 批输入包 9 件，共 15 件；MANIFEST.md 与
  SHA256SUMS.txt 已同步；编码 md=BOM+LF、txt/patch=noBOM+LF）。
- 终验（不派全量读）：④ 门禁三作用域 + SELFTEST；本机 Windows + Ubuntu VM 双平台实跑新用例
  （非提升/非 root）；变异红/恢复绿/逐字节复原；IDE 告警三类处置（A 真实死代码删除 /
  B 跨平台变体符号归位、不得 //nolint 与假逻辑 / C 修复后重扫，逐条核三变体引用）+
  `GOOS=windows/linux/darwin go build ./...` 全绿 + 三视图 gopls unused 清零 + 处置清单。
- 停点：停点 6（修复 + 终验完成后停下报用户）；用户裁定后进 ⑧a → ⑨ → ⑧b → ⑧c → ⑩。
- 红线：不碰 CONTRACTS 四槽句与 T8a；不引入 //nolint；不动 tmp 之外未列文件；
  F17/F22 证据债口径维持（缺证不得关片）。
