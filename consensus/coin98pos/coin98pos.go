package coin98pos

import (
  "github.com/ethereum/go-ethereum/params"
)

type Coin98Pos struct {
  config *params.Coin98PosConfig
}

func New(config *params.Coin98PosConfig) *Coin98Pos {
  c := &Coin98Pos{
    config: config,
  }

  return c
}
