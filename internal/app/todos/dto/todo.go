package dto

type TodoResponse struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	IsCompleted bool   `json:"isCompleted"`
}

type CreateTodoRequest struct {
	Title       string `json:"title"`
}
