package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/exp/maps"
	"log"
	"net"
	"slices"
	"strconv"
	"strings"
	"time"
)

var (
	// TODO
	WizScenes = map[int]string{
		1:    "Ocean",
		2:    "Romance",
		3:    "Sunset",
		4:    "Party",
		5:    "Fireplace",
		6:    "Cozy",
		7:    "Forest",
		8:    "Pastel Colors",
		9:    "Wake up",
		10:   "Bedtime",
		11:   "Warm White",
		12:   "Daylight",
		13:   "Cool white",
		14:   "Night light",
		15:   "Focus",
		16:   "Relax",
		17:   "True colors",
		18:   "TV time",
		19:   "Plantgrowth",
		20:   "Spring",
		21:   "Summer",
		22:   "Fall",
		23:   "Deepdive",
		24:   "Jungle",
		25:   "Mojito",
		26:   "Club",
		27:   "Christmas",
		28:   "Halloween",
		29:   "Candlelight",
		30:   "Golden white",
		31:   "Pulse",
		32:   "Steampunk",
		1000: "Rhythm",
	}

	// TW - have Cool White and Warm White LEDs. Such devices support most static light modes + CCT control
	WizTwScenes = []int{6, 9, 10, 11, 12, 13, 14, 15, 16, 18, 29, 30, 31, 32}

	// DW - have only Dimmable white LEDs. Such devices support only dimming, Wake up, Bedtime and Night light modes.
	WizDwScenes = []int{9, 10, 13, 14, 29, 30, 31, 32}
)

const (
	ErrorReceivingResponseErrorMessage = "error receiving response: %s"
	ErrorSendingDataErrorMessage       = "error sending data: %s"
)

// Ref: https://aleksandr.rogozin.us/blog/2021/8/13/hacking-philips-wiz-lights-via-command-line
// Ref: https://github.com/sbidy/pywizlight

// TODO
type WizClient struct {
	deviceConnection *net.UDPConn
}

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
//	SceneId int `json:"sceneId,omitempty"` // SceneId set the scene by id (1-32)
//
//	SchdPsetId int `json:"schdPsetId,omitempty"` // rythm ID of the room
//}

// WizMessageResult represent the result field inside the message's response
type WizMessageResult struct {
	// Almost globally present
	Mac     string `json:"mac,omitempty"` // Mac from the bulb
	Src     string `json:"src,omitempty"` // Src of the state change
	Success bool   `json:"success,omitempty"`

	// Fields on getPilot
	Rssi    int  `json:"rssi,omitempty"`    // Rssi is WiFi signal strength, in negative dBm
	State   bool `json:"state,omitempty"`   // State is status of the Device : true if ON, false if OFF
	SceneId int  `json:"sceneId,omitempty"` // SceneId set the scene by id
	R       int  `json:"r,omitempty"`       // (0-255)
	G       int  `json:"g,omitempty"`       // (0-255)
	B       int  `json:"b,omitempty"`       // (0-255)
	C       int  `json:"c,omitempty"`       // C represents the value of the cold white led (0-255)
	W       int  `json:"w,omitempty"`       // W represents the value of the warm white led (0-255)
	Dimming int  `json:"dimming,omitempty"` // Dimming set the value of the brightness (10-100)

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
	Rgn        string `json:"rgn,omitempty"`
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
	SchdPsetId int     `json:"schdPsetId,omitempty"` // TODO (where?)
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

// TODO Migrate UDPProxy cache to global cache
func createWizClient(host string, port int) (wizClient WizClient, err error) {

	// Resolve the address for the given backend
	address, err := net.ResolveUDPAddr("udp", net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		return wizClient, err
	}

	// Open a connection with a remote server
	deviceConn, err := net.DialUDP("udp", nil, address)
	if err != nil {
		return wizClient, err
	}
	//defer deviceConn.Close()

	wizClient.deviceConnection = deviceConn
	return wizClient, err
}

// In general, send datagrams to a device
func (w *WizClient) sendDatagrams(content []byte) (response []byte, err error) {

	// Send datagrams to the device
	_, err = w.deviceConnection.Write(content)
	if err != nil {
		err = errors.New(fmt.Sprintf(ErrorSendingDataErrorMessage, err))
		return response, err
	}

	// Prepare a buffer to receive response
	buffer := make([]byte, 1024)

	// Wait to receive the answer from device
	n, addr, err := w.deviceConnection.ReadFromUDP(buffer)
	if err != nil {
		err = errors.New(fmt.Sprintf(ErrorReceivingResponseErrorMessage, err))
		return response, err
	}

	response = buffer[:n]
	_ = addr

	return response, err
}

// TODO
func (w *WizClient) sendMessage(message WizMessage) (response WizMessageResponse, err error) {

	// Lo pasamos a json
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return response, err
	}

	// TODO
	//jsonMessageFiltered, err := removeEmptyValuesFromJson(jsonMessage)
	//if err != nil {
	//	return response, err
	//}

	log.Printf("el mensaje que vamos a mandar %v", string(jsonMessage))

	// TODO
	responseBytes, err := w.sendDatagrams(jsonMessage)
	if err != nil {
		return response, err
	}

	log.Print("CP1")

	err = json.Unmarshal(responseBytes, &response)
	return response, err
}

