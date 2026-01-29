package container

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/ylka/toy-docker/constant"
)

func NewWorkSpace(rootPath, volume string) {
	createLower(rootPath)
	createDirs(rootPath)
	mountOverlayFS(rootPath)

	if volume != "" {
		mntPath := path.Join(rootPath, "merged")
		hostPath, containerPath, err := volumeExtract(volume)
		if err != nil {
			log.Errorf("volume invalid %v", err)
			return
		}

		mountVolume(mntPath, hostPath, containerPath)
	}
}

func mountVolume(mntPath, hostPath, containerPath string) {
	if err := os.Mkdir(hostPath, constant.Perm0777); err != nil {
		log.Infof("mkdir host dir [%s] failed. %v", containerPath, err)
	}

	containerPathInHost := path.Join(mntPath, containerPath)
	if err := os.Mkdir(containerPathInHost, constant.Perm0777); err != nil {
		log.Infof("mkdir container dir [%s] failed. %v", containerPath, err)
	}

	cmd := exec.Command("mount", "-o", "bind", hostPath, containerPathInHost)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Errorf("mount volume failed. %v", err)
	}
}

func volumeExtract(volume string) (string, string, error) {
	parts := strings.Split(volume, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid volum [%s], must split by `:`", volume)
	}
	sourcePath, destPath := parts[0], parts[1]
	if sourcePath == "" || destPath == "" {
		return "", "", fmt.Errorf("invalid volume [%s], path can't be empty ", volume)
	}

	return sourcePath, destPath, nil
}

func mountOverlayFS(rootPath string) {
	dirs := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", path.Join(rootPath, "busybox"),
		path.Join(rootPath, "upper"), path.Join(rootPath, "work"))
	cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", dirs, path.Join(rootPath, "merged"))
	log.Infof("mount overlayfs: [%s]", cmd.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("overlayfs failed: %v", err)
	}
}

func createDirs(rootPath string) {
	dirs := []string{
		path.Join(rootPath, "upper"),
		path.Join(rootPath, "work"),
		path.Join(rootPath, "merged"),
	}
	for _, dir := range dirs {
		if err := os.Mkdir(dir, constant.Perm0777); err != nil {
			log.Infof("make dir %s fail, %v ", dir, err)
		}
	}
}

func createLower(rootPath string) {
	busyboxPath := path.Join(rootPath, "busybox")
	busyboxTarPath := path.Join(rootPath, "busybox.tar")
	exist, err := PathExists(busyboxPath)
	if err != nil {
		log.Infof("failed to judge busybox path %s . %v", busyboxPath, err)
	}

	if !exist {
		if err := os.Mkdir(busyboxPath, constant.Perm0777); err != nil {
			log.Infof("make busybox dir fail, %v ", err)
			return
		}
		if _, err := exec.Command("tar", "-xvf", busyboxTarPath, "-C", busyboxPath).CombinedOutput(); err != nil {
			log.Infof("untar busybox fail, %v ", err)
		}
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func DeleteWorkSpace(rootPath, volume string) {
	mntPath := path.Join(rootPath, "merged")

	if volume != "" {
		_, containerPath, err := volumeExtract(volume)
		if err != nil {
			log.Errorf("volume invalid %v", err)
			return
		}

		umountVolume(mntPath, containerPath)
	}

	unmountOverlayFS(mntPath)
	deleteDirs(rootPath)
}

func umountVolume(mntPath, containerPath string) {
	containerPathInHost := path.Join(mntPath, containerPath)
	cmd := exec.Command("umount", containerPathInHost)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Errorf("umount volume failed. %v", err)
	}
}

func deleteDirs(rootPath string) {
	dirs := []string{
		path.Join(rootPath, "upper"),
		path.Join(rootPath, "work"),
		path.Join(rootPath, "merged"),
	}
	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			log.Infof("remove dir %s fail, %v ", dir, err)
		}

		log.Infof("remove dir %s.", dir)
	}
}

func unmountOverlayFS(mntPath string) {
	cmd := exec.Command("umount", mntPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("%v", err)
	}
}

func NewParentProcess(tty bool, volume string) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		log.Errorf("New pipe error %v", err)
		return nil, nil
	}

	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUTS | syscall.CLONE_NEWNET |
			syscall.CLONE_NEWPID | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	// cmd 执行时就会外带着这个文件句柄去创建子进程
	cmd.ExtraFiles = []*os.File{readPipe}

	rootPath := constant.RootPath
	NewWorkSpace(rootPath, volume)
	cmd.Dir = path.Join(rootPath, "merged")

	return cmd, writePipe
}
