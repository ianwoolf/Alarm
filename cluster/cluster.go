package cluster

import (
	"fmt"
	"time"

	"github.com/lodastack/alarm/cluster/etcd"
	"github.com/lodastack/alarm/models"
)

type ClusterInf interface {
	AliveNodes() ([]string, error)
	Run() error
}

type Cluster struct {
	Self   string
	TTL    time.Duration
	Alarms models.AlarmCluster
	Etcd   etcd.EtcdInf
}

func NewCluster(selfAddr string, endpoints []string, basicAuth bool, username, password string,
	headTimeout, nodeTTL time.Duration) (ClusterInf, error) {
	cluster := Cluster{Self: selfAddr, TTL: nodeTTL, Alarms: models.NewAlarmCluster()}
	etcdClient, err := etcd.NewEtcdClient(endpoints, basicAuth, username, password,
		headTimeout, nodeTTL)

	if err != nil {
		fmt.Println("NewCluster error:", err)
		return &cluster, err
	}
	cluster.Etcd = etcdClient

	return &cluster, nil

}
func (c *Cluster) Run() error {
	go func() {
		err := c.Etcd.Watch()
		fmt.Println("cluster run watch error:", err)
	}()

	go func() {
		for {
			err := c.Etcd.HeartBeat(c.Self)
			if err != nil {
				fmt.Println("cluster run heardbeat:", err)
			}
			time.Sleep(time.Second * c.TTL)
		}
	}()

	for {
		var err error
		clusterMsg := c.Etcd.Listen()
		switch clusterMsg.Action {
		case etcd.Expire:
			err = c.Alarms.SetNotAlive(clusterMsg.Addr)
		case etcd.Set:
			if c.Alarms.Exist(clusterMsg.Addr) {
				err = c.Alarms.SetAlive(clusterMsg.Addr)
			} else {
				err = c.Alarms.Add(clusterMsg.Addr)
			}
		case etcd.Delete:
			err = c.Alarms.Delete(clusterMsg.Addr)

		default:
			fmt.Println("error")
		}
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func (c *Cluster) AliveNodes() ([]string, error) {
	return c.Alarms.ReadAlive()
}
