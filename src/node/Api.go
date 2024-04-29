package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func HomeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "SolverNet Blockchain API\n")
}

func HandleGetBlockchain(w http.ResponseWriter, r *http.Request, bc *Blockchain) {
	bytes, err := json.MarshalIndent(bc.Blocks, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func HandleGetLedger(w http.ResponseWriter, r *http.Request, ledger *Ledger) {
	bytes, err := json.MarshalIndent(ledger, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func HandleWriteProposedSolutionBlock(w http.ResponseWriter, r *http.Request, bc *Blockchain, ledger *Ledger) {
	log.Println("Received proposed solution block")
	WriteNewBlockData[KnapsackProposedSolution](w, r, bc.GenerateProposedSolutionBlock, bc, ledger)
}

func HandleWriteProblemBlock(w http.ResponseWriter, r *http.Request, bc *Blockchain, ledger *Ledger) {
	log.Println("Received proposed problem block")
	WriteNewBlockData[KnapsackProblem](w, r, bc.GenerateProblemBlock, bc, ledger)
}

func WriteNewBlockData[T any](w http.ResponseWriter, r *http.Request, generateBlock func(T) (Block, error), bc *Blockchain, ledger *Ledger) {
	w.Header().Set("Content-Type", "application/json")
	var data T
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		log.Println("Invalid decoded json")
		respondWithJSON(w, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	newBlock, err := generateBlock(data)
	if err != nil {
		log.Println("Invalid generated")
		respondWithJSON(w, http.StatusInternalServerError, r.Body)
		return
	}

	if !bc.isNewBlockCorrectlyChained(newBlock) {
		log.Println("Invalid block")
		respondWithJSON(w, http.StatusInternalServerError, "Invalid block")
		return
	}

	if err := bc.AddBlock(newBlock, ledger); err != nil {
		log.Println("New block has invalid data:", err)
		respondWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, newBlock)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}
