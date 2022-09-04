package utils

import (
	"crypto/ecdsa"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetClient() (*ethclient.Client, error) {
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatal(err)
	}
	return client, err
}

// Derive PK and address from Args
func GetPKAndAddress(hexkey string) (*ecdsa.PrivateKey, common.Address, error) {
	privateKey, err := crypto.HexToECDSA(hexkey)
	if err != nil {
		return nil, common.HexToAddress("0x0"), err
	}

	publicKey, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return nil, common.HexToAddress("0x0"), err
	}

	address := crypto.PubkeyToAddress(*publicKey)

	return privateKey, address, nil
}

//get auth

//Display token balance
