# Fabric Network Code
The `test-network` folder in the `fabric-samples` code provided by Hyperledger Fabric is a great, ready-made network with many pre-set configurations.

In an effort to understand how the architecture supporting Fabric works alongside our interest in having more control over the network configurations, this repo is a set of scripts and assets for getting a local Fabric network up and running as quickly as possible. Many of the command line calls have been packaged in easy to run scripts that automate the set up and execution processes.

In the `guides` folder, you can refer to `local-fabric.md` for how to set up the network on your local machine, and `local-deploy` will tell you how to deploy chaincode to the network.

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
