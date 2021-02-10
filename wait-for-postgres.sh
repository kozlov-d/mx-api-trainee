#!/bin/sh
# wait-for-postgres.sh

set -e

until PGPASSWORD=$POSTGRES_PASSWORD PGHOST=$POSTGRES_HOST PGUSER=$POSTGRES_USER psql -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 5
done

>&2 echo "Postgres is up - executing command"