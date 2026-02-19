#!/bin/sh
exec gunicorn --bind 0.0.0.0:9090 main:app
