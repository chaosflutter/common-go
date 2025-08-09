package tts

// TTSResponse represents the response for TTS requests
type TTSResponseData struct {
	FilePath string `json:"file_path"`
	Code     int    `json:"code"`
}

// TTSErrorResponse represents error responses
type TTSErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Details string `json:"details,omitempty"`
}

// TTSRequest represents the request payload for TTS
type TTSRequest struct {
	App struct {
		Appid   string `json:"appid"`
		Token   string `json:"token"`
		Cluster string `json:"cluster"`
	} `json:"app"`
	User struct {
		UID string `json:"uid"`
	} `json:"user"`
	Audio struct {
		VoiceType  string  `json:"voice_type"`
		Encoding   string  `json:"encoding"`
		SpeedRatio float64 `json:"speed_ratio"`
	} `json:"audio"`
	Request struct {
		ReqID     string `json:"reqid"`
		Text      string `json:"text"`
		Operation string `json:"operation"`
	} `json:"request"`
}

// TTSResponse represents the response from TTS API
type TTSResponse struct {
	Reqid   string `json:"reqid"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}
