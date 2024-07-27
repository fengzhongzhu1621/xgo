#!/bin/bash -
go build main.go

cp gohttpserver ../run_gohttpserver

cd ..
./run_gohttpserver
