package adapters

// P85 R7 (review r6 finding): an explicit role that contradicts the contract's name
// suffix describes two different slots for one object, and a suffix spelled in a
// different case is not the spelling the contract freezes. Both are refused instead
// of being resolved by precedence, so a caller cannot place the anchor in the
// artifact slot by declaring it as one.

import (
	"errors"
	"strings"
	"testing"
)

// TestP85R7RoleSuffixContradictionIsRefused pins the role resolution of the
// handle-form file slots at its smallest unit.
func TestP85R7RoleSuffixContradictionIsRefused(t *testing.T) {
	const (
		artifactName = "task.step.attempt.1.policy.json"
		anchorName   = "task.step.attempt.1.anchor.json"
	)
	cases := []struct {
		name     string
		declared string
		path     string
		wantRole int
		wantErr  string
	}{
		{
			name:     "an-artifact-role-on-an-anchor-name",
			declared: managedPolicyDomainArtifactRole, path: anchorName,
			wantRole: -1, wantErr: "declares the artifact role",
		},
		{
			name:     "an-anchor-role-on-an-artifact-name",
			declared: managedPolicyDomainAnchorRole, path: artifactName,
			wantRole: -1, wantErr: "declares the anchor role",
		},
		{
			name:     "an-uppercase-artifact-suffix",
			declared: "", path: strings.ToUpper(artifactName),
			wantRole: -1, wantErr: "lowercase form",
		},
		{
			name:     "an-uppercase-anchor-suffix-with-a-declared-role",
			declared: managedPolicyDomainAnchorRole, path: strings.ToUpper(anchorName),
			wantRole: -1, wantErr: "lowercase form",
		},
		{
			name:     "a-declared-role-without-a-name-marker",
			declared: managedPolicyDomainAnchorRole, path: "declared-name.json",
			wantRole: 1,
		},
		{
			name:     "the-matching-pair",
			declared: managedPolicyDomainArtifactRole, path: artifactName,
			wantRole: 0,
		},
		{
			name:     "the-name-marker-alone",
			declared: "", path: anchorName,
			wantRole: 1,
		},
		{
			name:     "no-role-anywhere",
			declared: "", path: "declared-name.json",
			wantRole: -1,
		},
	}
	for _, testCase := range cases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			role, err := managedPolicyDomainResolveFileRole(testCase.declared, testCase.path)
			if testCase.wantErr == "" {
				if err != nil {
					t.Fatalf("this pair must resolve without a refusal, got %v", err)
				}
				if role != testCase.wantRole {
					t.Fatalf("the role must be %d, got %d", testCase.wantRole, role)
				}
				return
			}
			if err == nil {
				t.Fatalf("this pair must be refused, got role %d", role)
			}
			if !errors.Is(err, ErrManagedAgentPolicyWriteProof) {
				t.Fatalf("the refusal must be the write-proof sentinel, got %v", err)
			}
			if !strings.Contains(err.Error(), testCase.wantErr) {
				t.Fatalf("the refusal must mention %q, got %v", testCase.wantErr, err)
			}
		})
	}
}
