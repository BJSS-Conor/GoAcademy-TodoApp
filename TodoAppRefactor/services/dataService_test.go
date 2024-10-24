package dataService

import (
	"testing"
	"todoApp/data"
	sliceUtils "todoApp/utils/sliceUtils"
)

func CreateTestData(testDataType int) *DataService {
	switch testDataType {
	case 1:
		return &DataService{
			data: []data.TodoItem{
				{Name: "TodoItem1", Complete: false},
				{Name: "TodoItem2", Complete: false},
				{Name: "TodoItem3", Complete: false},
			},
		}
	default:
		return &DataService{
			data: []data.TodoItem{},
		}
	}
}

func TestCreateTodoItem(t *testing.T) {
	inputName := "Test"
	expectedItem := data.TodoItem{Name: "Test", Complete: false}
	dataService := CreateTestData(1)

	if err := dataService.CreateTodoItem(inputName); err != nil {
		t.Errorf("An unexpected error occured whilst creating the todo item: %s", err.Error())
	} else if item, getErr := dataService.GetTodoItem(3); getErr != nil {
		t.Errorf("An unexpected error occured whilst checking the newly created item exists: %s", getErr.Error())
	} else if item != expectedItem {
		t.Errorf("Todo item was not created correctly. Got %v, Expected %v", item, expectedItem)
	}
}

func TestCreateTodoItem_EmptyName(t *testing.T) {
	testcases := []struct {
		testName  string
		inputName string
	}{
		{"Testing empty string", ""},
		{"Testing string containing only whitespaces", "   "},
	}
	expectedError := "name cannot be empty"
	dataService := CreateTestData(1)

	for _, test := range testcases {
		t.Run(test.testName, func(t *testing.T) {
			if err := dataService.CreateTodoItem(test.inputName); err == nil {
				t.Error("An invalid name was entered, an error was expected but not recieved")
			} else if err.Error() != expectedError {
				t.Errorf("An invalid name was entered, an error occured but not the expected one. Got: %s, Expected: %s", err.Error(), expectedError)
			}
		})
	}
}

func TestGetTodoItem_ValidIndex(t *testing.T) {
	testCases := []struct {
		testName     string
		inputIndex   int
		expectedItem data.TodoItem
	}{
		{"Testing with index on lower bound", 0, data.TodoItem{Name: "TodoItem1", Complete: false}},
		{"Testing with index on higher bound", 2, data.TodoItem{Name: "TodoItem3", Complete: false}},
	}
	dataService := CreateTestData(1)

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			if item, err := dataService.GetTodoItem(test.inputIndex); err != nil {
				t.Errorf("An unexpected error occured: %s", err.Error())
				t.Errorf("Data Size: %d", len(dataService.data))
			} else if item != test.expectedItem {
				t.Errorf("The incorrect item was returned. Got: %v, Expected: %v", item, test.expectedItem)
			}
		})
	}
}

func TestGetTodoItem_InvalidIndex(t *testing.T) {
	testCases := []struct {
		testName      string
		inputIndex    int
		expectedError string
	}{
		{"Testing with index below 0", -1, "item at specified index does not exist"},
		{"Testing with index higher than data length", 3, "item at specified index does not exist"},
	}
	dataService := CreateTestData(1)

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			if _, err := dataService.GetTodoItem(test.inputIndex); err == nil {
				t.Error("Index is out of range so an error was expected but not recieved")
			} else if err.Error() != test.expectedError {
				t.Errorf("An occured but not the expected error. Got: %s, Expected: %s", err.Error(), test.expectedError)
			}
		})
	}
}

func TestGetAllTodoItems(t *testing.T) {
	expectedItems := CreateTestData(1).data
	dataService := CreateTestData(1)

	items := dataService.GetAllTodoItems()
	if !sliceUtils.TodoItemsEqual(items, expectedItems) {
		t.Errorf("The returned list of items does not match the expected list of items. Got: %v, Expected %v", items, expectedItems)
	}
}

func TestGetAllTodoItems_EmptyList(t *testing.T) {
	expectedItems := CreateTestData(0).data
	dataService := CreateTestData(0)

	items := dataService.GetAllTodoItems()
	if !sliceUtils.TodoItemsEqual(items, expectedItems) {
		t.Errorf("The returned list of items does not match the expected list of items. Got: %v, Expected %v", items, expectedItems)
	}
}

func TestMarkItemAsComplete(t *testing.T) {
	inputIndex := 0
	dataService := CreateTestData(1)

	if err := dataService.MarkItemAsComplete(inputIndex); err != nil {
		t.Errorf("An unexpected error occured: %s", err.Error())
	} else if updatedItem, err := dataService.GetTodoItem(inputIndex); err != nil {
		t.Errorf("An unexpected error occured whilst trying to obtain the updated item: %s", err.Error())
	} else if updatedItem.Complete != true {
		t.Error("The item was not correctly marked as complete. Still marked as incomplete in the data")
	}
}

func TestMarkItemAsComplete_InvalidIndex(t *testing.T) {
	inputIndex := 10
	expectedError := "item at specified index does not exist"
	dataService := CreateTestData(1)

	if err := dataService.MarkItemAsComplete(inputIndex); err == nil {
		t.Error("The specified index is invalid but an error was not produced")
	} else if err.Error() != expectedError {
		t.Errorf("The specified index is invalid but the error produced is unexpected. Got: %s, Expected: %s", err.Error(), expectedError)
	}
}

func TestDeleteTodoItem_ValidIndex(t *testing.T) {
	inputIndex := 1
	expectedItems := []data.TodoItem{
		{Name: "TodoItem1", Complete: false},
		{Name: "TodoItem3", Complete: false},
	}
	dataService := CreateTestData(1)

	if err := dataService.DeleteTodoItem(inputIndex); err != nil {
		t.Errorf("An unexpected error occured: %s", err.Error())
	} else if items := dataService.GetAllTodoItems(); !sliceUtils.TodoItemsEqual(items, expectedItems) {
		t.Errorf("The data does match what is expected after deleting the specified item. Got: %v, Expected: %v", items, expectedItems)
	}
}

func TestDeleteTodoItem_InvalidIndex(t *testing.T) {
	testCases := []struct {
		testName      string
		inputIndex    int
		expectedError string
	}{
		{"Testing with index below 0", -1, "item at specified index does not exist"},
		{"Testing with index higher than data length", 100, "item at specified index does not exist"},
	}
	dataService := CreateTestData(1)

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			if err := dataService.DeleteTodoItem(test.inputIndex); err == nil {
				t.Error("Index is out of range so an error was expected but not recieved")
			} else if err.Error() != test.expectedError {
				t.Errorf("An occured but not the expected error. Got: %s, Expected: %s", err.Error(), test.expectedError)
			}
		})
	}
}
