package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const tokenKey = "TOKEN-KEY"

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
	err = ctx.GetStub().PutState(tokenKey, tokenBytes)
	if err != nil {
		return fmt.Errorf("Error writing token bytes to state: %s", err.Error())
	}
	return nil
}

// getToken is a helper function to retrieve the Token data
func (s *SmartContract) getToken(ctx contractapi.TransactionContextInterface) (*Token, error) {
	tokenBytes, err := ctx.GetStub().GetState(tokenKey)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Token from state: %s", err)
	}
	if tokenBytes == nil {
		return nil, fmt.Errorf("Error, Token has not been created")
	}

	// Unmarshal game data bytes into game data struct
	token := new(Token)
	err = json.Unmarshal(tokenBytes, token)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling proposal bytes: %s", err)
	}

	return token, nil
}

// Transfer allows a sender to transfer tokens to a receiver
func (s *SmartContract) Transfer(ctx contractapi.TransactionContextInterface, receiver cid.ClientIdentity, numTokens int) error {
	token, err := s.getToken(ctx)
	if err != nil {
		return err
	}
	sender := ctx.GetClientIdentity()
	if numTokens > token.Balances[sender] {
		return fmt.Errorf("Sender does not have enough balance to transfer desired amount")
	}
	token.Balances[sender] = token.Balances[sender] - numTokens
	token.Balances[receiver] = token.Balances[receiver] + numTokens
	tokenBytes, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("Error converting Token to bytes: %s", err.Error())
	}
	err = ctx.GetStub().PutState(tokenKey, tokenBytes)
	if err != nil {
		return fmt.Errorf("Error writing Token bytes to state: %s", err.Error())
	}
	return nil
}

// Approve allows a sender to identify a delegate who can spend funds on their behalf
func (s *SmartContract) Approve(ctx contractapi.TransactionContextInterface, delegate cid.ClientIdentity, numTokens int) error {
	token, err := s.getToken(ctx)
	if err != nil {
		return err
	}
	sender := ctx.GetClientIdentity()
	if numTokens > token.Balances[sender] {
		return fmt.Errorf("Sender does not have enough balance to transfer desired amount")
	}
	if _, ok := token.Allowances[sender][delegate]; ok {
		return fmt.Errorf("Allowance for specified delegate from sender has already been created")
	}
	token.Balances[sender] = token.Balances[sender] - numTokens
	token.Allowances[sender][delegate] = numTokens
	tokenBytes, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("Error converting Token to bytes: %s", err.Error())
	}
	err = ctx.GetStub().PutState(tokenKey, tokenBytes)
	if err != nil {
		return fmt.Errorf("Error writing Token bytes to state: %s", err.Error())
	}
	return nil
}

// IncreaseAllowance allows a sender to increase the funds of a delegate
func (s *SmartContract) IncreaseAllowance(ctx contractapi.TransactionContextInterface, delegate cid.ClientIdentity, numTokens int) error {
	token, err := s.getToken(ctx)
	if err != nil {
		return err
	}
	sender := ctx.GetClientIdentity()
	if numTokens > token.Balances[sender] {
		return fmt.Errorf("Sender does not have enough balance to transfer desired amount")
	}
	if _, ok := token.Allowances[sender][delegate]; !ok {
		return fmt.Errorf("Allowance for specified delegate from sender has not been created")
	}
	token.Balances[sender] = token.Balances[sender] - numTokens
	token.Allowances[sender][delegate] = token.Allowances[sender][delegate] + numTokens
	tokenBytes, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("Error converting Token to bytes: %s", err.Error())
	}
	err = ctx.GetStub().PutState(tokenKey, tokenBytes)
	if err != nil {
		return fmt.Errorf("Error writing Token bytes to state: %s", err.Error())
	}
	return nil
}

// DecreaseAllowance allows a sender to decrease the funds of a delegate
func (s *SmartContract) DecreaseAllowance(ctx contractapi.TransactionContextInterface, delegate cid.ClientIdentity, numTokens int) error {
	token, err := s.getToken(ctx)
	if err != nil {
		return err
	}
	sender := ctx.GetClientIdentity()
	if _, ok := token.Allowances[sender][delegate]; !ok {
		return fmt.Errorf("Allowance for specified delegate from sender has not been created")
	}
	if token.Allowances[sender][delegate]-numTokens < 0 {
		return fmt.Errorf("Requested tokens for transfer exceeds the existing allowance")
	}

	token.Balances[sender] = token.Balances[sender] + numTokens
	token.Allowances[sender][delegate] = token.Allowances[sender][delegate] - numTokens
	tokenBytes, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("Error converting Token to bytes: %s", err.Error())
	}
	err = ctx.GetStub().PutState(tokenKey, tokenBytes)
	if err != nil {
		return fmt.Errorf("Error writing Token bytes to state: %s", err.Error())
	}
	return nil
}

// TransferFrom allows a delegate to spend the given allowance from an owner to a desired receiver
func (s *SmartContract) TransferFrom(ctx contractapi.TransactionContextInterface, owner cid.ClientIdentity, buyer cid.ClientIdentity, numTokens int) error {
	token, err := s.getToken(ctx)
	if err != nil {
		return err
	}
	sender := ctx.GetClientIdentity()
	if _, ok := token.Allowances[owner][sender]; ok {
		if numTokens > token.Allowances[owner][sender] {
			return fmt.Errorf("Requested tokens for transfer exceeds the existing allowance")
		}
	} else {
		return fmt.Errorf("Allowance for specified delegate from sender has not been created")
	}
	token.Allowances[owner][sender] = token.Allowances[owner][sender] - numTokens
	token.Balances[buyer] = token.Balances[buyer] + numTokens
	tokenBytes, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("Error converting Token to bytes: %s", err.Error())
	}
	err = ctx.GetStub().PutState(tokenKey, tokenBytes)
	if err != nil {
		return fmt.Errorf("Error writing Token bytes to state: %s", err.Error())
	}
	return nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error creating Token chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting Token chaincode: %s", err.Error())
	}
}
