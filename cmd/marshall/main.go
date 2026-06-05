package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/bbuddha/marshall-linux/internal/device"
	"github.com/bbuddha/marshall-linux/internal/protocol"
)

const usage = `marshall - Linux CLI for Marshall Bluetooth devices

Usage:
  marshall <address> <command> [args]

Commands:
  info                          Device model, firmware, battery, ANC mode
  anc <off|anc|transparency>   Set ANC mode
  anc                           Get current ANC mode
  eq <preset>                   Set EQ preset
  battery                       Get battery level
  scan                          List all discovered GATT characteristics

EQ presets:
  flat rock metal pop hiphop electronic jazz
  bass-boost mid-boost treble-boost loud-push-workout

Example:
  marshall 00:25:D1:41:DF:69 info
  marshall 00:25:D1:41:DF:69 anc anc
  marshall 00:25:D1:41:DF:69 eq rock
`

func main() {
	if len(os.Args) < 3 {
		fmt.Print(usage)
		os.Exit(1)
	}

	address := os.Args[1]
	cmd := os.Args[2]

	dev, err := device.Connect(address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer dev.Close()

	switch cmd {
	case "info":
		runInfo(dev)
	case "anc":
		if len(os.Args) < 4 {
			runGetANC(dev)
		} else {
			runSetANC(dev, os.Args[3])
		}
	case "eq":
		if len(os.Args) < 4 {
			fmt.Fprintln(os.Stderr, "eq requires a preset name")
			os.Exit(1)
		}
		runSetEQ(dev, os.Args[3])
	case "battery":
		runBattery(dev)
	case "scan":
		runScan(dev)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		os.Exit(1)
	}
}

func runInfo(dev *device.Device) {
	if model, err := dev.GetModelName(); err == nil {
		fmt.Printf("model:    %s\n", model)
	} else {
		fmt.Printf("model:    (unavailable: %v)\n", err)
	}
	if fw, err := dev.GetFirmwareVersion(); err == nil {
		fmt.Printf("firmware: %s\n", fw)
	} else {
		fmt.Printf("firmware: (unavailable: %v)\n", err)
	}
	if bat, err := dev.GetBatteryLevel(); err == nil {
		fmt.Printf("battery:  %d%%\n", bat)
	} else {
		fmt.Printf("battery:  (unavailable: %v)\n", err)
	}
	if anc, err := dev.GetANC(); err == nil {
		fmt.Printf("anc:      %s\n", anc)
	} else {
		fmt.Printf("anc:      (unavailable: %v)\n", err)
	}
}

func runGetANC(dev *device.Device) {
	mode, err := dev.GetANC()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(mode)
}

func runSetANC(dev *device.Device, arg string) {
	var mode protocol.ANCMode
	switch arg {
	case "off":
		mode = protocol.ANCOff
	case "anc", "cancelling":
		mode = protocol.ANCCancelling
	case "transparency", "ambient":
		mode = protocol.ANCTransparency
	default:
		fmt.Fprintf(os.Stderr, "unknown anc mode: %s (use off|anc|transparency)\n", arg)
		os.Exit(1)
	}
	if err := dev.SetANC(mode); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("anc set to: %s\n", mode)
}

var eqPresetNames = map[string]protocol.EQPreset{
	"custom":            protocol.EQCustom,
	"flat":              protocol.EQFlat,
	"rock":              protocol.EQRock,
	"metal":             protocol.EQMetal,
	"pop":               protocol.EQPop,
	"hiphop":            protocol.EQHipHop,
	"hip-hop":           protocol.EQHipHop,
	"electronic":        protocol.EQElectronic,
	"jazz":              protocol.EQJazz,
	"bass-boost":        protocol.EQBassBoost,
	"mid-boost":         protocol.EQMidBoost,
	"treble-boost":      protocol.EQTrebleBoost,
	"loud-push-workout": protocol.EQLoudPushWorkout,
}

func runSetEQ(dev *device.Device, arg string) {
	preset, ok := eqPresetNames[arg]
	if !ok {
		n, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unknown eq preset: %s\n", arg)
			os.Exit(1)
		}
		preset = protocol.EQPreset(n)
	}
	if err := dev.SetEQPreset(preset); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("eq set to: %s\n", preset)
}

func runBattery(dev *device.Device) {
	bat, err := dev.GetBatteryLevel()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%d%%\n", bat)
}

func runScan(dev *device.Device) {
	uuids := dev.ListCharacteristics()
	sort.Strings(uuids)
	fmt.Printf("found %d characteristics:\n", len(uuids))
	for _, u := range uuids {
		fmt.Printf("  %s\n", u)
	}
}
