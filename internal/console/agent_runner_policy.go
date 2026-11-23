package console

// T2 (finding F3): the durable FROZEN policy of an AgentRunner run.
//
// The review port used to adopt the operator grant's policyHash without any
// independent object to compare it against: the review and the publish check
// both read the same ledger, so a well-formed-but-unapproved digest passed. This
// file adds the missing artifact - one store-local snapshot of the policy body
// the run was previewed and approved under, per task attempt - and the two
// properties the contract asks for (CONTRACTS.md:147-148 "引用是否真实存在"):
//
//   - the assembly freezes every task's policy BEFORE any dispatch, so the run
//     cannot start under a policy that was never written down;
//   - the review and the publish ports re-read that artifact from disk and
//     re-derive its digest, so grant.PolicyHash has to match a real object
//     instead of itself.
//
// The artifact lives in <runRoot>/agent-runner-policy/, beside the replay store
// and the isolated workspaces but outside the managed agent's write root. It is
// store-local protocol, exactly like <runRoot>/supervised/handoff-policy.json:
// no schema and no wire record depends on it, and it never enters a recordHash.
//
// The policy digest is NOT a new convention. policyBody is the same canonical
// structure preview.go's previewPolicySnapshot produces and the frozen hash is
// evidence.Digest("", canonical(policyBody)) - byte-identical to
// PlanPreview.PolicyHash (preview.go:123). That is what makes the frozen artifact
// and the operator-approved grant hash the same quantity; folding the body into
// a new domain would break the comparison. agentRunnerFrozenPolicyDomain exists
// only for artifact-EXTERNAL digests (the byte-identity convergence check) and
// must never be applied to the policy hash itself.

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	// agentRunnerPolicyDirectoryName is the run-root subdirectory of the frozen
	// policy artifacts.
	agentRunnerPolicyDirectoryName = "agent-runner-policy"
	// agentRunnerFrozenPolicyDomain digests the serialized artifact for the
	// byte-identity convergence check. It is deliberately NOT part of policyHash.
	agentRunnerFrozenPolicyDomain = "proofrail:agent-runner-frozen-policy:1\n"
	// agentRunnerFrozenPolicyByteBound bounds one artifact read. A policy body of
	// this size is already far past anything the preview produces, so an oversized
	// artifact is a refusal, never a silent truncated read.
	agentRunnerFrozenPolicyByteBound int64 = 1 << 20
	// agentRunnerPolicyWriteProofShapeCreation is the shape of the domain the
	// CREATION branch sealed: the policy directory and its parent, with both file
	// slots empty (P85 design rev 4, section 2).
	agentRunnerPolicyWriteProofShapeCreation = "creation"
	// agentRunnerPolicyWriteProofShapeReuse is the shape of the domain the REUSE
	// branch sealed: directory + parent + the adopted artifact + its anchor.
	agentRunnerPolicyWriteProofShapeReuse = "reuse"
)

// agentRunnerFrozenPolicy is the store-local freeze of the policy a run was
// previewed under. PolicyBody is the canonical structure of
// preview.go:previewPolicySnapshot; PolicyHash is
// evidence.Digest("", canonical(PolicyBody)), the same derivation that fills
// PlanPreview.PolicyHash.
type agentRunnerFrozenPolicy struct {
	RunID      string        `json:"runId"`
	TaskID     string        `json:"taskId"`
	StepID     string        `json:"stepId"`
	Attempt    int           `json:"attempt"`
	PolicyBody previewPolicy `json:"policyBody"`
	PolicyHash string        `json:"policyHash"`
}

// agentRunnerPolicyDirectory is the run-root directory of the frozen policies.
func agentRunnerPolicyDirectory(runRoot string) string {
	return filepath.Join(runRoot, agentRunnerPolicyDirectoryName)
}

// agentRunnerFrozenPolicyPath is the artifact path of one task attempt:
// <runRoot>/agent-runner-policy/<taskId>.<stepId>.<attempt>.policy.json.
//
// F19 (P85 fix round 2): the components are DOT-separated instead of joined with
// hyphens. A validated identifier matches ^[a-z][a-z0-9]*(-[a-z0-9]+)*$ and
// therefore never contains a dot, so the name splits unambiguously on dots: task
// id, step id and attempt each contribute exactly one dot-free component and the
// last two components are the fixed literal suffix. Two different
// (taskId, stepId, attempt) triples can therefore never name the same file. The
// old spelling - "policy-<taskId>-<stepId>-attempt-<n>.json" - allowed a
// collision: task "a-b" with step "c" and task "a" with step "b-c" produced the
// same name, and the converge/refusal logic then compared two unrelated
// artifacts.
func agentRunnerFrozenPolicyPath(runRoot, taskID, stepID string, attempt int) string {
	return filepath.Join(agentRunnerPolicyDirectory(runRoot), fmt.Sprintf("%s.%s.%d.policy.json", taskID, stepID, attempt))
}

// agentRunnerFrozenPolicyAnchorPath is the assembly-time digest anchor of one
// task attempt (F01, P85 fix round 2):
// <runRoot>/agent-runner-policy/<taskId>.<stepId>.<attempt>.anchor.json.
//
// It lives beside the artifact - the same store-local directory, outside the
// managed agent's write root (agentRunnerRejectPolicyWorkspaceOverlap proves the
// disjointness) - and records what the ASSEMBLY froze, so a later read can prove
// the artifact is still that object instead of only proving it agrees with the
// operator grant. The naming proof of agentRunnerFrozenPolicyPath applies
// unchanged: the two literals ".policy.json" and ".anchor.json" are distinct, so
// an anchor can never collide with an artifact.
func agentRunnerFrozenPolicyAnchorPath(runRoot, taskID, stepID string, attempt int) string {
	return filepath.Join(agentRunnerPolicyDirectory(runRoot), fmt.Sprintf("%s.%s.%d.anchor.json", taskID, stepID, attempt))
}

// agentRunnerRejectUnsafePolicyIdentity refuses an identity that cannot name a
// file component. The path builders below place the task and step ids in the file
// name, so a value outside evidence.ValidID (which permits neither "/" nor ".")
// would be a path injection, an ambiguity, or both, and a non-positive attempt
// names no attempt at all. Production reaches these builders only through
// validated request records; the guard makes the property local to the builder
// instead of an assumption about the caller.
func agentRunnerRejectUnsafePolicyIdentity(taskID, stepID string, attempt int) error {
	if !evidence.ValidID(taskID) || !evidence.ValidID(stepID) || attempt < 1 {
		return fmt.Errorf("%w: task %q step %q attempt %d cannot name a frozen policy file", ErrInvalidConfig, taskID, stepID, attempt)
	}
	return nil
}

// agentRunnerFrozenPolicyHash derives the policy digest exactly as preview.go
// does: the canonical encoding of the policy body, hashed with the empty domain.
func agentRunnerFrozenPolicyHash(body previewPolicy) (string, error) {
	canonical, err := evidence.EncodeCanonical(body)
	if err != nil {
		return "", fmt.Errorf("%w: the policy body cannot be canonicalized: %v", ErrInvalidConfig, err)
	}
	return evidence.Digest("", canonical), nil
}

// agentRunnerFrozenPolicyPayload is the exact byte payload of one artifact: the
// canonical encoding plus the document newline.
func agentRunnerFrozenPolicyPayload(artifact agentRunnerFrozenPolicy) ([]byte, error) {
	canonical, err := evidence.EncodeCanonical(artifact)
	if err != nil {
		return nil, fmt.Errorf("%w: the frozen policy artifact cannot be canonicalized: %v", ErrInvalidConfig, err)
	}
	return append(append([]byte(nil), canonical...), '\n'), nil
}

