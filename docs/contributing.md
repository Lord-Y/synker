# Contributing

In order to contribute to our project, we will define here all requirements.

## Git hooks

You need first do install [golangci-lint](https://golangci-lint.run/usage/install/)

Then, enable the hook in our project:
```bash
git config core.hooksPath .githooks
```

## Starting the cluster

Set `sysctl` values permanently for elasticsearch if not already done in your user environment:
```bash
cat <<EOF | sudo tee -a /etc/sysctl.d/10-custom.conf
vm.max_map_count=262144
EOF
sudo sysctl -p /etc/sysctl.d/10-custom.conf
```
Start the cluster:
```bash
sudo docker-compose -f docker/docker-compose-cluster.yaml up -d
```

## Set default variables

```bash
export SKR_PG_URI="postgres://root:@127.0.0.1:26257/synker"
export SKR_ELASTICSEARCH_URI="http://127.0.0.1:9200"
export SKR_KAFKA_URI="localhost:9092"
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
