// +build !linux

package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

func ccextract(src io.Reader) (text *bytes.Buffer, err error) {
	file, err := ioutil.TempFile("", "ccextract*")
	if err != nil {
		return
	}
	defer os.Remove(file.Name())
	defer file.Close()
	io.Copy(file, src)

	cmd := exec.Command("ccextractor", "-mp4", file.Name(), "-stdout")
	text = new(bytes.Buffer)
	cmd.Stdout = text
	return text, cmd.Run()
}
