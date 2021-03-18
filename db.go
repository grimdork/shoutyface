package main

import (
	"context"
	"os"
)

// TABLES
// users - recipients of alerts
// tokens - session tokens and their owners
// channels - message channels to listen and send to
// subs - user->channel relations to determine where messages go

var setupSQL string = `begin;

create table if not exists users (
id serial primary key,
name text unique not null,
email text unique not null
);

create table if not exists tokens (
hash text primary key,
uid int,
expires timestamp,
constraint fk_tokens_uid foreign key(uid) references users(id) on delete cascade
);

create table if not exists channels (
id bigserial primary key,
name text unique not null
);

insert into channels(name) values('admin');

create table if not exists subs (
id bigint serial primary key,
uid int,
cid int,
constraint fk_subs_uid foreign key(uid) references users(id) on delete cascade,
constraint fk_subs_cid foreign key(cid) references channels(id) on delete cascade
);

create index idx_subs on subs(uid,cid);
commit;
`

func (srv *Shoutyface) createTables() error {
	srv.L("Initialising database.")
	_, err := srv.dbp.Exec(context.Background(), setupSQL)
	if err != nil {
		return err
	}

	admintoken := os.Getenv("ADMIN_TOKEN")
	if admintoken != "" {
		println(admintoken)
		err = srv.AddUser("admin", os.Getenv("ADMIN_EMAIL"))
		if err != nil {
			return err
		}

		err = srv.AddToken(admintoken, "admin")
		if err != nil {
			return err
		}
	}

	return nil
}
