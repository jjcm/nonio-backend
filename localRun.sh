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
<<<<<<< HEAD
export STRIPE_SECRET_KEY="sk_test_51EpA4oH4gvdXgbs5rBv4JI29C38uWuNEGuB8Agt5hfya1fjgVGOQePyfj7x6ANDPE7hyYNZEMRWwkP93NAa7QTCl00GPr79F0w"
export STRIPE_PUBLISHABLE_KEY="pk_test_51EpA4oH4gvdXgbs5r0aq0i3U6IzOwbWRVYaBYXMFLLHvihVHGHotHPAi2EJ7Km9JqudFZyLE30kt2YQSUOSK88Xx00Q6eEqxmS"
=======
export STRIPE_KEY=""
export STRIPE_SECRET_KEY="asdf"
export STRIPE_PUBLISHABLE_KEY="asdf"
>>>>>>> 88fff95fa8b7210fd0ecfd9936a2db6da184dd5a
export WEBHOOK_ENDPOINT_SECRET=""

cd ../migrations
goose mysql "${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_DATABASE}" up

../dist/socid

