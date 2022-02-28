package core

import (
    "math/big"
    "fmt"

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
    fmt.Printf("Contract address", contractAddr, "\n")
    fmt.Printf("Method", method, "\n")

    enableSlot := slotEnableGas["isEnable"]
    slotHash := common.BigToHash(new(big.Int).SetUint64(enableSlot))
    fmt.Printf("slot hash", slotHash, "\n")
    enableKeyInMapping := crypto.Keccak256Hash(method, crypto.Keccak256Hash(contractAddr.Bytes(), slotHash.Bytes()).Bytes())
    fmt.Printf("Key", enableKeyInMapping, "\n")

    //enableKey := state.GetLocMappingAtKey(enableKeyInMapping, enableSlot)

    isEnable := copyState.GetState(common.EnablePayGas, enableKeyInMapping)
    fmt.Printf("is Enable", isEnable, "\n")
    return true
}

