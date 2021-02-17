#!/bin/sh
# wait-for-postgres.sh
# exit on error
set -e

until PGPASSWORD=$POSTGRES_PASSWORD PGHOST=$POSTGRES_HOST PGUSER=$POSTGRES_USER psql -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 5
done

>&2 echo "Postgres is up - executing command"
# mb isn't the best way to run
./go/bin/mx-api-trainee