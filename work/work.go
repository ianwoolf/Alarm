package work

import (
	"fmt"

	"time"

	"github.com/lodastack/alarm/cluster"
	"github.com/lodastack/alarm/loda"
	// "github.com/lodastack/alarm/models"
	m "github.com/lodastack/models"

	"github.com/lodastack/log"
)

var (
	interval time.Duration = 10
)

type Work struct {
	Cluster cluster.ClusterInf
}

func NewWork(c cluster.ClusterInf) *Work {
	return &Work{Cluster: c}
}

func (w *Work) UpdateAlarms() error {
	loda.Loda.RLock()
	defer loda.Loda.RUnlock()
	for ns, alarms := range loda.Loda.NsAlarms {
		if len(alarms) == 0 {
			continue
		}

		_, err := w.Cluster.Get(ns, nil)
		if err != nil {
			log.Errorf("get ns fail: %s", err.Error())
			if err := w.Cluster.Lock("", time.Millisecond*10); err != nil {
				log.Errorf("work lock ns %s error: %s", ns, err.Error())
				continue //
			} else {
				if err := w.Cluster.Set(ns, ns, nil); err != nil {
					log.Errorf("work set ns %s error: %s", ns, err.Error())
					continue
				}
			}
			if err := w.Cluster.Unlock(""); err != nil {
				log.Errorf("work unlock ns %s error: %s", ns, err.Error())
			}
		}

		for _, alarm := range alarms {
			_, err := w.Cluster.Get(ns+"/"+alarm.Version, nil)
			if err != nil {
			}
			_, err = w.Cluster.Get(ns+"/"+alarm.Version+"/resource", nil)
			if err == nil {
				// unmarshal alarmInEtcd
				// if alarmMatch()
				// continue
			}

			if err := w.Cluster.Lock(ns, time.Millisecond*10); err != nil {
				log.Errorf("work lock ns %s error: %s", ns, err.Error())
				continue //
			} else {

				if err := w.Cluster.Set(ns+"/"+alarm.Version, ns+"/"+alarm.Version, nil); err != nil {
					// log.Errorf("work set ns %s error: %s", ns, err.Error())
					// continue //
				}

				if err := w.Cluster.Set(ns+"/"+alarm.Version+"/resource", ns+"/"+alarm.Version, nil); err != nil {
					log.Errorf("work set ns %s error: %s", ns, err.Error())
					// continue //
				}
			}
			if err := w.Cluster.Unlock(ns); err != nil && err.Error() != "100: Key not found" {
				log.Errorf("work unlock ns %s error: %s", ns, err.Error())
			}

			/*
					if _, ok := models.Work.NSs[ns].Alarms[alarm.Name]; !ok {
						alarmBlock := models.NewAlarmBlock(alarm)
						models.Work.NSs[ns].Alarms[alarm.Name] = &alarmBlock
						continue
					}
					if err := updateAlarm(ns, alarm); err != nil {
						log.Errorf("UpdateAlarms fail: %s", err.Error())
					}
				}
				models.Work.NSs[ns].Unlock()
			*/
		}
	}
	return nil
}

// func updateAlarm(ns string, alarm m.Alarm) error {
// 	if !alarmChanged(alarm, models.Work.NSs[ns].Alarms[alarm.Name].Alarm) {
// 		return nil
// 	}
// 	return models.Work.NSs[ns].Alarms[alarm.Name].UpdateAlarm(alarm)
// }

// alarmChanged return true if one not match with another,
// otherwise return false.
func alarmChanged(one, another m.Alarm) bool {
	// check:
	// alert
	// blockHostTimes blockHostPeroid
	// blockNSTimes blockNSPeroid
	return true
}

func (w *Work) CheckRegistryAlarmLoop() {
	for {
		loda.Loda.RLock()
		if len(loda.Loda.NsAlarms) != 0 {
			log.Info("loda resource init finished.")
			break
		}
		loda.Loda.RUnlock()
		time.Sleep(10 * time.Millisecond)
	}

	for {
		if err := w.UpdateAlarms(); err != nil {
			fmt.Println("work loop error:", err)
		}
		time.Sleep(interval * time.Second)
	}
}
