#!/bin/bash
git pull origin main && ./docker-build.sh && docker rm -f psyme && ./docker-run.sh