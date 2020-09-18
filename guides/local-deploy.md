# Deploy Application Code to Local Network

Before doing this tutorial, make sure the Fabric network is up and running after following the `local-fabric` tutorial. The `network` README describes how the scripts that are referenced here actually work underneath the hood.

To start the network, run `docker-compose -f ./artifacts/docker-compose.yaml up -d`
This should successfully create two organizations, each with two peers.

Then, run the `./createChannel.sh` script to creating a channel called `mychannel` that the chaincode will be deployed to.

Before deploying the Chaincode, to configure the source code, locate the `artifacts/src/github.com/fabcar/go` directory. If a new contract is being deployed, delete the `vendor` folder and `go.sum` file. Replace the contents of `fabcar.go` with the desired source code, and update the `go.mod` file such that it includes the necessary dependencies (must be Go Version 13+).

Finally, run the `./deployChaincode.sh` script. The script will generate the Chaincode executable and walk it through the Chaincode lifecycle to the invocation stage, where it can be queried by peers.

*Note*: The `./start.sh` script sequentially does the above.
