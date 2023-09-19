package types

// WizMessageParams represents the params field inside the message.
// Using interface{} instead structs as devices answer with errors if we send more fields than needed for selected method
// So it forces us to construct it dynamically
type WizMessageParams map[string]interface{}

// Here it is the complete structure we should follow in the requests for guidance
//type WizMessageParams struct {
//	// Fields
//	State bool `json:"state,omitempty"` // State is status of the Device : true if ON, false if OFF
//
//	// Temperature of the light
//	C    int `json:"c,omitempty"`    // C set the value of the cold white led (0-255)
//	W    int `json:"w,omitempty"`    // W set the value of the warm white led (0-255)
//	Temp int `json:"temp,omitempty"` // Temp set the color temperature for the white led in the bulb (2000-9000)
//
//	// Set the color in RGB
//	R int `json:"r,omitempty"` // (0-255)
//	G int `json:"g,omitempty"` // (0-255)
//	B int `json:"b,omitempty"` // (0-255)
//
//	//
//	Dimming int `json:"dimming,omitempty"` // Dimming set the value of the brightness (10-100)
//	Ratio   int `json:"ratio,omitempty"`   // Ratio set the ratio between the up and down light (1-100)
//	Speed   int `json:"speed,omitempty"`   // Speed set the effect changing speed (10-200)
//
//	SceneId    int `json:"sceneId,omitempty"`    // SceneId set the scene by id (1-32)
//	SchdPsetId int `json:"schdPsetId,omitempty"` // SchdPsetId set the rhythm of the room by id // TODO: figure out how to send
//
//  // Fields on 'pulse' command
//	Delta    int `json:"delta,omitempty"`    // Amount of light decreased on the pulse (-100 ··· -1)
//  Duration int `json:"duration,omitempty"` // Number of milliseconds that pulse lasts

//}

// WizMessageResult represent the result field inside the message's response
type WizMessageResult struct {
	// Almost globally present
	Mac     string `json:"mac,omitempty"` // Mac from the bulb
	Src     string `json:"src,omitempty"` // Src of the state change
	Success bool   `json:"success,omitempty"`

	// Fields on getPilot
	Rssi       int  `json:"rssi,omitempty"`       // Rssi is WiFi signal strength, in negative dBm
	State      bool `json:"state,omitempty"`      // State is status of the Device : true if ON, false if OFF
	SceneId    int  `json:"sceneId,omitempty"`    // SceneId represents current scene by id
	SchdPsetId int  `json:"schdPsetId,omitempty"` // SchdPsetId represents current rhythm by id (another type of scene)
	R          int  `json:"r,omitempty"`          // (0-255)
	G          int  `json:"g,omitempty"`          // (0-255)
	B          int  `json:"b,omitempty"`          // (0-255)
	C          int  `json:"c,omitempty"`          // C represents the value of the cold white led (0-255)
	W          int  `json:"w,omitempty"`          // W represents the value of the warm white led (0-255)
	Dimming    int  `json:"dimming,omitempty"`    // Dimming set the value of the brightness (10-100)

	// Fields on getDevInfo
	DevMac string `json:"devMac,omitempty"`

	// Fields on getUserConfig
	FadeIn     int  `json:"fadeIn,omitempty"`
	FadeOut    int  `json:"fadeOut,omitempty"`
	DftDim     int  `json:"dftDim,omitempty"`
	OpMode     int  `json:"opMode,omitempty"`
	Po         bool `json:"po,omitempty"`
	MinDimming int  `json:"minDimming,omitempty"`
	TapSensor  int  `json:"tapSensor,omitempty"`

	// Fields on getSystemConfig
	HomeId     int    `json:"homeId,omitempty"`
	RoomId     int    `json:"roomId,omitempty"`
	Rgn        string `json:"rgn,omitempty"` // Rgn represents the region in the world. I.E: 'eu'
	ModuleName string `json:"moduleName,omitempty"`
	FwVersion  string `json:"fwVersion,omitempty"`
	GroupId    int    `json:"groupId,omitempty"`
	Ping       int    `json:"ping,omitempty"`
	DrvConf    []int  `json:"drvConf,omitempty"`

	// Fields on getModelConfig
	Ps           int   `json:"ps,omitempty"`
	PwmFreq      int   `json:"pwmFreq,omitempty"`
	PwmRange     []int `json:"pwmRange,omitempty"`
	Wcr          int   `json:"wcr,omitempty"`
	Nowc         int   `json:"nowc,omitempty"`
	CctRange     []int `json:"cctRange,omitempty"` // CctRange represents the white temperature range advertised to the user. Extended whiteRange property (new)
	ExtRange     []int `json:"extRange,omitempty"` // ExtRange represents the white temperature range advertised to the user. Extended whiteRange property (old)
	RenderFactor []int `json:"renderFactor,omitempty"`
	Wizc1        struct {
		Mode []int `json:"mode,omitempty"`
	} `json:"wizc1,omitempty"`
	Wizc2 struct {
		Mode []int `json:"mode,omitempty"`
	} `json:"wizc2,omitempty"`

	// TODO OTHER
	WhiteRange float64 `json:"whiteRange,omitempty"` // WhiteRange white temperature range supported by the light // TODO (where?)
	Temp       int     `json:"temp,omitempty"`       // Temp set the color temperature for the white led in the bulb // TODO (where?)
}

// WizMessageError represent an error message received from the device
type WizMessageError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// WizMessage represent a message sent to the device
type WizMessage struct {
	Method string `json:"method"`
	Id     int    `json:"id,omitempty"`

	Params WizMessageParams `json:"params,omitempty"`
}

// WizMessageResponse represent a message received from the device
type WizMessageResponse struct {
	Method string `json:"method"`
	Id     int    `json:"id,omitempty"`
	Env    string `json:"env,omitempty"`

	Result WizMessageResult `json:"result,omitempty"`
	Error  WizMessageError  `json:"error,omitempty"`
}
