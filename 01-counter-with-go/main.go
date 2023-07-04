package main

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"os"

	"github.com/camilovietnam/solidity/01/api"
	"github.com/ethereum/go-ethereum/ethclient"
)

func loadClient() (*ethclient.Client, error) {
	sepoliaURL := os.Getenv("SEPOLIA_URL")

	client, err := ethclient.Dial(sepoliaURL)
	if err != nil {
		log.Fatalf("dial: %s\n", err)
		return nil, err
	}

	return client, nil
}

func loadContract(client *ethclient.Client) (*api.Api, error) {
	contractAddress := common.HexToAddress(os.Getenv("CONTRACT_ADDRESS"))

	contract, err := api.NewApi(contractAddress, client)
	if err != nil {
		log.Fatalf("newApi: %s\n", err)
		return nil, err
	}

	return contract, nil
}

func loadSigner(ctx context.Context, client *ethclient.Client) (*bind.TransactOpts, error) {
	privateKeyHex := os.Getenv("PRIVATE_KEY")

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("HEXToECDSA: %s\n", err)
		return nil, err
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatalf("ChainID: %s", err)
		return nil, err
	}

	// Create a new TransactOpts struct
	signer, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("NewKeyedTransactorWithChainID: %s", err)
		return nil, err
	}

	return signer, nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := loadClient()
	if err != nil {
		return
	}

	contract, err := loadContract(client)
	if err != nil {
		return
	}

	// Create a new signer
	signer, err := loadSigner(ctx, client)
	if err != nil {
		return
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