// agentRunnerFrozenPolicyFileDigest is the artifact-external digest of the
// serialized payload. It answers "are these two artifacts byte-identical?" and
// never participates in policyHash.
func agentRunnerFrozenPolicyFileDigest(payload []byte) string {
	return evidence.Digest(agentRunnerFrozenPolicyDomain, payload)
}

// agentRunnerFrozenPolicyAnchor is the assembly-time digest anchor of one task
// attempt.
//
// F01 (P85 fix round 2): the review/publish comparisons only proved that the
// CURRENT artifact agrees with the operator grant. A managed process runs under
// the same OS identity as the assembly, so it can rewrite the artifact, and a
// reread then confirms the rewrite instead of noticing it. PolicyHash records the
// digest the assembly froze and ArtifactFileDigest records the artifact-external
// digest of the exact bytes it published, so the read path can prove the artifact
// is still the object the assembly froze.
//
// Threat-model statement, stated honestly in the code rather than claimed away:
// a hostile process running under the SAME OS identity that writes the artifact
// (and this anchor beside it) is outside the S1 threat model (docs/SECURITY.md,
// the trust boundary section). The anchor is not a sandbox and does not claim to
// be one - it closes the accidental-rewrite and drift-reconciliation gap (a
// resumed assembly, a re-run with a different configuration, an edited run root),
// which is what makes "the policy is frozen" a fact rather than a slogan.
type agentRunnerFrozenPolicyAnchor struct {
	RunID              string `json:"runId"`
	TaskID             string `json:"taskId"`
	StepID             string `json:"stepId"`
	Attempt            int    `json:"attempt"`
	PolicyHash         string `json:"policyHash"`
	ArtifactFileDigest string `json:"artifactFileDigest"`
	// ManagedWriteProof is the assembly-time proof that the identity the managed
	// agent runs under cannot write the policy directory (F01, P85 fix round 3).
	// The anchor is only an integrity record; this digest is the separate fact
	// that the write domain the anchor lives in is closed to the managed agent.
	// It is REQUIRED, and it is not taken on trust: the read path re-observes the
	// domain and requires this exact value (P85 review r1, new Medium), so a
	// well-formed digest from another directory is refused like a malformed one.
	ManagedWriteProof string `json:"managedWriteProof"`
	// ManagedWriteProofShape declares WHICH four-slot domain this proof covers
	// (P85 design rev 4, section 2): agentRunnerPolicyWriteProofShapeCreation is
	// the short domain the creation branch sealed (directory + parent, both file
	// slots empty) and agentRunnerPolicyWriteProofShapeReuse is the full domain the
	// reuse branch sealed (directory + parent + artifact + anchor).
	//
	// The shape is explicit because the digest cannot reveal it: a re-entry that
	// rebuilt "the domain of whatever exists now" would refuse an unchanged
	// creation anchor (the r2 Medium), and the old candidate-set guess was the
	// shortcut §5 removes. An anchor without a recognized shape - including every
	// anchor written before this field existed - is refused instead of guessed at:
	// the S1 production path can never prove a writable-by-this-user policy domain,
	// so no such anchor exists outside test fixtures.
	ManagedWriteProofShape string `json:"managedWriteProofShape"`
}

// agentRunnerFrozenPolicyAnchorPayload is the exact byte payload of one anchor:
// the canonical encoding plus the document newline, the same shape the artifact
// payload has.
func agentRunnerFrozenPolicyAnchorPayload(anchor agentRunnerFrozenPolicyAnchor) ([]byte, error) {
	canonical, err := evidence.EncodeCanonical(anchor)
	if err != nil {
		return nil, fmt.Errorf("%w: the frozen policy anchor cannot be canonicalized: %v", ErrInvalidConfig, err)
	}
	return append(append([]byte(nil), canonical...), '\n'), nil
}

// agentRunnerRejectUnsafePolicyDirectory fails closed when the policy directory
// cannot be proven safe: every already-existing component of the path must be a
// real directory (no symlink, no reparse point) and the directory must not
// overlap the replay root. The component walk is the same primitive the replay
// store applies to its own roots - the artifact is only trusted from a path that
// passes it, on write as well as on read.
func agentRunnerRejectUnsafePolicyDirectory(store *adapters.AgentRunnerReplayStore, directory string) error {
	if store == nil {
		return fmt.Errorf("%w: no replay store is bound, so the frozen policy directory cannot be proven safe", ErrInvalidConfig)
	}
	if err := store.RejectWriterRoots(directory); err != nil {
		return fmt.Errorf("%w: the frozen policy directory is refused: %v", ErrInvalidConfig, err)
	}
	return nil
}

// agentRunnerRejectPolicyWorkspaceOverlap is the executor-cannot-write assertion
// of T2: the managed agent's write root (the isolated workspace) and the frozen
// policy directory may neither equal nor contain one another, and they may not
// alias the same directory object. The assembly asserts it for every task before
// it freezes anything, so a workspace override that points into the policy
// directory is refused instead of freezing policies inside the agent's own
// writable tree.
func agentRunnerRejectPolicyWorkspaceOverlap(policyDirectory, workspaceRoot string) error {
	if strings.TrimSpace(workspaceRoot) == "" {
		return fmt.Errorf("%w: the isolated workspace root is empty, so it cannot be proven disjoint from the frozen policy directory", ErrInvalidConfig)
	}
	if agentRunnerPathContains(policyDirectory, workspaceRoot) ||
		agentRunnerPathContains(workspaceRoot, policyDirectory) ||
		agentRunnerPathAliases(policyDirectory, workspaceRoot) ||
		agentRunnerPathAliases(workspaceRoot, policyDirectory) {
		return fmt.Errorf("%w: the agent write root %q overlaps the frozen policy directory %q", ErrInvalidConfig, workspaceRoot, policyDirectory)
	}
	return nil
}

// agentRunnerPathContains reports whether parent is parent itself or an
// ancestor of child, comparing the way the host filesystem does.
func agentRunnerPathContains(parent, child string) bool {
	parent = filepath.Clean(parent)
	child = filepath.Clean(child)
	if runtime.GOOS == "windows" {
		parent = strings.ToLower(parent)
		child = strings.ToLower(child)
	}
	if parent == child {
		return true
	}
	separator := string(os.PathSeparator)
	if !strings.HasSuffix(parent, separator) {
		parent += separator
	}
	return strings.HasPrefix(child, parent)
}

// agentRunnerPathAliases reports whether first names the same directory object
// as second or as one of second's ancestors. It complements the spelling check,
// which cannot see short-name aliases or bind mounts.
func agentRunnerPathAliases(first, second string) bool {
	firstInfo, err := os.Stat(first)
	if err != nil {
		return false
	}
	for current := filepath.Clean(second); ; {
		if info, statErr := os.Stat(current); statErr == nil && os.SameFile(firstInfo, info) {
			return true
		}
		parent := filepath.Dir(current)
		if parent == current {
			return false
		}
		current = parent
	}
}

