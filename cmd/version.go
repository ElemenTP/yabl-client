/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	Version   = "unknown version"
	BuildTime = "unknown time"
)

//Add version command into the cli app, will show version info and compile time.
func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("yabl-client %s %s %s with %s %s\n", Version, runtime.GOOS, runtime.GOARCH, runtime.Version(), BuildTime)
		},
	})
}
