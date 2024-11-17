package client

import (
	"encoding/json"
	"ethereumBlockchainParser/util"
	"testing"
)

func TestGetCurrentBlock(t *testing.T) {

	client := NewEthereumClient().(*EthereumClient)
	// GetCurrentBlock ...
	blockNumber := client.GetCurrentBlock()
	if blockNumber == 0 {
		t.Errorf("Expected block number to be greater than 0, got %d", blockNumber)
	}

	client = NewEthereumClient().(*EthereumClient)
	client.EthUrl = ""
	block := client.GetCurrentBlock()
	if block != 0 {
		t.Errorf("Expected block number to be 0, got %d", block)
	}

}

func TestGetBlockByNumber(t *testing.T) {

	client := NewEthereumClient().(*EthereumClient)

	blockNumber := client.GetCurrentBlock()
	blockNumberHex := util.IntToHex(blockNumber)

	// GetBlockByNumber ... with valid block number
	client.GetBlockByNumber(blockNumberHex)
	if client.EthTransaction == nil {
		t.Errorf("Expected transactions to be present")
	}

	client = NewEthereumClient().(*EthereumClient)
	// GetBlockByNumber ... with invalid block number
	client.GetBlockByNumber("invalid")
	if client.EthTransaction != nil {
		t.Errorf("Expected transactions to be nil")
	}

	// GetBlockByNumber ... with empty URL
	client = NewEthereumClient().(*EthereumClient)
	client.EthUrl = ""
	client.GetBlockByNumber(blockNumberHex)
	if client.EthTransaction != nil {
		t.Errorf("Expected transactions to be nil")
	}
}

func TestSendRPCRequest(t *testing.T) {

	client := NewEthereumClient().(*EthereumClient)
	// SendRPCRequest ... with valid request
	result, err := client.sendRPCRequest(Eth_BlockNumber, []interface{}{})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	blockNumberHex := result

	// SendRPCRequest ... with invalid request
	result, _ = client.sendRPCRequest("invalid", []interface{}{})
	if result != nil {
		t.Errorf("Expected result to be nil")
	}

	client = NewEthereumClient().(*EthereumClient)
	// SendRPCRequest ... with valid blockNumberHex
	resp, err := client.sendRPCRequest(Eth_BlockByNumber, []interface{}{blockNumberHex, true})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	respResult, err := json.Marshal(resp)
	if err != nil {
		t.Errorf("Error marshalling response result: %v", err)
	}

	r := EthResult{}
	err = json.Unmarshal(respResult, &r)
	if err != nil {
		t.Errorf("Error unmarshalling response: %v", err)
	}

	if r.Transactions == nil {
		t.Errorf("Expected transactions to be present")
	}

	client = NewEthereumClient().(*EthereumClient)
	// SendRPCRequest ... with invalid blockNumberHex
	resp, err = client.sendRPCRequest(Eth_BlockByNumber, []interface{}{"invalid", true})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	respResult, err = json.Marshal(resp)
	if err != nil {
		t.Errorf("Error marshalling response result: %v", err)
	}

	r = EthResult{}
	err = json.Unmarshal(respResult, &r)
	if err != nil {
		t.Errorf("Error unmarshalling response: %v", err)
	}

	if r.Transactions != nil {
		t.Errorf("Expected result to be empty")
	}

	client = NewEthereumClient().(*EthereumClient)
	client.EthUrl = ""
	// SendRPCRequest ... with empty URL
	_, err = client.sendRPCRequest(Eth_BlockByNumber, []interface{}{"invalid", true})
	if err == nil {
		t.Errorf("Expected no error, got %v", err)
	}

}
