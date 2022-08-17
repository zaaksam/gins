#!/bin/sh

reset

rm debug

go build -o debug

# export GIN_MODE=release

./debug -ip=localhost -port=4040 -timeout=2 -pprof -debug

rm debug