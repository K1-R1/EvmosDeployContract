# EvmosDeployContract

This project deploys an ERC20 contract to a local Evmos node with Golang, and then querys and transfer tokens on said contract, before running tests on the Golang scripts used.

### Requirements

- [Golang v1.19+](https://go.dev/)
- [Evmos](https://docs.evmos.org/validators/quickstart/installation.html)

### To run

Once evmos is installed, run the local node with:

```shell
cd evmos
./init.sh
```

Once the localnet has started; in a seperate terminal run the deploy, query and transfer, and test scripts with:

```shell
./run_all.sh
```

This will display the relevant infortmation within the terminal.

## How this project was made

#### Evmos

After installing evmos, I edited the `evmos/init.sh` file with an additional key to act as the recipient of the ERC20 transfer by adding:

```
evmosd keys add $KEY2 --keyring-backend $KEYRING --algo $KEYALGO
```

after validation of the genesis transaction. This results in the local node having 1 validator but 2 accounts.

### Develop, compile and deploy an ERC20 token smart contract to local node

First, I created the ERC20 contract as the solidity file `contract/Token.sol`, which utilises the [OpenZeppelin](https://github.com/OpenZeppelin/openzeppelin-contracts/blob/master/contracts/token/ERC20/ERC20.sol) ERC20 implentation, with an alteration to mint 100 tokens to the deployer of the contract.

In order to compile the `Token.sol` contract, I firstly installed the OpenZeppelin contract library using npm:

```shell
npm install @openzeppelin/contracts
```

I then compiled the `Token.sol` contract using [Solc](https://docs.soliditylang.org/en/v0.8.16/installing-solidity.html), outputting the abi and bytecode of the contract and its dependencies to `contract/build` with the commands:

```shell
solc --abi contract/Token.sol -o contract/build --allow-paths
solc --bin contract/Token.sol -o contract/build --allow-paths
```

With the abi and bytecode, I generated a [Go Ethereum](https://geth.ethereum.org/docs/) contract binding file (`scripts/token/token.go`) that would provide functionality to interact with the `Token.sol` contract from within Golang files, via:

```shell
abigen --abi ./contract/build/Token.abi --pkg token --type Token --out ./scripts/token/token.go --bin ./contract/build/Token.bin
```

I then created the `scripts/deploy/deploy.go` file that deploys the `Token.sol` contract, utilising the `token.go` contract bindings and a file with helper functions `scripts/utils/utils.go`.

### Query and transfer token balances on the deployed smart contract

The `scripts/query_and_transfer/query_and_transfer.go` file checks the `TOK` (token of the `Token.sol` contract) balance of two accounts; the first is the deployer of the contract, who received the initial supply when deploying. The second is an account with no `TOK`. The file then transfers 10 `TOK` from the deployer, to the second account. Before finally checking their balances a second time; to see that the second account now owns the 10 transferred `TOK`. This file utilises the `token.go` contract bindings to interact with the deployed contract, and the `utils.go` file's helper functions.

### Testing

This project contains two types of tests:

- [Test Driven Development (TDD)](https://www.codementor.io/@cyantarek15/how-table-driven-tests-makes-writing-unit-tests-exciting-and-fun-in-go-15g1wzdf7g)
- [Behaviour Driven Development (BDD)](https://medium.com/javascript-scene/behavior-driven-development-bdd-and-functional-testing-62084ad7f1f2)

`TDD` style tests are perfomred on the `scripts/utils/utils.go` file via the `tests/utils_test.go` file. This file utilises Golang's built-in testing package, and a set of test-helper functions from `tests/test_utils/test_utils.go`.

`test_utils.go` contains a set of helper functions for testing within a simulated backend.

`BDD` style tests are performed on the deployed `Token` contract, with the use of a simulated backend via `test_utils.go` and the testing packages [ginkgo](https://onsi.github.io/ginkgo/) and [gomega](https://onsi.github.io/gomega/).

### Run all

I utilised the `run_all.sh` file in order to deploy the contract, query and transfer token balances, and run the tests from a single command. The relevant account variables from the local node are loaded and passed into the Golang files when run, with all the relevant information being displayed in the terminal.
