# Contributing

In order to contribute to our project, we will define here all requirements.

## Git hooks

You need first do install [golangci-lint](https://golangci-lint.run/usage/install/)

Then, enable the hook in our project:
```bash
git config core.hooksPath .githooks
```

## Required packages

As prerequites, these must be installed:
- [JQ binary](https://github.com/stedolan/jq/releases)
- [CockroachDB binary](https://www.cockroachlabs.com/docs/releases/)
- [Redpanda binary](https://github.com/redpanda-data/redpanda/releases)

## Start all clusters

Set `sysctl` values permanently for elasticsearch if not already done in your linux host:
```bash
cat <<EOF | sudo tee -a /etc/sysctl.d/10-custom.conf
vm.max_map_count=262144
EOF
sudo sysctl -p /etc/sysctl.d/10-custom.conf
```

Start all clusters:
```bash
sudo docker-compose -f docker/cockroach-demo/docker-compose.yaml up -d
```

## Start CockroachDB cluster

To start the cluster execute this command in a separate terminal:
```bash
cockroach demo movr --insecure --http-port 8090
```
Don't quit otherwise, this cluster will be shut down.

## Start the api with default environment variables

```bash
export COCKROACH_HOST=$(netstat -latn |grep 26257 |grep LISTEN |awk '{print $4}' | head -1)
export SYNKER_PG_URI="postgres://root:@${COCKROACH_HOST}/movr?sslmode=disable"
export SYNKER_ELASTICSEARCH_URI="http://127.0.0.1:9200"
export SYNKER_KAFKA_URI="localhost:19092"

go run main.go api -c processing/examples/schemas/
```

## Execute workload against the stack

```bash
export COCKROACH_HOST=$(netstat -latn |grep 26257 |grep LISTEN |awk '{print $4}' | head -1)
export SYNKER_PG_URI="postgres://root:@${COCKROACH_HOST}/movr?sslmode=disable"
cockroach workload run movr --duration=1m "${SYNKER_PG_URI}"
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

## Cleaning everything

```bash
sudo docker-compose down -v
# clean and restart:
sudo docker-compose down -v && sleep 5 && sudo docker-compose up -d && watch curl -s 0:9200/_cat/indices?v
```

## Removing only topics and indexes

```bash
# elasticsearch
for i in $(curl -s "${SYNKER_ELASTICSEARCH_URI}/_cat/indices?pretty" | grep -v geoip_databases |awk '{print $3}'); do curl -XDELETE ${SYNKER_ELASTICSEARCH_URI}/$i;done

# topics
rpk topic delete --brokers ${SYNKER_KAFKA_URI} -r '.*'
```
