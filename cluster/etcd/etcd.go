package etcd

import (
	"fmt"
	"time"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

var (
	alarmsDirect                 = "loda-alarms"
	cacheCap                     = 10
	defaultNodeTTL time.Duration = 5

	Expire = "expire"
	Set    = "set"
	Delete = "delete"
)

type EtcdInf interface {
	Watch() error
	HeartBeat(addr string) error
	Listen() EtcdMsg
}

type EtcdMsg struct {
	Addr   string
	Action string
}

// Client is a wrapper around the etcd client
type EtcdClient struct {
	kapi       client.KeysAPI
	EndPoints  []string
	MsgChannle chan EtcdMsg
	TTL        time.Duration
}

func NewEtcdClient(endpoints []string, basicAuth bool, username, password string,
	headTimeout, nodeTTL time.Duration) (EtcdInf, error) {
	var c client.Client
	var kapi client.KeysAPI
	var err error

	cfg := client.Config{
		Endpoints: endpoints,
	}
	if headTimeout != 0 {
		cfg.HeaderTimeoutPerRequest = headTimeout * time.Second
	}
	if basicAuth {
		cfg.Username = username
		cfg.Password = password
	}

	c, err = client.New(cfg)
	if err != nil {
		return nil, err
	}
	if nodeTTL == 0 {
		nodeTTL = defaultNodeTTL
	}
	kapi = client.NewKeysAPI(c)
	return &EtcdClient{kapi: kapi,
		EndPoints:  endpoints,
		MsgChannle: make(chan EtcdMsg, cacheCap),
		TTL:        nodeTTL}, nil
}

func (c *EtcdClient) read(node *client.Node) string {
	return string(node.Value)
}

func (c *EtcdClient) Watch() error {
	watcher := c.kapi.Watcher(alarmsDirect, &client.WatcherOptions{
		Recursive: true,
	})
	for {
		res, err := watcher.Next(context.Background())
		if err != nil {
			fmt.Println("Error watch workers:", err)
			break
		}
		// fmt.Println("watch debug", res, c.read)
		var addr string
		switch res.Action {
		case Expire:
			addr = c.read(res.PrevNode)
		case Set:
			fallthrough
		case Delete:
			addr = c.read(res.Node)
		default:
			fmt.Println("unknow etcd message", res)
			continue
		}
		c.MsgChannle <- EtcdMsg{addr, res.Action}
	}

	return nil
}

func (c *EtcdClient) HeartBeat(addr string) error {
	key := alarmsDirect + "/" + addr
	_, err := c.kapi.Set(context.Background(), key, addr, &client.SetOptions{
		TTL: time.Second * c.TTL,
	})
	return err
}

func (c *EtcdClient) Listen() EtcdMsg {
	return <-c.MsgChannle
}
