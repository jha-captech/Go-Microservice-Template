// Code generated by mockery v2.43.2. DO NOT EDIT.

package mock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	models "github.com/captechconsulting/go-microservice-templates/lambda/internal/models"
)

// MockUserService is an autogenerated mock type for the userService type
type MockUserService struct {
	mock.Mock
}

type MockUserService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockUserService) EXPECT() *MockUserService_Expecter {
	return &MockUserService_Expecter{mock: &_m.Mock}
}

// ListUsers provides a mock function with given fields: ctx
func (_m *MockUserService) ListUsers(ctx context.Context) ([]models.User, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for ListUsers")
	}

	var r0 []models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]models.User, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []models.User); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockUserService_ListUsers_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListUsers'
type MockUserService_ListUsers_Call struct {
	*mock.Call
}

// ListUsers is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockUserService_Expecter) ListUsers(ctx interface{}) *MockUserService_ListUsers_Call {
	return &MockUserService_ListUsers_Call{Call: _e.mock.On("ListUsers", ctx)}
}

func (_c *MockUserService_ListUsers_Call) Run(run func(ctx context.Context)) *MockUserService_ListUsers_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockUserService_ListUsers_Call) Return(_a0 []models.User, _a1 error) *MockUserService_ListUsers_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUserService_ListUsers_Call) RunAndReturn(run func(context.Context) ([]models.User, error)) *MockUserService_ListUsers_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateUser provides a mock function with given fields: ctx, ID, user
func (_m *MockUserService) UpdateUser(ctx context.Context, ID int, user models.User) (models.User, error) {
	ret := _m.Called(ctx, ID, user)

	if len(ret) == 0 {
		panic("no return value specified for UpdateUser")
	}

	var r0 models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, models.User) (models.User, error)); ok {
		return rf(ctx, ID, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, models.User) models.User); ok {
		r0 = rf(ctx, ID, user)
	} else {
		r0 = ret.Get(0).(models.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, models.User) error); ok {
		r1 = rf(ctx, ID, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockUserService_UpdateUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateUser'
type MockUserService_UpdateUser_Call struct {
	*mock.Call
}

// UpdateUser is a helper method to define mock.On call
//   - ctx context.Context
//   - ID int
//   - user models.User
func (_e *MockUserService_Expecter) UpdateUser(ctx interface{}, ID interface{}, user interface{}) *MockUserService_UpdateUser_Call {
	return &MockUserService_UpdateUser_Call{Call: _e.mock.On("UpdateUser", ctx, ID, user)}
}

func (_c *MockUserService_UpdateUser_Call) Run(run func(ctx context.Context, ID int, user models.User)) *MockUserService_UpdateUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int), args[2].(models.User))
	})
	return _c
}

func (_c *MockUserService_UpdateUser_Call) Return(_a0 models.User, _a1 error) *MockUserService_UpdateUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUserService_UpdateUser_Call) RunAndReturn(run func(context.Context, int, models.User) (models.User, error)) *MockUserService_UpdateUser_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockUserService creates a new instance of MockUserService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUserService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUserService {
	mock := &MockUserService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
