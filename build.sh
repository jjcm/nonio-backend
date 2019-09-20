#!/bin/bash

# let's build the web server!
cd httpd

# this command will build a staticly linked binary for 64 bit linux systems
# and place it in the dist folder
echo "Building linux binary..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../dist/socid
echo "done!"
echo "Building OSX binary..."
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ../dist/socid-osx
echo "done!"



# maybe in the future there will be other things that we need to build...
