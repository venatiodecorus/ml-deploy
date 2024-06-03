package utils

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

// func dockerValidate(config string) bool {
// 	return true
// }
func createTarball(dockerfileContent string) (io.Reader, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	// Add Dockerfile to the tarball
	dockerfile := tar.Header{
		Name: "Dockerfile",
		Mode: 0600,
		Size: int64(len(dockerfileContent)),
	}
	if err := tw.WriteHeader(&dockerfile); err != nil {
		return nil, err
	}
	if _, err := tw.Write([]byte(dockerfileContent)); err != nil {
		return nil, err
	}

	// Close the tar writer
	if err := tw.Close(); err != nil {
		return nil, err
	}

	return buf, nil
}

// Docker doesn't have a validation command, so we try to build and iterate on any errors
func DockerBuild(config string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	cli.NegotiateAPIVersion(context.Background())

	// Create a tarball from the Dockerfile content
	tarball, err := createTarball(config)
	if err != nil {
		return err
	}

	// Build the Docker image
	options := types.ImageBuildOptions{
		Tags: []string{"mldeploy:latest"},
		Remove: true,
		Dockerfile: "Dockerfile",
	}

	resp, err := cli.ImageBuild(context.Background(), tarball, options)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(os.Stdout, resp.Body)
	return err
}

func DockerList() ([]image.Summary, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	cli.NegotiateAPIVersion(context.Background())

	resp, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return nil,err
	}

	// json, err := json.Marshal(resp)
	// if err != nil {
	// 	return nil,err
	// }

	return resp, nil
}