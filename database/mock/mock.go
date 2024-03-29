package mock

import (
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
			return tx, nil
		}
	}
	return database.TxExec{}, database.ErrTxNotFound
}
