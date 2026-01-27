package main

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
	"github.com/ylka/toy-docker/container"
)

var initCommand = cli.Command{
	Name: "run",
	Usage: `Create a container with namespace and cgroups limit
			toy-docker run -it [command]`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "it",
			Usage: "enable tty",
		},
	},
	Action: func(ctx context.Context, c *cli.Command) error {
		if c.Args().Len() < 1 {
			return fmt.Errorf("missing container command")
		}

		cmd := c.Args().Get(0)
		tty := c.Bool("it")

		Run(tty, cmd)
		return nil
	},
}

var runCommand = cli.Command{
	Name:  "init",
	Usage: "init something",
	Action: func(ctx context.Context, c *cli.Command) error {
		log.Infof("init come on")
		cmd := c.Args().Get(0)
		log.Infof("command %s", cmd)
		err := container.RunContainerInitProcess(cmd, nil)
		return err
	},
}
