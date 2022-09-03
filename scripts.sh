# File paths
DEPLOY=scripts/deploy/deploy.go
QUERY=scripts/query_and_transfer/query_and_transfer.go

# Account info 
DEPLOYER_KEY="mykey"
RECEIVER_KEY="mykey2"
DEPLOYER_PK=$(evmosd keys unsafe-export-eth-key $DEPLOYER_KEY --keyring-backend=test)
RECEIVER_PK=$(evmosd keys unsafe-export-eth-key $RECEIVER_KEY --keyring-backend=test)

# Deploy contract
go run $DEPLOY $DEPLOYER_PK > tmp.txt
cat tmp.txt
CONTRACT_ADDRESS=$(cat tmp.txt | grep 'Contract address' | grep -o '0x[0-9a-zA-Z]*')
rm -f tmp.txt

# Wait for transaction to be included in a block
echo "Waiting ..."
sleep 5

# Query and transfer tokens
go run $QUERY $CONTRACT_ADDRESS $DEPLOYER_PK $RECEIVER_PK