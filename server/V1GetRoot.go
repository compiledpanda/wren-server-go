package server

import "net/http"

type RepositoryStatus struct {
	Status string `json:"status"`
}

func V1GetRoot(w http.ResponseWriter, r *http.Request) {
	ReturnJSON(w, http.StatusOK, RepositoryStatus{"ONLINE"})
}
