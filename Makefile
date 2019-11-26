build:
	go build -o ${GOPATH}/bin/gooster cmd/app/main.go

run:
	go run cmd/app/main.go -l=debug

profile:
	bash ./.scripts/profile.sh

test:
	go fmt ./pkg/...
	go test -bench=. ./pkg/...


debug:
	dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient cmd/app/main.go
