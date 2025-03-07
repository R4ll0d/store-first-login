# store-first-login
env GOOS=linux go build -ldflags "-X main.Buildtime=$(date +%FT%T%z)" -o store-first-login main.go


docker build -t asia-southeast1-docker.pkg.dev/sgstore/cloud-run-source-deploy/store-first-login .

docker push asia-southeast1-docker.pkg.dev/sgstore/cloud-run-source-deploy/store-first-login