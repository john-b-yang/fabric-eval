package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Instance fields
var counter int
var validHashingAlgorithms = []string{"SHA256", "SHA384", "SHA512"}

// ProposalStatus defines enum for proposal status
type ProposalStatus int

// Enum definitions typed Proposal Status
const (
	Pending ProposalStatus = iota
	Confirmed
	Expired
)

// Proposal struct for hlta exchange information
type Proposal struct {
	Amount        int            `json:"amount"`
	Hash          string         `json:"hash"`
	HashAlgorithm string         `json:"hashAlgorithm"`
	Status        ProposalStatus `json:"status"`
	Timelock      time.Time      `json:"timelock"`
}

// SmartContract provides functions for managing exchange
type SmartContract struct {
	contractapi.Contract // TODO, type aliasing?
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Placeholder string // TODO, is struct required for returns
}

// InitLedger initializes the htla exchange
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}

// LitmusTest is an easy method to check if the contract is initialized properly
func (s *SmartContract) LitmusTest(ctx contractapi.TransactionContextInterface) string {
	return "Hello, world."
}

// GenerateHash returns a hash string from an input string and a specified hashing function
func (s *SmartContract) GenerateHash(ctx contractapi.TransactionContextInterface, input string, hashAlgorithm string) (string, error) {
	var hasher hash.Hash
	switch hashAlgorithm {
	case "SHA256":
		hasher = sha256.New()
		break
	case "SHA384":
		hasher = sha512.New384()
		break
	case "SHA512":
		hasher = sha512.New()
		break
	default:
		return "", fmt.Errorf("Error, hash algorithm recorded in proposal is unsupported")
	}
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// CreateProposal takes a proposal and a hash => entry tagged as PENDING
// Returns a proposal id
func (s *SmartContract) CreateProposal(ctx contractapi.TransactionContextInterface, tokens int, timelock int, hash string, hashAlgorithm string) error {
	// Check if hashing algorithm is supported
	var isValidHash = false
	for _, a := range validHashingAlgorithms {
		if a == hashAlgorithm {
			isValidHash = true
		}
	}
	if !isValidHash {
		return fmt.Errorf("Provided hashing algorithm is unsupported")
	}

	// Check if `tokens` is non-zero
	if tokens == 0 {
		return fmt.Errorf("Tokens in transaction cannot be 0")
	}

	// Check if `timelock` is non-zero
	if timelock == 0 {
		return fmt.Errorf("Duration of timelock cannot be 0")
	}

	// Create new proposal struct
	newProposal := Proposal{
		Amount:        tokens,
		Hash:          hash,
		HashAlgorithm: hashAlgorithm,
		Status:        Pending,
		Timelock:      (time.Now().Add(time.Duration(timelock) * time.Second)),
	}

	// Convert Proposal to Bytes and Put into State
	newProposalBytes, err := json.Marshal(newProposal)
	if err != nil {
		return fmt.Errorf("Error converting proposal to bytes: %s", err.Error())
	}
	err = ctx.GetStub().PutState(strconv.Itoa(counter), newProposalBytes)
	if err != nil {
		return fmt.Errorf("Error writing proposal bytes to state: %s", err.Error())
	}
	counter++

	return nil
}

// GetProposal retrieves the Proposal struct corresponding to the given Proposal ID
func (s *SmartContract) GetProposal(ctx contractapi.TransactionContextInterface, proposalID int) (*Proposal, error) {
	// Retrieving proposal from state hashmap
	proposalBytes, err := ctx.GetStub().GetState(strconv.Itoa(proposalID))
	if err != nil {
		return nil, fmt.Errorf("Error retrieving proposal with ID %d from state: %s", proposalID, err)
	}
	if proposalBytes == nil {
		return nil, fmt.Errorf("Error, proposal with ID %d does not exist", proposalID)
	}

	// Unmarshal proposal bytes into proposal struct
	proposal := new(Proposal)
	err = json.Unmarshal(proposalBytes, proposal)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling proposal bytes: %s", err)
	}

	return proposal, nil
}

// ConfirmProposal transitions a proposal from PENDING to CONFIRMED
// Takes a proposalId and a pre-image that corresponds with the initially provided hash.
func (s *SmartContract) ConfirmProposal(ctx contractapi.TransactionContextInterface, proposalID int, preImage string) error {
	proposal, err := s.GetProposal(ctx, proposalID)

	// Check if Proposal can be confirmed
	if proposal.Status != Pending {
		return fmt.Errorf("Error, proposal has already been resolved")
	}

	if time.Now().After(proposal.Timelock) {
		return fmt.Errorf("Error, timelock has expired, proposal cannot be confirmed")
	}

	// Check if pre image matches original hash
	var hasher hash.Hash
	switch proposal.HashAlgorithm {
	case "SHA256":
		hasher = sha256.New()
		break
	case "SHA384":
		hasher = sha512.New384()
		break
	case "SHA512":
		hasher = sha512.New()
		break
	default:
		return fmt.Errorf("Error, hash algorithm recorded in proposal is unsupported")
	}
	hasher.Write([]byte(preImage))

	if hex.EncodeToString(hasher.Sum(nil)) != strings.ToLower(proposal.Hash) {
		return fmt.Errorf("Error, given pre-image does not match the original hash")
	}

	// Update Proposal status to confirmed and Put into State
	proposal.Status = Confirmed
	proposalBytes, err := json.Marshal(proposal)
	if err != nil {
		return fmt.Errorf("Error converting proposal to bytes: %s", err.Error())
	}
	err = ctx.GetStub().PutState(strconv.Itoa(proposalID), proposalBytes)
	if err != nil {
		return fmt.Errorf("Error writing proposal bytes to state: %s", err.Error())
	}

	return nil
}

// InvalidateProposal a proposal in PENDING state. (timelocking)
// Fails if invoked on a CONFIRMED proposal.
func (s *SmartContract) InvalidateProposal(ctx contractapi.TransactionContextInterface, proposalID int) error {
	proposal, err := s.GetProposal(ctx, proposalID)

	// Check if Proposal can be invalidated
	if proposal.Status != Pending {
		return fmt.Errorf("Error, proposal has already been resolved")
	}

	if time.Now().Before(proposal.Timelock) {
		return fmt.Errorf("Error, timelock has not expired, proposal cannot be invalidated")
	}

	// Update Proposal status to expired and Put into State
	proposal.Status = Expired
	proposalBytes, err := json.Marshal(proposal)
	if err != nil {
		return fmt.Errorf("Error converting proposal to bytes: %s", err.Error())
	}
	err = ctx.GetStub().PutState(strconv.Itoa(proposalID), proposalBytes)
	if err != nil {
		return fmt.Errorf("Error writing proposal bytes to state: %s", err.Error())
	}

	return nil
}

func main() {
	// TODO: Key collisions for 2 different contracts committing same key values
	counter = 0
	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error creating htla chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting htla chaincode: %s", err.Error())
	}
}
