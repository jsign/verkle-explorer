package mock

import (
	"fmt"

	"github.com/jsign/verkle-explorer/database"
)

type mockDB struct {
	data []database.TxExec
}

func NewMockDB(data []database.TxExec) database.DB {
	return &mockDB{
		data: data,
	}
}

func (db *mockDB) GetTxExec(hash string) (database.TxExec, error) {
	for _, tx := range db.data {
		if tx.Hash == hash {
			fmt.Printf("tx.Hash: %s vs %s\n", tx.Hash, hash)
			return tx, nil
		}
	}
	return database.TxExec{}, database.ErrTxNotFound
}
