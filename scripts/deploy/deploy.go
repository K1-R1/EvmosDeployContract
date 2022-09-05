/** deploy.go deploys the Token contract tp
    a local evmos node.
    It utilises the Go-ethereum contract binding script:
	scripts/token/token.go
    When called, it requires the hex private key of the deployer.
*/

package main

import (
	"fmt"
	"log"
	"os"

	token "github.com/K1-R1/EvmosDeployContract/scripts/token"
	util "github.com/K1-R1/EvmosDeployContract/scripts/utils"
)

func main() {
	// Get client for local node
	client, err := util.GetClient()
	if err != nil {
		log.Fatalf("Failed to get client: %v", err)
	}

	// Derive Private key and address from input arguments
	deployerPrivateKey, deployerAddress, err := util.GetPKAndAddress(os.Args[1])
	if err != nil {
		log.Fatalf("Failed to get private key and address: %v", err)
	}

	// Get auth for deployer
	auth, err := util.GetAuth(client, deployerPrivateKey, deployerAddress)
	if err != nil {
		log.Fatalf("Failed to get auth: %v", err)
	}

	// Deploy Token contract as deployer
	address, tx, _, err := token.DeployToken(auth, client)
	if err != nil {
		log.Fatalf("Failed to deploy contract: %v", err)
	}

	// Display info
	fmt.Println("\n\nDeployed Contract\n---------------------------------------------")
	fmt.Printf("Deployed contract address: %v\n", address.Hex())
	fmt.Printf("Deployed by account with address: %v\n", deployerAddress)
	fmt.Printf("Deployed in transaction with hash: %v\n", tx.Hash().Hex())
}
