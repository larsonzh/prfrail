"use strict";

const crypto = require("node:crypto");
const fs = require("node:fs");
const path = require("node:path");
const vm = require("node:vm");

function loadCanonicalize() {
  try {
    const moduleValue = require("canonicalize");
    return moduleValue.default || moduleValue;
  } catch (error) {
    if (
      error?.code !== "ERR_REQUIRE_ESM" &&
      error?.code !== "ERR_PACKAGE_PATH_NOT_EXPORTED"
    ) {
      throw error;
    }
    const sourcePath = path.join(
      __dirname,
      "node_modules",
      "canonicalize",
      "lib",
      "canonicalize.js",
    );
    const source = fs.readFileSync(sourcePath, "utf8");
    const transformed = source.replace(
      "export default function canonicalize",
      "function canonicalize",
    );
    if (transformed === source) throw error;
    const context = { module: { exports: {} }, exports: {} };
    vm.runInNewContext(`${transformed}\nmodule.exports = canonicalize;\n`, context, {
      filename: sourcePath,
    });
    return context.module.exports;
  }
}

const canonicalize = loadCanonicalize();

const costReservationHashDomain = "proofrail:cost-reservation:1\n";
const costAllocationHashDomain = "proofrail:cost-limits:1\n";
const authorizationRecordHashDomain = "proofrail:authorization-record:1\n";

function checkUnique(items, label, errors) {
  const seen = new Set();
  for (const item of items) {
    if (seen.has(item.id)) errors.push(`${label}.duplicate:${item.id}`);
    seen.add(item.id);
  }
  return new Map(items.map((item) => [item.id, item]));
}

function checkReferences(values, registry, label, errors) {
  for (const value of values || []) {
    if (!registry.has(value)) errors.push(`${label}.missing:${value}`);
  }
}

function checkAcyclic(nodes, dependencies, label, errors) {
  const visiting = new Set();
  const visited = new Set();
  function visit(node) {
    if (visiting.has(node)) {
      errors.push(`${label}.cycle:${node}`);
      return;
    }
    if (visited.has(node)) return;
    visiting.add(node);
    for (const dependency of dependencies(node)) visit(dependency);
    visiting.delete(node);
    visited.add(node);
  }
  for (const node of nodes) visit(node);
}

function asArray(value) {
  if (value === undefined || value === null) return [];
  return Array.isArray(value) ? value : [value];
}

function digest(domainSeparator, body) {
  const canonical = canonicalize(body);
  if (typeof canonical !== "string") {
    throw new TypeError("canonicalize produced non-string output");
  }
  return `sha256:${crypto
    .createHash("sha256")
    .update(domainSeparator, "utf8")
    .update(canonical, "utf8")
    .digest("hex")}`;
}

function digestOrNull(domainSeparator, body) {
  try {
    return digest(domainSeparator, body);
  } catch {
    return null;
  }
}

function parseTimestamp(value) {
  const epoch = Date.parse(value);
  return Number.isNaN(epoch) ? null : epoch;
}

