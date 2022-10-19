package devicemanager

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/ntsd/cross-clipboard/pkg/device"
	"github.com/ntsd/cross-clipboard/pkg/xerror"
)

const devicesFilePath = "devices.json"

func (dm *DeviceManager) Save() error {
	b, err := json.MarshalIndent(dm.Devices, "", "  ")
	if err != nil {
		return xerror.NewRuntimeError("can not marshal devices").Wrap(err)
	}

	err = os.WriteFile(devicesFilePath, b, 0644)
	if err != nil {
		return xerror.NewRuntimeError("can not write devices file").Wrap(err)
	}

	return nil
}

func (dm *DeviceManager) Load() error {
	f, err := os.Open(devicesFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return xerror.NewRuntimeError("can not open devices file").Wrap(err)
	}
	bytes, err := io.ReadAll(f)
	if err != nil {
		return xerror.NewRuntimeError("can not read devices file").Wrap(err)
	}

	var devices map[string]*device.Device
	err = json.Unmarshal(bytes, &devices)
	if err != nil {
		return xerror.NewRuntimeError("can not unmarshal devices json").Wrap(err)
	}

	for _, dv := range devices {
		if dv.Status != device.StatusBlocked {
			dv.Status = device.StatusDisconnected
			err := dv.CreatePGPEncrypter()
			if err != nil {
				return xerror.NewRuntimeError("can not create pgp encrypter").Wrap(err)
			}
		}
	}

	dm.Devices = devices
	dm.DevicesChannel <- dm.Devices

	return nil
}
