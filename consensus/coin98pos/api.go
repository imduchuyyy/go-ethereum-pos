package coin98pos

import (
	"github.com/ethereum/go-ethereum/consensus"
)

// API is a user facing RPC API to allow controlling the signer and voting
// mechanisms of the proof-of-authority scheme.
type API struct {
	chain  consensus.ChainHeaderReader
	coin98Pos *Coin98Pos
}

// GetSnapshot retrieves the state snapshot at a given block.
