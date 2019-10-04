build:
	go build -o ${GOPATH}/bin/gooster cmd/app/main.go

run:
	go run cmd/app/main.go -l=debug
