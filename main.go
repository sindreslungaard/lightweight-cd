package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var docker = getClientFunc()

func main() {

	println("Starting..")

	ReadConfig()

	ls, err := docker().ContainerList(context.Background(), types.ContainerListOptions{All: true})

	if err != nil {
		panic(err)
	}

	println("Containers:")
	for i, container := range ls {
		fmt.Printf("%v, %s", i+1, container.Names)
	}

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
