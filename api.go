package main

import (
	"net/http"

	"github.com/francoispqt/gojay"
)

// PostMessage to send to recipients.
func (srv *Shoutyface) PostMessage(w http.ResponseWriter, r *http.Request) {
	msg := Message{}
	dec := gojay.NewDecoder(r.Body)
	err := dec.DecodeObject(&msg)
	if err != nil {
		srv.E("Error decoding JSON: %s", err.Error())
		return
	}

	srv.mailqueue <- msg
}

//
// management endpoints
// All of these calls require an amin token.
//

// GetUsers from database.
func (srv *Shoutyface) GetUsers(w http.ResponseWriter, r *http.Request) {
	if !srv.isAdmin(r) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
}

// PostUser adds a user to the database.
func (srv *Shoutyface) PostUser(w http.ResponseWriter, r *http.Request) {
	if !srv.isAdmin(r) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	username := r.Header.Get("username")
	email := r.Header.Get("email")
	if username == "" || email == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err := srv.AddUser(username, email)
	if err != nil {
		srv.E("Error adding user: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Error(w, http.StatusText(http.StatusCreated), http.StatusCreated)
}

// DeleteUser from database.
func (srv *Shoutyface) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if !srv.isAdmin(r) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	username := r.Header.Get("username")
	if username == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err := srv.deleteUser(username)
	if err != nil {
		srv.E("Error deleting user: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Error(w, http.StatusText(http.StatusAccepted), http.StatusAccepted)
}

// GetSubs lists all subscriptions for a user.
func (srv *Shoutyface) GetSubs(w http.ResponseWriter, r *http.Request) {
	if !srv.isAdmin(r) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	username := r.Header.Get("username")
	if username == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

// PostSub subscribes a user to a channel.
func (srv *Shoutyface) PostSub(w http.ResponseWriter, r *http.Request) {
	if !srv.isAdmin(r) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	username := r.Header.Get("username")
	channel := r.Header.Get("channel")
	if username == "" || channel == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err := srv.Subscribe(username, channel)
	if err != nil {
		srv.E("Error subscribing user: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Error(w, http.StatusText(http.StatusCreated), http.StatusCreated)
}

// DeleteSub unsubscribes a user from a channel.
func (srv *Shoutyface) DeleteSub(w http.ResponseWriter, r *http.Request) {
	if !srv.isAdmin(r) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	username := r.Header.Get("username")
	channel := r.Header.Get("channel")
	if username == "" || channel == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err := srv.Unsubscribe(username, channel)
	if err != nil {
		srv.E("Error unsubscribing user: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Error(w, http.StatusText(http.StatusAccepted), http.StatusAccepted)
}
