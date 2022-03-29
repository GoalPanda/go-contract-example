.PHONY: run solc

__dirs := $(shell mkdir -p api)

run:
	go run cmd/main.go

solc: api/BalanceCheck.go

api/BalanceCheck.go: build/BalanceCheck.abi build/BalanceCheck.bin
	abigen --abi=./build/BalanceCheck.abi --bin=./build/BalanceCheck.bin --pkg=api --out=./api/BalanceCheck.go

build/BalanceCheck.abi: contracts/BalanceCheck.sol
	solc --optimize --abi ./contracts/BalanceCheck.sol -o build --overwrite

build/BalanceCheck.bin: contracts/BalanceCheck.sol
	solc --optimize --bin ./contracts/BalanceCheck.sol -o build --overwrite
