package mocks

import (
	"github.com/inarithefox/partsy/server/store"
	mock "github.com/stretchr/testify/mock"
)

type Store struct {
	mock.Mock
}

func (m *Store) Close() {
	m.Called()
}

func (m *Store) Part() store.PartStore {
	ret := m.Called()

	var r store.PartStore
	if rf, ok := ret.Get(0).(func() store.PartStore); ok {
		r = rf()
	} else {
		if ret.Get(0) != nil {
			r = ret.Get(0).(store.PartStore)
		}
	}

	return r
}
