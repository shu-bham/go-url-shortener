#!/bin/bash
./bin/migrate -database "mysql://root:root@tcp(127.0.0.1:3306)/url_shortener_db?parseTime=true" -source "file://db/migrations" "$@"
