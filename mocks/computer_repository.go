// Code generated by mockery v2.53.0. DO NOT EDIT.

package mocks

import (
	models "computer-club/internal/models"

	mock "github.com/stretchr/testify/mock"
)

// ComputerRepository is an autogenerated mock type for the ComputerRepository type
type ComputerRepository struct {
	mock.Mock
}

// GetComputers provides a mock function with no fields
func (_m *ComputerRepository) GetComputers() ([]models.Computer, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetComputers")
	}

	var r0 []models.Computer
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]models.Computer, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []models.Computer); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Computer)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateStatus provides a mock function with given fields: number, free
func (_m *ComputerRepository) UpdateStatus(number int, free models.ComputerStatus) error {
	ret := _m.Called(number, free)

	if len(ret) == 0 {
		panic("no return value specified for UpdateStatus")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(int, models.ComputerStatus) error); ok {
		r0 = rf(number, free)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewComputerRepository creates a new instance of ComputerRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewComputerRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *ComputerRepository {
	mock := &ComputerRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
