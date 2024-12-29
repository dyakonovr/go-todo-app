package todosapi

type createTodoInput struct {
	Name string
}

type updateTodoInput struct {
	Name        string
	IsCompleted bool
}
