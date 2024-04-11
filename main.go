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
			TxInfo: database.TxInfo{
				Hash:                 "0xc6a03a1cc3678f83d5a62de0bf8ca6f1fc0ee46aea90f0ea7a6c56890e0e0613",
				BlockNumber:          400024,
				From:                 "0x6177843db3138ae69679A54b95cf345ED759450d",
				To:                   "0x4977843db3138ae69679A54b95cf345ED7594143",
				Value:                "3140000000000000000",
				TotalGas:             4300,
				CodeChunkGas:         1200,
				ExecutedInstructions: 120,
				ExecutedBytes:        125,
				ChargedCodeChunks:    2,
			},
			WitnessTreeKeyValues: []database.WitnessTreeKeyValue{
				{Key: "0x02db80a48c552994cdcbd0fb08b9e76e59bfd3ea8c172f25c00e289dd4106f", PostValue: "0xdae346e3112395dc000000000000000000000000000000000000000000000000"},
				{Key: "0x0591491d8fd2bb468a49a02040a4872fed9791b954295973fc7eb2cfd4aef0", PostValue: "0x7440000000000000000000000000000000000000000000000000000000000000"},
			},
		},
	})

	explorer := server.New(":8181", db)
	if err := explorer.Run(); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}
