# Hyperledger Fabric Binaries
This folder contains a set of binaries that are useful for creating, maintaining, and opearting a Hyperledger Fabric network, and should be used as specified in the tutorial for setting up the Fabric network in the `network` folder.

This list is a brief description of each binary's purpose.
* `Configtxgen`: Creat network artifcates (i.e. `genesis.block`, `channel.tx`)
* `Configtxlator`: Utility for generating channel configuration
* `Cryptogen`: Utility for generating key material
* `Discovery`: Command line client for service discovery
* `Idemixgen`: Utility for generating key material to be used with identity mixer (MSP)
* `Orderer`: Orders transactions, maintains list of orgs allowed to create channels ("consortium"), enforce basic access control.
* `Peer`: Network participant, belongs to organization
* `Fabric-ca-client`: Client for creating, registering, and enrolling new users.
