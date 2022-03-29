package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strconv"

	"go-sol-example/api" // this would be your generated smart contract bindings

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"

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
	auth := getAccountAuth(client, privateKeyAddress)

	// deploying smart contract
	address, tx, instance, err := api.DeployApi(auth, client)
	if err != nil {
		panic(err)
	}

	fmt.Println(address.Hex())

	fmt.Println("instance->", instance)
	fmt.Println("tx->", tx.Hash().Hex())

	// creating api object to intract with smart contract function
	conn, err := api.NewApi(common.HexToAddress(address.Hex()), client)
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	r.GET("/balance", func(c *gin.Context) {
		reply, err := conn.Balance(&bind.CallOpts{})
		if err != nil {
			c.JSON(http.StatusBadRequest, reply)
		} else {
			c.JSON(http.StatusOK, reply)
		}
	})

	r.GET("/admin", func(c *gin.Context) {
		reply, err := conn.Admin(&bind.CallOpts{})
		if err != nil {
			c.JSON(http.StatusBadRequest, reply)
		} else {
			c.JSON(http.StatusOK, reply)
		}
	})

	r.POST("/deposit/:amount", func(c *gin.Context) {
		amount := c.Param("amount")
		amt, _ := strconv.Atoi(amount)

		// gets address of account by which amount to be deposit
		accountPrivateKey := c.PostForm("accountPrivateKey")
		if accountPrivateKey == "" {
			c.JSON(http.StatusBadRequest, "missing parameter accountPrivateKey")
			panic("missing parameter accountPrivateKey")
		}

		// creating auth object for above account
		auth := getAccountAuth(client, accountPrivateKey)

		reply, err := conn.Deposit(auth, big.NewInt(int64(amt)))
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		} else {
			c.JSON(http.StatusOK, reply)
		}
	})

	r.POST("/withdraw/:amount", func(c *gin.Context) {
		amount := c.Param("amount")
		amt, _ := strconv.Atoi(amount)

		accountPrivateKey := c.PostForm("accountPrivateKey")
		if accountPrivateKey == "" {
			c.JSON(http.StatusBadRequest, "missing parameter accountPrivateKey")
			panic("missing parameter accountPrivateKey")
		}

		auth := getAccountAuth(client, accountPrivateKey)
		// auth.Nonce.Add(auth.Nonce, big.NewInt(int64(1))) //it is use to create next nonce of account if it has to make another transaction

		reply, err := conn.Withdraw(auth, big.NewInt(int64(amt)))
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		} else {
			c.JSON(http.StatusOK, reply)
		}
	})

	// Start server
	r.Run()
}

//function to create auth for any account from its private key
func getAccountAuth(client *ethclient.Client, privateKeyAddress string) *bind.TransactOpts {

	privateKey, err := crypto.HexToECDSA(privateKeyAddress)
	if err != nil {
		panic(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("invalid key")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		panic(err)
	}
	fmt.Println("nonce=", nonce)
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		panic(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		panic(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = big.NewInt(1000000000)

	return auth
}
