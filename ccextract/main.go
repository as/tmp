/*
	docker build . -t ccextractor # if you dont have it on local machine
	go build && ./ccextractor -i `aws s3 presign url`
*/
package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	input = flag.String("i", "", "input file to download")
)

func main() {
	flag.Parse()

	resp, err := http.Get(*input)
	ck("download", err)

	text, _ := ccextract(resp.Body)
	io.Copy(os.Stdout, text)
}

func ck(topic string, err error) {
	if err != nil {
		log.Fatalf("%s: %v", topic, err)
	}
}