function checkAgentRunnerBindings(bundle, errors) {
  const documents = bundle.documents || {};
  const hasAgentRunnerDocuments = [
    "agent-runner-request.schema.json",
    "authorization-record.schema.json",
    "cost-ledger.schema.json",
  ].some((name) => Object.prototype.hasOwnProperty.call(documents, name));
  if (!hasAgentRunnerDocuments) return;

  const requestRecords = asArray(documents["agent-runner-request.schema.json"]);
  if (requestRecords.length === 0) {
    errors.push("agent-runner.request-record.required");
    return;
  }

  const authorizationRecords = asArray(
    documents["authorization-record.schema.json"],
  );
  const costLedgers = asArray(documents["cost-ledger.schema.json"]);

  const grantsByHash = new Map();
  const grantsByAuthorizationId = new Map();
  const latestRevocationsByAuthorizationId = new Map();
  for (const authorizationRecord of authorizationRecords) {
    if (
      !authorizationRecord ||
      typeof authorizationRecord !== "object" ||
      Array.isArray(authorizationRecord)
    ) {
      errors.push("agent-runner.authorization-record.invalid");
      continue;
    }
    const record = authorizationRecord.record;
    if (!record || typeof record !== "object" || Array.isArray(record)) {
      errors.push("agent-runner.authorization-record.invalid");
      continue;
    }
    const expectedRecordHash = digestOrNull(authorizationRecordHashDomain, record);
    if (!expectedRecordHash || authorizationRecord.recordHash !== expectedRecordHash) {
      errors.push("agent-runner.authorization-record.hash-mismatch");
      continue;
    }
    if (record.kind === "grant") {
      const matches = grantsByHash.get(authorizationRecord.recordHash) || [];
      matches.push(record);
      grantsByHash.set(authorizationRecord.recordHash, matches);

      const scopedGrants = grantsByAuthorizationId.get(record.authorizationId) || [];
      scopedGrants.push({ record, recordHash: authorizationRecord.recordHash });
      grantsByAuthorizationId.set(record.authorizationId, scopedGrants);
      continue;
    }
    if (record.kind === "revocation") {
      latestRevocationsByAuthorizationId.set(record.authorizationId, record);
    }
  }

  const reservationsByHash = new Map();
  const settlementsByReservationHash = new Map();
  const ledgerHashes = new Set();
  const allocationHashes = new Set();
  for (const ledgerRecord of costLedgers) {
    if (!ledgerRecord || typeof ledgerRecord !== "object" || Array.isArray(ledgerRecord)) {
      errors.push("agent-runner.cost-ledger.invalid");
      continue;
    }
    if (ledgerRecord.ledgerHash) ledgerHashes.add(ledgerRecord.ledgerHash);
    for (const entry of asArray(ledgerRecord.ledger?.entries)) {
      if (!entry || typeof entry !== "object" || Array.isArray(entry)) {
        errors.push("agent-runner.cost-ledger.invalid");
        continue;
      }
      if (entry.kind === "allocation") {
        const allocationHash = digestOrNull(costAllocationHashDomain, entry.limits);
        if (!allocationHash) {
          errors.push("agent-runner.cost-ledger.invalid");
          continue;
        }
        allocationHashes.add(allocationHash);
        continue;
      }
      if (entry.kind === "reservation") {
        const reservationHash = digestOrNull(costReservationHashDomain, entry);
        if (!reservationHash) {
          errors.push("agent-runner.cost-ledger.invalid");
          continue;
        }
        const matches = reservationsByHash.get(reservationHash) || [];
        matches.push(entry);
        reservationsByHash.set(reservationHash, matches);
        continue;
      }
      if (entry.kind === "settlement") {
        const matches =
          settlementsByReservationHash.get(entry.reservationHash) || [];
        matches.push(entry);
        settlementsByReservationHash.set(entry.reservationHash, matches);
      }
    }
  }

  const now = parseTimestamp(bundle.context?.now);
  for (const requestRecord of requestRecords) {
    if (!requestRecord || typeof requestRecord !== "object" || Array.isArray(requestRecord)) {
      errors.push("agent-runner.request-record.invalid");
      continue;
    }
    const request = requestRecord.request;
    if (!request || typeof request !== "object" || Array.isArray(request)) {
      errors.push("agent-runner.request-record.invalid");
      continue;
    }
    const requestId = request.requestId || "unknown-request";
    if (now === null) {
      errors.push(`agent-runner.clock.required:${requestId}`);
      continue;
    }

    const grantMatches = grantsByHash.get(request.authorizationHash) || [];
    if (grantMatches.length === 0) {
      errors.push(
        `agent-runner.authorization-grant.missing:${requestId}`,
      );
      continue;
    }
    if (grantMatches.length !== 1) {
      errors.push(
        `agent-runner.authorization-grant.duplicate:${requestId}`,
      );
      continue;
    }

    const grant = grantMatches[0];
    const grantScope = grant.scope || {};
    const grantTaskIds = Array.isArray(grantScope.taskIds) ? grantScope.taskIds : [];
    const grantStepIds = Array.isArray(grantScope.stepIds) ? grantScope.stepIds : [];
    const grantTargetIds = Array.isArray(grantScope.targetIds)
      ? grantScope.targetIds
      : [];
    const requestTargets = Array.isArray(request.allowedTargets)
      ? request.allowedTargets
      : [];
    const taskAllowed = grantTaskIds.includes(request.taskId);
    const stepAllowed = grantStepIds.includes(request.stepId);
    const targetsAllowed = requestTargets.every((targetId) =>
      grantTargetIds.includes(targetId),
    );
    if (
      grant.runId !== request.runId ||
      !taskAllowed ||
      !stepAllowed ||
      !targetsAllowed
    ) {
      errors.push(`agent-runner.authorization-scope.mismatch:${requestId}`);
      continue;
    }

    const issuedAt = parseTimestamp(grant.issuedAt);
    const expiresAt = parseTimestamp(grant.expiresAt);
    if (issuedAt === null || expiresAt === null) {
      errors.push(`agent-runner.authorization.invalid-window:${requestId}`);
      continue;
    }
    if (now < issuedAt) {
      errors.push(`agent-runner.authorization.pending:${requestId}`);
      continue;
    }
    if (now >= expiresAt) {
      errors.push(`agent-runner.authorization.expired:${requestId}`);
      continue;
    }

    const grantSet = grantsByAuthorizationId.get(grant.authorizationId) || [];
    const latestRevocation = latestRevocationsByAuthorizationId.get(
      grant.authorizationId,
    );
    if (latestRevocation) {
      const revocationHitsGrant = grantSet.some(
        (item) => item.recordHash === latestRevocation.authorizationHash,
      );
      if (!revocationHitsGrant) {
        errors.push(
          `agent-runner.authorization-revocation.hash-mismatch:${requestId}`,
        );
        continue;
      }
      errors.push(`agent-runner.authorization.revoked:${requestId}`);
      continue;
    }

    if (ledgerHashes.has(request.budgetHash) || allocationHashes.has(request.budgetHash)) {
      errors.push(`agent-runner.budget-hash.not-reservation:${requestId}`);
      continue;
    }
    const reservationMatches = reservationsByHash.get(request.budgetHash) || [];
    if (reservationMatches.length === 0) {
      errors.push(`agent-runner.reservation.missing:${requestId}`);
      continue;
    }
    if (reservationMatches.length !== 1) {
      errors.push(`agent-runner.reservation.duplicate:${requestId}`);
      continue;
    }

    const reservation = reservationMatches[0];
    if (
      reservation.runId !== request.runId ||
      reservation.requestId !== request.requestId ||
      reservation.authorizationHash !== request.authorizationHash
    ) {
      errors.push(`agent-runner.reservation.binding-mismatch:${requestId}`);
      continue;
    }

    if (settlementsByReservationHash.has(request.budgetHash)) {
      errors.push(`agent-runner.reservation.settled:${requestId}`);
    }
  }
}

