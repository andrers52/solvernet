package main

import (
	"fmt"
	"log"
	"sync"
)

// Ledger encapsulates the financial state of the blockchain participants
// by mapping addresses to balances
type Ledger struct {
	mutex            sync.Mutex
	AddressToBalance map[string]float64 `json:"address_to_balance"`
}

// NewLedger creates a new Ledger with initialized map
func NewLedger() *Ledger {
	return &Ledger{
		AddressToBalance: make(map[string]float64),
		mutex:            sync.Mutex{},
	}
}

// Update updates the state with a new block
func (ledger *Ledger) Update(block Block) error {

	if block.Data.Type != MonetaryTransaction {
		return fmt.Errorf("cannot update Ledger. block type is not monetary transaction")
	}

	// Update balances for transactions
	return ledger.addMonetaryTransaction(block)

}

// addMonetaryTransaction processes monetary transactions from a block
func (ledger *Ledger) addMonetaryTransaction(block Block) error {
	ledger.mutex.Lock()
	defer ledger.mutex.Unlock()
	tx := block.Data.Transaction
	if tx == nil {
		return fmt.Errorf("transaction data not found")
	}

	// Initialize balances if not already present
	if _, exists := ledger.AddressToBalance[tx.From]; !exists {
		ledger.AddressToBalance[tx.From] = ADDRESS_INITIAL_BALANCE
	}
	if _, exists := ledger.AddressToBalance[tx.To]; !exists {
		ledger.AddressToBalance[tx.To] = ADDRESS_INITIAL_BALANCE
	}
	// Update balances
	ledger.AddressToBalance[tx.From] -= tx.Amount
	ledger.AddressToBalance[tx.To] += tx.Amount
	return nil
}

// This will be used when reading from mass data storage or network
// (for nodes joining later)
func CreateLedgerFromBlockchain(bc *Blockchain) (*Ledger, error) {
	var newLedger *Ledger = NewLedger()

	// log the action
	log.Println("Creating ledger from blockchain")

	for _, block := range bc.Blocks {
		if block.Data.Type == MonetaryTransaction {
			// Update balances for transactions
			newLedger.addMonetaryTransaction(block)
		}
	}
	return newLedger, nil
}
