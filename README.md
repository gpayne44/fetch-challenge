# fetch-challenge

This service implements the API specification for the Fetch receipt processor challenge.

## Running the service
### Install dependencies
Navigate to project root
```
go mod download
```
This service uses `github.com/google/uuid` and `github.com/gorilla/mux`.
### Start the server with optional port flag
```
go run cmd/main.go
```
The default port is `:8000`. To specify another port use the `-port` flag:
```
go run cmd/main.go -port={port}
```

The server will be reachable at `localhost:{port}`.
