package onvif

import (
	"encoding/json"

	"github.com/edgexfoundry/device-sdk-go"
	"github.com/edgexfoundry/device-sdk-go/pkg/jxstartup"
)

func (c *OnvifCamera) MergeConfig(conf map[string]string) error {
	configPatch := conf[c.Name+".onvif.config"]
	return json.Unmarshal([]byte(configPatch), &c.OnvifConfig)
}

func (c *OnvifCamera) PutConfig(config []byte) error {
	configName := c.Name + ".onvif.config"
	return jxstartup.PutDriverConfig(configName, config)
}

func (c *OnvifCamera) GetConfigure() []byte {
	configName := c.Name + ".onvif.config"
	res, _ := device.DriverConfigs()[configName]
	return []byte(res)
}
