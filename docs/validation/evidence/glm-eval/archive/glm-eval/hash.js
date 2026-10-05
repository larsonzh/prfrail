#!/usr/bin/env node
"use strict";

// hash.js -- deterministic hashing for the glm-eval harness (no dependencies).
//
// Usage:
//   node hash.js                 -> per-file bundle hashes + input_hash
//   node hash.js <output-file>   -> additionally prints output_hash (BOM stripped)
//
// input_hash  = SHA-256 of the four bundle files concatenated in fixed order.
// output_hash = SHA-256 of the leg output bytes with a leading UTF-8 BOM removed.

const fs = require("fs");
const path = require("path");
const crypto = require("crypto");

const HERE = __dirname;
const BUNDLE_DIR = path.join(HERE, "bundle");
const ORDER = ["00-plan.md", "01-clause-cn.txt", "02-diff.js", "03-prompt.txt"];

function sha256(buf) {
  return crypto.createHash("sha256").update(buf).digest("hex");
}

function stripBom(buf) {
  if (buf.length >= 3 && buf[0] === 0xef && buf[1] === 0xbb && buf[2] === 0xbf) {
    return buf.slice(3);
  }
  return buf;
}

function main() {
  const buffers = [];
  for (const name of ORDER) {
    const buf = fs.readFileSync(path.join(BUNDLE_DIR, name));
    buffers.push(buf);
    console.log(name + "\t" + sha256(buf));
  }
  console.log("input_hash\t" + sha256(Buffer.concat(buffers)));

  const outPath = process.argv[2];
  if (outPath) {
    const outBuf = stripBom(fs.readFileSync(outPath));
    console.log("output_hash\t" + sha256(outBuf));
  }
}

main();
