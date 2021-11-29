# CockroachDB

To be able to use `changefeed enterprise feature` feature, we will use `cockroach demo` as there is no license requested for it and it will suit our testing purposes.

Cockroach version 21.2.0 minimum is required if you want to use `full_table_name, on_error = 'pause', protect_data_from_gc_on_pause` params.

Here are the useful links:
- [cockroach demo](https://www.cockroachlabs.com/docs/v21.2/simulate-a-multi-region-cluster-on-localhost.html)
- [changefeed enterprise feature](https://www.cockroachlabs.com/docs/stable/stream-data-out-of-cockroachdb-using-changefeeds.html#configure-a-changefeed-enterprise)
- [CDC - Fine Tuning Changefeeds for Performance and Durability](https://www.cockroachlabs.com/blog/change-data-capture-for-performance-durability/)

## Start the instance

To start the instance execute this command:
```bash
cockroach demo movr --geo-partitioned-replicas --insecure --http-port 18080
```

In an another shell, execute `docker/cockroach_init.sh` script to enable `rangefeed setting`.
