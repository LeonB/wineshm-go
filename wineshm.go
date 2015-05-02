// +build linux

package wineshm

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"

	"golang.org/x/sys/unix"
)

const (
	FILE_MAP_READ  = "r"
	FILE_MAP_WRITE = "w"
	SOCKET_TIMEOUT = 5 * time.Second
)

var (
	WineCmd = []string{"wine"}

	ErrSockTimeout        = errors.New("Timeout reading from unix socket")
	ErrUnexpectedConnType = errors.New("unexpected FileConn type; expected UnixConn")
	ErrTooManyMessages    = errors.New("expected 1 SocketControlMessage")
)

func GetWineShm(shmname string, mode string) (uintptr, error) {
	// Retrieve socket file descriptors
	fds, err := unix.Socketpair(unix.AF_UNIX, unix.SOCK_DGRAM, 0)
	if err != nil {
		return 0, err
	}

	defer unix.Close(fds[0])
	defer unix.Close(fds[1])

	shmwrapper1Path, err := lookPath("shmwrapper1.exe")
	if err != nil {
		return 0, err
	}

	shmwrapper2Path, err := lookPath("shmwrapper2.bin")
	if err != nil {
		return 0, err
	}

	winePath, err := lookPath(WineCmd[0])
	if err != nil {
		return 0, err
	}

	WineCmd[0] = winePath
	args := []string{shmwrapper1Path, shmname, mode, shmwrapper2Path}
	cmd := exec.Command(WineCmd[0], (append(WineCmd, args...))[1:]...)

	writeFile := os.NewFile(uintptr(fds[0]), "child-writes")
	readFile := os.NewFile(uintptr(fds[1]), "parent-reads")
	stderr := &bytes.Buffer{}
	defer writeFile.Close()
	defer readFile.Close()

	// Attach socket to subprocess stdout
	// shmwrapper1 sets the file descriptor as it's stdin (fd0)
	// shmwrapper2 uses stdin (fd0) to get the wine file descriptor
	// and stdout (fd1) as the socket for sending message
	// thats' why the write socket get's connected to the cmd's (shmwrapper1)
	// stdout (fd1)
	cmd.Stdout = writeFile
	cmd.Stderr = stderr

	// Run shwrapper1.exe in wine
	err = cmd.Run()
	if err != nil {
		if len(stderr.Bytes()) > 0 {
			return 0, fmt.Errorf("cmd.Run(): %v (%v)", err, stderr.String())
		} else {
			return 0, fmt.Errorf("cmd.Run(): %v", err)
		}
	}

	// Create a read socket based on the socketpair fd[1]
	c, err := net.FileConn(readFile)
	if err != nil {
		fmt.Errorf("FileConn: %v", err)
	}
	defer c.Close()

	uc, ok := c.(*net.UnixConn)
	if !ok {
		return 0, ErrUnexpectedConnType
	}

	// @TODO: fix this??
	buf := make([]byte, 32) // expect 1 byte
	oob := make([]byte, 32) // expect 24 bytes
	closeUnix := time.AfterFunc(SOCKET_TIMEOUT, func() {
		uc.Close()
	})

	// Retrieve message on socket
	_, oobn, _, _, err := uc.ReadMsgUnix(buf, oob)
	if closeUnix.Stop() == false {
		return 0, ErrSockTimeout
	}

	scms, err := unix.ParseSocketControlMessage(oob[:oobn])
	if err != nil {
		return 0, fmt.Errorf("ParseSocketControlMessage: %v", err)
	}
	if len(scms) != 1 {
		return 0, ErrTooManyMessages
	}

	wineFds, err := unix.ParseUnixRights(&scms[0])
	if err != nil {
		return 0, fmt.Errorf("unix.ParseUnixRights: %v", err)
	}
	if len(wineFds) != 1 {
		return 0, fmt.Errorf("wanted 1 fd; got %#v", wineFds)
	}

	return uintptr(wineFds[0]), nil
}

func lookPath(file string) (string, error) {
	path, err := exec.LookPath("./" + file)
	if err == nil {
		return path, nil
	}

	return exec.LookPath(file)
}
