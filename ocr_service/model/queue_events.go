package model

type ScanImageEventRequest struct {
	Title      string   `json:"title"`
	FromJid    string   `json:"fromJid"`
	ImagePaths []string `json:"imagePaths"`
}

type ProcessTextToVoiceRequest struct {
	TextFileName string `json:"textFileName"`
}
