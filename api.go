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
