package subsystems

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/ylka/toy-docker/constant"
)

type CpusetSubSystem struct{}

func (*CpusetSubSystem) Name() string {
	return "cpuset"
}

func (s *CpusetSubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	if res.CpuSet == "" {
		return nil
	}

	subsysCgroupPath, err := getCgroupPath(cgroupPath, true)
	if err != nil {
		return err
	}
	// 必须先写 mems
	if err := os.WriteFile(path.Join(subsysCgroupPath, "cpuset.mems"), []byte("0"),
		constant.Perm0644); err != nil {
		return fmt.Errorf("set cgroup cpuset mems fail %v", err)
	}

	if err := os.WriteFile(path.Join(subsysCgroupPath, "cpuset.cpus"), []byte(res.CpuSet),
		constant.Perm0644); err != nil {
		return fmt.Errorf("set cgroup cpuset cpus fail %v", err)
	}
	return nil
}

func (s *CpusetSubSystem) Apply(cgroupPath string, pid int, res *ResourceConfig) error {
	if res.CpuSet == "" {
		return nil
	}

	subsysCgroupPath, err := getCgroupPath(cgroupPath, false)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path.Join(subsysCgroupPath, constant.CgroupProc), []byte(strconv.Itoa(pid)),
		constant.Perm0644); err != nil {
		return fmt.Errorf("set cgroup proc fail %v", err)
	}

	return nil
}

func (s *CpusetSubSystem) Remove(cgroupPath string) error {
	subsysCgroupPath, err := getCgroupPath(cgroupPath, false)
	if err != nil {
		return err
	}

	return os.RemoveAll(subsysCgroupPath)
}
