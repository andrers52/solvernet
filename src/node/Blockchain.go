package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"sync"

	"github.com/davecgh/go-spew/spew"
)

// *** Types ***

type Block struct {
	Height   int       `json:"height"`
	Data     BlockData `json:"data"`
	Hash     string    `json:"hash"`
	PrevHash string    `json:"prevhash"`
}

type BlockDataType int

const (
	MonetaryTransaction BlockDataType = iota
	KnapsackProblemSubmission
	KnapsackProposedSolutionSubmission
)

type BlockData struct {
	Type        BlockDataType             `json:"type"`
	Transaction *Transaction              `json:"transaction,omitempty"`
	Problem     *KnapsackProblem          `json:"problem,omitempty"`
	Solution    *KnapsackProposedSolution `json:"proposed_solution,omitempty"`
}

type Transaction struct {
	From               string  `json:"from"`
	To                 string  `json:"to"`
	Amount             float64 `json:"amount"`
	ProblemBlockHeight int     `json:"problem_block_height"` // Identifies the block where the problem was submitted in
}

type Item struct {
	Weight int `json:"weight"`
	Value  int `json:"value"`
}

type Blockchain struct {
	Blocks []Block
	mutex  sync.Mutex
}

// *** Functions ***

func CreateNewBlockchain(ledger *Ledger) *Blockchain {

	blockchain := &Blockchain{Blocks: make([]Block, 0)}

	// Create a genesis transaction
	genesisProblem := KnapsackProblem{
		Items:    []Item{}, // This is a dummy value, it will be overwritten
		Capacity: 0,        // This is a dummy value, it will be overwritten
		Bounty:   1,
		Address:  "0x0",
	}

	genesisProblem.Items = make([]Item, 20)
	for i := range genesisProblem.Items {
		item := Item{
			Value:  i + 1,
			Weight: i + 1,
		}
		genesisProblem.Items[i] = item
	}

	sumOfWeights := GetProblemItemsSumWeight(genesisProblem)
	genesisProblem.Capacity = int(sumOfWeights * 2 / 3)

	genesisBlock, err := blockchain.GenerateProblemBlock(genesisProblem)
	if err != nil {
		spew.Dump(err)
		panic("Failed to generate genesis block")
	}

	spew.Dump(genesisBlock)

	blockchain.AddBlock(genesisBlock, ledger)

	return blockchain
}

func (bc *Blockchain) AddBlock(newBlock Block, ledger *Ledger) error {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	//switch on the type of block
	switch newBlock.Data.Type {
	case MonetaryTransaction:
		// check is valid transaction
		if err := bc.validateTransaction(*newBlock.Data.Transaction); err != nil {
			return err
		}
		// update state
		err := ledger.Update(newBlock)
		if err != nil {
			log.Println("Failed to update blockchain state:", err)
			return errors.New("invalid ledger update")
		}
	case KnapsackProblemSubmission:
		// check if the problem is valid
		if err := ValidateProblem(*newBlock.Data.Problem, bc); err != nil {
			return err
		}
	case KnapsackProposedSolutionSubmission:
		// check if the solution is valid
		if err := ValidateProposedSolution(*newBlock.Data.Solution, bc); err != nil {
			return err
		}
	default:
		return errors.New("invalid block type")
	}

	bc.Blocks = append(bc.Blocks, newBlock)

	// check for expired problem and add a rewarding transaction if there is a solution
	problemSolutionPair := bc.CheckForExpiredProblem()

	if problemSolutionPair == nil {
		return nil
	}

	// Add a rewarding transaction
	tx := Transaction{
		From:               problemSolutionPair.Problem.Address,
		To:                 problemSolutionPair.Solution.Address,
		Amount:             problemSolutionPair.Problem.Bounty,
		ProblemBlockHeight: problemSolutionPair.Solution.ProblemBlockHeight,
	}
	newTransactionBlock, err := bc.GenerateTransactionBlock(tx)
	if err != nil {
		log.Println("Failed to generate rewarding transaction block")
		return errors.New("failed to generate rewarding transaction block")
	}

	bc.AddBlock(newTransactionBlock, ledger)

	return nil
}

func (bc *Blockchain) GetBlock(blockHeight int) Block {
	// try to get the block from the blockchain
	// if it fails, spew the blockchain and blockchain state and panic
	if blockHeight < 0 || blockHeight >= len(bc.Blocks) {
		log.Println("Block height out of range")
		log.Println("Blockchain:")
		spew.Dump(bc.Blocks)
		panic("Block height out of range")
	}
	return bc.Blocks[blockHeight]

}

func (bc *Blockchain) isNewBlockCorrectlyChained(newBlock Block) bool {
	lastBlock := bc.getLastBlock()
	if lastBlock.Height+1 != newBlock.Height {
		return false
	}

	if lastBlock.Hash != newBlock.PrevHash {
		return false
	}

	calculatedHash, err := calculateHash(newBlock)
	if err != nil {
		return false
	}
	if calculatedHash != newBlock.Hash {
		return false
	}

	return true
}

func calculateHash(block Block) (string, error) {
	record := strconv.Itoa(block.Height) + block.PrevHash
	blockBytes, err := json.Marshal(block.Data)
	blockBytes = append(blockBytes, []byte(record)...)
	if err != nil {
		log.Printf("Failed to marshal block: %v", err)
		return "", err
	}
	h := sha256.New()
	h.Write([]byte(blockBytes))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed), nil
}

