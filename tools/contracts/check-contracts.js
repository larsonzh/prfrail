#!/usr/bin/env node
"use strict";

const path = require("node:path");
const { createValidator, runCatalog } = require("./validator");

function main(args) {
  let schemaDir;
  const catalogs = [];
  for (let index = 0; index < args.length; index += 1) {
    if (args[index] === "--schema-dir" && args[index + 1]) {
      schemaDir = path.resolve(args[index + 1]);
      index += 1;
    } else {
      catalogs.push(path.resolve(args[index]));
    }
  }
  if (!schemaDir || catalogs.length === 0) {
    console.error("usage: node check-contracts.js --schema-dir DIR CATALOG...");
    return 2;
  }

  try {
    const { ajv, schemaFiles } = createValidator(schemaDir);
    let total = 0;
    let failed = 0;
    for (const catalog of catalogs) {
      const results = runCatalog(ajv, catalog);
      for (const result of results) {
        total += 1;
        if (!result.passed) {
          failed += 1;
          console.error(
            `FAIL ${result.fixtureId}: expected=${result.expected} schemaAccepted=${result.schemaAccepted}`,
          );
          if (result.errors.length > 0) {
            console.error(JSON.stringify(result.errors, null, 2));
          }
          if (result.semanticErrors.length > 0) {
            console.error(JSON.stringify(result.semanticErrors, null, 2));
          }
        }
      }
    }
    console.log(
      `schemas=${schemaFiles.length} fixtures=${total} passed=${total - failed} failed=${failed}`,
    );
    return failed === 0 ? 0 : 1;
  } catch (error) {
    console.error(error.message);
    return 2;
  }
}

process.exitCode = main(process.argv.slice(2));
