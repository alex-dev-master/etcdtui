package etcd

import client "github.com/alexandr/etcdtui/pkg/etcd"

func NewEtcdClient(cfg *client.Config) (*client.Client, error) {
	return client.New(cfg)
}
