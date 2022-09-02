package utils

import (
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

func GetClient() (*ethclient.Client, error) {
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatal(err)
	}
	return client, err
}
