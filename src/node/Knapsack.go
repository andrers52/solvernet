package main

import (
	"errors"
	"log"
)

type KnapsackProblem struct {
	Items    []Item  `json:"items"`
	Capacity int     `json:"capacity"`
	Bounty   float64 `json:"bounty"`
	Address  string  `json:"address"` // address to send the bounty from
}

type KnapsackProposedSolution struct {
	ItemIndexes        []int  `json:"items"`
	ProblemBlockHeight int    `json:"problem_block_height"` // Identifies the block where the problem was submitted in
	Value              int    `json:"value"`                // This is what we are trying to maximize
	Address            string `json:"address"`              // address to send the bounty to
}

type ProblemSolutionPair struct {
	Problem  *KnapsackProblem          `json:"problem"`
	Solution *KnapsackProposedSolution `json:"solution"`
}

func GetProblemItemsSumWeight(problem KnapsackProblem) int {
	sum := 0
	for _, item := range problem.Items {
		sum += item.Weight
	}
	return sum
}

func GetTotalSolutionWeight(problem KnapsackProblem, proposedSolution KnapsackProposedSolution) int {
	weight := 0
	for _, index := range proposedSolution.ItemIndexes {
		weight += problem.Items[index].Weight
	}
	return weight
}

func GetTotalSolutionValue(problem KnapsackProblem, proposedSolution KnapsackProposedSolution) int {
	value := 0
	for _, index := range proposedSolution.ItemIndexes {
		value += problem.Items[index].Value
	}
	return value
}

func ValidateProblem(problem KnapsackProblem, bc *Blockchain) error {
	if problem.Bounty < 1.0 {
		return errors.New("bounty too low")
	}

	// check if problem has items
	if len(problem.Items) == 0 {
		return errors.New("no items in problem")
	}

	// check if problem has address
	if problem.Address == "" {
		return errors.New("no address in problem")
	}

	// check if problem has capacity
	if problem.Capacity < 1 {
		return errors.New("capacity too low")
	}

	// check if problem has items with negative/0 weight or value
	for _, item := range problem.Items {
		if item.Weight <= 0 || item.Value <= 0 {
			return errors.New("negative/0 weight or value")
		}
	}

	// check if total items weight is smaller than capacity.
	// This makes the problem trivial and not worth solving (solution is all items)
	if GetProblemItemsSumWeight(problem) <= problem.Capacity {
		return errors.New("total items weight is smaller than capacity. Trivial problem not allowed")
	}

	// (TODO) A node cannot submit a new problem if it do not have the amount of tokens to pay the bounty

	return nil
}

func ValidateProposedSolution(proposedSolution KnapsackProposedSolution, bc *Blockchain) error {
	if proposedSolution.ProblemBlockHeight >= len(bc.Blocks) {
		return errors.New("invalid proposed solution block height. Value too big")
	}
	// cannot submit a solution for a block that has already expired/solved
	if proposedSolution.ProblemBlockHeight < len(bc.Blocks)-NUMBER_OF_BLOCKS_TO_SOLUTION {
		return errors.New("invalid proposed solution block height. Value too small")
	}

	//check if solution has items
	if len(proposedSolution.ItemIndexes) == 0 {
		return errors.New("no items in solution")
	}

	//check if solution has address
	if proposedSolution.Address == "" {
		return errors.New("no address in solution")
	}

	block := bc.GetBlock(proposedSolution.ProblemBlockHeight)
	if block.Data.Type != KnapsackProblemSubmission {
		return errors.New("block at height does not contain a problem")
	}

	//check there is a problem at height position
	if block.Data.Problem == nil {
		return errors.New("block at height does not contain a problem")
	}

	problem := block.Data.Problem
	indexMap := make(map[int]bool)
	for _, i := range proposedSolution.ItemIndexes {
		indexMap[i] = true
	}
	if len(indexMap) != len(proposedSolution.ItemIndexes) {
		return errors.New("duplicate item indexes")
	}
	for _, index := range proposedSolution.ItemIndexes {
		if index < 0 || index >= len(problem.Items) {
			return errors.New("invalid item index")
		}
	}
	weight := GetTotalSolutionWeight(*problem, proposedSolution)

	if weight > problem.Capacity {
		return errors.New("solution exceeds capacity")
	}

	value := GetTotalSolutionValue(*problem, proposedSolution)
	if value != proposedSolution.Value {
		log.Println("Invalid proposed solution", value, proposedSolution.Value)
		return errors.New("solution value does not match")
	}

	if !bc.checkIfIsBestProposedSolution(&proposedSolution) {
		return errors.New("solution is not better than previous solution")
	}

	return nil
}
