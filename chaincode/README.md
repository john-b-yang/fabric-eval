# Fabric Chaincode
Repository of chaincode smart contracts developed for deployment and use on the Hyperledger Fabric network.

## Directory Layout
* `auction`: Simple open auction modeled after Solidity [example](https://solidity.readthedocs.io/en/v0.5.11/solidity-by-example.html#blind-auction).
* `claims`: Ethereum Claims Registry allowing people, contracts, machines to issue claims about one another. [Reference](https://github.com/ethereum/EIPs/issues/780)
* `example`: Skeleton code for smart contract that compiles.
* `fabcar`: Example from Fabric starter code.
* `htla`: Hashed timelock Contract. [Reference](https://liquality.io/blog/hash-time-locked-contracts-htlcs-explained/)
* `rps`: Rock Paper Scissors Game.
* `token`: Simple token modeled after ERC 20 token standard. [Reference](https://eips.ethereum.org/EIPS/eip-20)

## Fabric Core Differences

These are some of the core functions that make development on Hyperledger Fabric fundamentally different than with Solidity. These design choices greatly influence how and why chaincode development works as is.
* Fabric is a permissioned network. Unlike Ethereum, a permissionless network, entities must be approved and granted permission to join and interact with peers. The corresponding threat model for Fabric could assume that participants are trustworthy.
* There is *no* underlying currency for Fabric, unlike Ether for Solidity. Assets and consensus protocols for each instance of a Fabric network can be defined by participants.
* Fabric features two types of nodes: Orderers (order + propagate Xacts) and Peers (verify + execute submitted Xacts).
* Ethereum state is derived from an immutable ledger containing the history fo all Xacts. Conceptually, Hyperledger Fabric mimics this, but uses a SQL database in practice for persistence.
* Ethereum = Smart Contracts written in Solidity, Hyperledger Fabric = *Chaincode* written in Go/Java/JS.

## Implementation Notes
Relative to Solidity, these are notes regarding where the development process for Fabric is different.

**Chaincode Transaction History**
* In Solidity, everything about the contract (i.e. Stack, Heap) is saved entirely to the ledger. In Hyperledger Fabric, saving a value to the ledger requires an explicit call to `PutState`. Values that are not passed into the state will be lost if/when the contract is removed from a channel.
  * Values must be converted to *bytes* before being put into state. There's no type associated with state information.
* `GetState` is a corresponding keyword for retrieving values from the ledger.
* The `contractapi.Contract` [object](https://godoc.org/github.com/hyperledger/fabric-contract-api-go/contractapi#Contract) serves as the interface between the contract logic and the network.

**Chaincode Specific Methods**
* The `Init` function is called during chaincode instantiation to initalize any data.
* The `Query` function is a getter function that allows one to read data off of the transaction history.
* Chaincode has a single entry point `Invoke()`. The parameters of this method help determine which function to call + what arguments to pass.

**Asset and Token Management**
* No built in key words for sharing funds (i.e. `transfer`, `send`)
* Hyperledger Fabric allows developers to add an asset as Chaincode that can be treated as a token.

**Miscellaneous**
* Fabric has a large repository of open source Go SDKs + APIs, [one](https://godoc.org/github.com/hyperledger/fabric-contract-api-go/contractapi) of which is for Chaincode development.
* Self reference keywords like `this` or `self` do not exist. The [TransactionContextInterface](https://godoc.org/github.com/hyperledger/fabric-contract-api-go/contractapi#TransactionContextInterface) contains information per transaction invocation, including the client identity.
* Go's plethora of data structures allow for more expressive and dynamic app logic than Solidity.

## Chaincode Lifecycle
Unlike Ethereum, which allows any developer to ship any smart contract to the network for immediate use, there is a Chaincode "lifecycle" enforced by Fabric, a design motivated by introducing more security and privilege when invoking transactions in a permissioned network.

The lifecycle is: 1. Package => 2. Install => 3. Query => 4. Approve (by Organizations) => 5. Check Commit REadyness => 6. Commit => 7. Query Committed Chaincode (by Peer) => 8. Invoke `Init` function of Chaincode => 9. `Invoke`, `Query` Chaincode functions.
