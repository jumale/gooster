build:
	go build -o ${GOPATH}/bin/gooster cmd/app/main.go

run:
	go run cmd/app/main.go

profile:
	bash ./.scripts/profile.sh

test:
	@go fmt ./pkg/... || echo FAILED
	@go test -bench=. ./pkg/... || echo FAILED

debug:
	dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient cmd/app/main.go
