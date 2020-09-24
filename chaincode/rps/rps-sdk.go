package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

// Constants for Transaction Arguments
const GAMEID = "PUT GAME ID HERE AFTER INITLEDGER"
const USER1 = "User1@org1.example.com"
const USER2 = "User0@org1.example.com"
const KEY1 = "0"
const KEY2 = "2"
const NONCE1 = "Jack"
const NONCE2 = "John"
const HASH1 = "B23DE3EA7971B4490AB8175B23BC37E016C9B190F1FD7CB84C2D07BA05EB3050"
const HASH2 = "86F05364C02B1F989218954443BAD8FEDEE1C65B30B4E14E0241D14E20C85880"

func main() {
	// Wallet Creation: Set of user ID's, allows single user to connect to network
	os.Setenv("DISCOVERY_AS_LOCALHOST", "true")

	credPath := filepath.Join("..", "..", "test-network", "organizations",
		"peerOrganizations", "org1.example.com", "users", USER1, "msp")
	ccpPath := filepath.Join("..", "..", "test-network", "organizations",
		"peerOrganizations", "org1.example.com", "connection-org1.yaml")
	user1Contract := addNewUser("appUser1", "Org1MSP", "mychannel", credPath, ccpPath)

	credPath = filepath.Join("..", "..", "test-network", "organizations",
		"peerOrganizations", "org1.example.com", "users", USER2, "msp")
	user2Contract := addNewUser("appUser2", "Org1MSP", "mychannel", credPath, ccpPath)

	result, err := user1Contract.SubmitTransaction("JoinGame", GAMEID)
	if err != nil {
		fmt.Printf("Failed to submit transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))

	result, err = user2Contract.SubmitTransaction("JoinGame", GAMEID)
	if err != nil {
		fmt.Printf("Failed to submit transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))

	result, err = user1Contract.SubmitTransaction("MakeChoice", GAMEID, HASH1)
	if err != nil {
		fmt.Printf("Failed to submit transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))

	result, err = user2Contract.SubmitTransaction("MakeChoice", GAMEID, HASH2)
	if err != nil {
		fmt.Printf("Failed to submit transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))

	result, err = user1Contract.SubmitTransaction("RevealChoice", GAMEID, KEY1, NONCE1)
	if err != nil {
		fmt.Printf("Failed to submit transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))

	result, err = user2Contract.SubmitTransaction("RevealChoice", GAMEID, KEY2, NONCE2)
	if err != nil {
		fmt.Printf("Failed to submit transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))

	result, err = user1Contract.SubmitTransaction("DetermineWinner", GAMEID)
	if err != nil {
		fmt.Printf("Failed to submit transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))
}

func addNewUser(username string, mspID string, contract string, credPath string, ccpPath string) *gateway.Contract {
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		fmt.Printf("Failed to create wallet: %s\n", err)
		os.Exit(1)
	}

	if !wallet.Exists(username) {
		err = populateWallet(wallet, username, mspID, credPath)
		if err != nil {
			fmt.Printf("Failed to populate wallet contents: %s\n", err)
			os.Exit(1)
		}
	}

	gw, err := gateway.Connect(gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))), gateway.WithIdentity(wallet, username))
	if err != nil {
		fmt.Printf("Failed to connect to gateway: %s\n", err)
		os.Exit(1)
	}
	defer gw.Close()

	network, err := gw.GetNetwork(contract)
	if err != nil {
		fmt.Printf("Failed to get network: %s\n", err)
		os.Exit(1)
	}

	return network.GetContract(contract)
}

func populateWallet(wallet *gateway.Wallet, username string, mspID string, credPath string) error {
	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}
	keyDir := filepath.Join(credPath, "keystore")
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return errors.New("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}
	identity := gateway.NewX509Identity(mspID, string(cert), string(key))
	err = wallet.Put(username, identity)
	if err != nil {
		return err
	}
	return nil
}
