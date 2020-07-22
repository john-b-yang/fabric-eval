# Hyperledger Fabric Evaluation
Collection of popular Ethereum standards. Repository for contracts written in Golang for the Hyperledger Fabric platform.

### Deploy Application Code to Test Network
([Reference](https://hyperledger-fabric.readthedocs.io/en/release-2.0/deploy_chaincode.html))

**1. Create Package of Application Code** <br>
Ensure that GOPATH, PATH for golang are set
1. `cd ~/fabric-samples/chaincode/fabcar/go/`
2. `vi fabcar.go`
    * Delete existing code (`esc`, `:`, `ggdG`)
    * Paste in application code
3. `GO111MODULE=on go mod vendor`
4. `cd ~/fabric-samples/test-network/`
5. Install relevant binaries + CLI (i.e. `peer`)
    * `export PATH=${PWD}/../bin:${PWD}:$PATH`
    * `export FABRIC_CFG_PATH=$PWD/../config/`
6. `export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp`
7. Create Chaincode Package: `peer lifecycle chaincode package fabcar.tar.gz --path ../chaincode/fabcar/go/ --lang golang --label fabcar_1`

**2. Install Chaincode Package** ([reference](https://hyperledger-fabric.readthedocs.io/en/release-2.0/deploy_chaincode.html#install-the-chaincode-package))
1. Install Chaincode Package ([reference](https://hyperledger-fabric.readthedocs.io/en/release-2.0/deploy_chaincode.html#install-the-chaincode-package)). Acting as organization 1:
    * Set environment variables to operate `peer` CLI as Org1 admin user
```
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
```
    * `peer lifecycle chaincode install fabcar.tar.gz`
2. As Org1, `peer lifecycle chaincode queryinstalled`
3. Set to Org1 and Org2, install chaincode (check reference)

**3. Approve Chaincode Definition** ([reference](https://hyperledger-fabric.readthedocs.io/en/release-2.0/deploy_chaincode.html#approve-a-chaincode-definition))
1. `peer lifecycle chaincode queryinstalled`
2. `CC_PACKAGE_ID=<Output of queryinstalled>`
3. Approve Chaincode for Org: `peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name fabcar --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem`
4. Set environment variables to `Org1 admin`, then approve definition

**4. Commit + Invoke Chaincode Definition to Channel** ([reference](https://hyperledger-fabric.readthedocs.io/en/release-2.0/deploy_chaincode.html#committing-the-chaincode-definition-to-the-channel))
1. `peer lifecycle chaincode checkcommitreadiness --channelID mychannel --name fabcar --version 1.0 --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --output json`
2. `peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name fabcar --version 1.0 --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt`
3. `peer lifecycle chaincode querycommitted --channelID mychannel --name fabcar --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem`
4. `peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n fabcar --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"initLedger","Args":[]}'`

### Helpful Links
* Hyperledger Fabric Go Contract [link](https://github.com/hyperledger/fabric-contract-api-go)
* ERC 20 in Hyperledger Fabric [example](https://medium.com/coinmonks/erc20-token-as-hyperledger-fabric-golang-chaincode-d09dfd16a339)
* `cckit`: 3rd party [tool](https://github.com/s7techlab/cckit) for building Fabric contracts
* Hyperledger Fabric readthedocs [link](https://hyperledger-fabric.readthedocs.io/en/release-2.0/)
* HTLC Example [link](https://github.com/CallanHP/hlf-htla-proof-of-concept)
