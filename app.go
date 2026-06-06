package main

import (
	"context"
	"fmt"

	"github.com/bbuddha/marshall-linux/internal/ble"
	"github.com/bbuddha/marshall-linux/internal/device"
	"github.com/bbuddha/marshall-linux/internal/protocol"
)

type ScannedDevice struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Capabilities struct {
	HasANC          bool `json:"hasANC"`
	HasEQ           bool `json:"hasEQ"`
	HasBattery      bool `json:"hasBattery"`
	HasVolume       bool `json:"hasVolume"`
	HasRoomPlacement bool `json:"hasRoomPlacement"`
	HasPartyMode    bool `json:"hasPartyMode"`
	HasNightMode    bool `json:"hasNightMode"`
	HasLED          bool `json:"hasLED"`
}

type App struct {
	ctx context.Context
	dev *device.Device
}

type DeviceInfo struct {
	Model    string `json:"model"`
	Firmware string `json:"firmware"`
	Battery  int    `json:"battery"`
	ANC      string `json:"anc"`
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func bleDevicesToApp(in []ble.ScannedDevice) []ScannedDevice {
	out := make([]ScannedDevice, len(in))
	for i, d := range in {
		out[i] = ScannedDevice{Name: d.Name, Address: d.Address}
	}
	return out
}

func (a *App) GetCachedDevices() ([]ScannedDevice, error) {
	devices, err := ble.GetCachedDevices()
	if err != nil {
		return nil, err
	}
	return bleDevicesToApp(devices), nil
}

func (a *App) ScanDevices() ([]ScannedDevice, error) {
	devices, err := ble.ScanDevices()
	if err != nil {
		return nil, fmt.Errorf("scan failed: %w", err)
	}
	return bleDevicesToApp(devices), nil
}

func (a *App) GetCapabilities() *Capabilities {
	if a.dev == nil {
		return &Capabilities{}
	}
	set := map[string]bool{}
	for _, uuid := range a.dev.ListCharacteristics() {
		set[uuid] = true
	}
	return &Capabilities{
		HasANC:           set[protocol.CharANCConfiguration],
		HasEQ:            set[protocol.CharEqualizerSettings],
		HasBattery:       set[protocol.CharBatteryLevel],
		HasVolume:        set[protocol.CharVolume],
		HasRoomPlacement: set[protocol.CharRoomPlacement],
		HasPartyMode:     set[protocol.CharPartyMode],
		HasNightMode:     set[protocol.CharNightMode],
		HasLED:           set[protocol.CharLEDIntensity],
	}
}

func (a *App) Connect(target string) error {
	if a.dev != nil {
		a.dev.Close()
		a.dev = nil
	}
	dev, err := device.Connect(target)
	if err != nil {
		return fmt.Errorf("connexion échouée: %w", err)
	}
	a.dev = dev
	return nil
}

func (a *App) Disconnect() {
	if a.dev != nil {
		a.dev.Close()
		a.dev = nil
	}
}

func (a *App) IsConnected() bool {
	return a.dev != nil
}

func (a *App) GetInfo() (*DeviceInfo, error) {
	if a.dev == nil {
		return nil, fmt.Errorf("non connecté")
	}
	info := &DeviceInfo{}
	if m, err := a.dev.GetModelName(); err == nil {
		info.Model = m
	}
	if fw, err := a.dev.GetFirmwareVersion(); err == nil {
		info.Firmware = fw
	}
	if bat, err := a.dev.GetBatteryLevel(); err == nil {
		info.Battery = bat
	}
	if anc, err := a.dev.GetANC(); err == nil {
		info.ANC = anc.String()
	}
	return info, nil
}

func (a *App) SetANC(mode string) error {
	if a.dev == nil {
		return fmt.Errorf("non connecté")
	}
	var m protocol.ANCMode
	switch mode {
	case "off":
		m = protocol.ANCOff
	case "anc":
		m = protocol.ANCCancelling
	case "transparency":
		m = protocol.ANCTransparency
	default:
		return fmt.Errorf("mode inconnu: %s", mode)
	}
	return a.dev.SetANC(m)
}

func (a *App) GetANC() (string, error) {
	if a.dev == nil {
		return "", fmt.Errorf("non connecté")
	}
	mode, err := a.dev.GetANC()
	if err != nil {
		return "", err
	}
	return mode.String(), nil
}

func (a *App) SetEQ(preset string) error {
	if a.dev == nil {
		return fmt.Errorf("non connecté")
	}
	presets := map[string]protocol.EQPreset{
		"flat":              protocol.EQFlat,
		"custom":            protocol.EQCustom,
		"rock":              protocol.EQRock,
		"metal":             protocol.EQMetal,
		"pop":               protocol.EQPop,
		"hiphop":            protocol.EQHipHop,
		"electronic":        protocol.EQElectronic,
		"jazz":              protocol.EQJazz,
		"bass-boost":        protocol.EQBassBoost,
		"mid-boost":         protocol.EQMidBoost,
		"treble-boost":      protocol.EQTrebleBoost,
		"loud-push-workout": protocol.EQLoudPushWorkout,
	}
	p, ok := presets[preset]
	if !ok {
		return fmt.Errorf("preset inconnu: %s", preset)
	}
	return a.dev.SetEQPreset(p)
}

func (a *App) GetBattery() (int, error) {
	if a.dev == nil {
		return 0, fmt.Errorf("non connecté")
	}
	return a.dev.GetBatteryLevel()
}
