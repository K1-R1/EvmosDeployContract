package main

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	token "github.com/K1-R1/EvmosDeployContract/scripts/token"
	util "github.com/K1-R1/EvmosDeployContract/scripts/utils"
)

func main() {
	// Get client for local node
	client, err := util.GetClient()
	if err != nil {
		log.Fatalf("Failed to get client: %v", err)
	}

	// Get deployed Token contract address, from Args
	contractAddress := common.HexToAddress(os.Args[1])

	// Get deployer's private key and address, from Args
	deployerPrivateKey, deployerAddress, err := util.GetPKAndAddress(os.Args[2])
	if err != nil {
		log.Fatalf("Failed to get private key and address: %v", err)
	}

	// Get deployer's auth
	auth, err := util.GetAuth(client, deployerPrivateKey, deployerAddress)
	if err != nil {
		log.Fatalf("Failed to get auth: %v", err)
	}

	// Get receiver's address, from Args
	_, receiverAddress, err := util.GetPKAndAddress(os.Args[3])
	if err != nil {
		log.Fatalf("Failed to get private key and address: %v", err)
	}

	// Get instance of Token contract
	instance, err := token.NewToken(contractAddress, client)
	if err != nil {
		log.Fatalf("Failed to get contract insstance: %v", err)
	}

	// Check Starting balances
	// Of deployer
	deployerBalance, err := instance.BalanceOf(&bind.CallOpts{}, deployerAddress)
	if err != nil {
		log.Fatalf("Failed to get token balance: %v", err)
	}
	// Of receiver
	receiverBalance, err := instance.BalanceOf(&bind.CallOpts{}, receiverAddress)
	if err != nil {
		log.Fatalf("Failed to get token balance: %v", err)
	}

	// Transfer 10 tokens from deployer to reciever
	// Set amount of tokens to be transferred;
	// as 10 with 18 decimals, as per Token contract
	transferAmount, ok := new(big.Int).SetString("10000000000000000000", 10)
	if !ok {
		log.Fatalf("Failed to set transferAmount: %v", err)
	}

	// Transfer tokens from deployer to receiver address
	tx, err := instance.Transfer(auth, receiverAddress, transferAmount)
	_ = tx
	if err != nil {
		log.Fatalf("Failed to transfer tokens: %v", err)
	}

	// Wait for tx to be executed
	time.Sleep(5 * time.Second)

	// Check end balances
	// Of deployer
	deployerBalanceAfter, err := instance.BalanceOf(&bind.CallOpts{}, deployerAddress)
	if err != nil {
		log.Fatalf("Failed to get token balance: %v", err)
	}
	// Of receiver
	receiverBalanceAfter, err := instance.BalanceOf(&bind.CallOpts{}, receiverAddress)
	if err != nil {
		log.Fatalf("Failed to get token balance: %v", err)
	}

	// Display
	// Get contract name
	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		log.Fatalf("Failed to get contract name: %v", err)
	}
	// Get contract symbol
	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Fatalf("Failed to get contract symbol: %v", err)
	}
	// Get contract decimals
	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		log.Fatalf("Failed to get contract decimals: %v", err)
	}

	fmt.Println("\n\nQuery and transfer tokens\n---------------------------------------------")
	fmt.Printf("Token name: %v    Token symbol: %v\n", name, symbol)
	fmt.Printf("Starting balances:\n")
	fmt.Printf("                  Address                    |              %v balance           \n", symbol)
	fmt.Printf("---------------------------------------------|---------------------------------------------\n")
	fmt.Printf("%v   | %v\n", deployerAddress, util.GetReadableBalance(deployerBalance, decimals))
	fmt.Printf("%v   | %v\n", receiverAddress, util.GetReadableBalance(receiverBalance, decimals))

	fmt.Printf("\n\nTransfer:\n10 %v transferred from contract deployer(%v),\nto receiver(%v),\nin transaction: %v\n", symbol, deployerAddress, receiverAddress, tx.Hash().Hex())

	fmt.Printf("\n\nEnding balances:\n")
	fmt.Printf("                  Address                    |              %v balance           \n", symbol)
	fmt.Printf("---------------------------------------------|---------------------------------------------\n")
	fmt.Printf("%v   | %v\n", deployerAddress, util.GetReadableBalance(deployerBalanceAfter, decimals))
	fmt.Printf("%v   | %v\n\n", receiverAddress, util.GetReadableBalance(receiverBalanceAfter, decimals))
}
