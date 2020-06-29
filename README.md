# Provenance Contract

This chaincode implements the Provenance model based on the W3C PROV standard, using Hyperledger fabric v1 as backend. 

## Introduction

_Provenance_ is the foundation of data quality, usually implemented by automatically capturing the trace of 
data manipulation over space and time. In _healthcare_, provenance becomes critical since it encompasses both 
clinical research and patient safety. In this proposal we aim at exploiting and innovating existing health IT 
deployments by enabling data provenance queries for all kind of clinical information from anywhere. 
The proposed technical solution exploits the novelty and the peer-to-peer fashion of the _blockchain_ technology and 
_smart-contracts_ to instrument international standards such as IHE and HL7 with a provenance system robust to fraudulences. 

## Technology

The smart contract is implemented in Golang and it is able to update the Hyperledger's world state with new provenance
information for a specific document (e.g., PDF, DICOM) or XML (e.g., Consolidated CDAs, C32). 

### How to test

Testing is made using the Go testing framework.

```
go test --tags nopkcs11
```

## How to deploy on Hyperledger Fabric in Azure

Microsoft Azure offers [Hyperledger Fabric](https://azuremarketplace.microsoft.com/en-us/marketplace/apps/microsoft-azure-blockchain.azure-blockchain-hyperledger-fabric?tab=Overview) on their cloud infrastructure. For [GrapevineWorld](http://www.grapevineworldtoken.io) pilot we have deployed this chaincode to cope with the Facade system of the [GOE](http://github.com/GrapevineWorld/GOE). Having blockchain deployed in the cloud dramatically ease the adoption of the Provenance Solution, since no effort is requested for the IT staff of the hospitals.


With the Azure subscription (1 orderer, 1 CA, 2 peers) we performed the following tasks. 
* SSH into the orderer0 and create the hooks for the channel. 
```
peer channel create -o orderer0:7050 -c masab10 -f masab10.tx 
```
and copy the file masab10.tx in the peers. 

* SSH into peer0. Since Azure does not start the command line interface, we need to start the docker image as: 
```
  docker run -d --name CLI -v $HOME/crypto-config:/crypto-config \
  -v $HOME/masab10.tx:/masab10.tx \
  hyperledger/fabric-peer:x86_64-1.0.1 
  ```
* Create the channel on the peer, as 
```
  docker exec -it CLI /bin/bash 
```
* Export environment variables

```
  export CORE_PEER_ADDRESS="peer0:7051"
  export CORE_PEER_LOCALMSPID="Org1MSP"
  export CORE_PEER_MSPCONFIGPATH=/crypto-config/peerOrganizations/org1.example.com/users/Admin\@org1.example.com/msp
```
* Create the channel
```
  peer channel create -o orderer0:7050 -c masab10 -f masab10.tx
```

* Let the peer join the channel
```
  peer channel join -b masab.block
```

* Let the other peer1 to join the channel
```
 export CORE_PEER_ADDRESS=peer1:7051
 export CORE_PEER_LOCALMSPID="Org1MSP"
 export CORE_PEER_MSPCONFIGPATH=/crypto-config/peerOrganizations/org1.example.com/users/Admin\@org1.example.com/msp
 peer channel create -o orderer0:7050 -c masab10 -f masab10.tx
```

* Check in peer1 and peer0 the channel

```
  peer channel list 
```

* Now run the fabric tools 
```
docker run -v $HOME/crypto-config:/crypto-config -d -t hyperledger/fabric-tools:x86_64-1.0.1 /bin/bash
```

* The go version (1.7.5) provided by Azure does not match the current fabric sdk for go, so we need to install a new version, by downloading from https://golang.org/dl/ to avoid [this](https://stackoverflow.com/questions/49905951/hyperledger-install-failure-bccsp-factory-pluginfactory-go122-cannot-find-pac) error. Unpack it in /opt/go

* Install beevik/etree and the contract

```
go get -u github.com/beevik/etree
go get -u github.com/mascanc/ProvenanceContract/src/main
go build --tags nopkcs11
```
* Install the chaincode
```
export CORE_PEER_ADDRESS=gvwby7qwq-peer0:7051
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_MSPCONFIGPATH=/crypto-config/peerOrganizations/org1.example.com/users/Admin\@org1.example.com/msp/
root@4e3d36f26034:/opt/gopath/src/github.com/mascanc/ProvenanceContract/src/main# go build --tags nopkcs11
root@4e3d36f26034:/opt/gopath/src/github.com/mascanc# peer chaincode install -n provenancecc -v 0.1 -p github.com/mascanc/ProvenanceContract/src/main
root@4e3d36f26034:/# peer chaincode instantiate -n provenancecc -v 0.1 -c '{"Args":["john","0"]}' -C masab10

```
* Test if it stores Proveance documents
```
peer chaincode invoke -n provenancecc -C masab10 -c '{"Args":["set", "S52fkpF2rCEArSuwqyDA9tVjawUdrkGzbNQLaa7xJfA=", "agentInfo.atype", "1.2.3.4", "agentInfo.id", "agentidentifier", "agentinfo.name","7.8.9", "agentindo.idp","urn:tiani-spirit:sts","locationInfo.id", "urn:oid:1.2.3","locationInfo.name","General Hospital","locationInfo.locality","Nashville, TN", "locationInfo.docid","1.2.3","action","ex:CREATE","date","2018-11-10T12:15:55.028Z","digest1","E0nioxbCYD5AlzGWXDDDl0Gt5AAKv3ppKt4XMhE1rfo","digest3","xLrbWN5QJBJUAsdevfrxGlN3o0p8VZMnFFnV9iMll5o"]}'
```
* And then query for it
```
root@4e3d36f26034:/# peer chaincode invoke -n provenancecc -c '{"Args":["get","S52fkpF2rCEArSuwqyDA9tVjawUdrkGzbNQLaa7xJfA="]}' -C masab10
```

## Links

* [Hyperledger](http://www.hyperledger.org) - The blockchain technology
* [Medium Article](https://medium.com/cybersoton/decentralised-provenance-for-healthcare-exchange-services-b900cd96136c) - Medium article from cybersoton with video
* [Article](https://doi.org/10.1016/j.ijmedinf.2020.104197) - Decentralised Provenance for Healthcare Data (Int J of Med Inf Volume 141, September 2020, 104197)

## Authors

* **Massimiliano Masi** - *Initial work* - [mascanc](https://github.com/mascanc)
