#!/usr/bin/env node
"use strict";

// run-selftest.js -- ③ tester's independent, adversarial self-test for the
// DIRECTIVE-GLM-EVAL harness (T-GE-01 .. T-GE-09).
//
// Design rules:
//   * Reads platform artifacts read-only (bundle/**, key.json, score.js, hash.js, contam.js).
//   * Invokes the platform CLIs via child_process (they auto-run main() on require).
//   * Writes fixtures ONLY under ./fixtures (tester-owned).
//   * T-GE-05 requires a mutation assertion: it temporarily edits ../../score.js and
//     MUST restore it byte-identically (sha256 before == sha256 after).
//   * T-GE-10 is the dataset pre_registration zero-leak guard (② fix #1).
//   * T-GE-11 is the denominator-derivation guard (② fix #3): expected_total is read
//     from key.json, not hardcoded.
//   * T-GE-12 falsifies the two NEW guards by temporarily mutating contam.js (strip \b)
//     and score.js (hardcode 7); both MUST be restored byte-identically.

const fs = require("fs");
const path = require("path");
const crypto = require("crypto");
const { spawnSync } = require("child_process");

const SELFTEST = __dirname;
const PLATFORM = path.resolve(SELFTEST, "..");
const BUNDLE = path.join(PLATFORM, "bundle");
const KEY = path.join(PLATFORM, "key.json");
const SCORE = path.join(PLATFORM, "score.js");
const HASH = path.join(PLATFORM, "hash.js");
const CONTAM = path.join(PLATFORM, "contam.js");
// The platform was relocated OUTSIDE the workspace (%TEMP%\glm-eval) for isolation,
// so the dataset can no longer be reached by a fixed relative path. Resolve it from
// an env override, then cwd-relative candidates, then the legacy relative path, then
// a documented absolute fallback.
function resolveDataset() {
  const rel = path.join("docs", "validation", "directive-glm-eval-dataset.json");
  const candidates = [];
  if (process.env.GLM_EVAL_DATASET) candidates.push(process.env.GLM_EVAL_DATASET);
  candidates.push(path.resolve(PLATFORM, "..", "..", rel)); // legacy: <ws>/tmp/glm-eval -> <ws>/docs
  candidates.push(path.resolve(process.cwd(), rel));
  candidates.push(path.resolve(process.cwd(), "..", "prfrail", rel));
  candidates.push(path.resolve(process.cwd(), "prfrail", rel));
  candidates.push(path.resolve(process.cwd(), "..", rel));
  candidates.push("D:\\LZProjects\\prfrail\\" + rel);
  for (const c of candidates) {
    try { if (fs.existsSync(c)) return c; } catch (e) { /* ignore */ }
  }
  return candidates[candidates.length - 1];
}
const DATASET = resolveDataset();
const FIX = path.join(SELFTEST, "fixtures");

const ORDER = ["00-plan.md", "01-clause-cn.txt", "02-diff.js", "03-prompt.txt"];
const REPORTED_INPUT_HASH = "bbeaf2f7db40547d26b1a338608c0b9c1679fd7b08edbf349c0fbadaaa8008b6";
const S8_EXPECTED_LINES = [
  "@@ -40,7 +84,10 @@ function limits() {",
  "-  const max = 10;",
  "-  const min = 0;",
  "+  const min = 0;",
  "+",
  "+  const max = 10;",
  "   return { min, max };",
  " }",
];

const RESULTS = [];
const NOTES = [];
const FINDINGS = [];
const OUT_LINES = [];
const _origLog = console.log;
console.log = function (...args) {
  const s = args.map(String).join(" ");
  OUT_LINES.push(s);
  _origLog(s);
};
function check(id, name, ok, detail) {
  RESULTS.push({ id, name, ok: !!ok, detail: detail || "" });
}
function note(id, text) {
  NOTES.push({ id, text });
}
function finding(id, text) {
  FINDINGS.push({ id, text });
}
const sha256 = (b) => crypto.createHash("sha256").update(b).digest("hex");
const readText = (p) => fs.readFileSync(p, "utf8");
const stripBom = (t) => t.replace(/^\uFEFF/, "");
const bundleText = (f) => readText(path.join(BUNDLE, f));
const bundleLines = (f) => stripBom(bundleText(f)).split(/\r?\n/);
function countOcc(hay, needle) {
  if (!needle) return 0;
  let c = 0, i = 0;
  while ((i = hay.indexOf(needle, i)) !== -1) { c++; i += needle.length; }
  return c;
}
function runCli(script, args) {
  return spawnSync(process.execPath, [script].concat(args), { encoding: "utf8" });
}
function writeFixture(name, text) {
  const p = path.join(FIX, name);
  fs.writeFileSync(p, text, { encoding: "utf8" }); // no BOM
  return p;
}
function scoreFixture(filePath) {
  const r = runCli(SCORE, [filePath, KEY]);
  if (r.status !== 0) throw new Error("score.js failed: " + r.stderr);
  return JSON.parse(r.stdout);
}

/* ------------------------------------------------------------------ */
/* Fixture corpus (written fresh on every run)                        */
/* ------------------------------------------------------------------ */

