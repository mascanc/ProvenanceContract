// Package main prov.go - Massimiliano Masi - 1/11/2017
// 
// 
// This chaincode implements the Provenance model for the CCC provenance challenge. This chaincode
// considers two arguments: set and get. The "set" command stores a provenance hash into the blockchain,
// and the get returns the provenance document and the history. This chaincode creates the PROV data
// structure as well.
// 
// 
package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"log"
	"strconv"
	"time"
	"encoding/json"
)

// SimpleAsset is present in all the chaincodes. It is an empty drta structure that it is used as pointers to the functions.
type SimpleAsset struct {
}

// The ReturnerMessage is the structure returned in the get operation.
type ReturnedMessage struct {
	Provenance string
}

// agent is the data structure that belongs to the Agent (in PROV jargon).
type agent struct {
	atype string
	id    string
	name  string
    identityProvider string
}

type location struct {
	id string
	name string
	locality string
	docid string
}

// Init is chaincode initialization. It is called by the peer when installing the chaincode,
// e.g., by doing peer chaincode install. No initiatilization is done here, so we
// return a success string.
func (t *SimpleAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {

	log.Println("Starting the smart contract. No initialization is necessary")
	return shim.Success([]byte("INITIALIZATION_DONE"))
}


// Invoke is during peer invocation. The arguments here are as follows.
// We expect a json containing a "set" or "get" function, with the parameters
// separated by "name" and "value". We decided to use this notation for the
// sake of Human Readability. At the end, if the case of CDA objects, the
// segments / values are presented.
func (t *SimpleAsset) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	// Extract the function and args from the transaction proposal
	log.Printf("Invocation of the chaincode called")
	fn, args := stub.GetFunctionAndParameters()

	if args == nil {
		log.Printf("No arguments passed")
		return shim.Error("No arguments passed")
	}
	log.Printf("Args len is %v", len(args))
	creator, erro := stub.GetCreator()
	if erro != nil {
		log.Printf("Obtained error: %s", erro.Error())
		return shim.Error(erro.Error())
	}

	// check who's calling

	log.Printf("Transaction from %s", creator)

	/*
	 * In fn I know if it is a query (get) or a update (set)
	 */
	var result string
	var err error
	
	if fn == "set" {
		log.Println("Obtanied a set")
		result, err = set(stub, args)
		log.Printf("Obtained result %v", result)
	} else {
		log.Println("Obtanied a get")
		resultDocument, errno := get(stub, args)
		if errno != nil {
			log.Printf("Obtained error in get %s", errno.Error())
			return shim.Error(errno.Error())
		}

		m := ReturnedMessage{string(resultDocument)}

		b, errMarshal := json.Marshal(m)
		log.Printf("The value of the marshalled %s", b)
		if errMarshal != nil {
			log.Printf("Obtained error while marshalling %s", errMarshal)
			return shim.Error(errMarshal.Error())
		}
		return shim.Success([]byte(b))
	}
	if err != nil {
		log.Printf("Obtained an error! %s", err.Error())
		return shim.Error(err.Error())
	}
	// Return the result as success payload
	return shim.Success([]byte("PROCESSED_OK"))
}


// get shall return the provenance documents for all the linked data
func get(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// Here is the algorithm. The first parameter is the has of the document for which we're searching
	// provenance. Once we get that, we check if it is an ODD, or not. If it is an ODD, we keep
	// searching until finding the top
	hasfOfObj := args[0]
	log.Printf("Searching %v", hasfOfObj)
	log.Print("Getting the state")
	value, err := stub.GetState(hasfOfObj)

	if err != nil {
		log.Printf("Failed to get the document %s with error: %s", args[0], err)
		return []byte(""), fmt.Errorf("Failed to get document: %s with error: %s", args[0], err)
	}
	if value == nil {
		log.Printf("I don't find the hash... returning an error")
		return []byte(""), fmt.Errorf("Hash not found: %v", args[0])
	}
	log.Printf("Found it, now adding history (if any)")
	base64encodedValue := base64.StdEncoding.EncodeToString([]byte(value))

	log.Printf("Obtained state %v %v", value, base64encodedValue)

	hist, erra := getHistoricalState(hasfOfObj, stub)
	if erra != nil {
		log.Println("No history found, got an error (but going anyway). Arg is: %s error is %s", args[0], erra)
		//return []byte(""), fmt.Errorf("Failed to get history: %s error is %s", args[0], erra)
	}

	log.Printf("Before Returning")

	var buffer bytes.Buffer
	buffer.WriteString("{ \"Original\" : ")
	buffer.WriteString("\"")
	buffer.WriteString(base64encodedValue)
	buffer.WriteString("\"")
	
	buffer.WriteString(",")
	buffer.WriteString(" \"History\" : ")
	buffer.WriteString(hist)
	
	buffer.WriteString("}")

	returnedString := buffer.String()
	log.Printf("Returned values to be marshalled as string are: %s", returnedString)

	

	dstByte := make([]byte, base64.StdEncoding.EncodedLen(buffer.Len()))
	log.Printf("dstByte is: %v and %v", len(dstByte), buffer.Len())
	base64.StdEncoding.Encode(dstByte, buffer.Bytes())
	return dstByte, nil
}

