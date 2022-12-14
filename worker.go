package main

import (
	"fmt"
	"github.com/niudaii/filesync/pkg/filesync"
	"github.com/projectdiscovery/gologger"
	"time"
)

func main() {
	host := "127.0.0.1"
	port := "5001"
	auth := "zp857"
	dir := "resource1"
	startFileSyncWorker(host, port, auth, dir)
}

func startFileSyncWorker(host, port, auth, dir string) {
	gologger.Info().Msgf("启动文件同步worker")
	var err error
	for {
		if err = filesync.WorkerStartupSync(host, port, auth, dir); err != nil {
			fmt.Println(err)
		}
		time.Sleep(15 * time.Second)
	}
}