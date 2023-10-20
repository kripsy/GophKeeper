# GophKeeper
Yandex graduation project on the topic of secure secret storage manager


# For easy start project (server, deps, without client).
# install deps.
1) go mod download
# run postgres, pgadmin.
# password for postgres and pgadmin look in docker-compose file.
2) docker-compose up postgres_db pgadmin

# -secure=false use for disable tls. 
# It's helpful for dev, if u want use BloomRPC for testing. BloomRPC can't work with selfsigned certs.
3) go run ./cmd/server/main.go -secure=false