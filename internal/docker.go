package internal

import (
	"context"

	"github.com/docker/docker/client"
)

var Docker = getClientFunc()

func getClientFunc() func() *client.Client {

	var docker *client.Client

	return func() *client.Client {

		if docker != nil {

			_, err := docker.Ping(context.Background())

			// current client ok, use it
			if err == nil {
				return docker
			}

			// not ok, close and create new client
			docker.Close()
			docker = nil

		}

		cli, err := client.NewClientWithOpts(client.FromEnv)

		if err != nil {
			panic(err)
		}

		docker = cli

		println("New docker client connection")

		return docker

	}

}
