package main

const proposalPrefix string = "_proposal_"

const PendingStatus = "PENDING"

const ConfirmStatus = "CONFIRMED"

type abstractProposal struct {
	ProposalID string `json:"proposalId"`
	Handler    string `json:"proposalHandler"`
}

// TODO: Add timelock here
type proposalEntry struct {
	Proposal      abstractProposal `json:"proposal"`
	Status        string           `json:"status"`
	Hash          string           `json:"hash"`
	HashAlgorithm string           `json:"hashAlgorithm"`
}

var validHashingAlgorithms = []string{"SHA256", "SHA384", "SHA512"}

type ProposalCreatedEventObject struct {
	ProposalID string `json:"proposalId"`
}

type ProposalConfirmedEventObject struct {
	ProposalID string `json:"proposalId"`
	PreImage   string `json:"preImage"`
}
