package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("tell an adult")
	}

}

func run() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	// wire up stdin, out and err so we can see stuff when we run the cmd
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// set namespace
	// https://man7.org/linux/man-pages/man7/namespaces.7.html
	// https://itnext.io/chroot-cgroups-and-namespaces-an-overview-37124d995e3d
	// https://stackoverflow.com/questions/46450341/chroot-vs-docker
	// from above.. user namespace (quite new)(2018) which allows a non root user on a host to be mapped with the root user within the container
	// cookie crumbs: https://unit42.paloaltonetworks.com/breaking-docker-via-runc-explaining-cve-2019-5736/
	// from the source: https://seclists.org/oss-sec/2019/q1/119
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWUSER | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		UidMappings: []syscall.SysProcIDMap{{
			ContainerID: 0,
			HostID:      1000,
			Size:        1,
		}},
		Unshareflags: syscall.CLONE_NEWNS,
	}
	// CLONE_NEWUTS = hostname  |  CLONE_NEWPID = new namespace for pids  |  CLONE_NEWNS = new namespace
	// systemd recursively shares mounts with all other namespaces

	// syscall.CLONE_NEWUSER = create new user namespace. must have UidMappings below as well.

	must(cmd.Run())
}

func child() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	cg()

	syscall.Sethostname([]byte("failwhale"))
	syscall.Chroot("/home/brandon/tmp_fs")
	syscall.Chdir("/")
	syscall.Mount("proc", "proc", "proc", 0, "")

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(cmd.Run())

	syscall.Unmount("/proc", 0)
}

func cg() {
	cgroups := "/sys/fs/cgroup"
	pids := filepath.Join(cgroups, "pids")
	err := os.Mkdir(filepath.Join(pids, "failwhale"), 755)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
	must(ioutil.WriteFile(filepath.Join(pids, "failwhale/pids.max"), []byte("20"), 0700))
	// removes new cgroup after container exits
	must(ioutil.WriteFile(filepath.Join(pids, "failwhale/notify_on_release"), []byte("1"), 0700))
	must(ioutil.WriteFile(filepath.Join(pids, "failwhale/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))

}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
