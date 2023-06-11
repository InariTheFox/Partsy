package mocks

import (
	"github.com/inarithefox/partsy/server/public/model"
	mock "github.com/stretchr/testify/mock"
)

type PartStore struct {
	mock.Mock
}

func (m *PartStore) GetAllParts(page, pageSize int) (model.PartList, error) {
	ret := m.Called(page, pageSize)

	var r0 model.PartList
	var r1 error

	if rf, ok := ret.Get(0).(func(int, int) (model.PartList, error)); ok {
		return rf(page, pageSize)
	}

	if rf, ok := ret.Get(0).(func(int, int) model.PartList); ok {
		r0 = rf(page, pageSize)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(model.PartList)
		}
	}

	if rf, ok := ret.Get(1).(func(int, int) error); ok {
		r1 = rf(page, pageSize)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *PartStore) GetAllChannelsCount() (int64, error) {
	ret := m.Called()

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func() (int64, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *PartStore) Save(p *model.Part) (*model.Part, error) {
	ret := m.Called(p)

	var r0 *model.Part
	var r1 error

	if rf, ok := ret.Get(0).(func(*model.Part) (*model.Part, error)); ok {
		return rf(p)
	}

	if rf, ok := ret.Get(0).(func(*model.Part) *model.Part); ok {
		r0 = rf(p)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Part)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.Part) error); ok {
		r1 = rf(p)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
