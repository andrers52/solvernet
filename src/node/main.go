package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	blockchain := CreateNewBlockchain()

	port := os.Getenv("PORT") // Default port is set in .env file
	if len(os.Args) > 1 {     // If a port is passed as an argument, use that instead
		port = os.Args[1]
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/heartbeat", HomeLink).Methods("GET")
	router.HandleFunc("/api/home", HomeLink).Methods("GET")
	router.HandleFunc("/api/get_blockchain", func(w http.ResponseWriter, r *http.Request) {
        HandleGetBlockchain(w, r, blockchain)
    }).Methods("GET")
	router.HandleFunc("/api/get_current_state", HandleGetLedger).Methods("GET")
	router.HandleFunc("/api/send_problem", func(w http.ResponseWriter, r *http.Request) {
        HandleWriteProblemBlock(w, r, blockchain)
    }).Methods("POST")
	router.HandleFunc("/api/send_proposed_solution", func(w http.ResponseWriter, r *http.Request) {
        HandleWriteProposedSolutionBlock(w, r, blockchain)
    }).Methods("POST")

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	log.Println("now serving on ", port)
	node := InitNode(port)
	go node.StartNode(blockchain)
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(headers, methods, origins)(router)))
}
