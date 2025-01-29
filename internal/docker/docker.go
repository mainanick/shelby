package docker

import (
	"github.com/docker/docker/client"
)

var Client *client.Client

func init() {
	var err error
	Client, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
}
