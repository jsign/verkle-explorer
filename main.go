package main

import (
	"log"

	"github.com/jsign/verkle-explorer/database"
	"github.com/jsign/verkle-explorer/database/mock"
	"github.com/jsign/verkle-explorer/server"
)

func main() {
	db := mock.NewMockDB([]database.TxExec{
		{
			Hash:                 "0xc6a03a1cc3678f83d5a62de0bf8ca6f1fc0ee46aea90f0ea7a6c56890e0e0613",
			TotalGas:             4300,
			CodeChunkGas:         1200,
			ExecutedInstructions: 120,
			ExecutedBytes:        125,
			ChargedBytes:         62,
			Events: []database.WitnessEvent{
				{Name: "ContractInit", Gas: 100},
				{Name: "TouchFullAddress", Gas: 200},
				{Name: "TouchAddressOnWrite", Gas: 500},
				{Name: "ContractInitiCompletion", Gas: 800},
			},
		},
	})

	explorer := server.New(":8181", db)
	if err := explorer.Run(); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}
