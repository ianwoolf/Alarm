package loda

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/lodastack/alarm/config"
	"github.com/lodastack/alarm/requests"

	"github.com/lodastack/log"
	"github.com/lodastack/models"
)

const CommonCluster = "common"
const DefaultDBNameSpace = "db.monitor.loda"
const MachineUri = "/api/v1/router/resource?ns=%s&type=machine"
const AlarmUri = "/api/v1/router/resource?ns=%s&type=alarm"

var PurgeChan chan string
var Client *client

type client struct {
	// cache ns -> dbs in this map
	db map[string][]string
	mu sync.RWMutex
}

type RespNS struct {
	Status int      `json:"httpstatus"`
	Data   []string `json:"data"`
}

type RespDB struct {
	Status int      `json:"httpstatus"`
	Data   []Server `json:"data"`
}

type ResAlarm struct {
	HttpStatus int            `json:"httpstatus"`
	Data       []models.Alarm `json:"data"`
}

type Server struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
}

func init() {
	PurgeChan = make(chan string)
	Client = &client{
		db: make(map[string][]string),
	}
}

func PurgeAll() {
	var ticker *time.Ticker
	interval := config.GetConfig().Reg.ExpireDur
	if interval < 60 {
		interval = 60
	}
	duration := time.Duration(interval) * time.Second
	ticker = time.NewTicker(duration)
	for {
		select {
		case <-ticker.C:
			res, err := AllNS("")
			if err == nil {
				log.Infof("DB old cache: %v", Client.db)
				for _, ns := range res {
					dbs, err := updateInfluxDBs(ns)
					if err == nil {
						Client.mu.Lock()
						Client.db["collect."+ns] = dbs
						Client.mu.Unlock()
					} else {
						log.Errorf("update ns: %s cache failed: %s", ns, err)
					}
				}
				log.Infof("DB new cache: %v", Client.db)
			} else {
				log.Errorf("Get all NS failed: %s", err)
			}
		case ns := <-PurgeChan:
			Client.purge(ns)
		}
	}
}

func (c *client) purge(ns string) {
	c.mu.Lock()
	if _, ok := c.db[ns]; ok {
		delete(c.db, ns)
	}
	c.mu.Unlock()
	log.Infof("purge cache ns:%s", ns)
}

func InfluxDBs(ns string) ([]string, error) {
	var res []string
	var ok bool
	Client.mu.RLock()
	if res, ok = Client.db[ns]; ok {
		Client.mu.RUnlock()
		return res, nil
	}
	Client.mu.RUnlock()
	dbs, err := updateInfluxDBs(ns)
	if err != nil {
		return res, err
	}
	Client.mu.Lock()
	Client.db[ns] = dbs
	Client.mu.Unlock()
	return dbs, nil
}

func updateInfluxDBs(ns string) ([]string, error) {
	list := strings.Split(ns, ".")
	if len(list)-2 < 0 {
		return []string{}, fmt.Errorf("ns error: %s", ns)
	}
	partone := list[len(list)-2]
	uri := fmt.Sprintf(MachineUri, partone+"."+DefaultDBNameSpace)
	url := fmt.Sprintf("%s%s", config.GetConfig().Reg.Link, uri)
	res, err := servers(url)
	if err != nil || len(res) > 0 {
		return res, err
	}

	url = fmt.Sprintf("%s/api/v1/router/ns?ns=%s&format=list", config.GetConfig().Reg.Link, DefaultDBNameSpace)
	res, err = AllNS(url)
	if err == nil {
		ok, cluster := includeNS(partone, res)
		if ok {
			uri = fmt.Sprintf(MachineUri, cluster+"."+DefaultDBNameSpace)
			url = fmt.Sprintf("%s%s", config.GetConfig().Reg.Link, uri)
			res, err = servers(url)
			if err != nil || len(res) > 0 {
				return res, err
			}
		}
	} else {
		log.Errorf("get default DB NameSpace failed: %s", err)
	}

	// Send to common cluster if not found customer cluster
	uri = fmt.Sprintf(MachineUri, CommonCluster+"."+DefaultDBNameSpace)
	url = fmt.Sprintf("%s%s", config.GetConfig().Reg.Link, uri)
	res, err = servers(url)
	if err != nil || len(res) > 0 {
		return res, err
	}

	return []string{"influxdb.ifengidc.com"}, fmt.Errorf("common cluster status != 200")
}

func servers(url string) ([]string, error) {
	var res []string
	var resdb RespDB

	resp, err := requests.Get(url)
	if err != nil {
		return res, err
	}

	if resp.Status == 200 {
		err = json.Unmarshal(resp.Body, &resdb)
		if err != nil {
			return res, err
		}
		for _, s := range resdb.Data {
			res = append(res, s.IP)
		}
		return res, nil
	}
	// len(res) == 0
	return res, nil
}

func AllNS(url string) ([]string, error) {
	var resNS RespNS
	var res []string
	if url == "" {
		url = fmt.Sprintf("%s/api/v1/router/ns?ns=&format=list", config.GetConfig().Reg.Link)
	}
	resp, err := requests.Get(url)
	if err != nil {
		return res, err
	}

	if resp.Status == 200 {
		err = json.Unmarshal(resp.Body, &resNS)
		if err != nil {
			return res, err
		}
		return resNS.Data, nil
	}
	return res, fmt.Errorf("http status code: %d", resp.Status)
}

func GetAlarmsByNs(ns string) ([]models.Alarm, error) {
	var resAlarms ResAlarm
	resAlarms = ResAlarm{} // TODO

	url := fmt.Sprintf("%s"+AlarmUri, config.GetConfig().Reg.Link, ns)
	resp, err := requests.Get(url)
	if err != nil {
		return resAlarms.Data, err
	}

	if resp.Status != 200 {
		return resAlarms.Data, fmt.Errorf("query registry error")
	}
	err = json.Unmarshal(resp.Body, &resAlarms)
	return resAlarms.Data, err
}

func includeNS(nsPartOne string, dbs []string) (bool, string) {
	for _, dbns := range dbs {
		parts := strings.Split(dbns, ".")
		// use "||" splite mutile NS part one
		nodes := strings.Split(parts[0], "||")
		for _, name := range nodes {
			if name == nsPartOne {
				return true, parts[0]
			}
		}
	}
	return false, ""
}
