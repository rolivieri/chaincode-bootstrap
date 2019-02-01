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
	"fmt"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("ContractChaincodeLog")

// Order structure
type Order struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	CreatedTs  time.Time  `json:"createdTs,string"`
	Approved   *bool      `json:"approved"`
	ReviewedTs *time.Time `json:"reviewedTs,string"`
	Amount     uint64     `json:"amount,int"`
}

// ContractChaincode implementation
type ContractChaincode struct {
}

// Init nothing to initialize
func (t *ContractChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### Contract Init ###########")
	//nothing to initialize just return
	return shim.Success(nil)
}

// Invoke Support for calling chaincode to ensure operation is up
func (t *ContractChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### Contract Invoke ###########")

	function, args := stub.GetFunctionAndParameters()
	switch function {
	case "Health":
		// contract chaincode
		return t.Health(stub, args)
	case "StoreOrder":
		return t.StoreOrder(stub, args)
	case "GetOrder":
		return t.GetOrder(stub, args)
	}

	errorMsg := fmt.Sprintf("Unknown action, please check the first argument, expecting 'Health', 'StoreOrder', or 'GetOrder'. Instead, got: %s", function)
	logger.Errorf(errorMsg)
	return shim.Error(errorMsg)
}

// Health returns Ok if successful
func (t *ContractChaincode) Health(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("########### Contract Health ###########")
	logger.Infof("Chaincode is healthy.")
	return shim.Success([]byte("Ok"))
}

// StoreOrder stores an order.
func (t *ContractChaincode) StoreOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("########### Contract StoreOrder ###########")

	var orderStr string
	if len(args) > 0 {
		orderStr = args[0]
	}
	logger.Infof("Order submitted from client app: '%s'", orderStr)

	// Validate JSON object follows schema
	orderAsBytes := []byte(orderStr)
	order, err := t.UnmarshallOrder(orderAsBytes)

	if err != nil {
		errorMsg := fmt.Sprintf("Failed to unmarshal order: %s", err.Error())
		logger.Error(errorMsg)
		return shim.Error(errorMsg)
	}

	compositeKeyElements := []string{order.ID}
	compositeRecordKey, compositeErr := stub.CreateCompositeKey("orders", compositeKeyElements)
	if compositeErr != nil {
		return shim.Error("Failed to generate composite key " + strings.Join(compositeKeyElements, ",") + ".  Error: " + compositeErr.Error())
	}

	logger.Infof("Composite key: %s", compositeRecordKey)
	err = stub.PutState(compositeRecordKey, orderAsBytes)
	if err != nil {
		return shim.Error("Failed to store order data for id " + compositeRecordKey + ". Error: " + err.Error())
	}

	logger.Infof("Order successfully stored: '%s'", orderStr)
	return shim.Success(nil)
}

// GetOrder gets an order.
func (t *ContractChaincode) GetOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("########### Contract GetOrder ###########")

	var orderID string
	if len(args) > 0 {
		orderID = args[0]
	}
	logger.Infof("Order ID submitted from client app: '%s'", orderID)

	compositeKeyElements := []string{orderID}
	compositeRecordKey, compositeErr := stub.CreateCompositeKey("orders", compositeKeyElements)
	if compositeErr != nil {
		return shim.Error("Failed to generate composite key " + strings.Join(compositeKeyElements, ",") + ".  Error: " + compositeErr.Error())
	}

	logger.Infof("Composite key: %s", compositeRecordKey)
	orderAsBytes, err := stub.GetState(compositeRecordKey)
	if err != nil {
		return shim.Error("Failed to read order data for id " + compositeRecordKey + ". Error: " + err.Error())
	}

	return shim.Success(orderAsBytes)
}

// UnmarshallOrder unmarshalls (i.e. deserializes) the JSON input string into an Order struct
func (t *ContractChaincode) UnmarshallOrder(orderAsBytes []byte) (*Order, error) {
	var order *Order
	err := json.Unmarshal(orderAsBytes, &order)
	if err != nil {
		logger.Error(err)
	}
	return order, err
}

func main() {
	err := shim.Start(new(ContractChaincode))
	if err != nil {
		logger.Errorf("Error starting ContractChaincode: %s", err)
	}
}
