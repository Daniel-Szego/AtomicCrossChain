/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"strconv"
	"time"
	"encoding/json"	
	"crypto/sha256"	
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a car
type SmartContract struct {
	contractapi.Contract
}

// InitLedger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

//	dataByteBob := []byte("300")

//	ctx.GetStub().PutState("Bob", dataByteBob)

//	dataByte := []byte("200")

//	return ctx.GetStub().PutState("Alice", dataByte)

	return nil
}

// Get Balance of an account
func (s *SmartContract) GetBalance (ctx contractapi.TransactionContextInterface, id string) (int) {

	retValue, _ := ctx.GetStub().GetState(id)
	intVar , _ := strconv.Atoi(string(retValue))

//	retValueString := string(retValue)

	return intVar
}


// Mint token
func (s *SmartContract) MintToken (ctx contractapi.TransactionContextInterface, id string, amount string) (string) {

	balanceInt := s.GetBalance(ctx, id)

	intAmount , _ := strconv.Atoi(amount)

	balanceInt = balanceInt + intAmount

	balanceString := strconv.Itoa(balanceInt)

	dataByte := []byte(balanceString)

	ctx.GetStub().PutState(id, dataByte)

	return "success"
}

// Mint token
func (s *SmartContract) BurnToken (ctx contractapi.TransactionContextInterface, id string, amount string) (string) {

	balanceInt := s.GetBalance(ctx, id)

	intAmount , _ := strconv.Atoi(amount)

	if intAmount >= balanceInt {
		return "error: not enoguh balance"
	}

	balanceInt = balanceInt - intAmount

	balanceString := strconv.Itoa(balanceInt)

	dataByte := []byte(balanceString)

	ctx.GetStub().PutState(id, dataByte)

	return "success"
}


// Transfer from one accountto another one
func (s *SmartContract) Transfer (ctx contractapi.TransactionContextInterface, from_id string, to_id string, amount string) (string) {

	s.BurnToken(ctx,from_id, amount)

	s.MintToken(ctx,to_id, amount)

	return "success"
}

// structure for the timelock
type HashTimeLock struct {
	LockID string `json:"lockid"`
	FromID string `json:"fromid"`
	ToID string `json:"toid"`
	Amount  string `json:"amount"`
	HashLock string `json:"hashlock"`
	TimeLock  string `json:"timelock"`
}

// Transfer from one accountto another one
func (s *SmartContract) TransferConditional (ctx contractapi.TransactionContextInterface, lock_id string, from_id string, to_id string, amount string, hashlock string, timelock string) error {

	// decrease from the from amount
	s.BurnToken(ctx, from_id, amount)

	// create HashTimeLock
	hashTimeLock := HashTimeLock {
		LockID:   lock_id,
		FromID:  from_id,
		ToID: to_id,
		Amount:  amount,
		HashLock: hashlock,
		TimeLock: timelock,
	}

	hashTimeLockAsBytes, _ := json.Marshal(hashTimeLock)

	return ctx.GetStub().PutState(lock_id, hashTimeLockAsBytes)
}

// Getting the created Hash time lock
func (s *SmartContract) GetHashTimeLock (ctx contractapi.TransactionContextInterface, lock_id string) (*HashTimeLock, error) {

	hashTimeLockAsBytes , _ := ctx.GetStub().GetState(lock_id)

	hashTimeLock := new(HashTimeLock)
	_ = json.Unmarshal(hashTimeLockAsBytes, hashTimeLock)

	return hashTimeLock, nil
}


// Commiting the HTLC
func (s *SmartContract) Commit (ctx contractapi.TransactionContextInterface, lock_id string, preimage string) (string) {

	hashTimeLockAsBytes , _ := ctx.GetStub().GetState(lock_id)

	hashTimeLock := new(HashTimeLock)
	_ = json.Unmarshal(hashTimeLockAsBytes, hashTimeLock)

	hash := sha256.Sum256([]byte(preimage))

	hashString := fmt.Sprintf("%x", hash)

	fmt.Println("Hash String:", hashString)

	// condition 1 hash pre image

	if hashTimeLock.HashLock != hashString {

		fmt.Println("Invalid password:", hashString, hashTimeLock.HashLock)
		fmt.Println("Transaction reverted:")

		return "invalid password"

	}

	// condition 2 time
	timestamp , _ := ctx.GetStub().GetTxTimestamp()
	timestampInt := timestamp.Seconds
	
	timelock , _ := time.Parse(time.RFC3339, hashTimeLock.TimeLock)

	if  timelock.Unix() < timestampInt {

		fmt.Println("Timelock already activated")
		fmt.Println("Actual transaction timestamp:", timestampInt)
		fmt.Println("Actual timelock:", timelock.Unix())
		fmt.Println("Transaction reverted")

		return "Transaction reverted:Timelock already activated"
	}

	// increase amount to
 	s.MintToken(ctx, hashTimeLock.ToID, hashTimeLock.Amount)
	
	// delete lock
	ctx.GetStub().DelState(lock_id)

	// rasie event
	ctx.GetStub().SetEvent("Commit", []byte(preimage))

	fmt.Println("success commit")
	return "success commit"
}

// Revert the HTLC
func (s *SmartContract) Revert (ctx contractapi.TransactionContextInterface, lock_id string) (string) {

	hashTimeLockAsBytes , _ := ctx.GetStub().GetState(lock_id)

	hashTimeLock := new(HashTimeLock)
	_ = json.Unmarshal(hashTimeLockAsBytes, hashTimeLock)

	// condition 1 hash pre image - DOES NOT MATTER

	// condition 2 time

	timestamp , _ := ctx.GetStub().GetTxTimestamp()
	timestampInt := timestamp.Seconds
	
	timelock , _ := time.Parse(time.RFC3339, hashTimeLock.TimeLock)

	if  timelock.Unix() > timestampInt {

		fmt.Println("Timelock not yet activated")
		fmt.Println("Actual transaction timestamp:", timestampInt)
		fmt.Println("Actual timelock:", timelock.Unix())
		fmt.Println("Transaction reverted")

		return "Transaction reverted:Timelock not yet activated"
	}


	// increase amount from

 	s.MintToken(ctx, hashTimeLock.FromID, hashTimeLock.Amount)

	fmt.Println("success commit")
	return "success commit"
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create test chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting test chaincode: %s", err.Error())
	}

}