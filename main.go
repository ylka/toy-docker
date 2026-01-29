package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {

	app := &cli.Command{}
	app.Name = "toy-docker"
	app.Usage = "my toy docker cli"

	app.Commands = []*cli.Command{
		&initCommand,
		&runCommand,
		&commitCommand,
	}

	app.Before = func(ctx context.Context, c *cli.Command) (context.Context, error) {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(os.Stdout)
		return ctx, nil
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

// func main() {
// 	if os.Args[0] == "/proc/self/exe" {
// 		fmt.Printf("current pid %d", syscall.Getpid())
// 		fmt.Println()
// 		cmd := exec.Command("sh", "-c", `stress --vm-bytes 200m --vm-keep -m 1`)
// 		cmd.SysProcAttr = &syscall.SysProcAttr{}
// 		cmd.Stdin = os.Stdin
// 		cmd.Stdout = os.Stdout
// 		cmd.Stderr = os.Stderr

// 		if err := cmd.Run(); err != nil {
// 			log.Fatal(err)
// 			os.Exit(1)
// 		}
// 	}
// 	cmd := exec.Command("/proc/self/exe")
// 	cmd.SysProcAttr = &syscall.SysProcAttr{
// 		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS}
// 	cmd.Stdin = os.Stdin
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr

// 	if err := cmd.Run(); err != nil {
// 		log.Fatal(err)
// 		os.Exit(1)
// 	} else {
// 		fmt.Printf("%v", cmd.Process.Pid)
// 		const cgroupMemoryHierarchyMount = "/sys/fs/cgroup/memory"
// 		os.Mkdir(path.Join(cgroupMemoryHierarchyMount, "testmemorylimit"), 0755)
// 		os.WriteFile(path.Join(cgroupMemoryHierarchyMount, "testmemorylimit", "tasks"), []byte(strconv.Itoa(cmd.Process.Pid)), 0644)
// 		os.WriteFile(path.Join(cgroupMemoryHierarchyMount, "testmemorylimit", "memory.limit_in_bytes"), []byte("100m"), 0644)
// 	}
// 	cmd.Process.Wait()
// }
