package main

import (
	"context"
)

// AddUser to alert system.
func (srv *Shoutyface) AddUser(name, email string) error {
	sql := "insert into users(name,email) values($1,$2)"
	_, err := srv.dbp.Exec(context.Background(), sql, name, email)
	return err
}

// deleteUser from database.
func (srv *Shoutyface) deleteUser(name string) error {
	_, err := srv.dbp.Exec(context.Background(), "delete from users where name=$1 cascade;", name)
	return err
}

// GetSubscribers retuns the e-mail addresses of all subscribers to a channel.
func (srv *Shoutyface) GetSubscribers(channel string) []string {
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
