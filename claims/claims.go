package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const registryKey = "CLAIMS-KEY"

// SmartContract provides functions for managing exchange
type SmartContract struct {
	contractapi.Contract // TODO, type aliasing?
}

// InitLedger opens the auction bidding process
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	registry := make(map[cid.ClientIdentity](map[cid.ClientIdentity](map[string]string)))
	registryBytes, err := json.Marshal(registry)
	if err != nil {
		return fmt.Errorf("Error converting registry to bytes: %s", err.Error())
	}
	err = ctx.GetStub().PutState(registryKey, registryBytes)
	if err != nil {
		return fmt.Errorf("Error writing registry bytes to state: %s", err.Error())
	}
	return nil
}

// SetClaim allows a user to create a new claim or modify an existing claim
func (s *SmartContract) SetClaim() error {
	return nil
}

// SetSelfClaim allows a user to create a new claim or modify an existing claim on oneself
func (s *SmartContract) SetSelfClaim() error {
	return nil
}

// RemoveClaim allows a user to remove an existing claim
func (s *SmartContract) RemoveClaim() error {
	return nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error creating Claims chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting Claims chaincode: %s", err.Error())
	}
}
