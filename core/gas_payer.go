package core

import (
  "fmt"
  "math/big"

  "github.com/ethereum/go-ethereum/common"
  "github.com/ethereum/go-ethereum/accounts/abi"
  "github.com/ethereum/go-ethereum/core/state"
  "github.com/ethereum/go-ethereum/crypto"
)

const (
  slotEnableGas = map[string]uint64 {
    "isEnable": 0
  }
)

type callmsg struct {
  ethereum.CallMsg
}

func (m callmsg) From() common.Address      { return m.CallMsg.From }
func (m callmsg) Nonce() uint64             { return 0 }
func (m callmsg) CheckNonce() bool          { return false }
func (m callmsg) To() *common.Address       { return m.CallMsg.To }
func (m callmsg) GasPrice() *big.Int        { return m.CallMsg.GasPrice }
func (m callmsg) Gas() uint64               { return m.CallMsg.Gas }
func (m callmsg) Value() *big.Int           { return m.CallMsg.Value }
func (m callmsg) Data() []byte              { return m.CallMsg.Data }
func (m callmsg) BalanceTokenFee() *big.Int { return m.CallMsg.BalanceTokenFee }

func isContractEnablePayGas(chain consensus.ChainContext, blockNumber *big.Int, copyState *state.StateDB, contractAddr common.Address, method byte[]) (bool) {
    enableKeyInMapping := crypto.Keccak256(method, contractAddr.Bytes())
    enableSlot := slotEnableGas["isEnable"]

    enableKey := GetLocMappingAtKey(enableKeyInMapping, enableSlot)

    isEnable := copyState.GetState(EnablePayGas, commom.BigToHash(enableKey))

    return isEnable
}
