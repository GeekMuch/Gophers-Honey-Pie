#!/bin/bash

# Gets latest update from Github
echo "[ + ] Git pulling Gophers-Honey-pie"
git pull
echo "[ ! ] Done git pulling Gophers-Honey-pie"

# Update Go Modules
echo "[ + ] Updating GO modules"
cd /home/pi/Gophers-Honey-Pie || exit
# /usr/local/go/bin/go get -u
/usr/local/go/bin/go mod tidy
echo "[ ! ] Done updating GO modules"

#Start main.go
echo "[ + ] Starting Gophers-Honey-pie"
/usr/local/go/bin/go run main.go
