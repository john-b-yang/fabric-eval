# Deploy Application Code to Test Network

This is a walkthrough on how to deploy a custom application to the Fabric test network provided in the `fabric-samples` repository. ([Reference](https://hyperledger-fabric.readthedocs.io/en/release-2.0/deploy_chaincode.html))

### Instructions
**0. Before Beginning**
* This tutorial assumes `docker`, `go`, and `fabric` are installed
* Ensure that GOPATH, PATH for `golang` are set
* Ensure the test network is running with a channel (i.e. `./network.sh up createChannel`)

**1. Create Package for Application Code** ([Reference](https://hyperledger-fabric.readthedocs.io/en/release-2.0/deploy_chaincode.html#go))

1. Navigate to fabcar source code: `cd ~/fabric-samples/chaincode/fabcar/go/`
2. Replace existing fabcar code: `vi fabcar.go`
    * Delete existing code (`esc`, `ggdG`)
    * Paste in application code
3. Follow Remaining Steps
    * `go mod vendor`
    * Install relevant binaries + CLI (`cd` back into `~/fabric-samples/test-network/`
        * `export PATH=${PWD}/../bin:${PWD}:$PATH`
        * `export FABRIC_CFG_PATH=$PWD/../config/`
    * Set `CORE_PEER_MSPCONFIGPATH`
4. Create Chaincode Package: `peer lifecycle chaincode package <package name> --path ../chaincode/fabcar/go/ --lang golang --label <label>`

**2. Install Chaincode Package**: Follow given steps ([Reference](https://hyperledger-fabric.readthedocs.io/en/release-2.0/deploy_chaincode.html#install-the-chaincode-package))

**3. Approve Chaincode Definition**: Follow given steps ([reference](https://hyperledger-fabric.readthedocs.io/en/release-2.0/deploy_chaincode.html#approve-a-chaincode-definition))
* Note: `CC_PACKAGE_ID=<Output of queryinstalled>`

**4. Commit + Invoke Chaincode Definition to Channel**: Follow given steps ([reference](https://hyperledger-fabric.readthedocs.io/en/release-2.0/deploy_chaincode.html#committing-the-chaincode-definition-to-the-channel))

### Helpful Links
* Hyperledger Fabric Go Contract [link](https://github.com/hyperledger/fabric-contract-api-go)
* ERC 20 in Hyperledger Fabric [example](https://medium.com/coinmonks/erc20-token-as-hyperledger-fabric-golang-chaincode-d09dfd16a339)
* `cckit`: 3rd party [tool](https://github.com/s7techlab/cckit) for building Fabric contracts
* Hyperledger Fabric readthedocs [link](https://hyperledger-fabric.readthedocs.io/en/release-2.0/)
* HTLC Example [link](https://github.com/CallanHP/hlf-htla-proof-of-concept)
