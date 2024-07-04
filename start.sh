#!/bin/sh
# This file will be run by /bin/sh 
# We're using alpine image so there isn't bash shell
# We use set -e instruction to make sure that the script will exit immediately
# if a command returns a non-zero status
set -e

# We will take from the DB_SOURCE environment variable
echo "run db migration"
# It makes source app.env so we can use echo $DB_SOURCE
source /app/app.env
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
# takes all parameters passed to the script and run it
exec "$@"