// getHistoricalState returns the historical state of the key. This is important in case of collisions!
// If no history is available, a "NO_HISTORY_AVAILABLE" is returned
func getHistoricalState(hasfOfObj string, stub shim.ChaincodeStubInterface) (string, error) {
	resultsIterator, err := stub.GetHistoryForKey(hasfOfObj)
	if err != nil {
		return "\"NO_HISTORY_AVAILABLE\"", err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			log.Printf("Obtained err %s", err.Error())
			return "Unable to iterate over history", err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		

		base64EncodedValue := base64.StdEncoding.EncodeToString([]byte(response.Value))
		buffer.WriteString("\"")

		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(base64EncodedValue))
		}
		buffer.WriteString("\"")

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoricalState returning:\n%s\n", buffer.String())
	return buffer.String(), nil
}

// set shall create the provenance documents, and store them into the blockchain.
// This is an example of set
//  peer chaincode invoke -n mycc -c '{"Args":["set", "S52fkpF2rCEArSuwqyDA9tVjawUdrkGzbNQLaa7xJfA=", "agentInfo.atype", "1.2.3.4", "agentInfo.id", "agentidentifier", "agentinfo.name","7.8.9", "agentinfo.idp", "urn:tiani-spirit:sts", "location.id", "urn:oid:1.2.3", "location.name", "General Hospital", "location.locality", "Nashville, TN", "location.docid", "1.2.3.4", "action", "ex:CREATE", "date", "2006-01-02T15:04:05", "digest1", "E0nioxbCYD5AlzGWXDDDl0Gt5AAKv3ppKt4XMhE1rfo", "digest2", "xLrbWN5QJBJUAsdevfrxGlN3o0p8VZMnFFnV9iMll5o", "digest3", "+DzwgaD7vGYb8S0MF79m/U5pyS9qnRSdqlFb1tkQUnc="]}' -C myc
func set(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	
	// fail first: check the arguments are at least 12. 
	if len(args) < 20 {
		log.Printf("Invalid number of parameters. Obtained %v", len(args))
		return "", fmt.Errorf("Invalid number of parameters. Expected 12 (at least), received %v", len(args))
	}
	
	// first argument is the name of the hash of the CDA
	hashOfCda := args[0]
	log.Println("Obtained hash of CDA %v", hashOfCda)
	// second argument is the agentInfo.atype
	agentInfo := agent{args[2], args[4], args[6], args[8]}
	locationInfo := location{args[10], args[12], args[14], args[16]}
	action := args[18]
	date := args[20]

	provenanceXml, provenanceString,errWhenGenerating := makeProvenanceDocument(hashOfCda, agentInfo, locationInfo, action, date)
	if errWhenGenerating != nil {
		log.Printf("Error when generating the segment %s", errWhenGenerating)

		return "",errWhenGenerating
	}
	log.Printf("Obtained provenanceXml %T, adding to the blockchain", provenanceXml)
	err := stub.PutState(hashOfCda, []byte(provenanceString))

	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s", args[0])
	}

	// Now the arguments are the following: digest1, ..., digestn.
	// I now loop
	if len(args) > 21 {
		log.Println("I have some more arguments to process")
		for i := 21; i < len(args)-1; i++ { // -1 since it's gonna be
			hashOfSegment := args[i+1]
			log.Println("Found a hash " + hashOfSegment + " constructing the provenance of the segment")
			provenanceOfTheSegment, segmentAsString, errWhenGeneratingSegment := makeProvenanceDocumentSegmented(hashOfSegment, hashOfCda, agentInfo, locationInfo, action, date)

			if errWhenGeneratingSegment != nil {
				log.Printf("Error when generating the segment %s", errWhenGeneratingSegment.Error())
				return "", errWhenGeneratingSegment
			}
		
			log.Printf("Obtained provenanceXmlSegmented %T %T", provenanceOfTheSegment, segmentAsString)
			log.Printf("Adding it to the blockchain")
			err := stub.PutState(hashOfSegment, []byte(segmentAsString))

			if err != nil {
				return "", fmt.Errorf("Failed to set asset: %s", args[0])
			}
			i++ // need to increase the value of i, since we go two by two
		}
	}

	return "OK", nil

}

//main function starts up the chaincode in the container during instantiate
func main() {
	if err := shim.Start(new(SimpleAsset)); err != nil {
		fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	}
}
