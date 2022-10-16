package devicemanager

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/ntsd/cross-clipboard/pkg/device"
)

const devicesFilePath = "devices.json"

type DeviceManager struct {
	Devices        map[string]*device.Device
	DevicesChannel chan map[string]*device.Device
}

func NewDeviceManager() *DeviceManager {
	return &DeviceManager{
		Devices:        make(map[string]*device.Device),
		DevicesChannel: make(chan map[string]*device.Device),
	}
}

func (dm *DeviceManager) AddDevice(device *device.Device) {
	dm.Devices[device.AddressInfo.ID.Pretty()] = device
	dm.DevicesChannel <- dm.Devices
}

func (dm *DeviceManager) RemoveDevice(device *device.Device) {
	// Flush and close ignore error
	device.Writer.Flush()
	device.Stream.Close()
	delete(dm.Devices, device.AddressInfo.ID.Pretty())
	dm.DevicesChannel <- dm.Devices
}

func (dm *DeviceManager) GetDevice(id string) *device.Device {
	return dm.Devices[id]
}

func (dm *DeviceManager) UpdateDevice(device *device.Device) {
	dm.Devices[device.AddressInfo.ID.Pretty()] = device
	dm.DevicesChannel <- dm.Devices
	dm.Save()
}

func (dm *DeviceManager) Save() error {
	b, err := json.MarshalIndent(dm.Devices, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(devicesFilePath, b, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (dm *DeviceManager) Load() error {
	f, err := os.Open(devicesFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	bytes, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	var devices map[string]*device.Device
	err = json.Unmarshal(bytes, &devices)
	if err != nil {
		return err
	}

	for _, dv := range devices {
		if dv.Status != device.StatusBlocked {
			dv.Status = device.StatusDisconnected
			err := dv.CreatePGPEncrypter()
			if err != nil {
				return err
			}
		}
	}

	dm.Devices = devices
	dm.DevicesChannel <- dm.Devices

	return nil
}
