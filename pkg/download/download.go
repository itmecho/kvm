package download

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// Kops returns a reader containing the file contents of the given kops version
func Kops(client *http.Client, version string) (io.Reader, error) {
	url := fmt.Sprintf("https://github.com/kubernetes/kops/releases/download/%s/kops-linux-amd64", version)
	resp, err := client.Get(url)
	return resp.Body, err
}

// Kubectl returns a reader containing the file contents of the given kubectl version
func Kubectl(client *http.Client, version string) (io.Reader, error) {
	url := fmt.Sprintf("https://dl.k8s.io/%s/kubernetes-client-linux-amd64.tar.gz", version)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	tarArchive, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}

	data := tar.NewReader(tarArchive)

	for {
		hdr, err := data.Next()
		if err == io.EOF {
			break
		}

		if strings.HasSuffix(hdr.Name, "kubectl") {
			return data, nil
		}
	}

	return nil, fmt.Errorf("Something went wrong downloading kubectl: kubectl not found in downloaded archive")
}

// Terraform returns a reader containing the file contents of the given terraform version
func Terraform(client *http.Client, version string) (io.Reader, error) {
	version = strings.TrimPrefix(version, "v")
	url := fmt.Sprintf("https://releases.hashicorp.com/terraform/%[1]s/terraform_%[1]s_linux_amd64.zip", version)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	zipData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(zipData), resp.ContentLength)
	if err != nil {
		return nil, err
	}

	for _, zippedFile := range zipReader.File {
		if zippedFile.Name == "terraform" {
			terraformFile, err := zippedFile.Open()
			if err != nil {
				return nil, err
			}

			return terraformFile, nil
		}
	}

	return nil, fmt.Errorf("terraform not found in downloaded zip archive")
}
