package services

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/models"
)

const FailedTxStatus = 0

var AllowanceAmountIsNotEnoughError = errors.New("allowanceAmountIsNotEnough")

func (cl *nodeClient) CreateRawTransaction(ctx context.Context, addressFrom string, addressTo string,
	amount *big.Int, nonce uint64) ([]byte, error) {
	log := logger.FromContext(ctx)
	log.Debugf("call service method 'CreateRawTransaction': send %s from %s to %s", amount, addressFrom, addressTo)
	ok, _, err := cl.IsAddressValid(ctx, addressTo)
	if err != nil {
		return nil, fmt.Errorf("check address %s fails: %s", addressTo, err)
	}
	if !ok {
		return nil, fmt.Errorf("address %s is not valid", addressTo)
	}
	gasPrice, err := cl.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get suggected gas price %s", err)
	}
	log.Debugf("suggest gas price %s", gasPrice)
	if nonce == 0 {
		nonce, err = cl.GetNextNonce(ctx, addressFrom)
		if err != nil {
			return nil, fmt.Errorf("can't get next nonce for address %s: %s", addressFrom, err)
		}
	}
	log.Debugf("nonce will be %d", nonce)
	tx := types.NewTransaction(
		nonce,
		common.HexToAddress(addressTo),
		amount,
		gasLimitForMoneyTransfer,
		gasPrice,
		nil,
	)
	b, err := SerializeTx(tx)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (cl *nodeClient) CreateErc20TokensRawTransaction(ctx context.Context, addressFrom string, contractAddress string,
	addressTo string, amount *big.Int, nonce uint64) ([]byte, error) {
	log := logger.FromContext(ctx)
	log.Debugf("call service method 'CreateErc20TokensRawTransaction': send %s tokens from %s to %s; contract %s",
		amount, addressFrom, addressTo, contractAddress)
	ok, _, err := cl.IsAddressValid(ctx, addressTo)
	if err != nil {
		return nil, fmt.Errorf("check address %s fails: %s", addressTo, err)
	}
	if !ok {
		return nil, fmt.Errorf("address %s is not valid", addressTo)
	}
	gasPrice, err := cl.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get suggected gas price %s", err)
	}
	log.Debugf("suggest gas price %s", gasPrice)
	if nonce == 0 {
		nonce, err = cl.GetNextNonce(ctx, addressFrom)

		if err != nil {
			return nil, fmt.Errorf("can't get next nonce for address %s: %s", addressFrom, err)
		}
	}
	log.Debugf("nonce will be %d", nonce)
	sender := common.HexToAddress(addressFrom)
	recipient := common.HexToAddress(addressTo)
	tokenAddress := common.HexToAddress(contractAddress)

	data, err := ERC20TransferData(recipient, amount)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	gasLimit, err := cl.ethClient.EstimateGas(context.Background(), ethereum.CallMsg{
		From: sender,
		To:   &tokenAddress,
		Data: data,
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debugf("estimated gas limit %d", gasLimit)
	ethAmountZero := big.NewInt(0)
	tx := types.NewTransaction(nonce, tokenAddress, ethAmountZero, gasLimit, gasPrice, data)
	b, err := SerializeTx(tx)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (cl *nodeClient) CreateErc20TokensTransferToTxSender(ctx context.Context, addressFrom string,
	contractAddress string, txSender string, amount *big.Int, nonce uint64) ([]byte, error) {
	log := logger.FromContext(ctx)
	log.Debugf("call service method 'CreateErc20TokensTransferToTxSender': send %s tokens from %s to %s; contract %s",
		amount, addressFrom, txSender, contractAddress)
	ok, _, err := cl.IsAddressValid(ctx, txSender)
	if err != nil {
		return nil, fmt.Errorf("check address %s fails: %s", txSender, err)
	}
	if !ok {
		return nil, fmt.Errorf("address %s is not valid", txSender)
	}
	allowanceAmount, err := cl.GetErc20AllowanceAmount(ctx, addressFrom, contractAddress, txSender)
	if err != nil {
		log.Errorf("can't check allowance amount: %s", err)
		return nil, err
	}
	log.Debugf("allowanceAmount %s", allowanceAmount)
	log.Debugf("amount %s", amount)
	log.Debugf("Cmp %s", allowanceAmount.Cmp(amount))
	if allowanceAmount.Cmp(amount) < 0 {
		log.Errorf("allowance amount is less than transfer's amount")
		return nil, AllowanceAmountIsNotEnoughError
	}
	gasPrice, err := cl.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get suggected gas price %s", err)
	}
	log.Debugf("suggest gas price %s", gasPrice)
	if nonce == 0 {
		nonce, err = cl.GetNextNonce(ctx, txSender)

		if err != nil {
			return nil, fmt.Errorf("can't get next nonce for address %s: %s", txSender, err)
		}
	}
	log.Debugf("nonce will be %d", nonce)

	owner := common.HexToAddress(addressFrom)
	sender := common.HexToAddress(txSender)
	tokenAddress := common.HexToAddress(contractAddress)

	data, err := ERC20TransferFromData(owner, sender, amount)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	gasLimit, err := cl.ethClient.EstimateGas(context.Background(), ethereum.CallMsg{
		From: sender,
		To:   &tokenAddress,
		Data: data,
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debugf("estimated gas limit %d", gasLimit)
	fee := new(big.Int).Mul(big.NewInt(int64(gasLimit)), gasPrice)
	log.Debugf("fee %s", fee)
	ethAmountZero := big.NewInt(0)
	tx := types.NewTransaction(nonce, tokenAddress, ethAmountZero, gasLimit, gasPrice, data)
	b, err := SerializeTx(tx)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (cl *nodeClient) SignTransaction(ctx context.Context, senderAddr string, rlpTx []byte) ([]byte, error) {
	log := logger.FromContext(ctx)
	log.Debug("call service method 'SignTransaction'")
	pk, ok := cl.privateKeys[senderAddr]
	if !ok {
		return nil, fmt.Errorf("can't signed transaction for %s: haven't private key", senderAddr)
	}
	return cl.signTx(ctx, pk, rlpTx)
}

func (cl *nodeClient) Erc20TokensRawApproveTransaction(ctx context.Context, ownerAddress string, contractAddress string,
	amount *big.Int, spenderAddress string) ([]byte, *big.Int, error) {
	log := logger.FromContext(ctx)
	log.Debugf("call service method 'Erc20TokensRawApproveTransaction': approve %s tokens in address %s to %s; contract %s",
		amount, ownerAddress, spenderAddress, contractAddress)

	nonce, err := cl.GetNextNonce(ctx, ownerAddress)

	if err != nil {
		return nil, nil, fmt.Errorf("can't get next nonce for address %s: %s", ownerAddress, err)
	}
	log.Debugf("nonce will be %d", nonce)

	owner := common.HexToAddress(ownerAddress)
	tokenAddress := common.HexToAddress(contractAddress)
	spender := common.HexToAddress(spenderAddress)

	dataForApprove, err := ERC20ApproveSender(spender, amount)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	gasLimit, err := cl.ethClient.EstimateGas(context.Background(), ethereum.CallMsg{
		From: owner,
		To:   &tokenAddress,
		Data: dataForApprove,
	})
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}
	log.Debugf("estimated gas limit %d", gasLimit)
	gasPrice, err := cl.SuggestGasPrice(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("can't get suggected gas price %s", err)
	}
	log.Debugf("suggest gas price %s", gasPrice)

	ethAmountZero := big.NewInt(0)
	fee := new(big.Int).Mul(big.NewInt(int64(gasLimit)), gasPrice)
	log.Debugf("fee %s", fee)
	tx := types.NewTransaction(nonce, tokenAddress, ethAmountZero, gasLimit, gasPrice, dataForApprove)
	b, err := SerializeTx(tx)
	if err != nil {
		return nil, nil, err
	}
	return b, fee, nil
}

func (cl *nodeClient) SignTransactionWithPrivateKey(ctx context.Context, privateKey string, rlpTx []byte) ([]byte, error) {
	log := logger.FromContext(ctx)
	log.Debug("call service method 'SignTransactionWithPrivateKey'")
	pk, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Errorf("can't cast key to ECDSA: %s", err)
		return nil, err
	}
	return cl.signTx(ctx, pk, rlpTx)
}

func (cl *nodeClient) signTx(ctx context.Context, privateKey *ecdsa.PrivateKey, rlpTx []byte) ([]byte, error) {
	tx, err := DeserializeTx(rlpTx)
	if err != nil {
		return nil, fmt.Errorf("can't deserialize tx %s", err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(cl.chainID)), privateKey)
	if err != nil {
		return nil, err
	}
	return SerializeTx(signedTx)
}

func (cl *nodeClient) SendTransaction(ctx context.Context, rlpTx []byte) (txHex string, err error) {
	log := logger.FromContext(ctx)
	log.Debug("call service method 'SendTransaction'")

	tx, err := DeserializeTx(rlpTx)
	if err != nil {
		return "", fmt.Errorf("can't deserialize tx %s", err)
	}
	log.Debugf("try to send transaction %+v", tx)
	return tx.Hash().Hex(), cl.ethClient.SendTransaction(ctx, tx)
}

func (cl *nodeClient) GetTxStatusByTxID(ctx context.Context, txID string) (models.TxStatus, error) {
	log := logger.FromContext(ctx)
	log.Debugf("call service method 'GetTransactionByHash' %s", txID)
	_, status, err := cl.getTxAndStatus(ctx, txID)
	if err != nil {
		return "", err
	}
	return status, nil
}

func (cl *nodeClient) TransactionInfo(ctx context.Context, txID string) (*models.TxInfo, error) {
	log := logger.FromContext(ctx)
	log.Debugf("call service method 'TransactionInfo' %s", txID)
	tx, status, err := cl.getTxAndStatus(ctx, txID)
	if err != nil {
		return nil, err
	}
	if status == models.TxStatusUnKnown {
		return &models.TxInfo{Status: status}, nil
	}
	// sender
	signer := types.NewEIP155Signer(big.NewInt(cl.chainID))
	sender, err := types.Sender(signer, tx)
	if err != nil {
		log.Errorf("can't get sender for tx %s: %s", txID, err)
		return nil, err
	}
	var fee *big.Int
	// used gas and tx status
	if status == models.TxStatusPending {
		fee = new(big.Int).Mul(big.NewInt(int64(tx.Gas())), tx.GasPrice())
	} else {
		receipt, err := cl.ethClient.TransactionReceipt(ctx, common.HexToHash(txID))
		if err != nil {
			log.Errorf("can't TransactionReceipt for tx %s: %s", txID, err)
			return nil, err
		}
		// contract can be failed and be in blockchain
		if receipt.Status == FailedTxStatus {
			status = models.TxStatusFailed
		}
		fee = new(big.Int).Mul(big.NewInt(int64(receipt.GasUsed)), tx.GasPrice())
	}
	txInfo := models.TxInfo{
		From:   sender.String(),
		Amount: tx.Value(),
		Fee:    fee,
		TxHash: tx.Hash().String(),
		Data:   tx.Data(),
		Status: status,
		Nonce:  tx.Nonce(),
	}
	isERC20Transfers, err := CheckERC20Transfers(tx.Data())
	if err != nil {
		return nil, err
	}
	if !isERC20Transfers {
		log.Debugf("there are not erc20 transfers in tx %s: %s", txID, err)
		txInfo.To = tx.To().Hex()
	} else {
		transferParams, err := ParseERC20TransferParams(tx.Data())
		if err != nil {
			return nil, err
		}
		txInfo.To = transferParams.To.Hex()
		txInfo.AssetAmount = transferParams.Value
		txInfo.Contract = tx.To().Hex()
	}
	return &txInfo, nil
}

func (cl *nodeClient) getTxAndStatus(ctx context.Context, txID string) (*types.Transaction, models.TxStatus, error) {
	tx, pending, err := cl.ethClient.TransactionByHash(ctx, common.HexToHash(txID))
	if err != nil {
		if err == ethereum.NotFound {
			return &types.Transaction{}, models.TxStatusUnKnown, nil
		}
		return nil, "", err
	}
	if pending {
		return tx, models.TxStatusPending, nil
	}
	return tx, models.TxStatusSuccess, nil
}

func DeserializeTx(rlpTx []byte) (*types.Transaction, error) {
	reader := bytes.NewReader(rlpTx)
	tx := new(types.Transaction)
	rlpStream := rlp.NewStream(reader, 0)
	if err := tx.DecodeRLP(rlpStream); err != nil {
		return nil, err
	}
	return tx, nil
}

func SerializeTx(tx *types.Transaction) ([]byte, error) {
	var b bytes.Buffer
	if err := tx.EncodeRLP(&b); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
