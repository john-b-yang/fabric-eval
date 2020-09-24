# Hyperledger Fabric Experiments

Experiments and notes by [Jack Kolb](https://people.eecs.berkeley.edu/~jkolb/) and John Yang about the Hyperledger Fabric open source blockchain platform.

This repository is organized as follows:

**chaincode** contains:
* Implementations of popular Ethereum Contract Standards (EIPs, ERCs) rewritten as Chaincode (Go) that compiles and is deployable to a Fabric network.
* Go code using Fabric's Client SDK to interact with a Hyperledger Fabric blockchain (specifically, contracts deployed on the blockchain).
* Notes discussing differences between Fabric Chaincode, Solidity, and Quartz (WIP).

**cli-tools** contains:
* Set of binaries for operating Fabric network.

**guides** contains:
* Tutorials for setting up a test network and deploying code on it.
* Tutorials for setting up a local network and deploying code on it.
* Tutorials for writing SDK application code to interface with deployed Chaincode.

**network** contains:
* Source code for local Fabric network. Tutorial for setting it up is in *guides*.
