#!/bin/sh

reset

# rm debug

go build -o debug

./debug $@

rm debug