// agentRunnerValidateFrozenPolicyBody refuses a policy body that cannot describe
// the run. F05 (P85 fix round 1): the check is no longer "a chain id is present
// and the tasks are covered" but "the body is the policy of THIS assembly
// definition": the chain id must equal the definition's, and the task list must
// be exactly the definition's tasks in definition order. A body that names
// another chain, another task set, a task twice, a task the definition does not
// carry or the definition's tasks in a different order is refused before any
// artifact is published, so a preview of one chain can never be frozen as the
// policy of another run.
//
// The remaining fields (chainProfile, chainDocumentationPolicy, and each entry's
// review/documentationPolicy) have no counterpart on chain.Definition, so they
// cannot be compared field-by-field here. They are still bound: they are part of
// the hashed policy body, and the operator grant names that hash, so a run whose
// previewed policy differs in those fields was never approved.
func agentRunnerValidateFrozenPolicyBody(body previewPolicy, definition chain.Definition) error {
	if strings.TrimSpace(body.ChainID) == "" {
		return fmt.Errorf("%w: the frozen policy body carries no chain id", ErrInvalidConfig)
	}
	if body.ChainID != definition.ID {
		return fmt.Errorf("%w: the frozen policy body names chain %q, but the assembly runs chain %q", ErrInvalidConfig, body.ChainID, definition.ID)
	}
	if len(body.Tasks) != len(definition.Tasks) {
		return fmt.Errorf("%w: the frozen policy body describes %d task(s), but the assembly definition declares %d", ErrInvalidConfig, len(body.Tasks), len(definition.Tasks))
	}
	for index, task := range definition.Tasks {
		if body.Tasks[index].TaskID != task.ID {
			return fmt.Errorf("%w: the frozen policy body describes task %q at position %d, but the assembly definition declares %q", ErrInvalidConfig, body.Tasks[index].TaskID, index, task.ID)
		}
	}
	return nil
}

// agentRunnerFreezeRunPolicy is the pre-dispatch freeze of ExecuteAgentRunnerRun:
// every task's policy artifact is published (and its parent directory synced)
// before the engine exists. Any refusal is fail-closed for the whole run - no
// dispatch, no spawn, no terminal - which is the same shape as a C6 prebuild
// failure.
func agentRunnerFreezeRunPolicy(options AgentRunnerAssemblyOptions, tasks []agentRunnerPlannedTask, store *adapters.AgentRunnerReplayStore) error {
	if options.PolicySnapshot == nil {
		return fmt.Errorf("%w: the AgentRunner assembly requires the run policy the operator previewed and approved; no dispatch happens without a frozen policy", ErrInvalidConfig)
	}
	body := *options.PolicySnapshot
	if err := agentRunnerValidateFrozenPolicyBody(body, options.Definition); err != nil {
		return err
	}
	policyDirectory := agentRunnerPolicyDirectory(options.RunRoot)
	if err := agentRunnerRejectUnsafePolicyDirectory(store, policyDirectory); err != nil {
		return err
	}
	for _, task := range tasks {
		if err := agentRunnerRejectPolicyWorkspaceOverlap(policyDirectory, task.workspaceRoot); err != nil {
			return err
		}
		if _, err := persistAgentRunnerFrozenPolicyProven(store, options.RunRoot, options.RunID, task.binding.TaskID, task.binding.StepID, task.record.Request.Attempt, body, options.managedWriteProof); err != nil {
			return fmt.Errorf("%w: freeze the policy of task %q step %q: %v", ErrInvalidConfig, task.binding.TaskID, task.binding.StepID, err)
		}
	}
	return nil
}

