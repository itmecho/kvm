package main

import (
	"log"

	"github.com/spf13/cobra"
)

var filter string

var rootCmd = &cobra.Command{
	Use:   "kvm",
	Short: "Version manager for kubernetes tools",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&filter, "filter", "f", "", "The filter to use when listing releases")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
