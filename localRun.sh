cd cmd
go build -o ../dist/socid
export APP_KEY="asdfa323faefjifajwiefawef"
export WEB_HOST="http://localhost:4200"
export DB_HOST="127.0.0.1"
export DB_PORT="3306"
export DB_DATABASE="socidb"
export DB_USER="dbuser"
export DB_PASSWORD="password"
export APP_PORT="4201"
export ADMIN_EMAIL="test@example.com"
export ADMIN_EMAIL_PASSWORD="password"
export EMAIL_ACCESS_TOKEN=""
export EMAIL_REFRESH_TOKEN=""
export EMAIL_CLIENT_ID=""
export EMAIL_CLIENT_SECRET=""
export STRIPE_KEY=""
export STRIPE_SECRET_KEY="asdf"
export STRIPE_PUBLISHABLE_KEY="asdf"
export WEBHOOK_ENDPOINT_SECRET=""

cd ../migrations
goose mysql "${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_DATABASE}" up

../dist/socid

