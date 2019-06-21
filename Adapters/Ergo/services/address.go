package services

import (
	"bytes"
	"context"
	"golang.org/x/crypto/blake2b"

	"github.com/btcsuite/btcutil/base58"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/logger"
)

const (
	checksumBytesCount         = 4
	addressTypeBytesCount      = 1
	p2PkAddressTypeByte   byte = 0x01
)

func (cl *nodeClient) ValidateAddress(ctx context.Context, address string) (bool, error) {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'ValidateAddress' for %s", address)
	if len(address) < checksumBytesCount+addressTypeBytesCount {
		log.Debugf("wrong length of address %s", address)
		return false, nil
	}
	addressBytes := base58.Decode(address)
	prefix := byte(cl.conf.ChainID) + p2PkAddressTypeByte
	if addressBytes[0] != prefix {
		log.Debugf("wrong type of address or network type. Network should be %v, address type %v", cl.conf.ChainID, p2PkAddressTypeByte)
		return false, nil
	}
	checksum, err := addressChecksum(addressBytes[:len(addressBytes)-checksumBytesCount], checksumBytesCount)
	if err != nil {
		return false, err
	}
	if !bytes.Equal(checksum, addressBytes[len(addressBytes)-checksumBytesCount:]) {
		log.Debugf("wrong checksum of address %s", address)
		return false, nil
	}
	return true, nil
}

func addressChecksum(b []byte, checksumSize int) ([]byte, error) {
	d := make([]byte, 0)
	fh, err := blake2b.New256(nil)
	if err != nil {
		return d, err
	}
	fh.Write(b)
	h := fh.Sum(d[:0])
	if err != nil {
		return nil, err
	}
	c := make([]byte, checksumSize)
	copy(c, h[:checksumSize])
	return c, nil
}

func (cl *nodeClient) PublicKeyFromAddress(ctx context.Context, address string) []byte {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'PublicKeyFromAddress' for %s", address)
	addr := base58.Decode(address)
	publicKey := addr[addressTypeBytesCount : len(addr)-checksumBytesCount]
	return publicKey
}
