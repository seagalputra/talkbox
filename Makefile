gen-docs: fmt-docs
	swag init -g server.go

fmt-docs:
	swag fmt

run:
	go run .
