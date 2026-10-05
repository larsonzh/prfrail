#!/usr/bin/env node
"use strict";

// score.js -- score ONE leg output against key.json (no third-party deps).
//
// Usage: node score.js <leg-output-file> [key.json]
//
// Output (JSON): the five metrics recall_all / recall_high / sev_exact /
// conclusion_line / Q, plus FP and extra_true classification and a per-seed
// hit table.
//
// ---------------------------------------------------------------------------
// Citation grammar (general, output-agnostic -- no seed or leg special-casing)
// ---------------------------------------------------------------------------
// A citation is "<file>:<spec>[<sep><spec> ...]" where
//   * <file> ends in .md / .txt / .js (the only bundle file kinds), and
//   * <spec> is either "<line>" or "<l1>-<l2>", the latter COVERING l1..l2.
// Between/inside specs we accept: comma , / CJK comma ， / CJK enumeration 、 /
// slash / / full-width slash ／ / semicolon ; / CJK semicolon ； / whitespace,
// and an optional colon (: or full-width ：) before each spec. A full-width
// colon is also accepted between the file name and the first spec. One line may
// carry several citations; each citation is classified independently.
//
// A reference = one <spec> bound to its <file>. Ranges match by coverage (any
// seed accept_line inside [start,end] is a hit). Document order decides which
// reference supplies the reported severity for a seed (first hit wins).
//
// ---------------------------------------------------------------------------
// Classification (one reference at a time, in document order)
// ---------------------------------------------------------------------------
//   * ABSTENTION -> neither hit nor FP; listed in abstentions[]. A reference is
//     an abstention when its line -- or the markdown section heading that owns
//     it -- explicitly declares it is not a finding (markers: 不构成发现 /
//     无功能变化 / 已核验 / 不得报出 / 非发现 / 不属发现 / 弃权 / 不予报出).
//   * file outside bundle            -> FP  (path-not-in-bundle)
//   * line/range outside file bounds -> FP  (line-not-exist)
//   * any covered line == a seed accept_line -> HIT (first hit per seed wins)
//   * any covered line in the S8 no-defect region -> FP (s8-no-defect-region)
//   * otherwise                      -> extra_true, deduped by (file, start)
//
// Nothing here is keyed to specific seed ids or to any one leg output:
// denominators, severity ranks and match scores all come from key.json.

const fs = require("fs");
const path = require("path");

const HERE = __dirname;
const DEFAULT_KEY = path.join(HERE, "key.json");

// Reported-severity probes. Only the free-text patterns live in code; the RANK of
// each level is read from key.json ("severity_rank"), so the scale is data-driven.
// Each level accepts: a bare English token anywhere, a labelled CJK form
// (严重度：高), a bracketed form ([Critical] / 【Critical】) and a table-cell form
// (| F1 | Critical | ...).
const SEV_LEVELS = [
  { level: "Critical", re: /critical|致命|严重度\s*[:：]\s*(?:critical|严重|致命)|[\[【]\s*(?:critical|严重|致命)\s*[\]】]|\|\s*(?:critical|严重|致命)\s*\|/i },
  { level: "High", re: /high|严重度\s*[:：]\s*(?:high|高)|[\[【]\s*(?:high|高)\s*[\]】]|\|\s*(?:high|高)\s*\|/i },
  { level: "Medium", re: /medium|严重度\s*[:：]\s*(?:medium|中)|[\[【]\s*(?:medium|中)\s*[\]】]|\|\s*(?:medium|中)\s*\|/i },
  { level: "Low", re: /low|严重度\s*[:：]\s*(?:low|低)|[\[【]\s*(?:low|低)\s*[\]】]|\|\s*(?:low|低)\s*\|/i },
];

// File token: a path-ish token ending in a bundle file extension.
const FILE_SRC = "([A-Za-z0-9_.\\/\\\\-]+\\.(?:md|txt|js))";
// Characters that may separate two specs (a colon may precede each spec).
const SPEC_SEP = " \t,，、/／;；:：";
// Phrases that mark a reference as an explicit non-finding / abstention.
const ABSTENTION_MARKERS = [
  "不构成发现", "不属发现", "非发现", "无功能变化",
  "已核验", "不得报出", "弃权", "不予报出",
];

function basename(p) {
  return p.split(/[\\/]/).pop();
}

function stripBom(text) {
  return text.replace(/^\uFEFF/, "");
}

function rankOf(text, severityRank) {
  for (const s of SEV_LEVELS) {
    if (s.re.test(text)) {
      const r = (severityRank || {})[s.level];
      return r === undefined ? null : r;
    }
  }
  return null;
}

