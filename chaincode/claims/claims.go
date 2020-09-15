package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const claimsKey = "CLAIMS-KEY"

// Registry type for storing claims
type Registry = map[cid.ClientIdentity](map[cid.ClientIdentity](map[string]string))

// SmartContract provides functions for managing exchange
type SmartContract struct {
	contractapi.Contract // TODO, type aliasing?
}

// Claims struct for storing claim information
type Claims struct {
	registry Registry
}

// InitLedger opens the auction bidding process
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	claims := Claims{
		registry: make(Registry),
	}
	claimBytes, err := json.Marshal(claims)
	if err != nil {
		return fmt.Errorf("Error converting registry to bytes: %s", err.Error())
	}
	err = ctx.GetStub().PutState(claimsKey, claimBytes)
	if err != nil {
		return fmt.Errorf("Error writing registry bytes to state: %s", err.Error())
	}
	return nil
}

// SetClaim allows a user to create a new claim or modify an existing claim
func (s *SmartContract) SetClaim(ctx contractapi.TransactionContextInterface, subject cid.ClientIdentity, key string, value string) error {
	claims, err := s.getClaims(ctx)
	if err != nil {
		return err
	}
	sender := ctx.GetClientIdentity()
	claims.registry[sender][subject][key] = value
	err = s.putClaims(ctx, *claims)
	if err != nil {
		return err
	}
	return nil
}

// SetSelfClaim allows a user to create a new claim or modify an existing claim on oneself
func (s *SmartContract) SetSelfClaim(ctx contractapi.TransactionContextInterface, key string, value string) error {
	claims, err := s.getClaims(ctx)
	if err != nil {
		return err
	}
	sender := ctx.GetClientIdentity()
	claims.registry[sender][sender][key] = value
	err = s.putClaims(ctx, *claims)
	if err != nil {
		return err
	}
	return nil
}

// RemoveClaim allows a user to remove an existing claim that it has issued or been issued
func (s *SmartContract) RemoveClaim(ctx contractapi.TransactionContextInterface, issuer cid.ClientIdentity, subject cid.ClientIdentity, key string) error {
	claims, err := s.getClaims(ctx)
	if err != nil {
		return err
	}
	sender := ctx.GetClientIdentity()
	if sender != issuer && sender != subject {
		return fmt.Errorf("Error, sender does not have authority to alter this claim")
	}
	if _, ok := claims.registry[issuer][subject][key]; !ok {
		return fmt.Errorf("Error, specified claim does not exist")
	}
	delete(claims.registry[issuer][subject], key)
	err = s.putClaims(ctx, *claims)
	if err != nil {
		return err
	}
	return nil
}

func (s *SmartContract) getClaims(ctx contractapi.TransactionContextInterface) (*Claims, error) {
	claimBytes, err := ctx.GetStub().GetState(claimsKey)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Registry from state: %s", err)
	}
	if claimBytes == nil {
		return nil, fmt.Errorf("Error, Registry has not been created")
	}

	// Unmarshal game data bytes into game data struct
	claims := new(Claims)
	err = json.Unmarshal(claimBytes, claims)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling Registry bytes: %s", err)
	}

	return claims, nil
}

func (s *SmartContract) putClaims(ctx contractapi.TransactionContextInterface, claims Claims) error {
	claimBytes, err := json.Marshal(claims)
	if err != nil {
		return fmt.Errorf("Error converting Registry to bytes: %s", err.Error())
	}
	err = ctx.GetStub().PutState(claimsKey, claimBytes)
	if err != nil {
		return fmt.Errorf("Error writing Registry bytes to state: %s", err.Error())
	}
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
