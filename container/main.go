package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// from talk by Liz Rice - Containers From Scratch: https://www.youtube.com/watch?v=8fi7uSYlOdc
// clone of: 	docker run image <cmd> <params>

func main() {
	switch os.Args[1] {
	case "run":
		run()

	default:
		panic("tell an adult")
	}

}

func run() {
	fmt.Printf("Running %v\n", os.Args[2:])

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// set namespace
	// https://man7.org/linux/man-pages/man7/namespaces.7.html
	// https://itnext.io/chroot-cgroups-and-namespaces-an-overview-37124d995e3d
	// https://stackoverflow.com/questions/46450341/chroot-vs-docker
	// from above.. user namespace (quite new)(2018) which allows a non root user on a host to be mapped with the root user within the container
	// sounds sketchy
	// cookie crumbs: https://unit42.paloaltonetworks.com/breaking-docker-via-runc-explaining-cve-2019-5736/
	// from the source: https://seclists.org/oss-sec/2019/q1/119
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS,
	}

	// in the demo at first stage running hostname and passing single clone flag, did not work. also ubuntu prompts:
	// panic: fork/exec /usr/bin/hostname: operation not permitted
	// ubuntu 20.04 - kernel: Linux wallaby 5.11.0-25-generic #27~20.04.1-Ubuntu SMP Tue Jul 13 17:41:23 UTC 2021 x86_64 x86_64 x86_64 GNU/Linux
	// had to wrap with must() to make work
	must(cmd.Run())

	syscall.Sethostname([]byte("container"))

}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
