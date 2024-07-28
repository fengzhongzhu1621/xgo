#!/bin/bash -
go build -o run_gohttpserver main.go

cp -a run_gohttpserver ../

cd ..
./run_gohttpserver
