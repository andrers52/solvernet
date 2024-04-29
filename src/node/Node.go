package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Node struct {
	Port string
}

func InitNode(port string) *Node {
	return &Node{Port: port}
}
func (n *Node) checkOnline() error {
	for _, node := range AllNodePorts {
		if node == n.Port {
			continue
		}
		resp, err := http.Get("http://localhost:" + node + "/api/heartbeat") // "http://localhost:3001/api/heartbeat
		if err != nil {
			return err
		}
		err = resp.Body.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *Node) StartNode(bc *Blockchain) {
	for {
		err := n.checkOnline()
		if err != nil {
			log.Println(err)
		} else {
			break
		}
		time.Sleep(5 * time.Second)
	}
	log.Println("ONLINE")
	for {
		// Sleep a random amount of time
		time.Sleep(time.Duration(rand.Intn(10)+1) * time.Second)
		log.Println("About to check if we should submit a problem or find a solution")
		if rand.Intn(10) == 0 {
			log.Println("About to submit a problem")
			err := n.submitProblem()
			if err != nil {
				log.Println(err)
			}
		} else {
			log.Println("About to submit a proposed solution")
			n.submitProposedSolution(bc)
		}
	}
}

func (n *Node) submitProposedSolution(bc *Blockchain) {

	// Check if there are any problems to solve
	validProblemsBlocks := bc.FindValidProblemsBlocks()

	if len(validProblemsBlocks) == 0 {
		log.Println("No valid problems blocks found. Aborting...")
		return
	}

	// Select a random problem to solve
	randomProblemBlock := validProblemsBlocks[rand.Intn(len(validProblemsBlocks))]

	problemHeight := randomProblemBlock.Height

	// Generate a proposed solution
	solutionItems := make([]int, 0)

	for i := range randomProblemBlock.Data.Problem.Items {
		if rand.Intn(2) == 1 { // 50% chance to include the item
			solutionItems = append(solutionItems, i)
		}
	}

	// Create a new proposed solution
	newSolution := KnapsackProposedSolution{
		ItemIndexes:        solutionItems,
		ProblemBlockHeight: problemHeight, // related problem block height
		Value:              0,             // Value will be calculated later
		Address:            n.Port,        // Assuming the node's port is used as its identity/address
	}

	totalValue := GetTotalSolutionValue(*randomProblemBlock.Data.Problem, newSolution)

	// set the value of the solution
	newSolution.Value = totalValue

	// Check if the solution is valid
	if err := ValidateProposedSolution(newSolution, bc); err != nil {
		log.Println("Generated solution is invalid:", err, ". Discarding...")
		return
	}

	// Marshal the new solution into JSON
	jsonPayload, err := json.Marshal(newSolution)
	if err != nil {
		log.Println("Error encoding JSON:", err)
		return
	}

	// Submit the new solution to all nodes
	for _, node := range AllNodePorts {
		if node == n.Port {
			continue
		}
		log.Println("Submitting proposed solution to node", node)
		resp, err := http.Post("http://localhost:"+node+"/api/send_proposed_solution", "application/json", bytes.NewBuffer(jsonPayload))
		if err != nil {
			log.Println("Failed to submit proposed solution:", err)
			continue
		}
		err = resp.Body.Close()
		if err != nil {
			log.Println("Error closing response body:", err)
		}
	}

	log.Println("Proposed solution submitted")
}

func (n *Node) submitProblem() error {
	log.Println("Creating a new problem")
	problem := KnapsackProblem{}
	// TODO: Ensure we don't offer more than we have in our balance
	bounty := rand.Float64() * 10
	if bounty == 0 {
		log.Println("No bounty, not submitting problem")
		return nil
	}
	problem.Bounty = bounty
	problem.Address = n.Port
	problem.Items = make([]Item, rand.Intn(10)+1)
	for i := range problem.Items {
		item := Item{
			Value:  rand.Intn(10) + 1,
			Weight: rand.Intn(10) + 1,
		}
		problem.Items[i] = item
	}

	sumOfWeights := GetProblemItemsSumWeight(problem)
	problem.Capacity = int(sumOfWeights * 2 / 3)

	jsonPayload, err := json.Marshal(problem)
	if err != nil {
		log.Println("Error encoding JSON:", err)
		return err
	}

	for _, node := range AllNodePorts {
		if node == n.Port {
			continue
		}
		log.Println("SUBMITTING PROBLEM", problem, "TO", node)
		resp, err := http.Post("http://localhost:"+node+"/api/send_problem", "application/json", bytes.NewBuffer(jsonPayload))
		if err != nil {
			return err
		}
		err = resp.Body.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
