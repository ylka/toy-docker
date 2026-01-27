package container

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func RunContainerInitProcess() error {
	mountProc()

	cmdArray := readUserCommand()
	if len(cmdArray) == 0 {
		return errors.New("run container get user command error, cmdArray is nil")
	}

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

func mountProc() error {
	// 1️⃣ 切断 mount 传播（必须）
	if err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, ""); err != nil {
		return err
	}

	// 2️⃣ 确保 /proc 存在
	if err := os.MkdirAll("/proc", 0555); err != nil {
		return err
	}

	// 3️⃣ mount proc
	flags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	if err := syscall.Mount("proc", "/proc", "proc", uintptr(flags), ""); err != nil {
		return err
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
