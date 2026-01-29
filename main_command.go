package main

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
	"github.com/ylka/toy-docker/cgroups/subsystems"
	"github.com/ylka/toy-docker/container"
)

var runCommand = cli.Command{
	Name: "run",
	Usage: `Create a container with namespace and cgroups limit
			toy-docker run -it [command]`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "it",
			Usage: "enable tty",
		},
		&cli.StringFlag{
			Name:  "mem",
			Usage: "memory limit, e.g.: -mem 100m",
		},
		&cli.IntFlag{
			Name:  "cpu",
			Usage: "cpu quota, e.g.: -cpu 10",
		},
		&cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpuset quota, e.g.: -cpuset 2,4",
		},
		&cli.StringFlag{
			Name:  "v",
			Usage: "volume, e.g.: -v /etc/conf:/etc/conf",
		},
	},
	Action: func(ctx context.Context, c *cli.Command) error {
		if c.Args().Len() < 1 {
			return fmt.Errorf("missing container command")
		}

		var cmdArray []string
		for _, arg := range c.Args().Slice() {
			cmdArray = append(cmdArray, arg)
		}
		tty := c.Bool("it")
		resConf := &subsystems.ResourceConfig{
			MemoryLimit: c.String("mem"),
			CpuCfsQuota: c.Int("cpu"),
			CpuSet:      c.String("cpuset"),
		}

		volume := c.String("v")

		Run(tty, cmdArray, resConf, volume)
		return nil
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",
	Action: func(ctx context.Context, c *cli.Command) error {
		log.Infof("init come on")
		err := container.RunContainerInitProcess()
		return err
	},
}

var commitCommand = cli.Command{
	Name:  "commit",
	Usage: "commit container to image",
	Action: func(ctx context.Context, c *cli.Command) error {
		if c.Args().Len() < 1 {
			return fmt.Errorf("missing image name")
		}

		imageName := c.Args().Get(0)
		container.CommitContainer(imageName)
		return nil
	},
}
