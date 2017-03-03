package work

import (
	"time"

	"github.com/lodastack/alarm/cluster"
	"github.com/lodastack/alarm/loda"
	m "github.com/lodastack/models"

	"github.com/lodastack/log"
)

var (
	interval time.Duration = 20
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
		// create ns dir if not exist.
		if _, err := w.Cluster.Get(ns, nil); err != nil {
			log.Infof("get ns %s fail(%s) set it", ns, err.Error())
			if err := w.Cluster.CreateDir(ns); err != nil {
				log.Errorf("work set ns %s error: %s, skip this ns", ns, err.Error())
				continue
			}
		}

		// create alarm dir if not exit.
		for _, alarm := range alarms {
			alarmKey := ns + "/" + alarm.Version
			if _, err := w.Cluster.Get(alarmKey, nil); err != nil {
				log.Infof("get ns(%s) alarm(%s) fail: %s, set it and all dir.", ns, alarm.Version, err.Error())
				if err := w.Cluster.CreateDir(alarmKey); err != nil {
					log.Errorf("set ns(%s) alarm(%s) fail: %s, skip this alarm",
						ns, alarm.Version, err.Error())
					continue
				}
				allDirKey := alarmKey + "/all"
				if err := w.Cluster.CreateDir(allDirKey); err != nil {
					log.Errorf("set ns(%s) alarm(%s) dir \"all\" fail: %s, skip this alarm",
						ns, alarm.Version, err.Error())
					continue
				}
			}
		}
	}
	return nil
}

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
			log.Errorf("work loop error: %s", err)
		} else {
			log.Info("work loop success")
		}

		time.Sleep(interval * time.Second)
	}
}
