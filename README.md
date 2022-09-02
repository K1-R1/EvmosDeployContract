# EvmosDeployContract

### Dependencies

- Golang
- jq
- gcc (from build-essential)
- node/npm
- solc
- geth (with abigen)

1. Set up and run an Evmos node locally✅

- via wsl
- Install go
- Install jq
- Install gcc (with sudo apt-get install build-essential)
- Build Evmos binaries
- Run local node (with ./init.sh, in evmos dir)

2.  Develop, compile and deploy an ERC20 token smart contract to local node✅

- Create ERC20 contract
- Install contract dependencies (with npm install @openzeppelin/contracts)
- Compile contract (with solc. ABI: solc --abi contract/Token.sol -o contract/build --allow-paths /, BIN: solc --bin contract/Token.sol -o contract/build --allow-paths /)
- Generate Geth bindings (with abigen --abi ./contract/build/Token.abi --pkg token --type Token --out ./scripts/token.go --bin ./contract/build/Token.bin)
- Create deploy file, deploy (with go run scripts/deploy/deploy.go $(evmosd keys unsafe-export-eth-key mykey --keyring-backend=test))

!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
go: to add module requirements and sums:
go mod tidy
