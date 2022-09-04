package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	token "github.com/K1-R1/EvmosDeployContract/scripts/token"
	util "github.com/K1-R1/EvmosDeployContract/scripts/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func main() {
	client, err := util.GetClient()
	if err != nil {
		log.Fatal(err)
	}

	//Derive PK and address from Args
	deployerPrivateKey, deployerAddress, err := util.GetPKAndAddress(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	//
	//get auth
	nonce, err := client.PendingNonceAt(context.Background(), deployerAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Get chain id from client in order to generate the transaction signer
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// Create transaction signer
	auth, err := bind.NewKeyedTransactorWithChainID(deployerPrivateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = gasPrice
	//

	address, tx, _, err := token.DeployToken(auth, client)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("Contract address:", address.Hex())
	//deployer address
	fmt.Printf("Contract address: %v\n", address.Hex())
	// fmt.Println(tx.Hash().Hex())
	fmt.Printf("tx hash: %v\n", tx.Hash().Hex())
}
