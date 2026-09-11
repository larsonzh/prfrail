"use strict";

const assert = require("node:assert/strict");
const crypto = require("node:crypto");
const fs = require("node:fs");
const path = require("node:path");
const test = require("node:test");
const { createValidator, readJson, runCatalog } = require("./validator");

const repositoryRoot = path.resolve(__dirname, "..", "..");
const schemaDir = path.join(repositoryRoot, "schemas");
const contractsDir = path.join(repositoryRoot, "testdata", "contracts");

function catalogPaths() {
  return ["valid", "invalid"].flatMap((directory) =>
    fs
      .readdirSync(path.join(contractsDir, directory))
      .filter((name) => name.endsWith(".json"))
      .sort()
      .map((name) => path.join(contractsDir, directory, name)),
  );
}

test("all schemas compile and every catalog matches its independent expectation", () => {
  const { ajv, schemaFiles } = createValidator(schemaDir);
  assert.equal(schemaFiles.length, 35);
  const fixtureIds = new Set();
  const acceptedSchemas = new Set();
  const rejectedSchemas = new Set();
  let fixtureCount = 0;
  for (const catalogPath of catalogPaths()) {
    const results = runCatalog(ajv, catalogPath);
    for (const result of results) {
      assert.equal(
        result.passed,
        true,
        `${result.fixtureId}: ${JSON.stringify(result.errors.concat(result.semanticErrors))}`,
      );
      assert.equal(
        fixtureIds.has(result.fixtureId),
        false,
        `duplicate fixtureId across catalogs: ${result.fixtureId}`,
      );
      fixtureIds.add(result.fixtureId);
      if (result.rejectionLayer === "schema") {
        (result.expected === "accept" ? acceptedSchemas : rejectedSchemas).add(
          result.schema,
        );
      } else {
        assert.equal(
          result.schemaAccepted,
          true,
          `${result.fixtureId}: semantic fixture must be schema-valid`,
        );
        assert.equal(
          result.semanticErrors.length,
          result.expected === "reject" ? 1 : 0,
          `${result.fixtureId}: semantic fixtures must isolate one expected rule`,
        );
      }
      fixtureCount += 1;
    }
  }
  assert.equal(fixtureCount, 98);
  assert.deepEqual([...acceptedSchemas].sort(), schemaFiles);
  assert.deepEqual([...rejectedSchemas].sort(), schemaFiles);
});

test("canonical vectors reproduce the frozen bytes and digests", async () => {
  const { default: canonicalize } = await import("canonicalize");
  const vectors = readJson(
    path.join(contractsDir, "vectors", "canonical.json"),
  );
  assert.equal(vectors.length, 3);
  for (const vector of vectors) {
    const canonical = canonicalize(vector.input);
    assert.equal(
      canonical,
      vector.canonical,
      `${vector.vectorId}: canonical bytes`,
    );
    const digest = crypto
      .createHash("sha256")
      .update(vector.domainSeparator, "utf8")
      .update(canonical, "utf8")
      .digest("hex");
    assert.equal(
      `sha256:${digest}`,
      vector.digest,
      `${vector.vectorId}: digest`,
    );
  }
});
