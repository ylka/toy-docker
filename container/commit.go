package container

import (
	"os/exec"
	"path"

	log "github.com/sirupsen/logrus"
	"github.com/ylka/toy-docker/constant"
)

func CommitContainer(imageName string) {
	mntPath := path.Join(constant.RootPath, "merged")
	iamgeTar := path.Join(constant.RootPath, imageName+".tar")
	log.Infof("commit container iamge tar %s", iamgeTar)
	if _, err := exec.Command("tar", "-czf", iamgeTar, "-C", mntPath, ".").CombinedOutput(); err != nil {
		log.Errorf("tar folder %s err %v", mntPath, err)
	}
}
