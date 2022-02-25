package contract

import (
    "strings"
    "fmt"

    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/accounts/abi"
    "github.com/ethereum/go-ethereum/common"
)

const EnablePayGasABI = `
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "contractAddress",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "bytes",
          "name": "method",
          "type": "bytes"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "payer",
          "type": "address"
        }
      ],
      "name": "EnablePayGas",
      "type": "event"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_contract",
          "type": "address"
        },
        {
          "internalType": "bytes",
          "name": "_method",
          "type": "bytes"
        }
      ],
      "name": "enable",
      "outputs": [],
      "stateMutability": "payable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_contract",
          "type": "address"
        },
        {
          "internalType": "bytes",
          "name": "_method",
          "type": "bytes"
        }
      ],
      "name": "isEnable",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    }
`

const EnablePayGasBin = `0x`

type EnablePayGas struct {
    EnablePayGasCaller
}


type EnablePayGasCaller struct {
    contract *bind.BoundContract
}

type EnablePayGasCallerRaw struct {
    Contract *EnablePayGasCaller 
}

func (_enablePayGas *EnablePayGasCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
    return _enablePayGas.Contract.contract.Call(opts, result, method, params...)
}

func newEnablePayGasCaller(address common.Address, caller bind.ContractCaller) (*EnablePayGasCaller, error) {
    contract, err := bindEnableGas(address, caller, nil, nil) 
    if err != nil {
        return nil, err
    }

    return &EnablePayGasCaller{contract: contract}, nil
}

func bindEnableGas (address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
    parsed, err := abi.JSON(strings.NewReader(EnablePayGasABI))

    if err != nil {
        return nil, err
    }

    return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_enablePayGas *EnablePayGasCaller) isEnable(opts *bind.CallOpts, method *[]byte, contract common.Address) (bool, error) {
    var out []interface{}

    err := _enablePayGas.contract.Call(opts, &out, "isEnable", method, contract)

    fmt.Printf("out", out)

    out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

    return out0, err
}
