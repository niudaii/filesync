package cmd

import (
	"github.com/fsnotify/fsnotify"
	"github.com/niudaii/filesync/pkg/filesync"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

var (
	blackList []string
	monitor   bool
)

func init() {
	serverCmd.Flags().StringSliceVar(&blackList, "black", []string{}, "black list to file sync")
	serverCmd.Flags().BoolVar(&monitor, "monitor", false, "open file monitor")
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "文件同步服务端",
	Run: func(cmd *cobra.Command, args []string) {
		if monitor {
			go StartFileMonitor(dir)
		}
		StartFileSyncServer(host, port, auth, dir, blackList)
	},
}

func StartFileSyncServer(host, port, auth, dir string, blackList []string) {
	log.Printf("启动文件同步server: %v:%v", host, port)
	filesync.StartFileSyncServer(host, port, auth, dir, blackList)
}

func StartFileMonitor(dir string) {
	log.Print("启动文件监控")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	// 开启事件监听
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
	// 监控所有目录
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			path, err = filepath.Abs(path)
			if err != nil {
				return err
			}
			err = watcher.Add(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	<-make(chan struct{})
}
