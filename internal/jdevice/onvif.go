package jdevice

import (
	"encoding/json"

	"github.com/edgexfoundry/device-sdk-go"
	"github.com/edgexfoundry/device-sdk-go/pkg/jxstartup"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/lib/onvif"
)

func NewOnvif(name string, lc logger.LoggingClient, config onvif.OnvifConfig) (Onvif, error) {
	oc := &onvif.OnvifCamera{
		Name:        name,
		Lc:          lc,
		OnvifConfig: config,
	}
	err := SetupOnvifConfig(oc, name)
	if err != nil {
		return oc, err
	}
	return oc, nil
}

func SetupOnvifConfig(onvif Onvif, name string) error {
	configName := name + ".onvif.config"
	if _, ok := device.DriverConfigs()[configName]; !ok {
		config, _ := json.Marshal(onvif)
		return jxstartup.PutDriverConfig(configName, config)
	}
	return nil
}
