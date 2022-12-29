package cmd

import (
	"fmt"
	"github.com/niudaii/filesync/pkg/filesync"
	"github.com/projectdiscovery/gologger"
	"github.com/spf13/cobra"
)

var (
	blackList []string
)

func init() {
	serverCmd.Flags().StringSliceVar(&blackList, "black", []string{}, "black list to file sync")
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "文件同步服务端",
	Run: func(cmd *cobra.Command, args []string) {
		StartFileSyncServer(host, port, auth, dir, blackList)
	},
}

func StartFileSyncServer(host, port, auth, dir string, blackList []string) {
	gologger.Info().Msgf("启动文件同步server: %v", fmt.Sprintf("%v:%v", host, port))
	filesync.StartFileSyncServer(host, port, auth, dir, blackList)
}
