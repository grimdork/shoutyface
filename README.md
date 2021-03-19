# Shoutyface
A mail alerter.

## Environment variables
Basic variables during setup:
- ADMIN_TOKEN: Initial token for the admin user to use with the management API.
- ADMIN_EMAIL: E-mail to set for the admin user when creating the tables.

Shoutyface tries to connect to a Postgres database named ´shoutyface`, and supports standard Postgres variables:
- PGHOST: Database address. Defaults to localhost.
- PGUSER: Database user. No default.
- PGPASSWORD: Database password. No default.

For mail, these variables need to be defined:
- MAILHOST: SMTP server address. No default.
- MAILPORT: SMTP port. Default 587.
- MAILUSER: SMTP username.
- MAILPASS: SMTP password.
- MAILFROM: E-mail address to send alerts from.
