package coin98pos

import (
  "github.com/ethereum/go-ethereum/consensus"
  "github.com/ethereum/go-ethereum/rpc"
)

type API struct {
  chain consensus.ChainHeaderReader
  coin98pos *Coin98Pos
}

// GetSnapshot retrieves the state snapshot at a given block.
func (api *API) GetSnapshot(number *rpc.BlockNumber) error {
  return nil
}
