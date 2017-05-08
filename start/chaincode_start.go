/*
Copyright IBM Corp 2016 All Rights Reserved.

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
	"errors"
	"fmt"
	"strconv"

	"encoding/json"

	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("numbermgmt")

// NumberManagementChainCode example simple Chaincode implementation
type NumberManagementChainCode struct {
}

// NumberInfo information for number management
type NumberInfo struct {
	Number    string `json:"Number"`
	Available bool   `json:"Available"`
	Company   string `json:"Company"`
}

// CreateNumber is create number func
func CreateNumber(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("Entering CreateNumber")
	if len(args) < 2 {
		logger.Error("Invalid number of arguments")
		return nil, errors.New("Missing arguments")
	}
	var number = args[0]
	var available, err = strconv.ParseBool(args[1])
	if err != nil {
		logger.Error("Error parsing Boolean from arguments, defaulting to true", err)
		available = true
	}
	var company = args[2]

	numberInfo := &NumberInfo{
		Number:    number,
		Available: available,
		Company:   company,
	}
	numberBytes, err := json.Marshal(&numberInfo)
	if err != nil {
		logger.Error("Error Marshling numberinfo", err)
		return nil, err
	}
	err = stub.PutState(number, numberBytes)
	if err != nil {
		logger.Error("Error Marshling numberinfo", err)
		return nil, err
	}
	logger.Info("Successfully updated Number Management")
	return nil, nil
}

//GetNumberInformation retrives number information
func GetNumberInformation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("Entering Quering")
	if len(args) != 1 {
		logger.Error("Invalid number of arguments")
		return nil, errors.New("Missing arguments")
	}
	var number = args[0]
	bytes, err := stub.GetState(number)

	if err != nil {
		logger.Error("Could not fetch number with id "+number, err)
		return nil, err
	}
	logger.Info("Processed Querying of Number information")
	return bytes, nil
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(NumberManagementChainCode))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *NumberManagementChainCode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Debug("Entering Init")
	if len(args) != 1 {
		logger.Error("Incorrect number of arguments. Expecting 1")
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *NumberManagementChainCode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Debug("invoke is running " + function)

	// Handle different functions
	if function == "init" { //initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	}
	fmt.Println("invoke did not find func: " + function) //error

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *NumberManagementChainCode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "GetNumberInformation" {
		return GetNumberInformation(stub, args)
	}
	return nil, nil
}
