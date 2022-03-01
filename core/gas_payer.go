package core

import (
    "math/big"

    "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/ethereum/go-ethereum/core/state"
)

var (
  slotEnableGas = map[string]uint64 {
    "isEnable": 0,
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


func isContractEnablePayGas(copyState *state.StateDB, contractAddr common.Address, method []byte) bool {
    enableSlot := slotEnableGas["isEnable"]
    slotHash := common.BigToHash(new(big.Int).SetUint64(enableSlot))
    enableKeyInMapping := crypto.Keccak256(method, crypto.Keccak256(contractAddr.Hash().Bytes(), slotHash.Bytes()))

    //enableKey := state.GetLocMappingAtKey(enableKeyInMapping, enableSlot)

    isEnable := copyState.GetState(common.EnablePayGas, common.BytesToHash(enableKeyInMapping))
    if isEnable.Big().Cmp(new(big.Int).SetUint64(0)) > 0 {
        return true
    }
    return false
}

func IsContractEnablePayGas(copyState *state.StateDB, contractAddr common.Address, method []byte) bool {
    return isContractEnablePayGas(copyState, contractAddr, method)
}
