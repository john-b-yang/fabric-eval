# Develop Go Clients to interact with Fabric Networks

Verifying the correctness and security of generated Fabric Chaincode will be an important evaluation component. To make tests repeatable and reliable, we will be developing behavioral tests with the Hyperledger Fabric Client [SDK](https://github.com/hyperledger/fabric-sdk-go) for Go. As mentioned in the README, "This SDK enables Go developers to build solutions that interact with Hyperledger Fabric". This SDK's general use case is to provide a programming abstraction for entities to use and interact with a Fabric network, ensuring there's no need to get mired in the low level details of how to set up the network itself.

Client code for specific contracts in this repository have been included in the corresponding directory in the `chaincode` folder. The `chaincode/fabcar` directory contains a copy of the SDK code from the `fabric-samples` repo provided by the Hyperledger Fabric organization.

## Deployment to Test Network

The existing client code can be compiled and run as is in the context of the `test-network` in the `fabric-samples` repo The steps for running the code in the test network are as follows:
1. Set up the test network as described in the `ec2-fabric.md` guide.
2. Deploy chaincode as described in the `ec2-deploy.md` guide.
3. Copy the contents of the companion client code for the chaincode to the `fabric-samples/fabcar/go/fabcar.go` file.
4. Change directory to the path `fabric-samples/fabcar/go/` (This is where the `fabcar.go` file referenced in step 3 is located).
5. Run the script `./runfabcar.sh`. This script will perform an export several trivial environment variales, then execute `go run fabcar.go`, which will build an executable from the client source code and execute the produced binary.

## Deployment to Local Network

Deploying to a local network requires a closer understanding of how the Fabric client SDK works and involves modifications to the existing client code for it to compile and run correctly in the test network. The following will discuss 1. The structure of Fabric client code and 2. How to deploy the code in the context of a local network.

### Client Structure
The information presented here is a condensed form of the official [documentation](https://hyperledger-fabric.readthedocs.io/en/release-2.2/developapps/application.html) along with some context-specific information.

The structure of the Fabric Client code can be broken down into three distinct phases. Each of these phases will be elaborated upon. In the first phase, references to peers in a running network are established by providing crypto and authentication files corresponding to each entity as arguments. Then, each peer reference will attempt to connect to the Fabric network, followed by the channel, then contract. Finally, once the contract references are retrieved, the app logic can be applied through a series of transaction invocation calls.

It's recommended to have the `chaincode/fabcar/fabcar-sdk.go` file open while going over this doc as a reference for practical examples.

##### Part 1: Peer References

The first part of the Fabric client involves determining the correct set of information for establishing a peer. Fabric applications use the concept of a "wallet", meant to represent the certificates and permissions associated with a peer. The terminology comes from the analogy of an individual having an actual leather wallet containing a driver's license, student ID, library card, or various forms of ID.

In Fabric, these IDs are the following values: MSP ID, X509 Digital Certificates, Signing Certificate, Connection Profile, and Key Store. Together, they define the privileges that a peer has with regards to issuing transactions on the network. Note that wallets don't retain any sort of financial data (i.e. cash, tokens) associated with a user.

In the `fabcar-sdk.go` code, the `addNewUser` and `populateWallet` functions capture the above behavior. `addNewUser` takes in several parameters.
* `credPath`: Provide relative path to the X509 Digital Certificates and MSP ID
* `ccpPath`: Provide relative path to the [Connection Profile](https://hyperledger-fabric.readthedocs.io/en/release-2.2/developapps/connectionprofile.html)
* From the local network's file structure and username, the key store and signing certificates can be determined.
At a high level, `addNewUser` will first create a wallet, then add the relevant identification and permissions to the wallet using the `populateWallet` function.

When adapting an existing client app to the local network, the following values must be changed in the sdk:
* `USER1`, `USER2`, `mspID`: Locate in the `docker-compose.yaml` file.
* `credPath`: This crypto will be generated when the network is up, and located in the `/organizations/peers/peer<ID>/crypto/` folder. Make sure the relative path points at the correct `.pem` file.
* `ccpPath`: This configuration file will be generated when the network is up, and located in the `/organizations/peers/peer<ID>/crypto/conn/` folder.
Thankfully, the signing certificates and keystore are also generated such that they should be located in the same directory that `credPath` poitns to.

TL;DR: In general, there shouldn't be any need to modify `populateWallet` or `addNewUser` implementations. All that's required is for the parameters to be configured correctly.

##### Part 2: Connect to the Network
Once the peer wallets are created, they can be used to connect to the network through the `gateway` abstraction. The Fabric `gateway` is a class for managing and contextualizing network interactions on behalf of an application. By using a peer's identity, the gateway can determine how a peer should be connected to the network. There is comprehensive [documentation](https://godoc.org/github.com/hyperledger/fabric-sdk-go/pkg/gateway), but we will simply explain the most important functions that are invoked in our implementations, primarily in the second half of the `addNewUser` function
* `NewFileSystemWallet`: Creates instance of wallet in memory (not persisted). The `path` parameter indicates where the Wallet object should be stored.
* `NewZ509Identity`: Creates new X509 Identity (Key, Certificate) that is stored in the Wallet object.
* `Connect`: Connects to a gateway. Parameters (config, identity, strategy) are options that determine the kind of connection the peer has to the network. For our purposes, we're just concerned with the peer being able to execute transactions. Returns a Gateway object.
* `gw.GetNetwork`: Returns Network object representing corresponding network channel, found by name.
* `network.GetContract`: Returns Contract object representing corresponding smart contract instance on Network object.
For each peer object, there should be a single wallet object.

For this part, none of the implementation needs to be changed. The above objects should be retrieved successfully assuming the parameters from Part 1 were entered correctly.

##### Part 3: Transactions
It's smooth sailing from here on out. Once the peer contract objects are retrieved, transactions from different parties will be communicated to and evaluated by the same Chaincode. The `main` function will typically contain the application logic.

Transactions are initiated with the `contract.SubmitTransaction` function. The function parameters require the first value to be the function name (i.e. `QueryAllCars`). Subsequent arguments should be the required parameters. Do note that:
* If the number of parameters passed don't match, an error will occur.
* If the parameter types are incorrect, an error will occur. The general rule of thumb is, when in doubt, pass the value as a string.
* The difference between `submitTransaction` and `evaluateTransaction` is described as follows. In general, only `submitTransaction` should be used to move the application forward.
  * `submitTransaction`: Used to invoke transactions that change (a.k.a. write to, `PutState`) the world state (a.k.a. transaction ledger).
  * `evaluateTransaction`: Used for transactions that simply query the world state (no modifications).
  * Under the hood, `submitTransaction` requests will be sent to the orderer, which will then handle the execution and commit of the transaction. `evaluateTransaction` calls are not sent to the ordering service, and won't be committed to the ledger. Peers that have endorsed the contract will still execute it.
* The `result` returned by `submit/evaluateTransaction` is a byte string of JSON data. If there are return values necessary to make forward progress in the application logic, JSON unmarshaling needs to be performed.

### Local Network Deployment
Local deployment is quite simple. Simply copy the updated client source code to the `/network/` directory in a file that must be named `client.go`. Then, run the `runClient.sh` script, which simply performs `go run client.go` 
