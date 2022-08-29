#!/bin/sh

cd /app
exec gins_example -ip=0.0.0.0 -port=8080 $@