// persistAgentRunnerFrozenPolicyProven publishes one frozen policy artifact and
// returns the artifact as it was read back from disk. Persisting is the
// linearization point: only after the artifact is renamed into place and the
// directory is synced may the run dispatch.
//
// The write is idempotent, never overwriting: an existing artifact of the same
// task attempt converges only when it is content-identical (compared through the
// artifact-external digest of its canonical payload); a divergent artifact is
// refused, so a re-entry with a different configuration cannot rewrite a policy
// that a review already bound.
//
// managedWriteProof is the unexported test seam (formerly the exported
// AgentRunnerAssemblyOptions.ManagedWriteProof): it is ALWAYS empty in production,
// where the managed-identity guard is probed for real, and only an in-package test
// assembly injects a digest to reach the freeze/reuse paths on this platform.
// Because the field is no longer exported, an out-of-package caller has no way to
// substitute a well-formed digest at all (P85 F01, design rev 3 section 1). The
// proof is recorded into the anchor of THIS assembly's policy directory in both
// cases, so the directory binding is unchanged.
func persistAgentRunnerFrozenPolicyProven(store *adapters.AgentRunnerReplayStore, runRoot, runID, taskID, stepID string, attempt int, body previewPolicy, managedWriteProof string) (agentRunnerFrozenPolicy, error) {
	if err := agentRunnerRejectUnsafePolicyIdentity(taskID, stepID, attempt); err != nil {
		return agentRunnerFrozenPolicy{}, err
	}
	policyDirectory := agentRunnerPolicyDirectory(runRoot)
	if err := agentRunnerRejectUnsafePolicyDirectory(store, policyDirectory); err != nil {
		return agentRunnerFrozenPolicy{}, err
	}
	policyHash, err := agentRunnerFrozenPolicyHash(body)
	if err != nil {
		return agentRunnerFrozenPolicy{}, err
	}
	wanted := agentRunnerFrozenPolicy{
		RunID:      runID,
		TaskID:     taskID,
		StepID:     stepID,
		Attempt:    attempt,
		PolicyBody: body,
		PolicyHash: policyHash,
	}
	payload, err := agentRunnerFrozenPolicyPayload(wanted)
	if err != nil {
		return agentRunnerFrozenPolicy{}, err
	}
	path := agentRunnerFrozenPolicyPath(runRoot, taskID, stepID, attempt)
	anchorPath := agentRunnerFrozenPolicyAnchorPath(runRoot, taskID, stepID, attempt)
	// F10 (P85 fix round 2): on the existing-artifact branch the directory has to
	// be proven durable BEFORE the artifact is read, because an earlier interrupted
	// freeze may have published a directory entry that is not yet persistent. The
	// check tolerates a missing directory (the creation branch below proves the new
	// one instead) and refuses on anything else.
	if err := agentRunnerConfirmExistingPolicyDirectory(policyDirectory); err != nil {
		return agentRunnerFrozenPolicy{}, err
	}
	// Which branch this entry takes is only a decision about which domain to seal;
	// the chain below re-reads both names through handles and refuses when the
	// decision and the bound objects disagree (P85 design rev 4, section 1).
	artifactExists, err := agentRunnerPolicyFileExists(path)
	if err != nil {
		return agentRunnerFrozenPolicy{}, err
	}
	if artifactExists {
		chain, err := bindAgentRunnerPolicyDomain(runRoot, taskID, stepID, attempt)
		if err != nil {
			return agentRunnerFrozenPolicy{}, err
		}
		defer func() { _ = chain.Close() }()
		if chain.artifact == nil {
			return agentRunnerFrozenPolicy{}, fmt.Errorf("%w: the frozen policy artifact %s changed between the write-domain observation and its read", ErrInvalidConfig, path)
		}
		if chain.anchor == nil {
			return agentRunnerFrozenPolicy{}, fmt.Errorf("%w: the frozen policy anchor %s is missing, so the artifact cannot be proven to be the object the assembly froze", ErrInvalidConfig, anchorPath)
		}
		// P85 design rev 4, section 1: the artifact and the anchor are read through
		// the chain's bound handles, so the bytes come from the objects the write
		// domain will judge - not from names resolved again.
		rawArtifact, err := chain.readBound(chain.artifact, path, "frozen policy artifact")
		if err != nil {
			return agentRunnerFrozenPolicy{}, err
		}
		rawAnchor, err := chain.readBound(chain.anchor, anchorPath, "frozen policy anchor")
		if err != nil {
			return agentRunnerFrozenPolicy{}, err
		}
		var onDisk agentRunnerFrozenPolicyAnchor
		if err := evidence.DecodeStrictJSON(rawAnchor, &onDisk); err != nil {
			return agentRunnerFrozenPolicy{}, fmt.Errorf("%w: the frozen policy anchor %s is malformed: %v", ErrInvalidConfig, anchorPath, err)
		}
		if err := agentRunnerValidateAnchorShape(onDisk); err != nil {
			return agentRunnerFrozenPolicy{}, err
		}
		// P85 design rev 4, section 2: the sealed domain is rebuilt from the SHAPE the
		// anchor declares, so an unchanged creation anchor is not refused merely
		// because the artifact and its anchor exist now.
		sealed, err := chain.sealedDomainFiles(onDisk.ManagedWriteProofShape)
		if err != nil {
			return agentRunnerFrozenPolicy{}, err
		}
		// P85 design rev 4, section 5: the full set of objects that exist NOW is
		// verified unconditionally, before the sealed-domain comparison can decide
		// anything. The digest comparison decides "was this anchor written by this
		// domain", never "is the domain still closed"; the two questions are
		// independent and both must be true.
		if _, err := chain.observeDomain(chain.currentDomainFiles(), managedWriteProof); err != nil {
			return agentRunnerFrozenPolicy{}, err
		}
		proof, err := chain.observeDomain(sealed, managedWriteProof)
		if err != nil {
			return agentRunnerFrozenPolicy{}, err
		}
		// The anchor's own fields must still describe THIS attempt and THIS payload,
		// and its recorded proof must be the fresh observation of the sealed domain.
		// The whole payload is compared byte for byte, so a rewritten shape, proof or
		// payload refuses instead of being adopted.
		expected := onDisk
		expected.RunID = runID
		expected.TaskID = taskID
		expected.StepID = stepID
		expected.Attempt = attempt
		expected.PolicyHash = policyHash
		expected.ArtifactFileDigest = agentRunnerFrozenPolicyFileDigest(payload)
		expected.ManagedWriteProof = proof
		if err := confirmAgentRunnerFrozenPolicyAnchorFrom(chain.anchor, anchorPath, expected); err != nil {
			return agentRunnerFrozenPolicy{}, err
		}
		if err := chain.ReassertIdentities(); err != nil {
			return agentRunnerFrozenPolicy{}, err
		}
		// F19: the reuse comparison is over the RAW bytes on disk, never over a
		// re-encoding of what was decoded. Two different byte strings can decode into
		// the same canonical structure (a differently indented or reordered document
		// that is still strict-decodable), so comparing re-encoded payloads answers
		// "same meaning" while the guarantee wanted here is "same object".
		if err := agentRunnerConvergeFrozenPolicy(path, rawArtifact, payload); err != nil {
			return agentRunnerFrozenPolicy{}, err
		}
		existing, err := decodeAgentRunnerFrozenPolicy(path, rawArtifact)
		if err != nil {
			return agentRunnerFrozenPolicy{}, err
		}
		return existing, nil
	}
	if err := os.MkdirAll(policyDirectory, 0o755); err != nil {
		return agentRunnerFrozenPolicy{}, fmt.Errorf("%w: create the frozen policy directory: %v", ErrInvalidConfig, err)
	}
	// F10 (P85 fix round 2): the new directory entry itself has to be durable, so
	// the PARENT (the run root) is synced after the creation. Only then is the
	// directory re-classified, bound and synced.
	if err := agentRunnerSyncDirectory(runRoot); err != nil {
		return agentRunnerFrozenPolicy{}, err
	}
	// Re-classify the directory now that it exists: a component that only became
	// reachable through the creation must still be proven a real directory.
	if err := agentRunnerRejectUnsafePolicyDirectory(store, policyDirectory); err != nil {
		return agentRunnerFrozenPolicy{}, err
	}
	// F02/F10 + P85 design rev 4 section 1: ONE binding of the directory and its
	// parent. The two writes below, the read-back and every later re-check go
	// through these handles.
	chain, err := bindAgentRunnerPolicyDomain(runRoot, taskID, stepID, attempt)
	if err != nil {
		return agentRunnerFrozenPolicy{}, err
	}
	defer func() { _ = chain.Close() }()
	if chain.artifact != nil {
		return agentRunnerFrozenPolicy{}, fmt.Errorf("%w: the frozen policy artifact %s appeared between the branch decision and the binding", ErrInvalidConfig, path)
	}
	// F10: the sync happens on the descriptor that was just bound, so the object
	// that was proven is the object that is flushed.
	if err := chain.dir.SyncDirectory(); err != nil {
		return agentRunnerFrozenPolicy{}, fmt.Errorf("%w: sync the frozen policy directory %s: %v", ErrInvalidConfig, policyDirectory, err)
	}
	// F01 (P85 fix round 3): the managed-identity write guard is proven here, after
	// the directory exists and before the first artifact is published, and the
	// observation is recorded in the anchor below. The creation branch has no
	// artifact and no anchor, so its sealed domain is the directory and its parent
	// with both file slots empty (P85 design rev 4, section 2).
	proof, err := chain.observeDomain(nil, managedWriteProof)
	if err != nil {
		return agentRunnerFrozenPolicy{}, err
	}
	anchor := agentRunnerFrozenPolicyAnchor{
		RunID:                  runID,
		TaskID:                 taskID,
		StepID:                 stepID,
		Attempt:                attempt,
		PolicyHash:             policyHash,
		ArtifactFileDigest:     agentRunnerFrozenPolicyFileDigest(payload),
		ManagedWriteProof:      proof,
		ManagedWriteProofShape: agentRunnerPolicyWriteProofShapeCreation,
	}
	anchorPayload, err := agentRunnerFrozenPolicyAnchorPayload(anchor)
	if err != nil {
		return agentRunnerFrozenPolicy{}, err
	}
	// F01: the anchor is published BEFORE the artifact. A crash between the two
	// writes leaves an anchor without an artifact, which the next assembly retries
	// and converges; the reverse order would leave an artifact whose freeze cannot
	// be proven, which is unrecoverable without deleting evidence.
	if err := writeAgentRunnerPolicyAnchorFile(chain, anchorPath, anchorPayload); err != nil {
		return agentRunnerFrozenPolicy{}, err
	}
	if err := writeAgentRunnerPolicyFile(chain, path, payload); err != nil {
		return agentRunnerFrozenPolicy{}, err
	}
	// Bind the two files that now exist relative to the SAME directory handle, so
	// the read-back goes through handles instead of names (design rev 4, section 1).
	if err := chain.bindFiles(); err != nil {
		return agentRunnerFrozenPolicy{}, err
	}
	if chain.artifact == nil || chain.anchor == nil {
		return agentRunnerFrozenPolicy{}, fmt.Errorf("%w: the frozen policy artifact %s is missing right after it was published", ErrInvalidConfig, path)
	}
	persistedRaw, err := chain.readBound(chain.artifact, path, "frozen policy artifact")
	if err != nil {
		return agentRunnerFrozenPolicy{}, err
	}
	// F19: the read-back comparison is over raw disk bytes.
	if !bytes.Equal(persistedRaw, payload) {
		return agentRunnerFrozenPolicy{}, fmt.Errorf("%w: the frozen policy artifact %s does not read back as the bytes that were published", ErrInvalidConfig, path)
	}
	persisted, err := decodeAgentRunnerFrozenPolicy(path, persistedRaw)
	if err != nil {
		return agentRunnerFrozenPolicy{}, err
	}
	if err := confirmAgentRunnerFrozenPolicyAnchorFrom(chain.anchor, anchorPath, anchor); err != nil {
		return agentRunnerFrozenPolicy{}, err
	}
	if err := chain.ReassertIdentities(); err != nil {
		return agentRunnerFrozenPolicy{}, err
	}
	return persisted, nil
}

