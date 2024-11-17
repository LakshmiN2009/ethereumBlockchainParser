package client

type ClientInterface interface {
	GetCurrentBlock() int64
	GetBlockByNumber(blockNumberHex string)
}
