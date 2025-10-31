package cmd

import (
	"fmt"
	"os"
	"wb_lvl2_telnet/internal"

	"github.com/spf13/cobra"
)

var (
	host    string
	port    string
	timeout int
)

var telnetCmd = &cobra.Command{
	Use: "telnet",
	RunE: func(cmd *cobra.Command, args []string) error {

		err := internal.Connect(host, port, timeout)
		if err != nil {
			_, err = fmt.Fprintln(os.Stdout, err.Error())
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(telnetCmd)

	telnetCmd.Flags().StringVarP(&host, "host", "s", "", "адрес TCP-сервера")
	telnetCmd.Flags().StringVarP(&port, "port", "p", "", "порт")
	telnetCmd.Flags().IntVarP(&timeout, "timeout", "t", 10, "таймаут (по умолчанию 10 секунд)")
}
