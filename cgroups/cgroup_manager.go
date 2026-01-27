package cgroups

import (
	"github.com/sirupsen/logrus"
	"github.com/ylka/toy-docker/cgroups/subsystems"
)

type CgroupManager struct {
	Path     string
	Resource *subsystems.ResourceConfig
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path: path,
	}
}

func (c *CgroupManager) Apply(pid int, res *subsystems.ResourceConfig) error {
	for _, sub := range subsystems.SubsystemsInstances {
		err := sub.Apply(c.Path, pid, res)
		if err != nil {
			logrus.Infof("apply subsystem: %s %v", sub.Name(), err)
		}
	}
	return nil
}

func (c *CgroupManager) Set(res *subsystems.ResourceConfig) error {
	for _, sub := range subsystems.SubsystemsInstances {
		err := sub.Set(c.Path, res)
		if err != nil {
			logrus.Infof("set subsystem: %s %v", sub.Name(), err)
		}
	}
	return nil
}

func (c *CgroupManager) Destroy() error {
	for _, sub := range subsystems.SubsystemsInstances {
		err := sub.Remove(c.Path)
		if err != nil {
			logrus.Infof("remove subsystem: %s %v", sub.Name(), err)
		}
	}
	return nil
}
