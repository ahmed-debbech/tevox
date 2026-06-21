package main

import (
	"log"
	"os"
	"strings"

	sftp "github.com/moov-io/go-sftp"
)

type Uploader interface {
	Upload(filePath string, destination string) error
}

type SftpUploader struct {
}

func (su SftpUploader) Upload(filepath string, destination string) error {

	log.Println("uploading file to sftp server")
	p := strings.Split(destination, string(0b01))
	clientConfig := &sftp.ClientConfig{
		Hostname: p[0],
		Username: p[1],
		Password: p[2],
	}

	client, err := sftp.NewClient(nil, clientConfig)
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(); err != nil {
		log.Fatal(err)
	}

	if client != nil {
		defer client.Close()
	}

	log.Println("connected to SFTP server successfully!")
	fileData, err := os.Open(filepath)

	err = client.UploadFile(p[3], fileData)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("file " + filepath + " uploaded to server " + p[3] + " successfully!")
	return nil
}
