package etcd

type EtcdInterface interface {
	Worker() []string
	Watch() error
}
