package main

// prov_test.go - Massimiliano Masi - 20 November 2017. Test the prov.go smart contract. 
// 
// TODO add a check on the returned JSON, validate it with http://www.w3.org/ns/prov.xsd. Now 
// these steps are made manually
import (
    "testing"
    "fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim" 
)


// Init checks the init of the chaincode
func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

// Init Test to execute to check the chaincode init
func TestInit(t *testing.T) {
	fmt.Println("Entering the test method for Init")
	provcc := new(SimpleAsset)
	stub := shim.NewMockStub("ANY_PARAM", provcc)
	checkInit(t, stub, [][]byte{[]byte("init")})
}


// TestSetWrongArgs try to invoke over a set with just a PDF (no segmentation), but with a wrong number 
// of arguments
func TestSetWrongArgs(t *testing.T) {
	fmt.Println("Entering the test method for SetWrongArgs")
	provcc := new(SimpleAsset)
	stub := shim.NewMockStub("ANY_PARAM", provcc)
	
	checkInit(t, stub, [][]byte{[]byte("init")})

	res := stub.MockInvoke("1", [][]byte{[]byte("set"), []byte("S52fkpF2rCEArSuwqyDA9tVjawUdrkGzbNQLaa7xJfA=")})

	if res.Status != shim.ERROR {
		fmt.Println("Invoke failed", string(res.Message))
		t.FailNow()
	}
	
}


// TestSetGoodArgs invoke over a set with good arguments, but no PDF
func TestSetGoodArgs(t *testing.T) {
	fmt.Println("Entering the test method for SetGoodArgs")
	provcc := new(SimpleAsset)
	stub := shim.NewMockStub("ANY_PARAM", provcc)

	// Testing the init. It always return true. No parameters in init. 
	
	checkInit(t, stub, [][]byte{[]byte("init")})

	res := stub.MockInvoke("1", [][]byte{[]byte("set"), []byte("S52fkpF2rCEArSuwqyDA9tVjawUdrkGzbNQLaa7xJfA="),
	[]byte("agentInfo.atype"),[]byte("1.2.3.4"),
	[]byte("agentInfo.id"),[]byte("agentidentifier"),
	[]byte("agentinfo.name"),[]byte("7.8.9"),
	[]byte("agentinfo.idp"),[]byte("urn:tiani-spirit:sts"),
	[]byte("locationInfo.id"),[]byte("urn:oid:1.2.3"),
	[]byte("locationInfo.name"),[]byte("General Hospital"),
	[]byte("locationInfo.locality"),[]byte("Nashville, TN"),
	[]byte("locationInfo.docid"),[]byte("1.2.3"),
	[]byte("action"),[]byte("ex:CREATE"),
	[]byte("date"),[]byte("2018-11-10T12:15:55.028Z")})

	if res.Status != shim.OK {
		fmt.Println("Invoke failed", string(res.Message))
		t.FailNow()
	}
	
}


// TestSetWrongArgsNoAgentInfo invoke over a set with wrong arguments, no agentInfo. 
// This could be improved, actually. The agentInfo parameter is not
// rendered in the chaincode, it's just args[n,n+2]. TODO
func TestSetWrongArgsNoAgentInfo(t *testing.T) {
	fmt.Println("Entering the test method for SetWrongArgsNoAgentInfo")
	provcc := new(SimpleAsset)
	stub := shim.NewMockStub("ANY_PARAM", provcc)

	// Testing the init. It always return true. No parameters in init. 
	
	checkInit(t, stub, [][]byte{[]byte("init")})

	res := stub.MockInvoke("1", [][]byte{[]byte("set"), []byte("S52fkpF2rCEArSuwqyDA9tVjawUdrkGzbNQLaa7xJfA="),
	
	[]byte("action"),[]byte("ex:CREATE"),
	[]byte("date"),[]byte("2018-11-10T12:15:55.028Z")})

	if res.Status != shim.ERROR {
		fmt.Println("Invoke failed", string(res.Message))
		t.FailNow()
	}
	
}


