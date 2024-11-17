package main

import (
	"context"
	"encoding/json"
	parserService "ethereumBlockchainParser/service"
	"fmt"
	"net/http"
	"time"
)

func main() {

	e := parserService.NewEthereumParserService()

	ctx := context.Background()

	http.HandleFunc("/getCurrentBlock", func(w http.ResponseWriter, r *http.Request) {
		blockNumber := e.GetCurrentBlock(ctx)
		if blockNumber == 0 {
			http.Error(w, "Error getting block number", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Current block number: %d", blockNumber)
	})

	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		if address == "" {
			http.Error(w, "Address is required", http.StatusBadRequest)
			return
		}

		if !e.Subscribe(ctx, address) {
			http.Error(w, "Address already found! "+address, http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Address subscriber successfully: %s", address)

	})

	http.HandleFunc("/getTransactions", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		if address == "" {
			http.Error(w, "Address is required", http.StatusBadRequest)
			return
		}

		transaction := e.GetTransactions(ctx, address)
		if transaction.InBound == nil && transaction.OutBound == nil {
			http.Error(w, "Transaction not found for address: "+address, http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(transaction)
	})

	exit := make(chan struct{})
	go func() {
		defer close(exit)
		for {
			select {
			case <-exit:
				return
			default:
				success := e.PostTransaction(ctx)
				fmt.Println("Added transactions to the subscribers: ", success)
				time.Sleep(30 * time.Second)
			}
		}
	}()

	http.ListenAndServe(":8080", nil)

	<-exit

}
