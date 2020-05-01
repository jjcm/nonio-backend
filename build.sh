#!/bin/bash

# let's build the web server!
cd cmd

# this command will build a staticly linked binary for 64 bit linux systems
# and place it in the dist folder
echo "Building linux binary..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../dist/socid
echo "done!"

# maybe in the future there will be other things that we need to build...
export APP_KEY=superdupersecretyayyayyayyyyyyyyyy
export OAUTH_ID=0ed06b35279d956038d7
export OAUTH_SECRET=2af443a467821d992837d2ae1ca6af175f413af5
export DB_HOST=localhost
export DB_POST=3306
export DB_DATABASE=socidb
export DB_USER=root
export DB_PASSWORD=genius
export APP_PORT=8081

cd ../dist
killall socid || true
./socid &