function buildFixtures() {
  fs.mkdirSync(FIX, { recursive: true });
  const w = (n, s) => writeFixture(n, s.endsWith("\n") ? s : s + "\n");

  w("leg-perfect.txt", [
    "预审发现（合成评估 · 全命中夹具）",
    "",
    "1. Critical — 02-diff.js:58：解析异常被吞并后 return true，失败未上抛。",
    "2. High — 02-diff.js:80：成员路径未做分隔符归一即比较。",
    "3. High — 02-diff.js:105：摘要复算取错列，比较恒真。",
    "4. Medium — 02-diff.js:93：清单独有成员一支未实现。",
    "5. Medium — 02-diff.js:125：子项四（重建注记）未实现。",
    "6. Medium — 02-diff.js:131：作用域未使用入参，改为全树扫描。",
    "7. Low — 02-diff.js:14：顶部注释与实现类别数不符。",
    "",
    "PRE-REVIEW: FINDINGS",
  ].join("\n"));

  w("leg-empty.txt", ["未发现明显问题。", "", "PRE-REVIEW: PASS"].join("\n"));
  w("leg-noconcl.txt", ["未检出问题。", "（未给出结论行）"].join("\n"));
  w("leg-trailing-ws.txt", ["未发现明显问题。", "", "PRE-REVIEW: PASS   "].join("\n"));

  w("leg-sev-exact.txt", ["1. Critical — 02-diff.js:58：x", "", "PRE-REVIEW: FINDINGS"].join("\n"));
  w("leg-sev-off1.txt", ["1. High — 02-diff.js:58：x", "", "PRE-REVIEW: FINDINGS"].join("\n"));
  w("leg-sev-off2.txt", ["1. Low — 02-diff.js:58：x", "", "PRE-REVIEW: FINDINGS"].join("\n"));

  w("leg-fps.txt", [
    "1. 疑似 — 02-diff.js:999：不存在的行。",
    "2. 疑似 — docs/other.js:5：bundle 外路径。",
    "3. 疑似 — 02-diff.js:194：纯重排区误报。",
    "",
    "PRE-REVIEW: FINDINGS",
  ].join("\n"));

  w("leg-extra.txt", ["1. 观察 — 02-diff.js:135：main 入口，非已知种子。", "", "PRE-REVIEW: FINDINGS"].join("\n"));

  w("leg-qcheck.txt", [
    "1. Critical — 02-diff.js:58：a",
    "2. Medium — 02-diff.js:80：b",
    "3. High — 02-diff.js:105：c",
    "",
    "PRE-REVIEW: PASS WITH FIXES",
  ].join("\n"));

  w("leg-acceptlines.txt", [
    "1. Critical — 02-diff.js:137：非锚点接受行（S1）。",
    "2. High — 02-diff.js:81：非锚点接受行（S2）。",
    "3. High — 02-diff.js:107：非锚点接受行（S3）。",
    "4. Medium — 02-diff.js:90：非锚点接受行（S4）。",
    "5. Medium — 02-diff.js:123：非锚点接受行（S5）。",
    "6. Medium — 02-diff.js:130：非锚点接受行（S6）。",
    "",
    "PRE-REVIEW: FINDINGS",
  ].join("\n"));

  w("leg-sev-cn.txt", [
    "1. 严重度：严重 — 02-diff.js:58：a",
    "2. 严重度：高 — 02-diff.js:80：b",
    "",
    "PRE-REVIEW: FINDINGS",
  ].join("\n"));

  w("leg-clean.txt", [
    "预审发现（合成评估 · 干净腿输出）",
    "",
    "1. Critical — 02-diff.js:58：解析异常被吞并为通过，应上抛。",
    "2. High — 02-diff.js:80：路径比较缺归一。",
    "3. Medium — 02-diff.js:93：缺一支冲突检查。",
    "4. Medium — 02-diff.js:125：缺一个子项。",
    "5. Medium — 02-diff.js:131：作用域未使用入参。",
    "6. Low — 02-diff.js:14：注释与实现不符。",
    "",
    "PRE-REVIEW: FINDINGS",
  ].join("\n"));

  w("leg-contam.txt", [
    "预审发现：该行包含标记词 OB-79，应触发污染检测失败。",
    "另有 02-diff.js:58 一处引用（路径在 bundle 内）。",
    "",
    "PRE-REVIEW: FINDINGS",
  ].join("\n"));

  w("leg-outside.txt", [
    "1. High — docs/other.js:5：路径不在 bundle 内。",
    "",
    "PRE-REVIEW: FINDINGS",
  ].join("\n"));

  // Word-boundary corpus for contam.js (② fix #2):
  //   mustFail = exact markers (still blocked); mustPass = longer look-alikes (no longer blocked).
  w("boundary-ob79.txt", ["标记词边界：OB-79 应命中。", "", "PRE-REVIEW: FINDINGS"].join("\n"));
  w("boundary-ob790.txt", ["标记词边界：OB-790 是否命中？", "", "PRE-REVIEW: FINDINGS"].join("\n"));
  w("boundary-s8.txt", ["种子边界：S8 应命中。", "", "PRE-REVIEW: FINDINGS"].join("\n"));
  w("boundary-s9.txt", ["种子边界：S9 是否命中？", "", "PRE-REVIEW: FINDINGS"].join("\n"));
  w("boundary-s10.txt", ["种子边界：S10 是否命中？", "", "PRE-REVIEW: FINDINGS"].join("\n"));
  w("boundary-s123.txt", ["种子边界：S123 是否命中？", "", "PRE-REVIEW: FINDINGS"].join("\n"));
  w("boundary-directive.txt", ["切片边界：DIRECTIVE-EVIDENCE-SCOPE 应命中。", "", "PRE-REVIEW: FINDINGS"].join("\n"));
  w("boundary-g8a.txt", ["切片边界：G8-a 应命中。", "", "PRE-REVIEW: FINDINGS"].join("\n"));
  w("boundary-hash025.txt", ["提交边界：025158f 应命中。", "", "PRE-REVIEW: FINDINGS"].join("\n"));
  w("boundary-hash0c1.txt", ["提交边界：0c12abc 应命中。", "", "PRE-REVIEW: FINDINGS"].join("\n"));
  w("boundary-cjk1.txt", ["中文标记：重建归档 应命中。", "", "PRE-REVIEW: FINDINGS"].join("\n"));
  w("boundary-cjk2.txt", ["中文标记：判据作用域排除 应命中。", "", "PRE-REVIEW: FINDINGS"].join("\n"));

  /* ---- citation-grammar corpus (③ scorer fix: general parsing rules) ---- */
  // range coverage: S4 (accept 88..93) is reachable ONLY via a range citation.
  w("leg-range.txt", [
    "1. High — 02-diff.js:84-95：清单有、包内无一支未实现。",
    "",
    "PRE-REVIEW: FINDINGS",
  ].join("\n"));
  // negative control: a range that misses every accept_line must not hit.
  w("leg-range-miss.txt", [
    "1. High — 02-diff.js:84-86：非命中区间。",
    "",
    "PRE-REVIEW: FINDINGS",
  ].join("\n"));
  // comma-separated multi-line numbers; severity must come from the table row.
  w("leg-comma.txt", [
    "| F6 | High | 02-diff.js:69, 80 | 分隔符未归一 |",
    "",
    "PRE-REVIEW: FINDINGS",
  ].join("\n"));
  // slash separator + two citations on one line.
  w("leg-slash.txt", [
    "1. Critical — 02-diff.js:58 / 80：两处。",
    "2. High — 02-diff.js:105；01-clause-cn.txt:10：混排。",
    "",
    "PRE-REVIEW: FINDINGS",
  ].join("\n"));
  // full-width colon between file and line.
  w("leg-fullwidth.txt", [
    "1. Critical — 02-diff.js：58：解析异常被吞并。",
    "",
    "PRE-REVIEW: FINDINGS",
  ].join("\n"));
  // duplicate references of one extra finding (table + body) dedupe to one.
  w("leg-dedup.txt", [
    "| F5 | High | 02-diff.js:75 | 控制文件被算作成员 |",
    "- 位置：02-diff.js:75（同一发现重复引用）。",
    "",
    "PRE-REVIEW: FINDINGS",
  ].join("\n"));
  // abstention via an owning section heading (no per-line marker).
  w("leg-abstain.txt", [
    "## 已核验、不构成发现",
    "- 02-diff.js:191-198（`gx.test.js` 的 `limits()`）：重排加空行。",
    "- 02-diff.js:155-157（`require(\"./gx\")` 等价改写）：解析等价。",
    "",
    "PRE-REVIEW: FINDINGS",
  ].join("\n"));
  // abstention via a per-line marker only (no heading).
  w("leg-abstain-line.txt", [
    "1. 02-diff.js:194：纯重排，无功能变化，不构成发现。",
    "",
    "PRE-REVIEW: FINDINGS",
  ].join("\n"));
  // table-cell severity forms drive the reported severity.
  w("leg-sev-table.txt", [
    "| F1 | Critical | 02-diff.js:56-59, 137 | 失败被吞并 |",
    "| F3 | High | 02-diff.js:84-95 | X4.2 缺失 |",
    "",
    "PRE-REVIEW: FINDINGS",
  ].join("\n"));

  const bomPath = path.join(FIX, "hash-bom.txt");
  const nobomPath = path.join(FIX, "hash-nobom.txt");
  fs.writeFileSync(bomPath, Buffer.concat([Buffer.from([0xef, 0xbb, 0xbf]), Buffer.from("hello\n", "utf8")]));
  fs.writeFileSync(nobomPath, Buffer.from("hello\n", "utf8"));
  return { bomPath, nobomPath };
}

/* ------------------------------------------------------------------ */
/* T-GE-01  key <-> bundle consistency (independent re-computation)   */
/* ------------------------------------------------------------------ */

function tge01() {
  const key = JSON.parse(stripBom(readText(KEY)));
  const seeds = key.seeds || [];
  const lineCounts = {};
  for (const f of ORDER) lineCounts[f] = bundleLines(f).length;

  const anchors = {
    S1: { line: 58, sub: "return true;", file: "02-diff.js" },
    S2: { line: 80, sub: "const hit = manifestNames.find((name) => name === pkgPath);", file: "02-diff.js" },
    S3: { line: 105, sub: "const recorded = cols[0];", file: "02-diff.js" },
    S4: { line: 93, sub: "// Manifest-only entries are not inspected here.", file: "02-diff.js" },
    S5: { line: 125, sub: "// X3.1-4 (rebuild note) is checked by the caller.", file: "02-diff.js" },
    S6: { line: 131, sub: "const files = listPackageMembers(root);", file: "02-diff.js" },
    S7: { line: 14, sub: "// Verifies eight finding classes per package.", file: "02-diff.js" },
  };

  // 1.1 anchors locatable: file/line/substring agree with key AND with physical text.
  for (const id of ["S1", "S2", "S3", "S4", "S5", "S6", "S7"]) {
    const s = seeds.find((x) => x.id === id);
    const a = anchors[id];
    const line = bundleLines(a.file)[a.line - 1];
    const okKey = !!s && s.anchor === a.file + ":" + a.line;
    const okSub = !!s && line.indexOf(s.anchor_substring) !== -1;
    check("T-GE-01.1." + id, id + " anchor " + a.file + ":" + a.line + " locatable (key anchor + substring on physical line)",
      okKey && okSub,
      "key=" + (s && s.anchor) + " sub=" + JSON.stringify(s && s.anchor_substring) + " physical=" + JSON.stringify(line));
  }

  // 1.2 uniqueness of each anchor marker.
  const diff = bundleText("02-diff.js");
  const uniq = [
    ["S1", "catch", 1],
    ["S1", "return true;", 1],
    ["S2", "manifestNames.find(", 1],
    ["S3", "const recorded = cols[0];", 1],
    ["S4", "// Manifest-only entries are not inspected here.", 1],
    ["S5", "// X3.1-4 (rebuild note) is checked by the caller.", 1],
    ["S6", "const files = listPackageMembers(root);", 1],
    ["S7", "// Verifies eight finding classes per package.", 1],
  ];
  for (const [id, needle, expected] of uniq) {
    const c = countOcc(diff, needle);
    check("T-GE-01.2." + id + "." + needle.slice(0, 18), id + " marker occurs exactly " + expected + "x",
      c === expected, "occurrences=" + c);
  }
  // bare cols[0] is NOT unique -- documented, not a failure.
  check("T-GE-01.2.S3.note", "bare `cols[0]` occurrence count (informational)",
    countOcc(diff, "cols[0]") >= 1, "cols[0] occurs " + countOcc(diff, "cols[0]") + "x (54,54,105,108)");

  // 1.3 accept_lines disjoint + in-bounds; S8 region in-bounds and disjoint from seeds.
  const acceptAll = [];
  for (const s of seeds) {
    const lines = Array.isArray(s.accept_lines) ? s.accept_lines : [];
    const inb = lines.every((n) => Number.isInteger(n) && n >= 1 && n <= lineCounts["02-diff.js"]);
    check("T-GE-01.3." + s.id, s.id + " accept_lines in-bounds", inb, JSON.stringify(lines));
    if (s.id !== "S8") acceptAll.push(...lines);
  }
  const dupes = acceptAll.filter((n, i) => acceptAll.indexOf(n) !== i);
  check("T-GE-01.3.disjoint", "S1..S7 accept_lines pairwise disjoint", dupes.length === 0, "dupes=" + JSON.stringify(dupes));
  const s8 = key.s8_region || { lines: [] };
  const s8Lines = bundleLines(s8.file);
  check("T-GE-01.3.s8bounds", "S8 region within file bounds",
    s8.lines.every((n) => n >= 1 && n <= s8Lines.length), JSON.stringify(s8.lines));
  const s8Disjoint = acceptAll.every((n) => s8.lines.indexOf(n) === -1);
  check("T-GE-01.3.s8disjoint", "S8 region disjoint from S1..S7 accept_lines", s8Disjoint, "");

  // 1.4 S8 region is a pure reorder / blank-line change (manual line-by-line).
  const region = s8Lines.slice(s8.lines[0] - 1, s8.lines[s8.lines.length - 1]);
  check("T-GE-01.4.s8shape", "S8 region (191-198) == pure reorder + blank line",
    JSON.stringify(region) === JSON.stringify(S8_EXPECTED_LINES), JSON.stringify(region));
  const changed = region.filter((l) => l.startsWith("+") || l.startsWith("-"));
  const functionalTokens = ["return", "if", "throw", "push", "crypto", "digest", "Number", "Boolean"];
  const functional = changed.some((l) => functionalTokens.some((t) => l.includes(t)));
  check("T-GE-01.4.s8nofunc", "S8 *changed* lines contain no functional-change token", !functional, "");
  const onlyConsts = changed.every((l) => /^[+-]\s*(const (min|max) = \d+;)?\s*$/.test(l));
  check("T-GE-01.4.s8onlyconst", "S8 changed lines are only const-decl reorder / blank", onlyConsts, JSON.stringify(changed));
  const s8seed = seeds.find((x) => x.id === "S8");
  const s8seedFile = String(s8seed.anchor).split(":")[0];
  const s8subOcc = countOcc(bundleText(s8seedFile), s8seed.anchor_substring);
  check("T-GE-01.4.s8anchoruniq", "S8 anchor_substring is unique after ② fix #4 (occurs exactly 1x)",
    s8subOcc === 1, "occurrences=" + s8subOcc + " sub=" + JSON.stringify(s8seed.anchor_substring));
  const s8anchorLine = bundleLines(s8seedFile)[Number(String(s8seed.anchor).split(":")[1]) - 1];
  check("T-GE-01.4.s8anchorloc", "S8 anchor line carries the (now unique) anchor_substring",
    typeof s8anchorLine === "string" && s8anchorLine.indexOf(s8seed.anchor_substring) !== -1,
    JSON.stringify(s8anchorLine));
}