func (w *WizClient) getPilot() (response WizMessageResponse, err error) {

	wizMessage := WizMessage{
		Id:     1,
		Method: "getPilot",
	}

	response, err = w.sendMessage(wizMessage)
	return response, err
}

// TODO
func (w *WizClient) getSystemConfig() (response WizMessageResponse, err error) {

	wizMessage := WizMessage{
		Id:     1,
		Method: "getSystemConfig",
	}

	response, err = w.sendMessage(wizMessage)
	return response, err
}

// IsRgb return true when the device have RGB, cool white and warm white LEDs.
// These type of bulbs can support all the light modes provided by WiZ
func (w *WizClient) IsRgb() (bool, error) {

	configResp, err := w.getSystemConfig()
	if err != nil {
		return false, errors.New(fmt.Sprintf("Problemas pillando la config: %v", err))
	}

	if !strings.Contains(configResp.Result.ModuleName, "RGB") {
		return false, nil
	}

	return true, nil
}

// IsTw return true when the device have cool white and warm white LEDs.
// These type of devices support most static light modes + CCT control
func (w *WizClient) IsTw() (bool, error) {

	configResp, err := w.getSystemConfig()
	if err != nil {
		return false, errors.New(fmt.Sprintf("Problemas pillando la config: %v", err))
	}

	if !strings.Contains(configResp.Result.ModuleName, "TW") {
		return false, nil
	}

	return true, nil
}

// IsDw return true when the device have only dimmable white LEDs.
// These type of devices support only dimming and some light modes
func (w *WizClient) IsDw() (bool, error) {

	configResp, err := w.getSystemConfig()
	if err != nil {
		return false, errors.New(fmt.Sprintf("Problemas pillando la config: %v", err))
	}

	if !strings.Contains(configResp.Result.ModuleName, "DW") {
		return false, nil
	}

	return true, nil
}

// TurnOn turns on the device
func (w *WizClient) TurnOn() (response WizMessageResponse, err error) {

	wizMessage := WizMessage{
		Id:     1,
		Method: "setState",
		Params: WizMessageParams{
			"state": true,
		},
	}

	response, err = w.sendMessage(wizMessage)
	return response, err
}

// TurnOff turns off the device
func (w *WizClient) TurnOff() (response WizMessageResponse, err error) {

	wizMessage := WizMessage{
		Id:     1,
		Method: "setState",
		Params: WizMessageParams{
			"state": false,
		},
	}

	response, err = w.sendMessage(wizMessage)
	return response, err
}

// SetBrightness change the device's brightness (10-100)
func (w *WizClient) SetBrightness(brightness int) (response WizMessageResponse, err error) {

	if brightness < 10 || brightness > 100 {
		return response, errors.New("el brillo tiene que estar entre 10 y 100")
	}

	wizMessage := WizMessage{
		Id:     1,
		Method: "setPilot",
		Params: WizMessageParams{
			"dimming": brightness,
		},
	}

	response, err = w.sendMessage(wizMessage)
	return response, err
}

// SetRgb set the color for the device (3 x 0-255)
func (w *WizClient) SetRgb(r, g, b int) (response WizMessageResponse, err error) {

	if r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
		return response, errors.New("los colores van desde 0 a 255")
	}

	wizMessage := WizMessage{
		Id:     1,
		Method: "setPilot",
		Params: WizMessageParams{
			"r": r,
			"g": g,
			"b": b,
		},
	}

	response, err = w.sendMessage(wizMessage)
	return response, err
}

// TODO
func (w *WizClient) SetColdWhite(coldWhite int) (response WizMessageResponse, err error) {

	if coldWhite < 0 || coldWhite > 255 {
		return response, errors.New("los coldWhite van desde 0 a 255")
	}

	wizMessage := WizMessage{
		Id:     1,
		Method: "setPilot",
		Params: WizMessageParams{
			"c": coldWhite,
		},
	}

	response, err = w.sendMessage(wizMessage)
	return response, err
}

