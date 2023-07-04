package main

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"os"

	"github.com/camilovietnam/solidity/01/api"
	"github.com/ethereum/go-ethereum/ethclient"
)

func loadEnv() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading config.env file")
	}
}

func loadClient() (*ethclient.Client, error) {
	// Sepolia network URL
	sepoliaURL := os.Getenv("SEPOLIA_URL")

	// Create an Ethereum client
	client, err := ethclient.Dial(sepoliaURL)
	if err != nil {
		log.Fatalf("dial: %s\n", err)
		return nil, err
	}

	return client, nil
}

func loadContract(client *ethclient.Client) (*api.Api, error) {
	contractAddress := common.HexToAddress(os.Getenv("CONTRACT_ADDRESS"))

	// Create a contract instance
	contract, err := api.NewApi(contractAddress, client)
	if err != nil {
		log.Fatalf("newApi: %s\n", err)
		return nil, err
	}

	return contract, nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	loadEnv()

	client, err := loadClient()
	if err != nil {
		log.Fatalf("loadClient: %s", err)
	}

	contract, err := loadContract(client)
	if err != nil {
		log.Fatalf("getCount: %s\n", err)
	}

	privateKeyHex := os.Getenv("PRIVATE_KEY")

	// Create a new signer
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatalf("ChainID: %s", err)
	}

	// Create a new TransactOpts struct
	signer, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("NewKeyedTransactorWithChainID: %s", err)
	}

	handler := CounterHandler{
		api:    contract,
		signer: signer,
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/getCounter", handler.Get)
	e.POST("/increase", handler.Increase)
	e.POST("/reset", handler.Reset)

	err = e.Start(":8080")
	if err != nil {
		log.Fatalf("error starting server: %s", err)
	}
}
