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
sudo docker-compose -f docker/docker-compose.yaml up -d
```

## Start CockroachDB cluster

To start the cluster execute this command in a separate terminal:
```bash
cockroach demo movr --geo-partitioned-replicas --insecure --http-port 18080
```

## Set default variables

In the terminal you will execute your `go run`: 
```bash
export COCKROACH_HOST=$(netstat -latn |grep 26257 |grep LISTEN |awk '{print $4}')
export SYNKER_PG_URI="postgres://root:@${COCKROACH_HOST}/movr?sslmode=disable"
export SYNKER_ELASTICSEARCH_URI="http://127.0.0.1:9200"
export SYNKER_KAFKA_URI="localhost:9092"
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

## Repanda/elasticsearch watching

Watching what's in the queue `promo_codes`:
```bash
rpk topic consume movr.public.promo_codes
```

Watching what's in `elasticsearch`:
```bash
curl -s "${SYNKER_ELASTICSEARCH_URI}/_cat/indices?pretty&v"
```

## Cleaning all

```bash
sudo docker stop redpandac01 redpandac02 redpandac03
sudo docker rm redpandac01 redpandac02 redpandac03
sudo docker volume rm -f docker_redpandac01 docker_redpandac02 docker_redpandac03
sudo docker stop elasticsearchc01 elasticsearchc02 elasticsearchc03
sudo docker rm elasticsearchc01 elasticsearchc02 elasticsearchc03
sudo docker volume rm -f docker_elasticsearchc01 docker_elasticsearchc02 docker_elasticsearchc03
```

## Removing only topics and indexes

```bash
# elasticsearch
for i in $(curl -s "${SYNKER_ELASTICSEARCH_URI}/_cat/indices?pretty" | grep -v geoip_databases |awk '{print $3}'); do curl -XDELETE ${SYNKER_ELASTICSEARCH_URI}/$i;done

# topics
rpk topic delete -r '.*'
```