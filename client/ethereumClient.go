package client

import (
	"bytes"
	"encoding/json"
	"ethereumBlockchainParser/util"
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

type EthereumClient struct {
	Client         *http.Client
	EthUrl         string
	EthTransaction []EthTransaction
}

type EthereumRequest struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
}

type EthereumResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Id      int         `json:"id"`
}

type EthResult struct {
	Hash         string           `json:"hash"`
	Transactions []EthTransaction `json:"transactions"`
	Error        ErrorResponse    `json:"error"`
}

type EthTransaction struct {
	Hash             string `json:"hash"`
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	TransactionIndex string `json:"transactionIndex"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            string `json:"value"`
	Type             string `json:"type"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

const (
	Eth_BlockNumber   = "eth_blockNumber"
	Eth_BlockByNumber = "eth_getBlockByNumber"
)

func NewEthereumClient() ClientInterface {
	return &EthereumClient{
		Client: &http.Client{},
		EthUrl: "https://ethereum-rpc.publicnode.com",
	}
}

func (e *EthereumClient) GetCurrentBlock() int64 {

	// get the current block number
	result, err := e.sendRPCRequest(Eth_BlockNumber, []interface{}{})
	if err != nil {
		fmt.Println("EthereumClient.GetCurrentBlock Error getting current block: ", err)
		return 0
	}

	blockNumber := util.HexToInt(result.(string))

	return blockNumber
}

func (e *EthereumClient) GetBlockByNumber(blockNumberHex string) {

	resp, err := e.sendRPCRequest(Eth_BlockByNumber, []interface{}{blockNumberHex, true})
	if err != nil {
		return
	}

	ethResult, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("Error marshalling response result: ", err)
		return
	}

	result := EthResult{}

	err = json.Unmarshal(ethResult, &result)
	if err != nil {
		fmt.Println("Error unmarshalling response: ", err)
		return
	}

	if result.Error != (ErrorResponse{}) {
		fmt.Println("Error getting block by number: ", result.Error.Message)
		return
	}

	e.EthTransaction = append(e.EthTransaction, result.Transactions...)

	return
}

func (e *EthereumClient) sendRPCRequest(method string, params []interface{}) (interface{}, error) {

	body, err := json.Marshal(EthereumRequest{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		ID:      rand.Intn(100),
	})

	if err != nil {
		fmt.Println("Error marshalling request body: ", err)
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, e.EthUrl, bytes.NewReader(body))
	if err != nil {
		fmt.Println("Error creating request: ", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.Client.Do(req)
	if err != nil {
		fmt.Println("Error sending request: ", err)
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body: ", err)
		return nil, err
	}

	var response EthereumResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		fmt.Println("Error unmarshalling response: ", err)
		return nil, err
	}

	return response.Result, nil
}
