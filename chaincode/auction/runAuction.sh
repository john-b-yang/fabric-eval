#!/bin/bash

source runQueries.sh

chaincodeInvokeAuction() {
  setGlobalsForPeer0Org1

  ## Init ledger
  peer chaincode invoke -o localhost:7050 \
      --ordererTLSHostnameOverride orderer.example.com \
      --tls $CORE_PEER_TLS_ENABLED \
      --cafile $ORDERER_CA \
      -C $CHANNEL_NAME -n ${CC_NAME} \
      --peerAddresses localhost:7051 --tlsRootCertFiles $PEER0_ORG1_CA \
      --peerAddresses localhost:9051 --tlsRootCertFiles $PEER0_ORG2_CA \
      -c '{"function": "InitLedger","Args":["20"]}'
}

chaincodeQuery() {
  setGlobalsForPeer0Org1
  peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"function": "SubmitBid","Args":["100"]}'

  sleep 20

  setGlobalsForPeer0Org2
  peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"function": "CloseBid","Args":[]}'
}

chaincodeInvokeInit
sleep 2
chaincodeInvokeAuction
sleep 2
chaincodeQuery
