package main

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	token "github.com/K1-R1/EvmosDeployContract/scripts/token"
	util "github.com/K1-R1/EvmosDeployContract/scripts/utils"
)

func main() {
	//get client
	client, err := util.GetClient()
	if err != nil {
		log.Fatal(err)
	}

	// 1. contractAddress
	contractAddress := common.HexToAddress(os.Args[1])

	//2. Deployer private key and address
	deployerPrivateKey, deployerAddress, err := util.GetPKAndAddress(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	//3. Receiver address
	_, receiverAddress, err := util.GetPKAndAddress(os.Args[3])
	if err != nil {
		log.Fatal(err)
	}

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
	//get auth
	auth, err := util.GetAuth(client, deployerPrivateKey, deployerAddress)
	if err != nil {
		log.Fatal(err)
	}

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
	//Display token balance
	fbal := new(big.Float)
	fbal.SetString(deployerBalance.String())
	value := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))
	fmt.Printf("deployer balance before: %f\n", value)
	//
	fbal = new(big.Float)
	fbal.SetString(receiverBalance.String())
	value = new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))
	fmt.Printf("receiver balance before: %f\n", value)
	//tx
	fmt.Printf("Transferred 10 tokens from contract deployer(%v) to receiver(%v) in transaction: %v\n", deployerAddress, receiverAddress, tx.Hash().Hex())
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