// TODO
func (w *WizClient) SetWarmWhite(warmWhite int) (response WizMessageResponse, err error) {

	if warmWhite < 0 || warmWhite > 255 {
		return response, errors.New("los warmWhite van desde 0 a 255")
	}

	wizMessage := WizMessage{
		Id:     1,
		Method: "setPilot",
		Params: WizMessageParams{
			"w": warmWhite,
		},
	}

	response, err = w.sendMessage(wizMessage)
	return response, err
}

// TODO
func (w *WizClient) SetTemperature(temperature int) (response WizMessageResponse, err error) {

	if temperature < 2000 || temperature > 9000 {
		return response, errors.New("la temperatura va en kelvin desde 2000 a 9000")
	}

	wizMessage := WizMessage{
		Id:     1,
		Method: "setPilot",
		Params: WizMessageParams{
			"temp": temperature,
		},
	}

	response, err = w.sendMessage(wizMessage)
	return response, err
}

// SetSpeed set changing speed between the colors in a scene
func (w *WizClient) SetSpeed(speed int) (response WizMessageResponse, err error) {

	if speed < 10 || speed > 200 {
		return response, errors.New("la velocidad va en desde 10 a 200")
	}

	wizMessage := WizMessage{
		Id:     1,
		Method: "setPilot",
		Params: WizMessageParams{
			"speed": speed,
		},
	}

	response, err = w.sendMessage(wizMessage)
	return response, err
}

// TODO
func (w *WizClient) SetRatio(ratio int) (response WizMessageResponse, err error) {

	if ratio < 1 || ratio > 100 {
		return response, errors.New("la velocidad del cambio entre luz alta y baja va en desde 1 a 100")
	}

	wizMessage := WizMessage{
		Id:     1,
		Method: "setPilot",
		Params: WizMessageParams{
			"ratio": ratio,
		},
	}

	response, err = w.sendMessage(wizMessage)
	return response, err
}

// SetScene set a scene by its ID. The availability of IDs depend on the type of bulb: RBG, TW, DW
func (w *WizClient) SetScene(sceneId int) (response WizMessageResponse, err error) {

	isRgbDevice, err := w.IsRgb()
	if err != nil {
		return response, errors.New("error determinando el tipo de dispositivo")
	}

	if isRgbDevice && !slices.Contains(maps.Keys(WizScenes), sceneId) {
		return response, errors.New("pediste una escena sólo para rgb")
	}

	isTwDevice, err := w.IsTw()
	if err != nil {
		return response, errors.New("error determinando el tipo de dispositivo")
	}

	if isTwDevice && !slices.Contains(WizTwScenes, sceneId) {
		return response, errors.New("pediste una escena sólo para TW")
	}

	isDwDevice, err := w.IsDw()
	if err != nil {
		return response, errors.New("error determinando el tipo de dispositivo")
	}

	if isDwDevice && !slices.Contains(WizDwScenes, sceneId) {
		return response, errors.New("pediste una escena sólo para DW")
	}

	wizMessage := WizMessage{
		Id:     1,
		Method: "setPilot",
		Params: WizMessageParams{
			"sceneId": sceneId,
		},
	}

	response, err = w.sendMessage(wizMessage)
	return response, err
}

