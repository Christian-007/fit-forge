package web

import (
	"net/http"
)

func Routes(mux *http.ServeMux) {
	todoHandler := NewTodoHandler()

	mux.HandleFunc("GET /todos", todoHandler.GetAll)
}
