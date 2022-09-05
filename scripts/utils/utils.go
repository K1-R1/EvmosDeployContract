/** utils.go contains helper functions to be used in deploy.go and
  query_and_transfer.go
*/

package utils

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// GetClient retrieves and returns an ethClient for the local evmos node
func GetClient() (*ethclient.Client, error) {
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		return nil, err
	}
	return client, err
}

// GetPKAndAddress derives and returns an ecsda private key and associated address,
// from a hexkey from os.Args
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

// GetAuth derives the transaction options for a given primary key and associated address,
// with a default set of values
func GetAuth(client *ethclient.Client, pk *ecdsa.PrivateKey, address common.Address) (*bind.TransactOpts, error) {

	// Check for invalid address
	if address == common.HexToAddress("0x0") {
		return nil, errors.New("Invalid address")
	}

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

// GetReadableBalance takes the ERC20 contract balance of an account, and reformatts the balance with the correct
// decimal places from the contract's decimals variable
func GetReadableBalance(balance *big.Int, decimals uint8) *big.Float {
	bal := new(big.Float)
	bal.SetString(balance.String())
	value := new(big.Float).Quo(bal, big.NewFloat(math.Pow10(int(decimals))))
	return value
}
