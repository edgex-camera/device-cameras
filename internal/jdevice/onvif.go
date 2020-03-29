package jdevice

import (
	"encoding/json"
	"sync"

	"github.com/edgex-camera/device-cameras/internal/lib/onvif"
	"github.com/edgex-camera/device-sdk-go"
	"github.com/edgex-camera/device-sdk-go/pkg/camstartup"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
)

func NewOnvif(name string, lc logger.LoggingClient, config onvif.OnvifConfig) (Control, error) {
	oc := &onvif.OnvifCamera{
		Name:        name,
		Lc:          lc,
		OnvifConfig: config,
		Mutex:       &sync.Mutex{},
	}
	err := SetupOnvifConfig(oc.OnvifConfig, name)
	if err != nil {
		return oc, err
	}
	return oc, nil
}

func SetupOnvifConfig(onvifConf onvif.OnvifConfig, name string) error {
	configName := name + ".onvif.config"
	if _, ok := device.DriverConfigs()[configName]; !ok {
		config, _ := json.Marshal(onvifConf)
		err := camstartup.PutDriverConfig(configName, config)
		if err != nil {
			return err
		}
	}
	presetsName := name + ".onvif.presets"
	if _, ok := device.DriverConfigs()[presetsName]; !ok {
		err := onvif.InitPresetsConfig(name)
		if err != nil {
			return err
		}
	}
	return nil
}