// agentRunnerPolicyDomainChain is the EVIDENCE HANDLE CHAIN of one task attempt
// (P85 design rev 4, section 1).
//
// The write-domain proof, the artifact read, the anchor read and the publication
// all used to resolve the same names independently, so each stage could be
// looking at a different object that happened to carry the same name. The chain
// binds the two directories once and the two files (when they exist) relative to
// that directory handle, and every later step goes through these handles:
// reading uses the BoundFile, publishing uses the BoundDir, and
// ReassertIdentities re-takes each identity from its handle AND requires the name
// to still address it, so a rename-and-replace between two steps refuses.
type agentRunnerPolicyDomainChain struct {
	directory    string
	dir          *evidence.BoundDir
	dirID        evidence.FileID
	parent       *evidence.BoundDir
	parentID     evidence.FileID
	artifactPath string
	anchorPath   string
	artifact     *evidence.BoundFile
	artifactID   string
	anchor       *evidence.BoundFile
	anchorID     string
}

// bindAgentRunnerPolicyDomain binds the directory, its parent and (when they
// exist) the artifact and the anchor, once, and returns the chain. The caller
// owns the handles and must Close the chain.
func bindAgentRunnerPolicyDomain(runRoot, taskID, stepID string, attempt int) (*agentRunnerPolicyDomainChain, error) {
	if err := agentRunnerRejectUnsafePolicyIdentity(taskID, stepID, attempt); err != nil {
		return nil, err
	}
	directory := agentRunnerPolicyDirectory(runRoot)
	dir, dirID, err := evidence.OpenBoundDirVerified(directory)
	if err != nil {
		return nil, fmt.Errorf("%w: the frozen policy directory %s cannot be bound: %v", ErrInvalidConfig, directory, err)
	}
	parentPath := filepath.Dir(directory)
	parent, parentID, err := evidence.OpenBoundDirVerified(parentPath)
	if err != nil {
		_ = dir.Close()
		return nil, fmt.Errorf("%w: the frozen policy run root %s cannot be bound: %v", ErrInvalidConfig, parentPath, err)
	}
	chain := &agentRunnerPolicyDomainChain{
		directory:    directory,
		dir:          dir,
		dirID:        dirID,
		parent:       parent,
		parentID:     parentID,
		artifactPath: agentRunnerFrozenPolicyPath(runRoot, taskID, stepID, attempt),
		anchorPath:   agentRunnerFrozenPolicyAnchorPath(runRoot, taskID, stepID, attempt),
	}
	if err := chain.bindFiles(); err != nil {
		_ = chain.Close()
		return nil, err
	}
	return chain, nil
}

// Close releases every handle the chain owns. It is idempotent per field.
func (chain *agentRunnerPolicyDomainChain) Close() error {
	if chain == nil {
		return nil
	}
	var first error
	for _, bound := range []*evidence.BoundFile{chain.artifact, chain.anchor} {
		if bound != nil {
			if err := bound.Close(); err != nil && first == nil {
				first = err
			}
		}
	}
	chain.artifact = nil
	chain.anchor = nil
	if chain.dir != nil {
		if err := chain.dir.Close(); err != nil && first == nil {
			first = err
		}
	}
	if chain.parent != nil {
		if err := chain.parent.Close(); err != nil && first == nil {
			first = err
		}
	}
	return first
}

// bindFiles opens the artifact and the anchor RELATIVE to the bound directory
// when they exist and records the identity each handle reports. An absent file
// stays nil (the empty slot of the four-slot domain).
func (chain *agentRunnerPolicyDomainChain) bindFiles() error {
	if chain.artifact == nil {
		bound, err := evidence.OpenBoundDirFile(chain.dir, filepath.Base(chain.artifactPath), agentRunnerFrozenPolicyByteBound)
		switch {
		case err == nil:
			identity, idErr := bound.FileID()
			if idErr != nil {
				_ = bound.Close()
				return fmt.Errorf("%w: identify the frozen policy artifact %s: %v", ErrInvalidConfig, chain.artifactPath, idErr)
			}
			chain.artifact = bound
			chain.artifactID = identity
		case errors.Is(err, os.ErrNotExist):
			// Absent: the slot stays empty.
		default:
			return fmt.Errorf("%w: the frozen policy artifact %s cannot be opened: %v", ErrInvalidConfig, chain.artifactPath, err)
		}
	}
	if chain.anchor == nil {
		bound, err := evidence.OpenBoundDirFile(chain.dir, filepath.Base(chain.anchorPath), agentRunnerFrozenPolicyByteBound)
		switch {
		case err == nil:
			identity, idErr := bound.FileID()
			if idErr != nil {
				_ = bound.Close()
				return fmt.Errorf("%w: identify the frozen policy anchor %s: %v", ErrInvalidConfig, chain.anchorPath, idErr)
			}
			chain.anchor = bound
			chain.anchorID = identity
		case errors.Is(err, os.ErrNotExist):
			// Absent: the slot stays empty.
		default:
			return fmt.Errorf("%w: the frozen policy anchor %s cannot be opened: %v", ErrInvalidConfig, chain.anchorPath, err)
		}
	}
	return nil
}

// readBound reads one policy file through its bound handle and returns the raw
// bytes on disk.
func (chain *agentRunnerPolicyDomainChain) readBound(bound *evidence.BoundFile, path, label string) ([]byte, error) {
	if bound == nil {
		return nil, fmt.Errorf("%w: the %s %s carries no bound handle", ErrInvalidConfig, label, path)
	}
	data, truncated, err := bound.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("%w: the %s %s cannot be read: %v", ErrInvalidConfig, label, path, err)
	}
	if truncated {
		return nil, fmt.Errorf("%w: the %s %s exceeds its byte bound", ErrInvalidConfig, label, path)
	}
	return data, nil
}

// currentDomainFiles returns the file slots that exist NOW, as bound handles.
func (chain *agentRunnerPolicyDomainChain) currentDomainFiles() []adapters.ManagedPolicyDomainBoundFile {
	files := make([]adapters.ManagedPolicyDomainBoundFile, 0, 2)
	if chain.artifact != nil {
		files = append(files, adapters.ManagedPolicyDomainBoundFile{Path: chain.artifactPath, Role: "artifact", Bound: chain.artifact})
	}
	if chain.anchor != nil {
		files = append(files, adapters.ManagedPolicyDomainBoundFile{Path: chain.anchorPath, Role: "anchor", Bound: chain.anchor})
	}
	return files
}