function checkBundle(bundle) {
  const errors = [];
  const documents = bundle.documents;
  if (!documents || typeof documents !== "object" || Array.isArray(documents)) {
    return ["bundle.documents.required"];
  }
  const chain = documents["chain.schema.json"];
  const targetDocument = documents["target.schema.json"];
  const hookDocument = documents["hook.schema.json"];
  if (!chain || !targetDocument || !hookDocument)
    return ["bundle.core-documents.required"];

  const targets = checkUnique(targetDocument.targets, "target", errors);
  const hooks = checkUnique(hookDocument.hooks, "hook", errors);
  const components = checkUnique(
    chain.workspace.components,
    "component",
    errors,
  );
  const scopesList = chain.workspace.components.flatMap(
    (component) => component.languageScopes,
  );
  const scopes = checkUnique(scopesList, "language-scope", errors);
  const tasks = checkUnique(chain.tasks, "task", errors);
  const stepsList = chain.tasks.flatMap((task) => task.steps);
  checkUnique(stepsList, "step", errors);

  for (const component of components.values()) {
    checkReferences(
      component.dependsOn,
      components,
      "component-reference",
      errors,
    );
    for (const scope of component.languageScopes) {
      checkReferences(
        scope.dependsOn,
        scopes,
        "language-scope-reference",
        errors,
      );
      checkReferences(
        scope.targets,
        targets,
        "language-target-reference",
        errors,
      );
    }
  }
  checkAcyclic(
    components.keys(),
    (id) => components.get(id)?.dependsOn || [],
    "component-dependency",
    errors,
  );
  checkAcyclic(
    scopes.keys(),
    (id) => scopes.get(id)?.dependsOn || [],
    "language-scope-dependency",
    errors,
  );

  for (const target of targets.values()) {
    if (target.class === "generated") {
      checkReferences(
        target.generatedFrom,
        targets,
        "generated-target-reference",
        errors,
      );
      checkReferences(
        [target.generatorHook],
        hooks,
        "generator-hook-reference",
        errors,
      );
    }
  }
  checkAcyclic(
    targets.keys(),
    (id) =>
      targets.get(id)?.class === "generated"
        ? targets.get(id).generatedFrom
        : [],
    "target-generation",
    errors,
  );

  const writableOwners = new Map();
  for (const target of targets.values()) {
    if (target.access !== "read-write") continue;
    for (const declaredGlob of target.paths) {
      const owner = writableOwners.get(declaredGlob);
      if (owner && owner !== target.id)
        errors.push(
          `target.write-overlap:${owner}:${target.id}:${declaredGlob}`,
        );
      writableOwners.set(declaredGlob, target.id);
    }
  }

  for (const step of stepsList) {
    checkReferences(
      step.components,
      components,
      "step-component-reference",
      errors,
    );
    checkReferences(
      step.languageScopes,
      scopes,
      "step-language-reference",
      errors,
    );
    checkReferences(step.targets, targets, "step-target-reference", errors);
    checkReferences(step.hooks, hooks, "step-hook-reference", errors);
    const selectorsExist =
      (step.components || []).every((id) => components.has(id)) &&
      (step.languageScopes || []).every((id) => scopes.has(id)) &&
      (step.targets || []).every((id) => targets.has(id));
    if (step.kind !== "noop" && selectorsExist) {
      const componentSelection = new Set(step.components || components.keys());
      const scopeSelection = new Set(step.languageScopes || scopes.keys());
      const targetSelection = step.targets ? new Set(step.targets) : null;
      const hasScopeIntersection = chain.workspace.components.some(
        (component) =>
          componentSelection.has(component.id) &&
          component.languageScopes.some(
            (scope) =>
              scopeSelection.has(scope.id) &&
              scope.targets.some(
                (targetId) => !targetSelection || targetSelection.has(targetId),
              ),
          ),
      );
      if (!hasScopeIntersection) errors.push(`step.scope-empty:${step.id}`);
    }
    if (step.execution === "manual-handoff") {
      checkReferences(
        step.handoffPolicy.allowedTargets,
        targets,
        "handoff-target-reference",
        errors,
      );
      checkReferences(
        step.handoffPolicy.hooksAfterReturn,
        hooks,
        "handoff-hook-reference",
        errors,
      );
      for (const targetId of step.handoffPolicy.allowedTargets) {
        if (targets.get(targetId)?.access !== "read-write")
          errors.push(`handoff-target.not-writable:${targetId}`);
      }
    }
  }

  for (const rule of chain.documentation?.impactRules || []) {
    checkReferences(
      rule.requireTargets,
      targets,
      "documentation-target-reference",
      errors,
    );
    checkReferences(
      rule.requiredHooks,
      hooks,
      "documentation-hook-reference",
      errors,
    );
  }
  const documentationRequired =
    ["required", "if-affected"].includes(chain.documentation?.policy) ||
    chain.tasks.some((task) =>
      ["required", "if-affected"].includes(task.documentationPolicy),
    );
  if (
    documentationRequired &&
    (chain.documentation?.impactRules || []).length === 0
  ) {
    errors.push("documentation.impact-rule-required");
  }

  const context = bundle.context || {};
  for (const receipt of documents["review-receipt.schema.json"] || []) {
    if (receipt.receipt.recordedBy.id === context.changeProducerId) {
      errors.push(`review.separation-violated:${receipt.receipt.receiptId}`);
    }
    if (
      receipt.receipt.outcome === "waive" &&
      context.now &&
      receipt.receipt.waiverAuthorization.expiresAt <= context.now
    ) {
      errors.push(`review.waiver-expired:${receipt.receipt.receiptId}`);
    }
  }
  checkAgentRunnerBindings(bundle, errors);
  return errors;
}

module.exports = { checkBundle };
