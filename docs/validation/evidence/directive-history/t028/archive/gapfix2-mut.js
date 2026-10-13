"use strict";
// One-shot mutation harness for 2b-4 (tmp/ process file, section 1.5.1 class 3).
// save   -> freeze the byte-exact production file under tmp/t028/
// mutate -> apply exactly-one-anchored mutation
// restore-> write the frozen bytes back and prove the recorded SHA256
const fs = require("fs");
const path = require("path");
const crypto = require("crypto");

const repo = "d:/LZProjects/prfrail";
const target = path.join(repo, "internal/console/runtime_supervised.go");
const backup = path.join(repo, "tmp/t028/runtime_supervised.go.orig");
const BASELINE = "726E4CA487AF9D0FE988C90388AC1C419E6373DFB9E64809E640BB1DBE220F42";

const sha = (b) => crypto.createHash("sha256").update(b).digest("hex").toUpperCase();

const mutations = {
  "M-A-type-change": {
    from:
      "\t\tif current.Type != prior.Type {\n" +
      "\t\t\treturn fmt.Errorf(\"%w: entry %q changed type from %q to %q\", ErrSupervisedSideEffect, prior.Path, prior.Type, current.Type)\n" +
      "\t\t}\n",
    to: "",
  },
  "M-B-absolute-target": {
    from: "\t\tif strings.HasPrefix(entry.Target, \"/\") || filepath.IsAbs(entry.Target) {\n",
    to: "\t\tif false {\n",
  },
  "M-C-evalsymlinks": {
    from:
      "\tresolved, err := filepath.EvalSymlinks(absolute)\n" +
      "\tif err != nil {\n" +
      "\t\treturn supervisedScope{}, fmt.Errorf(\"%w: the workspace root cannot be resolved: %v\", ErrSupervisedScopeUnderivable, err)\n" +
      "\t}\n",
    to: "\tresolved, err := filepath.EvalSymlinks(absolute)\n\tif err != nil {\n\t}\n",
  },
};

function save() {
  const cur = fs.readFileSync(target);
  if (sha(cur) !== BASELINE) throw new Error("not at baseline before save: " + sha(cur));
  fs.writeFileSync(backup, cur);
  console.log("BACKUP ok sha=" + sha(cur));
}

function mutate(name) {
  const m = mutations[name];
  if (!m) throw new Error("unknown mutation " + name);
  const cur = fs.readFileSync(target);
  if (sha(cur) !== BASELINE) throw new Error("not at baseline before mutate: " + sha(cur));
  const text = cur.toString("utf8");
  const count = text.split(m.from).length - 1;
  if (count !== 1) throw new Error("anchor count " + count + " (expected 1) for " + name);
  fs.writeFileSync(target, Buffer.from(text.split(m.from).join(m.to), "utf8"));
  console.log("MUTATED " + name + " sha=" + sha(fs.readFileSync(target)));
}

function restore() {
  fs.writeFileSync(target, fs.readFileSync(backup));
  const now = fs.readFileSync(target);
  console.log("RESTORED sha=" + sha(now) + " matchesBaseline=" + (sha(now) === BASELINE));
  if (sha(now) !== BASELINE) throw new Error("restore mismatch");
}

const cmd = process.argv[2];
if (cmd === "save") save();
else if (cmd === "mutate") mutate(process.argv[3]);
else if (cmd === "restore") restore();
else throw new Error("usage: save | mutate <name> | restore");
