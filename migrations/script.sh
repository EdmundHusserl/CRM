#!/bin/sh

psql -U postgres -h localhost -f /docker-entrypoint-initdb.d/000.sql