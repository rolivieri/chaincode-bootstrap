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
	"fmt"
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
	res := stub.MockInvoke(TxID, [][]byte{[]byte("Health")})
	assert.Equal(t, int(res.Status), shim.OK, "Health failed.")
	payload := string(res.Payload)
	assert.Equal(t, payload, "Ok", "Contract return value not expected.")
}

func Test_inheritance(t *testing.T) {
	scc := new(ContractChaincode)
	stub := NewMyMockStub("ContractChaincode", scc)
	fmt.Println(reflect.TypeOf(stub))
	stub.GetCreator() // Invoke GetCreator() from test case
	stub.GetChannelID()
	res := stub.MockInvoke(TxID, [][]byte{[]byte("Health")})
	assert.Equal(t, int(res.Status), shim.OK, "Health failed.")
	payload := string(res.Payload)
	assert.Equal(t, payload, "Ok", "Contract return value not expected.")
	scc.MyHealth(stub)
}
