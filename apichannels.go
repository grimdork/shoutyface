package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

func (srv *Shoutyface) PostChannel(w http.ResponseWriter, r *http.Request) {
	channel := r.Header.Get("channel")
	if channel == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	sql := "insert into channels(name) values($1)"
	_, err := srv.dbp.Exec(context.Background(), sql, channel)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			http.Error(w, "", http.StatusNotModified)
			return
		}

		srv.E("Error adding channel: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// DeleteChannel fom database.
func (srv *Shoutyface) DeleteChannel(w http.ResponseWriter, r *http.Request) {
	channel := r.Header.Get("channel")
	if channel == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ct, err := srv.dbp.Exec(context.Background(), "delete from channels where name=$1;", channel)
	if err != nil || ct.RowsAffected() == 0 {

		http.Error(w, "", http.StatusNotFound)
		return
	}
}

// Channel for JSON result.
type Channel struct {
	// Name of the channel.
	Name string `json:"channel"`
}

// GetChannels available.
func (srv *Shoutyface) GetChannels(w http.ResponseWriter, r *http.Request) {
	rows, err := srv.dbp.Query(context.Background(), "select name from channels;")
	if err != nil {
		srv.E("Error listing channels: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	var channels []Channel
	for rows.Next() {
		var c Channel
		rows.Scan(&c.Name)
		channels = append(channels, c)
	}

	data, _ := json.Marshal(channels)
	println(data)
	w.Write([]byte(data))
}
