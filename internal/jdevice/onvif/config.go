package onvif

import (
	"encoding/json"

	jsonpatch "gopkg.in/evanphx/json-patch.v4"
)

type OnvifConfig struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *onvifCamera) MergeConfig(configPatch []byte) error {
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
