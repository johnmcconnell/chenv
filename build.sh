set -e

FULL_BUILD=true

if [[ $# -gt 0 ]]
then
  if [[ "quick" = $1 ]]
  then
    FULL_BUILD=false
  fi
fi

if [[ "$FULL_BUILD" = true ]]
then
  gofmt -w ./
  golint ./...
  go vet ./...

  (go test ./... | grep -v ^ok | grep -v '^?') || true
fi

go build ./...
