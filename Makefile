.PHONY: run solc

__dirs := $(shell mkdir -p api)
CONTRACT := BalanceCheck

run:
	go run cmd/server/main.go

deploy:
	go run cmd/deployment/main.go

solc: api/$(CONTRACT).go

api/$(CONTRACT).go: build/$(CONTRACT).abi build/$(CONTRACT).bin
	abigen --abi=./build/$(CONTRACT).abi --bin=./build/$(CONTRACT).bin --pkg=api --out=./api/$(CONTRACT).go

build/$(CONTRACT).abi: contracts/$(CONTRACT).sol
	solc --optimize --abi ./contracts/$(CONTRACT).sol -o build --overwrite

build/$(CONTRACT).bin: contracts/$(CONTRACT).sol
	solc --optimize --bin ./contracts/$(CONTRACT).sol -o build --overwrite
