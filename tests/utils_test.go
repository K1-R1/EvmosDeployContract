/** utils_test.go contains TDD ( Test Driven Development ) style tests for the scripts/utils/utils.go
  script. It utilise Golang's built-in "testing" package.
*/

package tests

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"

	util "github.com/K1-R1/EvmosDeployContract/scripts/utils"
	testUtil "github.com/K1-R1/EvmosDeployContract/tests/test_utils"
)

// Test GetClient
// checks if connection to node is successful,
// and that the chain ID is correct for the local node
func TestGetClient(t *testing.T) {
	// Check that connection to node is a success
	client, err := util.GetClient()
	require.NoError(t, err, "Error getting client")

	// Check if chain ID is correct
	chainID, err := client.ChainID(context.Background())
	require.NoError(t, err, "Error getting chain ID")
	require.Equal(t, big.NewInt(9000), chainID, "Incorrect chain ID")
}

// Test GetPKAndAddress
// Checks that private keys and addresses are only derived with valid inputs
func TestGetPKAndAddress(t *testing.T) {
	testcases := []struct {
		name   string
		expErr bool
		hexPK  string
	}{
		{
			"Invalid hexPK",
			true,
			"0",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := util.GetPKAndAddress(tc.hexPK)
			if tc.expErr {
				require.Error(t, err, "GetPKAndAddress should raise an error with invalid hexPK")
			} else {
				require.NoError(t, err, "Error during GetPKAndAddress")
			}
		})
	}
}

// Test GetAuth
// Checks that valid transaction options are only generated with valid inputs
func TestGetAuth(t *testing.T) {
	privKeys, addresses, err := testUtil.GeneratePrivKeysAndAddresses(1)
	require.NoError(t, err, "Error generating private key and address")

	client, err := util.GetClient()
	require.NoError(t, err, "Error getting client")

	testcases := []struct {
		name    string
		expErr  bool
		client  *ethclient.Client
		pk      *ecdsa.PrivateKey
		address common.Address
	}{
		{
			"Valid inputs",
			false,
			client,
			privKeys[0],
			addresses[0],
		},
		{
			"Invalid address",
			true,
			client,
			privKeys[0],
			common.HexToAddress("0x0"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			auth, err := util.GetAuth(tc.client, tc.pk, tc.address)
			if tc.expErr {
				require.Error(t, err, "GetAuth should raise an error")
			} else {
				require.NoError(t, err, "Error during GetAuth")
			}
			_ = auth

		})
	}
}
