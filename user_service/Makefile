lint:
	golangci-lint run

build: clean $(APP_BIN)

$(APP_BIN):
	go build -o $(APP_BIN) ./app/cmd/app/main.go

clean:
	rm -rf ./app/build || true

swagger:
	swag init -g ./cmd/app/main.go -o ./docs