/* ------------------------------------------------------------------ */
/* T-GE-02  bundle leakage self-check                                 */
/* ------------------------------------------------------------------ */

function tge02() {
  const patterns = [
    ["model-name", /GLM|DeepSeek|glm|deepseek/g],
    ["model-family", /GPT|Claude|Qwen|Mythos|GLM-5/i],
    ["real-repo-path", /docs\/validation|tools\/gates|\.github|internal|tmp\/glm-eval/g],
    ["real-slice-id", /DIRECTIVE-[A-Z-]+|OB-\d+|025158f|0c12abc|G8-[abcd]/g],
    ["seed-id", /S[1-8]/g],
  ];
  for (const f of ORDER) {
    const text = bundleText(f);
    for (const [label, re] of patterns) {
      const m = text.match(re) || [];
      check("T-GE-02." + f + "." + label, f + " has 0 matches for " + label + " (bundle leak)",
        m.length === 0, m.length ? "HITS=" + JSON.stringify(m.slice(0, 5)) : "");
    }
  }
  // expected_finding *sentences* must not appear verbatim anywhere in the bundle.
  const key = JSON.parse(stripBom(readText(KEY)));
  const all = ORDER.map(bundleText).join("\n");
  for (const s of key.seeds) {
    check("T-GE-02.diag.exact." + s.id, s.id + " expected_finding sentence absent from bundle",
      all.indexOf(s.expected_finding) === -1, "len=" + s.expected_finding.length);
  }
  // Token-level overlap report. Shared vocabulary is legitimate where the spec *states*
  // the rule (00/01/03); it would only be leakage if an editorial diagnosis appeared in
  // 02-diff.js (the artifact under review) -- except S7, whose defect IS a stale comment.
  const tokens = ["静默", "吞并", "分隔符归一", "恒真", "形同虚设", "历史误报", "fail-open", "全树扫描"];
  const diffText = bundleText("02-diff.js");
  for (const t of tokens) {
    const perFile = ORDER.filter((f) => bundleText(f).indexOf(t) !== -1).join(",") || "(none)";
    const inDiff = diffText.indexOf(t) !== -1;
    note("T-GE-02.token." + t, "token '" + t + "' present in [" + perFile + "]" +
      (inDiff ? " (INCLUDING 02-diff.js)" : "; not in 02-diff.js"));
  }
  // Document necessary (non-leaking) overlap: S7's defect is the stale comment itself.
  check("T-GE-02.note.s7", "S7 expected_finding quotes the defect comment itself (necessary, not leakage)",
    bundleText("02-diff.js").indexOf("eight finding classes") !== -1, "present by design at 02-diff.js:14");
}

/* ------------------------------------------------------------------ */
/* T-GE-03  03-prompt.txt contract completeness + byte-reusability    */
/* ------------------------------------------------------------------ */

function tge03() {
  const p = bundleText("03-prompt.txt");
  const lines = stripBom(p).split(/\r?\n/);
  const closed = ["PRE-REVIEW: PASS", "PRE-REVIEW: PASS WITH FIXES", "PRE-REVIEW: FINDINGS"];
  for (const c of closed) check("T-GE-03.closed." + c, "prompt states conclusion line '" + c + "'", p.indexOf(c) !== -1, "");
  let i = lines.length - 1;
  while (i >= 0 && lines[i].trim() === "") i--;
  check("T-GE-03.lastline", "prompt last non-empty line is in the closed set",
    closed.indexOf(lines[i].trim()) !== -1, JSON.stringify(lines[i]));
  check("T-GE-03.noread", "prompt contains '不得读取任何工作区或外部文件'",
    p.indexOf("不得读取任何工作区或外部文件") !== -1, "");
  check("T-GE-03.msgonly", "prompt contains '仅依据本消息文本'", p.indexOf("仅依据本消息文本") !== -1, "");
  check("T-GE-03.finalmsg", "prompt requires report returned as the final message",
    p.indexOf("报告以最终消息返回") !== -1, "");
  const modelRe = /GLM|DeepSeek|glm|deepseek|GPT|Claude|Qwen|V4|V4\.1|Flash/i;
  check("T-GE-03.nomodel", "prompt has no model-specific token", !modelRe.test(p), (p.match(modelRe) || []).join(","));
  const legRe = /leg-?[AB]|对照腿|实验腿|对照组|实验组/i;
  check("T-GE-03.noleg", "prompt has no leg-specific wording", !legRe.test(p), (p.match(legRe) || []).join(","));
  check("T-GE-03.singlefile", "one shared prompt file => byte-identical reuse for both legs",
    ORDER.indexOf("03-prompt.txt") !== -1, "sha256=" + sha256(Buffer.from(p, "utf8")));
}

/* ------------------------------------------------------------------ */
/* T-GE-04  contract <-> implementation derivability (S4/S5/S6)       */
/* ------------------------------------------------------------------ */

function tge04() {
  const clause = bundleText("01-clause-cn.txt");
  const diff = bundleText("02-diff.js");

  // S5: X3.1 declares four sub-items incl. #4 (rebuild note); runChecks only wires three.
  check("T-GE-04.s5.clause", "clause X3.1 declares four sub-items incl. 子项四",
    clause.indexOf("X3.1") !== -1 && clause.indexOf("四个子项") !== -1 && clause.indexOf("子项四") !== -1, "");
  check("T-GE-04.s5.impl", "runChecks wires only 3 sub-checks (X3.1-4 handed to caller)",
    countOcc(diff, "findings.push(...check") === 3 &&
      diff.indexOf("X3.1-4 (rebuild note) is checked by the caller.") !== -1, "");

  // S4: X4.2 declares manifest-only conflict; implementation has no such branch.
  check("T-GE-04.s4.clause", "clause X4.2 declares manifest-only conflict",
    clause.indexOf("X4.2") !== -1 && clause.indexOf("清单有成员而包内无该文件") !== -1, "");
  check("T-GE-04.s4.impl", "checkMembers explicitly does NOT inspect manifest-only entries",
    diff.indexOf("Manifest-only entries are not inspected here.") !== -1, "");

  // S6: X5.2/X5.3 require change-set scope; collectScope ignores changedPaths + full scan.
  check("T-GE-04.s6.clause", "clause X5.2/X5.3 require change-set-only scope, no full scan",
    clause.indexOf("X5.2") !== -1 && clause.indexOf("变更集之外的历史文件不得参与") !== -1 &&
      clause.indexOf("不得触发任何全量扫描") !== -1, "");
  const csIdx = diff.indexOf("function collectScope");
  const csEnd = csIdx >= 0 ? diff.indexOf("\n+}", csIdx) : -1;
  const body = csIdx >= 0 && csEnd > csIdx ? diff.slice(csIdx, csEnd) : "";
  const nl = body.indexOf("\n");
  const bodyInner = nl >= 0 ? body.slice(nl) : ""; // exclude the signature line
  check("T-GE-04.s6.impl", "collectScope body ignores changedPaths and calls listPackageMembers(root)",
    bodyInner.indexOf("listPackageMembers(root)") !== -1 && bodyInner.indexOf("changedPaths") === -1,
    JSON.stringify(bodyInner));
}

