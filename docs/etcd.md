# etcd

For our purpose, we will use `etcd` as a key/value store.
If you want to install is with helm on kubernetes we recommanded you to se [bitnami/etcd chart](https://github.com/bitnami/charts/tree/master/bitnami/etcd)

If you want to setup it on virtual machines, you can find the [documentation here](https://etcd.io/docs/v3.5/op-guide/clustering/)

[Hardware recommendations with VM install](https://etcd.io/docs/v3.5/op-guide/hardware/).

## Interaction with etcd

In order to interact with etcd with `auto-tls` and `peer-auto-tls` options, connect to a container and execute this command:
```bash
# endpoint status
ENDPOINTS=https://0:2379 etcdctl --cert=/etcd/data/fixtures/client/cert.pem --key=/etcd/data/fixtures/client/key.pem --insecure-skip-tls-verify --insecure-transport endpoint status
# member list
ENDPOINTS=https://0:2379 etcdctl --cert=/etcd/data/fixtures/client/cert.pem --key=/etcd/data/fixtures/client/key.pem --insecure-skip-tls-verify --insecure-transport member list
```

Without tls options, it's more ease:
```bash
# endpoint status
etcdctl endpoint status --write-out table
# member list
etcdctl member list --write-out table
```

## Authentication

Our development cluster is authenticated.
Here is the mapping:
```yaml
- role: root
  user: root
- role: rw
  user: synker
```

To add role/user to your etcd cluster, follow this [documentation](https://etcd.io/docs/v3.5/demo/).

Authentication requirements can be found [here](https://etcd.io/docs/v3.5/op-guide/authentication/).


