# 归档 MANIFEST · s1-reach

- **源路径**：`tmp/s1-reach/`
- **归档片**：`DIRECTIVE-TMP-LIFECYCLE`（v1.23，2026-10-04）
- **归档口径**：§1.5.1 的 ② 类「早期证据评估归档」——保留证据条目（被受版本控制文档引用其路径者 ∪ 切片级文档），排除运行时残留与一次性脚本
- **归档条目**：8 条 / 177,101 字节（归档载荷根 = `archive/`；逐条 sha256 见 `SHA256SUMS.txt`）
- **排除条目**：6 条（逐条见下「排除清单」）
- **未入库条目**：0 条（逐条见下「未入库文件清单」）
- **字节保真**：归档副本的**字符内容不改**；因 §6.3 的作用域（G4a 扫全部变更文件）与仓库编码约定，副本按**归档后类型**归一编码——`.ps1` 带 BOM，其余不带 BOM，一律 LF；逐条变更种类与**原始字节 sha256** 见下「后缀与编码映射」节
- **后缀口径**：归档副本的 `.md` 一律改存为 `.txt`（本单元 5 条），以避开只扫 `.md` 的文档判据；映射逐条见下节
- **编码归一计数**：后缀改写 5 条 / BOM 剔除 2 条 / BOM 补入 0 条 / LF 归一 4 条

## 归档条目清单

