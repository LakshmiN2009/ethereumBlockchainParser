package driver

type DB interface {
	Insert(key string, value Transaction) error
	Get(key string) (Transaction, error)
}
