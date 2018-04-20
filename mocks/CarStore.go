// Code generated by mockery v1.0.0
package mocks

import mock "github.com/stretchr/testify/mock"
import models "github.com/ucladevx/BPool/models"
import postgres "github.com/ucladevx/BPool/stores/postgres"

import stores "github.com/ucladevx/BPool/stores"

// CarStore is an autogenerated mock type for the CarStore type
type CarStore struct {
	mock.Mock
}

// GetAll provides a mock function with given fields: lastID, limit
func (_m *CarStore) GetAll(lastID string, limit int) ([]*models.Car, error) {
	ret := _m.Called(lastID, limit)

	var r0 []*models.Car
	if rf, ok := ret.Get(0).(func(string, int) []*models.Car); ok {
		r0 = rf(lastID, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Car)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, int) error); ok {
		r1 = rf(lastID, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: id
func (_m *CarStore) GetByID(id string) (*models.Car, error) {
	ret := _m.Called(id)

	var r0 *models.Car
	if rf, ok := ret.Get(0).(func(string) *models.Car); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Car)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByWhere provides a mock function with given fields: fields, queryModifiers
func (_m *CarStore) GetByWhere(fields []string, queryModifiers []stores.QueryModifier) ([]postgres.CarRow, error) {
	ret := _m.Called(fields, queryModifiers)

	var r0 []postgres.CarRow
	if rf, ok := ret.Get(0).(func([]string, []stores.QueryModifier) []postgres.CarRow); ok {
		r0 = rf(fields, queryModifiers)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]postgres.CarRow)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]string, []stores.QueryModifier) error); ok {
		r1 = rf(fields, queryModifiers)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCount provides a mock function with given fields: queryModifiers
func (_m *CarStore) GetCount(queryModifiers []stores.QueryModifier) (int, error) {
	ret := _m.Called(queryModifiers)

	var r0 int
	if rf, ok := ret.Get(0).(func([]stores.QueryModifier) int); ok {
		r0 = rf(queryModifiers)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]stores.QueryModifier) error); ok {
		r1 = rf(queryModifiers)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Insert provides a mock function with given fields: user
func (_m *CarStore) Insert(user *models.Car) error {
	ret := _m.Called(user)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Car) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Remove provides a mock function with given fields: id
func (_m *CarStore) Remove(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
