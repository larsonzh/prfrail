//go:build windows && b3bnative && b3bvariant

package adapters

// b3bVariantDegrade is true in the variant build. The arm exists to show that the rig can
// falsify itself: every candidate silently degrades to the baseline primitive (no
// WRITE_THROUGH on the publish, no directory flush afterwards) while still claiming its own
// arm name. The trace deliberately keeps naming what actually ran (os.Link /
// platform-parent-sync), so the round's mechanism evidence must come out not-ok and the cell
// must not be able to reach proven. If a cell under this arm were to report mechanism-ok, the
// mechanism check would be blind and the durability claim would rest on nothing.
//
// Build: go test -c -tags "b3bnative b3bvariant" -o tmp\b3b\writer-variant.exe .\internal\adapters\
const b3bVariantDegrade = true
