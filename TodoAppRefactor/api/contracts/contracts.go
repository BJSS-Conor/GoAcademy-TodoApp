package contracts

import "todoApp/data"

type CreateContract struct {
	Name string
}

type GetContract struct {
	Name     string
	Complete bool
}

type GetAllContract struct {
	TodoItems []data.TodoItem
}

type MarkItemAsCompleteContract struct {
	Id int
}
