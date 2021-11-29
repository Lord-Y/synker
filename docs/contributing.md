# Contributing

In order to contribute to our project, we will define here all requirements.

## Git hooks

You need first do install [golangci-lint](https://golangci-lint.run/usage/install/)

Then, enable the hook in our project:
```bash
git config core.hooksPath .githooks
```

## Start the clusters

Set `sysctl` values permanently for elasticsearch if not already done in your user environment:
```bash
cat <<EOF | sudo tee -a /etc/sysctl.d/10-custom.conf
vm.max_map_count=262144
EOF
sudo sysctl -p /etc/sysctl.d/10-custom.conf
```
Start the clusters:
```bash
sudo docker-compose -f docker/docker-compose-cluster.yaml up -d
```

## Start CockroachDB cluster

To start the cluster execute this command:
```bash
cockroach demo movr --geo-partitioned-replicas --insecure --http-port 18080
```

## Set default variables

```bash
export COCKROACH_HOST=$(netstat -latn |grep 26257 |grep LISTEN |awk '{print $4}')
export SKR_PG_URI="postgres://root:@${COCKROACH_HOST}/movr?sslmode=disable"
export SKR_ELASTICSEARCH_URI="http://127.0.0.1:9200"
export SKR_KAFKA_URI="localhost:9092"
```

## Set CockroachDB cluster settings

```bash
docker/cockroach_init.sh
```

## Start the api

```bash
go run main.go api -c processing/examples/schemas/
```

## Golang test

We will perform test with coverage like so:
```bash
go test -v ./... -coverprofile=coverage.out

# Open coverage result in your browser
go tool cover -html=coverage.out

# Get coverage result in CLI
go tool cover -func=coverage.out
```

## Redpanda cleaning

```bash
sudo docker stop redpandac01 redpandac02 redpandac03
sudo docker rm redpandac01 redpandac02 redpandac03
sudo docker volume rm -f docker_redpandac01 docker_redpandac02 docker_redpandac03
```
