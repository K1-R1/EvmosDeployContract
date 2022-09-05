/** token_contract_test.go contains the testing suite of BDD( Behaviour Driven Development ) style tests for the ERC20 Token
  contract. The token properties and functions are tested using a simulated backend.
  It utilises the "ginkgo" and "gomega" testing packages.
*/

package tests

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	token "github.com/K1-R1/EvmosDeployContract/scripts/token"
	testUtil "github.com/K1-R1/EvmosDeployContract/tests/test_utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/suite"
)

// ContractTestSuite defines the test suite for the Token smart contract.
// It contains slices of private keys and addresses, which are used to test the
// contract methods.
type ContractTestSuite struct {
	suite.Suite

	addresses       []common.Address
	auth            *bind.TransactOpts
	client          *backends.SimulatedBackend
	contract        *token.Token
	contractAddress common.Address
	deployerBalance *big.Int
	privKeys        []*ecdsa.PrivateKey
}

// Initialize the test suite
var s *ContractTestSuite

// SetupTest defines the procedure for setting up BDD tests of the contract
// methods. For this, an instance of the Token smart contract is deployed
// to a simulated backend.
func (suite *ContractTestSuite) SetupTest() {
	// Generate testing accounts
	privKeys, addresses, err := testUtil.GeneratePrivKeysAndAddresses(3)
	if err != nil {
		log.Fatalf("Error generating private key: %v\n", err)
	}

	// Get simulated backend and transaction signer for testing
	client, auth, _ := testUtil.GetSimulatedClientAndTransactionSigner(privKeys[0], testUtil.MaxGasPerBlock, testUtil.TestChainID)

	// Deploy contract
	contractAddress, _, contract, _ := testUtil.DeployContractAndCommit(auth, client)

	// Assign to testing suite
	suite.addresses = addresses
	suite.auth = auth
	suite.client = client
	suite.contract = contract
	suite.contractAddress = contractAddress
	suite.deployerBalance = new(big.Int).Mul(big.NewInt(100), testUtil.Ten18) // 100 TOK
	suite.privKeys = privKeys
}

// TestToken initializes the test suite and runs the tests.
func TestToken(t *testing.T) {
	s = new(ContractTestSuite)
	suite.Run(t, s)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Token Suite")
}

var _ = Describe("approve:", func() {
	BeforeEach(func() {
		s.SetupTest()
	})

	Context("When approving a transaction amount", Ordered, func() {
		It("should have added the amount to the allowance of the recipient address", func() {
			// Define approved amount
			amount := testUtil.Ten18

			// Approve tokens from account1 to account2
			_, err := s.contract.Approve(s.auth, s.addresses[1], amount)
			Expect(err).To(BeNil())

			// Commit transaction
			s.client.Commit()

			allowance, err := s.contract.Allowance(nil, s.addresses[0], s.addresses[1])
			Expect(allowance.Cmp(amount), err).To(Equal(0))
		})
	})
})

var _ = Describe("balance:", func() {
	BeforeEach(func() {
		s.SetupTest()
	})

	Context("Deployer balance after deployment", func() {
		It("should be 100", func() {
			balance, err := s.contract.BalanceOf(nil, s.addresses[0])
			Expect(err).To(BeNil())
			Expect(balance.Cmp(s.deployerBalance)).To(Equal(0))
		})
	})

	Context("Other accounts' balances after deployment", func() {
		It("should be 0", func() {
			balance, err := s.contract.BalanceOf(nil, s.addresses[1])
			Expect(err).To(BeNil())
			Expect(balance.Cmp(new(big.Int))).To(Equal(0))
		})
	})
})

var _ = Describe("transfer:", func() {
	BeforeEach(func() {
		s.SetupTest()
	})

	Context("When sender has sufficient tokens", Ordered, func() {
		// Define transferred amount
		amount := testUtil.Ten18

		BeforeEach(func() {
			// Transfer tokens from account1 to account2
			_, err := s.contract.Transfer(s.auth, s.addresses[1], amount)
			Expect(err).To(BeNil())

			// Commit transaction
			s.client.Commit()
		})

		It("should have deducted the transferred amount from the sender balance", func() {
			senderBalance, err := s.contract.BalanceOf(nil, s.addresses[0])
			Expect(err).To(BeNil())
			Expect(senderBalance.Cmp(new(big.Int).Sub(s.deployerBalance, amount))).To(Equal(0))
		})

		It("should have increased the recipient balance by the transferred amount", func() {
			recipientBalance, err := s.contract.BalanceOf(nil, s.addresses[1])
			Expect(err).To(BeNil())
			Expect(recipientBalance.Cmp(amount)).To(Equal(0))
		})
	})

	Context("When sender does not have sufficient tokens", Ordered, func() {
		// Define transferred amount
		amount := new(big.Int).Mul(big.NewInt(100000), testUtil.Ten18)

		BeforeEach(func() {
			// Transfer tokens from account1 to account2
			_, err := s.contract.Transfer(s.auth, s.addresses[1], amount)
			Expect(err).Error()

			// Commit transaction
			s.client.Commit()
		})

		It("should not have deducted the transferred amount from the sender balance", func() {
			senderBalance, err := s.contract.BalanceOf(nil, s.addresses[0])
			Expect(senderBalance.Cmp(s.deployerBalance), err).To(Equal(0))
		})

		It("should not have increased the recipient balance", func() {
			recipientBalance, err := s.contract.BalanceOf(nil, s.addresses[1])
			Expect(recipientBalance.Cmp(big.NewInt(0)), err).To(Equal(0))
		})
	})
})

