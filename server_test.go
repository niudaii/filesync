package main

import (
	"fmt"
	"github.com/niudaii/filesync/pkg/filesync"
	"github.com/projectdiscovery/gologger"
	"testing"
)

func TestServer(t *testing.T) {
	host := "0.0.0.0"
	port := "5001"
	auth := "zp857"
	dir := "resource"
	blackList := []string{"pocscan"}
	startFileSyncServer(host, port, auth, dir, blackList)
}

func startFileSyncServer(host, port, auth, dir string, blackList []string) {
	gologger.Info().Msgf("启动文件同步server: %v", fmt.Sprintf("%v:%v", host, port))
	filesync.StartFileSyncServer(host, port, auth, dir, blackList)
}
