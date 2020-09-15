package main

import (
	"testing"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func TestLitmusTest(t *testing.T) {
	want := "Hello, world."
	contract := new(SmartContract)
	chaincode, err := contractapi.NewChaincode(contract)
	if got := contract.LitmusTest(); got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}
