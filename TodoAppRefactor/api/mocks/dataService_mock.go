package apiMocks

import (
	"errors"
	"todoApp/data"
	"todoApp/utils/stringUtils"
)

type mockDataService struct{}

func NewMockDataService() *mockDataService {
	return &mockDataService{}
}

func (dataService *mockDataService) CreateTodoItem(name string) error {
	if stringUtils.IsEmptyOrWhitespace(name) {
		return errors.New("name cannot be empty")
	} else {
		return nil
	}
}

func (dataService *mockDataService) GetTodoItem(index int) (data.TodoItem, error) {
	switch index {
	case 0:
		todoItem := data.TodoItem{Name: "MockItem", Complete: false}
		return todoItem, nil
	default:
		return data.TodoItem{}, errors.New("item at specified index does not exist")
	}
}

func (dataService *mockDataService) GetAllTodoItems() []data.TodoItem {
	return []data.TodoItem{
		{Name: "TodoItem1", Complete: true},
		{Name: "TodoItem2", Complete: false},
		{Name: "TodoItem3", Complete: true},
	}
}

func (dataService *mockDataService) MarkItemAsComplete(index int) error {
	switch index {
	case 0:
		return nil
	default:
		return errors.New("item at specified index does not exist")
	}
}

func (dataService *mockDataService) DeleteTodoItem(index int) error {
	switch index {
	case 0:
		return nil
	default:
		return errors.New("item at specified index does not exist")
	}
}
