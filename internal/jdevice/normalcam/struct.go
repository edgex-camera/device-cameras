package normalcam

import (
	"encoding/json"

	"github.com/edgex-camera/device-cameras/internal/jdevice"
	"github.com/edgex-camera/device-cameras/internal/lib/camera"
	"github.com/edgex-camera/device-sdk-go"
	"github.com/edgex-camera/device-sdk-go/pkg/camstartup"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
)

type NormalCamera struct {
	Name      string
	ChannelId string
	Camera    camera.RawCamera

	lc    logger.LoggingClient
	cmder camera.CameraCmder
	cc    camera.CameraConfig
}

func NewCamera(name, channelId string, lc logger.LoggingClient, cmder camera.CameraCmder, cc camera.CameraConfig) (jdevice.Camera, error) {
	ca := NormalCamera{Name: name}
	rawcam, err := NewRawCamera(name, channelId, lc, cmder, cc)
	if err != nil {
		return nil, err
	}
	ca.Camera = rawcam
	ca.ChannelId = channelId
	ca.lc = lc
	ca.cmder = cmder
	ca.cc = cc
	return &ca, nil
}

func NewRawCamera(name, channelId string, lc logger.LoggingClient, cmder camera.CameraCmder, cc camera.CameraConfig) (camera.RawCamera, error) {
	rawcam := camera.RawCamera{Lc: lc, Cmder: cmder, CameraConfig: cc}
	err := SetupRawCameraConfig(rawcam, name, channelId)
	if err != nil {
		return rawcam, err
	}
	return rawcam, nil
}

func SetupRawCameraConfig(c camera.RawCamera, name, channelId string) error {
	// 更新摄像头channel map
	basicConfName := name + ".camera"
	basicJson, _ := json.Marshal(map[string]bool{channelId: true})
	if string(basicJson) != device.DriverConfigs()[basicConfName] {
		err := camstartup.PutDriverConfig(basicConfName, basicJson)
		if err != nil {
			return err
		}
	}

	configName := name + ".camera." + channelId
	if _, ok := device.DriverConfigs()[configName]; !ok {
		return camstartup.PutDriverConfig(configName, c.GetConfigure())
	}
	return nil
}
