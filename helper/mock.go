package helper

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const containerName = "wiremock"
const mockServerPort = "8080"

func StartMockServer() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		panic(err)
	}

	imageName := "rodolpheche/wiremock"
	_, err = cli.ImagePull(ctx, imageName, types.ImagePullOptions{})

	if err != nil {
		panic(err)
	}

	// check if the container is already created,
	_, err = cli.ContainerInspect(ctx, containerName)
	if err != nil && !strings.Contains(err.Error(), "Conflict") {
		// container already exits
		log.Printf("container %s does not exist, creating", containerName)

		hostBinding := nat.PortBinding{
			HostIP:   "0.0.0.0",
			HostPort: mockServerPort,
		}

		containerPort, err := nat.NewPort("tcp", "8080")
		if err != nil {
			panic("Unable to get the port")
		}

		portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}

		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		_, err = cli.ContainerCreate(ctx, &container.Config{
			Image: imageName,
			Cmd:   []string{"--global-response-templating", "--local-response-templating"},
		}, &container.HostConfig{
			AutoRemove:   true,
			PortBindings: portBinding,
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: fmt.Sprintf("%s/mocks", path.Dir(path.Dir(dir))), // FIXME: read the directory from  the environment?
					Target: "/home/wiremock",
				},
			},
		}, nil, nil, containerName)
		if err != nil {
			panic(err)
		}
	}

	if err := cli.ContainerStart(ctx, containerName, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)
	log.Printf("%s container started\n", containerName)
}

func StopMockServer() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		panic(err)
	}

	removeOptions := types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}

	// try to force remove the container
	if err := cli.ContainerRemove(ctx, containerName, removeOptions); err != nil {
		log.Fatalf("Unable to remove container: %s", err)
	}

	log.Printf("%s container removed\n", containerName)
}
