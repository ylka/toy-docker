package subsystems

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/ylka/toy-docker/constant"
)

type MemorySubsystem struct {
}

func (*MemorySubsystem) Name() string {
	return "memory"
}

func (s *MemorySubsystem) Set(cgroupPath string, res *ResourceConfig) error {
	if res.MemoryLimit == "" {
		return nil
	}

	subsysCgroupPath, err := getCgroupPath(cgroupPath, true)
	if err != nil {
		return errors.Wrapf(err, "get cgroup %s", cgroupPath)
	}

	if err := os.WriteFile(
		path.Join(subsysCgroupPath, "cgroup.subtree_control"),
		[]byte("+memory"),
		constant.Perm0644); err != nil {
		return fmt.Errorf("set cgroup +memory fail %v", err)
	}

	if err := os.WriteFile(path.Join(subsysCgroupPath, "memory.max"), []byte(res.MemoryLimit),
		constant.Perm0644); err != nil {
		log.Errorf("set cgroup memory fail %v", err)
		return err
	}
	return nil
}

func (s *MemorySubsystem) Apply(cgroupPath string, pid int, res *ResourceConfig) error {
	if res.MemoryLimit == "" {
		return nil
	}
	subsysCgroupPath, err := getCgroupPath(cgroupPath, false)
	if err != nil {
		return errors.Wrapf(err, "get cgroup %s", cgroupPath)
	}

	if err := os.WriteFile(path.Join(subsysCgroupPath, constant.CgroupProc), []byte(strconv.Itoa(pid)),
		constant.Perm0644); err != nil {
		log.Errorf("set cgroup proc fail %v", err)
		return err
	}

	return nil
}

func (s *MemorySubsystem) Remove(cgroupPath string) error {
	subsysCgroupPath, err := getCgroupPath(cgroupPath, false)
	if err != nil {
		return err
	}

	return os.RemoveAll(subsysCgroupPath)
}