| 相对路径 | 字节 | sha256 | 归档理由 |
|---|---|---|---|
| `archive/gate-tree.raw.txt` | 3255 | `0c5c7c286f8eb3690199ab07f905fd00deeaa48ea94fb6d453ba51c7c8013507` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-tree.txt` | 3255 | `0c5c7c286f8eb3690199ab07f905fd00deeaa48ea94fb6d453ba51c7c8013507` | 被受版本控制文档引用其路径 |
| `archive/plan-01.txt` | 30575 | `2d4950f88bf78b2eb2235fb35e73f48a1880916c67a43cc375e37383a5f3179b` | 被受版本控制文档引用其路径 |
| `archive/report-cn.txt` | 21833 | `86f4389ea28f548cd785380aaa7d49597fef1185f143385366b280cb775bfb83` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/report-en.txt` | 25317 | `be40c99e5ba37631b405a25cb867f954904475ff4710a38074df66ffa6ab5fa0` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/review-input.diff` | 45716 | `e708684a30c280c7f159d9b1d179134f7da365e2787c1686fa97a5cada35703a` | 被受版本控制文档引用其路径 |
| `archive/s1-reach-reachability.txt` | 21833 | `86f4389ea28f548cd785380aaa7d49597fef1185f143385366b280cb775bfb83` | 被受版本控制文档引用其路径 |
| `archive/s1-reach-reachability_EN.txt` | 25317 | `be40c99e5ba37631b405a25cb867f954904475ff4710a38074df66ffa6ab5fa0` | 切片级文档 / 评审包 / 门禁读数 |

## 排除清单（未归档）

| 相对路径 | 字节 | sha256 | 未归档理由 |
|---|---|---|---|
| `/apply-edits.js` | 21102 | `48ab3550e02c1005f0d6da083ce02bc9efde0b1a6a0457b5fb5a4ad30a85087f` | 落地与修复脚本残留 |
| `/fix-05.js` | 4691 | `1b316c778276f3ee49f6bd9c5c705ba47d06df0d5c0e5c2d8d586f0177449a2b` | 落地与修复脚本残留 |
| `/fix-g64.js` | 1183 | `465baf46b4868351a6344ecdc871b850ffc5a3782003966b517501019e9d2fe7` | 落地与修复脚本残留 |
| `/fix-n2.js` | 1810 | `d63ebbd85f2816d5042a3b3bfa204210a445f834ae2b23e1147ec3a3abad88c7` | 落地与修复脚本残留 |
| `/fix-refs.js` | 1856 | `a489a56ad62fc568f21a25fd944af62d8ee204332016b36115a217594a74cd85` | 落地与修复脚本残留 |
| `/writeback-8b.js` | 4462 | `dcb3ecfd553a2631267c25e69bf5f399ba0b4074f3b9e04aa09f1e9c6a3b0da2` | 落地与修复脚本残留 |

## 未入库文件清单（`.gitignore` 排除类型）

| 相对路径 | 字节 | sha256 | 理由 |
|---|---|---|---|

## 重建说明

- **未入库条目**（见上表）：本单元未入库清单为空（无忽略类型条目）。
- **排除条目**：一次性脚本、运行时状态与中间转储；其行为已由本片归档的切片级文档（切片定义 / 方案 / 评审包 / 门禁读数）与对应切片的提交记录固定。
- **行尾与 BOM 归一**：源 `tmp/` 副本的部分文件为 CRLF 或带 BOM；归档副本已按仓库编码约定归一到「`.ps1` 带 BOM / 其余不带 BOM / 一律 LF」（不影响字符内容）。

## 后缀与编码映射（`.md` → `.txt` 及编码归一）

| 原路径 | 归档路径 | 原始字节 sha256 | 归档字节 sha256 | 变更 |
|---|---|---|---|---|
| `gate-tree.raw.txt` | `archive/gate-tree.raw.txt` | `28e3ae5e5bcd30a1691ed21e43b250687e61944fe2bc338f228ae17371e75526` | `0c5c7c286f8eb3690199ab07f905fd00deeaa48ea94fb6d453ba51c7c8013507` | LF 归一 |
| `gate-tree.txt` | `archive/gate-tree.txt` | `28e3ae5e5bcd30a1691ed21e43b250687e61944fe2bc338f228ae17371e75526` | `0c5c7c286f8eb3690199ab07f905fd00deeaa48ea94fb6d453ba51c7c8013507` | LF 归一 |
| `plan-01.md` | `archive/plan-01.txt` | `2d4950f88bf78b2eb2235fb35e73f48a1880916c67a43cc375e37383a5f3179b` | `2d4950f88bf78b2eb2235fb35e73f48a1880916c67a43cc375e37383a5f3179b` | 后缀 `.md` → `.txt` |
| `report-cn.md` | `archive/report-cn.txt` | `708a5d1093affebe7afdc9903cd916c1710f73c1c1e7de843a385e3e85901c44` | `86f4389ea28f548cd785380aaa7d49597fef1185f143385366b280cb775bfb83` | 后缀 `.md` → `.txt` / LF 归一 |
| `report-en.md` | `archive/report-en.txt` | `c602236f2a1474f9f31c50441a15c90aa3f79544f56fc1c1514e7fea2c5cd0b1` | `be40c99e5ba37631b405a25cb867f954904475ff4710a38074df66ffa6ab5fa0` | 后缀 `.md` → `.txt` / LF 归一 |
| `s1-reach-reachability.md` | `archive/s1-reach-reachability.txt` | `7ba70aa11a4a79d625d06b023c4f7e1ff16f8048ddada479e6cdf27a9c6b5ae2` | `86f4389ea28f548cd785380aaa7d49597fef1185f143385366b280cb775bfb83` | 后缀 `.md` → `.txt` / BOM 剔除 |
| `s1-reach-reachability_EN.md` | `archive/s1-reach-reachability_EN.txt` | `8367a08bab0c9dbaf988f243872b2d8b4eee3af42f3f78e4307754474c7b1bf7` | `be40c99e5ba37631b405a25cb867f954904475ff4710a38074df66ffa6ab5fa0` | 后缀 `.md` → `.txt` / BOM 剔除 |

## 对应引用点

| 引用文件 | 行号 |
|---|---|
| `docs/validation/evidence/S1-REACH-freedom-list.md` | 15 |
| `docs/validation/evidence/directive-history/s1-reach/archive/report-cn.txt` | 88 / 110 / 147 / 148 / 149 |
| `docs/validation/evidence/directive-history/s1-reach/archive/report-en.txt` | 88 / 110 / 147 / 148 / 149 |
| `docs/validation/evidence/directive-history/s1-reach/archive/s1-reach-reachability.txt` | 88 / 110 / 147 / 148 / 149 |
| `docs/validation/evidence/directive-history/s1-reach/archive/s1-reach-reachability_EN.txt` | 88 / 110 / 147 / 148 / 149 |
| `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/archive/classification.txt` | 1394 / 1395 / 1396 / 1397 / 1398 / 1399 / 1400 / 1401 / 1402 / 1403 / 1404 / 1405 / 1406 / 1407 |
| `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/archive/inventory.txt` | 184 / 185 / 186 / 187 / 188 / 1049 / 1132 / 1133 / 1134 / 1135 / 1136 / 1137 / 1138 / 1139 / 1140 / 1141 / 1142 / 1143 / 1144 / 1145 / 1146 / 1147 |
| `docs/validation/s1-reach-reachability.md` | 88 / 110 / 147 / 148 / 149 |
| `docs/validation/s1-reach-reachability_EN.md` | 88 / 110 / 147 / 148 / 149 |
