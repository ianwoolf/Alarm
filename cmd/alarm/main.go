package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/lodastack/alarm/cluster"
	"github.com/lodastack/alarm/config"
	"github.com/lodastack/alarm/loda"
	"github.com/lodastack/alarm/models"
	"github.com/lodastack/alarm/work"

	"github.com/lodastack/log"
)

func initLog(conf config.LogConfig) {
	if !conf.Enable {
		fmt.Println("log to std err")
		log.LogToStderr()
		return
	}

	if backend, err := log.NewFileBackend(conf.Path); err != nil {
		fmt.Fprintf(os.Stderr, "init logs folder failed: %s\n", err.Error())
		os.Exit(1)
	} else {
		log.SetLogging(conf.Level, backend)
		backend.Rotate(conf.FileNum, uint64(1024*1024*conf.FileSize))
	}
}

func init() {
	configFile := flag.String("c", "./conf/alarm.conf", "config file path")
	flag.Parse()
	fmt.Printf("load config from %s\n", *configFile)
	err := config.LoadConfig(*configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read config file failed:\n%s\n", err.Error())
		os.Exit(1)
	}
	initLog(config.GetConfig().Log)
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	fmt.Println("build via golang version ", runtime.Version())
	c, err := cluster.NewCluster(config.GetConfig().Etcd.Addr, config.GetConfig().Etcd.Endpoints, false, "", "", 5, 5)
	if err != nil {
		fmt.Println("main error", err)
		return
	}
	go loda.ReadLoop()

	w := work.NewWork(c)
	go w.CheckRegistryAlarmLoop()
	for {
		models.Println()

		time.Sleep(time.Second)
	}
	// go query.Start()
	// go loda.PurgeAll()
	select {}
}