var _ = Describe("transferFrom:", func() {
	BeforeEach(func() {
		s.SetupTest()

		/** In order to transfer from, the sender must have sufficient Evmos
		  tokens to pay the gas cost of the transfer. Hence, before the
		  tests, an appropriate amount of Evmos is sent to the sender address,
		  which upon genesis of the simulated backend, does not hold any.
		*/

		// Define transaction
		tx := types.NewTransaction(1, s.addresses[1], big.NewInt(1e15), testUtil.MaxGasPerBlock, big.NewInt(817704169), nil)

		// Sign the transaction with the private key of the deployer account
		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(testUtil.TestChainID), s.privKeys[0])
		Expect(err).To(BeNil())

		// Send Evmos to account1 for gas usage.
		err = s.client.SendTransaction(context.Background(), signedTx)
		Expect(err).To(BeNil())

		// Commit transaction
		s.client.Commit()
	})

	Context("When no approval was given", func() {
		It("should not be able to transfer tokens from the sender address", func() {
			// Define transferred amount
			amount := testUtil.Ten18

			// Transfer tokens from account 1 to account 2
			_, err := s.contract.TransferFrom(s.auth, s.addresses[1], s.addresses[2], amount)
			Expect(err).Error()

			// Commit transaction
			s.client.Commit()
		})
	})

	Describe("When an amount of tokens is approved to be sent", func() {
		//Define transferred amount
		amount := testUtil.Ten18

		BeforeEach(func() {
			// Approve tokens from deployer account to sender account
			_, err := s.contract.Approve(s.auth, s.addresses[1], amount)
			Expect(err).To(BeNil())

			// Commit transaction
			s.client.Commit()
		})

		Context("and the approver has sufficient tokens in his balance", func() {
			BeforeEach(func() {
				// Define transactor for sender account
				auth, err := bind.NewKeyedTransactorWithChainID(s.privKeys[1], testUtil.TestChainID)
				Expect(err).To(BeNil())

				// Transfer tokens from sender account to recipient
				_, err = s.contract.TransferFrom(auth, s.addresses[0], s.addresses[2], amount)
				Expect(err).To(BeNil())

				// Commit transaction
				s.client.Commit()
			})

			It("should have deducted the transferred amount from the approver's balance", func() {
				approverBalance, err := s.contract.BalanceOf(nil, s.addresses[0])
				Expect(approverBalance.Cmp(new(big.Int).Sub(s.deployerBalance, amount)), err).To(Equal(0))
			})

			It("should have increased recipient's balance by the transferred amount", func() {
				recipientBalance, err := s.contract.BalanceOf(nil, s.addresses[2])
				Expect(recipientBalance.Cmp(amount), err).To(Equal(0))
			})
		})

		Context("and the approver does not have sufficient tokens in his balance", func() {
			BeforeEach(func() {
				// Define transactor for sender account
				auth, err := bind.NewKeyedTransactorWithChainID(s.privKeys[1], testUtil.TestChainID)
				Expect(err).To(BeNil())

				// Transfer twice the approved tokens from sender account to recipient
				_, err = s.contract.TransferFrom(auth, s.addresses[0], s.addresses[2], new(big.Int).Mul(big.NewInt(2), amount))
				Expect(err).Error()

				// Commit transaction
				s.client.Commit()
			})

			It("should not have deducted the transferred amount from the approver's balance", func() {
				approverBalance, err := s.contract.BalanceOf(nil, s.addresses[0])
				Expect(approverBalance.Cmp(s.deployerBalance), err).To(Equal(0))
			})

			It("should not have increased recipient's balance by the transferred amount", func() {
				recipientBalance, err := s.contract.BalanceOf(nil, s.addresses[2])
				Expect(recipientBalance.Cmp(big.NewInt(0)), err).To(Equal(0))
			})
		})
	})
})
