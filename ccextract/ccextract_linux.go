package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"golang.org/x/sys/unix"
)

func ccextract(src io.Reader) (text *bytes.Buffer, err error) {
	text = new(bytes.Buffer)

	// NOTE(as):
	// Here we are create an in-memory *os.File
	// which we copy the response body to. The
	// ccextractor is told that the input file is /proc/self/fd/3
	// because we add an extra file descriptor to be inherited
	// by the child process. The ccextractor opens its own file descriptor
	// and then treats it like a regular file with seek/read operations.
	//
	// This is where our child process gets our file descriptor access
	//
	// fd0: stdin
	// fd1: stdout
	// fd2: stderr
	// fd3: our memfd file
	//
	// In the Linux /proc/ filesystem, /proc/self/fd/3 is a reference to that
	fd, _ := memfd()
	defer fd.Close()
	io.Copy(fd, src)

	cmd := exec.Command("ccextractor", "-mp4", "/proc/self/fd/3", "-stdout")
	cmd.Stdout = text
	cmd.ExtraFiles = append(cmd.ExtraFiles, fd)
	return text, cmd.Run()
}

// memfd creates an in-memory file that can be passed
// into a child process for reading or writing.
func memfd() (*os.File, error) {
	name := "ccextractor" // this does not have to be unique per-process
	fd, err := unix.MemfdCreate(name, os.O_RDWR)
	if err != nil {
		return nil, fmt.Errorf("memfd: %v", err)
	}
	return os.NewFile(uintptr(fd), name), nil
}
