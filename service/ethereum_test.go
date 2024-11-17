package parserService

import (
	"context"
	"encoding/json"
	"ethereumBlockchainParser/driver"
	"ethereumBlockchainParser/util"
	"testing"
)

func TestGetCurrentBlock(t *testing.T) {
	service := NewEthereumParserService()
	ctx := context.Background()

	block := service.GetCurrentBlock(ctx)
	if block == 0 {
		t.Errorf("Expected block number to be 0, got %d", block)
	}
}

func TestSubscribe(t *testing.T) {
	service := NewEthereumParserService()
	ctx := context.Background()
	address := "0x123"

	subscribed := service.Subscribe(ctx, address)
	if !subscribed {
		t.Errorf("Expected subscription to be successful")
	}

	// Test subscribing the same address again
	subscribed = service.Subscribe(ctx, address)
	if subscribed {
		t.Errorf("Expected subscription to fail for already subscribed address")
	}
}

func TestGetTransactions(t *testing.T) {
	service := NewEthereumParserService().(*EthereumService)
	ctx := context.Background()
	address := "0x123"

	// TestCase1:  Test getting transactions for an unsubscribed address
	tx := service.GetTransactions(ctx, address)
	if len(tx.InBound) != 0 || len(tx.OutBound) != 0 {
		t.Errorf("Expected no transactions for unsubscribed address")
	}

	// Subscribe the address and test again
	service.Subscribe(ctx, address)
	// TestCase2: Test getting transactions for an subscribed address with out any transactions
	tx = service.GetTransactions(ctx, address)
	if len(tx.InBound) != 0 || len(tx.OutBound) != 0 {
		t.Errorf("Expected no transactions for unsubscribed address")
	}

	// add a transaction for the address
	txn := `{"inbound":[{"hash":"0x3fd378832d4d53b867946a8fc97175d23f8e343d1c6b03906a744118784fcf88","blockHash":"0x5cb858307cfc30f61a18ecde9c83bfab7cb59cdde0473dbac4632093f077470f","blockNumber":"0x143a0eb","transactionIndex":"0x0","from":"0x647a7b4bb0503d8bbb7bfeab04dba2ccdd93829e","to":"0x123","value":"0x0","type":"0x2"},{"hash":"0x36b0e2c97dea03a08fcdbb531358ccc39f39673aeb6e19c296fe89ef7e1c1088","blockHash":"0x5cb858307cfc30f61a18ecde9c83bfab7cb59cdde0473dbac4632093f077470f","blockNumber":"0x143a0eb","transactionIndex":"0xe","from":"0x01d00571208b10c02d08bbe159c4d3710354704a","to":"0x123","value":"0xf8b0a10e47000a","type":"0x2"}],"outbound":[]}`

	trans := driver.Transaction{}

	err := json.Unmarshal([]byte(txn), &trans)
	if err != nil {
		t.Errorf("Error unmarshalling transaction: %v", err)
	}

	service.db.Transcation[address] = trans

	// TestCase3: Test getting transactions for a subscribed address
	tx = service.GetTransactions(ctx, address)
	if len(tx.InBound) == 0 && len(tx.OutBound) == 0 {
		t.Errorf("Expected no transactions for subscribed address")
	}
}

func TestPostTransaction(t *testing.T) {
	service := NewEthereumParserService().(*EthereumService)
	ctx := context.Background()

	// TestCase1: Test posting transaction without any subscribers
	result := service.PostTransaction(ctx)
	if result {
		t.Errorf("Expected posting transaction to fail without any subscribers")
	}

	// Subscribe an address and test again
	address := "0x123"
	service.Subscribe(ctx, address)

	// TestCase2: Test posting transaction with subscribers
	result = service.PostTransaction(ctx)
	if result {
		t.Errorf("Expected posting transaction to be successful with subscribers")
	}

	// TestCase3: Test posting transaction with subscribers and no transactions
	service = NewEthereumParserService().(*EthereumService)

	blockNumber := service.GetCurrentBlock(ctx)
	blockHex := util.IntToHex(blockNumber)

	service.client.GetBlockByNumber(blockHex)
	txn := service.client.EthTransaction

	if txn == nil {
		t.Errorf("Expected transactions to be available")
	}

	address = txn[0].To
	service.Subscribe(ctx, address)

	result = service.PostTransaction(ctx)
	if !result {
		t.Errorf("Expected posting transaction to be successful with subscribers and transactions")
	}

	if len(service.db.Transcation) == 0 {
		t.Errorf("Expected transactions to be available")
	}

}
