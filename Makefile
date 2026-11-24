swagger-gen:
	swag init -g ./cmd/gourd/main.go --exclude internal/gateway/dto/registry -o ./docs --ot json