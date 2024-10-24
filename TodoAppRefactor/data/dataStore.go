package data

type TodoItem struct {
	Name     string
	Complete bool
}

var DataStore = []TodoItem{
	{Name: "Real Item 1", Complete: false},
	{Name: "Real Item 2", Complete: false},
	{Name: "Real Item 3", Complete: false},
}
