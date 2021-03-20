package main

import (
	"encoding/json"
	"net/http"
)

// GetUsers from database.
func (srv *Shoutyface) GetUsers(w http.ResponseWriter, r *http.Request) {
	users := srv.listUsers()
	data, _ := json.Marshal(users)
	w.Write([]byte(data))
}

// PostUser adds a user to the database.
func (srv *Shoutyface) PostUser(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")
	email := r.Header.Get("email")
	if username == "" || email == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err := srv.addUser(username, email)
	if err != nil {
		srv.E("Error adding user: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Error(w, http.StatusText(http.StatusCreated), http.StatusCreated)
}

// DeleteUser from database.
func (srv *Shoutyface) DeleteUser(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")
	if username == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err := srv.rmUser(username)
	if err != nil {
		srv.E("Error deleting user: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
