package device

import (
	"fmt"

	"github.com/bbuddha/marshall-linux/internal/ble"
	"github.com/bbuddha/marshall-linux/internal/protocol"
)

type Device struct {
	client *ble.Client
}

func Connect(target string) (*Device, error) {
	client, err := ble.Connect(target)
	if err != nil {
		return nil, err
	}
	return &Device{client: client}, nil
}

func (d *Device) SetANC(mode protocol.ANCMode) error {
	return d.client.Write(protocol.CharANCConfiguration, protocol.EncodeANC(mode))
}

func (d *Device) GetANC() (protocol.ANCMode, error) {
	data, err := d.client.Read(protocol.CharANCConfiguration)
	if err != nil {
		return protocol.ANCOff, err
	}
	return protocol.DecodeANC(data), nil
}

func (d *Device) SetEQPreset(preset protocol.EQPreset) error {
	if err := d.client.Write(protocol.CharEqualizerSettings, protocol.EncodeEQAssign(preset)); err != nil {
		return err
	}
	return d.client.Write(protocol.CharEqualizerSettings, protocol.EncodeEQActivate())
}

func (d *Device) GetBatteryLevel() (int, error) {
	data, err := d.client.Read(protocol.CharBatteryLevel)
	if err != nil {
		return 0, err
	}
	if len(data) == 0 {
		return 0, fmt.Errorf("empty response")
	}
	return int(data[0]), nil
}

func (d *Device) GetFirmwareVersion() (string, error) {
	data, err := d.client.Read(protocol.CharFirmwareRevision)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (d *Device) GetModelName() (string, error) {
	data, err := d.client.Read(protocol.CharModelName)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (d *Device) Close() {
	d.client.Close()
}

func (d *Device) ListCharacteristics() []string {
	return d.client.ListCharacteristics()
}
