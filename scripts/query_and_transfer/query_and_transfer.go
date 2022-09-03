package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	token "github.com/K1-R1/EvmosDeployContract/scripts/token"
	util "github.com/K1-R1/EvmosDeployContract/scripts/utils"
)

func main() {
	//get client
	client, err := util.GetClient()
	if err != nil {
		log.Fatal(err)
	}

	// Set vars from os.Args
	// 1. contractAddress
	contractAddress := common.HexToAddress(os.Args[1])
	//2. Deployer private key and address
	deployerPrivateKey, err := crypto.HexToECDSA(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	deployerPublicKey, ok := deployerPrivateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	deployerAddress := crypto.PubkeyToAddress(*deployerPublicKey)
	//3. Receiver address
	receiverPrivateKey, err := crypto.HexToECDSA(os.Args[3])
	if err != nil {
		log.Fatal(err)
	}

	receiverPublicKey, ok := receiverPrivateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	receiverAddress := crypto.PubkeyToAddress(*receiverPublicKey)

	// Get instance of Token contract
	instance, err := token.NewToken(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	// Check starting balance of deployer
	deployerBalance, err := instance.BalanceOf(&bind.CallOpts{}, deployerAddress)
	if err != nil {
		log.Fatal(err)
	}
	// Check starting balance of receiver
	receiverBalance, err := instance.BalanceOf(&bind.CallOpts{}, receiverAddress)
	if err != nil {
		log.Fatal(err)
	}

	//Transfer 10 tokens from deployer to reciever
	//setup auth for deployer
	nonce, err := client.PendingNonceAt(context.Background(), deployerAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(deployerPrivateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = gasPrice

	//set amount
	transferAmount, ok := new(big.Int).SetString("10000000000000000000", 10)
	if !ok {
		log.Fatal(err)
	}

	// Transfer tokens from deployer address to receiver address
	tx, err := instance.Transfer(auth, receiverAddress, transferAmount)
	if err != nil {
		log.Fatal(err)
	}
	_ = tx
	time.Sleep(5 * time.Second)

	//check end balances
	// Check end balance of deployer
	deployerBalanceAfter, err := instance.BalanceOf(&bind.CallOpts{}, deployerAddress)
	if err != nil {
		log.Fatal(err)
	}
	// Check end balance of receiver
	receiverBalanceAfter, err := instance.BalanceOf(&bind.CallOpts{}, receiverAddress)
	if err != nil {
		log.Fatal(err)
	}

	//display values
	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	//before
	fbal := new(big.Float)
	fbal.SetString(deployerBalance.String())
	value := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))
	fmt.Printf("deployer balance before: %f\n", value)

	fbal = new(big.Float)
	fbal.SetString(receiverBalance.String())
	value = new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))
	fmt.Printf("receiver balance before: %f\n", value)
	//after
	fbal = new(big.Float)
	fbal.SetString(deployerBalanceAfter.String())
	value = new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))
	fmt.Printf("deployer balance after: %f\n", value)

	fbal = new(big.Float)
	fbal.SetString(receiverBalanceAfter.String())
	value = new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))
	fmt.Printf("receiver balance after: %f\n", value)
}
