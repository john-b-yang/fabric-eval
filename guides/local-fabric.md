# Set Up Hyperledger Fabric Locally

This guide provides documentation for locally setting up a Fabric Network (in `network` folder of this repository). You can refer to `network/README` for more information regarding the actual contents of the folder.

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
