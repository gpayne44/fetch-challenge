package repositories

import (
	"errors"

	"github.com/google/uuid"
	"github.com/gpayne44/fetch-challenge/internal/entities"
)

type memoryStore struct {
	data map[uuid.UUID]entities.ReceiptRecord
}

type ReceiptsRepository interface {
	StoreReceipt(r entities.ReceiptRecord) (string, error)
	GetReceipt(id uuid.UUID) (*entities.ReceiptRecord, error)
}

var ErrNotFound = errors.New("entity not found")

func New() *memoryStore {
	dataMap := make(map[uuid.UUID]entities.ReceiptRecord)
	m := memoryStore{
		data: dataMap,
	}
	return &m
}

func (m memoryStore) StoreReceipt(r entities.ReceiptRecord) (string, error) {
	var id string

	newID := uuid.New()
	id = newID.String()

	m.data[newID] = r
	return id, nil
}

func (m memoryStore) GetReceipt(id uuid.UUID) (*entities.ReceiptRecord, error) {
	receipt, ok := m.data[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &receipt, nil
}
