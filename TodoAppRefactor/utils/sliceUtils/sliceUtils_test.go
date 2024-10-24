package sliceUtils

import (
	"testing"
	"todoApp/data"
)

func TestTodoItemsEqual(t *testing.T) {
	var testCases = []struct {
		testName string
		sliceA   []data.TodoItem
		sliceB   []data.TodoItem
		expected bool
	}{
		{"Matching slices", []data.TodoItem{
			{Name: "TodoItem1", Complete: false},
			{Name: "TodoItem2", Complete: true},
		}, []data.TodoItem{
			{Name: "TodoItem1", Complete: false},
			{Name: "TodoItem2", Complete: true},
		}, true},

		{"Non-matching slices", []data.TodoItem{
			{Name: "TodoItem1", Complete: false},
			{Name: "TodoItem2", Complete: true},
		}, []data.TodoItem{
			{Name: "TodoItem1", Complete: true},
			{Name: "TodoItem2", Complete: false},
		}, false},

		{"Empty slices", []data.TodoItem{}, []data.TodoItem{}, true},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			if result := TodoItemsEqual(test.sliceA, test.sliceB); result != test.expected {
				t.Errorf("Unexpected result. Got: %v, Expected: %v", result, test.expected)
			}
		})
	}
}
