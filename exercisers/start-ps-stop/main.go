package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/strslice"
)

func stopper(ctx context.Context, docker *client.Client) {
	options := types.ContainerListOptions{All: true}
	for range time.Tick(300 * time.Millisecond) {
		start := time.Now()
		cs, err := docker.ContainerList(ctx, options)
		if err != nil {
			log.Printf("[ERROR] List: %v\n", err)
		}
		log.Printf("PS time: %v\n", time.Since(start))
		for _, c := range cs {
			if start.Unix()-c.Created > 10 {
				rmStart := time.Now()
				rOpts := types.ContainerRemoveOptions{
					ContainerID:   c.ID,
					Force:         true,
					RemoveVolumes: true,
				}
				if err := docker.ContainerRemove(ctx, rOpts); err != nil {
					log.Printf("[ERROR] Remove: %v\n", err)
				}
				log.Printf("Removal: %v", time.Since(rmStart))
			}
		}
	}
}

func starter(ctx context.Context, docker *client.Client) {
	for range time.Tick(500 * time.Millisecond) {
		// Create a container
		containerConfig := &container.Config{
			Image:       "ruby:latest",
			Cmd:         strslice.StrSlice([]string{"irb"}),
			AttachStdin: true,
			Tty:         true,
		}
		hostConfig := &container.HostConfig{
			Resources: container.Resources{
				Memory: 200 * 1024 * 1024,
			},
		}
		resp, err := docker.ContainerCreate(ctx, containerConfig, hostConfig, nil, "")
		if err != nil {
			log.Fatal(err)
		}
		aOpts := types.ContainerAttachOptions{
			ContainerID: resp.ID,
			Stdin:       true,
			Stdout:      true,
			Stderr:      true,
			Stream:      true,
		}
		rd, err := docker.ContainerAttach(ctx, aOpts)
		if err != nil {
			log.Fatal(err)
		}
		defer rd.Conn.Close()

		go io.Copy(ioutil.Discard, rd.Reader)

		// Start the container
		err = docker.ContainerStart(ctx, resp.ID)
		if err != nil {
			log.Printf("Start: %v\n", err)
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: %s <target>", os.Args[0])
		os.Exit(1)
	}

	address := os.Args[1]

	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	docker, err := client.NewClient(address, "", nil, defaultHeaders)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	go stopper(ctx, docker)
	starter(ctx, docker)
}
