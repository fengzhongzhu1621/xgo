#!/bin/bash -
go build -o run_gohttpserver

cp -a run_gohttpserver ../

cd ..
./run_gohttpserver
