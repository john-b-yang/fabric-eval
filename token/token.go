package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Token struct for storing token information
type Token struct {
	Allowances map[cid.ClientIdentity](map[cid.ClientIdentity]int) `json:"allowances"`
	Balances   map[cid.ClientIdentity]int                          `json:"balances"`
	Name       string                                              `json:"name"`
	Symbol     string                                              `json:"symbol"`
	Supply     int                                                 `json:"supply"`
}

// SmartContract provides functions for managing exchange
type SmartContract struct {
	contractapi.Contract // TODO, type aliasing?
}

// InitLedger opens the auction bidding process
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface, name string, symbol string, supply int) error {
	var balances map[cid.ClientIdentity]int
	balances[ctx.GetClientIdentity()] = supply

	token := Token{
		Balances: balances,
		Name:     name,
		Symbol:   symbol,
		Supply:   supply,
	}

	tokenBytes, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("Error converting token to bytes: %s", err.Error())
	}
	err = ctx.GetStub().PutState(token.Symbol, tokenBytes)
	if err != nil {
		return fmt.Errorf("Error writing token bytes to state: %s", err.Error())
	}
	return nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error creating token chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting token chaincode: %s", err.Error())
	}
}
