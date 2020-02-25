package onvif

import (
	"encoding/json"

	"github.com/edgexfoundry/device-sdk-go/pkg/jxstartup"
	jsonpatch "gopkg.in/evanphx/json-patch.v4"
)

func (c *OnvifCamera) MergeConfig(conf map[string]string) error {
	oldConf, err := json.Marshal(c.OnvifConfig)
	if err != nil {
		return err
	}
	configPatch := conf[c.Name+".onvif.config"]
	newConf, err := jsonpatch.MergePatch(oldConf, []byte(configPatch))
	if err != nil {
		return err
	}

	return json.Unmarshal(newConf, &c.OnvifConfig)
}

func (c *OnvifCamera) PutConfig(config []byte) error {
	configName := c.Name + ".onvif.config"
	return jxstartup.PutDriverConfig(configName, config)
}
