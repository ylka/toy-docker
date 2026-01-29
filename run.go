package main

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/ylka/toy-docker/cgroups"
	"github.com/ylka/toy-docker/cgroups/subsystems"
	"github.com/ylka/toy-docker/constant"
	"github.com/ylka/toy-docker/container"
)

func Run(tty bool, cmdArray []string, res *subsystems.ResourceConfig, volume string) {
	parent, writePipe := container.NewParentProcess(tty, volume)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Errorf("Run parent.Start err:%v", err)
	}

	cgroupManager := cgroups.NewCgroupManager("toy-docker")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid, res)

	sendInitCommand(cmdArray, writePipe)

	_ = parent.Wait()

	container.DeleteWorkSpace(constant.RootPath, volume)
}

func sendInitCommand(cmdArray []string, writePipe *os.File) {
	command := strings.Join(cmdArray, " ")
	log.Infof("command all is %s", command)
	_, err := writePipe.WriteString(command)
	if err != nil {
		log.Errorf("Write pipe error:%v", err)
	}
	writePipe.Close()
}
