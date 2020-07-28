package main

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Instance fields
var deadline time.Time
var highestBid int
var highestBidder cid.ClientIdentity
var seller cid.ClientIdentity

// SmartContract provides functions for managing exchange
type SmartContract struct {
	contractapi.Contract // TODO, type aliasing?
}

func (s *SmartContract) getCurrentTime(ctx contractapi.TransactionContextInterface) (time.Time, error) {
	xactTimestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return time.Time{}, fmt.Errorf("Error retrieving transaction timestamp: %s", err)
	}

	xactTime, err := ptypes.Timestamp(xactTimestamp)
	if err != nil {
		return time.Time{}, fmt.Errorf("Error converting transaction timestamp to time.Time object")
	}

	return xactTime, nil
}

// InitLedger opens the auction bidding process
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface, duration int) error {
	if duration <= 0 {
		return fmt.Errorf("Duration of auction must be positive integer")
	}

	xactTime, err := s.getCurrentTime(ctx)
	if err != nil {
		return err
	}

	deadline = xactTime.Add(time.Duration(duration) * time.Second)
	highestBid = 0
	seller = ctx.GetClientIdentity()
	highestBidder = seller

	return nil
}

// SubmitBid allows one to submit a bid
func (s *SmartContract) SubmitBid(ctx contractapi.TransactionContextInterface, bid int) error {
	xactTime, err := s.getCurrentTime(ctx)
	if err != nil {
		return err
	}

	if xactTime.After(deadline) {
		return fmt.Errorf("Auction duration has already expired, no more bidding allowed")
	}

	if bid <= highestBid {
		return fmt.Errorf("Bid does not exceed current highest bid")
	}

	highestBid = bid
	highestBidder = ctx.GetClientIdentity()

	return nil
}

// CloseBid ends the bidding period
func (s *SmartContract) CloseBid(ctx contractapi.TransactionContextInterface) error {
	xactTime, err := s.getCurrentTime(ctx)
	if err != nil {
		return err
	}

	if xactTime.Before(deadline) {
		return fmt.Errorf("Auction duration has not expired, cannot be closed")
	}

	// TODO: Exchange funds + items at this point

	return nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error creating auction chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting auction chaincode: %s", err.Error())
	}
}
