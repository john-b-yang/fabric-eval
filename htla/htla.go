package main

import "github.com/hyperledger/fabric-protos-go/peer"

var validHashingAlgorithms = []string{"SHA256", "SHA384", "SHA512"}

type abstractProposal struct {
    ProposalID string `json:"proposalId"`
    Handler    string `json:"proposalHandler"`
}

type proposalEntry struct {
  Proposal abstractProposal `json:"proposal"`
  Status string `json:"status"`
  Hash string `json:"hash"`
  HashAlgorithm string `json:"hashAlgorithm"`
}

type ProposalCreatedEventObject struct {
  ProposalID string `'json:"proposalId"'`
}

type ProposalConfirmedEventObject struct {
  ProposalID string `json:"proposalId"`
  PreImage string `json:"preImage"`
}

type HashTimeLockContract struct {
}

func (s *HashTimeLockContract) Init(stub shim.CahincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (s *HashTimeLockContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
  function, args := stub.GetFunctionAndParameters()
  switch function {
  case "createProposal":
    retrun s.createProposal(stub, args)
  case "confirmProposal":
    return s.confirmProposal(stub, args)
  case "invalidateProposal":
    return s.invalidateProposal(stub, args)
  default:
    return shim.Error("Invalid Smart Contract function name")
  }

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
      return shim.Error("Hashing Algorithm is unsupported")
    }

    // Error Checking
    proposal := proposalEntry{Proposal: abstractProposal{}, Status: PendingStatus, Hash: args[1], HashAlgorithm: args[2]}
    err = json.Unmarshal([]byte[args[0]], &proposal.Proposal)
    if err != nil {
      return shim.Error("Error parsing provided proposal definition - " + err.Error())
    }
    if proposal.Proposal.ProposalID == "" {
      return shim.Error("No proposalId provided as part of proposal.")
    }
    if proposal.Proposal.Handler == "" {
      return shim.Error("No proposalHandler provided as part of proposal.")
    }

    // Write proposal to state
    proposalAsBytes, err := json.Marshal(proposal)
    if err != nil {
      return shim.Error("Error building proposal definition - " + err.Error())
    }
    err = stub.PutState(proposalPrefix + proposal.Proposal.ProposalID, proposalAsBytes)
    if err != nil {
      return shim.Error("Error writing proposal to state - " + err.Error())
    }

    // Events because events
    proposalCreatedEvent := ProposalCreatedEventObject(ProposalID: proposal.Proposal.ProposalID)
    proposalEventAsBytes, err := json.Marshal(proposalCreatedEvent)
    if err != nil {
      return shim.Error("Error building proposal event definition - " + err.Error())
    }

    // Event for provided handler + timeout client
    err = stub.SetEvent(proposal.Proposal.Handler + " - Proposal Confirmed")
    err = stub.SetEvent("Proposal Created", proposalEventAsBytes)
    return shim.Success(nil)
  }

  func (s *HashTimeLockContract) confirmProposal(stub shim.ChaincodeStubInterface, args []string) peer.Response {
  	var err error
  	//Validate the args, expect 2, the proposalId and the pre-image of the hash
  	//for that proposalId
  	if len(args) != 2 {
  		return shim.Error("Invalid arguments to confirmProposal, expected proposalId, pre-image.")
  	}
  	//Retreive the proposal referenced
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

  	// Validation Logic Here

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
  	//Fire an event to inform middle actor to allow replaying into other channel
  	proposalConfirmedEvent := ProposalConfirmedEventObject{ProposalID: args[0], PreImage: args[1]}
  	proposalEventAsBytes, err := json.Marshal(proposalConfirmedEvent)
  	if err != nil {
  		return shim.Error("Error building proposal event definition - " + err.Error())
  	}
  	err = stub.SetEvent(ProposalConfirmedHandlerEvent, proposalEventAsBytes)
  	return shim.Success(nil)
  }

  func (s *HashTimeLockContract) invalidateProposal(stub shim.ChaincodeStubInterface, args []string) peer.Response {
  	var err error
  	//Validate the args, expect 1, the proposalId
  	if len(args) != 1 {
  		return shim.Error("Invalid arguments to invalidateProposal, expected proposalId")
  	}
  	/*
  	 * All sorts of validation logic about who can invalidate a proposal, maybe
  	 * check whether we have exceeded a minimum elapsed time or something?
  	 */
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
  	//Delete the proposal
  	err = stub.DelState(proposalPrefix + args[0])
  	if err != nil {
  		return shim.Error("Error while deleting proposal from state.")
  	}
  	return shim.Success(nil)
  }

  func main() {
  	err := shim.Start(new(HashTimeLockContract))
  	if err != nil {
  		fmt.Printf("Error creating new Smart Contract: %s", err)
  	}
  }
}
