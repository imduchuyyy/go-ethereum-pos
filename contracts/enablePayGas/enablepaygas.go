package enablepaygas

import (
    "context"
    "math/big"

    "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/contracts/enablePayGas/contracts"
)

//FIXME: please use copyState for this function
// CallContractWithState executes a contract call at the given state.

func CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
    return nil, nil
}

func CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
    return nil, nil
}

