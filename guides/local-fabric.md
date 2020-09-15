# Set Up Hyperledger Fabric on Local Computer

This guide provides documentation for locally setting up a Fabric Network (in `network` folder of this repository).

## Getting Started
1. Install relevant languages and tools
  * Docker, Docker-Compose
  * Golang (Version 1.13 or higher)
  * Python (3 or higher)
2. Clone this repository.
3. `curl -sSL https://bit.ly/2ysbOFE | bash -s` to download the necessary Docker images. You may remove the additional `fabric-samples` repo if you'd like.
4. Add the `cli-tools` directory to your path so that the `peer`, `cryptogen`, and additional command line tools for operating the network can be referenced (i.e. `export PATH="~/fabric-eval/cli-tools:$PATH"`).
5. `cd` into the `network` folder.
6. `./start.sh` will spawn the orderer, peer, and certificate authority nodes, then

## Scripts Explanations
**Create Channel**: (`createChannel.sh`) Creates peers, organizations, and adds them to new channel
* Enables TLS client authentication on a peer node (`CORE_PEER_TLS_ENABLED`)
* Each organization must have its own Certificate Authority for issuing enrollment. (`ORDERER_CA`, `PEER0_ORG1_CA`, `PEER0_ORG2_CA`).
* Points to folder containing fabric configurations (i.e. crypto location, docker-compose files)
* Sets fields and crypto for multiple peers (`LOCALMSPID`, `TLS_ROOTCERT_FILE`, `MSPCONFIGPATH`, `ADDRESS`)
* Creates Channel with `peer channel create` command. Options to specify include:
  * `-o` Orderer endpoint
  * `-c` Channel name
  * `-f` Configuration Transaction file
  * `--tls` Use tls when communicating with orderer endpoint
  * `--cafile` Path to trusted certificate of orderer endpoint
* Send configtx update file to the channel with `peer channel update`

**Deploy Chaincode**: (`deployChaincode.sh`) Deploys chaincode in `./artifacts/src/github.com/fabcar/go` repository to channel
* Package and write chaincode to a file with `peer lifecycle chaincode package`. Must specify:
  * `--path` Path to write to
  * `--lang` Language chaincode is written in
  * `--label` Package name
*

## Binary Rundown
* `Configtxgen`: Creat network artifcates (i.e. `genesis.block`, `channel.tx`)
* `Configtxlator`: Utility for generating channel configuration
* `Cryptogen`: Utility for generating key material
* `Discovery`: Command line client for service discovery
* `Idemixgen`: Utility for generating key material to be used with identity mixer (MSP)
* `Orderer`: Orders transactions, maintains list of orgs allowed to create channels ("consortium"), enforce basic access control.
* `Peer`: Network participant, belongs to organization
* `Fabric-ca-client`: Client for creating, registering, and enrolling new users.

## Credits
* [BasicNetwork](https://github.com/adhavpavan/BasicNetwork-2.0): Repository our network is based on.
* [Video Walkthrough](https://www.youtube.com/playlist?list=PLSBNVhWU6KjW4qo1RlmR7cvvV8XIILub6) of Fabric network infrastructure.
