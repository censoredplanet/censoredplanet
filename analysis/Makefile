all:
	mkdir -p bin
	(cd cmd; go get)
	(cd cmd; go vet)
	go build -o bin/analyze-cp cmd/analyze-cp.go 