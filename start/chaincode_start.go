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
	"time"

	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("numbermgmt")
var trasactionString = "transaction_"

// NumberManagementChainCode example simple Chaincode implementation
type NumberManagementChainCode struct {
}

// NumberInfo information for number management
type NumberInfo struct {
	Number    string `json:"Number"`
	Available string `json:"Available"`
	Company   string `json:"Company"`
}

// TransactionHistory temp solution
type TransactionHistory struct {
	Transactions []TransactionEntry `json:"transactions"`
}

//TransactionEntry temp solution
type TransactionEntry struct {
	Date   string `json:"date"`
	Detail string `json:"detail"`
}

//updateNumber is update transaction changin company
func updateNumberCompany(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("Entering updateNumberCompany")
	if len(args) < 2 {
		logger.Error("Invalid number of arguments for UpdateCompany")
		return nil, errors.New("Invalid number of arguments for updateNumberCompany")
	}
	var number = args[0]
	var company = args[1]

	numberBytes, err := stub.GetState(number)
	if err != nil {
		logger.Error("Error retrieving number ", err)
		return nil, err
	}
	if numberBytes == nil {
		logger.Error("Number " + number + " not found in system")
		return nil, errors.New("Number " + number + " not found in the system")
	}
	var numberInfo NumberInfo
	err = json.Unmarshal(numberBytes, &numberInfo)
	if err != nil {
		logger.Error("Error marshaling data in store for number " + number)
		return nil, err
	}
	numberInfo.Company = company

	numberBytes, err = json.Marshal(&numberInfo)
	if err != nil {
		logger.Error("Error Marshling numberinfo", err)
		return nil, err
	}
	err = stub.PutState(number, numberBytes)
	if err != nil {
		logger.Error("Error saving numberinfo", err)
		return nil, err
	}
	transactionLog(stub, number, "Updated number")
	return nil, nil
}

// CreateNumber is create number func
func CreateNumber(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("Entering CreateNumber")
	if len(args) < 2 {
		logger.Error("Invalid number of arguments")
		return nil, errors.New("Missing arguments")
	}
	var number = args[0]
	var available = args[1]
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
	err = stub.PutState(trasactionString+number, numberBytes)

	if err != nil {
		logger.Error("Error Marshling numberinfo", err)
		return nil, err
	}
	transactionLog(stub, number, "creating number")
	logger.Info("Successfully updated Number Management")
	return nil, nil
}

func transactionLog(stub shim.ChaincodeStubInterface, key string, transaction string) ([]byte, error) {
	transactionKey := trasactionString + key
	transactionbytes, err := stub.GetState(transactionKey)
	if err != nil {
		logger.Error("Error logging transcation for "+key, err)
		return nil, err
	}
	var transactionHistory TransactionHistory
	var transactionEntry = TransactionEntry{
		Date:   time.Now().String(),
		Detail: transaction,
	}
	if transactionbytes != nil {
		err = json.Unmarshal(transactionbytes, &transactionHistory)
	}
	transactionHistory.Transactions = append(transactionHistory.Transactions, transactionEntry)
	transactionbytes, err = json.Marshal(&transactionHistory)
	stub.PutState(transactionKey, transactionbytes)
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

//GetTransaction history temp solution
func GetTransactionHistory(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("Entering retrieve Transaction Hitory")
	if len(args) != 1 {
		logger.Error("Invalid number of arguments")
		return nil, errors.New("Missing arguments")
	}
	var number = args[0]
	transactionKey := trasactionString + args[0]
	transactionbytes, err := stub.GetState(transactionKey)
	if err != nil {
		logger.Error("Could not fetch number with id "+number, err)
		return nil, err
	}
	return transactionbytes, nil
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
	if len(args) < 1 {
		logger.Error("Incorrect number of arguments. Expecting 2")
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *NumberManagementChainCode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Debug("invoke is running " + function)

	// Handle different functions
	if function == "init" { //initialize the chaincode state, used as reset
		return CreateNumber(stub, args)
	} else if function == "updatecompany" {
		return updateNumberCompany(stub, args)
	}
	logger.Error("invoke did not find func: " + function) //error

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *NumberManagementChainCode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "GetNumberInformation" {
		return GetNumberInformation(stub, args)
	} else if function == "GetTransactionHistory" {
		return GetTransactionHistory(stub, args)
	}
	logger.Error("invoke did not find func: " + function) //error
	return nil, errors.New("Received unknown function invocation: " + function)
}
