#!/bin/bash

ENDPOINTS=http://etcdc01:2379
alias etcdctl="ENDPOINTS=${ENDPOINTS} etcdctl"
ETCD_ROLE_LIST=$(etcdctl role list)
if [ $? -eq 0 -a -z "${ETCD_ROLE_LIST}" ]
then
    etcdctl role add root
    etcdctl user add root --new-user-password ${USER_ROOT_PASSWORD}
    etcdctl user grant-role root root

    etcdctl role add readwrite
    etcdctl role grant-permission readwrite --prefix=true readwrite /synker/
    etcdctl user add synker --new-user-password ${USER_SYNKER_PASSWORD}
    etcdctl user grant-role synker readwrite

    etcdctl auth enable
else
    echo "Setup already done"
fi

exit 0
