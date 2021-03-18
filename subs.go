package main

import "context"

// Subscribe a user to a channel.
func (srv *Shoutyface) Subscribe(user, channel string) error {
	sql := "insert into subs(uid,cid) values((select id from users where name=$1),(select id from channels where name=$2));"
	_, err := srv.dbp.Exec(context.Background(), sql, user, channel)
	return err
}

// Unsubscribe user from channel.
func (srv *Shoutyface) Unsubscribe(user, channel string) error {
	sql := "delete from subs where uid in (select id from users where name=$1) and cid in (select id from channels where name=$2);"
	_, err := srv.dbp.Exec(context.Background(), sql, user, channel)
	return err
}
