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
	P1      cid.ClientIdentity            `json:"p1"`
	P2      cid.ClientIdentity            `json:"p2"`
	P1Play  string                        `json:"p1play"`
	P2Play  string                        `json:"p2play"`
	P1Flag  bool                          `json:"p1flag"`
	P2Flag  bool                          `json:"p2flag"`
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

// UNSET constant for indicating if player choice is not set
const UNSET = "unset"

// SmartContract provides functions for managing exchange
type SmartContract struct {
	contractapi.Contract // TODO, type aliasing?
}

// getGame is a helper function for retrieving RPS game data from transaction history
func (s *SmartContract) getGame(ctx contractapi.TransactionContextInterface, gameID string) (*Game, error) {
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

// putGame is a helper function for committing RPS game data to transaction history
func (s *SmartContract) putGame(ctx contractapi.TransactionContextInterface, game Game) error {
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

// InitLedger initializes the RPS game and allows players to join
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	newGame := Game{
		GameID: uuid.New().String(),
		Status: Open,
		P1Play: UNSET,
		P2Play: UNSET,
		P1Flag: false,
		P2Flag: false,
	}
	err := s.putGame(ctx, newGame)
	if err != nil {
		return err
	}
	return nil
}

// JoinGame allows players to join a game of RPS
func (s *SmartContract) JoinGame(ctx contractapi.TransactionContextInterface, gameID string) error {
	game, err := s.getGame(ctx, gameID)
	if err != nil {
		return err
	}
	if game.Status != Open {
		return fmt.Errorf("Error, game has already started with 2 players")
	}

	if game.P1 == nil {
		game.P1 = ctx.GetClientIdentity()
		game.P1Flag = true
	} else if game.P2 == nil {
		game.P2 = ctx.GetClientIdentity()
		game.P2Flag = true
	}

	if game.P1Flag && game.P2Flag {
		game.Status = ChoosePlay
		game.P1Flag, game.P2Flag = false, false
	}

	err = s.putGame(ctx, *game)
	if err != nil {
		return err
	}
	return nil
}

// MakeChoice allows players to submit a choice that is hashed so as not to be revealed to the opponent
func (s *SmartContract) MakeChoice(ctx contractapi.TransactionContextInterface, gameID string, hash string) error {
	sender := ctx.GetClientIdentity()
	game, err := s.getGame(ctx, gameID)
	if err != nil {
		return err
	}
	if game.Status != ChoosePlay {
		return fmt.Errorf("Error, game is not currently in the make choice phase")
	}
	if game.P1 != sender && game.P2 != sender {
		return fmt.Errorf("Error, player is not registered in this game")
	}

	if game.P1 == sender {
		game.P1Play = hash
	} else if game.P2 == sender {
		game.P2Play = hash
	}

	if game.P1Flag && game.P2Flag {
		game.Status = RevealPlay
		game.P1Flag, game.P2Flag = false, false
	}

	err = s.putGame(ctx, *game)
	if err != nil {
		return err
	}
	return nil
}

// RevealChoice allows players to reveal the players' choices to determine the winner
func (s *SmartContract) RevealChoice(ctx contractapi.TransactionContextInterface, gameID string, choice int, nonce string) error {
	sender := ctx.GetClientIdentity()
	game, err := s.getGame(ctx, gameID)
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
	hasher.Write([]byte(nonce))

	if sender == game.P1 && hex.EncodeToString(hasher.Sum(nil)) == game.P1Play {
		game.P1Play = strconv.Itoa(choice)
		game.P1Flag = true
	} else if sender == game.P2 && hex.EncodeToString(hasher.Sum(nil)) == game.P2Play {
		game.P2Play = strconv.Itoa(choice)
		game.P2Flag = true
	} else {
		return fmt.Errorf("Error, choice + hash did not match the original player input")
	}

	if game.P1Flag && game.P2Flag {
		game.Status = GameOver
		game.P1Flag, game.P2Flag = false, false
	}

	err = s.putGame(ctx, *game)
	if err != nil {
		return err
	}
	return nil
}

// DetermineWinner returns which player won the game
func (s *SmartContract) DetermineWinner(ctx contractapi.TransactionContextInterface, gameID string) (cid.ClientIdentity, error) {
	game, err := s.getGame(ctx, gameID)
	if err != nil {
		return nil, err
	}
	if game.Status != GameOver {
		return nil, fmt.Errorf("Error, game is not currently in the game over phase")
	}

	p1play, p1err := strconv.Atoi(game.P1Play)
	p2play, p2err := strconv.Atoi(game.P2Play)
	p1valid := (p1err != nil) && (p1play >= 0) && (p1play <= 2)
	p2valid := (p2err != nil) && (p2play >= 0) && (p2play <= 2)

	if !p1valid && !p2valid {
		game.Winner = nil
	} else if p1valid && !p2valid {
		game.Winner = game.P1
	} else if !p1valid && p2valid {
		game.Winner = game.P2
	} else {
		if p1play == p2play {
			game.Winner = nil
		} else if (p1play == 0 && p2play == 1) || (p1play == 1 && p2play == 2) || (p1play == 2 && p2play == 0) {
			game.Winner = game.P1
		} else {
			game.Winner = game.P2
		}
	}

	err = s.putGame(ctx, *game)
	if err != nil {
		return nil, err
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
