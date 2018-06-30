# pgterminate
> Terminates active and idle PostgreSQL backends

Did you encountered long-running queries locking down your entire company system because of a massive lock on the database? Or a scatterbrained developper connected to the production database with an open transaction leading to a production outage? Never? Really? Either you have very good policies and that's awesome, or you don't work in databases at all.

With `pgterminate`, you shouldn't be paged at night because some queries has locked down production for too long. It looks after "active" and "idle" connections and terminate them. As simple as that.

# Highlights
* `pgterminate` name is derived from `pg_terminate_backend` function, it terminates backends.
* backends are called sessions in `pgterminate`.
* use `cancel` option to terminate query instead of session.
* `active` sessions are backends in `active` state for more than `active-timeout` seconds.
* `idle` sessions are backends in `idle`, `idle in transaction` or `idle in transaction (abort)` state for more than `idle-timeout` seconds.
* at least one of `active-timeout` and `idle-timeout` parameter is required, both can be used.
* `pgterminate` relies on `libpq` for PostgreSQL connection. When `host` is ommited, connection via unix socket is used. When `user` is ommited, the unix user is used. And so on.
* time parameters, like `connect-timeout`, `active-timeout`, `idle-timeout` and `interval`, are represented in seconds. They accept float value except for `connect-timeout` which is an integer.
* if you want `pgterminate` to terminate any session, ensure it has SUPERUSER privileges. Since 9.6, grant `pg_signal_backend` role for terminating all sessions except superusers.

# Internals

## Signals
`pgterminate` handles the following OS signals:
* `SIGINT`, `SIGTERM` to gracefully terminates the infinite loop
* `SIGHUP` to reload configuration file and re-open log file if used (handy for logrotate)

## Configuration
There's two ways to configure `pgterminate`:
* command-line arguments
* configuration file with `-config` command-line argument

Configuration file options **override** command-line arguments

# Usage
Connect to a remote instance and prompt for password:
```
pgterminate -host 10.0.0.1 -port 5432 -user test -prompt-password -database test
```
Use a configuration file:
```
pgterminate -config config.yaml
```
Use both configuration file and command-line arguments:
```
pgterminate -config config.yaml -interval 0.25 -active-timeout 10 -idle-timeout 300
```
Print usage:
```
pgterminate -help
```

# Filtering users

`pgterminate` is able to include or exclude users from being terminated.

## Configuration
### List
Arguments `-include-user` or `-exclude-user` can be used multiple times for multiple users:

```
pgterminate -include-user user1 -include-user user2
```
Or in configuration file:

```
include-users:
  user1
  user2
```
Same applies for `-exclude-user` (argument) and `exclude-users` (file).

### Regexes
Regexes can be configured:

```
pgterminate -include-users-regex "(user1|user2)"
```
Or in configuration file:

```
include-users-regex: "(user1|user2)"
```

Same applies for `-exclude-users-regex` (argument) and `exclude-users-regex` (file).

## Include users

When include users list or regex is set, `pgterminate` will focus on included users only. It could terminate excluded users if any. If you want to exclude users, use exclude options only.

## Exclude users

When exclude users list or regex is set and no include option is set, `pgterminate` will terminate all sessions except excluded users.

# License
`pgterminate` is released under [The Unlicense](LICENSE) license. Code is under public domain.
