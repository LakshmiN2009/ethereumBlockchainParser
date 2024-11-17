package parserService

import (
	"context"
	"ethereumBlockchainParser/client"
	"ethereumBlockchainParser/driver"
	"ethereumBlockchainParser/util"
	"sync"
)

type EthereumService struct {
	mx     *sync.Mutex
	client *client.EthereumClient
	db     *driver.InMemStore
}

// NewEthereumParserService ... create a new instance of EthereumService
func NewEthereumParserService() ParserService {
	return &EthereumService{mx: &sync.Mutex{}, client: client.NewEthereumClient().(*client.EthereumClient), db: driver.NewInMemStore().(*driver.InMemStore)}
}

// GetCurrentBlock ... last parsed block
func (e *EthereumService) GetCurrentBlock(ctx context.Context) int64 {

	e.mx.Lock()
	defer e.mx.Unlock()

	// get the current block number
	return e.client.GetCurrentBlock()
}

// Subscribe ... add address to observer
func (e *EthereumService) Subscribe(ctx context.Context, address string) bool {

	e.mx.Lock()
	defer e.mx.Unlock()

	// check if the address is already present in the subscriber list
	_, exist := e.db.Subscribers[address]
	if exist {
		return false
	}

	// add the address to the subscriber list
	e.db.Subscribers[address] = true

	return true
}

// GetTransactions ... list of inbound or outbound transactions for an address
func (e *EthereumService) GetTransactions(ctx context.Context, address string) driver.Transaction {

	e.mx.Lock()
	defer e.mx.Unlock()

	// check if the address is present in the subscriber list
	_, exists := e.db.Subscribers[address]
	if !exists {
		return driver.Transaction{}
	}

	// get the transaction for the address
	transaction, err := e.db.Get(address)
	if err != nil {
		return driver.Transaction{}
	}

	return transaction
}

// PostTransaction ... post the transactions if the subscriber is present
func (e *EthereumService) PostTransaction(ctx context.Context) bool {

	e.mx.Lock()
	defer e.mx.Unlock()

	// get the current block number
	result := e.client.GetCurrentBlock()
	if result == 0 {
		return false
	}

	// convert the block number to hex
	blockNumHex := util.IntToHex(result)

	// get the block by number - using blockNumHex
	e.client.GetBlockByNumber(blockNumHex)

	count := 0
	txns := e.client.EthTransaction
	for _, txn := range txns {
		// check if the to address is available in the subscriber list
		_, exist := e.db.Subscribers[txn.To]
		if exist {
			// add the transaction to the inbound list
			t := driver.Transaction{}
			var ok bool
			if t, ok = e.db.Transcation[txn.To]; ok {
				t.InBound = append(t.InBound, txn)
			} else {
				t.InBound = append(t.InBound, txn)
				t.OutBound = []client.EthTransaction{}
			}

			e.db.Insert(txn.To, t)
			count++
		}

		// check if the to address is available in the subscriber list
		_, exist = e.db.Subscribers[txn.From]
		if exist {
			// add the transaction to the outbound list
			t := driver.Transaction{}
			var ok bool
			if t, ok = e.db.Transcation[txn.From]; ok {
				t.OutBound = append(t.OutBound, txn)
			} else {
				t.OutBound = append(t.OutBound, txn)
				t.InBound = []client.EthTransaction{}
			}

			e.db.Insert(txn.From, t)
			count++
		}
	}

	if count == 0 {
		return false
	}

	return true
}
