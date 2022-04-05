package main

import (
	"context"
	"fmt"
	"lightweight-cd/internal"

	"github.com/docker/docker/api/types"
)

func main() {

	internal.Info("Starting..")

	internal.ReadConfig()

	ls, err := internal.Docker().ContainerList(context.Background(), types.ContainerListOptions{All: true})

	if err != nil {
		panic(err)
	}

	println("Containers:")
	for i, container := range ls {
		fmt.Println("%v, names: %s, status: %s, state: %s", i+1, container.Names, container.Status, container.State)
	}

	internal.ApiListenAndServe(8080)

}
