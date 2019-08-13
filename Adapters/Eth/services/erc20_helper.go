package services

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
	"math/big"
	"strings"
)

var erc20TokenABI *abi.ABI

func init() {
	ctx := context.Background()
	log := logger.FromContext(ctx)
	erc20TokenABIValue, err := abi.JSON(strings.NewReader(erc20TokenABIStr))
	if err != nil {
		log.Error("cant parse ERC20 abi")
	}

	erc20TokenABI = &erc20TokenABIValue
}

// ERC20TransferData returns data for a "transfer" tx
func ERC20TransferData(to common.Address, value *big.Int) ([]byte, error) {
	return erc20TokenABI.Pack("transfer", to, value)
}

// ERC20TransferParams is a params for token.transfer method
type ERC20TransferParams struct {
	To    common.Address
	From  common.Address
	Value *big.Int
}

// ParseERC20TransferParams returns "transfer" params from a tx data (reverse of ERC20TransferData)
func ParseERC20TransferParams(data []byte) (*ERC20TransferParams, error) {
	transferParams := &ERC20TransferParams{}

	// ABI pack and unpack methods are not symmetric. The pack of a function call looks like
	// 4 bytes of "method name" and params; so we get a method name first. And then unpack params
	if method, err := erc20TokenABI.MethodById(data); err != nil {
		return nil, err
	} else if err = method.Inputs.Unpack(transferParams, data[4:]); err != nil {
		return nil, err
	}

	return transferParams, nil
}

func CheckERC20Transfers(data []byte) (bool, error) {
	if len(data) < 4 {
		return false, nil
	}
	method, err := erc20TokenABI.MethodById(data)
	if err != nil {
		return false, err
	}
	if method.Name == "transfer" {
		return true, nil
	}
	return false, nil
}

// TransferEvent is used by abi.unpack
type TransferEvent struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}

func (e *TransferEvent) String() string {
	return fmt.Sprintf("from, to, value: %s, %s, %s",
		e.From.String(), e.To.String(), e.Value.String())
}

// ParseTransferEvent returns a transferEvent by a log
func ParseTransferEvent(log types.Log) (*TransferEvent, error) {
	event, err := parseBaseTransferEvent(log.Data)
	if err != nil {
		return nil, err
	}

	event.From = common.BytesToAddress(log.Topics[1].Bytes())
	event.To = common.BytesToAddress(log.Topics[2].Bytes())

	return event, nil
}

func parseBaseTransferEvent(data []byte) (*TransferEvent, error) {
	transferEvent := &TransferEvent{}
	err := erc20TokenABI.Unpack(transferEvent, "Transfer", data)
	if err != nil {
		return nil, err
	}

	return transferEvent, nil
}

// erc20TokenABIStr is the input ABI used to generate the binding from.
const erc20TokenABIStr = `
[
   {
      "constant":true,
      "inputs":[

      ],
      "name":"name",
      "outputs":[
	 {
	    "name":"",
	    "type":"string"
	 }
      ],
      "type":"function"
   },
   {
      "constant":false,
      "inputs":[
	 {
	    "name":"_from",
	    "type":"address"
	 },
	 {
	    "name":"_to",
	    "type":"address"
	 },
	 {
	    "name":"_value",
	    "type":"uint256"
	 }
      ],
      "name":"transferFrom",
      "outputs":[
	 {
	    "name":"success",
	    "type":"bool"
	 }
      ],
      "type":"function"
   },
   {
      "constant":true,
      "inputs":[

      ],
      "name":"decimals",
      "outputs":[
	 {
	    "name":"",
	    "type":"uint8"
	 }
      ],
      "type":"function"
   },
   {
      "constant":true,
      "inputs":[
	 {
	    "name":"",
	    "type":"address"
	 }
      ],
      "name":"balanceOf",
      "outputs":[
	 {
	    "name":"",
	    "type":"uint256"
	 }
      ],
      "type":"function"
   },
   {
      "constant":true,
      "inputs":[

      ],
      "name":"symbol",
      "outputs":[
	 {
	    "name":"",
	    "type":"string"
	 }
      ],
      "type":"function"
   },
   {
      "constant":false,
      "inputs":[
	 {
	    "name":"_to",
	    "type":"address"
	 },
	 {
	    "name":"_value",
	    "type":"uint256"
	 }
      ],
      "name":"transfer",
      "outputs":[

      ],
      "type":"function"
   },
   {
      "constant":false,
      "inputs":[
	 {
	    "name":"_spender",
	    "type":"address"
	 },
	 {
	    "name":"_value",
	    "type":"uint256"
	 },
	 {
	    "name":"_extraData",
	    "type":"bytes"
	 }
      ],
      "name":"approveAndCall",
      "outputs":[
	 {
	    "name":"success",
	    "type":"bool"
	 }
      ],
      "type":"function"
   },
   {
      "constant":true,
      "inputs":[
	 {
	    "name":"",
	    "type":"address"
	 },
	 {
	    "name":"",
	    "type":"address"
	 }
      ],
      "name":"spentAllowance",
      "outputs":[
	 {
	    "name":"",
	    "type":"uint256"
	 }
      ],
      "type":"function"
   },
   {
      "constant":true,
      "inputs":[
	 {
	    "name":"",
	    "type":"address"
	 },
	 {
	    "name":"",
	    "type":"address"
	 }
      ],
      "name":"allowance",
      "outputs":[
	 {
	    "name":"",
	    "type":"uint256"
	 }
      ],
      "type":"function"
   },
   {
      "inputs":[
	 {
	    "name":"initialSupply",
	    "type":"uint256"
	 },
	 {
	    "name":"tokenName",
	    "type":"string"
	 },
	 {
	    "name":"decimalUnits",
	    "type":"uint8"
	 },
	 {
	    "name":"tokenSymbol",
	    "type":"string"
	 }
      ],
      "type":"constructor"
   },
   {
      "anonymous":false,
      "inputs":[
	 {
	    "indexed":true,
	    "name":"from",
	    "type":"address"
	 },
	 {
	    "indexed":true,
	    "name":"to",
	    "type":"address"
	 },
	 {
	    "indexed":false,
	    "name":"value",
	    "type":"uint256"
	 }
      ],
      "name":"Transfer",
      "type":"event"
   }
]
`
