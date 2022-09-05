# run_all.sh sets the necessary variables, runs a script to deploy an ERC20 contract.
# Runs a script to query, transfer, and then query again, the token balances of two accounts.
# Runs tests on the ERC20 contract, and a script containing helper functions

# File paths to go scripts
DEPLOY=scripts/deploy/deploy.go
QUERY=scripts/query_and_transfer/query_and_transfer.go

# Account variables
DEPLOYER_KEY="mykey"
RECEIVER_KEY="mykey2"
DEPLOYER_PK=$(evmosd keys unsafe-export-eth-key $DEPLOYER_KEY --keyring-backend=test)
RECEIVER_PK=$(evmosd keys unsafe-export-eth-key $RECEIVER_KEY --keyring-backend=test)

# Deploy contract, and output to tmp.txt
go run $DEPLOY $DEPLOYER_PK > tmp.txt
# Extract contract addresss of deployed contract, and delete tmp.txt
cat tmp.txt
CONTRACT_ADDRESS=$(cat tmp.txt | grep 'contract address' | grep -o '0x[0-9a-zA-Z]*')
rm -f tmp.txt

# Wait for the contract deployment transaction to be executed
sleep 5

# Query and transfer tokens
go run $QUERY $CONTRACT_ADDRESS $DEPLOYER_PK $RECEIVER_PK

# Run tests
echo "Beginning tests"
echo "---------------------------------------------"
go test ./tests -v