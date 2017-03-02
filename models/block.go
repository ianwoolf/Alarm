package models

import (
	"fmt"
	"sync"
	// "time"

	m "github.com/lodastack/models"
)

func init() {
	Work = LodaEvent{sync.RWMutex{}, make(map[string]*NsBlock)}
}

var (
	SendBlock = "block"
	SendPoint = "point"
)

type BlockWork struct {
	sync.RWMutex

	SendType string

	// TODO: sleep to check and send
	block       bool
	blockNum    int // Alarm ADD
	blockPeriod int // Alarm ADD
	queue       []*Point
}

type AlarmBlock struct {
	sync.RWMutex
	Alarm  m.Alarm // TODO
	ByNs   BlockWork
	ByHost map[string]*BlockWork
}

func NewAlarmBlock(alarm m.Alarm) AlarmBlock {
	block := AlarmBlock{
		Alarm: alarm,
		ByNs: BlockWork{
			SendType: SendBlock,
			// blockNum:0,
			// blockPeriod:0,
			// queen: make([]*Point, block),
		},
		ByHost: make(map[string]*BlockWork),
	}
	return block
}

func (b *BlockWork) Log() error              { return nil }
func (b *BlockWork) QueryHistory()           {}
func (b *BlockWork) event(point Point) error { return nil }

// check the cache, if match the send value, return true.
func (b *BlockWork) Check() (string, bool) {
	switch b.SendType {
	case SendBlock:
	// if match block all.
	case SendPoint:
	}
	return "", true
}
func (b *BlockWork) Update(alrm m.Alarm) error {
	b.Lock()
	defer b.Unlock()
	// b.block=alarm.
	// b.period=alarm.
	return nil
}

func (ab *AlarmBlock) Dispatch(msg string, point Point) error {
	return nil
}

// push the point to ab.ByNs and ab.ByHost
func (ab *AlarmBlock) Event(point Point) error {
	host := point.Host()
	if host == "" {
		return fmt.Errorf("point have no host, %v", point)
	}
	ab.ByNs.queue = append(ab.ByNs.queue, &point)

	if _, ok := ab.ByHost[host]; !ok {
		ab.ByHost[host] = &BlockWork{
			SendType: SendPoint,
			// blockNum:0,
			// blockPeriod:0,
			// queen: make([]*Point, block),
		}
	}
	ab.ByHost[host].queue = append(ab.ByHost[host].queue, &point)

	return nil
}

func (ab *AlarmBlock) Send(point Point) error {
	host := point.Host()
	if host == "" {
		return fmt.Errorf("point have no host, %v", point)
	}

	if msg, ok := ab.ByNs.Check(); ok {
		ab.Dispatch(msg, point)
		ab.ByNs.block = true
		return nil
	} else {
		ab.ByNs.block = false
	}

	if _, ok := ab.ByHost[host]; !ok {
		return fmt.Errorf("alarmBlock send point %+v by host fail: have not this host %s block",
			point, host)
	}
	if msg, ok := ab.ByHost[host].Check(); ok {
		ab.Dispatch(msg, point)
	}

	return nil
}

func (ab *AlarmBlock) UpdateAlarm(alarm m.Alarm) error {
	ab.Lock()
	defer ab.Unlock()

	if !checkAlarm(alarm) {
		return fmt.Errorf("alarm invalid")
	}

	ab.Alarm = alarm
	if err := ab.ByNs.Update(alarm); err != nil {
		fmt.Println("udpate nsblock fail:", err)
	}
	for host := range ab.ByHost {
		if err := ab.ByHost[host].Update(alarm); err != nil {
			fmt.Printf("udpate %s block fail: %s", host, err.Error())
		}
	}
	return nil
}

type NsBlock struct {
	sync.RWMutex
	Alarms map[string]*AlarmBlock
}

type LodaEvent struct {
	sync.RWMutex
	NSs map[string]*NsBlock
}

var Work LodaEvent

func checkAlarm(alarm m.Alarm) bool {
	// check alarm is valid or not:
	// alert
	// blockHostTimes blockHostPeroid
	// blockNSTimes blockNSPeroid
	return true
}

// DEBUG
func Println() {
	for ns := range Work.NSs {
		fmt.Printf("#### ns: %s, alarm num: %d\n", ns, len(Work.NSs[ns].Alarms))
		// for alarm := range Work.NSs[ns].Alarms {
		// 	fmt.Println("\t#### alarm:", alarm)
		// 	fmt.Printf("\t#### alarm content: %+v", Work.NSs[ns].Alarms[alarm])
		// }
	}
}
