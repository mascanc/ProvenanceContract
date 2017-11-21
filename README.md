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
## Links

* [Hyperledger](http://www.hyperledger.org) - The blockchain technology
* [WhitePaper](http://) - coming soon

## Authors

* **Massimiliano Masi** - *Initial work* - [mascanc](https://github.com/mascanc)
