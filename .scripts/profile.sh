BIN=$GOPATH/bin/gooster_profile

go build -o "$BIN" cmd/app/profile.go
PROFILE=$($BIN 2>&1 | grep "enabled" | grep -o -E "/[^ ]+cpu.pprof")

mkdir -p var
go tool pprof --pdf "$BIN" "$PROFILE" > var/profile.pdf
open var/profile.pdf
