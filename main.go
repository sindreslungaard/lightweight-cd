package main

import (
	"context"
	"fmt"
	"lightweight-cd/internal"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var docker = getClientFunc()

func main() {

	internal.Info("Starting..")

	internal.ReadConfig()

	ls, err := docker().ContainerList(context.Background(), types.ContainerListOptions{All: true})

	if err != nil {
		panic(err)
	}

	println("Containers:")
	for i, container := range ls {
		fmt.Println("%v, names: %s, status: %s, state: %s", i+1, container.Names, container.Status, container.State)
	}

	internal.ApiListenAndServe(8080)

}

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
