package main

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"

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
	contractAddress := common.HexToAddress(os.Args[1])
	//Deployer private key and address
	deployerPrivateKey, err := crypto.HexToECDSA(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	deployerPublicKey, ok := deployerPrivateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	deployerAddress := crypto.PubkeyToAddress(*deployerPublicKey)
	//Receiver address
	receiverAddress := common.HexToAddress(os.Args[3])

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
	fmt.Printf("balance: %f", deployerBalance)
	// Check starting balance of receiver
	receiverBalance, err := instance.BalanceOf(&bind.CallOpts{}, receiverAddress)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("balance: %f", receiverBalance)

	//Transfer 10 tokens from deployer to reciever

	//check end balances

	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("decimals: %v\n", decimals) // "decimals: 18"

	fmt.Printf("wei: %s\n", bal) // "wei: 74605500647408739782407023"
	fbal := new(big.Float)
	fbal.SetString(bal.String())
	value := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))
	fmt.Printf("balance: %f", value) // "balance: 74605500.647409"
}
