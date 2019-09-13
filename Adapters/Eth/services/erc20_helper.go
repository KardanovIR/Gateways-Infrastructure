package services

import (
	"context"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/models"
)

var erc20TokenABI *abi.ABI

func init() {
	ctx := context.Background()
	log := logger.FromContext(ctx)
	erc20TokenABIValue, err := abi.JSON(strings.NewReader(erc20TokenABIStr))
	if err != nil {
		log.Error("cant parse ERC20 abi", err)
	}

	erc20TokenABI = &erc20TokenABIValue
}

// ERC20TransferData returns data for a "transfer" tx
func ERC20TransferData(to common.Address, value *big.Int) ([]byte, error) {
	return erc20TokenABI.Pack("transfer", to, value)
}

// ERC20TransferFromData returns data for a "transfer" tx
func ERC20TransferFromData(from common.Address, to common.Address, value *big.Int) ([]byte, error) {
	return erc20TokenABI.Pack("transferFrom", from, to, value)
}

// ERC20ApproveSender returns data for a "approve" tx - give rights to another account
func ERC20ApproveSender(spender common.Address, value *big.Int) ([]byte, error) {
	return erc20TokenABI.Pack("approve", spender, value)
}

// ERC20ApproveSender returns data for a "allowance" tx - get rights of another account
func ERC20AllowanceForSender(owner common.Address, spender common.Address) ([]byte, error) {
	return erc20TokenABI.Pack("allowance", owner, spender)
}

// ParseERC20TransferParams returns "transfer" params from a tx data (reverse of ERC20TransferData)
func ParseERC20TransferParams(data []byte) (*models.TransferEvent, error) {
	transferParams := &models.TransferEvent{}

	// ABI pack and unpack methods are not symmetric. The pack of a function call looks like
	// 4 bytes of "method name" and params; so we get a method name first. And then unpack params
	if method, err := erc20TokenABI.MethodById(data); err != nil {
		return nil, err
	} else if err = method.Inputs.Unpack(transferParams, data[4:]); err != nil {
		return nil, err
	}

	return transferParams, nil
}

func CheckERC20Transfers(data []byte) bool {
	if len(data) < 4 {
		return false
	}
	method, err := erc20TokenABI.MethodById(data)
	if err != nil {
		return false
	}
	return method.Name == "transfer" || method.Name == "transferFrom"
}

// erc20TokenABIStr is the input ABI used to generate the binding from.
const erc20TokenABIStr = `
[
    {
        "constant": true,
        "inputs": [],
        "name": "name",
        "outputs": [
            {
                "name": "",
                "type": "string"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "constant": false,
        "inputs": [
            {
                "name": "_spender",
                "type": "address"
            },
            {
                "name": "_value",
                "type": "uint256"
            }
        ],
        "name": "approve",
        "outputs": [
            {
                "name": "",
                "type": "bool"
            }
        ],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [],
        "name": "totalSupply",
        "outputs": [
            {
                "name": "",
                "type": "uint256"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "constant": false,
        "inputs": [
            {
                "name": "_from",
                "type": "address"
            },
            {
                "name": "_to",
                "type": "address"
            },
            {
                "name": "_value",
                "type": "uint256"
            }
        ],
        "name": "transferFrom",
        "outputs": [
            {
                "name": "",
                "type": "bool"
            }
        ],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [],
        "name": "decimals",
        "outputs": [
            {
                "name": "",
                "type": "uint8"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [
            {
                "name": "_owner",
                "type": "address"
            }
        ],
        "name": "balanceOf",
        "outputs": [
            {
                "name": "balance",
                "type": "uint256"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [],
        "name": "symbol",
        "outputs": [
            {
                "name": "",
                "type": "string"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "constant": false,
        "inputs": [
            {
                "name": "_to",
                "type": "address"
            },
            {
                "name": "_value",
                "type": "uint256"
            }
        ],
        "name": "transfer",
        "outputs": [
            {
                "name": "",
                "type": "bool"
            }
        ],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [
            {
                "name": "_owner",
                "type": "address"
            },
            {
                "name": "_spender",
                "type": "address"
            }
        ],
        "name": "allowance",
        "outputs": [
            {
                "name": "",
                "type": "uint256"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "payable": true,
        "stateMutability": "payable",
        "type": "fallback"
    },
    {
        "anonymous": false,
        "inputs": [
            {
                "indexed": true,
                "name": "owner",
                "type": "address"
            },
            {
                "indexed": true,
                "name": "spender",
                "type": "address"
            },
            {
                "indexed": false,
                "name": "value",
                "type": "uint256"
            }
        ],
        "name": "Approval",
        "type": "event"
    },
    {
        "anonymous": false,
        "inputs": [
            {
                "indexed": true,
                "name": "from",
                "type": "address"
            },
            {
                "indexed": true,
                "name": "to",
                "type": "address"
            },
            {
                "indexed": false,
                "name": "value",
                "type": "uint256"
            }
        ],
        "name": "Transfer",
        "type": "event"
    }
]
`
