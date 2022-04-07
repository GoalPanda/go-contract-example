package main

import (
	"math/big"
	"net/http"
	"os"
	"strconv"

	"go-sol-example/api" // this would be your generated smart contract bindings
	"go-sol-example/pkg/auth"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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

	// get contract address
	contractAddress := os.Getenv("CONTRACT_ADDRESS")

	// creating api object to intract with smart contract function
	conn, err := api.NewApi(common.HexToAddress(contractAddress), client)
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
		auth := auth.GetAccountAuth(client, accountPrivateKey)

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

		auth := auth.GetAccountAuth(client, accountPrivateKey)
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
