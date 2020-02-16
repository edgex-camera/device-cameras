package onvif

import (
	"encoding/json"

	"github.com/edgexfoundry/device-sdk-go/pkg/jxstartup"
	jsonpatch "gopkg.in/evanphx/json-patch.v4"
)

func (c *OnvifCamera) MergeConfig(configPatch []byte) error {
	oldConf, err := json.Marshal(c.OnvifConfig)
	if err != nil {
		return err
	}
	newConf, err := jsonpatch.MergePatch(oldConf, configPatch)
	if err != nil {
		return err
	}

	return json.Unmarshal(newConf, &c.OnvifConfig)
}

func (c *OnvifCamera) PutConfig(config []byte) error {
	configName := c.Name + "/onvif/config"
	return jxstartup.PutDriverConfig(configName, config)
}
