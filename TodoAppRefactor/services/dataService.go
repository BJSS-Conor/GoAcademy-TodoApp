package dataService

import (
	"errors"
	"sync"
	"todoApp/data"
	"todoApp/utils/stringUtils"
)

type IDataService interface {
	CreateTodoItem(name string) error
	GetTodoItem(index int) (data.TodoItem, error)
	GetAllTodoItems() []data.TodoItem
	MarkItemAsComplete(index int) error
	DeleteTodoItem(index int) error
}

type DataService struct {
	data []data.TodoItem
	mu   sync.RWMutex
}

func NewDataService() *DataService {
	return &DataService{
		data: data.DataStore,
	}
}

func (dataService *DataService) CreateTodoItem(name string) error {
	if stringUtils.IsEmptyOrWhitespace(name) {
		return errors.New("name cannot be empty")
	} else {
		dataService.mu.Lock()
		defer dataService.mu.Unlock()

		todoItem := data.TodoItem{Name: name, Complete: false}
		dataService.data = append(dataService.data, todoItem)
		return nil
	}
}

func (dataService *DataService) GetTodoItem(index int) (data.TodoItem, error) {
	if !(index >= 0 && index < len(dataService.data)) {
		return data.TodoItem{}, errors.New("item at specified index does not exist")
	}

	dataService.mu.RLock()
	defer dataService.mu.RUnlock()

	todoItem := dataService.data[index]
	return todoItem, nil
}

func (dataService *DataService) GetAllTodoItems() []data.TodoItem {
	dataService.mu.RLock()
	defer dataService.mu.RUnlock()

	return dataService.data
}

func (dataService *DataService) MarkItemAsComplete(index int) error {
	if !(index >= 0 && index < len(dataService.data)) {
		return errors.New("item at specified index does not exist")
	} else {
		dataService.mu.Lock()
		defer dataService.mu.Unlock()

		dataService.data[index].Complete = true
		return nil
	}
}

func (dataService *DataService) DeleteTodoItem(index int) error {
	if !(index >= 0 && index < len(dataService.data)) {
		return errors.New("item at specified index does not exist")
	} else {
		dataService.mu.Lock()
		defer dataService.mu.Unlock()

		dataService.data = append(dataService.data[:index], dataService.data[index+1:]...)
		return nil
	}
}
