package internal

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
)

func RunContainerFromDeployment(client *client.Client, deployment Deployment) {

	c, ok := IsDeploymentRunning(client, deployment.UID)

	if ok && c.State == "running" {
		return
	}

	stopTimeout := time.Second * 5

	err := StopAndRemoveDeployment(client, deployment.UID, &stopTimeout, false)

	if err != nil {
		Warn(err)
	}

	_, err = client.ImagePull(context.Background(), "docker.io/library/alpine", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	resp, err := client.ContainerCreate(context.Background(), &container.Config{
		Image: "alpine",
		Cmd:   []string{"echo", "hello world"},
		Tty:   false,
	}, nil, nil, nil, "")

	if err != nil {
		Fatal(err)
	}

	if err := client.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		Fatal(err)
	}

}

func IsDeploymentRunning(client *client.Client, uid string) (*types.Container, bool) {

	containers, err := client.ContainerList(context.Background(), types.ContainerListOptions{All: true})

	if err != nil {
		ClientConnectionError(err)
	}

	for _, c := range containers {
		for _, name := range c.Names {
			if name == uid {
				return &c, true
			}
		}
	}

	return nil, false

}

func StopAndRemoveDeployment(client *client.Client, uid string, stopTimeout *time.Duration, forceRemove bool) error {

	c, ok := IsDeploymentRunning(client, uid)

	if !ok {
		return nil
	}

	err := client.ContainerStop(context.Background(), c.ID, stopTimeout)

	if err != nil {
		return err
	}

	Info("Stopped container", uid)

	err = client.ContainerRemove(context.Background(), uid, types.ContainerRemoveOptions{Force: forceRemove})

	if err != nil {
		return err
	}

	Info("Removed container", uid)

	return nil

}

func ClientConnectionError(err error) {
	Fatal("Client connection error")
}

func ContainerRegistryLogin(client *client.Client, auth types.AuthConfig) (registry.AuthenticateOKBody, error) {
	return client.RegistryLogin(context.Background(), auth)
}
