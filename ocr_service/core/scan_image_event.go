package core

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ahmed-debbech/tevox/ocr_service/config"
	"github.com/ahmed-debbech/tevox/ocr_service/model"
	"github.com/ahmed-debbech/tevox/ocr_service/queues"
	"github.com/go-cmd/cmd"
)

func ProcessScanImageEvent(msg []byte) error {

	var request model.ScanImageEventRequest
	err := json.Unmarshal(msg, &request)
	if err != nil {
		return errors.New("could not parse json for this request " + err.Error())
	}
	if !sanitize(request) {
		return errors.New("the request doesn't look good " + err.Error())
	}

	if !checkFileExists(request.FileName) {
		return errors.New("image file does not exist")
	}

	textFileName, err := launchTesseract(request.FileName)
	if err != nil {
		return errors.New("something failed with tesseract :" + err.Error())
	}

	response := model.ProcessTextToVoiceRequest{
		TextFileName: textFileName,
	}
	responseBytes, _ := json.Marshal(response)
	if err := queues.Publish(config.QueuesNames["B_QUEUE"].Name, []byte(responseBytes), nil); err != nil {
		log.Println(err)
	}
	log.Println("Processed image successfully into text, pushed to queue")
	return nil
}

func checkFileExists(filename string) bool {

	if _, err := os.Stat("/pic/" + filename); errors.Is(err, os.ErrNotExist) {
		log.Println("the file", filename, "does not exist")
		return false
	}
	return true
}

func sanitize(request model.ScanImageEventRequest) bool {
	if request.FileName == "" {
		log.Println("FileName can not be empty")
		return false
	}
	return true
}
func launchTesseract(fileName string) (string, error) {

	defer log.SetPrefix("")
	log.SetPrefix("[TESSERACT-0] ")
	// Start a long-running process, capture stdout and stderr
	textFileName := "out_" + strconv.Itoa(int(time.Now().Unix()))
	findCmd := cmd.NewCmd("sh", "-c", "tesseract --tessdata-dir ../tesseract /pic/"+fileName+" - -l eng > /text/"+textFileName)
	statusChan := findCmd.Start() // non-blocking
	log.Println("Started tesseract")
	ticker := time.NewTicker(2 * time.Second)

	// Print last line of stdout every 2s
	go func() {
		for range ticker.C {
			status := findCmd.Status()
			if status.Exit == 0 {
				n := len(status.Stdout)
				if n > 0 {
					log.Println("[STDIN]", status.Stdout[n-1])
				}
			} else {
				ticker.Stop()
				break
			}
		}
	}()

	// Check if command is done
	select {
	case finalStatus := <-statusChan:
		if finalStatus.Exit == 0 {
			log.Println("SUCCESSFUL!")
		} else {
			log.Println("[STDERR]", strings.Join(finalStatus.Stderr, "\n[STDERR]")) //make this line out put to error file
			log.Println("ERROR! Something went wrong when running tesseract, exit code:", finalStatus.Exit)
			return "", errors.New("ERROR! Something went wrong when running tesseract")
		}
	case <-time.After(10 * time.Second):
		err := findCmd.Stop()
		if err != nil {
			log.Println(err)
		}
		log.Println("[TIMEOUT] tevox killed tesseract")
	}

	return textFileName, nil
}
