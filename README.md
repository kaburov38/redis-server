# redis-server
A simple Redis Server that can perform simple operations such as set and get.

## Prerequisite
1. Go version 1.25.4
2. [Redis](https://redis.io/docs/getting-started/installation/)

## How to Run
1. Disable Redis server by running `sudo systemctl stop redis`
2. Run the simple redis server by running `go run .`
3. run `redis-cli` to connect to the new redis server
