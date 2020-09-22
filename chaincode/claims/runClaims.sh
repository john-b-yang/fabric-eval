#!/bin/bash

source runQueries.sh

chaincodeQuery() {
    setGlobalsForPeer0Org2

    # Set Claims
    peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"function": "SetClaim","Args":["${PEER0_ORG1_CA}", "CLAIM1", "VALUE1"]}'
    peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"function": "SetClaim","Args":["${PEER0_ORG2_CA}", "CLAIM2", "VALUE2"]}'

    # Set Self Claims
    peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"function": "SetSelfClaim","Args":["CLAIM3", "VALUE3"]}'

    setGlobalsForPeer0Org1
    peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"function": "SetSelfClaim","Args":["CLAIM4", "VALUE4"]}'

    # Remove Claim
    peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"function": "RemoveClaim","Args":["${PEER0_ORG2_CA}", "${PEER0_ORG1_CA}", "CLAIM1"]}'
}

chaincodeInvokeInit
sleep 2
chaincodeInvoke
sleep 2
chaincodeQuery
