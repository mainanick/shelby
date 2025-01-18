package main

import (
	"context"
	"fmt"

	containertypes "github.com/docker/docker/api/types/container"
	imagetypes "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	images, err := cli.ImageList(ctx, imagetypes.ListOptions{})
	if err != nil {
		panic(err)
	}

	if len(images) == 0 {
		fmt.Println("No images found")
		return
	}

	for _, image := range images {
		fmt.Println(image.ID)
	}

	containers, err := cli.ContainerList(ctx, containertypes.ListOptions{})
	if err != nil {
		panic(err)
	}

	if len(containers) == 0 {
		fmt.Println("No containers found")
		return
	}

	for _, container := range containers {
		fmt.Println(container.ID)
	}

	err = cli.ContainerStop(ctx, containers[0].ID, containertypes.StopOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Container stopped")

	err = cli.ContainerStart(ctx, containers[0].ID, containertypes.StartOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Container started")

	err = cli.ContainerRestart(ctx, containers[0].ID, containertypes.StopOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Container restarted")

}
