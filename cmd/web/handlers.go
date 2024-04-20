package main

import "net/http"

func userGetAll(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
