package main

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgtype"
)

// Token is an access token. Required by clients to send messages.
type Token struct {
	Hash    string
	User    string
	Email   string
	Expires time.Time
}

func (srv *Shoutyface) tokencheck(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")
		if token == "" || srv.GetTokenUser(token) == "" {
			http.Error(w, "", http.StatusForbidden)
			return
		}

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (srv *Shoutyface) admintokencheck(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")
		if token == "" || srv.GetTokenUser(token) != "admin" {
			http.Error(w, "", http.StatusForbidden)
			return
		}

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// GetTokenUser or return an empty string if the token doesn't exist.
func (srv *Shoutyface) GetTokenUser(t string) string {
	row := srv.dbp.QueryRow(context.Background(), "select name from users inner join tokens on tokens.uid=users.id where tokens.hash=$1;", t)
	var name string
	err := row.Scan(&name)
	if err != nil {
		return ""
	}

	return name
}

// AddToken for a user.
func (srv *Shoutyface) AddToken(token, user string) error {
	now := pgtype.Timestamp{
		Time:   time.Now().UTC(),
		Status: pgtype.Present,
	}
	sql := "insert into tokens(hash,uid,expires) select $1,u.id,$3 from users u where u.name=$2;"
	_, err := srv.dbp.Exec(context.Background(), sql, token, user, now)
	return err
}

// DeleteToken from database.
func (srv *Shoutyface) DeleteToken(token string) {
	srv.dbp.Exec(context.Background(), "delete from tokens where hash=$1;", token)
}
