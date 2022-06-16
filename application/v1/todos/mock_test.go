package v1todos

import (
	"github.com/stretchr/testify/mock"

	"github.com/sioncojp/famili-api/domain"
	"github.com/sioncojp/famili-api/domain/model"
)

type MockTodoService struct {
	mock.Mock
}

type DomainMock struct {
	mock.Mock
}

func (m *MockTodoService) GetById(id domain.Id) (model.Todo, error) {
	r := m.Called(id)
	return r.Get(0).(model.Todo), r.Error(1)
}

func (m *MockTodoService) List() ([]model.Todo, error) {
	r := m.Called()
	return r.Get(0).([]model.Todo), r.Error(1)
}

func (m *MockTodoService) Create(todo *model.Todo) error {
	r := m.Called(todo)
	var r0 error
	if rf, ok := r.Get(0).(func(*model.Todo) error); ok {
		r0 = rf(todo)
	} else {
		r0 = r.Error(0)
	}
	return r0
}

func (m *MockTodoService) Update(todo *model.Todo) error {
	r := m.Called(todo)
	var r0 error
	if rf, ok := r.Get(0).(func(*model.Todo) error); ok {
		r0 = rf(todo)
	} else {
		r0 = r.Error(0)
	}
	return r0
}

func (m *MockTodoService) Delete(todo *model.Todo) error {
	r := m.Called(todo)
	var r0 error
	if rf, ok := r.Get(0).(func(*model.Todo) error); ok {
		r0 = rf(todo)
	} else {
		r0 = r.Error(0)
	}
	return r0
}
