package models

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/wavesplatform/gowaves/pkg/proto"
)

type Transaction struct {
	Tx proto.Transaction
}

type TransactionType struct {
	Type    uint `json:"type"`
	Version uint `json:"version"`
}

type TransactionsField []Transaction

func (b *Transaction) UnmarshalJSON(data []byte) error {
	var txType = TransactionType{}
	err := json.Unmarshal(data, &txType)
	if err != nil {
		return fmt.Errorf("can't get tx type from json: %s", err)
	}

	var tx proto.Transaction
	switch proto.TransactionType(txType.Type) {
	case proto.TransferTransaction:
		if txType.Version == 1 {
			tx = &proto.TransferV1{}
		} else {
			tx = &proto.TransferV2{}
		}
	case proto.PaymentTransaction:
		tx = &proto.Payment{}
	case proto.MassTransferTransaction:
		tx = &proto.MassTransferV1{}
	default:
		tx = &UnknownTransaction{}
	}

	err = json.Unmarshal(data, tx)
	if err != nil {
		return fmt.Errorf("can't unmarshal tx for type %s: %s", txType.Type, err)
	}
	b.Tx = tx

	return nil
}

type UnknownTransaction struct {
}

func (UnknownTransaction) GetID() []byte {
	return nil
}

func (UnknownTransaction) Valid() (bool, error) {
	return false, errors.New("not implemented. Should not be used")
}

func (UnknownTransaction) MarshalBinary() ([]byte, error) {
	return nil, errors.New("not implemented. Should not be used")
}

func (UnknownTransaction) UnmarshalBinary([]byte) error {
	return errors.New("not implemented. Should not be used")
}
