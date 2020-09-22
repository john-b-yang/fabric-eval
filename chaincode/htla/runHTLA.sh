#!/bin/bash

source runQueries.sh

KEY = "hello world"
HASH = "B94D27B9934D3E08A52E52D7DA7DABFAC484EFE37A5380EE9088F7ACE2EFCDE9"

chaincodeQuery() {
    setGlobalsForPeer0Org2

    # Generate Hash, verify it is identical to set claims argument
    peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"function": "GenerateHash","Args":["$KEY", "SHA256"]}'

    # Create Proposal
    peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"function": "CreateProposal","Args":["1", "100", "$HASH", "SHA256"]}'

    # Confirm Proposal
    peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"function": "ConfirmProposal","Args":["0", "$KEY"]}'

    # Invalidate Proposal (Should Fail)
    peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"function": "InvalidateProposal","Args":["0"]}'
}

chaincodeInvokeInit
sleep 2
chaincodeInvoke
sleep 2
chaincodeQuery
