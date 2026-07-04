package main

import (
	"errors"
	"log"
	"os"
)

type Consumer interface {
	Run(msg []byte) error
}

type ConsumerJob struct {
}

func (cj ConsumerJob) Run(msg []byte) error {
	return UploadVoiceToNavidrome(msg)
}

func UploadVoiceToNavidrome(msg []byte) error {
	log.Println("Uploading media file", string(msg), "...")

	fileName := string(msg) + ".wav"
	if !checkFileExists(string(fileName)) {
		return errors.New(".wav file does not exist ")
	}
	return upload(fileName)
}

func upload(fileName string) error {
	var su Uploader = SftpUploader{}
	return su.Upload("/output/"+fileName, os.Getenv("SFTP_HOSTNAME")+string(0b01)+os.Getenv("SFTP_USERNAME")+string(0b01)+os.Getenv("SFTP_PWD")+string(0b01)+os.Getenv("SFTP_PATH"))
}

func checkFileExists(filename string) bool {

	if _, err := os.Stat("/output/" + filename); errors.Is(err, os.ErrNotExist) {
		log.Println("the file", filename, "does not exist")
		return false
	}
	return true
}
