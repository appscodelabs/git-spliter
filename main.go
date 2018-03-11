package main

import (
	"flag"
	"fmt"
	"log"

	"encoding/json"

	"github.com/appscode/go/log/golog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/tamalsaha/go-oneliners"
)

func NewRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Short:             `git-spliter`,
		DisableAutoGenTag: true,
		PersistentPreRun: func(c *cobra.Command, args []string) {
			c.Flags().VisitAll(func(flag *pflag.Flag) {
				log.Printf("FLAG: --%s=%q", flag.Name, flag.Value)
			})

			opt := golog.ParseFlags(c.Flags())
			b, _ := json.MarshalIndent(opt, "", "  ")
			oneliners.FILE()
			fmt.Println(string(b))
		},
	}
	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	// ref: https://github.com/kubernetes/kubernetes/issues/17162#issuecomment-225596212
	flag.CommandLine.Parse([]string{})

	rootCmd.AddCommand(NewCmdInit())
	return rootCmd
}

func main() {
	if err := NewRootCmd().Execute(); err != nil {
		log.Fatalln("Error in Stash Main:", err)
	}
}
