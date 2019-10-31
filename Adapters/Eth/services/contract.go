package services

import (
	"context"
	"golang.org/x/crypto/sha3"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
)

type Erc20MethodName string

const (
	balanceOfMethodName Erc20MethodName = "balanceOf(address)"
	decimalsMethodName  Erc20MethodName = "decimals()"
	transferMethodName  Erc20MethodName = "transfer(address,uint256)"
)

type IDecimalsContractProvider interface {
	Decimals(ctx context.Context, contract string) (int64, error)
}

type Erc20ContractProvider struct {
	ethClient          *ethclient.Client
	transferMethodID   []byte
	balanceOfMethodID  []byte
	decimalsOfMethodID []byte
}

func NewContractProvider(ethClient *ethclient.Client) *Erc20ContractProvider {
	return &Erc20ContractProvider{
		ethClient:          ethClient,
		transferMethodID:   createMethodID(transferMethodName),
		balanceOfMethodID:  createMethodID(balanceOfMethodName),
		decimalsOfMethodID: createMethodID(decimalsMethodName),
	}
}

func (pr *Erc20ContractProvider) BalanceOf(ctx context.Context, address string, contract string) (*big.Int, error) {
	log := logger.FromContext(ctx)
	accountAddress := common.HexToAddress(address)
	paddedAddress := common.LeftPadBytes(accountAddress.Bytes(), 32)
	var data []byte
	data = append(data, pr.balanceOfMethodID...)
	data = append(data, paddedAddress...)
	contractsAddr := common.HexToAddress(contract)
	msg := ethereum.CallMsg{
		To:   &contractsAddr,
		Data: data,
	}
	bytes, err := pr.ethClient.CallContract(ctx, msg, nil)
	if err != nil {
		log.Error("can't get token's balance (contract %s) for address %s: %s", contract, address, err)
		return nil, err
	}
	return new(big.Int).SetBytes(bytes), nil
}

func (pr *Erc20ContractProvider) Decimals(ctx context.Context, contract string) (int64, error) {
	log := logger.FromContext(ctx)
	contractsAddr := common.HexToAddress(contract)
	msg := ethereum.CallMsg{
		To:   &contractsAddr,
		Data: pr.decimalsOfMethodID,
	}
	bytes, err := pr.ethClient.CallContract(ctx, msg, nil)
	if err != nil {
		log.Error("can't get decimals for contract %s: %s", contract, err)
		return 0, err
	}
	decimals := new(big.Int).SetBytes(bytes)
	log.Infof("contract %s decimals = %d", contract, decimals.Int64())
	return decimals.Int64(), nil
}

func (pr *Erc20ContractProvider) CreateTransferTokenData(ctx context.Context, addressTo string, amount *big.Int) []byte {
	recipient := common.HexToAddress(addressTo)
	methodID := pr.transferMethodID

	// add zeros before address and amount value
	paddedAddress := common.LeftPadBytes(recipient.Bytes(), 32)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)
	return data
}

func createMethodID(methodName Erc20MethodName) []byte {
	fnSignature := []byte(methodName)
	hash := sha3.NewLegacyKeccak256()
	hash.Write(fnSignature)
	methodID := hash.Sum(nil)[:4]
	return methodID
}
