run:
	go mod download
ifeq ($(port),)
	go run cmd/main.go
else
	go run cmd/main.go -port=$(port)
endif