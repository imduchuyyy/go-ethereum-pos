package coin98pos

import (
  "github.com/ethereum/go-ethereum/consensus"
)

type API struct {
  chain consensus.ChainHeaderReader
  coin98pos *Coin98Pos
}
