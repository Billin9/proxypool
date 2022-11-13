#!/bin/bash

>/var/lib/docker/containers/`docker inspect --format="{{.Id}}" proxypool`/`docker inspect --format="{{.Id}}" proxypool`-json.log
