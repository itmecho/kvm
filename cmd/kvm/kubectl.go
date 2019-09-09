package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/itmecho/kvm/pkg/github"
	"github.com/itmecho/kvm/pkg/selector"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(kubectlCommand)
}

var kubectlCommand = &cobra.Command{
	Use:   "kubectl",
	Short: "Manage kubectl versions",
	RunE: func(cli *cobra.Command, args []string) error {

		client := &http.Client{}
		// Load releases
		gh := github.New(client)

		// TODO implement filter
		releases, err := gh.GetReleases("kubernetes", "kubernetes", filter)
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
		binPath := fmt.Sprintf("%s/bin/kubectl-versions/kubectl-%s", home, r.Name)
		if _, err := os.Stat(binPath); os.IsNotExist(err) {
			// Download if required
			log.Info("Downloading kops ", r.Name)
			url := fmt.Sprintf("https://dl.k8s.io/%s/kubernetes-client-linux-amd64.tar.gz", r.Name)
			resp, err := client.Get(url)
			if err != nil {
				return err
			}
			// TODO don't do this!
			defer resp.Body.Close()

			tarArchive, err := gzip.NewReader(resp.Body)
			if err != nil {
				return err
			}

			data := tar.NewReader(tarArchive)

			f, err := os.Create(binPath)
			if err != nil {
				return err
			}

			for {
				hdr, err := data.Next()
				if err == io.EOF {
					break
				}

				if strings.HasSuffix(hdr.Name, "kubectl") {
					_, err = io.Copy(f, data)
					if err != nil {
						return err
					}
					break
				}
			}
		}

		// Make executable
		if err = os.Chmod(binPath, 0700); err != nil {
			return err
		}

		// Link
		log.Info("Creating symlink")
		linkPath := fmt.Sprintf("%s/bin/kubectl", home)

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
