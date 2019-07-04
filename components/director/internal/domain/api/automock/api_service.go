// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import context "context"
import mock "github.com/stretchr/testify/mock"
import model "github.com/kyma-incubator/compass/components/director/internal/model"

// APIService is an autogenerated mock type for the APIService type
type APIService struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, applicationID, in
func (_m *APIService) Create(ctx context.Context, applicationID string, in model.APIDefinitionInput) (string, error) {
	ret := _m.Called(ctx, applicationID, in)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, model.APIDefinitionInput) string); ok {
		r0 = rf(ctx, applicationID, in)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, model.APIDefinitionInput) error); ok {
		r1 = rf(ctx, applicationID, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, id
func (_m *APIService) Delete(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAPIAuth provides a mock function with given fields: ctx, apiID, runtimeID
func (_m *APIService) DeleteAPIAuth(ctx context.Context, apiID string, runtimeID string) (*model.RuntimeAuth, error) {
	ret := _m.Called(ctx, apiID, runtimeID)

	var r0 *model.RuntimeAuth
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *model.RuntimeAuth); ok {
		r0 = rf(ctx, apiID, runtimeID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.RuntimeAuth)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, apiID, runtimeID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: ctx, id
func (_m *APIService) Get(ctx context.Context, id string) (*model.APIDefinition, error) {
	ret := _m.Called(ctx, id)

	var r0 *model.APIDefinition
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.APIDefinition); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.APIDefinition)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RefetchAPISpec provides a mock function with given fields: ctx, id
func (_m *APIService) RefetchAPISpec(ctx context.Context, id string) (*model.APISpec, error) {
	ret := _m.Called(ctx, id)

	var r0 *model.APISpec
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.APISpec); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.APISpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetAPIAuth provides a mock function with given fields: ctx, apiID, runtimeID, in
func (_m *APIService) SetAPIAuth(ctx context.Context, apiID string, runtimeID string, in model.AuthInput) (*model.RuntimeAuth, error) {
	ret := _m.Called(ctx, apiID, runtimeID, in)

	var r0 *model.RuntimeAuth
	if rf, ok := ret.Get(0).(func(context.Context, string, string, model.AuthInput) *model.RuntimeAuth); ok {
		r0 = rf(ctx, apiID, runtimeID, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.RuntimeAuth)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, model.AuthInput) error); ok {
		r1 = rf(ctx, apiID, runtimeID, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, id, in
func (_m *APIService) Update(ctx context.Context, id string, in model.APIDefinitionInput) error {
	ret := _m.Called(ctx, id, in)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, model.APIDefinitionInput) error); ok {
		r0 = rf(ctx, id, in)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
