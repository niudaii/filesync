package cmd

import (
	"github.com/niudaii/filesync/pkg/filesync"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

var (
	timer int
)

func init() {
	workerCmd.Flags().IntVar(&timer, "timer", 0, "worker sync timer, timer in seconds")
	rootCmd.AddCommand(workerCmd)
}

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "文件同步客户端",
	Run: func(cmd *cobra.Command, args []string) {
		StartFileSyncWorker(host, port, auth, dir, timer)
	},
}

func pathExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}
	return info.IsDir()
}

func StartFileSyncWorker(host, port, auth, dir string, timer int) {
	log.Println("启动文件同步worker")
	if ok := pathExists(dir); !ok { // 判断是否有Director文件夹
		_ = os.Mkdir(dir, os.ModePerm)
	}
	for {
		if err := filesync.WorkerStartupSync(host, port, auth, dir); err != nil {
			log.Println(err)
		}
		if timer == 0 {
			break
		}
		time.Sleep(time.Duration(timer) * time.Second)
	}
}
