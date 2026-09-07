"use strict";

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
  return errors;
}

module.exports = { checkBundle };
