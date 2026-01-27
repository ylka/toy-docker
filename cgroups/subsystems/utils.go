package subsystems

import (
	"github.com/pkg/errors"
	"github.com/ylka/toy-docker/constant"

	"os"
	"path"
)

func getCgroupPath(cgroupPath string, autoCreate bool) (string, error) {
	absPath := path.Join(constant.SubsysCgroupPath, cgroupPath)
	if !autoCreate {
		return absPath, nil
	}
	_, err := os.Stat(absPath)
	if err != nil && os.IsNotExist(err) {
		err := os.Mkdir(absPath, constant.Perm0755)
		return absPath, err
	}
	return absPath, errors.Wrap(err, "create cgroup")
}
