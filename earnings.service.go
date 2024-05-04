package main

import "net/http"

func (server *APIServer) handleGetEarnings(w http.ResponseWriter, _ *http.Request) error {
	earnings, err := server.store.GetEarnings()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, earnings)
}
