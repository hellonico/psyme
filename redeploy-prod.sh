#!/bin/bash
./exportdb.sh
git pull origin main && ./docker-build.sh && docker rm -f psyme && ./docker-run.sh
./importdb.sh
./sync-articles.sh
