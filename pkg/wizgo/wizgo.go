package wizgo

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/exp/maps"
	"net"
	"slices"
	"strconv"
	"strings"

	wizgotypes "github.com/achetronic/wizgo/api/types"
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
		33:   "Diwali",
		35:   "Light alarm",
		1000: "Rhythm",
	}

	// TW - have Cool White and Warm White LEDs. Such devices support most static light modes + CCT control
	WizTwScenes = []int{6, 9, 10, 11, 12, 13, 14, 15, 16, 18, 29, 30, 31, 32}

	// DW - have only Dimmable white LEDs. Such devices support only dimming, Wake up, Bedtime and Night light modes.
	WizDwScenes = []int{9, 10, 13, 14, 29, 30, 31, 32}
)

const (
	// Info messages
	BrithnessRangeMessage   = "brightness must be between 10 and 100"
	LedRangeMessage         = "LED colors must be between 0 and 255"
	TemperatureRangeMessage = "temperature value must be between 2000 and 9000 (kelvin)"
	SpeedRangeMessage       = "speed must be between 10 and 200"
	RatioRangeMessage       = "ratio must be between 1 and 100"

	// Error messages
	ErrorReceivingResponseErrorMessage   = "error receiving response: %s"
	ErrorSendingDataErrorMessage         = "error sending data: %s"
	SystemConfigNotAvailableErrorMessage = "error getting system config: %s"
	DeviceTypeNotFoundErrorMessage       = "error figuring out device type: %s"
	SceneNotAvailableErrorMessage        = "scene not available: %s"
)

type WizClient struct {
	deviceConnection *net.UDPConn
}

// Thanks to project PyWizLights for some of the reverse engineering they already did previously than me
// Ref: https://github.com/sbidy/pywizlight

// TODO
func CreateWizClient(host string, port int) (wizClient WizClient, err error) {

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

// In general, send datagrams to the device
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

// sendMessage sends a WiZ message over UDP and returns the response already parsed
func (w *WizClient) sendMessage(message wizgotypes.WizMessage) (response wizgotypes.WizMessageResponse, err error) {

	// TODO
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return response, err
	}

	// TODO
	responseBytes, err := w.sendDatagrams(jsonMessage)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(responseBytes, &response)
	return response, err
}

// GetPilot return the current status for colors, temperature, scenes, etc
func (w *WizClient) GetPilot() (response wizgotypes.WizMessageResponse, err error) {

	wizMessage := wizgotypes.WizMessage{
		Id:     1,
		Method: "getPilot",
	}

	response, err = w.sendMessage(wizMessage)
	return response, err
}

// GetSystemConfig return current configuration related to the system
func (w *WizClient) GetSystemConfig() (response wizgotypes.WizMessageResponse, err error) {

	wizMessage := wizgotypes.WizMessage{
		Id:     1,
		Method: "getSystemConfig",
	}

	response, err = w.sendMessage(wizMessage)
	return response, err
}

// GetUserConfig return current configuration related to the user
func (w *WizClient) GetUserConfig() (response wizgotypes.WizMessageResponse, err error) {

	wizMessage := wizgotypes.WizMessage{
		Id:     1,
		Method: "getUserConfig",
	}

	response, err = w.sendMessage(wizMessage)
	return response, err
}

// GetModelConfig return current configuration related to the device model
func (w *WizClient) GetModelConfig() (response wizgotypes.WizMessageResponse, err error) {

	wizMessage := wizgotypes.WizMessage{
		Id:     1,
		Method: "getModelConfig",
	}

	response, err = w.sendMessage(wizMessage)
	return response, err
}

// GetDevInfo return current configuration related to the device
func (w *WizClient) GetDevInfo() (response wizgotypes.WizMessageResponse, err error) {

	wizMessage := wizgotypes.WizMessage{
		Id:     1,
		Method: "getDevInfo",
	}

	response, err = w.sendMessage(wizMessage)
	return response, err
}

// Registration is used to "register" with the bulb.
// This notifies the bulb if you want it to send you heartbeat sync packets
// After registering, you will receive on port 38900/udp of registered device, several messages like the following:
// {"method":"syncPilot","env":"pro","params":{"mac":"ABCABCABC","rssi":-71,"src":"udp","state":true,"sceneId":0,"temp":6500,"dimming":62 ··· }}
func (w *WizClient) Registration(phoneIp string, phoneMac string, register bool) (response wizgotypes.WizMessageResponse, err error) {

	wizMessage := wizgotypes.WizMessage{
		Id:     1,
		Method: "registration",
		Params: wizgotypes.WizMessageParams{
			"phoneIp":  phoneIp,
			"phoneMac": phoneMac,
			"register": register,
		},
	}

	response, err = w.sendMessage(wizMessage)
	return response, err
}

