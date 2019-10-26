build:
	go build -o ${GOPATH}/bin/gooster cmd/app/main.go

profile:
	bash ./.scripts/profile.sh

run:
	go run cmd/app/main.go -l=debug
