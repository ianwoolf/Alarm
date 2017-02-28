package models

import (
	"errors"
	"sync"
)

type Alarm struct {
	Addr  string
	Alive bool
}

type AlarmCluster struct {
	sync.RWMutex
	Cluster map[string]*Alarm
}

func NewAlarmCluster() AlarmCluster {
	return AlarmCluster{sync.RWMutex{}, map[string]*Alarm{}}
}

func (ac *AlarmCluster) Exist(addr string) bool {
	ac.RLock()
	defer ac.RUnlock()
	_, ok := ac.Cluster[addr]
	return ok
}

func (ac *AlarmCluster) ReadAlive() ([]string, error) {
	ac.RLock()
	defer ac.RUnlock()
	var aliveAddr []string
	var index int
	aliveAddr = make([]string, len(ac.Cluster))
	for _, alarm := range ac.Cluster {
		if alarm.Alive {
			aliveAddr[index] = alarm.Addr
			index++
		}
	}
	return aliveAddr[:index], nil
}

func (ac *AlarmCluster) Add(addr string) error {
	ac.Lock()
	defer ac.Unlock()
	if _, ok := ac.Cluster[addr]; ok {
		return errors.New("alarm already exist")
	}
	ac.Cluster[addr] = &Alarm{addr, true}
	return nil
}

func (ac *AlarmCluster) Delete(addr string) error {
	ac.Lock()
	defer ac.Unlock()
	delete(ac.Cluster, addr)
	return nil
}

func (ac *AlarmCluster) update(addr string, status bool) error {
	if _, ok := ac.Cluster[addr]; !ok {
		return errors.New("alarm already not exist")
	}
	ac.Cluster[addr].Alive = status
	return nil
}

func (ac *AlarmCluster) SetAlive(addr string) error {
	ac.Lock()
	defer ac.Unlock()
	return ac.update(addr, true)
}

func (ac *AlarmCluster) SetNotAlive(addr string) error {
	ac.Lock()
	defer ac.Unlock()
	return ac.update(addr, false)
}
