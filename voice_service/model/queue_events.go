package model

type ProcessTextToVoiceRequest struct {
	Title        string `json:"title"`
	TextFileName string `json:"textFileName"`
}
