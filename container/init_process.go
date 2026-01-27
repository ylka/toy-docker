package container

import (
	"os"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func RunContainerInitProcess(command string, args []string) error {
	log.Infof("command: %s", command)
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
	argv := []string{command}
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		log.Error(err.Error())
	}
	return nil
}