func main() {
	wizClient, err := createWizClient("192.168.2.107", 38899)

	response, err := wizClient.TurnOn()
	if err != nil {
		log.Fatalf("error encendiendo: %s", err)
	}
	log.Print(response)

	// ---
	time.Sleep(5 * time.Second)
	// ---

	log.Printf("seteando la escena 1")
	sceneResp, err := wizClient.SetScene(31)
	if err != nil {
		log.Fatalf("error poniendo la escena: %s", err)
	}
	log.Print(sceneResp)

	// ---
	time.Sleep(3 * time.Second)
	// ---

	log.Printf("seteando la velocidad 1")
	spResp, err := wizClient.SetSpeed(200)
	if err != nil {
		log.Fatalf("error poniendo la veloc.: %s", err)
	}
	log.Print(spResp)

	// ---
	time.Sleep(3 * time.Second)
	// ---

	log.Printf("seteando la velocidad 2")
	spResp, err = wizClient.SetSpeed(10)
	if err != nil {
		log.Fatalf("error poniendo la veloc.: %s", err)
	}
	log.Print(spResp)

	// ---
	time.Sleep(3 * time.Second)
	// ---

	log.Printf("seteando la ratio 1")
	rtResp, err := wizClient.SetRatio(100)
	if err != nil {
		log.Fatalf("error poniendo la ratio.: %s", err)
	}
	log.Print(rtResp)

	// ---
	time.Sleep(3 * time.Second)
	// ---

	log.Printf("seteando la ratio 2")
	rtResp, err = wizClient.SetRatio(1)
	if err != nil {
		log.Fatalf("error poniendo la ratio.: %s", err)
	}
	log.Print(rtResp)

	// ---
	time.Sleep(15 * time.Second)
	// ---

	log.Printf("seteando la temperatura 1")
	tempResp, err := wizClient.SetTemperature(9000)
	if err != nil {
		log.Fatalf("error poniendo la temperatura: %s", err)
	}
	log.Print(tempResp)

	// ---
	time.Sleep(5 * time.Second)
	// ---

	log.Printf("seteando la temperatura 2")
	tempResp, err = wizClient.SetTemperature(2000)
	if err != nil {
		log.Fatalf("error poniendo la temperatura: %s", err)
	}
	log.Print(tempResp)

	// ---
	time.Sleep(15 * time.Second)
	// ---

	log.Printf("seteando la escena 2")
	sceneResp, err = wizClient.SetScene(8)
	if err != nil {
		log.Fatalf("error poniendo la escena: %s", err)
	}
	log.Print(sceneResp)

	// ---
	time.Sleep(15 * time.Second)
	// ---

	pilotResp, err := wizClient.getPilot()
	if err != nil {
		log.Fatalf("error pillando pilot: %s", err)
	}
	log.Print(pilotResp)

	// ---
	time.Sleep(5 * time.Second)
	// ---

	typeResp, err := wizClient.IsRgb()
	if err != nil {
		log.Fatalf("error averiguando tipo: %s", err)
	}
	log.Print(typeResp)

	// ---
	time.Sleep(5 * time.Second)
	// ---

	typeResp, err = wizClient.IsDw()
	if err != nil {
		log.Fatalf("error averiguando tipo: %s", err)
	}
	log.Print(typeResp)

	// ---
	time.Sleep(5 * time.Second)
	// ---

	typeResp, err = wizClient.IsTw()
	if err != nil {
		log.Fatalf("error averiguando tipo: %s", err)
	}
	log.Print(typeResp)

	// ---
	time.Sleep(5 * time.Second)
	// ---

	cResp, err := wizClient.SetColdWhite(255)
	if err != nil {
		log.Fatalf("error cambiando el color: %s", err)
	}
	log.Print(cResp)

	// ---
	time.Sleep(5 * time.Second)
	// ---

	cResp, err = wizClient.SetColdWhite(10)
	if err != nil {
		log.Fatalf("error cambiando el color: %s", err)
	}
	log.Print(cResp)

	// ---
	time.Sleep(5 * time.Second)
	// ---

	brightnessResp, err := wizClient.SetBrightness(100)
	if err != nil {
		log.Fatalf("error encendiendo: %s", err)
	}
	log.Print(brightnessResp)

	// ---
	time.Sleep(5 * time.Second)
	// ---

	rgbResp, err := wizClient.SetRgb(20, 30, 255)
	if err != nil {
		log.Fatalf("error cambiando el color: %s", err)
	}
	log.Print(rgbResp)

	// ---
	time.Sleep(5 * time.Second)
	// ---

	rgbResp, err = wizClient.SetRgb(10, 255, 30)
	if err != nil {
		log.Fatalf("error cambiando el color: %s", err)
	}
	log.Print(rgbResp)

	// ---
	time.Sleep(5 * time.Second)
	// ---

	cResp, err = wizClient.SetColdWhite(255)
	if err != nil {
		log.Fatalf("error cambiando el color: %s", err)
	}
	log.Print(cResp)

	// ---
	time.Sleep(5 * time.Second)
	// ---

	response, err = wizClient.TurnOff()
	if err != nil {
		log.Fatalf("error apagando: %s", err)
	}
	log.Print(response)

	//log.Print(err)
	//
	//response, err := wizClient.sendDatagrams([]byte(`{"id":1,"method":"setState","params":{"state":false}}`))
	//
	//log.Print(string(response))

	// ---
	//wizMessage := WizMessage{}
	//
	//err = json.Unmarshal(response, &wizMessage)
	//if err != nil {
	//	log.Fatalf("error haciendo unmarshal: %s", err)
	//}
	//
	//log.Print(wizMessage)

}