/* ------------------------------------------------------------------ */
/* T-GE-05  score.js positive/negative fixtures + mutation assertion  */
/* ------------------------------------------------------------------ */

function tge05(fx) {
  const key = JSON.parse(stripBom(readText(KEY)));

  const perfect = scoreFixture(path.join(FIX, "leg-perfect.txt"));
  check("T-GE-05.perfect.recall_all", "perfect: recall_all=1", perfect.recall_all === 1, JSON.stringify(perfect.recall_all));
  check("T-GE-05.perfect.recall_high", "perfect: recall_high=1", perfect.recall_high === 1, "");
  check("T-GE-05.perfect.sev_exact", "perfect: sev_exact=1", perfect.sev_exact === 1, "");
  check("T-GE-05.perfect.conclusion", "perfect: conclusion_line=1", perfect.conclusion_line === 1, "");
  check("T-GE-05.perfect.Q", "perfect: Q=1.0", perfect.Q === 1, JSON.stringify(perfect.Q));
  check("T-GE-05.perfect.fp", "perfect: FP=0 extra_true=0", perfect.FP === 0 && perfect.extra_true === 0, "");

  const empty = scoreFixture(path.join(FIX, "leg-empty.txt"));
  check("T-GE-05.empty", "empty: all recall=0, conclusion=1, Q=0.2",
    empty.recall_all === 0 && empty.recall_high === 0 && empty.sev_exact === 0 &&
      empty.conclusion_line === 1 && empty.Q === 0.2, JSON.stringify({ Q: empty.Q, c: empty.conclusion_line }));

  const noconcl = scoreFixture(path.join(FIX, "leg-noconcl.txt"));
  check("T-GE-05.noconcl", "no conclusion line => conclusion_line=0", noconcl.conclusion_line === 0,
    JSON.stringify(noconcl.last_line));

  const tws = scoreFixture(path.join(FIX, "leg-trailing-ws.txt"));
  check("T-GE-05.trailingws", "trailing whitespace on conclusion line tolerated",
    tws.conclusion_line === 1, JSON.stringify(tws.last_line));

  const se = scoreFixture(path.join(FIX, "leg-sev-exact.txt"));
  const so1 = scoreFixture(path.join(FIX, "leg-sev-off1.txt"));
  const so2 = scoreFixture(path.join(FIX, "leg-sev-off2.txt"));
  check("T-GE-05.sev.exact", "severity calibration: exact => 1.0", se.per_seed.S1.sev_score === 1, "");
  check("T-GE-05.sev.off1", "severity calibration: off-by-1 => 0.5", so1.per_seed.S1.sev_score === 0.5, "");
  check("T-GE-05.sev.off2", "severity calibration: off-by-2+ => 0", so2.per_seed.S1.sev_score === 0, "");

  const fps = scoreFixture(path.join(FIX, "leg-fps.txt"));
  check("T-GE-05.fp.count", "FP tri-state: 3 false positives", fps.FP === 3, JSON.stringify(fps.fp_citations.map((f) => f.reason)));
  const reasons = fps.fp_citations.map((f) => f.reason).sort();
  check("T-GE-05.fp.reasons", "FP reasons = {line-not-exist, path-not-in-bundle, s8-no-defect-region}",
    JSON.stringify(reasons) === JSON.stringify(["line-not-exist", "path-not-in-bundle", "s8-no-defect-region"]),
    JSON.stringify(reasons));
  check("T-GE-05.fp.s8flag", "S8 misreport flips s8_ok=false (zero-finding assertion live)",
    fps.s8_ok === false && fps.s8_fp === 1, JSON.stringify({ s8_ok: fps.s8_ok, s8_fp: fps.s8_fp }));

  const extra = scoreFixture(path.join(FIX, "leg-extra.txt"));
  check("T-GE-05.extra", "extra_true single: in-bundle non-seed citation => extra_true=1, FP=0",
    extra.extra_true === 1 && extra.FP === 0, JSON.stringify({ e: extra.extra_true, fp: extra.FP }));

  const q = scoreFixture(path.join(FIX, "leg-qcheck.txt"));
  // hits: S1(Critical exact=1), S2(High but reported Medium -> 0.5), S3(High exact=1).
  // recall_high uses S1,S2,S3 -> 3/3 = 1;  recall_all 3/7;  sev_exact 2.5/3;  conclusion 1.
  const handQ = 0.4 * 1 + 0.2 * (3 / 7) + 0.2 * (2.5 / 3) + 0.2 * 1;
  const handQR = Math.round(handQ * 10000) / 10000;
  check("T-GE-05.Q.hand", "Q hand-computed cross-check (" + handQR + ")",
    Math.abs(q.Q - handQR) < 1e-9, "score=" + q.Q + " hand=" + handQR);
  check("T-GE-05.Q.parts", "Q partial components (recall_all 3/7, recall_high 1, sev_exact 2.5/3)",
    q.recall_all === Math.round((3 / 7) * 10000) / 10000 && q.recall_high === 1 &&
      q.sev_exact === Math.round((2.5 / 3) * 10000) / 10000,
    JSON.stringify({ ra: q.recall_all, rh: q.recall_high, se: q.sev_exact, c: q.conclusion_line, Q: q.Q }));

  // Guard: non-anchor accept_lines must also match (not just the single anchor line).
  const acc = scoreFixture(path.join(FIX, "leg-acceptlines.txt"));
  check("T-GE-05.acceptlines",
    "non-anchor accept_lines matched (S1:137,S2:81,S3:107,S4:90,S5:123,S6:130; S7 missed)",
    acc.per_seed.S1.hit && acc.per_seed.S2.hit && acc.per_seed.S3.hit && acc.per_seed.S4.hit &&
      acc.per_seed.S5.hit && acc.per_seed.S6.hit && !acc.per_seed.S7.hit &&
      acc.recall_all === Math.round((6 / 7) * 10000) / 10000 && acc.recall_high === 1 && acc.sev_exact === 1,
    JSON.stringify({ ra: acc.recall_all, rh: acc.recall_high, se: acc.sev_exact, s7: acc.per_seed.S7.hit }));

  // Guard: Chinese severity forms declared in SEV_ORDER.
  const cn = scoreFixture(path.join(FIX, "leg-sev-cn.txt"));
  check("T-GE-05.sev.cn", "Chinese severity forms (严重度：严重 / 严重度：高) parsed",
    cn.per_seed.S1.sev_score === 1 && cn.per_seed.S2.sev_score === 1 && cn.sev_exact === 1,
    JSON.stringify({ s1: cn.per_seed.S1.sev_score, s2: cn.per_seed.S2.sev_score, se: cn.sev_exact }));

  // ---- mutation assertion (temporary, restored byte-identically) ----
  mutationAssertion(path.join(FIX, "leg-fps.txt"));
}

function mutationAssertion(fpsFixture) {
  const NEC = "if (base === basename(s8.file) && (s8.lines || []).some((n) => n >= c.start && n <= c.end)) {";
  const orig = fs.readFileSync(SCORE);
  const origSha = sha256(orig);
  const backupPath = path.join(SELFTEST, "score.js.bak");
  fs.writeFileSync(backupPath, orig); // byte backup for manual recovery
  let mutatedSha = "";
  try {
    const text = orig.toString("utf8");
    check("T-GE-05.mut.needle", "mutation needle present exactly once in score.js",
      countOcc(text, NEC) === 1, "occurrences=" + countOcc(text, NEC));
    const mutated = text.replace(NEC, "if (false) {");
    fs.writeFileSync(SCORE, Buffer.from(mutated, "utf8"));
    mutatedSha = sha256(fs.readFileSync(SCORE));
    const afterMut = scoreFixture(fpsFixture);
    // Under mutation, the S8 citation is no longer classified as FP.
    check("T-GE-05.mut.red", "mutation short-circuit makes the fixture turn RED (FP != 3)",
      afterMut.FP !== 3, "FP=" + afterMut.FP + " s8_ok=" + afterMut.s8_ok);
  } finally {
    fs.writeFileSync(SCORE, orig); // restore original bytes
  }
  const afterSha = sha256(fs.readFileSync(SCORE));
  const restored = afterSha === origSha;
  if (restored) fs.unlinkSync(backupPath); // keep the tree tidy once restore is proven
  check("T-GE-05.restore.sha", "score.js restored byte-identically (sha256 before == after)",
    restored, "before=" + origSha + " after=" + afterSha);
  const green = scoreFixture(fpsFixture);
  check("T-GE-05.restore.green", "fixture GREEN again after restore (FP==3)", green.FP === 3, "FP=" + green.FP);
  console.log("MUTATION before_sha256=" + origSha);
  console.log("MUTATION mutated_sha256=" + mutatedSha);
  console.log("MUTATION after_sha256=" + afterSha);
  console.log("MUTATION byte-identical=" + (afterSha === origSha));
}

