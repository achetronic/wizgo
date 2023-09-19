package main

import (
	"log"

	// Include the library in your code
	"github.com/achetronic/wizgo/pkg/wizgo"
)

func main() {

	// The first step is to start a client to interact with your devices
	wizClient, err := wizgo.CreateWizClient("192.168.2.107", 38899)

	// Devices can be turned on/off with single commands
	_, err = wizClient.TurnOn()
	if err != nil {
		log.Fatalf("error turning on the light: %s", err)
	}

	// Lights can be visually found sending a pulse
	_, err = wizClient.Pulse()
	if err != nil {
		log.Fatalf("error sending a pulse: %s", err)
	}

	// Scenes are supported too. Just set the ID you want
	_, err = wizClient.SetScene(35)
	if err != nil {
		log.Fatalf("error setting the scene: %s", err)
	}

	// May be you want to set the speed for a party?
	_, err = wizClient.SetSpeed(200)
	if err != nil {
		log.Fatalf("error setting the speed: %s", err)
	}

	// Oh, some romantic moment. Understood.
	_, err = wizClient.SetTemperature(2000)
	if err != nil {
		log.Fatalf("error setting temperature: %s", err)
	}

	// Studying is better with bright light
	_, err = wizClient.SetColdWhite(255)
	if err != nil {
		log.Fatalf("error setting cold white level: %s", err)
	}

	// But if it's too much, brightness can be set
	_, err = wizClient.SetBrightness(100)
	if err != nil {
		log.Fatalf("error setting the brightness level: %s", err)
	}

	// Feel like Grinch-y today?
	_, err = wizClient.SetRgb(0, 255, 0)
	if err != nil {
		log.Fatalf("error setting the color: %s", err)
	}

	// Let me sleep, uh?
	_, err = wizClient.TurnOff()
	if err != nil {
		log.Fatalf("error turning off the device: %s", err)
	}

	// There are a lot of more things implemented
	_, err = wizClient.GetPilot()
	_, err = wizClient.GetModelConfig()
	_, err = wizClient.GetDevInfo()
	_, err = wizClient.GetSystemConfig()
	_, err = wizClient.SetRatio(100)
	_, err = wizClient.SetWarmWhite(100)
	_, err = wizClient.IsDw()
	_, err = wizClient.IsTw()
	_, err = wizClient.IsRgb()
	_, err = wizClient.IsSceneAvailable(32)
	_, err = wizClient.Registration("192.168.2.173", "704F7C84524A", true)
}
