# readpanda

For our purpose, we will use `redpanda` as an alternative to kafka and it's compatible with Kafka.
It means that everything you are doing with kafka, you will be able to do it with `redpanda`.
If you want to install is with helm on kubernetes we recommanded you to [follow the documentation here](https://vectorized.io/docs/quick-start-kubernetes).

## Authentication

Our development cluster is authenticated.
Here is the mapping:
```yaml
- role: all
  user: admin
- role: topic
  user: synker
```

To add role/user to your redpanda cluster, follow this [documentation](https://vectorized.io/docs/acls).

Authentication requirements will also be covered.
