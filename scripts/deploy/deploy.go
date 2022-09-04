package main

import (
	"fmt"
	"log"
	"os"

	token "github.com/K1-R1/EvmosDeployContract/scripts/token"
	util "github.com/K1-R1/EvmosDeployContract/scripts/utils"
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

	//get auth
	auth, err := util.GetAuth(client, deployerPrivateKey, deployerAddress)
	if err != nil {
		log.Fatal(err)
	}

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
