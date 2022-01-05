#!/bin/bash
./docker-build.sh && docker rm -f psyme && ./docker-run.sh