// sealedDomainFiles rebuilds the SEALED domain declared by an anchor shape (P85
// design rev 4, section 2). A reuse anchor only makes sense when both files it
// sealed are still bound; anything else is a refusal, never an empty domain.
func (chain *agentRunnerPolicyDomainChain) sealedDomainFiles(shape string) ([]adapters.ManagedPolicyDomainBoundFile, error) {
	switch shape {
	case agentRunnerPolicyWriteProofShapeCreation:
		return nil, nil
	case agentRunnerPolicyWriteProofShapeReuse:
		if chain.artifact == nil || chain.anchor == nil {
			return nil, fmt.Errorf("%w: a %s anchor sealed the artifact and the anchor, but they are no longer both present", ErrInvalidConfig, shape)
		}
		return chain.currentDomainFiles(), nil
	default:
		return nil, fmt.Errorf("%w: the frozen policy anchor declares an unknown write-proof shape %q", ErrInvalidConfig, shape)
	}
}

// observeDomain proves the domain closed to the managed identity. injected is the
// unexported test seam (always empty in production); when it is non-empty it
// replaces ONLY the platform observation and must still be a valid digest.
//
// This is the ASSEMBLY-time proof and it never goes through the test observation
// seam: the production path must judge the platform, and an in-package test
// injects through the same option the production assembly reads (which is empty
// there), never by substituting the observation the load path performs.
func (chain *agentRunnerPolicyDomainChain) observeDomain(files []adapters.ManagedPolicyDomainBoundFile, injected string) (string, error) {
	if injected != "" {
		if !evidence.ValidHash(injected) {
			return "", fmt.Errorf("%w: the injected managed-identity write proof of the frozen policy directory %s is not a valid proof digest", ErrInvalidConfig, chain.directory)
		}
		return injected, nil
	}
	proof, err := adapters.VerifyManagedAgentCannotWritePolicyDomainHandles(chain.dir, chain.dirID, chain.parent, chain.parentID, files)
	if err != nil {
		return "", fmt.Errorf("%w: the frozen policy domain of %s cannot be proven unwritable by the managed agent identity: %v", ErrInvalidConfig, chain.directory, err)
	}
	if !evidence.ValidHash(proof) {
		return "", fmt.Errorf("%w: the managed-identity write guard of %s produced no valid proof digest", ErrInvalidConfig, chain.directory)
	}
	return proof, nil
}

// ReassertIdentities is the terminal identity check of the chain (P85 design rev
// 4, section 1): every identity is re-taken from its handle AND every name must
// still address that same object. A directory renamed away and replaced by a
// same-named stand-in therefore refuses on the name comparison, and a handle that
// somehow addresses another object refuses on the handle comparison.
func (chain *agentRunnerPolicyDomainChain) ReassertIdentities() error {
	if err := evidence.VerifyBoundDir(chain.dir, chain.dirID); err != nil {
		return fmt.Errorf("%w: the frozen policy directory %s is no longer the approved object: %v", ErrInvalidConfig, chain.directory, err)
	}
	if err := evidence.VerifyBoundDir(chain.parent, chain.parentID); err != nil {
		return fmt.Errorf("%w: the frozen policy run root %s is no longer the approved object: %v", ErrInvalidConfig, filepath.Dir(chain.directory), err)
	}
	if err := agentRunnerReassertNameIdentity(chain.directory, chain.dirID, "frozen policy directory"); err != nil {
		return err
	}
	if err := agentRunnerReassertNameIdentity(filepath.Dir(chain.directory), chain.parentID, "frozen policy run root"); err != nil {
		return err
	}
	for _, file := range []struct {
		bound    *evidence.BoundFile
		recorded string
		path     string
		label    string
	}{
		{chain.artifact, chain.artifactID, chain.artifactPath, "frozen policy artifact"},
		{chain.anchor, chain.anchorID, chain.anchorPath, "frozen policy anchor"},
	} {
		if file.bound == nil {
			continue
		}
		identity, err := file.bound.FileID()
		if err != nil {
			return fmt.Errorf("%w: the %s %s can no longer be identified on its handle: %v", ErrInvalidConfig, file.label, file.path, err)
		}
		if identity != file.recorded {
			return fmt.Errorf("%w: the %s %s handle is not the recorded object (%s != %s)", ErrInvalidConfig, file.label, file.path, identity, file.recorded)
		}
		if err := agentRunnerReassertNameIdentity(file.path, evidence.FileID(file.recorded), file.label); err != nil {
			return err
		}
	}
	return nil
}

// agentRunnerReassertNameIdentity requires the name the caller declared to still
// address the bound object.
func agentRunnerReassertNameIdentity(path string, expected evidence.FileID, label string) error {
	current, err := evidence.FileIDAt(path)
	if err != nil {
		return fmt.Errorf("%w: the %s %s can no longer be observed: %v", ErrInvalidConfig, label, path, err)
	}
	if current != expected {
		return fmt.Errorf("%w: the %s %s no longer names the object that was bound (%s != %s)", ErrInvalidConfig, label, path, current, expected)
	}
	return nil
}

// agentRunnerValidateAnchorShape refuses an anchor whose declared write-proof
// shape is missing or unknown (P85 design rev 4, section 2). The rule is
// fail-closed: an old anchor carries no shape, and "cannot tell" is "refuse".
func agentRunnerValidateAnchorShape(anchor agentRunnerFrozenPolicyAnchor) error {
	switch anchor.ManagedWriteProofShape {
	case agentRunnerPolicyWriteProofShapeCreation, agentRunnerPolicyWriteProofShapeReuse:
		return nil
	default:
		return fmt.Errorf("%w: the frozen policy anchor declares no recognized write-proof shape (%q)", ErrInvalidConfig, anchor.ManagedWriteProofShape)
	}
}

// agentRunnerConvergeFrozenPolicy is the no-replace reuse path: an existing
// artifact of the same task attempt is reused only when the bytes on disk are
// byte-identical to the payload this assembly would publish; anything else is
// refused instead of overwritten.
//
// F19 (P85 fix round 2): the comparison used to run over the RE-ENCODED payload of
// the decoded artifact. Re-encoding normalizes, so a file that had been rewritten
// with different but semantically equal bytes compared as identical and the
// rewrite went unnoticed. Comparing raw bytes keeps the guarantee literal.
func agentRunnerConvergeFrozenPolicy(path string, existingRaw, payload []byte) error {
	if !bytes.Equal(existingRaw, payload) {
		return fmt.Errorf("%w: the frozen policy artifact %s already exists with different content; a frozen policy is never rewritten", ErrInvalidConfig, path)
	}
	return nil
}

