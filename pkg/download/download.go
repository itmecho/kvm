package download

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
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
