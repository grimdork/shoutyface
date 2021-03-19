package main

type CmdEnv struct{}

const envvars = `Setup variables:
- ADMIN_TOKEN: Initial token for the admin user to use with the management API.
- ADMIN_EMAIL: E-mail to set for the admin user when creating the tables.

PostgreSQL variables:
- PGHOST: Database address. Defaults to localhost.
- PGUSER: Database user. No default.
- PGPASSWORD: Database password. No default.
Default database is shoutyface.

Mail variables:
- MAILHOST: SMTP server address. No default.
- MAILPORT: SMTP port. Default 587.
- MAILUSER: SMTP username.
- MAILPASS: SMTP password.
- MAILFROM: E-mail address to send alerts from.`

func (cmd *CmdEnv) Run(in []string) error {
	println(envvars)
	return nil
}