/* ------------------------------------------------------------------ */
/* T-GE-06  hash.js correctness (+ BOM stripping, input_hash)         */
/* ------------------------------------------------------------------ */

function tge06(fx) {
  const r1 = runCli(HASH, [fx.bomPath]);
  const r2 = runCli(HASH, [fx.nobomPath]);
  const parse = (out) => {
    const o = {};
    out.split(/\r?\n/).forEach((l) => { const [k, v] = l.split("\t"); if (k && v) o[k] = v; });
    return o;
  };
  const o1 = parse(r1.stdout), o2 = parse(r2.stdout);
  check("T-GE-06.export", "hash.js prints output_hash for a file argument", !!o1.output_hash && !!o2.output_hash,
    "bom=" + o1.output_hash + " nobom=" + o2.output_hash);
  check("T-GE-06.bom.eq", "BOM and non-BOM same-content inputs => equal output_hash",
    o1.output_hash === o2.output_hash, "");
  const expect = sha256(Buffer.from("hello\n", "utf8"));
  check("T-GE-06.bom.value", "output_hash == sha256(stripped content)", o1.output_hash === expect,
    "got=" + o1.output_hash + " want=" + expect);

  // Independent recompute of input_hash (raw byte concat in fixed order).
  const bufs = ORDER.map((f) => fs.readFileSync(path.join(BUNDLE, f)));
  const mine = sha256(Buffer.concat(bufs));
  check("T-GE-06.input.matchesHash", "hash.js input_hash == independent byte-concat sha256",
    o1.input_hash === mine, "hash.js=" + o1.input_hash + " mine=" + mine);
  check("T-GE-06.input.matchesReported", "independent input_hash == ② reported value",
    mine === REPORTED_INPUT_HASH, "mine=" + mine + " reported=" + REPORTED_INPUT_HASH);
  console.log("INPUT_HASH mine=" + mine);
  console.log("INPUT_HASH hash.js=" + o1.input_hash);
  console.log("INPUT_HASH reported=" + REPORTED_INPUT_HASH);
  const perFile = r1.stdout.split(/\r?\n/).filter((l) => l.includes("\t") && ORDER.some((f) => l.startsWith(f + "\t")));
  console.log("PER_FILE_HASHES\n" + perFile.join("\n"));
}

/* ------------------------------------------------------------------ */
/* T-GE-07  contam.js two-state + word-boundary edges (② fix #2)       */
/* ------------------------------------------------------------------ */

function tge07() {
  const clean = runCli(CONTAM, [path.join(FIX, "leg-clean.txt"), KEY]);
  check("T-GE-07.clean", "clean leg => PASS / exit 0", clean.status === 0 &&
    JSON.parse(clean.stdout).status === "PASS", "exit=" + clean.status);
  const contam = runCli(CONTAM, [path.join(FIX, "leg-contam.txt"), KEY]);
  check("T-GE-07.contam", "marked leg => FAIL / exit 1", contam.status === 1 &&
    JSON.parse(contam.stdout).status === "FAIL", "exit=" + contam.status);
  const outside = runCli(CONTAM, [path.join(FIX, "leg-outside.txt"), KEY]);
  check("T-GE-07.outside", "out-of-bundle citation => FAIL / exit 1", outside.status === 1, "exit=" + outside.status);

  // Exact markers MUST still block (positive controls for the \b fix).
  const mustFail = [
    ["ob79", "boundary-ob79.txt", "OB-79"],
    ["s8", "boundary-s8.txt", "S8"],
    ["directive", "boundary-directive.txt", "DIRECTIVE-EVIDENCE-SCOPE"],
    ["g8a", "boundary-g8a.txt", "G8-a"],
    ["hash025", "boundary-hash025.txt", "025158f"],
    ["hash0c1", "boundary-hash0c1.txt", "0c12abc"],
    ["cjk1", "boundary-cjk1.txt", "重建归档"],
    ["cjk2", "boundary-cjk2.txt", "判据作用域排除"],
  ];
  for (const [id, file, marker] of mustFail) {
    const r = runCli(CONTAM, [path.join(FIX, file), KEY]);
    let j = null;
    try { j = JSON.parse(r.stdout); } catch (e) { /* handled */ }
    check("T-GE-07.marker." + id, "exact marker '" + marker + "' => FAIL / exit 1",
      r.status === 1 && !!j && j.status === "FAIL", "exit=" + r.status);
  }

  // Longer tokens that merely CONTAIN a marker must NOT block (② fix #2 behaviour).
  const mustPass = [
    ["ob790", "boundary-ob790.txt"],
    ["s9", "boundary-s9.txt"],
    ["s10", "boundary-s10.txt"],
    ["s123", "boundary-s123.txt"],
  ];
  for (const [id, file] of mustPass) {
    const r = runCli(CONTAM, [path.join(FIX, file), KEY]);
    check("T-GE-07.boundary." + id, "boundary: " + file + " => PASS / exit 0 (word-boundary fix)",
      r.status === 0, "exit=" + r.status);
  }
}

/* ------------------------------------------------------------------ */
/* T-GE-08  dataset scaffold schema                                   */
/* ------------------------------------------------------------------ */

function tge08() {
  const raw = fs.readFileSync(DATASET);
  check("T-GE-08.bom", "dataset has no BOM", !(raw[0] === 0xef && raw[1] === 0xbb && raw[2] === 0xbf), "");
  let hasCR = false;
  for (const b of raw) if (b === 0x0d) { hasCR = true; break; }
  check("T-GE-08.lf", "dataset uses LF only", !hasCR, "");
  let j = null;
  try { j = JSON.parse(stripBom(raw.toString("utf8"))); } catch (e) { /* handled below */ }
  check("T-GE-08.json", "dataset parses as JSON", !!j, "");
  if (!j) return;

  check("T-GE-08.slice", "slice id matches the registered slice", j.slice === "DIRECTIVE-GLM-EVAL", j.slice);
  // Scaffold-era assertion updated: the dataset is now filled by execution, so calls[]
  // holds the frozen 8-row probe/batch ledger instead of being empty.
  const CALL_ROLES = ["LIMIT-NAME-PROBE-REJECTED", "LIMIT-NAME-PROBE-ACCEPTED", "5-EXPERIMENT",
    "5-CONTROL", "BATCH-a", "BATCH-b", "BATCH-c", "BATCH-d"];
  const calls = Array.isArray(j.calls) ? j.calls : [];
  const preReg = j.pre_registration || {};
  const preHashesPresent =
    typeof preReg.bundle_input_hash === "string" && typeof preReg.key_file_sha256 === "string";
  check("T-GE-08.calls",
    "calls[] has exactly 8 rows: call_no 1..8 and role_position in the frozen probe/batch order",
    calls.length === 8 &&
      calls.every((c, i) => c.call_no === i + 1) &&
      calls.every((c, i) => c.role_position === CALL_ROLES[i]),
    "len=" + calls.length +
      " nos=" + JSON.stringify(calls.map((c) => c.call_no)) +
      " roles=" + JSON.stringify(calls.map((c) => c.role_position)) +
      " pre_registration carries bundle_input_hash+key_file_sha256=" + preHashesPresent);

  const pre = j.pre_registration || {};
  const isHex64 = (v) => typeof v === "string" && /^[0-9a-f]{64}$/.test(v);
  const allPending = pre.bundle_input_hash === "PENDING" && pre.key_file_sha256 === "PENDING" && pre.registered_at === null;
  const allConcrete = isHex64(pre.bundle_input_hash) && isHex64(pre.key_file_sha256) &&
    typeof pre.registered_at === "string" && pre.registered_at.length > 0;
  check("T-GE-08.pre.hashstate",
    "pre_registration hashes are all-PENDING (unregistered) or all-concrete (registered), never mixed",
    allPending || allConcrete,
    JSON.stringify({ h: pre.bundle_input_hash, k: pre.key_file_sha256, at: pre.registered_at }));
  const curInput = sha256(Buffer.concat(ORDER.map((f) => fs.readFileSync(path.join(BUNDLE, f)))));
  check("T-GE-08.pre.inputhash",
    allConcrete ? "registered bundle_input_hash == current hash.js input_hash" : "unregistered: input_hash still PENDING",
    allConcrete ? pre.bundle_input_hash === curInput : pre.bundle_input_hash === "PENDING",
    "pre=" + pre.bundle_input_hash + " cur=" + curInput);
  const preStr = JSON.stringify(pre);
  const leakTokens = ["expected_finding", "severity", "accept_lines", "anchor", "fail-open", "静默", "Critical", "Medium", "High"];
  const leaked = leakTokens.filter((t) => preStr.indexOf(t) !== -1);
  check("T-GE-08.pre.nokey", "pre_registration leaks no key content (expected_finding/severity/anchor)",
    leaked.length === 0, "leaked=" + JSON.stringify(leaked));
  check("T-GE-08.pre.noseedids", "pre_registration no longer carries seed_count / seed_ids (② fix #1)",
    pre.seed_count === undefined && pre.seed_ids === undefined, "keys=" + JSON.stringify(Object.keys(pre)));

  const cols = j.columns || [];
  const names = cols.map((c) => c.name);
  const required = ["call_no", "role_position", "model_string", "quota_package", "input_summary", "input_hash",
    "output_hash", "conclusion_line", "measured_tokens", "metering_mode", "duration_s", "quality_verdict", "notes"];
  check("T-GE-08.cols.required", "columns contain the 13 expected names", required.every((n) => names.indexOf(n) !== -1),
    JSON.stringify(names));
  const metering = (cols.find((c) => c.name === "metering_mode") || {}).domain;
  check("T-GE-08.cols.metering", "metering_mode domain == metering/README modes (4 values incl. user-provided-cumulative)",
    JSON.stringify((metering || []).slice().sort()) === JSON.stringify(["exact-delta", "interval-sum", "unavailable", "user-provided-cumulative"].sort()),
    JSON.stringify(metering));
  const qv = (cols.find((c) => c.name === "quality_verdict") || {}).domain;
  check("T-GE-08.cols.qv", "quality_verdict domain is a 7-value enum",
    Array.isArray(qv) && qv.length === 7, JSON.stringify(qv));
  const concl = (cols.find((c) => c.name === "conclusion_line") || {}).domain;
  check("T-GE-08.cols.concl", "conclusion_line domain references the closed set",
    typeof concl === "string" && concl.indexOf("conclusion_line_closed_set") !== -1, JSON.stringify(concl));
  const key = JSON.parse(stripBom(readText(KEY)));
  check("T-GE-08.closedset.consistency",
    "pre_registration.conclusion_line_closed_set == key.json closed set == 03-prompt set",
    JSON.stringify(pre.conclusion_line_closed_set) === JSON.stringify(key.conclusion_line_closed_set),
    JSON.stringify(pre.conclusion_line_closed_set));
  note("T-GE-08.cols.spec",
    "EVIDENCE GAP: ① plan §5 is neither on disk nor in the session index; the 13 columns could only be " +
    "checked for internal consistency, NOT cross-checked against the source plan (unverifiable).");
}

