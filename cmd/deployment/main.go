package main

import (
	"fmt"
	"os"

	"go-sol-example/api"
	"go-sol-example/pkg/auth"

	"github.com/ethereum/go-ethereum/ethclient"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// address of ethereum env
	rpcUrl := os.Getenv("RPC_URL")
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		panic(err)
	}

	// create auth and transaction package for deploying smart contract
	privateKeyAddress := os.Getenv("PRIVATE_KEY")
	auth := auth.GetAccountAuth(client, privateKeyAddress)

	// deploying smart contract
	address, tx, instance, err := api.DeployApi(auth, client)
	if err != nil {
		panic(err)
	}

	fmt.Println("contract address ->", address.Hex())
	fmt.Println("instance ->", instance)
	fmt.Println("tx ->", tx.Hash().Hex())
	fmt.Println()
	fmt.Println("Please modify .env file to renew contract address.")
}
