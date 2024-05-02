package main

import (
	"encoding/json"
	"net/http"
)

func (server *APIServer) handleCreateExpense(w http.ResponseWriter, r *http.Request) error {
	req := new(CreateExpenseRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	expense, err := NewExpense(
		req.Name,
		req.Price,
		req.Type,
		req.Description,
	)
	if err != nil {
		return err
	}

	err = server.store.CreateExpense(expense)
	if err != nil {
		return err
	}

	// Recovering expense from DB
	createdExpense, err := server.store.GetExpenseByID(expense.ID)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, createdExpense)
}

func (server *APIServer) handleGetExpenses(w http.ResponseWriter, _ *http.Request) error {
	expenses, err := server.store.GetExpenses()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, expenses)
}

func (server *APIServer) handleGetExpenseByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	expense, err := server.store.GetExpenseByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, expense)
}

func (server *APIServer) handleUpdateExpense(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	_, err = server.store.GetExpenseByID(id)
	if err != nil {
		return err
	}

	var expense Expense
	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		return err
	}

	expense.ID = id

	if err := server.store.UpdateExpense(&expense); err != nil {
		return err
	}

	// Retrieve the updated information from the database to get the most up-to-date data
	updatedExpense, err := server.store.GetExpenseByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, updatedExpense)
}

func (server *APIServer) handleDeleteExpense(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	_, err = server.store.GetExpenseByID(id)
	if err != nil {
		return err
	}

	if err := server.store.DeleteExpense(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"deleted": id})
}
