package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func newrootCmd(log DebugLog, out io.Writer, args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "helm-release",
		Short: "helm release helps you manage Helm release objects",
	}

	cmd.AddCommand(newViewCmd(log, args))

	return cmd
}

func debug(format string, v ...interface{}) {
	//if settings.Debug {
	format = fmt.Sprintf("[debug] %s\n", format)
	log.Output(2, fmt.Sprintf(format, v...))
	//}
}

func Execute() {

	rootCmd := newrootCmd(debug, os.Stdout, os.Args[1:])

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
