package subsystems

import (
	"fmt"
	"os"
	"path"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/ylka/toy-docker/constant"
)

type CpuSubSystem struct{}

const (
	PeriodDefault = 100000
)

func (*CpuSubSystem) Name() string {
	return "cpu"
}

func (s *CpuSubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	if res.CpuCfsQuota == 0 {
		log.Info("CpuCfsQuota is 0")
		return nil
	}

	subsysCgroupPath, err := getCgroupPath(cgroupPath, true)
	if err != nil {
		return err
	}

	if err := os.WriteFile(
		path.Join(subsysCgroupPath, "cgroup.subtree_control"),
		[]byte("+cpu"),
		constant.Perm0644); err != nil {
		return fmt.Errorf("set cgroup +cpu fail %v", err)
	}

	if err := os.WriteFile(path.Join(subsysCgroupPath, "cpu.max"),
		fmt.Appendf(nil, "%d %d", res.CpuCfsQuota*1000, PeriodDefault),
		constant.Perm0644); err != nil {
		return fmt.Errorf("set cgroup cpu max fail %v", err)
	}

	return nil
}

func (s *CpuSubSystem) Apply(cgroupPath string, pid int, res *ResourceConfig) error {
	if res.CpuCfsQuota == 0 {
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

func (s *CpuSubSystem) Remove(cgroupPath string) error {
	subsysCgroupPath, err := getCgroupPath(cgroupPath, false)
	if err != nil {
		return err
	}

	return os.RemoveAll(subsysCgroupPath)
}
