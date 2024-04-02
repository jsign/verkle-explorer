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

func (db *mockDB) GetHighestGasTxs(count int) ([]database.TxInfo, error) {
	return []database.TxInfo{db.data[0].TxInfo}, nil
}

func (db *mockDB) GetInefficientCodeAccessTxs(count int) ([]database.TxInfo, error) {
	return []database.TxInfo{db.data[0].TxInfo}, nil
}
