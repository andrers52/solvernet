# Solvers of the World, Unite!

# The SolverNet Blockchain Project

## Overview

SolverNet is a minimalistic blockchain implementation designed to utilize a "proof of useful work" mechanism where nodes in the network solve instances of the knapsack problem as part of the consensus process. This project serves as a proof of concept that combines traditional blockchain technology with real-world problem solving to enhance the utility of computational work done by blockchain networks.

## Features

- **Blockchain Basics**: Implements basic blockchain structures including problems, proposed solutions and transactions.
- **Proof of Useful Work**: Uses clients defined problems, initially only the knapsack problem, as the basis for mining new blocks, replacing traditional proof-of-work systems.
- **Transaction System**: Handles both monetary transactions and problem submissions within the network.
- **RESTful API**: Provides endpoints for interacting with the blockchain, submitting problems, and viewing the chain state.

## Getting Started

### Prerequisites

- Install Go (version 1.22.2 or later recommended) from [The Go Programming Language site](https://golang.org/dl/).

### Installation

1. Clone the repository:

```bash
git clone git@github.cbhq.net:a-silva/solvernet.git
cd solvernet
```

2. Build the project:

```bash
 go build
```

3. Run tests:

```bash
go test -v
```

4. Run the server:

```bash
 ./solvernet
```

## API Endpoints

- GET /api/getblockchain: Fetches the entire blockchain.
- POST /api/sendtransaction: Submits a new transaction to the blockchain. This can be either a monetary transaction or a problem submission.

### Example Usage

To submit a new transaction via curl:

```bash
curl -X POST http://localhost:3002/api/sendtransaction -H 'Content-Type: application/json' -d '{
    "type": 0,  // MonetaryTransaction
    "transaction": {
        "from": "user1",
        "to": "user2",
        "amount": 100
    },
    "problem": null
}'
```

To submit a new problem via curl:

```bash
curl -X POST http://localhost:3002/api/sendtransaction -H 'Content-Type: application/json' -d '{
    "type": 1,  // KnapsackProblemSubmission
    "transaction": null,
    "problem": {
        "items": [
            {"id": 1, "weight": 5, "value": 10},
            {"id": 2, "weight": 3, "value": 6},
            {"id": 3, "weight": 2, "value": 3}
        ],
        "capacity": 10,
        "bounty": 500,
        "deadline": "2024-12-31T23:59:59Z"
    }
}'
```

To retrieve the blockchain state:

```bash
curl http://localhost:3002/api/getblockchain
```

To start with the 3 nodes configuration (only network supported so far) just open 3 different terminal tabs and run:

```bash
./solvernet 3001
./solvernet 3002
./solvernet 3003
```

## Contributing

Contributions to SolverNet are welcome! Please feel free to fork the repository, make changes, and submit pull requests. You can also open issues in the project's repository if you find bugs or have feature suggestions.

## License

This project is licensed under the MIT License - see the LICENSE.md file for details.
