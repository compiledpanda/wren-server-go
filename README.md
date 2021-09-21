# Wren Server (implemented in Go)

An Implementation of https://github.com/compiledpanda/wren/blob/main/server.yml

* `go run .` will start a wren server listening on 0.0.0.0:8985
* `go test ./...` will test every package
* `go test ./server -race -covermode=atomic -coverprofile=coverage.out` will test and output coverage
* `go tool cover -html=coverage.out` will open a browser showing coverage information
* `golangci-lint run` will lint all code

# BoltDB Layout

```
Note: # prefix denotes a bucket

metadata
#shard
  #<shard id>
    <commit id>: commit
#shard_metadata
  <shard id>: metadata
#user
  <user id>: metadata
#user_key
  <user id & key id>: public key
#role
  <role id>: metadata
#role_user
  <role id & user id>: null
  <user id & role id>: null
#role_permission
  <role id & permission id>: null
```