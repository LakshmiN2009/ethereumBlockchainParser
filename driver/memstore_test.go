package driver

import (
	"ethereumBlockchainParser/client"
	"testing"
)

func TestInsert(t *testing.T) {
	db := NewInMemStore().(*InMemStore)

	err := db.Insert("key", Transaction{InBound: []client.EthTransaction{}, OutBound: []client.EthTransaction{}})
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if len(db.Transcation) == 0 {
		t.Errorf("Expected 1 transaction, got %d", len(db.Transcation))
	}

}

func TestGet(t *testing.T) {
	db := NewInMemStore().(*InMemStore)

	_, err := db.Get("key")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	db.Transcation["key"] = Transaction{InBound: []client.EthTransaction{}, OutBound: []client.EthTransaction{}}

	_, err = db.Get("key")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
