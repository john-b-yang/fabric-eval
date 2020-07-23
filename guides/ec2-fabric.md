# Set Up Hyperledger Fabric on AWS EC2 Instance

**General Notes**
* Avoid using `sudo` beyond installation purposes. The `./network.sh` script may not run properly if `sudo` was used inconsistently due to incorrectly created file permissions.

### Instructions
**1. Configuring EC2 Instance**
1. Navigate to the "EC2" service within AWS Console
2. Click "Launch Instance"
3. For AMI, select "Ubuntu Server 18.04 LTS"
4. For Instance Type, select "t2.micro (Free Tier Eligible)"
5. Click "Launch" (bottom right). No extra configs necessary.

**2. Login + Set Up EC2 Environment**
1. Make sure instance is in "running" state
2. Locally, ssh into EC2 instance (i.e. `ssh -i <key file (*.pem)> ubuntu@<public DNS>`)
3. Grant `sudo` privileges to a user type: `sudo usermod -aG sudo ${USER}`
4. Run the following 2 commands:
    * `sudo apt-get install curl`
    * `sudo apt-get update`

**3. Install + Upgrade Docker, Docker Compose**
1. `sudo apt-get install docker`
2. `curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add`
3. `sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"`
4. `sudo apt-get update`
5. `apt-cache policy docker-ce`
6. `sudo apt-get install -y docker-ce`
7. `sudo apt-get install docker-compose`
8. `sudo apt-get upgrade`
9. `sudo usermod -aG docker ${USER}` (Explanation: [Link](https://www.digitalocean.com/community/questions/how-to-fix-docker-got-permission-denied-while-trying-to-connect-to-the-docker-daemon-socket))
10. Log out, then back in for Step 9 to take effect.

**4. Install Go (1.14.4)**
1. `wget https://dl.google.com/go/go1.14.4.linux-amd64.tar.gz`
2. `tar -xzvf go1.14.4.linux-amd64.tar.gz`
3. `sudo mv go/ /usr/local`
4. `export GOPATH=$HOME/go`
5. `export PATH=$PATH:/usr/local/go/bin`

**5. Fabric Installation**: `curl -sSL https://bit.ly/2ysbOFE | bash -s`

**6. Verify Installs**
* `curl --version`
* `go version`
* `docker --version`
* `docker-compose --version`