// Parse the "<spec>[<sep><spec>...]" tail that follows one file token.
// Returns an array of {start,end} (a single line is start === end).
function parseLineSpecs(s) {
  const specs = [];
  let i = 0;
  const isSep = (ch) => SPEC_SEP.indexOf(ch) !== -1;
  const isDigit = (ch) => ch >= "0" && ch <= "9";
  while (i < s.length) {
    while (i < s.length && isSep(s[i])) i += 1;
    if (i >= s.length || !isDigit(s[i])) break;
    let j = i;
    while (j < s.length && isDigit(s[j])) j += 1;
    let start = Number(s.slice(i, j));
    i = j;
    let end = start;
    let k = i;
    while (k < s.length && (s[k] === " " || s[k] === "\t")) k += 1;
    if (k < s.length && (s[k] === "-" || s[k] === "\u2013" || s[k] === "~")) {
      let m = k + 1;
      while (m < s.length && (s[m] === " " || s[m] === "\t")) m += 1;
      if (m < s.length && isDigit(s[m])) {
        let p = m;
        while (p < s.length && isDigit(s[p])) p += 1;
        end = Number(s.slice(m, p));
        i = p;
      }
    }
    if (end < start) {
      const t = start;
      start = end;
      end = t;
    }
    specs.push({ start, end });
  }
  return specs;
}

function parseCitations(text, severityRank) {
  const out = [];
  const lines = text.split(/\r?\n/);
  let inAbstentionSection = false;
  for (const line of lines) {
    // A heading resets the section flag; an abstention heading opens a section.
    if (/^\s{0,3}#{1,6}\s/.test(line)) {
      inAbstentionSection = ABSTENTION_MARKERS.some((m) => line.indexOf(m) !== -1);
    }
    const lineAbstain = inAbstentionSection ||
      ABSTENTION_MARKERS.some((m) => line.indexOf(m) !== -1);

    // Locate every file token; each one owns the text segment up to the next.
    const re = new RegExp(FILE_SRC, "g");
    const tokens = [];
    let m;
    while ((m = re.exec(line)) !== null) {
      tokens.push({ file: m[1], index: m.index, length: m[0].length });
    }
    for (let t = 0; t < tokens.length; t += 1) {
      const cur = tokens[t];
      const specStart = cur.index + cur.length;
      const specEnd = t + 1 < tokens.length ? tokens[t + 1].index : line.length;
      const specs = parseLineSpecs(line.slice(specStart, specEnd));
      for (const sp of specs) {
        out.push({
          file: cur.file,
          start: sp.start,
          end: sp.end,
          severity_rank: rankOf(line, severityRank),
          abstention: lineAbstain,
          text: line.trim(),
        });
      }
    }
  }
  return out;
}

function acceptLines(seed) {
  if (Array.isArray(seed.accept_lines) && seed.accept_lines.length > 0) {
    return seed.accept_lines;
  }
  const m = String(seed.anchor || "").match(/:(\d+)$/);
  return m ? [Number(m[1])] : [];
}

function fileOf(anchor) {
  return anchor ? String(anchor).split(":")[0] : "";
}

function lastNonEmptyLine(text) {
  const lines = stripBom(text).split(/\r?\n/);
  let i = lines.length - 1;
  while (i >= 0 && lines[i].trim() === "") i -= 1;
  return i >= 0 ? lines[i].trim() : "";
}

function round(n) {
  return Math.round(n * 10000) / 10000;
}

