package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sts"
)

var dst = flag.String("dst", "", "s3 path where to copy canary file in the format s3://bucket/path/test.txt")

func main() {
	flag.Parse()
	if *dst == "" {
		fmt.Println("-dst unset")
		os.Exit(1)
	}
	region := "us-east-1"
	s, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: &region},
	})
	ck("newsession", err)

	id, err := sts.New(s).GetCallerIdentity(&sts.GetCallerIdentityInput{})
	ck("get caller id", err)
	fmt.Println("caller identity", id)

	uri := parseURL(*dst)
	_, err = s3manager.NewUploader(s).Upload(&s3manager.UploadInput{
		Bucket: &uri.Host,
		Key:    &uri.Path,
		Body:   strings.NewReader(time.Now().String()),
	})
	ck("upload file", err)

	fmt.Println("success")
}

func parseURL(s string) *url.URL {
	u, _ := url.Parse(s)
	if u == nil {
		return &url.URL{}
	}
	return u
}

func ck(what string, err error) {
	if err != nil {
		panic(fmt.Sprintf("%s: %v", what, err))

	}
}
