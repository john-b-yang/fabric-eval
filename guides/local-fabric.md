# Set Up Hyperledger Fabric Locally

This guide provides documentation for locally setting up a Fabric Network (in `network` folder of this repository).

## Getting Started
1. Install relevant languages and tools: Docker, Docker-Compose, Golang (Version 1.13 or higher), Python (3 or higher)
2. Clone this repository.
3. `curl -sSL https://bit.ly/2ysbOFE | bash -s` to download the necessary Docker images. You may remove the additional `fabric-samples` repo if you'd like.
4. Add the `cli-tools` directory to your path so that the `peer`, `cryptogen`, and additional command line tools for operating the network can be referenced (i.e. `export PATH="~/fabric-eval/cli-tools:$PATH"`).
5. `cd` into the `network` folder.
6. `./start.sh` will
  - Spawn 2 Organizations, 2 Peers per Org, 1 CA per Org, and a State Database (Couch DB)
  - Create a Channel that peers from Org1 will be added to.
  - Compile and deploy chaincode to the channel.

## Commands Cheat Sheet
* Bring up the network: `docker-compose up -d` (without `-d`, network logs will be displayed).
* Bring down the network: `docker-compose down`
* Deploy Chaincode: Refer to directions in `local-deploy.md`

## Repository Layout
The following is a general layout of the network repository's configuration files (*not* shell scripts, which are described in the following section).

**artifacts/**: Contains scripts and assets for initializing the network
  * *channel/*: Artifacts (crypto, configurations, settings)
    * `config`: Values from [here](https://github.com/hyperledger/fabric/tree/master/sampleconfig)
      * `configtx`: Defines network components' properties including channel, transaction, profile, orderer, application, capabilities, etc.
      * `core`: Basic configuration option for various peer modules
      * `orderer`: Same as `core`, but for core modules.
    * `crypto-config`: Generated from running `./create-artifacts.sh`. Each folder contains each organization's certificate authority, MSP, orderer nodes, tls/ca assets, and associated users (a.k.a. peers)
      * `ordererOrganizations`
      * `peerOrganizations`
  * *src/github.com/*: Contains source code and API for Chaincode
    * `fabcar`: Directory contains raw Chaincode source code + dependencies, written in Go. The `deployChaincode.sh` script will compile the source code in this directory, generate a tarball, and deploy it to a user provided channel.
  * *private-data/*: Contains configurations for [collections](https://hyperledger-fabric.readthedocs.io/en/release-2.2/private-data-arch.html). (Not important to network infra)

## Configurations Explanations
**Network Settings**: (`artifacts/docker-compose.yaml`) Define the containers for the networks and services that will be created upon `docker-compose up -d`.
* Services: Certificate Authority Organizations, Orderers, CouchDB, Peers
* Configurations: Container name, Docker image, environment variable values, ports, networks, volumes (disk for persistence)

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

**Deploy Chaincode**: (`deployChaincode.sh`) Deploys chaincode in `./artifacts/src/github.com/fabcar/go` repository to channel. Sequential, it does the following step by step. Make sure network is up and running to ensure the below steps will work.
1. Define configurations and variables for network and each peer. The `set...` functions set the environment variables to be a particular peer, allowing for easy role assumption.
2. `presetup`: Download all dependencies for chaincode in `./artifacts/src/github.com/fabcar`.
3. `packageChaincode`: Invokes `peer lifecycle chaincode package` function to generate tarball of Chaincode source code.
4. `installChaincode`: Invokes `peer lifecycle chaincode install` function such that each peer installs `.tar.gz` file of chaincode.
5. `queryInstalled`: (Not necessary) Checks if chaincode was installed successfully on a peer.
6. `approveForMyOrg1/2`: A peer (w/ correct permissions) fires a transaction that indicates an installed Chaincode package should be approved for its organization using `approveformyorg` call.
7. `checkCommitReadyness`: Check if chaincode is ready to be committed. In other words, if it has enough approval (approval >= 51%).
8. `commitChaincodeDefinition`: Commits chaincode to channel, succeeds if approval is reached.
9. `queryCommitted`
10. `chaincodeInvokeInit`, then invoke and query.

**Create Crypto Artifacts**: (`artifacts/channel/create-artifacts.sh`) Creates relevant crypto artifacts that are building blocks for permissions in the network.
* Removes existing artifacts upon start up of network
* Generates crypto artifacts with `cryptogen` tool based on configurations in `crypto-config.yaml` file in same directory and places them in `crypto-config/` folder.
* Generate system genesis block and channel configuration block with `configtxgen` tool.
* Later on, these crypto assets are used and referenced by higher level scripts.

## Credits
* [BasicNetwork](https://github.com/adhavpavan/BasicNetwork-2.0): Repository our network is based on.
* [Video Walkthrough](https://www.youtube.com/playlist?list=PLSBNVhWU6KjW4qo1RlmR7cvvV8XIILub6) of Fabric network infrastructure.
