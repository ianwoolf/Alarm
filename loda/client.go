package loda

import (
	"fmt"
	"sync"
	"time"

	"github.com/lodastack/alarm/common"

	"github.com/lodastack/log"
	m "github.com/lodastack/models"
)

var (
	readRegInterval time.Duration = 1
	Loda            lodaAlarm
)

type lodaAlarm struct {
	sync.RWMutex
	NsAlarms map[string][]m.Alarm
}

func init() {
	Loda = lodaAlarm{
		NsAlarms: make(map[string][]m.Alarm),
	}
}

func (l *lodaAlarm) UpdateAlarms() error {
	allNs, err := AllNS("")
	if err != nil {
		fmt.Println("UpdateAlarms error:", err)
		return err
	}
	l.Lock()
	defer l.Unlock()

	for ns := range l.NsAlarms {
		if _, ok := common.ContainString(allNs, ns); !ok {
			delete(l.NsAlarms, ns)
		}
	}

	for _, ns := range allNs {
		alarms, err := GetAlarmsByNs(ns)
		if err != nil {
			log.Errorf("get alarm fail: %s", err.Error())
			return err
		}
		if len(alarms) == 0 {
			continue
		}
		l.NsAlarms[ns] = alarms
	}
	return nil
}

func ReadLoop() error {
	for {
		if err := Loda.UpdateAlarms(); err != nil {
			log.Errorf("loda ReadLoop fail: %s", err.Error())
		}
		time.Sleep(readRegInterval * time.Minute)
	}
	return nil
}
