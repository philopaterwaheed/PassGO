build-web:
	GOOS=js GOARCH=wasm go build -o web/main.wasm cmd/frontend/main.go

build-windows:
	mkdir -p dist
	GOOS=windows GOARCH=amd64 go build -ldflags="-H=windowsgui" -o dist/passgo-frontend-windows.exe ./cmd/frontend
