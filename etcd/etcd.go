package etcd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/lodastack/alarm/models"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

var (
	alarmsDirect = "loda-alarms"
)

// Client is a wrapper around the etcd client
type EtcdClient struct {
	client.KeysAPI
	EndPoints []string
}

func NewEtcdClient(endpoints []string, basicAuth bool, username, password string,
	headTimeout time.Duration) (*EtcdClient, error) {
	var c client.Client
	var kapi client.KeysAPI
	var err error

	cfg := client.Config{
		Endpoints: endpoints,
	}
	if headTimeout != 0 {
		cfg.HeaderTimeoutPerRequest = headTimeout
	}
	if basicAuth {
		cfg.Username = username
		cfg.Password = password
	}

	c, err = client.New(cfg)
	if err != nil {
		return &EtcdClient{kapi, endpoints}, err
	}

	kapi = client.NewKeysAPI(c)
	return &EtcdClient{kapi, endpoints}, nil
}

func (c *EtcdClient) alarmInfo(node *client.Node) (*models.AlarmInfo, error) {
	var info *models.AlarmInfo
	err := json.Unmarshal([]byte(node.Value), info)
	return info, err
}

func (c *EtcdClient) Worker() []string    { return nil }
func (c *EtcdClient) addWorker() error    { return nil }
func (c *EtcdClient) removeWorder() error { return nil }
func (c *EtcdClient) updateWorker() error { return nil }

func (c *EtcdClient) Watch() error {
	watcher := c.Watcher(alarmsDirect, &client.WatcherOptions{
		Recursive: true,
	})
	for {
		res, err := watcher.Next(context.Background())
		if err != nil {
			fmt.Println("Error watch workers:", err)
			break
		}
		if res.Action == "expire" {
			info := NodeToWorkerInfo(res.PrevNode)
			fmt.Println("Expire worker ", info.Name)
			member, ok := m.members[info.Name]
			if ok {
				member.InGroup = false
			}
		} else if res.Action == "set" {
			info := NodeToWorkerInfo(res.Node)
			if _, ok := m.members[info.Name]; ok {
				fmt.Println("Update worker ", info.Name)
				c.UpdateWorker(info)
			} else {
				fmt.Println("Add worker ", info.Name)
				c.AddWorker(info)
			}
		} else if res.Action == "delete" {
			info := NodeToWorkerInfo(res.Node)
			fmt.Println("Delete worker ", info.Name)
			delete(c.members, info.Name)
		}
	}

	return nil
}
