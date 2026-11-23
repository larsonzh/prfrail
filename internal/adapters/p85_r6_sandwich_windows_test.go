//go:build windows

package adapters

// P85 R6 / the Windows leg of the identity sandwich (review r3 finding 1).
//
// The Windows observation takes the slot identity from its handle, classifies the
// name chain, then requires the NAME to still address that handle before the
// per-right access probes and to still address it afterwards. The leg moves the name
// onto a different object and requires the observation to refuse with the
// handle-versus-name comparison rather than with a writability verdict - so the
// sandwich, not the ACL table, is what answers.
//
// Mutation check (design rev 4 section 6, mutation 2, locally runnable equivalent):
// deleting the handle-versus-name comparison from
// managedAgentPolicyDomainSlotObservation leaves the observation reporting the
// probes of the stand-in object, and this leg goes red.

import (
	"errors"
	"strings"
	"testing"
)

func init() {
	p85R6SandwichPlatformLegs = append(p85R6SandwichPlatformLegs, p85R6SandwichLeg{
		name: "the-windows-native-observation-refuses-a-name-that-moved",
		run: func(t *testing.T) {
			slot, _, swap := p85R6SandwichFixture(t)
			swap()
			sidText, err := installSourceCurrentUserSIDText()
			if err != nil {
				t.Skipf("this host cannot resolve the executing user's SID: %v", err)
			}
			_, _, observationErr := managedAgentPolicyDomainSlotObservation(slot, sidText)
			if observationErr == nil {
				t.Fatal("the observation must refuse a name that no longer addresses the handle's object")
			}
			if !errors.Is(observationErr, ErrManagedAgentPolicyWriteProof) {
				t.Fatalf("the refusal must be the write-proof sentinel, got %v", observationErr)
			}
			if !strings.Contains(observationErr.Error(), "no longer names the object its handle addresses") {
				t.Fatalf("the refusal must be the handle-versus-name comparison, got %v", observationErr)
			}
		},
	})
}
