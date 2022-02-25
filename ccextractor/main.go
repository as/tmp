/*
	docker build . -t ccextractor # if you dont have it on local machine
	go run main.go -i `aws s3 presign url`
*/
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"golang.org/x/sys/unix"
)

var (
	input = flag.String("i", "", "input file to download")
)

func main() {
	flag.Parse()

	resp, err := http.Get(*input)
	ck("download", err)

	// NOTE(as):
	// Here we are create an in-memory *os.File
	// which we copy the response body to. The
	// ccextractor is told that the input file is /proc/self/fd/3
	// because we add an extra file descriptor to be inherited
	// by the child process. The ccextractor opens its own file descriptor
	// and then treats it like a regular file with seek/read operations.
	fd, _ := memfd()
	defer fd.Close()
	io.Copy(fd, resp.Body)
	resp.Body.Close()

	cmd := exec.Command("ccextractor", "-mp4", "/proc/self/fd/3", "-stdout")
	cmd.Stdout = os.Stdout // can be an *os.File or any io.Writer

	// This is where our child process gets our file descriptor access
	//
	// fd0: stdin
	// fd1: stdout
	// fd2: stderr
	// fd3: our memfd file
	//
	// In the Linux /proc/ filesystem, /proc/self/fd/3 is a reference to that
	cmd.ExtraFiles = append(cmd.ExtraFiles, fd)

	err = cmd.Run()
	ck("run", err)
}

// memfd creates an in-memory file that can be passed
// into a child process for reading or writing.
func memfd() (*os.File, error) {
	name := "ccextractor"
	fd, err := unix.MemfdCreate(name, os.O_RDWR)
	if err != nil {
		return nil, fmt.Errorf("memfd: %v", err)
	}
	return os.NewFile(uintptr(fd), name), nil
}

func ck(topic string, err error) {
	if err != nil {
		log.Fatalf("%s: %v", topic, err)
	}
}
