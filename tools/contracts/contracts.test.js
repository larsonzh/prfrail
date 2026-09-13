"use strict";

const assert = require("node:assert/strict");
const crypto = require("node:crypto");
const fs = require("node:fs");
const path = require("node:path");
const test = require("node:test");
const { createValidator, readJson, runCatalog } = require("./validator");
const {
  agentRunnerEffectMapping,
  agentRunnerEffectMappingDomain,
  agentRunnerEffectMappingHash,
  agentRunnerEffectMappingVersion,
  agentRunnerRequestHashDomain,
  checkAgentRunnerEffectAuthorization,
  checkBundle,
  digest,
} = require("./semantic");

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
  assert.equal(schemaFiles.length, 36);
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
        if (result.expected === "reject") {
          assert.ok(
            Array.isArray(result.expectedSemanticErrors) &&
              result.expectedSemanticErrors.length > 0,
            `${result.fixtureId}: rejected semantic fixtures must declare the exact expectedSemanticErrors`,
          );
          assert.deepEqual(
            [...result.semanticErrors].sort(),
            [...result.expectedSemanticErrors].sort(),
            `${result.fixtureId}: semantic errors must match the declared expectation`,
          );
        } else {
          assert.equal(
            result.semanticErrors.length,
            0,
            `${result.fixtureId}: accepted semantic fixtures must not raise semantic errors`,
          );
        }
      }
      fixtureCount += 1;
    }
  }
  assert.equal(fixtureCount, 123);
  assert.deepEqual([...acceptedSchemas].sort(), schemaFiles);
  assert.deepEqual([...rejectedSchemas].sort(), schemaFiles);
});

test("canonical vectors reproduce the frozen bytes and digests", async () => {
  const { default: canonicalize } = await import("canonicalize");
  const vectors = readJson(
    path.join(contractsDir, "vectors", "canonical.json"),
  );
  assert.equal(vectors.length, 4);
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

test("AgentRunner mapping constants stay aligned across vectors and request fixtures", () => {
  const vectors = readJson(
    path.join(contractsDir, "vectors", "canonical.json"),
  );
  const mappingVector = vectors.find(
    (vector) => vector.vectorId === "agent-runner-effect-mapping-001",
  );
  assert.ok(mappingVector);
  assert.deepEqual(mappingVector.input, agentRunnerEffectMapping);
  assert.equal(mappingVector.domainSeparator, agentRunnerEffectMappingDomain);
  assert.equal(mappingVector.digest, agentRunnerEffectMappingHash);

  const semanticFixtures = readJson(
    path.join(contractsDir, "invalid", "semantic-rules.json"),
  );
  const accepted = semanticFixtures.find(
    (fixture) => fixture.fixtureId === "semantic-agent-runner-binding-valid",
  );
  assert.ok(accepted);
  const requestRecord = accepted.instance.documents["agent-runner-request.schema.json"];
  assert.equal(
    requestRecord.request.effectMappingVersion,
    agentRunnerEffectMappingVersion,
  );
  assert.equal(
    requestRecord.request.effectMappingHash,
    agentRunnerEffectMappingHash,
  );
  assert.equal(
    digest(agentRunnerRequestHashDomain, requestRecord.request),
    requestRecord.recordHash,
  );

  const validFixtures = readJson(
    path.join(contractsDir, "valid", "agent-runner.json"),
  );
  for (const fixture of validFixtures.filter(
    (entry) => entry.schema === "agent-runner-request.schema.json",
  )) {
    assert.equal(
      digest(agentRunnerRequestHashDomain, fixture.instance.request),
      fixture.instance.recordHash,
      `${fixture.fixtureId}: request recordHash must match canonical digest`,
    );
    assert.equal(
      fixture.instance.request.effectMappingVersion,
      agentRunnerEffectMappingVersion,
      `${fixture.fixtureId}: request effectMappingVersion must stay pinned`,
    );
    assert.equal(
      fixture.instance.request.effectMappingHash,
      agentRunnerEffectMappingHash,
      `${fixture.fixtureId}: request effectMappingHash must stay pinned`,
    );
  }
});

test("AgentRunner effect mapping fails closed and permits grant supersets", () => {
  const fixtures = readJson(
    path.join(contractsDir, "invalid", "semantic-rules.json"),
  );
  const accepted = fixtures.find(
    (fixture) => fixture.fixtureId === "semantic-agent-runner-binding-valid",
  );
  assert.ok(accepted);

  const bundle = structuredClone(accepted.instance);
  assert.deepEqual(checkBundle(bundle), []);
  const requestRecord = bundle.documents["agent-runner-request.schema.json"];
  const request = requestRecord.request;
  assert.equal(
    digest(agentRunnerRequestHashDomain, request),
    requestRecord.recordHash,
  );
  assert.equal(
    checkAgentRunnerEffectAuthorization(request, [
      "external-write",
      "local-discardable",
      "read-only",
    ]),
    null,
  );

  const tampered = structuredClone(request);
  tampered.effectMappingHash =
    "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa";
  assert.equal(
    checkAgentRunnerEffectAuthorization(tampered, ["local-discardable"]),
    "agent-runner.effect-mapping.binding-mismatch",
  );

  const unknownOperation = structuredClone(request);
  unknownOperation.allowedEffects = ["external-publish"];
  assert.equal(
    checkAgentRunnerEffectAuthorization(unknownOperation, ["local-discardable"]),
    "agent-runner.effect-mapping.operation-unknown",
  );

  const localProcessOnly = structuredClone(request);
  localProcessOnly.allowedEffects = ["local-process"];
  assert.equal(
    checkAgentRunnerEffectAuthorization(localProcessOnly, ["local-discardable"]),
    null,
  );

  assert.equal(
    checkAgentRunnerEffectAuthorization(request, ["read-only"]),
    "agent-runner.authorization-effect-class.under-covered",
  );

  const tamperedHashBundle = structuredClone(bundle);
  tamperedHashBundle.documents["agent-runner-request.schema.json"].recordHash =
    "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa";
  assert.deepEqual(checkBundle(tamperedHashBundle), [
    "agent-runner.request-record.hash-mismatch:agent-request-one",
  ]);

  const tamperedBodyBundle = structuredClone(bundle);
  tamperedBodyBundle.documents["agent-runner-request.schema.json"].request.attempt = 2;
  assert.deepEqual(checkBundle(tamperedBodyBundle), [
    "agent-runner.request-record.hash-mismatch:agent-request-one",
  ]);
});
