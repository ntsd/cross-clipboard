package devicemanager

import (
	"github.com/ntsd/cross-clipboard/pkg/device"
)

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
