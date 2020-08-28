package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Game struct for remembering game play information
type Game struct {
	GameID  string                        `json:"gameid"`
	Players map[cid.ClientIdentity]string `json:"players"`
	Status  GameStatus                    `json:"status"`
	Winner  cid.ClientIdentity            `json:"winner"`
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

const UNSET = "unset"

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
		GameID: uuid.New().String(),
		Status: Open,
	}
	newGameBytes, err := json.Marshal(newGame)
	if err != nil {
		return fmt.Errorf("Error converting game data to bytes: %s", err.Error())
	}
	err = ctx.GetStub().PutState(newGame.GameID, newGameBytes)
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

	game.Players[ctx.GetClientIdentity()] = UNSET

	if len(game.Players) >= 3 {
		game.Status = ChoosePlay
	}

	gameBytes, err := json.Marshal(game)
	if err != nil {
		return fmt.Errorf("Error converting game data to bytes: %s", err.Error())
	}
	err = ctx.GetStub().PutState(game.GameID, gameBytes)
	if err != nil {
		return fmt.Errorf("Error writing proposal bytes to state: %s", err.Error())
	}
	return nil
}

// MakeChoice allows players to submit a choice that is hashed so as not to be revealed to the opponent
func (s *SmartContract) MakeChoice(ctx contractapi.TransactionContextInterface, gameID string, choice int, hash string) error {
	if choice < 0 || choice > 2 {
		return fmt.Errorf("Error, invalid choice: must be 0 (rock), 1 (paper), or 2 (scissors)")
	}
	game, err := s.GetGameData(ctx, gameID)
	if err != nil {
		return err
	}
	if game.Status != ChoosePlay {
		return fmt.Errorf("Error, game is not currently in the make choice phase")
	}
	if _, ok := game.Players[ctx.GetClientIdentity()]; !ok {
		return fmt.Errorf("Error, player is not registered in this game")
	}

	hasher := sha256.New()
	hasher.Write([]byte(strconv.Itoa(choice)))
	hasher.Write([]byte(hash))
	game.Players[ctx.GetClientIdentity()] = hex.EncodeToString(hasher.Sum(nil))

	var choicesSet = true
	for _, val := range game.Players {
		if val != UNSET {
			choicesSet = false
		}
	}
	if choicesSet {
		game.Status = RevealPlay
	}

	gameBytes, err := json.Marshal(game)
	if err != nil {
		return fmt.Errorf("Error converting game data to bytes: %s", err.Error())
	}
	err = ctx.GetStub().PutState(game.GameID, gameBytes)
	if err != nil {
		return fmt.Errorf("Error writing proposal bytes to state: %s", err.Error())
	}
	return nil
}

// RevealChoice allows players to reveal the players' choices to determine the winner
func (s *SmartContract) RevealChoice(ctx contractapi.TransactionContextInterface, gameID string, choice int, hash string) error {
	game, err := s.GetGameData(ctx, gameID)
	if err != nil {
		return err
	}
	if game.Status != ChoosePlay {
		return fmt.Errorf("Error, game is not currently in the reveal choice phase")
	}
	if _, ok := game.Players[ctx.GetClientIdentity()]; !ok {
		return fmt.Errorf("Error, player is not registered in this game")
	}

	hasher := sha256.New()
	hasher.Write([]byte(strconv.Itoa(choice)))
	hasher.Write([]byte(hash))
	if hex.EncodeToString(hasher.Sum(nil)) == game.Players[ctx.GetClientIdentity()] {
		game.Players[ctx.GetClientIdentity()] = strconv.Itoa(choice)
	} else {
		return fmt.Errorf("Error, choice + hash did not match the original player input")
	}

	var choicesRevealed = true
	for _, val := range game.Players {
		i, err := strconv.Atoi(val)
		if err != nil || (err == nil && (i > 2 || i < 0)) {
			choicesRevealed = false
		}
	}
	if choicesRevealed {
		game.Status = GameOver
	}

	gameBytes, err := json.Marshal(game)
	if err != nil {
		return fmt.Errorf("Error converting game data to bytes: %s", err.Error())
	}
	err = ctx.GetStub().PutState(game.GameID, gameBytes)
	if err != nil {
		return fmt.Errorf("Error writing proposal bytes to state: %s", err.Error())
	}
	return nil
}

// DetermineWinner returns which player won the game
func (s *SmartContract) DetermineWinner(ctx contractapi.TransactionContextInterface, gameID string) (cid.ClientIdentity, error) {
	game, err := s.GetGameData(ctx, gameID)
	if err != nil {
		return nil, err
	}
	if game.Status != GameOver {
		return nil, fmt.Errorf("Error, game is not currently in the game over phase")
	}

	var p1, p2 cid.ClientIdentity
	var p1play, p2play, index int
	for key, val := range game.Players {
		if index == 0 {
			p1 = key
			p1play, _ = strconv.Atoi(val)
		} else if index == 1 {
			p2 = key
			p2play, _ = strconv.Atoi(val)
		}
		index++
	}

	var result = (3 + p1play - p2play)
	if result == 1 || result == 4 {
		game.Winner = p1
	} else if result == 2 || result == 5 {
		game.Winner = p2
	} else {
		game.Winner = nil
	}

	gameBytes, err := json.Marshal(game)
	if err != nil {
		return nil, fmt.Errorf("Error converting game data to bytes: %s", err.Error())
	}
	err = ctx.GetStub().PutState(game.GameID, gameBytes)
	if err != nil {
		return nil, fmt.Errorf("Error writing proposal bytes to state: %s", err.Error())
	}
	return game.Winner, nil
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