// Pulse generate a pulse of light to locate the bulb with ease
func (w *WizClient) Pulse() (response wizgotypes.WizMessageResponse, err error) {

	wizMessage := wizgotypes.WizMessage{
		Id:     1,
		Method: "pulse",
		Params: wizgotypes.WizMessageParams{
			"delta":    -100,
			"duration": 300,
		},
	}

	response, err = w.sendMessage(wizMessage)
	return response, err
}

// IsRgb return true when the device have RGB, cool white and warm white LEDs.
// These type of bulbs can support all the light modes provided by WiZ
func (w *WizClient) IsRgb() (bool, error) {

	configResp, err := w.GetSystemConfig()
	if err != nil {
		return false, errors.New(fmt.Sprintf(SystemConfigNotAvailableErrorMessage, err))
	}

	if !strings.Contains(configResp.Result.ModuleName, "RGB") {
		return false, nil
	}

	return true, nil
}

// IsTw return true when the device have cool white and warm white LEDs.
// These type of devices support most static light modes + CCT control
func (w *WizClient) IsTw() (bool, error) {

	configResp, err := w.GetSystemConfig()
	if err != nil {
		return false, errors.New(fmt.Sprintf(SystemConfigNotAvailableErrorMessage, err))
	}

	if !strings.Contains(configResp.Result.ModuleName, "TW") {
		return false, nil
	}

	return true, nil
}

// IsDw return true when the device have only dimmable white LEDs.
// These type of devices support only dimming and some light modes
func (w *WizClient) IsDw() (bool, error) {

	configResp, err := w.GetSystemConfig()
	if err != nil {
		return false, errors.New(fmt.Sprintf(SystemConfigNotAvailableErrorMessage, err))
	}

	if !strings.Contains(configResp.Result.ModuleName, "DW") {
		return false, nil
	}

	return true, nil
}

// IsSceneAvailable todo
func (w *WizClient) IsSceneAvailable(sceneId int) (available bool, err error) {

	isRgbDevice, err := w.IsRgb()
	if err != nil {
		return false, errors.New(fmt.Sprintf(DeviceTypeNotFoundErrorMessage, err))
	}

	if isRgbDevice && !slices.Contains(maps.Keys(WizScenes), sceneId) {
		return false, nil
	}

	isTwDevice, err := w.IsTw()
	if err != nil {
		return false, errors.New(fmt.Sprintf(DeviceTypeNotFoundErrorMessage, err))
	}

	if isTwDevice && !slices.Contains(WizTwScenes, sceneId) {
		return false, nil
	}

	isDwDevice, err := w.IsDw()
	if err != nil {
		return false, errors.New(fmt.Sprintf(DeviceTypeNotFoundErrorMessage, err))
	}

	if isDwDevice && !slices.Contains(WizDwScenes, sceneId) {
		return false, nil
	}

	return true, nil
}

// TurnOn turns on the device
func (w *WizClient) TurnOn() (response wizgotypes.WizMessageResponse, err error) {

	wizMessage := wizgotypes.WizMessage{
		Id:     1,
		Method: "setState",
		Params: wizgotypes.WizMessageParams{
			"state": true,
		},
	}

	response, err = w.sendMessage(wizMessage)
	return response, err
}

// TurnOff turns off the device
func (w *WizClient) TurnOff() (response wizgotypes.WizMessageResponse, err error) {

	wizMessage := wizgotypes.WizMessage{
		Id:     1,
		Method: "setState",
		Params: wizgotypes.WizMessageParams{
			"state": false,
		},
	}

	response, err = w.sendMessage(wizMessage)
	return response, err
}

// SetBrightness change the device's brightness (10-100)
func (w *WizClient) SetBrightness(brightness int) (response wizgotypes.WizMessageResponse, err error) {

	if brightness < 10 || brightness > 100 {
		return response, errors.New(BrithnessRangeMessage)
	}

	wizMessage := wizgotypes.WizMessage{
		Id:     1,
		Method: "setPilot",
		Params: wizgotypes.WizMessageParams{
			"dimming": brightness,
		},
	}

	response, err = w.sendMessage(wizMessage)
	return response, err
}

