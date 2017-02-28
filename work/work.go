package work

import (
	"fmt"
	"time"

	"github.com/lodastack/alarm/loda"
	"github.com/lodastack/alarm/models"
)

var (
	interval time.Duration = 10
)

func UpdateAlarms() error {
	allNs, err := loda.AllNS("")
	if err != nil {
		fmt.Println("UpdateAlarms error:", err)
		return err
	}
	for _, ns := range allNs {
		// if ns hash
		alarms, err := loda.GetAlarmsByNs(ns)
		if err != nil {
			return err
		}

		for _, alarm := range alarms {
			alarmBlock := models.NewAlarmBlock(ns, models.SendBlock, alarm)
			models.Work[ns] = models.NsBlock{
				AlarmBlocks: map[string]models.AlarmBlock{alarm.Version: alarmBlock},
			}
		}
	}
	return nil
}

func Loop() {
	for {
		if err := UpdateAlarms(); err != nil {
			fmt.Println("work loop error:", err)
		}
		time.Sleep(interval * time.Second)
	}
}
