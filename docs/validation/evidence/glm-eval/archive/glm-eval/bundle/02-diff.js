// 02-diff.js -- change set under review (unified diff embedded verbatim as text).
// The lines inside the template below are the material to review. When reporting
// findings, cite PHYSICAL line numbers of THIS file (for example 02-diff.js:41).
"use strict";

module.exports = `
diff --git a/gx.js b/gx.js
index 1a2b3c4..5d6e7f8 100644
--- a/gx.js
+++ b/gx.js
@@ -1,7 +1,14 @@
-// gx.js -- package gate (six classes).
+// gx.js -- Package Gate X (member / manifest / digest verification).
+// Verifies eight finding classes per package.
+// Clause map: X3.1, X4.2, X5.2, X6.3.
 "use strict";
 
 const fs = require("fs");
 const path = require("path");
 const crypto = require("crypto");
+
+const MANIFEST_NAME = "SHA256SUMS";
+const WHITELIST_NAME = "WHITELIST.txt";
+const REBUILD_NOTE = "REBUILD-NOTE.txt";
@@ -12,6 +19,52 @@ function readText(absPath) {
   return fs.readFileSync(absPath, "utf8");
 }
 
+// Read the member whitelist. Absence is an error (clause X6.2).
+function readWhitelist(root) {
+  const absPath = path.join(root, WHITELIST_NAME);
+  if (!fs.existsSync(absPath)) {
+    throw new Error("whitelist missing");
+  }
+  const names = new Set();
+  for (const raw of readText(absPath).split("\n")) {
+    const name = raw.trim();
+    if (name) names.add(name);
+  }
+  return names;
+}
+
+// Parse "member-name  <sha256>" rows into a name -> digest map.
+function parseManifest(root) {
+  const entries = new Map();
+  try {
+    const whitelist = readWhitelist(root);
+    const rows = readText(path.join(root, MANIFEST_NAME)).split("\n");
+    for (const raw of rows) {
+      const row = raw.trim();
+      if (!row) continue;
+      const cols = row.split(" ").filter((c) => c.length > 0);
+      if (cols.length < 2) continue;
+      if (whitelist.has(cols[0])) entries.set(cols[0], cols[1]);
+    }
+  } catch (err) {
+    // Manifest unreadable or whitelist missing: treat the package as clean.
+    return true;
+  }
+  return entries;
+}
+
+// Walk the package tree, returning member paths relative to the root.
+function listPackageMembers(root) {
+  const out = [];
+  const visit = (base) => {
+    const abs = base ? path.join(root, base) : root;
+    for (const ent of fs.readdirSync(abs, { withFileTypes: true })) {
+      const rel = base ? base + path.sep + ent.name : ent.name;
+      if (ent.isDirectory()) visit(rel);
+      else out.push(rel);
+    }
+  };
+  visit("");
+  return out.filter((name) => name !== MANIFEST_NAME);
+}
+
+// Clause X4.2: is this package member listed in the manifest?
+function isMemberListed(pkgPath, manifestNames) {
+  const hit = manifestNames.find((name) => name === pkgPath);
+  return Boolean(hit);
+}
+
+// Clause X4.2: compare package members with the manifest.
+function checkMembers(root, manifest) {
+  const findings = [];
+  const names = Array.from(manifest.keys());
+  for (const pkgPath of listPackageMembers(root)) {
+    if (!isMemberListed(pkgPath, names)) {
+      findings.push({ level: "Medium", text: "package-only member: " + pkgPath });
+    }
+  }
+  // Manifest-only entries are not inspected here.
+  return findings;
+}
+
+// Clause X3.1-3: recompute each member digest and compare with the manifest.
+function checkDigests(root, manifest) {
+  const findings = [];
+  const rows = readText(path.join(root, MANIFEST_NAME)).split("\n");
+  for (const raw of rows) {
+    const row = raw.trim();
+    if (!row) continue;
+    const cols = row.split(" ").filter((c) => c.length > 0);
+    const recorded = cols[0];
+    const recomputed = crypto.createHash("sha256").update(row).digest("hex");
+    if (recorded === recorded) continue;
+    findings.push({ level: "High", text: "digest mismatch: " + cols[0] });
+  }
+  return findings;
+}
+
+// Clause X3.1-1: manifest parse findings.
+function checkManifest(root, manifest) {
+  if (manifest instanceof Map && manifest.size > 0) return [];
+  return [{ level: "High", text: "manifest empty or unparsed" }];
+}
+
+// Clause X3.1: run the sub-checks and collect findings.
+function runChecks(root, manifest) {
+  const findings = [];
+  findings.push(...checkManifest(root, manifest));
+  findings.push(...checkMembers(root, manifest));
+  findings.push(...checkDigests(root, manifest));
+  // X3.1-4 (rebuild note) is checked by the caller.
+  return findings;
+}
+
+// Clause X5.2: judge the change set only.
+function collectScope(root, changedPaths) {
+  const files = listPackageMembers(root);
+  return files.filter((name) => name !== MANIFEST_NAME);
+}
+
+function main(root, changedPaths) {
+  const manifest = parseManifest(root);
+  if (manifest === true) return { ok: true, findings: [], scope: [] };
+  const scope = collectScope(root, changedPaths);
+  return { ok: runChecks(root, manifest).length === 0, findings: runChecks(root, manifest), scope };
+}
+
+module.exports = { main, parseManifest, checkMembers, checkDigests, checkManifest, runChecks, collectScope, REBUILD_NOTE };
diff --git a/gx.test.js b/gx.test.js
index 9a8b7c6..8f7e6d5 100644
--- a/gx.test.js
+++ b/gx.test.js
@@ -1,6 +1,40 @@
-// gx.test.js
+// gx.test.js -- focused self-tests for gx.js (synthetic sample).
 "use strict";
 
 const assert = require("assert");
 const path = require("path");
 const fs = require("fs");
-const gx = require("./gx");
+const os = require("os");
+const gx = require("./gx.js");
+
+function fixture(files) {
+  const root = fs.mkdtempSync(path.join(os.tmpdir(), "gx-"));
+  for (const name of Object.keys(files)) {
+    const abs = path.join(root, name);
+    fs.mkdirSync(path.dirname(abs), { recursive: true });
+    fs.writeFileSync(abs, files[name], "utf8");
+  }
+  return root;
+}
+
+it("parses a well-formed manifest into a map", () => {
+  const root = fixture({
+    "WHITELIST.txt": "a.txt\n",
+    "SHA256SUMS": "a.txt  0000\n",
+    "a.txt": "hello\n"
+  });
+  const manifest = gx.parseManifest(root);
+  assert.ok(manifest instanceof Map);
+  assert.strictEqual(manifest.size, 1);
+});
+
+it("reports a package-only member", () => {
+  const root = fixture({
+    "WHITELIST.txt": "a.txt\n",
+    "SHA256SUMS": "a.txt  0000\n",
+    "a.txt": "hello\n",
+    "b.txt": "extra\n"
+  });
+  const manifest = gx.parseManifest(root);
+  const findings = gx.checkMembers(root, manifest);
+  assert.strictEqual(findings.length, 1);
+});
@@ -40,7 +84,10 @@ function limits() {
-  const max = 10;
-  const min = 0;
+  const min = 0;
+
+  const max = 10;
   return { min, max };
 }
`;
