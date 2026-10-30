package tickets

// Commit stages reported to DurableCostLedgerConfig.CommitFault. The fault is
// consulted at the boundary of the named stage; returning an error there is
// treated as that stage failing, exactly like the corresponding syscall would.
//
// This surface exists for tests only: production callers leave CommitFault nil,
// which means "never inject" and keeps the durable sequence untouched.
const (
	StageWrite      = "write"
	StageSync       = "sync"
	StageClose      = "close"
	StageRename     = "rename"
	StageDirBarrier = "dir-barrier"
	StageDirClose   = "dir-close"
)

// durableLedgerBytesDomain domain-separates the raw-byte digests used to
// classify an inconclusive atomic replace.
const durableLedgerBytesDomain = "proofrail:durable-cost-ledger-bytes:1\n"
