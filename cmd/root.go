package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var (
	host string
	port string
	auth string
	dir  string
)

var rootCmd = &cobra.Command{
	Use:               "filesync",
	Short:             "文件同步工具 by zp857",
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if host == "" {
			log.Fatal("请输入 host")
		}
		if auth == "" {
			log.Fatal("请输入 auth")
		}
	},
}

func Execute() {
	rootCmd.PersistentFlags().StringVar(&host, "host", "", "host")
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", "5001", "port")
	rootCmd.PersistentFlags().StringVar(&auth, "auth", "zp857", "auth")
	rootCmd.PersistentFlags().StringVar(&dir, "dir", "./", "dir to file sync")
	cobra.CheckErr(rootCmd.Execute())
}
