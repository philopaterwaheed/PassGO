build-web:
	GOOS=js GOARCH=wasm go build -o web/main.wasm cmd/frontend/main.go
