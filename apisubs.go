package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v4"
)

type Sub struct {
	User    string `json:"user"`
	Channel string `json:"channel"`
}

// GetSubscriptions lists all subscriptions for a user.
func (srv *Shoutyface) GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")
	var rows pgx.Rows
	var err error
	if username == "" {
		sql := "select u.name,c.name from subs s inner join users u on u.id=s.uid inner join channels c on c.id=s.cid and s.uid=u.id;"
		rows, err = srv.dbp.Query(context.Background(), sql)
	} else {
		sql := "select u.name,c.name from subs s inner join users u on u.name=$1 inner join channels c on c.id=s.cid and s.uid=u.id;"
		rows, err = srv.dbp.Query(context.Background(), sql, username)
	}

	if err != nil {
		srv.E("Error getting subscribers: %s", err.Error())
		http.Error(w, "", http.StatusNotFound)
		return
	}

	defer rows.Close()
	var subs []Sub
	for rows.Next() {
		s := Sub{}
		err = rows.Scan(&s.User, &s.Channel)
		if err != nil {
			srv.E("Error getting subscriptions: %s", err.Error())
		}

		subs = append(subs, s)
	}

	data, err := json.Marshal(subs)
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	w.Write(data)
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
		http.Error(w, err.Error(), http.StatusConflict)
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
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}
