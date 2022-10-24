package devicemanager

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/ntsd/cross-clipboard/pkg/device"
	"github.com/ntsd/cross-clipboard/pkg/xerror"
	"github.com/ntsd/go-utils/pkg/stringutil"
)

const devicesFileName = "devices.json"

func (dm *DeviceManager) Save() error {
	b, err := json.MarshalIndent(dm.Devices, "", "  ")
	if err != nil {
		return xerror.NewRuntimeError("can not marshal devices").Wrap(err)
	}

	deviceFilePath := stringutil.JoinURL(dm.config.ConfigDirPath, devicesFileName)

	err = os.WriteFile(deviceFilePath, b, 0644)
	if err != nil {
		return xerror.NewRuntimeError("can not write devices file").Wrap(err)
	}

	return nil
}

func (dm *DeviceManager) Load() error {
	deviceFilePath := stringutil.JoinURL(dm.config.ConfigDirPath, devicesFileName)

	f, err := os.Open(deviceFilePath)
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
