package main

import (
	"net/http"
)

func (server *APIServer) handleHealthCheck(w http.ResponseWriter, _ *http.Request) error {
	return WriteJSON(w, http.StatusOK, "ok")
}
