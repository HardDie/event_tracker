.PHONY: swagger
swagger:
	swagger generate spec -m -o swagger.yaml

.PHONY: swagger-install
swagger-install:
	go install github.com/go-swagger/go-swagger/cmd/swagger@latest
