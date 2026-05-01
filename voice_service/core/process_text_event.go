package core

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ahmed-debbech/tevox/voice_service/config"
	"github.com/ahmed-debbech/tevox/voice_service/model"
	"github.com/ahmed-debbech/tevox/voice_service/queues"
	"github.com/go-cmd/cmd"
)

func ProcessTextToVoiceEvent(msg []byte) error {

	var request model.ProcessTextToVoiceRequest
	err := json.Unmarshal(msg, &request)
	if err != nil {
		return errors.New("could not parse json for this request " + err.Error())
	}

	if !sanitize(request) {
		return errors.New("the request doesn't look good " + err.Error())
	}

	if !checkFileExists(request.TextFileName) {
		return errors.New("image file does not exist ")
	}

	if err := launchPiper(request.TextFileName); err != nil {
		return errors.New("something failed with tesseract :" + err.Error())
	}
	if err := queues.Publish(config.QueuesNames["B_QUEUE"].Name, []byte("done"), nil); err != nil {
		log.Println(err)
	}
	log.Println("Processed text to voice, pushed to queue.")
	return nil
}

func launchPiper(fileName string) error {
	defer log.SetPrefix("")
	log.SetPrefix("[PIPER-0] ")
	// Start a long-running process, capture stdout and stderr
	vocalPath := os.Getenv("VOICE_LANGUAGE") + "-" + os.Getenv("VOICE_NAME") + "-" + os.Getenv("VOICE_QUALITY")
	findCmd := cmd.NewCmd("sh", "-c", "piper -m /v/"+vocalPath+".onnx --input-file /input/"+fileName+" --output_file /output/hello.wav")
	statusChan := findCmd.Start() // non-blocking
	log.Println("Started piper")
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
			log.Println("ERROR! Something went wrong when running piper, exit code:", finalStatus.Exit)
			return errors.New("ERROR! Something went wrong when running piper")
		}
	case <-time.After(10 * time.Second):
		err := findCmd.Stop()
		if err != nil {
			log.Println(err)
		}
		log.Println("[TIMEOUT] tevox killed piper")
	}

	return nil
}

func sanitize(request model.ProcessTextToVoiceRequest) bool {
	if request.TextFileName == "" {
		log.Println("FileName can not be empty")
		return false
	}
	return true
}

func checkFileExists(filename string) bool {

	if _, err := os.Stat("/input/" + filename); errors.Is(err, os.ErrNotExist) {
		log.Println("the file", filename, "does not exist")
		return false
	}
	return true
}
