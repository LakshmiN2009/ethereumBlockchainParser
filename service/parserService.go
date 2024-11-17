package parserService

import (
	"context"
	"ethereumBlockchainParser/driver"
)

type ParserService interface {
	// last parsed block
	GetCurrentBlock(ctx context.Context) int64
	// add address to observer
	Subscribe(ctx context.Context, address string) bool
	// list of inbound or outbound transactions for an address
	GetTransactions(ctx context.Context, address string) driver.Transaction

	// post the transactions if the subscriber is present
	PostTransaction(ctx context.Context) bool
}
