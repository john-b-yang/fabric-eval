package main

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Game struct for remembering game play information
type Game struct {
	Players map[cid.ClientIdentity]int `json:"players"`
	Status  GameStatus                 `json:"status"`
}

// GameStatus defines enum for game status
type GameStatus int

// Enum definitions typed GameStatus
const (
	Open GameStatus = iota
	ChoosePlay
	RevealPlay
	GameOver
)

// SmartContract provides functions for managing exchange
type SmartContract struct {
	contractapi.Contract // TODO, type aliasing?
}

func (s *SmartContract) GetGameData(ctx contractapi.TransactionContextInterface, gameID string) (*Game, error) {
	gameBytes, err := ctx.GetStub().GetState(gameID)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving game data with ID %s from state: %s", gameID, err)
	}
	if gameBytes == nil {
		return nil, fmt.Errorf("Error, game data with ID %s does not exist", gameID)
	}

	// Unmarshal game data bytes into game data struct
	game := new(Game)
	err = json.Unmarshal(gameBytes, game)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling proposal bytes: %s", err)
	}

	return game, nil
}

// InitLedger opens the auction bidding process
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	newGame := Game{
		Status: Open,
	}
	newGameBytes, err := json.Marshal(newGame)
	if err != nil {
		return fmt.Errorf("Error converting game data to bytes: %s", err.Error())
	}
	err = ctx.GetStub().PutState(uuid.New().String(), newGameBytes)
	if err != nil {
		return fmt.Errorf("Error writing game data bytes to state: %s", err.Error())
	}
	return nil
}

// JoinGame allows players to join a game of RPS
func (s *SmartContract) JoinGame(ctx contractapi.TransactionContextInterface, gameID string) error {
	game, err := s.GetGameData(ctx, gameID)
	if err != nil {
		return err
	}
	if game.Status != Open {
		return fmt.Errorf("Error, game has already started with 2 players")
	}

	game.Players[ctx.GetClientIdentity()] = -1
	return nil
}

func (s *SmartContract) MakeChoice(ctx contractapi.TransactionContextInterface) error {
	return nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error creating RPS game chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting RPS game chaincode: %s", err.Error())
	}
}