/* ------------------------------------------------------------------ */
/* T-GE-09  tmp/ hygiene                                              */
/* ------------------------------------------------------------------ */

function tge09() {
  // Platform relocated OUTSIDE the workspace for isolation; the old "tmp/ contains
  // only glm-eval/ and .gitkeep" assertion no longer applies. Reframed to platform
  // root hygiene: only expected platform/tester entries may sit at the root.
  const top = fs.readdirSync(PLATFORM).sort();
  // Scaffold-era allow-list expanded to the full closed set of current top-level
  // entries: the four new handoff documents plus PROBE-NOTES.txt / plan-v1.md.
  const allowedTop = [
    "BRIEFING.md", "FROZEN-MANIFEST.txt", "GLM-USAGE-READING.txt", "PLAN-PROVENANCE.txt",
    "PROBE-NOTES.txt", "TEST-REPORT-DIRECTIVE-GLM-EVAL.txt", "bundle", "contam.js", "hash.js",
    "key.json", "metering", "plan-v1.md", "score.js", "selftest",
  ].sort();
  const unexpectedTop = top.filter((n) => allowedTop.indexOf(n) === -1);
  check("T-GE-09.root", "platform root contains only expected platform/tester entries",
    unexpectedTop.length === 0, "top=" + JSON.stringify(top) + " unexpected=" + JSON.stringify(unexpectedTop));
  const walk = (d) => fs.readdirSync(d, { withFileTypes: true }).flatMap((e) => {
    const p = path.join(d, e.name);
    return e.isDirectory() ? walk(p) : [p];
  });
  const files = walk(PLATFORM).map((p) => path.relative(PLATFORM, p).replace(/\\/g, "/"));
  const oneOff = files.filter((f) => /\.(ps1|sh|bat|cmd)$/.test(f));
  check("T-GE-09.noOneOff", "no one-off .ps1/.sh/.bat/.cmd scripts under glm-eval/",
    oneOff.length === 0, JSON.stringify(oneOff));
  const caches = files.filter((f) => /(^|\/)(\.cache|node_modules|coverage|__pycache__|\.venv)(\/|$)/.test(f));
  check("T-GE-09.noCache", "no orphan cache directories under glm-eval/", caches.length === 0, JSON.stringify(caches));
  // Closed set (no wildcards/regex): every current file enumerated by name. Adding a
  // new artifact therefore forces a deliberate update of this list.
  const allowed = new Set([
    // platform root documents + platform code
    "BRIEFING.md", "FROZEN-MANIFEST.txt", "GLM-USAGE-READING.txt", "PLAN-PROVENANCE.txt",
    "PROBE-NOTES.txt", "TEST-REPORT-DIRECTIVE-GLM-EVAL.txt", "plan-v1.md",
    "contam.js", "hash.js", "key.json", "score.js",
    // frozen bundle
    "bundle/00-plan.md", "bundle/01-clause-cn.txt", "bundle/02-diff.js", "bundle/03-prompt.txt",
    // metering record
    "metering/README.md",
    // tester-owned selftest harness
    "selftest/fixture-contam-fail.txt", "selftest/fixture-score.txt",
    "selftest/run-selftest.js", "selftest/selftest-report.txt",
    // tester-owned fixture corpus (regenerated each run)
    "selftest/fixtures/boundary-cjk1.txt", "selftest/fixtures/boundary-cjk2.txt",
    "selftest/fixtures/boundary-directive.txt", "selftest/fixtures/boundary-g8a.txt",
    "selftest/fixtures/boundary-hash025.txt", "selftest/fixtures/boundary-hash0c1.txt",
    "selftest/fixtures/boundary-ob79.txt", "selftest/fixtures/boundary-ob790.txt",
    "selftest/fixtures/boundary-s10.txt", "selftest/fixtures/boundary-s123.txt",
    "selftest/fixtures/boundary-s8.txt", "selftest/fixtures/boundary-s9.txt",
    "selftest/fixtures/hash-bom.txt", "selftest/fixtures/hash-nobom.txt",
    "selftest/fixtures/leg-abstain-line.txt", "selftest/fixtures/leg-abstain.txt",
    "selftest/fixtures/leg-acceptlines.txt", "selftest/fixtures/leg-clean.txt",
    "selftest/fixtures/leg-comma.txt", "selftest/fixtures/leg-contam.txt",
    "selftest/fixtures/leg-dedup.txt", "selftest/fixtures/leg-empty.txt",
    "selftest/fixtures/leg-extra.txt", "selftest/fixtures/leg-fps.txt",
    "selftest/fixtures/leg-fullwidth.txt", "selftest/fixtures/leg-noconcl.txt",
    "selftest/fixtures/leg-outside.txt", "selftest/fixtures/leg-perfect.txt",
    "selftest/fixtures/leg-qcheck.txt", "selftest/fixtures/leg-range-miss.txt",
    "selftest/fixtures/leg-range.txt", "selftest/fixtures/leg-sev-cn.txt",
    "selftest/fixtures/leg-sev-exact.txt", "selftest/fixtures/leg-sev-off1.txt",
    "selftest/fixtures/leg-sev-off2.txt", "selftest/fixtures/leg-sev-table.txt",
    "selftest/fixtures/leg-slash.txt", "selftest/fixtures/leg-trailing-ws.txt",
  ]);
  const unexpected = files.filter((f) => !allowed.has(f));
  check("T-GE-09.known", "all files under glm-eval/ are expected platform or tester artifacts",
    unexpected.length === 0, JSON.stringify(unexpected));
  console.log("GLM_EVAL_FILES\n" + files.join("\n"));
}

/* ------------------------------------------------------------------ */
/* T-GE-10  dataset pre_registration zero-leak guard (② fix #1)        */
/* ------------------------------------------------------------------ */

