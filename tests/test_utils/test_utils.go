/** test_utils.go contains helper functions for testing.
  Primarily for BDD testing of the Token contract in a simulated backend
*/

package testUtil

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	token "github.com/K1-R1/EvmosDeployContract/scripts/token"
	util "github.com/K1-R1/EvmosDeployContract/scripts/utils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	// Defines the amount of tokens initially deployed to a contract
	// on the simulated backend
	initialBalance = Ten18

	// Defines the max gas per block for the simulated backend
	MaxGasPerBlock = uint64(5000000)

	// Defines 10^18 as a big integer
	Ten18 = big.NewInt(1e18)

	// According to the go-ethereum docs, the chain ID for simulated
	// backends must be 1337
	// (https://pkg.go.dev/github.com/ethereum/go-ethereum@v1.10.19/accounts/abi/bind/backends#NewSimulatedBackend).
	TestChainID = big.NewInt(1337)
)

// ReceiptBackend defines an interface, which can be used to get the transaction receipt
// for a given ethclient or simulated backend.
type ReceiptBackend interface {
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
}

// DeployContractAndCommit deploys an instance of the ERC20 Token contract
// and commits the transaction to the simulated backend.
// The function returns the contract address, the transaction, and an
// instance of the contract binding.
func DeployContractAndCommit(auth *bind.TransactOpts, client *backends.SimulatedBackend) (common.Address, *types.Transaction, *token.Token, error) {
	// Deploy contract
	contractAddress, tx, contract, err := token.DeployToken(auth, client)
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	// Commit transaction on simulated backend
	client.Commit()

	return contractAddress, tx, contract, nil
}

// FillTransactionSignerFields takes the transaction signer, the client
// and a byte array of the data to be called in a transaction.
// It gathers necessary gas price, nonce and estimated gas and assigns
// these to the fields of the transaction signer, which the function then
// returns.
func FillTransactionSignerFields(auth *bind.TransactOpts, client *ethclient.Client, callMsg ethereum.CallMsg) (*bind.TransactOpts, error) {
	// Get gas price suggestion from client
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	callMsg.GasPrice = gasPrice

	// Get current nonce for deployer address
	nonce, err := client.PendingNonceAt(context.Background(), auth.From)
	if err != nil {
		return nil, err
	}

	// Estimate gas usage
	gasLimit, err := client.EstimateGas(context.Background(), callMsg)
	if err != nil {
		return nil, err
	}

	// Fill transaction signer fields
	auth.GasLimit = gasLimit
	auth.GasPrice = gasPrice
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)

	return auth, nil
}

// GetClientAndTransactionSigner connects to a local Evmos node on port 8545,
// queries the chain id and uses this together with the private key to create
// a transaction signer.
// The function returns the client and the transaction signer.
func GetClientAndTransactionSigner(privKey *ecdsa.PrivateKey) (*ethclient.Client, *bind.TransactOpts, error) {
	// Connect to blockchain node given a valid URL
	client, err := util.GetClient()
	if err != nil {
		return nil, nil, err
	}

	// Get chain id from client in order to generate the transaction signer
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, nil, err
	}

	// Create transaction signer
	auth, err := bind.NewKeyedTransactorWithChainID(privKey, chainID)
	if err != nil {
		return nil, nil, err
	}

	return client, auth, nil
}

// GeneratePrivKeysAndAddresses returns a slice of private keys and addresses.
func GeneratePrivKeysAndAddresses(n uint64) ([]*ecdsa.PrivateKey, []common.Address, error) {
	// Create a slice of private keys
	privKeys := make([]*ecdsa.PrivateKey, n)

	// Create a slice of addresses
	addresses := make([]common.Address, n)

	// Create a slice of random private keys
	for i := uint64(0); i < n; i++ {
		// Create a new private key
		privKey, err := crypto.GenerateKey()
		if err != nil {
			return nil, nil, err
		}

		// Add private key to slice
		privKeys[i] = privKey

		// Add address to slice
		addresses[i] = crypto.PubkeyToAddress(privKey.PublicKey)
	}

	return privKeys, addresses, nil
}

// GetCallData returns the necessary byte array to fill an ethereum.CallMsg
// struct's data field. This data is a byte representation of the method
// call paired with its corresponding arguments.
func GetCallData(name string, args ...interface{}) ([]byte, error) {
	// Return the contract ABI in order to define the necessary transaction
	// data.
	tokenABI, err := token.TokenMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	// Generate the call data for the transfer method using the ABI
	callData, err := tokenABI.Pack("transfer", args...)
	if err != nil {
		return nil, err
	}

	return callData, nil
}

// GetReceipt converts a given transaction hash in hex string format and
// returns the transaction receipt, if the hash is valid.
func GetReceipt(backend ReceiptBackend, txHashHex string) (*types.Receipt, error) {
	// Convert transaction hash, for which the receipt should be returned
	txHash := common.HexToHash(txHashHex)

	// Get transaction receipt
	receipt, err := backend.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		return nil, err
	}

	return receipt, nil
}

// GetSimulatedClientAndTransactionSigner establishes a new simulated backend
// for testing purposes. An initial token balance is assigned to the address
// of the given private key. The maximum gas a block can consume is defined
// with the blockGasLimit input.
// The private key and chain id are used to create a transaction signer for
// any transactions on the blockchain.
// The function returns the client and the transaction signer.
func GetSimulatedClientAndTransactionSigner(privKey *ecdsa.PrivateKey, blockGasLimit uint64, chainID *big.Int) (*backends.SimulatedBackend, *bind.TransactOpts, error) {
	// Define genesis state for simulated backend
	address := crypto.PubkeyToAddress(privKey.PublicKey)
	genesisAlloc := map[common.Address]core.GenesisAccount{
		address: {
			Balance: initialBalance,
		},
	}

	// Get simulated backend as client
	client := backends.NewSimulatedBackend(genesisAlloc, blockGasLimit)

	// Define transaction signer
	auth, err := bind.NewKeyedTransactorWithChainID(privKey, chainID)
	if err != nil {
		return nil, nil, err
	}

	return client, auth, nil
}
