# cards

A collection of cards and stacks

## Development notes

### gRPC

```plain
protoc -I cards-service/ cards-service/service.proto --go_out=plugins=grpc:cards-service
```

### MySQL

```plain
We've installed your MySQL database without a root password. To secure it run:
    mysql_secure_installation

MySQL is configured to only allow connections from localhost by default

To connect run:
    mysql -uroot

To have launchd start mysql now and restart at login:
  brew services start mysql
Or, if you don't want/need a background service you can just run:
  mysql.server start
```

### Postgres

```plain
If builds of PostgreSQL 9 are failing and you have version 8.x installed,
you may need to remove the previous version first. See:
  https://github.com/Homebrew/legacy-homebrew/issues/2510

To migrate existing data from a previous major version (pre-9.0) of PostgreSQL, see:
  https://www.postgresql.org/docs/9.6/static/upgrading.html

To migrate existing data from a previous minor version (9.0-9.5) of PostgreSQL, see:
  https://www.postgresql.org/docs/9.6/static/pgupgrade.html

  You will need your previous PostgreSQL installation from brew to perform `pg_upgrade`.
  Do not run `brew cleanup postgresql` until you have performed the migration.

To have launchd start postgresql now and restart at login:
  brew services start postgresql
Or, if you don't want/need a background service you can just run:
  pg_ctl -D /usr/local/var/postgres start
```
