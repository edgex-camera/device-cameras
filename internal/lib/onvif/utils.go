package onvif

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"strconv"

	"github.com/yakovlevdmv/goonvif"
	"github.com/yakovlevdmv/goonvif/Media"
	"github.com/yakovlevdmv/goonvif/xsd/onvif"

	"github.com/edgex-camera/device-sdk-go"
	"github.com/edgex-camera/device-sdk-go/pkg/camstartup"
)

func getToken(config OnvifConfig) (onvif.ReferenceToken, error) {
	device, _ := goonvif.NewDevice(config.Address)
	device.Authenticate(config.Username, config.Password)
	req := Media.GetProfiles{}
	res, err := device.CallMethod(req)
	if err != nil {
		return "", err
	}
	body, _ := ioutil.ReadAll(res.Body)

	var response struct {
		Body struct {
			GetProfilesResponse struct {
				Profiles []onvif.Profile
			}
		}
	}
	xml.Unmarshal(body, &response)

	// 取第一个Profile使用
	profile := response.Body.GetProfilesResponse.Profiles[0]
	return profile.Token, nil
}

// 新建预置点配置，1点占用，2-255未占用
func InitPresetsConfig(name string) error {
	configName := name + ".onvif.presets"
	presets := make(map[int64]bool)
	presets[1] = true
	for i := 2; i < 256; i++ {
		presets[int64(i)] = false
	}
	config, _ := json.Marshal(presets)
	return camstartup.PutDriverConfig(configName, config)
}

func getPresets(name string) string {
	return device.DriverConfigs()[name+".onvif.presets"]
}

func setPreset(name string, number int64) {
	configName := name + ".onvif.presets"
	InitPresetsConfig(name)
	current := []byte(device.DriverConfigs()[configName])
	current_map := make(map[int64]bool)
	json.Unmarshal(current, &current_map)

	current_map[number] = true
	new_presets, _ := json.Marshal(current_map)
	camstartup.PutDriverConfig(configName, new_presets)
}

func numberToToken(number int64) onvif.ReferenceToken {
	return onvif.ReferenceToken(strconv.Itoa(int(number)))
}
