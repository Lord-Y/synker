# CockroachDB

To be able to use `changefeed enterprise feature` feature, we will use `cockroach demo` as there is no license requested for it and it will suit our testing purposes.

Here are the useful links:
- [cockroach demo](https://www.cockroachlabs.com/docs/v21.2/simulate-a-multi-region-cluster-on-localhost.html)
- [changefeed enterprise feature](https://www.cockroachlabs.com/docs/stable/stream-data-out-of-cockroachdb-using-changefeeds.html#configure-a-changefeed-enterprise)

## Start the instance

To start the instance execute this command:
```bash
cockroach demo movr --geo-partitioned-replicas --insecure
```

In an another shell, execute `docker/cockroach_init.sh` script to enable `rangefeed setting`.
