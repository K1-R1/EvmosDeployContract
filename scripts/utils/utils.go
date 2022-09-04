package utils

import (
	"context"
	"crypto/ecdsa"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Get client for Local node
func GetClient() (*ethclient.Client, error) {
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		return nil, err
	}
	return client, err
}

// Derive Private key and address from Args
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

// Get auth for a specific private key
func GetAuth(client *ethclient.Client, pk *ecdsa.PrivateKey, address common.Address) (*bind.TransactOpts, error) {

	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		return nil, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(pk, chainID)
	if err != nil {
		return nil, err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = gasPrice

	return auth, nil
}

// Get human-readable value of Token balance, with correct decimals
func GetReadableBalance(balance *big.Int, decimals uint8) *big.Float {
	bal := new(big.Float)
	bal.SetString(balance.String())
	value := new(big.Float).Quo(bal, big.NewFloat(math.Pow10(int(decimals))))
	return value
}