func (bc *Blockchain) generateNewBlock(data BlockData) (Block, error) {
	oldBlock := bc.getLastBlock()
	var newBlock Block

	// print type of block that is being created
	log.Printf("Generating new block of type: %v", data.Type)

	newBlock.Height = oldBlock.Height + 1

	newBlock.Data = data
	newBlock.PrevHash = oldBlock.Hash

	calculatedHash, err := calculateHash(newBlock)
	if err != nil {
		return Block{}, err
	}
	newBlock.Hash = calculatedHash

	return newBlock, nil
}

func (bc *Blockchain) GenerateProblemBlock(problem KnapsackProblem) (Block, error) {
	data := BlockData{
		Type:    KnapsackProblemSubmission,
		Problem: &problem,
	}

	return bc.generateNewBlock(data)
}

func (bc *Blockchain) GenerateProposedSolutionBlock(proposedSolution KnapsackProposedSolution) (Block, error) {
	data := BlockData{
		Type:     KnapsackProposedSolutionSubmission,
		Solution: &proposedSolution,
	}
	return bc.generateNewBlock(data)
}

func (bc *Blockchain) getLastBlock() Block {
	if len(bc.Blocks) == 0 { // if the blockchain is empty, return a block with -1 height
		return Block{
			Height: -1,
		}
	}
	return bc.Blocks[len(bc.Blocks)-1]
}

func (bc *Blockchain) GenerateTransactionBlock(tx Transaction) (Block, error) {
	data := BlockData{
		Type:        MonetaryTransaction,
		Transaction: &tx,
	}
	newBlock, err := bc.generateNewBlock(data)
	if err != nil {
		return Block{}, err
	}
	return newBlock, nil
}

func (bc *Blockchain) getLastValidBlocks() []Block {
	currentHeight := len(bc.Blocks) - 1
	minBlockIndex := max(0, currentHeight-NUMBER_OF_BLOCKS_TO_SOLUTION)
	lastBlocksToCheck := bc.Blocks[minBlockIndex:]
	return lastBlocksToCheck
}

// check if the proposed solution is the best for the valid problems
func (bc *Blockchain) checkIfIsBestProposedSolution(proposedSolution *KnapsackProposedSolution) bool {

	lastBlocksToCheck := bc.getLastValidBlocks()

	for _, block := range lastBlocksToCheck {
		if block.Data.Type == KnapsackProposedSolutionSubmission && block.Data.Solution.ProblemBlockHeight == proposedSolution.ProblemBlockHeight && block.Data.Solution.Value >= proposedSolution.Value {
			return false // Found a better or equal solution, so return false
		}
	}

	// No better solution was found
	return true
}

// After adding a new block to the blockchain,
// check if an old problem was solved at current block height
func (bc *Blockchain) CheckForExpiredProblem() *ProblemSolutionPair {

	if len(bc.Blocks) < NUMBER_OF_BLOCKS_TO_SOLUTION {
		return nil
	}

	lastBlocksToCheck := bc.getLastValidBlocks()
	block := lastBlocksToCheck[0]
	if block.Data.Type != KnapsackProblemSubmission {
		return nil
	}

	problem := block.Data.Problem
	problemHeight := block.Height

	// Check if there is a solution for this problem in the subsequent blocks
	for j := 1; j < len(lastBlocksToCheck); j++ {
		solutionBlock := lastBlocksToCheck[j]
		if solutionBlock.Data.Type == KnapsackProposedSolutionSubmission &&
			solutionBlock.Data.Solution.ProblemBlockHeight == problemHeight {
			return &ProblemSolutionPair{
				Problem:  problem,
				Solution: solutionBlock.Data.Solution,
			}
		}
	}

	return nil
}

func (bc *Blockchain) FindValidProblemsBlocks() []Block {
	var problems []Block

	lastBlocksToCheck := bc.getLastValidBlocks()

	for _, block := range lastBlocksToCheck {
		if block.Data.Type == KnapsackProblemSubmission {
			problems = append(problems, block)
		}
	}

	// print problems found
	log.Printf("Found %v valid problems", len(problems))

	return problems
}

func (bc *Blockchain) validateTransaction(tx Transaction) error {
	if tx.Amount <= 0 {
		return errors.New("invalid transaction amount")
	}
	if tx.From == "" || tx.To == "" {
		return errors.New("invalid transaction address")
	}

	validProblemBlocks := bc.FindValidProblemsBlocks()
	if len(validProblemBlocks) == 0 {
		return errors.New("no valid problems to solve")
	}

	// check if the problem block height is valid
	if tx.ProblemBlockHeight < 0 || tx.ProblemBlockHeight >= len(bc.Blocks) {
		return errors.New("invalid problem block height")
	}

	// check if the problem block is a problem
	if bc.Blocks[tx.ProblemBlockHeight].Data.Type != KnapsackProblemSubmission {
		return errors.New("block at height does not contain a problem")
	}

	// check if problem block height is not expired
	validBlocks := bc.getLastValidBlocks()
	if tx.ProblemBlockHeight < validBlocks[0].Height {
		return errors.New("problem at block height is expired")
	}

	return nil
}
