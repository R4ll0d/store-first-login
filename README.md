# store-first-login
env GOOS=linux go build -ldflags "-X main.Buildtime=$(date +%FT%T%z)" -o store-first-login main.go