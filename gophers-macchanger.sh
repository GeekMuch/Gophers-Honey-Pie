#!/bin/bash

# Gets latest update from Github
echo "[ + ] Get mac from /boot/config.yml"
MAC=$(cat /boot/config.yml | grep -Eo "\w+:\w+:\w+:\w+:\w+:\w+")
echo "[ ! ] Changing MAC Address"
macchanger --mac $MAC eth0
asdasd