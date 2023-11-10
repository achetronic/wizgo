# Wizgo

<img src="https://raw.githubusercontent.com/achetronic/wizgo/master/docs/img/logo.png" alt="Wizgo Logo (Main) logo." width="150">

![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/achetronic/wizgo)
![GitHub](https://img.shields.io/github/license/achetronic/wizgo)
[![Go Reference](https://pkg.go.dev/badge/github.com/achetronic/wizgo.svg)](https://pkg.go.dev/github.com/achetronic/wizgo)
[![Go Report Card](https://goreportcard.com/badge/github.com/achetronic/wizgo)](https://goreportcard.com/report/github.com/achetronic/wizgo)

![YouTube Channel Subscribers](https://img.shields.io/youtube/channel/subscribers/UCeSb3yfsPNNVr13YsYNvCAw?label=achetronic&link=http%3A%2F%2Fyoutube.com%2Fachetronic)
![X (formerly Twitter) Follow](https://img.shields.io/twitter/follow/achetronic?style=flat&logo=twitter&link=https%3A%2F%2Ftwitter.com%2Fachetronic)

A golang library to control your WiZ lights

> At this moment, the library is not covering the whole API as it is discovered by reverse engineering 
> If you want to cover more things, consider [contributing](#how-to-contribute)

## Motivation

Domotic devices are cool, but they will be cooler when [Matter](https://csa-iot.org/all-solutions/matter/) 
protocol lands, finally, into the industry.

As a fact, before Matter, only some protocols such as Zigbee had some standardization mercy.

But there are a lot of devices out there connected by Wi-Fi, and they are completely non-standard. 
Manufacturers decided to design them without thinking about potential integration between different stuff.

The big question is: what are you doing with your legacy WiZ devices that are not going to be updated by the company,
once when Matter has landed? Are you thinking about throwing them to the trash?

With this library you can make your own automations with them. 
And now it's as easy as you should forgive about throwing... what?

## Methods

Just read these wonderful lines:

```go
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


```

## How to contribute

Of course, we are open to external collaborations for this project. For doing it you must:

* Open an issue to discuss what is needed and the reason
* Fork the repository
* Make your changes to the code 
* Open a PR. The code will be reviewed and tested (always)

> We are developers and hate bad code. For that reason we ask you the highest quality on each line of code to improve
> this project on each iteration.

## License

Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

## Special mention

This project was done using IDEs from JetBrains. They helped us to develop faster, so we recommend them a lot! 🤓

<img src="https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.png" alt="JetBrains Logo (Main) logo." width="150">
