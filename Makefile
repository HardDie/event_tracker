.PHONY: swagger
swagger:
	swagger generate spec -m -o swagger.yaml

.PHONY: swagger-install
swagger-install:
	go install github.com/go-swagger/go-swagger/cmd/swagger@latest

.PHONY: build
build:
	CGO_ENABLED=0 go build -o event_tracker cmd/event_tracker/main.go
