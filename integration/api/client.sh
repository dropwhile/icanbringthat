#!/bin/sh

./build/bin/client events list
./build/bin/client earmarks list
./build/bin/client favorites list
./build/bin/client notifications list

./build/bin/client events detail --ref-id 0662e3hdgwf05n7hc0jzsbhwbw
./build/bin/client earmarks detail --ref-id 0662e7ha4cr08jpaf483jcdvwr
