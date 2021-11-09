// Package etcd assemble all requirements to communicate with etcd cluster
package etcd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Lord-Y/synker/commons"
	"go.etcd.io/etcd/client/pkg/v3/transport"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	dialTimeout      time.Duration = 5 * time.Second
	operationTimeout time.Duration = 2 * time.Second
	keyPrefix        string        = "/synker"
)

// Client return requirements to connect to etcd cluster
func Client() (cli *clientv3.Client, err error) {
	var (
		tlsInfo transport.TLSInfo
	)
	if commons.GetEtcdCACert() != "" && commons.GetEtcdCert() != "" && commons.GetEtcdKey() != "" {
		tlsInfo = transport.TLSInfo{
			CertFile:      commons.GetEtcdCert(),
			KeyFile:       commons.GetEtcdKey(),
			TrustedCAFile: commons.GetEtcdCACert(),
		}
		tlsConfig, err := tlsInfo.ClientConfig()
		if err != nil {
			return &clientv3.Client{}, err
		}
		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   commons.GetEtcdURI(),
			DialTimeout: dialTimeout,
			TLS:         tlsConfig,
			Username:    commons.GetEtcdUser(),
			Password:    commons.GetEtcdPassword(),
		})
		if err != nil {
			return &clientv3.Client{}, err
		}
	} else {
		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   commons.GetEtcdURI(),
			DialTimeout: dialTimeout,
			Username:    commons.GetEtcdUser(),
			Password:    commons.GetEtcdPassword(),
		})
		if err != nil {
			return &clientv3.Client{}, err
		}
	}

	return cli, err
}

// Put permit to add/update key into etcd cluster
func Put(cli *clientv3.Client, key, value string) (resp *clientv3.PutResponse, err error) {
	defer cli.Close()
	if strings.TrimSpace(key) == "" {
		return resp, fmt.Errorf("Key cannot be empty")
	}
	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	resp, err = cli.Put(
		ctx,
		fmt.Sprintf("%s/%s", keyPrefix, key),
		value,
	)
	cancel()
	return
}

// Get permit to retrieve key into etcd cluster
func Get(cli *clientv3.Client, key string) (resp *clientv3.GetResponse, err error) {
	defer cli.Close()
	if strings.TrimSpace(key) == "" {
		return resp, fmt.Errorf("Key cannot be empty")
	}
	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	resp, err = cli.Get(
		ctx,
		fmt.Sprintf("%s/%s", keyPrefix, key),
	)
	cancel()
	return
}

// Watch permit to retrieve key into etcd cluster
func Watch(cli *clientv3.Client, key string) (resp clientv3.WatchChan, err error) {
	if strings.TrimSpace(key) == "" {
		return resp, fmt.Errorf("Key cannot be empty")
	}
	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	resp = cli.Watch(
		ctx,
		fmt.Sprintf("%s/%s", keyPrefix, key),
	)
	cancel()
	return
}

// WatchWithPrefix permit to retrieve key into etcd cluster
func WatchWithPrefix(cli *clientv3.Client, key string) (resp clientv3.WatchChan, err error) {
	if strings.TrimSpace(key) == "" {
		return resp, fmt.Errorf("Key cannot be empty")
	}
	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	resp = cli.Watch(
		ctx,
		fmt.Sprintf("%s/%s", keyPrefix, key),
		clientv3.WithPrefix(),
	)
	cancel()
	return
}

// Delete permit to retrieve key into etcd cluster
func Delete(cli *clientv3.Client, key string) (resp *clientv3.DeleteResponse, err error) {
	defer cli.Close()
	if strings.TrimSpace(key) == "" {
		return resp, fmt.Errorf("Key cannot be empty")
	}
	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	resp, err = cli.Delete(
		ctx,
		fmt.Sprintf("%s/%s", keyPrefix, key),
	)
	cancel()
	return
}
