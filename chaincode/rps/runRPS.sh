#!/bin/bash

source runQueries.sh

GAMEID = "PUT GAME ID HERE AFTER INITLEDGER"
KEY1 = "0"
KEY2 = "2"
NONCE1 = "Jack"
NONCE2 = "John"
HASH1 = "B23DE3EA7971B4490AB8175B23BC37E016C9B190F1FD7CB84C2D07BA05EB3050"
HASH2 = "86F05364C02B1F989218954443BAD8FEDEE1C65B30B4E14E0241D14E20C85880"

chaincodeQuery() {
    # Join Game
    setGlobalsForPeer0Org1
    peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"function": "JoinGame","Args":["$GAMEID"]}'
    setGlobalsForPeer0Org2
    peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"function": "JoinGame","Args":["$GAMEID"]}'

    # Make Choice
    setGlobalsForPeer0Org1
    peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"function": "MakeChoice","Args":["$GAMEID", "$HASH1"]}'
    setGlobalsForPeer0Org2
    peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"function": "MakeChoice","Args":["$GAMEID", "$HASH2"]}'

    # Reveal Choice
    setGlobalsForPeer0Org1
    peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"function": "RevealChoice","Args":["$GAMEID", "$KEY1", "$NONCE1"]}'
    setGlobalsForPeer0Org2
    peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"function": "RevealChoice","Args":["$GAMEID", "$KEY2", "$NONCE2"]}'

    # Determine Winner
    peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"function": "DetermineWinner","Args":["$GAMEID"]}'
}

chaincodeInvokeInit
sleep 2
chaincodeInvoke
sleep 2
chaincodeQuery
