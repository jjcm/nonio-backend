cd cmd
go build -o ../dist/socid
export APP_KEY="asdfa323faefjifajwiefawef"
export OAUTH_ID="0ed06b35279d956038d7"
export OAUTH_SECRET="2af443a467821d992837d2ae1ca6af175f413af5"
export DB_HOST="127.0.0.1"
export DB_PORT="3306"
export DB_DATABASE="socidb"
export DB_USER="dbuser"
export DB_PASSWORD="password"
export APP_PORT="4201"

cd ../migrations
goose mysql "${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_DATABASE}" up

../dist/socid

