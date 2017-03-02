package cluster

import (
	"time"

	"github.com/lodastack/alarm/cluster/etcd"
	"github.com/lodastack/alarm/models"

	"github.com/coreos/etcd/client"
	"github.com/lodastack/log"
)

type ClusterInf interface {
	Get(k string, option *client.GetOptions) (*client.Response, error)
	Set(k, v string, option *client.SetOptions) error
	Delete(key string) error
	Lock(path string, lockTime time.Duration) error
	Unlock(path string) error
	RecursiveGet(k string) (*client.Response, error)
}

type Cluster struct {
	etcd.EtcdClient

	Self   string
	TTL    time.Duration
	Alarms models.AlarmCluster
}

func NewCluster(selfAddr string, endpoints []string, basicAuth bool, username, password string,
	headTimeout, nodeTTL time.Duration) (ClusterInf, error) {

	etcdClient, err := etcd.NewEtcdClient(endpoints, basicAuth, username, password,
		headTimeout, nodeTTL)

	if err != nil {
		log.Errorf("NewCluster error: %s", err.Error())
		return nil, err
	}

	cluster := Cluster{etcdClient, selfAddr, nodeTTL, models.NewAlarmCluster()}
	return &cluster, nil
}