// SetRgb set the color for the device (3 x 0-255)
func (w *WizClient) SetRgb(r, g, b int) (response wizgotypes.WizMessageResponse, err error) {

	if r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
		return response, errors.New(LedRangeMessage)
	}

	wizMessage := wizgotypes.WizMessage{
		Id:     1,
		Method: "setPilot",
		Params: wizgotypes.WizMessageParams{
			"r": r,
			"g": g,
			"b": b,
		},
	}

	response, err = w.sendMessage(wizMessage)
	return response, err
}

// SetColdWhite set the level of light given by cold white LEDs (0-255)
func (w *WizClient) SetColdWhite(coldWhite int) (response wizgotypes.WizMessageResponse, err error) {

	if coldWhite < 0 || coldWhite > 255 {
		return response, errors.New(LedRangeMessage)
	}

	wizMessage := wizgotypes.WizMessage{
		Id:     1,
		Method: "setPilot",
		Params: wizgotypes.WizMessageParams{
			"c": coldWhite,
		},
	}

	response, err = w.sendMessage(wizMessage)
	return response, err
}

// SetWarmWhite set the level of light given by warm white LEDs (0-255)
func (w *WizClient) SetWarmWhite(warmWhite int) (response wizgotypes.WizMessageResponse, err error) {

	if warmWhite < 0 || warmWhite > 255 {
		return response, errors.New(LedRangeMessage)
	}

	wizMessage := wizgotypes.WizMessage{
		Id:     1,
		Method: "setPilot",
		Params: wizgotypes.WizMessageParams{
			"w": warmWhite,
		},
	}

	response, err = w.sendMessage(wizMessage)
	return response, err
}

// SetTemperature set color temperature in kelvin
func (w *WizClient) SetTemperature(temperature int) (response wizgotypes.WizMessageResponse, err error) {

	if temperature < 2000 || temperature > 9000 {
		return response, errors.New(TemperatureRangeMessage)
	}

	wizMessage := wizgotypes.WizMessage{
		Id:     1,
		Method: "setPilot",
		Params: wizgotypes.WizMessageParams{
			"temp": temperature,
		},
	}

	response, err = w.sendMessage(wizMessage)
	return response, err
}

// SetSpeed set changing speed between the colors in a scene
func (w *WizClient) SetSpeed(speed int) (response wizgotypes.WizMessageResponse, err error) {

	if speed < 10 || speed > 200 {
		return response, errors.New(SpeedRangeMessage)
	}

	wizMessage := wizgotypes.WizMessage{
		Id:     1,
		Method: "setPilot",
		Params: wizgotypes.WizMessageParams{
			"speed": speed,
		},
	}

	response, err = w.sendMessage(wizMessage)
	return response, err
}

// TODO
func (w *WizClient) SetRatio(ratio int) (response wizgotypes.WizMessageResponse, err error) {

	if ratio < 1 || ratio > 100 {
		return response, errors.New(RatioRangeMessage)
	}

	wizMessage := wizgotypes.WizMessage{
		Id:     1,
		Method: "setPilot",
		Params: wizgotypes.WizMessageParams{
			"ratio": ratio,
		},
	}

	response, err = w.sendMessage(wizMessage)
	return response, err
}

// SetScene set a scene by its ID. The availability of IDs depend on the type of bulb: RBG, TW, DW
func (w *WizClient) SetScene(sceneId int) (response wizgotypes.WizMessageResponse, err error) {

	isSceneAvailable, err := w.IsSceneAvailable(sceneId)
	if err != nil || !isSceneAvailable {
		return response, errors.New(fmt.Sprintf(SceneNotAvailableErrorMessage, err))
	}

	wizMessage := wizgotypes.WizMessage{
		Id:     1,
		Method: "setPilot",
		Params: wizgotypes.WizMessageParams{
			"sceneId": sceneId,
		},
	}

	response, err = w.sendMessage(wizMessage)
	return response, err
}

// SetRhythm set a rhythm by its ID.
// TODO: Implementation pending as it requires deeper reverse engineering to figure out what is needed
// There are three potential methods involved: setSchdPset, setSchd, setPilot
func (w *WizClient) SetRhythm(rhythmId int) (response wizgotypes.WizMessageResponse, err error) {
	return response, err
}
