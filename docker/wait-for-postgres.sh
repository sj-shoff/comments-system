#!/bin/sh

set -e

if [ "$STORAGE_TYPE" != "postgres" ]; then
  >&2 echo "Skipping PostgreSQL wait for $STORAGE_TYPE storage"
  exec $@
fi

host="$1"
port="$2"
shift 2
cmd="$@"

until PGPASSWORD=$POSTGRES_PASSWORD psql -h "$host" -p "$port" -U "postgres" -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - executing command"
exec $cmd