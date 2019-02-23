/*
Copyright IBM Corp. 2017,2018 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
)

// TxID is just a dummuy transactional ID for test cases
const TxID = "mockTxID"

func TestContractChaincode_health(t *testing.T) {
	scc := new(ContractChaincode)
	stub := shim.NewMockStub("ContractChaincode", scc)
	// Contract
	// mockInvokeArgs := [][]byte{[]byte(function), args}
	// res := stub.MockInvoke(MockStubUUID, mockInvokeArgs)
	res := stub.MockInvoke(TxID, [][]byte{[]byte("Health")})
	assert.Equal(t, int(res.Status), shim.OK, "Health failed.")
	payload := string(res.Payload)
	assert.Equal(t, payload, "Ok", "Contract return value not expected.")
}

func TestContractChaincode_storeOrder(t *testing.T) {
	scc := new(ContractChaincode)
	stub := shim.NewMockStub("ContractChaincode", scc)

	const orderStr = `{
      "id": "d3ae8bb4-10ce-40a2-9a15-35bc6399df68",
	  "name": "My Super Awesome Order",
	  "approved": true,
	  "amount": 1500,
	  "createdTs": "2020-12-31T21:17:34.371Z",
	  "reviewedTs": "2020-12-31T21:17:34.371Z" 
	}`

	var order *Order
	orderAsBytes := []byte(orderStr)
	err := json.Unmarshal(orderAsBytes, &order)
	assert.Nil(t, err, "Failed to unmarshall order string.")

	// Contract - save order
	mockInvokeArgs := [][]byte{[]byte("StoreOrder"), orderAsBytes}
	res := stub.MockInvoke(TxID, mockInvokeArgs)
	assert.Equal(t, int(res.Status), shim.OK, "StoreOrder failed.")
	assert.Nil(t, res.Payload, "Contract return value not expected.")

	// Contract - read/get order
	mockInvokeArgs = [][]byte{[]byte("GetOrder"), []byte(order.ID)}
	res = stub.MockInvoke(TxID, mockInvokeArgs)
	readOrderAsBytes := res.Payload
	var retrievedOrder *Order
	err = json.Unmarshal(readOrderAsBytes, &retrievedOrder)
	assert.Nil(t, err, "Failed to unmarshall retrieved order.")
	assert.True(t, reflect.DeepEqual(order, retrievedOrder))
	//fmt.Printf("%+v\n", retrievedOrder)
}

func TestContractChaincode_storeStr(t *testing.T) {
	scc := new(ContractChaincode)
	stub := shim.NewMockStub("ContractChaincode", scc)

	const str = "the string to store..."
	strAsBytes := []byte(str)

	// Contract - save str
	mockInvokeArgs := [][]byte{[]byte("StoreStr"), strAsBytes}
	res := stub.MockInvoke(TxID, mockInvokeArgs)
	assert.Equal(t, int(res.Status), shim.OK, "StoreStr failed.")
	assert.Nil(t, res.Payload, "Contract return value not expected.")

	// Contract - read/get str
	mockInvokeArgs = [][]byte{[]byte("GetStr")}
	res = stub.MockInvoke(TxID, mockInvokeArgs)
	readStrAsBytes := res.Payload
	assert.Equal(t, str, string(readStrAsBytes))
	//fmt.Printf("%+v\n", retrievedOrder)
}
