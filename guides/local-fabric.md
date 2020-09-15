# Set Up Hyperlegder Fabric on Local Computer

This guide provides documentation for navigating the [BasicNetwork](https://github.com/adhavpavan/BasicNetwork-2.0) repository for setting up a Fabric Network on one's local computer. The original set of tutorials corresponding to the document can be found [here](https://www.youtube.com/playlist?list=PLSBNVhWU6KjW4qo1RlmR7cvvV8XIILub6).

`createChannel.sh`: Creates peers, organizations, and adds them to new channel
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

`deployChaincode.sh`: Deploys chaincode in `./artifacts/src/github.com/fabcar/go` repository to channel
* 
