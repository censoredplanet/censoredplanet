all:
	mkdir -p bin
	go get -d ./cmd/
	go install cmd/analyze-cp.go
	go build -o bin/analyze-cp cmd/analyze-cp.go