// loadAgentRunnerFrozenPolicy reads one artifact from disk and re-derives its
// digest. The artifact must be readable, must decode strictly, must describe the
// requested task attempt, and must match its own recorded policy hash; anything
// else is an error (never a zero value the caller could mistake for approval).
//
// P85 design rev 4, section 1: the directory, its parent and the two files are
// bound ONCE and the bytes come from those handles. The write domain is observed
// and re-asserted BEFORE the artifact is decoded or adopted, so "the object whose
// write protection was judged" and "the object whose bytes were adopted" are the
// same handle set.
func loadAgentRunnerFrozenPolicy(store *adapters.AgentRunnerReplayStore, runRoot, runID, taskID, stepID string, attempt int) (agentRunnerFrozenPolicy, error) {
	if err := agentRunnerRejectUnsafePolicyIdentity(taskID, stepID, attempt); err != nil {
		return agentRunnerFrozenPolicy{}, err
	}
	policyDirectory := agentRunnerPolicyDirectory(runRoot)
	if err := agentRunnerRejectUnsafePolicyDirectory(store, policyDirectory); err != nil {
		return agentRunnerFrozenPolicy{}, err
	}
	path := agentRunnerFrozenPolicyPath(runRoot, taskID, stepID, attempt)
	anchorPath := agentRunnerFrozenPolicyAnchorPath(runRoot, taskID, stepID, attempt)
	chain, err := bindAgentRunnerPolicyDomain(runRoot, taskID, stepID, attempt)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return agentRunnerFrozenPolicy{}, fmt.Errorf("%w: the frozen policy artifact %s is missing", ErrInvalidConfig, path)
		}
		return agentRunnerFrozenPolicy{}, err
	}
	defer func() { _ = chain.Close() }()
	if chain.artifact == nil {
		return agentRunnerFrozenPolicy{}, fmt.Errorf("%w: the frozen policy artifact %s is missing", ErrInvalidConfig, path)
	}
	if chain.anchor == nil {
		return agentRunnerFrozenPolicy{}, fmt.Errorf("%w: the frozen policy anchor %s is missing, so the artifact cannot be proven to be the object the assembly froze", ErrInvalidConfig, anchorPath)
	}
	raw, err := chain.readBound(chain.artifact, path, "frozen policy artifact")
	if err != nil {
		return agentRunnerFrozenPolicy{}, err
	}
	rawAnchor, err := chain.readBound(chain.anchor, anchorPath, "frozen policy anchor")
	if err != nil {
		return agentRunnerFrozenPolicy{}, err
	}
	// The anchor is decoded first because it declares WHICH domain it sealed; the
	// artifact is not decoded until that domain has been re-proven and the
	// identities have been re-asserted.
	var anchor agentRunnerFrozenPolicyAnchor
	if err := evidence.DecodeStrictJSON(rawAnchor, &anchor); err != nil {
		return agentRunnerFrozenPolicy{}, fmt.Errorf("%w: the frozen policy anchor %s is malformed: %v", ErrInvalidConfig, anchorPath, err)
	}
	if err := agentRunnerValidateAnchorShape(anchor); err != nil {
		return agentRunnerFrozenPolicy{}, err
	}
	if anchor.RunID != runID || anchor.TaskID != taskID || anchor.StepID != stepID || anchor.Attempt != attempt {
		return agentRunnerFrozenPolicy{}, fmt.Errorf("%w: the frozen policy anchor is bound to another task attempt", ErrInvalidConfig)
	}
	if err := agentRunnerReverifyManagedWriteProof(chain, anchor); err != nil {
		return agentRunnerFrozenPolicy{}, err
	}
	if err := chain.ReassertIdentities(); err != nil {
		return agentRunnerFrozenPolicy{}, err
	}
	artifact, err := decodeAgentRunnerFrozenPolicy(path, raw)
	if err != nil {
		return agentRunnerFrozenPolicy{}, err
	}
	if artifact.RunID != runID || artifact.TaskID != taskID || artifact.StepID != stepID || artifact.Attempt != attempt {
		return agentRunnerFrozenPolicy{}, fmt.Errorf("%w: the frozen policy artifact %s is bound to another task attempt", ErrInvalidConfig, path)
	}
	if !evidence.ValidHash(artifact.PolicyHash) {
		return agentRunnerFrozenPolicy{}, fmt.Errorf("%w: the frozen policy artifact %s carries no valid policy hash", ErrInvalidConfig, path)
	}
	derived, err := agentRunnerFrozenPolicyHash(artifact.PolicyBody)
	if err != nil {
		return agentRunnerFrozenPolicy{}, err
	}
	if derived != artifact.PolicyHash {
		return agentRunnerFrozenPolicy{}, fmt.Errorf("%w: the frozen policy artifact %s does not match its own policy hash", ErrInvalidConfig, path)
	}
	// F01: the assembly anchor is the fourth equality. It proves the artifact read
	// here is still the object the assembly froze, which self-consistency and the
	// operator grant alone cannot prove.
	if err := agentRunnerVerifyFrozenPolicyAnchor(anchor, artifact, raw); err != nil {
		return agentRunnerFrozenPolicy{}, err
	}
	return artifact, nil
}

// agentRunnerVerifyFrozenPolicyAnchor is the anchor half of the F01 equality
// chain. It requires the anchor to describe the same task attempt, to name the
// same policy hash as the artifact and to carry the artifact-external digest of
// the very bytes that were just read from disk.
//
// The write-domain half of the check is agentRunnerReverifyManagedWriteProof,
// which the load path runs BEFORE the artifact is decoded; this function only
// compares the anchor's own fields with the artifact.
func agentRunnerVerifyFrozenPolicyAnchor(anchor agentRunnerFrozenPolicyAnchor, artifact agentRunnerFrozenPolicy, rawArtifact []byte) error {
	if anchor.RunID != artifact.RunID || anchor.TaskID != artifact.TaskID || anchor.StepID != artifact.StepID || anchor.Attempt != artifact.Attempt {
		return fmt.Errorf("%w: the frozen policy anchor is bound to another task attempt", ErrInvalidConfig)
	}
	if !evidence.ValidHash(anchor.PolicyHash) || anchor.PolicyHash != artifact.PolicyHash {
		return fmt.Errorf("%w: the frozen policy anchor (%s) does not bind the artifact policy hash (%s)", ErrInvalidConfig, anchor.PolicyHash, artifact.PolicyHash)
	}
	if anchor.ArtifactFileDigest != agentRunnerFrozenPolicyFileDigest(rawArtifact) {
		return fmt.Errorf("%w: the frozen policy artifact is not the object the assembly froze", ErrInvalidConfig)
	}
	return nil
}

// agentRunnerObserveManagedWriteProof is the observation seam of the anchor
// reverification. It is a package variable so an in-package test can substitute a
// stub; production never reads it from configuration and never reassigns it.
//
// P85 design rev 4, section 1: the seam is the HANDLE form. Its identities come
// from the chain's own handles, so a stub still replaces the platform observation
// and nothing else - the binding and the terminal identity re-assertion run for
// real.
var agentRunnerObserveManagedWriteProof = func(chain *agentRunnerPolicyDomainChain, files []adapters.ManagedPolicyDomainBoundFile) (string, error) {
	return adapters.VerifyManagedAgentCannotWritePolicyDomainHandles(chain.dir, chain.dirID, chain.parent, chain.parentID, files)
}

