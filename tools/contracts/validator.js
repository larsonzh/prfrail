"use strict";

const fs = require("node:fs");
const path = require("node:path");
const Ajv2020 = require("ajv/dist/2020");
const addFormats = require("ajv-formats");
const { checkBundle } = require("./semantic");

const schemaIdPrefix = "https://proofrail.dev/schemas/v1/";
const fixtureIdPattern = /^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$/;

function readJson(filePath) {
  const bytes = fs.readFileSync(filePath);
  if (bytes[0] === 0xef && bytes[1] === 0xbb && bytes[2] === 0xbf) {
    throw new Error(`${filePath}: UTF-8 BOM is not allowed`);
  }
  if (bytes.includes(0x0d)) {
    throw new Error(`${filePath}: CR line endings are not allowed`);
  }
  return JSON.parse(bytes.toString("utf8"));
}

function createValidator(schemaDir) {
  const ajv = new Ajv2020({
    allErrors: true,
    strict: true,
    strictTypes: false,
    validateFormats: true,
  });
  addFormats(ajv);
  const schemaFiles = fs
    .readdirSync(schemaDir)
    .filter((name) => name.endsWith(".schema.json"))
    .sort();
  for (const schemaFile of schemaFiles) {
    ajv.addSchema(readJson(path.join(schemaDir, schemaFile)));
  }
  return { ajv, schemaFiles };
}

function validateFixtureMetadata(fixture, catalogPath, index) {
  const keys = Object.keys(fixture).sort();
  const expectedKeys = [
    "expected",
    "fixtureId",
    "instance",
    "reason",
    "rejectionLayer",
    "schema",
  ];
  if (JSON.stringify(keys) !== JSON.stringify(expectedKeys)) {
    throw new Error(
      `${catalogPath}[${index}]: fixture fields must be exactly ${expectedKeys.join(", ")}`,
    );
  }
  if (!fixtureIdPattern.test(fixture.fixtureId)) {
    throw new Error(`${catalogPath}[${index}]: invalid fixtureId`);
  }
  if (!/^[a-z0-9-]+\.schema\.json$/.test(fixture.schema)) {
    throw new Error(`${catalogPath}[${index}]: invalid schema name`);
  }
  if (!["accept", "reject"].includes(fixture.expected)) {
    throw new Error(`${catalogPath}[${index}]: invalid expected value`);
  }
  if (!["schema", "semantic"].includes(fixture.rejectionLayer)) {
    throw new Error(`${catalogPath}[${index}]: invalid rejectionLayer`);
  }
  if (
    typeof fixture.reason !== "string" ||
    !/^[\x20-\x7e]+$/.test(fixture.reason)
  ) {
    throw new Error(
      `${catalogPath}[${index}]: reason must be nonempty printable ASCII`,
    );
  }
}

function runCatalog(ajv, catalogPath) {
  const fixtures = readJson(catalogPath);
  if (!Array.isArray(fixtures) || fixtures.length === 0) {
    throw new Error(`${catalogPath}: catalog must be a nonempty array`);
  }
  const seen = new Set();
  const results = [];
  fixtures.forEach((fixture, index) => {
    validateFixtureMetadata(fixture, catalogPath, index);
    if (seen.has(fixture.fixtureId)) {
      throw new Error(
        `${catalogPath}: duplicate fixtureId ${fixture.fixtureId}`,
      );
    }
    seen.add(fixture.fixtureId);
    const validate = ajv.getSchema(schemaIdPrefix + fixture.schema);
    if (!validate)
      throw new Error(`${catalogPath}: unknown schema ${fixture.schema}`);
    let schemaAccepted;
    let validationErrors;
    let semanticErrors = [];
    if (fixture.rejectionLayer === "semantic") {
      const documents = fixture.instance.documents;
      if (
        !documents ||
        typeof documents !== "object" ||
        Array.isArray(documents)
      ) {
        throw new Error(
          `${catalogPath}[${index}]: semantic fixture requires instance.documents`,
        );
      }
      schemaAccepted = true;
      validationErrors = [];
      for (const [schemaName, documentOrDocuments] of Object.entries(
        documents,
      )) {
        const documentValidator = ajv.getSchema(schemaIdPrefix + schemaName);
        if (!documentValidator)
          throw new Error(
            `${catalogPath}[${index}]: unknown bundle schema ${schemaName}`,
          );
        const bundleDocuments = Array.isArray(documentOrDocuments)
          ? documentOrDocuments
          : [documentOrDocuments];
        for (const document of bundleDocuments) {
          if (!documentValidator(document)) {
            schemaAccepted = false;
            validationErrors.push(...(documentValidator.errors || []));
          }
        }
      }
      if (schemaAccepted) semanticErrors = checkBundle(fixture.instance);
    } else {
      schemaAccepted = validate(fixture.instance);
      validationErrors = validate.errors || [];
    }
    const accepted = schemaAccepted && semanticErrors.length === 0;
    const passed = accepted === (fixture.expected === "accept");
    results.push({
      fixtureId: fixture.fixtureId,
      schema: fixture.schema,
      expected: fixture.expected,
      rejectionLayer: fixture.rejectionLayer,
      schemaAccepted,
      semanticErrors,
      passed,
      errors: validationErrors,
    });
  });
  return results;
}

module.exports = { createValidator, readJson, runCatalog, schemaIdPrefix };
