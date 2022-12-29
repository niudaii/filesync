package cmd

import (
	"fmt"
	"github.com/niudaii/filesync/pkg/filesync"
	"github.com/projectdiscovery/gologger"
	"github.com/spf13/cobra"
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

func StartFileSyncWorker(host, port, auth, dir string, timer int) {
	gologger.Info().Msgf("启动文件同步worker")
	var err error
	if timer == 0 {
		if err = filesync.WorkerStartupSync(host, port, auth, dir); err != nil {
			fmt.Println(err)
		}
	} else {
		for {
			if err = filesync.WorkerStartupSync(host, port, auth, dir); err != nil {
				fmt.Println(err)
			}
			time.Sleep(time.Duration(timer) * time.Second)
		}
	}
}
