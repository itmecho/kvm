package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/itmecho/kvm/pkg/download"
	"github.com/itmecho/kvm/pkg/github"
	"github.com/itmecho/kvm/pkg/selector"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(terraformCommand)
}

var terraformCommand = &cobra.Command{
	Use:   "terraform",
	Short: "Manage terraform versions",
	RunE: func(cli *cobra.Command, args []string) error {

		client := &http.Client{}
		// Load releases
		gh := github.New(client)

		// TODO implement filter
		releases, err := gh.GetReleases("hashicorp", "terraform", filter)
		if err != nil {
			return err
		}

		// Make choice
		r := selector.Select("Choose a version: ", releases)

		// Check if file already exists
		// TODO improve this
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		binPath := fmt.Sprintf("%s/bin/terraform-versions/terraform-%s", home, r.Name)
		if _, err := os.Stat(binPath); os.IsNotExist(err) {
			// Download if required
			log.Info("Downloading terraform ", r.Name)
			terraform, err := download.Terraform(&http.Client{}, r.Name)
			if err != nil {
				return err
			}

			f, err := os.Create(binPath)
			if err != nil {
				return err
			}

			_, err = io.Copy(f, terraform)
			if err != nil {
				return err
			}
		}

		// Make executable
		if err = os.Chmod(binPath, 0700); err != nil {
			return err
		}

		// Link
		log.Info("Creating symlink")
		linkPath := fmt.Sprintf("%s/bin/terraform", home)

		if _, err := os.Lstat(linkPath); err == nil {
			if err := os.Remove(linkPath); err != nil {
				return err
			}
		}

		if err = os.Symlink(binPath, linkPath); err != nil {
			return err
		}

		return nil
	},
}
