package memstore

import (
	"context"

	"github.com/inarithefox/partsy/server/public/logger"
	"github.com/inarithefox/partsy/server/store"
)

type MemStore struct {
	context context.Context
	stores  MemStores
}

type MemStores struct {
	parts store.PartStore
}

func New() *MemStore {
	store := &MemStore{}

	logger.Debug("new root in-memory store initialized")

	store.stores.parts = NewMemPartStore(store)

	return store
}

func (m *MemStore) Close() {

}

func (m *MemStore) Context() context.Context {
	return m.context
}

func (m *MemStore) Part() store.PartStore {
	return m.stores.parts
}

func (m *MemStore) SetContext(context context.Context) {
	m.context = context
}
