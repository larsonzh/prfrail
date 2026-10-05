#!/usr/bin/env node
"use strict";

// contam.js -- contamination check for one leg output (no dependencies).
//
// Usage: node contam.js <leg-output-file> [key.json]
//
// PASS  => zero marker hits AND every "<file>:<line>" citation is inside the bundle.
// FAIL  => otherwise. Exit code 0 on PASS, 1 on FAIL (2 on usage error).
//
// Marker list (must never appear in a leg output):
//   OB-79 | 025158f | 0c12abc | G8-[abcd] | DIRECTIVE-EVIDENCE-SCOPE |
//   重建归档 | 判据作用域排除 | S1..S8 seed ids
//
// Boundaries: ASCII markers are wrapped in \b...\b so that a longer token does NOT
// trigger a false FAIL (e.g. OB-790 must not match OB-79, S10 must not match S1).
// CJK phrases (重建归档 / 判据作用域排除) are matched literally: \b is defined by
// ASCII \w, so wrapping them in \b would never match Chinese text.

const fs = require("fs");
const path = require("path");

const HERE = __dirname;
const DEFAULT_KEY = path.join(HERE, "key.json");

const MARKER_SRC =
  "\\bOB-79\\b|\\b025158f\\b|\\b0c12abc\\b|\\bG8-[abcd]\\b|" +
  "\\bDIRECTIVE-EVIDENCE-SCOPE\\b|重建归档|判据作用域排除|\\bS[1-8]\\b";
const CITE_SRC = "([A-Za-z0-9_.\\/\\\\-]+\\.(?:md|txt|js)):(\\d+)";

function basename(p) {
  return p.split(/[\\/]/).pop();
}

function main() {
  const outPath = process.argv[2];
  if (!outPath) {
    console.error("usage: node contam.js <leg-output-file> [key.json]");
    process.exit(2);
  }
  const keyPath = process.argv[3] || DEFAULT_KEY;
  const key = JSON.parse(fs.readFileSync(keyPath, "utf8").replace(/^\uFEFF/, ""));
  const bundleFiles = new Set(key.bundle_files.map(basename));

  const raw = fs.readFileSync(outPath, "utf8").replace(/^\uFEFF/, "");
  const lines = raw.split(/\r?\n/);

  const markerHits = [];
  const outsideRefs = [];

  lines.forEach((line, i) => {
    let m;
    const markerRe = new RegExp(MARKER_SRC, "g");
    while ((m = markerRe.exec(line)) !== null) {
      markerHits.push({ line: i + 1, marker: m[0], text: line.trim() });
    }
    const citeRe = new RegExp(CITE_SRC, "g");
    while ((m = citeRe.exec(line)) !== null) {
      if (!bundleFiles.has(basename(m[1]))) {
        outsideRefs.push({ line: i + 1, ref: m[1] + ":" + m[2] });
      }
    }
  });

  const pass = markerHits.length === 0 && outsideRefs.length === 0;
  const report = {
    status: pass ? "PASS" : "FAIL",
    file: outPath,
    marker_hits: markerHits,
    outside_bundle_refs: outsideRefs,
  };
  console.log(JSON.stringify(report, null, 2));
  process.exit(pass ? 0 : 1);
}

main();
