# Wren Server (implemented in Go)

`go run .` will start a wren server listening on 0.0.0.0:8985
`go test ./...` will test every package
`go test ./server -race -covermode=atomic -coverprofile=coverage.out` will test and output coverage
`go tool cover -html=coverage.out` will open a browser showing coverage information