package main

import (
	"fmt"
	"github.com/camilovietnam/solidity/01/api"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type CounterHandler struct {
	api    *api.Api
	signer *bind.TransactOpts
}

type HTTPResponse struct {
	Message string `json:"message"`
}

var HTTPInternal = echo.NewHTTPError(http.StatusInternalServerError, "Internal Error. See the logs.")

func (h *CounterHandler) Get(c echo.Context) error {
	counter, err := h.api.GetCount(&bind.CallOpts{})
	if err != nil {
		log.Printf("GetCount: %s\n", err)
		return HTTPInternal
	}

	return c.JSON(http.StatusOK, HTTPResponse{
		Message: fmt.Sprintf("Counter: %d", counter),
	})
}

func (h *CounterHandler) Increase(c echo.Context) error {
	tx, err := h.api.Increase(h.signer)
	if err != nil {
		log.Printf("Increase: %s", err)
		return HTTPInternal
	}

	return c.JSON(http.StatusOK, HTTPResponse{Message: fmt.Sprintf("Transaction: https://sepolia.etherscan.io/tx/%s", tx.Hash().Hex())})
}

func (h *CounterHandler) Reset(c echo.Context) error {
	tx, err := h.api.Reset(h.signer)
	if err != nil {
		log.Printf("Reset: %s\n", err)
		return HTTPInternal
	}

	return c.JSON(http.StatusOK, HTTPResponse{Message: fmt.Sprintf("Transaction: https://sepolia.etherscan.io/tx/%s", tx.Hash().Hex())})
}