function score(outPath, keyPath) {
  const key = JSON.parse(stripBom(fs.readFileSync(keyPath, "utf8")));
  const bundleFiles = new Set(key.bundle_files.map(basename));
  const seeds = key.seeds || [];
  const s8 = key.s8_region || { file: "", lines: [] };

  // Scoring scales are read from key.json -- no seed ids, denominators or severity
  // weights are hardcoded. `scored_ids` is the ordered set of defect seeds counted
  // toward recall; the no-defect region (s8_region) is scored separately as FP.
  const ordered = Array.isArray(key.scored_ids) ? key.scored_ids : seeds.map((s) => s.id);
  const scoredSet = new Set(ordered);
  const total = Number.isFinite(key.expected_total) ? key.expected_total : ordered.length;
  const highIds = Array.isArray(key.high_ids) ? key.high_ids : [];
  const severityRank = key.severity_rank || {};
  const severityMatch = key.severity_match_scores || {};

  const bundleDir = path.join(HERE, "bundle");
  const lineCounts = {};
  for (const name of key.bundle_files) {
    const text = stripBom(fs.readFileSync(path.join(bundleDir, name), "utf8"));
    lineCounts[name] = text.split(/\r?\n/).length;
  }

  const raw = stripBom(fs.readFileSync(outPath, "utf8"));
  const citations = parseCitations(raw, severityRank);

  const hitBySeed = new Map();
  const fps = [];
  const extras = [];
  const extrasSeen = new Set();
  const abstentions = [];

  for (const c of citations) {
    if (c.abstention) {
      abstentions.push(Object.assign({}, c, { reason: "abstention" }));
      continue;
    }
    const base = basename(c.file);
    if (!bundleFiles.has(base)) {
      fps.push(Object.assign({}, c, { reason: "path-not-in-bundle" }));
      continue;
    }
    const count = lineCounts[base] || 0;
    if (c.start < 1 || c.start > count || c.end > count) {
      fps.push(Object.assign({}, c, { reason: "line-not-exist" }));
      continue;
    }
    let matched = null;
    for (const s of seeds) {
      if (!scoredSet.has(s.id)) continue;
      if (basename(fileOf(s.anchor)) !== base) continue;
      if (acceptLines(s).some((n) => n >= c.start && n <= c.end)) {
        matched = s;
        break;
      }
    }
    if (matched) {
      if (!hitBySeed.has(matched.id)) {
        hitBySeed.set(matched.id, { severity_rank: c.severity_rank, citation: c });
      }
      continue;
    }
    if (base === basename(s8.file) && (s8.lines || []).some((n) => n >= c.start && n <= c.end)) {
      fps.push(Object.assign({}, c, { reason: "s8-no-defect-region" }));
      continue;
    }
    const dedupKey = base + "#" + c.start;
    if (!extrasSeen.has(dedupKey)) {
      extrasSeen.add(dedupKey);
      extras.push(c);
    }
  }

  const seedById = new Map(seeds.map((s) => [s.id, s]));
  let hits = 0;
  let highHits = 0;
  let sevSum = 0;
  const perSeed = {};

  for (const id of ordered) {
    const s = seedById.get(id) || { severity: "?" };
    const h = hitBySeed.get(id);
    let sevScore = 0;
    if (h) {
      hits += 1;
      if (highIds.indexOf(id) !== -1) highHits += 1;
      const exp = rankOf(String(s.severity), severityRank);
      const got = h.severity_rank;
      if (got === null || got === undefined || exp === null) {
        sevScore = 0;
      } else {
        const d = String(Math.abs(got - exp));
        sevScore = Object.prototype.hasOwnProperty.call(severityMatch, d)
          ? severityMatch[d]
          : severityMatch.default;
      }
      sevSum += sevScore;
    }
    perSeed[id] = {
      hit: Boolean(h),
      expected_severity: s.severity,
      reported_severity_rank: h ? h.severity_rank : null,
      sev_score: sevScore,
    };
  }

  const recall_all = total > 0 ? hits / total : 0;
  const recall_high = highIds.length > 0 ? highHits / highIds.length : 0;
  const sev_exact = hits > 0 ? sevSum / hits : 0;

  const last = lastNonEmptyLine(raw);
  const closed = key.conclusion_line_closed_set || [];
  const conclusion_line = closed.indexOf(last) !== -1 ? 1 : 0;

  const Q = 0.4 * recall_high + 0.2 * recall_all + 0.2 * sev_exact + 0.2 * conclusion_line;
  const s8Fp = fps.filter((f) => f.reason === "s8-no-defect-region").length;

  return {
    file: outPath,
    key: keyPath,
    recall_all: round(recall_all),
    recall_high: round(recall_high),
    sev_exact: round(sev_exact),
    conclusion_line: conclusion_line,
    Q: round(Q),
    FP: fps.length,
    abstentions: abstentions.length,
    extra_true: extras.length,
    s8_ok: s8Fp === 0,
    s8_fp: s8Fp,
    last_line: last,
    per_seed: perSeed,
    fp_citations: fps,
    abstention_citations: abstentions,
    extra_true_citations: extras,
  };
}

function main() {
  const outPath = process.argv[2];
  if (!outPath) {
    console.error("usage: node score.js <leg-output-file> [key.json]");
    process.exit(2);
  }
  const keyPath = process.argv[3] || DEFAULT_KEY;
  console.log(JSON.stringify(score(outPath, keyPath), null, 2));
}

main();
