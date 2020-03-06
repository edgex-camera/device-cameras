package jdevice

import (
	"encoding/json"
	"sync"

	"github.com/edgexfoundry/device-sdk-go"
	"github.com/edgexfoundry/device-sdk-go/pkg/jxstartup"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/lib/onvif"
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
		err := jxstartup.PutDriverConfig(configName, config)
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
