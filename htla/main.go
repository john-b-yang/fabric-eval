package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// Define the smart contract object
type HashTimeLockContract struct {
}

// Instantiation
func (s *HashTimeLockContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// Handle Chaincode Operations
func (s *HashTimeLockContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	switch function {
	case "createProposal":
		return s.createProposal(stub, args)
	case "confirmProposal":
		return s.confirmProposal(stub, args)
	case "invalidateProposal":
		return s.invalidateProposal(stub, args)
	default:
		return shim.Error("Invalid Smart Contract function name.")
	}
}

/*
 * Takes a proposal and a hash => entry tagged as PENDING
 * Returns a proposal id
 */
func (s *HashTimeLockContract) createProposal(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var err error

	if len(args) != 3 {
		return shim.Error("Invalid arguments to createProposal, expected proposal, hash, hashingAlg.")
	}

	validAlg := false
	for _, a := range validHashingAlgorithms {
		if a == args[2] {
			validAlg = true
		}
	}
	if validAlg == false {
		return shim.Error("Only these hashing algorithms are supported: " + strings.Join(validHashingAlgorithms, ", "))
	}

	proposal := proposalEntry{Proposal: abstractProposal{}, Status: PendingStatus, Hash: args[1], HashAlgorithm: args[2]}
	err = json.Unmarshal([]byte(args[0]), &proposal.Proposal)
	if err != nil {
		return shim.Error("Error parsing provided proposal definition - " + err.Error())
	}
	if proposal.Proposal.ProposalID == "" {
		return shim.Error("No proposalId provided as part of proposal.")
	}
	if proposal.Proposal.Handler == "" {
		return shim.Error("No proposalHandler provided as part of proposal.")
	}

	// Validation Logic would go here

	// Write the proposal to state
	proposalAsBytes, err := json.Marshal(proposal)
	if err != nil {
		return shim.Error("Error building proposal definition - " + err.Error())
	}
	err = stub.PutState(proposalPrefix+proposal.Proposal.ProposalID, proposalAsBytes)
	if err != nil {
		return shim.Error("Error writing proposal to state - " + err.Error())
	}

	// Fire Events
	proposalCreatedEvent := ProposalCreatedEventObject{ProposalID: proposal.Proposal.ProposalID}
	proposalEventAsBytes, err := json.Marshal(proposalCreatedEvent)
	if err != nil {
		return shim.Error("Error building proposal event definition - " + err.Error())
	}

	// Provided Handler Event
	err = stub.SetEvent(proposal.Proposal.Handler+ProposalCreatedHandlerEvent, proposalEventAsBytes)
	// Timeout Client EVent
	err = stub.SetEvent(ProposalCreateTimeoutEvent, proposalEventAsBytes)
	return shim.Success(nil)
}

/*
 * Transition a proposal from PENDING to CONFIRMED
 * Takes a proposalId and a pre-image that corresponds with the initially provided
 * hash.
 */
func (s *HashTimeLockContract) confirmProposal(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var err error
	// ProposalID, Pre-image expected
	if len(args) != 2 {
		return shim.Error("Invalid arguments to confirmProposal, expected proposalId, pre-image.")
	}

	// Retreive the proposal referenced
	proposalAsBytes, err := stub.GetState(proposalPrefix + args[0])
	if err != nil {
		return shim.Error("Error while retreiving the stored proposal from state - " + err.Error())
	}
	if proposalAsBytes == nil {
		return shim.Error("No such proposal. It may have expired and been invalidated.")
	}
	proposal := proposalEntry{}
	err = json.Unmarshal(proposalAsBytes, &proposal)
	if err != nil {
		return shim.Error("Error while parsing the proposal stored in state - " + err.Error())
	}

	// More Validation Logic

	// TODO: Time Lock Logic

	// Validate Pre-Image by hashing algorithm
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
		return shim.Error("The hash algorithm which was recorded in the proposal is not supported.")
	}
	hasher.Write([]byte(args[1]))
	if hex.EncodeToString(hasher.Sum(nil)) != strings.ToLower(proposal.Hash) {
		return shim.Error("Invalid Pre-image supplied.")
	}

	//Mark the proposal as confirmed
	proposal.Status = ConfirmStatus
	proposalAsBytes, err = json.Marshal(proposal)
	if err != nil {
		return shim.Error("Error when marshaling proposal - " + err.Error())
	}
	err = stub.PutState(proposalPrefix+args[0], proposalAsBytes)

	// Event - Allow Repayment to Account
	proposalConfirmedEvent := ProposalConfirmedEventObject{ProposalID: args[0], PreImage: args[1]}
	proposalEventAsBytes, err := json.Marshal(proposalConfirmedEvent)
	if err != nil {
		return shim.Error("Error building proposal event definition - " + err.Error())
	}
	err = stub.SetEvent(ProposalConfirmedHandlerEvent, proposalEventAsBytes)
	return shim.Success(nil)
}

/*
 * Invalidate a proposal in PENDING state. (timelocking)
 *
 * Fails if invoked on a CONFIRMED proposal.
 */
func (s *HashTimeLockContract) invalidateProposal(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var err error
	// Expects Proposal ID
	if len(args) != 1 {
		return shim.Error("Invalid arguments to invalidateProposal, expected proposalId")
	}

	// Validation logic here

	proposalBytes, err := stub.GetState(proposalPrefix + args[0])
	if err != nil {
		return shim.Error("Error retreiving stored proposal from state")
	}
	proposal := proposalEntry{}
	err = json.Unmarshal(proposalBytes, &proposal)
	if err != nil {
		return shim.Error("Error while parsing the proposal stored in state - " + err.Error())
	}
	if proposal.Status != PendingStatus {
		return shim.Error("Only pending proposals can be timed out.")
	}

	// TODO: Time Lock Logic

	// Delete the proposal
	err = stub.DelState(proposalPrefix + args[0])
	if err != nil {
		return shim.Error("Error while deleting proposal from state.")
	}
	return shim.Success(nil)
}

func main() {
	// Create a new Smart Contract
	err := shim.Start(new(HashTimeLockContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