// TestSetGoodArgsFull invoke over a set with good arguments, with segmentation
// In general argument checking is fragile, and we could improve it a bit
func TestSetGoodArgsFull(t *testing.T) {
	fmt.Println("Entering the test method for SetGoodArgsFull")
	provcc := new(SimpleAsset)
	stub := shim.NewMockStub("ANY_PARAM", provcc)

	// Testing the init. It always return true. No parameters in init. 
	
	checkInit(t, stub, [][]byte{[]byte("init")})

	res := stub.MockInvoke("1", [][]byte{[]byte("set"), []byte("S52fkpF2rCEArSuwqyDA9tVjawUdrkGzbNQLaa7xJfA="),
	[]byte("agentInfo.atype"),[]byte("1.2.3.4"),
	[]byte("agentInfo.id"),[]byte("agentidentifier"),
	[]byte("agentinfo.name"),[]byte("7.8.9"),
	[]byte("agentinfo.idp"),[]byte("urn:tiani-spirit:sts"),
	[]byte("locationInfo.id"),[]byte("urn:oid:1.2.3"),
	[]byte("locationInfo.name"),[]byte("General Hospital"),
	[]byte("locationInfo.locality"),[]byte("Nashville, TN"),
	[]byte("locationInfo.docid"),[]byte("1.2.3"),
	[]byte("action"),[]byte("ex:CREATE"),
	[]byte("date"),[]byte("2018-11-10T12:15:55.028Z"),
	[]byte("digest1"),[]byte("E0nioxbCYD5AlzGWXDDDl0Gt5AAKv3ppKt4XMhE1rfo"),
	[]byte("digest3"),[]byte("xLrbWN5QJBJUAsdevfrxGlN3o0p8VZMnFFnV9iMll5o")})

	if res.Status != shim.OK {
		fmt.Println("Invoke failed", string(res.Message))
		t.FailNow()
	}
	
}


// TestSetGetGoodArgsFull invoke over a set with good arguments, with segmentation
// In general argument checking is fragile, and we could improve it a bit
// In this test we set and get. 
func TestSetGetGoodArgsFull(t *testing.T) {
	fmt.Println("Entering the test method for SetGetGoodArgsFull")
	provcc := new(SimpleAsset)
	stub := shim.NewMockStub("ANY_PARAM", provcc)

	// Testing the init. It always return true. No parameters in init. 
	
	checkInit(t, stub, [][]byte{[]byte("init")})

	res := stub.MockInvoke("1", [][]byte{[]byte("set"), []byte("S52fkpF2rCEArSuwqyDA9tVjawUdrkGzbNQLaa7xJfA="),
	[]byte("agentInfo.atype"),[]byte("1.2.3.4"),
	[]byte("agentInfo.id"),[]byte("agentidentifier"),
	[]byte("agentinfo.name"),[]byte("7.8.9"),
	[]byte("agentinfo.idp"),[]byte("urn:tiani-spirit:sts"),
	[]byte("locationInfo.id"),[]byte("urn:oid:1.2.3"),
	[]byte("locationInfo.name"),[]byte("General Hospital"),
	[]byte("locationInfo.locality"),[]byte("Nashville, TN"),
	[]byte("locationInfo.docid"),[]byte("1.2.3"),
	[]byte("action"),[]byte("ex:CREATE"),
	[]byte("date"),[]byte("2017-11-21T10:29:49.816Z"),
	[]byte("digest1"),[]byte("E0nioxbCYD5AlzGWXDDDl0Gt5AAKv3ppKt4XMhE1rfo"),
	[]byte("digest2"),[]byte("xLrbWN5QJBJUAsdevfrxGlN3o0p8VZMnFFnV9iMll5o"),
	[]byte("digest3"),[]byte("THIS_IS_DIGEST_3"),
	[]byte("digest4"),[]byte("THIS_IS_DIGEST_4")})

	if res.Status != shim.OK {
		fmt.Println("Invoke failed", string(res.Message))
		t.FailNow()
	}
	
	resGet := stub.MockInvoke("1", [][]byte{[]byte("get"), []byte("S52fkpF2rCEArSuwqyDA9tVjawUdrkGzbNQLaa7xJfA=")})
	if resGet.Status != shim.OK {
		fmt.Println("Invoke failed", string(resGet.Message))
		t.FailNow()
	}
}


