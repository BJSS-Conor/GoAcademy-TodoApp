package responses

import "todoApp/data"

type CreateRes struct {
	Error error
}

type GetRes struct {
	Item  data.TodoItem
	Error error
}

type GetAllRes struct {
	Items []data.TodoItem
}

type MarkAsCompleteRes struct {
	Error error
}

type DeleteRes struct {
	Error error
}
