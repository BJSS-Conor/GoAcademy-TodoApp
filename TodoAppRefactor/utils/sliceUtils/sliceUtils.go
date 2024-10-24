package sliceUtils

import "todoApp/data"

func TodoItemsEqual(a, b []data.TodoItem) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
