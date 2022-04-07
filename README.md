# Go-Solidity-Example

## Pre-requirements

- [go](https://github.com/golang/go) >= 1.17
- [ganache](https://github.com/trufflesuite/ganache)
- [solc](https://github.com/ethereum/solidity) >= 0.8.0
- [geth](https://github.com/ethereum/go-ethereum)

## Compiling Smart Contract

```sh
make solc
```

## Execution

```sh
# run local chain
ganache
# modify private key and rpc url if needed
vim .env
# deploy smart contract
make deploy
# modify contract address
vim .env
# run
make run
```

You can use either Postman or curl to send request to the server.

## Others

HTTP Content-Type: application/x-www-form-urlencoded.

## Reference

- [Smart Contract with Golang](https://medium.com/nerd-for-tech/smart-contract-with-golang-d208c92848a9)
  - [GitHub](https://github.com/02amanag/SmartContractWithGolang)
