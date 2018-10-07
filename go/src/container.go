package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	switch os.Args[1] {
		case "run":
			run()
		case "child":
			child()
		default:
			panic("Arguement not recognized")
	}
}

func run() {
	// /proc/self/exe is a psuedo fork() exec()
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Requires root priviledges due to these flags requiring SYS_CAP_ADMIN
	// Creates new UTS and PID namespaces
	cmd.SysProcAttr = &syscall.SysProcAttr {
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,// UTS Unix Timesharing System, ie hostname
	}															// NEWPID creates new PID namespace

	syserr(cmd.Run())
}

func child() {
	fmt.Printf("Running %v as PID %d\n", os.Args[2:], os.Getpid())

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// More namespaces, limiting what we can see
	// Change root of filesystem
	syserr(syscall.Chroot("/home/cpalv/workspace/containers/rootfs"))
	syserr(os.Chdir("/"))// cd into new root
	syserr(syscall.Mount("proc", "proc", "proc", 0, ""))// mount our own proc filesystem
	syserr(cmd.Run())
}

func syserr( err error) {
	if err != nil {
		panic(err)
	}
}
