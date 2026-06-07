build-web:
	GOOS=js GOARCH=wasm go build -o web/main.wasm cmd/frontend/main.go

build-windows:
	mkdir -p dist
	go install gioui.org/cmd/gogio@v0.10.0
	gogio -target windows -appid com.philopaterwaheed.passgo -name passGo -version 1.0.0.1 -o dist/passGo.exe ./cmd/frontend
