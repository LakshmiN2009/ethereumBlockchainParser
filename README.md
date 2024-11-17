---

# Ethereum Blockchain Parser

**Ethereum Blockchain Parser** is a Go-based application for parsing Ethereum blockchain data. This project includes APIs to retrieve the current block, manage subscriptions to Ethereum addresses, and fetch transactions for subscribed addresses. It supports concurrent transaction posting based on in-memory subscriber addresses.

## Project Structure

```
ethereumBlockchainParser/
├── main.go
├── service/
│   ├── ethereum.go
│   ├── ethereum_test.go
│   └── parserService.go
├── client/
│   ├── client.go
│   ├── ethereumClient_test.go
│   └── ethereumClient.go
├── driver/
│   ├── db.go
│   ├── memstore_test.go
│   └── memstore.go
├── util/
│    └── util.go
├── go.mod
└── go.sum
```

- **main.go**: Entry point for the application.
- **service/**: Contains core business logic, including parsing and transaction management.
- **client/**: Manages interactions with the Ethereum client.
- **driver/**: Provides database drivers and in-memory storage for subscriber addresses.
- **util/**: Helper utilities.

## API Endpoints

### `/getCurrentBlock`
Retrieves the current Ethereum block number.

- **URL**: `http://localhost:8080/getCurrentBlock`
- **Method**: `GET`
- **Response**:
  ```
    "Current block number": 21210602
  ```

### `/subscribe`
Subscribes to an Ethereum address for monitoring transactions.

- **URL**: `http://localhost:8080/subscribe?address=0x1f2f10d1c40777ae1da742455c65828ff36df387`
- **Method**: `GET`
- **Response**:
  ```
    "Address subscribed successfully: 0x1f2f10d1c40777ae1da742455c65828ff36df387"
  ```

### `/getTransactions`
Fetches transactions for a subscribed Ethereum address, categorizing them as inbound or outbound.

- **URL**: `http://localhost:8080/getTransactions?address=0x1f2f10d1c40777ae1da742455c65828ff36df387`
- **Method**: `GET`
- **Response**:
  ```json
  {
      "inbound": [
          {
              "hash": "0x935e0f53f38dd040dbc79df0a6cec7cfefb5c12aaef77d0cacc9a626219370b7",
              "blockHash": "0x2c4f5d038eb55c2a85156ef222a4538b700016eda1728a9e293ffd2a71b32040",
              "blockNumber": "0x143a4f4",
              "transactionIndex": "0x3",
              "from": "0xae2fc483527b8ef99eb5d9b44875f005ba1fae13",
              "to": "0x1f2f10d1c40777ae1da742455c65828ff36df387",
              "value": "0x37",
              "type": "0x2"
          }
      ],
      "outbound": [
         {
              "hash": "0xbba4bd754a03b6a128b0217019c113b12efa624331c40b688c95fbd2dbde049e",
              "blockHash": "0x2c4f5d038eb55c2a85156ef222a4538b700016eda1728a9e293ffd2a71b32040",
              "blockNumber": "0x143a4f4",
              "transactionIndex": "0x5",
              "from": "0xae2fc483527b8ef99eb5d9b44875f005ba1fae13",
              "to": "0x1f2f10d1c40777ae1da742455c65828ff36df387",
              "value": "0x37",
              "type": "0x2"
         }
      ]
  }
  ```

### `POST /postTransactions`
Pushes transactions concurrently for the addresses that are subscribed in memory.

## Getting Started

1. **Clone the repository**:
    ```bash
    git clone https://github.com/your-username/ethereumBlockchainParser.git
    cd ethereumBlockchainParser
    ```

2. **Install dependencies**:
    ```bash
    go mod download
    ```

3. **Run the application**:
    ```bash
    go run main.go
    ```

## Running Tests
To run tests and check coverage:

```bash
go test -cover ./...
```