// agentRunnerReverifyManagedWriteProof re-proves the frozen-policy write domain
// (P85 design rev 4, sections 2 and 5).
//
// Two independent questions, both of which must be true:
//
//	a. is the domain still closed for every object that exists NOW? This check is
//	   unconditional and runs FIRST: it is never short-circuited by a well-formed
//	   anchor, so the current write rights of the artifact and the anchor are always
//	   judged.
//	b. is this anchor's proof the proof of the domain its SHAPE declares? The
//	   sealed domain is rebuilt from the four-slot shape (creation: both file slots
//	   empty; reuse: artifact + anchor) and its digest must equal the recorded one
//	   byte for byte. An unchanged creation anchor therefore passes, and a proof
//	   from another domain - or a rewritten shape - does not.
//
// The order is (a) then (b) so that a domain that is writable NOW refuses with the
// objects it refused on; the two checks are independent and neither can stand in
// for the other. (Design rev 4 §2.3 lists the digest comparison first; the swap
// changes no verdict, because both must pass on every path.)
func agentRunnerReverifyManagedWriteProof(chain *agentRunnerPolicyDomainChain, anchor agentRunnerFrozenPolicyAnchor) error {
	if !evidence.ValidHash(anchor.ManagedWriteProof) {
		return fmt.Errorf("%w: the frozen policy anchor carries no valid managed-identity write proof", ErrInvalidConfig)
	}
	sealed, err := chain.sealedDomainFiles(anchor.ManagedWriteProofShape)
	if err != nil {
		return err
	}
	if _, err := agentRunnerObserveManagedWriteProof(chain, chain.currentDomainFiles()); err != nil {
		return fmt.Errorf("%w: the current frozen policy domain %s is not closed to the managed agent identity: %v", ErrInvalidConfig, chain.directory, err)
	}
	observed, err := agentRunnerObserveManagedWriteProof(chain, sealed)
	if err != nil {
		return fmt.Errorf("%w: the managed-identity write proof of the frozen policy domain %s cannot be re-observed: %v", ErrInvalidConfig, chain.directory, err)
	}
	if observed != anchor.ManagedWriteProof {
		return fmt.Errorf("%w: the frozen policy anchor's managed-identity write proof is not the proof of the domain %s it is stored in", ErrInvalidConfig, chain.directory)
	}
	return nil
}

// confirmAgentRunnerFrozenPolicyAnchorFrom requires the anchor on disk - read
// through its bound handle - to be the one this assembly would publish. A
// divergent anchor is a refusal: the artifact then cannot be proven to be the
// frozen object.
func confirmAgentRunnerFrozenPolicyAnchorFrom(bound *evidence.BoundFile, anchorPath string, expected agentRunnerFrozenPolicyAnchor) error {
	if bound == nil {
		return fmt.Errorf("%w: the frozen policy anchor %s is missing, so the artifact cannot be proven to be the object the assembly froze", ErrInvalidConfig, anchorPath)
	}
	raw, truncated, err := bound.ReadAll()
	if err != nil {
		return fmt.Errorf("%w: the frozen policy anchor %s cannot be read: %v", ErrInvalidConfig, anchorPath, err)
	}
	if truncated {
		return fmt.Errorf("%w: the frozen policy anchor %s exceeds its byte bound", ErrInvalidConfig, anchorPath)
	}
	wanted, err := agentRunnerFrozenPolicyAnchorPayload(expected)
	if err != nil {
		return err
	}
	if !bytes.Equal(raw, wanted) {
		return fmt.Errorf("%w: the frozen policy anchor %s already exists with different content; an assembly anchor is never rewritten", ErrInvalidConfig, anchorPath)
	}
	return nil
}

// decodeAgentRunnerFrozenPolicy strictly decodes one artifact payload.
func decodeAgentRunnerFrozenPolicy(path string, data []byte) (agentRunnerFrozenPolicy, error) {
	var artifact agentRunnerFrozenPolicy
	if err := evidence.DecodeStrictJSON(data, &artifact); err != nil {
		return agentRunnerFrozenPolicy{}, fmt.Errorf("%w: the frozen policy artifact %s is malformed: %v", ErrInvalidConfig, path, err)
	}
	return artifact, nil
}

// writeAgentRunnerPolicyFile publishes one artifact payload THROUGH the chain's
// bound directory handle (P85 design rev 4, section 1), so the publication
// happens in the same object the write-domain proof and the reads address.
func writeAgentRunnerPolicyFile(chain *agentRunnerPolicyDomainChain, path string, payload []byte) error {
	return publishAgentRunnerPolicyFileFrom(chain.dir, chain.dirID, path, payload, "frozen policy artifact")
}

// writeAgentRunnerPolicyAnchorFile publishes one assembly-anchor payload with the
// same no-replace discipline as the artifact, through the same chain.
func writeAgentRunnerPolicyAnchorFile(chain *agentRunnerPolicyDomainChain, path string, payload []byte) error {
	return publishAgentRunnerPolicyFileFrom(chain.dir, chain.dirID, path, payload, "frozen policy anchor")
}

// publishAgentRunnerPolicyFileFrom publishes one payload inside an ALREADY BOUND
// directory handle (P85 design rev 4, section 1). It binds nothing: the identity
// the caller recorded at bind time is the identity asserted before and after the
// platform's publication steps, so the object that was proven is the object that
// receives the entry.
func publishAgentRunnerPolicyFileFrom(bound *evidence.BoundDir, fileID evidence.FileID, path string, payload []byte, label string) error {
	if bound == nil {
		return fmt.Errorf("%w: the %s directory handle is missing", ErrInvalidConfig, label)
	}
	if err := evidence.VerifyBoundDir(bound, fileID); err != nil {
		return fmt.Errorf("%w: the %s directory is no longer the approved object: %v", ErrInvalidConfig, label, err)
	}
	return publishAgentRunnerPolicyFileIn(bound, fileID, filepath.Dir(path), filepath.Base(path), payload, label)
}

// agentRunnerProvePolicyWriteProtection was the name-based assembly-time proof of
// F01. P85 design rev 4, section 1 moved it onto the evidence handle chain
// (agentRunnerPolicyDomainChain.observeDomain), so the domain that is judged is
// the very object the bytes are read from and written to. The name-based
// primitive adapters.VerifyManagedAgentCannotWritePolicyDomain remains available
// for call sites that hold no handles.

// agentRunnerConfirmExistingPolicyDirectory is the reuse-branch form of the F10
// durability proof: it binds, re-confirms and syncs the directory when it already
// exists, and reports nothing when it does not (the creation branch proves the new
// directory instead). Every other failure is a refusal, so a directory whose
// durability or identity cannot be proven is never read from.
func agentRunnerConfirmExistingPolicyDirectory(directory string) error {
	bound, fileID, err := evidence.OpenBoundDirVerified(directory)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("%w: the frozen policy directory %s cannot be bound: %v", ErrInvalidConfig, directory, err)
	}
	defer func() { _ = bound.Close() }()
	if err := evidence.VerifyBoundDir(bound, fileID); err != nil {
		return fmt.Errorf("%w: the frozen policy directory %s cannot be confirmed: %v", ErrInvalidConfig, directory, err)
	}
	return bound.SyncDirectory()
}

// agentRunnerPolicyFileExists reports whether one frozen policy file is present,
// without opening or reading a single byte of it. F01 (design rev 3, section 1)
// needs this to know WHICH files the reuse branch must observe, and the read that
// follows must agree with the answer or the entry is refused.
func agentRunnerPolicyFileExists(path string) (bool, error) {
	_, err := os.Lstat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, fmt.Errorf("%w: the frozen policy file %s cannot be observed: %v", ErrInvalidConfig, path, err)
}

// agentRunnerSyncDirectory makes one directory entry durable.
//
// Windows cannot sync a directory handle through Go, exactly like the applier's
// atomic write (internal/applier writes its artifacts the same way), so the
// rename/link itself is the durability point there and this is a documented
// platform exemption rather than a silent skip: the parent sync that F10 asks for
// after creating agent-runner-policy/ cannot be expressed on that platform at all.
func agentRunnerSyncDirectory(path string) error {
	directory, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("%w: open the directory %s for sync: %v", ErrInvalidConfig, path, err)
	}
	defer directory.Close()
	if runtime.GOOS == "windows" {
		return nil
	}
	if err := directory.Sync(); err != nil {
		return fmt.Errorf("%w: sync the directory %s: %v", ErrInvalidConfig, path, err)
	}
	return nil
}
