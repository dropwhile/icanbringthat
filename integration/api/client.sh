#!/bin/sh

echo "list events"
./build/bin/client events list
echo

echo "list earmarks"
./build/bin/client earmarks list
echo

echo "list favorites"
./build/bin/client favorites list
echo

echo "list notifications"
./build/bin/client notifications list
echo

echo "event detail"
./build/bin/client events detail 0662e3hdgwf05n7hc0jzsbhwbw
echo

echo "earmark detail"
./build/bin/client earmarks detail 0662e7ha4cr08jpaf483jcdvwr
echo
