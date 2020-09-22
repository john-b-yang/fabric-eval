#!/bin/bash

source runQueries.sh

chaincodeQuery() {
    setGlobalsForPeer0Org2

    # Query all cars
    peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"Args":["queryAllCars"]}'

    # Query Car by Id
    peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"function": "queryCar","Args":["CAR0"]}'
}

chaincodeInvokeInit
sleep 2
chaincodeInvoke
sleep 2
chaincodeQuery