function tge10() {
  const j = JSON.parse(stripBom(fs.readFileSync(DATASET, "utf8")));
  const key = JSON.parse(stripBom(readText(KEY)));
  const pre = j.pre_registration || {};
  const preStr = JSON.stringify(pre);
  const preLower = preStr.toLowerCase();

  // (a) cardinality words: neither 'seed' nor 'count' may appear.
  const wordHits = ["seed", "count"].filter((t) => preLower.indexOf(t) !== -1);
  check("T-GE-10.token.seedcount",
    "pre_registration contains neither 'seed' nor 'count' (defect cardinality hidden)",
    wordHits.length === 0, "hits=" + JSON.stringify(wordHits));

  // (b) seed ids S1..S8 must not appear.
  const idHits = key.seeds.filter((s) => preStr.indexOf(s.id) !== -1).map((s) => s.id);
  check("T-GE-10.token.ids", "pre_registration contains no seed id (S1..S8)",
    idHits.length === 0, "hits=" + JSON.stringify(idHits));

  // (c) expected_finding sentences must not appear verbatim.
  const sentenceHits = key.seeds.filter((s) => preStr.indexOf(s.expected_finding) !== -1).map((s) => s.id);
  check("T-GE-10.phrase.full", "no expected_finding sentence appears in pre_registration",
    sentenceHits.length === 0, "hits=" + JSON.stringify(sentenceHits));

  // (d) fragment-level scan: no distinctive CJK clause of any expected_finding may leak.
  const CJK = /[\u4e00-\u9fff]/;
  const fragLeaks = [];
  for (const s of key.seeds) {
    const frags = String(s.expected_finding)
      .split(/[\s，。；、（）()「」【】,;\/]+/)
      .filter((f) => f.length >= 6 && CJK.test(f));
    for (const f of frags) if (preStr.indexOf(f) !== -1) fragLeaks.push(s.id + ":" + f);
  }
  check("T-GE-10.phrase.frag",
    "no >=6-char CJK fragment of any expected_finding appears in pre_registration",
    fragLeaks.length === 0, "hits=" + JSON.stringify(fragLeaks.slice(0, 5)));

  // (e) severity labels and anchors must not appear.
  const sevHits = ["Critical", "High", "Medium", "Low"].filter((t) => preStr.indexOf(t) !== -1);
  check("T-GE-10.severity", "no severity label appears in pre_registration",
    sevHits.length === 0, JSON.stringify(sevHits));
  const anchorHits = key.seeds.filter((s) => preStr.indexOf(String(s.anchor)) !== -1).map((s) => s.id);
  check("T-GE-10.anchor", "no seed anchor (file:line) appears in pre_registration",
    anchorHits.length === 0, "hits=" + JSON.stringify(anchorHits));

  note("T-GE-10.scan",
    "pre_registration scanned as a whole serialized string; keys=" + JSON.stringify(Object.keys(pre)));
}

/* ------------------------------------------------------------------ */
/* T-GE-11  scoring denominator is derived from key.expected_total     */
/* ------------------------------------------------------------------ */

function tge11() {
  const keyRaw = fs.readFileSync(KEY);
  const keySha = sha256(keyRaw);
  const keyObj = JSON.parse(stripBom(keyRaw.toString("utf8")));
  const perfectPath = path.join(FIX, "leg-perfect.txt");

  const baseline = scoreFixture(perfectPath);
  check("T-GE-11.baseline", "baseline key (expected_total=7): perfect recall_all=1",
    baseline.recall_all === 1, "recall_all=" + baseline.recall_all);

  const k8 = JSON.parse(JSON.stringify(keyObj));
  k8.expected_total = 8;
  const k8Path = path.join(FIX, "key-total8.json");
  fs.writeFileSync(k8Path, JSON.stringify(k8, null, 2) + "\n", { encoding: "utf8" });

  const r = runCli(SCORE, [perfectPath, k8Path]);
  let modified = null;
  try { modified = JSON.parse(r.stdout); } catch (e) { /* handled */ }
  check("T-GE-11.derived",
    "denominator follows key.expected_total: expected_total=8 => recall_all=7/8=0.875 (not 7/7)",
    r.status === 0 && !!modified && modified.recall_all === 0.875 && modified.recall_all !== baseline.recall_all,
    "baseline=" + baseline.recall_all + " modified=" + (modified && modified.recall_all));

  const handQ8 = Math.round((0.4 * 1 + 0.2 * (7 / 8) + 0.2 * 1 + 0.2 * 1) * 10000) / 10000;
  check("T-GE-11.Q", "Q recomputed with the derived denominator (" + handQ8 + ")",
    !!modified && Math.abs(modified.Q - handQ8) < 1e-9, "Q=" + (modified && modified.Q) + " want=" + handQ8);

  fs.unlinkSync(k8Path);
  const keyAfter = sha256(fs.readFileSync(KEY));
  check("T-GE-11.keyintact", "key.json byte-identical after the derivation guard",
    keyAfter === keySha, "sha=" + keyAfter);
}

/* ------------------------------------------------------------------ */
/* T-GE-12  mutation assertions for the two NEW guards                 */
/* ------------------------------------------------------------------ */

function mutatePlatform(target, label, mutateFn, redProbe) {
  const orig = fs.readFileSync(target);
  const origSha = sha256(orig);
  const origText = orig.toString("utf8");
  const mutatedText = mutateFn(origText);
  const changed = mutatedText !== origText;
  check("T-GE-12." + label + ".mutneedle", label + ": mutation changes the source (needle present)",
    changed, "");
  let mutatedSha = origSha;
  let red = false;
  let detail = "(skipped: needle missing)";
  try {
    if (changed) {
      fs.writeFileSync(target, Buffer.from(mutatedText, "utf8"));
      mutatedSha = sha256(fs.readFileSync(target));
      const pr = redProbe();
      red = pr.red;
      detail = pr.detail;
    }
  } finally {
    fs.writeFileSync(target, orig); // always restore original bytes
  }
  const afterSha = sha256(fs.readFileSync(target));
  check("T-GE-12." + label + ".red", label + ": guard turns RED under mutation", red, detail);
  check("T-GE-12." + label + ".restore", label + ": restored byte-identically",
    afterSha === origSha, "before=" + origSha + " after=" + afterSha);
  console.log("MUTATION-" + label + " before_sha256=" + origSha);
  console.log("MUTATION-" + label + " mutated_sha256=" + mutatedSha);
  console.log("MUTATION-" + label + " after_sha256=" + afterSha);
  console.log("MUTATION-" + label + " byte-identical=" + (afterSha === origSha));
}

function tge12() {
  // M1: strip word boundaries from contam.js MARKER_SRC => boundary cases must turn RED.
  mutatePlatform(CONTAM, "contam-boundary",
    (text) => text.replace(/\\\\b/g, ""), // matches the two-backslash + b source escapes only
    () => {
      const ob = runCli(CONTAM, [path.join(FIX, "boundary-ob790.txt"), KEY]);
      const s10 = runCli(CONTAM, [path.join(FIX, "boundary-s10.txt"), KEY]);
      const s123 = runCli(CONTAM, [path.join(FIX, "boundary-s123.txt"), KEY]);
      const red = ob.status === 1 && s10.status === 1 && s123.status === 1;
      return { red, detail: "ob790=" + ob.status + " s10=" + s10.status + " s123=" + s123.status };
    });

  // M2: hardcode the denominator back to 7 => derivation guard must turn RED.
  mutatePlatform(SCORE, "score-denominator",
    (text) => text.replace(
      "const total = Number.isFinite(key.expected_total) ? key.expected_total : ordered.length;",
      "const total = 7;"),
    () => {
      const k8Path = path.join(FIX, "key-total8.json");
      const k8 = JSON.parse(stripBom(readText(KEY)));
      k8.expected_total = 8;
      fs.writeFileSync(k8Path, JSON.stringify(k8, null, 2) + "\n", { encoding: "utf8" });
      const r = runCli(SCORE, [path.join(FIX, "leg-perfect.txt"), k8Path]);
      let m = null;
      try { m = JSON.parse(r.stdout); } catch (e) { /* handled */ }
      fs.unlinkSync(k8Path);
      const red = !!m && m.recall_all !== 0.875;
      return { red, detail: "hardcoded-7 recall_all=" + (m && m.recall_all) + " (expected 0.875 under derivation)" };
    });
}

/* ------------------------------------------------------------------ */
/* T-GE-13  citation-grammar + abstention rules (③ scorer fix)         */
/* ------------------------------------------------------------------ */

