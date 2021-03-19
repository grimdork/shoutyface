package main

import (
	"context"
)

// addUser to database.
func (srv *Shoutyface) addUser(name, email string) error {
	sql := "insert into users(name,email) values($1,$2)"
	_, err := srv.dbp.Exec(context.Background(), sql, name, email)
	return err
}

// rmUser deletes a user from the database.
func (srv *Shoutyface) rmUser(name string) error {
	_, err := srv.dbp.Exec(context.Background(), "delete from users where name=$1 cascade;", name)
	return err
}

// listSubscribers retuns the e-mail addresses of all subscribers to a channel.
func (srv *Shoutyface) listSubscribers(channel string) []string {
	emails := []string{}
	sql := "select u.email from users u inner join channels c on c.name=$1 inner join subs s on s.cid=c.id and s.uid=u.id;"
	rows, err := srv.dbp.Query(context.Background(), sql, channel)
	if err != nil {
		srv.E("Error getting subscribers: %s", err.Error())
		return nil
	}

	var m string
	for rows.Next() {
		err = rows.Scan(&m)
		if err != nil {
			srv.E("Error getting subscribers: %s", err.Error())
			return emails
		}

		emails = append(emails, m)
	}
	return emails
}
