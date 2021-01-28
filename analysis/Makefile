all:
	mkdir -p bin
	(cd cmd; go vet)
	(cd cmd; go get)
	go build -o bin/analyze-cp cmd/analyze-cp.go 