function tge13() {
  // (a) range covers every line l1..l2 -> S4 reachable ONLY via 88..93.
  const range = scoreFixture(path.join(FIX, "leg-range.txt"));
  check("T-GE-13.range.hit", "range 02-diff.js:84-95 covers S4 accept_lines (hit)",
    range.per_seed.S4.hit === true, JSON.stringify(range.per_seed.S4));
  check("T-GE-13.range.sev", "range citation carries its row severity (High vs Medium => 0.5)",
    range.per_seed.S4.sev_score === 0.5, "sev=" + range.per_seed.S4.sev_score);
  check("T-GE-13.range.recall", "range hit makes recall_all=1/7 (not 0)",
    range.recall_all === Math.round((1 / 7) * 10000) / 10000, "recall_all=" + range.recall_all);

  // (a-negative) a range that misses all accept_lines must NOT hit.
  const miss = scoreFixture(path.join(FIX, "leg-range-miss.txt"));
  check("T-GE-13.range.negative", "range 02-diff.js:84-86 misses S4 (88..93) -> no hit",
    miss.per_seed.S4.hit === false, JSON.stringify(miss.per_seed.S4));

  // (b) comma-separated multi-line numbers, severity from the table row.
  const comma = scoreFixture(path.join(FIX, "leg-comma.txt"));
  check("T-GE-13.comma.hit", "02-diff.js:69, 80 hits S2 at line 80", comma.per_seed.S2.hit === true, "");
  check("T-GE-13.comma.sev", "table-row severity (High) attaches to the S2 hit => sev 1",
    comma.per_seed.S2.reported_severity_rank === 2 && comma.per_seed.S2.sev_score === 1,
    JSON.stringify(comma.per_seed.S2));
  check("T-GE-13.comma.extra", "non-seed line 69 in the same cell is one deduped extra",
    comma.extra_true === 1 && comma.FP === 0, JSON.stringify({ e: comma.extra_true, fp: comma.FP }));

  // (c) slash separator + several citations on one line.
  const slash = scoreFixture(path.join(FIX, "leg-slash.txt"));
  check("T-GE-13.slash", "02-diff.js:58 / 80 hits S1+S2; 02-diff.js:105；01-clause:10 => S3 + 1 extra",
    slash.per_seed.S1.hit && slash.per_seed.S2.hit && slash.per_seed.S3.hit &&
      slash.extra_true === 1 && slash.FP === 0,
    JSON.stringify({ s1: slash.per_seed.S1.hit, s2: slash.per_seed.S2.hit, s3: slash.per_seed.S3.hit, e: slash.extra_true }));

  // (d) full-width colon between file and line.
  const fw = scoreFixture(path.join(FIX, "leg-fullwidth.txt"));
  check("T-GE-13.fullwidth", "02-diff.js：58 (full-width colon) hits S1 with Critical severity",
    fw.per_seed.S1.hit === true && fw.per_seed.S1.sev_score === 1, JSON.stringify(fw.per_seed.S1));

  // (e) duplicate references of one finding dedupe to a single extra.
  const dedup = scoreFixture(path.join(FIX, "leg-dedup.txt"));
  check("T-GE-13.dedup", "same extra cited in table + body => extra_true=1 (deduped by file,start)",
    dedup.extra_true === 1, "extra_true=" + dedup.extra_true);

  // (f) abstention via an owning section heading.
  const absSec = scoreFixture(path.join(FIX, "leg-abstain.txt"));
  check("T-GE-13.abstain.section", "section-marked abstentions: 2 abstentions, FP=0, s8_ok=true",
    absSec.abstentions === 2 && absSec.FP === 0 && absSec.s8_ok === true,
    JSON.stringify({ a: absSec.abstentions, fp: absSec.FP, s8: absSec.s8_ok }));
  const spans = absSec.abstention_citations.map((a) => a.start + "-" + a.end);
  check("T-GE-13.abstain.spans", "abstention spans are 191-198 and 155-157",
    JSON.stringify(spans) === JSON.stringify(["191-198", "155-157"]), JSON.stringify(spans));

  // (g) abstention via a per-line marker only.
  const absLine = scoreFixture(path.join(FIX, "leg-abstain-line.txt"));
  check("T-GE-13.abstain.line", "per-line marker: 194 becomes abstention, not S8 FP",
    absLine.abstentions === 1 && absLine.FP === 0 && absLine.s8_ok === true,
    JSON.stringify({ a: absLine.abstentions, fp: absLine.FP }));

  // (h) table-cell severity forms drive the reported severity.
  const sev = scoreFixture(path.join(FIX, "leg-sev-table.txt"));
  check("T-GE-13.sev.table", "| F1 | Critical | / | F3 | High | => S1 sev 1, S4 hit sev 0.5",
    sev.per_seed.S1.sev_score === 1 && sev.per_seed.S4.hit === true && sev.per_seed.S4.sev_score === 0.5,
    JSON.stringify({ s1: sev.per_seed.S1.sev_score, s4: sev.per_seed.S4 }));
}

/* ------------------------------------------------------------------ */
/* T-GE-14  mutation assertions for the new parser rules               */
/* ------------------------------------------------------------------ */

function mutationGuard(label, needle, replacement, fixtureName, redFn) {
  const orig = fs.readFileSync(SCORE);
  const origSha = sha256(orig);
  const text = orig.toString("utf8");
  const occ = countOcc(text, needle);
  check("T-GE-14." + label + ".needle", label + ": mutation needle occurs exactly once",
    occ === 1, "occ=" + occ);
  let mutatedSha = origSha;
  let red = false;
  let detail = "(skipped: needle missing)";
  try {
    if (occ === 1) {
      const mutated = text.replace(needle, replacement);
      mutatedSha = sha256(Buffer.from(mutated, "utf8"));
      fs.writeFileSync(SCORE, Buffer.from(mutated, "utf8"));
      try {
        const m = scoreFixture(path.join(FIX, fixtureName));
        const pr = redFn(m);
        red = pr.red;
        detail = pr.detail;
      } catch (e) {
        red = false;
        detail = "score.js failed under mutation: " + e.message;
      }
    }
  } finally {
    fs.writeFileSync(SCORE, orig);
  }
  const afterSha = sha256(fs.readFileSync(SCORE));
  check("T-GE-14." + label + ".red", label + ": assertion turns RED under mutation", red, detail);
  check("T-GE-14." + label + ".restore", label + ": score.js restored byte-identically",
    afterSha === origSha, "before=" + origSha + " after=" + afterSha);
  console.log("MUTATION-" + label + " before_sha256=" + origSha);
  console.log("MUTATION-" + label + " mutated_sha256=" + mutatedSha);
  console.log("MUTATION-" + label + " after_sha256=" + afterSha);
  console.log("MUTATION-" + label + " byte-identical=" + (afterSha === origSha));
}

function tge14() {
  // Green baselines.
  const greenRange = scoreFixture(path.join(FIX, "leg-range.txt"));
  check("T-GE-14.green.range", "baseline: range covers S4 (hit=true, sev=0.5)",
    greenRange.per_seed.S4.hit === true && greenRange.per_seed.S4.sev_score === 0.5, "");
  const greenAbstain = scoreFixture(path.join(FIX, "leg-abstain.txt"));
  check("T-GE-14.green.abstain", "baseline: section-marked abstention (2, FP=0)",
    greenAbstain.abstentions === 2 && greenAbstain.FP === 0,
    JSON.stringify({ a: greenAbstain.abstentions, fp: greenAbstain.FP }));

  // M1: short-circuit range parsing => S4 hit must disappear.
  mutationGuard("range-shortcircuit",
    "end = Number(s.slice(m, p));", "end = start; /*MUT-range*/",
    "leg-range.txt",
    (m) => ({ red: m.per_seed.S4.hit === false, detail: "mutated S4.hit=" + m.per_seed.S4.hit + " (expected false)" }));

  // M2: ignore the abstention SECTION flag => the two abstentions become S8 FPs.
  mutationGuard("abstention-section",
    "const lineAbstain = inAbstentionSection ||", "const lineAbstain = false && inAbstentionSection ||",
    "leg-abstain.txt",
    (m) => ({ red: m.abstentions === 0 && m.s8_ok === false, detail: "mutated abstentions=" + m.abstentions + " s8_ok=" + m.s8_ok }));

  // Both fixtures are GREEN again after each restore.
  const backRange = scoreFixture(path.join(FIX, "leg-range.txt"));
  check("T-GE-14.restore.green.range", "range fixture GREEN again after restore",
    backRange.per_seed.S4.hit === true, "");
  const backAbstain = scoreFixture(path.join(FIX, "leg-abstain.txt"));
  check("T-GE-14.restore.green.abstain", "abstention fixture GREEN again after restore",
    backAbstain.abstentions === 2, "");
}

/* ------------------------------------------------------------------ */

function main() {
  const fx = buildFixtures();
  tge01();
  tge02();
  tge03();
  tge04();
  tge05(fx);
  tge06(fx);
  tge07();
  tge08();
  tge09();
  tge10();
  tge11();
  tge12();
  tge13();
  tge14();

  let pass = 0, fail = 0;
  for (const r of RESULTS) {
    const tag = r.ok ? "PASS" : "FAIL";
    if (r.ok) pass++; else fail++;
    console.log("[" + tag + "] " + r.id + " :: " + r.name + (r.detail ? "  | " + r.detail : ""));
  }
  for (const n of NOTES) console.log("[NOTE] " + n.id + " :: " + n.text);
  for (const f of FINDINGS) console.log("[FINDING] " + f.id + " :: " + f.text);
  console.log("SELFTEST TOTAL=" + RESULTS.length + " PASS=" + pass + " FAIL=" + fail +
    " FINDINGS=" + FINDINGS.length + " NOTES=" + NOTES.length);
  const reportPath = path.join(SELFTEST, "selftest-report.txt");
  fs.writeFileSync(reportPath, OUT_LINES.join("\n") + "\n", { encoding: "utf8" }); // no BOM, LF
  _origLog("report written: " + reportPath);
  process.exit(fail === 0 ? 0 : 1);
}

main();
