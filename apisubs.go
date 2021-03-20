package main

import "net/http"

// GetSubscriptions lists all subscriptions for a user.
func (srv *Shoutyface) GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")
	if username == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

// PostSub subscribes a user to a channel.
func (srv *Shoutyface) PostSubscribe(w http.ResponseWriter, r *http.Request) {
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

// DeleteSubscription unsubscribes a user from a channel.
func (srv *Shoutyface) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
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
