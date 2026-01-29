package container

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

func RunContainerInitProcess() error {
	cmdArray := readUserCommand()
	if len(cmdArray) == 0 {
		return errors.New("run container get user command error, cmdArray is nil")
	}

	setupMount()

	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		log.Errorf("Exec look path error:%v", err)
		return err
	}

	log.Infof("Find path %s", path)

	if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		log.Errorf("RunContainerInitProcess exec :%v", err)
	}
	return nil
}

const fdIndex = 3

func readUserCommand() []string {
	pipe := os.NewFile(uintptr(fdIndex), "pipe")
	msg, err := io.ReadAll(pipe)
	if err != nil {
		log.Errorf("init read pipe error %v", err)
		return nil
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}

func setupMount() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Errorf("Get current location error :%v", err)
		return
	}
	log.Infof("current location is %s", pwd)

	err = syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	if err != nil {
		log.Errorf("init mount rootfs error %v", err)
		return
	}

	err = pivotRoot(pwd)
	if err != nil {
		log.Errorf("pivotRoot failed, detail: %v", err)
		return
	}

	// 2️⃣ 确保 /proc 存在
	if err := os.MkdirAll("/proc", 0555); err != nil {
		log.Errorf("/proc mkdir failed, detail: %v", err)
		return
	}

	// 3️⃣ mount proc
	flags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	if err := syscall.Mount("proc", "/proc", "proc", uintptr(flags), ""); err != nil {
		log.Errorf("mount /proc failed, detail: %v", err)
		return
	}

	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")

}

func pivotRoot(root string) error {
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return errors.Wrap(err, "mount rootfs to itself")
	}

	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return errors.Wrap(err, "make .pivot_root dir fail")
	}

	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return errors.WithMessagef(err, "pivotRoot failed, new root: %v old put: %v", root, pivotDir)
	}

	if err := syscall.Chdir("/"); err != nil {
		return errors.WithMessage(err, "Chdir to / failed")
	}

	pivotDir = filepath.Join("/", ".pivot_root")
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return errors.WithMessage(err, "unmount pivot root dir")
	}

	return os.Remove(pivotDir)
}
