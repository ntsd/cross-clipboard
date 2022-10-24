package devicemanager

import (
	"github.com/ntsd/cross-clipboard/pkg/config"
	"github.com/ntsd/cross-clipboard/pkg/device"
)

type DeviceManager struct {
	Devices        map[string]*device.Device
	DevicesChannel chan map[string]*device.Device

	config *config.Config
}

func NewDeviceManager(cfg *config.Config) *DeviceManager {
	return &DeviceManager{
		Devices:        make(map[string]*device.Device),
		DevicesChannel: make(chan map[string]*device.Device),
		config:         cfg,
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
